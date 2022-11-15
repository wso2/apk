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

package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/adapter"
)

// invokes the code from the /internal and /pkg directories and nothing else.
func main() {
	logger.Info("Starting the Adapter")
	conf := config.ReadConfigs()
	adapter.Run(conf)
}
