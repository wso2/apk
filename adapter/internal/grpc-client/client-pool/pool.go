/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/wso2/apk/adapter/config"
	"google.golang.org/grpc"
) 


type Pool struct {
	serverAddress string;
	maxCapacity int;
	desiredCapacity int;
	dialOptions []grpc.DialOption;
	availableConnections *[]grpc.ClientConn;
	usedConnections *[]grpc.ClientConn;
	lock *sync.Mutex;
	retryPolicy RetryPolicy;
}

type RetryPolicy struct {
	MaxAttempts int;
	BackOffInMilliSeconds int;
	RetryableStatuses []string;
}


func Init(serverAddress string, maxCapacity int, desiredCapacity int, dialOptions []grpc.DialOption, retryPolicy RetryPolicy) (*Pool, error) {
	connectionPool := Pool{lock: &sync.Mutex{}};
	connectionPool.lock.Lock();
	defer connectionPool.lock.Unlock();
	connectionPool.maxCapacity = maxCapacity;
	connectionPool.desiredCapacity =  desiredCapacity;
	connectionPool.serverAddress = serverAddress;
	connectionPool.dialOptions = dialOptions;
	connectionPool.availableConnections = &[]grpc.ClientConn{};
	connectionPool.usedConnections = &[]grpc.ClientConn{};
	connectionPool.retryPolicy = retryPolicy;
	for i := 0; i < desiredCapacity; i++ {
		conn, err := createGRPCConnection(connectionPool);
		if (err != nil) {
			return nil,err;
		}
		
		availableConnections := connectionPool.availableConnections;
		*availableConnections = append(*availableConnections, *conn);
	}
	return &connectionPool, nil;
}

func InitWithConfig()  (*Pool, error) {
	conf, _ := config.ReadConfigs()
	address := conf.Adapter.GRPCClient.ManagementServerAddress + ":" + strconv.Itoa(conf.Adapter.GRPCClient.ManagementServerGRPCPort);
	return Init(address, conf.Adapter.GRPCClient.MaxCapacity, conf.Adapter.GRPCClient.DesiredCapacity, 
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock()}, 
		RetryPolicy{
			MaxAttempts : conf.Adapter.GRPCClient.MaxAttempts,
			BackOffInMilliSeconds : conf.Adapter.GRPCClient.BackOffInMilliSeconds,
			RetryableStatuses : []string{},
		})
}

func (connectionPool *Pool) GetConnection() (*grpc.ClientConn, error){
	connectionPool.lock.Lock();
	defer connectionPool.lock.Unlock();
	availableConnections := *connectionPool.availableConnections;
	availableConnectionLength := len(availableConnections);
	if (availableConnectionLength > 0) {
		availableConnection := &availableConnections[0];
		*connectionPool.availableConnections = (*connectionPool.availableConnections)[1:]
		*connectionPool.usedConnections = append(*connectionPool.usedConnections, *availableConnection);
		return availableConnection, nil;
	} else {
		totalConnectionLength := availableConnectionLength + len(*connectionPool.usedConnections);
		if (totalConnectionLength < connectionPool.maxCapacity) {
			return createGRPCConnection(*connectionPool);
		} else {
			return nil, errors.New("Maximum connection reached in the pool.")
		}
	}
}

func createGRPCConnection(connectionPool Pool) (*grpc.ClientConn, error) {
	return grpc.Dial(
		connectionPool.serverAddress, 
		connectionPool.dialOptions...	
	)
}

func (connectionPool *Pool) Close(connection *grpc.ClientConn) error{
	connectionPool.lock.Lock();
	defer connectionPool.lock.Unlock();
	var index int = -1;
	var connections *[]grpc.ClientConn;
	for k, v := range *connectionPool.usedConnections {
       if connection == &v {
           index = k;
		   connections = connectionPool.usedConnections;
		   break;
       }
    }
	if (index == -1) {
		for k, v := range *connectionPool.availableConnections {
			if connection == &v {
				index = k;
				connections = connectionPool.availableConnections;
				break;
			}
		}
	}
	if (index != -1) {
		len := len(*connections);
		(*connections)[index] = (*connections)[len-1]
		*connections = (*connections)[:len-1];
	}
	return connection.Close();
}

func (connectionPool *Pool) CloseAll() {
	for _, v := range *connectionPool.usedConnections {
       v.Close();
    }
	for _, v := range *connectionPool.availableConnections {
       v.Close();
    }
	*connectionPool.usedConnections = []grpc.ClientConn{};
	*connectionPool.availableConnections = []grpc.ClientConn{};
}

func (connectionPool *Pool) ExecuteGRPCCall(connection *grpc.ClientConn, call func() (interface{}, error)) (interface{}, error) {
	retries := 0;
	response, err := call();
	for {
		
		if (err != nil) {
			if (connectionPool.retryPolicy.MaxAttempts < 0) {
				retries = connectionPool.retryPolicy.MaxAttempts - 1;
			} else {
				retries++;
			}
			if (retries <= connectionPool.retryPolicy.MaxAttempts) {
				log.Print("Error occured while calling grpc server", err);
				time.Sleep(time.Duration(connectionPool.retryPolicy.BackOffInMilliSeconds) * time.Millisecond)
				response, err = call();
			} else {
				return response, err;
			}
		} else {
			return response, nil;
		}
	}
}
