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

package logging

import (
	logging "github.com/wso2/apk/adapter/pkg/logging"
)

func init() {
	Mapper = CombineMapper(internalMapper)
}

// Error Log Internal Adapter Constants
const (
	error1101 = 1101
	error1102 = 1102
	error1103 = 1103
	error1104 = 1104
)

var internalMapper = map[int]logging.ErrorDetails{
	error1101: {
		ErrorCode: error1101,
		Message:   "Failed to start XDS GRPS server",
		Severity:  "BLOCKER",
	},
	error1102: {
		ErrorCode: error1102,
		Message:   "Error reading the log configs",
		Severity:  "CRITICAL",
	},
	error1103: {
		ErrorCode: error1103,
		Message:   "Error while initializing authorization component, when intializing adapter REST API",
		Severity:  "BLOCKER",
	},
	error1104: {
		ErrorCode: error1104,
		Message:   "Readiness probe is not set as local api artifacts processing has failed.",
		Severity:  "CRITICAL",
	},
}
