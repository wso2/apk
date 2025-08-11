/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package logging contains the package references for log messages
// If a new package is introduced, the corresponding logger reference is need to be created as well.
package logging

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
	pkgMain              = "github.com/wso2/apk/config-deployer-service-go/cmd"
	pkgArtifactGenerator = "github.com/wso2/apk/config-deployer-service-go/internal/artifactgenerator"
	pkgConfig            = "github.com/wso2/apk/config-deployer-service-go/internal/config"
	pkgUtil              = "github.com/wso2/apk/config-deployer-service-go/internal/util"
)

// logger package references
var (
	LoggerMain                     logging.Log
	LoggerAPIClient                logging.Log
	LoggerArtifactGeneratorService logging.Log
	LoggerConfigGeneratorClient    logging.Log
	LoggerConfig                   logging.Log
	LoggerDefinitionParser         logging.Log
	LoggerGraphQLParser            logging.Log
	LoggerOasParser                logging.Log
	LoggerProtoParser              logging.Log
	LoggerRuntimeAPICommon         logging.Log
)

func init() {
	UpdateLoggers()
}

// UpdateLoggers initializes the logger package references
func UpdateLoggers() {
	LoggerMain = logging.InitPackageLogger(pkgMain)
	LoggerAPIClient = logging.InitPackageLogger(pkgArtifactGenerator)
	LoggerArtifactGeneratorService = logging.InitPackageLogger(pkgArtifactGenerator)
	LoggerConfigGeneratorClient = logging.InitPackageLogger(pkgArtifactGenerator)
	LoggerConfig = logging.InitPackageLogger(pkgConfig)
	LoggerDefinitionParser = logging.InitPackageLogger(pkgUtil)
	LoggerGraphQLParser = logging.InitPackageLogger(pkgUtil)
	LoggerOasParser = logging.InitPackageLogger(pkgUtil)
	LoggerProtoParser = logging.InitPackageLogger(pkgUtil)
	LoggerRuntimeAPICommon = logging.InitPackageLogger(pkgUtil)
	logrus.Info("Updated loggers")
}
