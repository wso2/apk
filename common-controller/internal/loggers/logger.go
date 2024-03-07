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
	pkgAPKOperator             = "github.com/wso2/apk/adapter/internal/operator"
	pkgEnforcerXdsCallbacks    = "github.com/wso2/apk/adapter/internal/discovery/xds/enforcercallbacks"
	pkgRouterXdsCallbacks      = "github.com/wso2/apk/adapter/internal/discovery/xds/routercallbacks"
	pkgGrpcClient              = "github.com/wso2/apk/adapter/internal/grpc-client"
	pkgXds                     = "github.com/wso2/apk/adapter/internal/discovery/xds"
	pkgOasParser               = "github.com/wso2/apk/adapter/internal/oasparser"
	pkgInterceptor             = "github.com/wso2/apk/adapter/internal/interceptor"
	pkgSvcDiscovery            = "github.com/wso2/apk/adapter/internal/svcdiscovery"
	pkgNotifier                = "github.com/wso2/apk/adapter/internal/notifier"
	pkgAPI                     = "github.com/wso2/apk/adapter/internal/api"
	pkgAPK                     = "github.com/wso2/apk/adapter/internal/adapter"
	pkgRateLimiterXdsCallbacks = "github.com/wso2/apk/adapter/internal/discovery/xds/ratelimitercallbacks"
	pkgDatabase                = "github.com/wso2/apk/adapter/internal/database"
)

// logger package references
var (
	LoggerAPKOperator             logging.Log
	LoggerEnforcerXdsCallbacks    logging.Log
	LoggerRouterXdsCallbacks      logging.Log
	LoggerXds                     logging.Log
	LoggerGRPCClient              logging.Log
	LoggerOasparser               logging.Log
	LoggerInterceptor             logging.Log
	LoggerSvcDiscovery            logging.Log
	LoggerNotifier                logging.Log
	LoggerAPI                     logging.Log
	LoggerAPK                     logging.Log
	LoggerRateLimiterXdsCallbacks logging.Log
	LoggerDatabase                logging.Log
)

func init() {
	UpdateLoggers()
}

// UpdateLoggers initializes the logger package references
func UpdateLoggers() {
	LoggerRouterXdsCallbacks = logging.InitPackageLogger(pkgRouterXdsCallbacks)
	LoggerEnforcerXdsCallbacks = logging.InitPackageLogger(pkgEnforcerXdsCallbacks)
	LoggerXds = logging.InitPackageLogger(pkgXds)
	LoggerGRPCClient = logging.InitPackageLogger(pkgGrpcClient)
	LoggerAPKOperator = logging.InitPackageLogger(pkgAPKOperator)
	LoggerOasparser = logging.InitPackageLogger(pkgOasParser)
	LoggerInterceptor = logging.InitPackageLogger(pkgInterceptor)
	LoggerSvcDiscovery = logging.InitPackageLogger(pkgSvcDiscovery)
	LoggerRateLimiterXdsCallbacks = logging.InitPackageLogger(pkgRateLimiterXdsCallbacks)
	LoggerNotifier = logging.InitPackageLogger(pkgNotifier)
	LoggerAPI = logging.InitPackageLogger(pkgAPI)
	LoggerAPK = logging.InitPackageLogger(pkgAPK)
	LoggerDatabase = logging.InitPackageLogger(pkgDatabase)
	logrus.Info("Updated loggers")
}
