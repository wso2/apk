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

// GetErrorByCode used to keep error details for error logs
func GetErrorByCode(code int, args ...interface{}) logging.ErrorDetails {
	errorLog, ok := Mapper[code]
	if !ok {
		errorLog = logging.ErrorDetails{
			ErrorCode: 0000,
			Message:   fmt.Sprintf("No error message found for error code: %v", code),
			Severity:  "BLOCKER",
		}
	}
	message := errorLog.Message
	message = fmt.Sprintf(message, args...)

	errorLog.Message = message
	return errorLog
}
