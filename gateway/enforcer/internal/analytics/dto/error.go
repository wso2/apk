/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 
package dto

// Error represents the error attributes in an analytics event.
type Error struct {
	ErrorCode    int              `json:"errorCode"`
	ErrorMessage FaultSubCategory `json:"errorMessage"`
}

// GetErrorCode returns the error code.
func (e *Error) GetErrorCode() int {
	return e.ErrorCode
}

// SetErrorCode sets the error code.
func (e *Error) SetErrorCode(errorCode int) {
	e.ErrorCode = errorCode
}

// GetErrorMessage returns the error message.
func (e *Error) GetErrorMessage() FaultSubCategory {
	return e.ErrorMessage
}

// SetErrorMessage sets the error message.
func (e *Error) SetErrorMessage(errorMessage FaultSubCategory) {
	e.ErrorMessage = errorMessage
}
