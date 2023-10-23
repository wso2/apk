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

package utils

import (
	"github.com/wso2/apk/common-controller/internal/config"
)

// GetEnvironment takes the environment of the API. If the value is empty,
// it will return the default environment that is set in the config of the common controller.
func GetEnvironment(environment string) string {
	if environment != "" {
		return environment
	}
	return config.ReadConfigs().CommonController.Environment
}
