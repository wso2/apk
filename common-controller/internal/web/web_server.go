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

package web

import (
	"github.com/gin-gonic/gin"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	config "github.com/wso2/apk/common-controller/internal/config"
	"fmt"
)

// StartWebServer starts the web server
func StartWebServer() {
	loggers.LoggerAPI.Info("Starting web server")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/notify", NotifyHandler)
	conf := config.ReadConfigs()
	certPath := conf.CommonController.Keystore.CertPath
	keyPath := conf.CommonController.Keystore.KeyPath
	port := conf.CommonController.WebServer.Port
	router.RunTLS(fmt.Sprintf(":%d", port), certPath, keyPath)
}
