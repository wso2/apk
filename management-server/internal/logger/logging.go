/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org)
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

// Package logger contains the package references for log messages
// If a new package is introduced, the corresponding logger reference is need to be created as well.
package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/wso2/product-microgateway/adapter/pkg/logging"
)

/* loggers should be initiated only for the main packages
 ********** Don't initiate loggers for sub packages ****************

When you add a new logger instance add the related package name as a constant
*/

// package name constants
const (
	pkgServer    = "github.com/wso2/apk/APKManagementServer"
	pkgXds       = "github.com/wso2/apk/APKManagementServer/xds"
	pkgXdsServer = "github.com/wso2/apk/APKManagementServer/xds/server"
)

// logger package references
var (
	LoggerServer    logging.Log
	LoggerXds       logging.Log
	LoggerXdsServer logging.Log
)

func init() {
	UpdateLoggers()
}

// UpdateLoggers initializes the logger package references
func UpdateLoggers() {
	LoggerServer = logging.InitPackageLogger(pkgServer)
	LoggerXds = logging.InitPackageLogger(pkgXds)
	LoggerXdsServer = logging.InitPackageLogger(pkgXdsServer)
	logrus.Info("Updated loggers")
}
