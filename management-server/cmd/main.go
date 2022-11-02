/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org).
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

	"github.com/wso2/apk/APKManagementServer/internal/logger"
	"github.com/wso2/apk/APKManagementServer/internal/xds"
)

func main() {
	logger.LoggerServer.Info("Hello, world.")
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go xds.InitAPKMgtServer()
	// todo(amaliMatharaarachchi) watch data updates and update snapshot accordingly.
	go xds.FeedData()

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
