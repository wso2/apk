/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

// Package loggers contains the package references for log messages
// If a new package is introduced, the corresponding logger reference is need to be created as well.
package loggers

import (
	"github.com/sirupsen/logrus"
	"github.com/wso2/apk/adapter/pkg/logging"
)

/* loggers should be initiated only for the main packages
 ********** Don't initiate loggers for sub packages ****************

When you add a new logger instance add the related package name as a constant
*/

// package name constants
const (
	apkOperator             = "github.com/wso2/apk/adapter/internal/operator"
	pkgEnforcerXdsCallbacks = "github.com/wso2/apk/adapter/internal/discovery/xds/enforcercallbacks"
	pkgRouterXdsCallbacks   = "github.com/wso2/apk/adapter/internal/discovery/xds/routercallbacks"
	pkgXds                  = "github.com/wso2/product-microgateway/adapter/internal/discovery/xds"
	gRPCClient              = "github.com/wso2/apk/adapter/internal/grpc-client"
)

// logger package references
var (
	LoggerAPKOperator          logging.Log
	LoggerEnforcerXdsCallbacks logging.Log
	LoggerRouterXdsCallbacks   logging.Log
	LoggerXds                  logging.Log
	LoggerGRPCClient           logging.Log
)

func init() {
	UpdateLoggers()
}

// UpdateLoggers initializes the logger package references
func UpdateLoggers() {
	LoggerRouterXdsCallbacks = logging.InitPackageLogger(pkgRouterXdsCallbacks)
	LoggerEnforcerXdsCallbacks = logging.InitPackageLogger(pkgEnforcerXdsCallbacks)
	LoggerXds = logging.InitPackageLogger(pkgXds)
	LoggerAPKOperator = logging.InitPackageLogger(apkOperator)
	LoggerGRPCClient = logging.InitPackageLogger(gRPCClient)
	logrus.Info("Updated loggers")
}
