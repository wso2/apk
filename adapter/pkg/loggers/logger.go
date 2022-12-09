/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	pkgAuth      = "github.com/wso2/apk/adapter/pkg/auth"
	pkgSync      = "github.com/wso2/apk/adapter/pkg/synchronizer"
	pkgTLSUtils  = "github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	pkgHealth    = "github.com/wso2/apk/adapter/pkg/health"
	pkgSoapUtils = "github.com/wso2/apk/adapter/pkg/utils/soaputils"
)

// logger package references
var (
	LoggerAuth      logging.Log
	LoggerSync      logging.Log
	LoggerTLSUtils  logging.Log
	LoggerHealth    logging.Log
	LoggerSoapUtils logging.Log
)

func init() {
	UpdateLoggers()
}

// UpdateLoggers initializes the logger package references
func UpdateLoggers() {
	LoggerAuth = logging.InitPackageLogger(pkgAuth)
	LoggerSync = logging.InitPackageLogger(pkgSync)
	LoggerTLSUtils = logging.InitPackageLogger(pkgTLSUtils)
	LoggerHealth = logging.InitPackageLogger(pkgHealth)
	LoggerSoapUtils = logging.InitPackageLogger(pkgSoapUtils)
	logrus.Info("Updated loggers")
}
