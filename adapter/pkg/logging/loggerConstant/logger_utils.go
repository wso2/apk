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
	"fmt"

	logging "github.com/wso2/apk/adapter/pkg/logging"
)

// Mapper used to keep error details for error logs
var Mapper map[int]logging.ErrorDetails = make(map[int]logging.ErrorDetails)

// CombineMapper used to keep error details for error logs
func CombineMapper(comMap map[int]logging.ErrorDetails) map[int]logging.ErrorDetails {
	for key, value := range comMap {
		Mapper[key] = value
	}
	return Mapper
}

// GetErrorByCode used to keep error details for error logs
func GetErrorByCode(code int, args ...interface{}) logging.ErrorDetails {
	errorLog := Mapper[error1101]
	message := errorLog.Message
	for item := range args {
		message += fmt.Sprintf(" %v", item)
	}
	errorLog.Message = message
	return errorLog
}
