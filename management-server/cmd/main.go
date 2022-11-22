/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org).
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package main

import (
	"os"
	"os/signal"

	"github.com/wso2/apk/management-server/internal/database"
	server "github.com/wso2/apk/management-server/internal/grpc-server"
	"github.com/wso2/apk/management-server/internal/logger"
	"github.com/wso2/apk/management-server/internal/xds"
)

func main() {
	logger.LoggerServer.Info("Starting Management server ...")
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	// connect to the postgres database
	database.ConnectToDB()
	defer database.CloseDBConn()
	go xds.InitAPKMgtServer()
	// todo(amaliMatharaarachchi) watch data updates and update snapshot accordingly.

	// temp data
	// config := config.ReadConfigs()
	// var arr = []*internal_types.ApplicationEvent{
	// 	{
	// 		Label:         config.ManagementServer.NodeLabels[0],
	// 		UUID:          "b9850225-c7db-444d-87fd-4feeb3c6b3cc",
	// 		IsRemoveEvent: false,
	// 	},
	// }
	// go xds.AddMultipleApplications(arr)
	go server.StartGRPCServer()

OUTER:
	for {
		select {
		case s := <-sig:
			switch s {
			case os.Interrupt:
				break OUTER
			}
		}
	}
}
