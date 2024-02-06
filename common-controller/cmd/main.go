/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	logger "github.com/sirupsen/logrus"
	commoncontroller "github.com/wso2/apk/common-controller/commoncontroller"
	config "github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/database"
	"github.com/wso2/apk/common-controller/internal/server"
	web "github.com/wso2/apk/common-controller/internal/web"
)

func main() {
	conf := config.ReadConfigs()
	logger.Info("Starting the Web server")
	go web.StartWebServer()
	go server.StartInternalServer()
	if conf.CommonController.Database.Enabled {
		logger.Info("Starting the Database connection")
		go startDB()
	}
	logger.Info("Starting the Common Controller")
	commoncontroller.InitCommonControllerServer(conf)

}

func startDB() {
	database.ConnectToDB()
	database.GetApplicationByUUID()
	defer database.CloseDBConn()
}
