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
	"sync"
	"time"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"google.golang.org/grpc"
	grpcStatus "google.golang.org/grpc/status"
) 


type Pool struct {
	// GRPC server address
	serverAddress string;
	// Maximum capacity of the pool
	maxCapacity int;
	// Initial active connection. Set zero if you need on demand connections
	desiredCapacity int;
	dialOptions []grpc.DialOption;
	availableConnections *[]grpc.ClientConn;
	usedConnections *[]grpc.ClientConn;
	lock *sync.Mutex;
	retryPolicy RetryPolicy;
}

type RetryPolicy struct {
	// Maximum number of time a failed grpc call will be retried. Set negative value to try indefinitely.
	MaxAttempts int;
	// Time delay between retries. (In milli seconds)
	BackOffInMilliSeconds int;
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
		*connectionPool.availableConnections = append(*connectionPool.availableConnections, *conn);
	}
	return &connectionPool, nil;
}

// Initialize a pool with default configuration
func InitWithConfig()  (*Pool, error) {
	conf, _ := config.ReadConfigs()
	address := conf.Adapter.GRPCClient.ManagementServerAddress;
	return Init(address, conf.Adapter.GRPCClient.MaxCapacity, conf.Adapter.GRPCClient.DesiredCapacity, 
		[]grpc.DialOption{
			// TODO use tls credentials.
			grpc.WithInsecure(),
			grpc.WithBlock()}, 
		RetryPolicy{
			MaxAttempts : conf.Adapter.GRPCClient.MaxAttempts,
			BackOffInMilliSeconds : conf.Adapter.GRPCClient.BackOffInMilliSeconds,
		})
}

func (connectionPool *Pool) GetConnection() (*grpc.ClientConn, error){
	connectionPool.lock.Lock();
	defer connectionPool.lock.Unlock();
	availableConnectionLength := len(*connectionPool.availableConnections);
	if (availableConnectionLength > 0) {
		availableConnection := &(*connectionPool.availableConnections)[0];
		*connectionPool.availableConnections = (*connectionPool.availableConnections)[1:]
		*connectionPool.usedConnections = append(*connectionPool.usedConnections, *availableConnection);
		return availableConnection, nil;
	} else {
		totalConnectionLength := availableConnectionLength + len(*connectionPool.usedConnections);
		if (totalConnectionLength < connectionPool.maxCapacity) {
			connection, err := createGRPCConnection(*connectionPool);
			if (err == nil) {
				*connectionPool.usedConnections = append(*connectionPool.usedConnections, *connection);
				return connection, nil;
			}
			return nil, err;
		} else {
			return nil, errors.New("maximum connection reached in the pool")
		}
	}
}

func createGRPCConnection(connectionPool Pool) (*grpc.ClientConn, error) {
	return grpc.Dial(
		connectionPool.serverAddress, 
		connectionPool.dialOptions...	
	)
}

// Close a specific connection.
func (connectionPool *Pool) Close(connection *grpc.ClientConn) error{
	connectionPool.lock.Lock();
	defer connectionPool.lock.Unlock();
	var indexOfConnectionInPool int = -1;
	var connections *[]grpc.ClientConn;
	// Try to find the connection in the usedConnection slice.
	for k, v := range *connectionPool.usedConnections {
		if connection == &v {
			indexOfConnectionInPool = k;
			connections = connectionPool.usedConnections;
			break;
		}
    }
	// If usedConnection slice does not contain connection find it in available connection.
	if (indexOfConnectionInPool == -1) {
		for k, v := range *connectionPool.availableConnections {
			if connection == &v {
				indexOfConnectionInPool = k;
				connections = connectionPool.availableConnections;
				break;
			}
		}
	}
	if (indexOfConnectionInPool != -1) {
		len := len(*connections);
		(*connections)[indexOfConnectionInPool] = (*connections)[len-1]
		*connections = (*connections)[:len-1];
	}
	return connection.Close();
}

// Close all connection in the pool.
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
			errStatus, _ := grpcStatus.FromError(err)
			logger.LoggerGRPCClient.Errorf("gRPC call failed. errorCode: %s errorMessage: %s", errStatus.Code().String(), errStatus.Message());
			if (connectionPool.retryPolicy.MaxAttempts < 0) {
				// If max attempts has a negative value, retry indefinitely by setting retry less than max attempts.
				retries = connectionPool.retryPolicy.MaxAttempts - 1;
			} else {
				retries++;
			}
			if (retries <= connectionPool.retryPolicy.MaxAttempts) {
				// Retry grpc call after BackOffInMilliSeconds
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
