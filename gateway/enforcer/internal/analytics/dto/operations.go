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

// Operation represents operation attributes in an analytics event.
type Operation struct {
	APIMethod           string `json:"apiMethod"`
	APIResourceTemplate string `json:"apiResourceTemplate"`
}

// GetAPIMethod returns the API method.
func (o *Operation) GetAPIMethod() string {
	return o.APIMethod
}

// SetAPIMethod sets the API method.
func (o *Operation) SetAPIMethod(apiMethod string) {
	o.APIMethod = apiMethod
}

// GetAPIResourceTemplate returns the API resource template.
func (o *Operation) GetAPIResourceTemplate() string {
	return o.APIResourceTemplate
}

// SetAPIResourceTemplate sets the API resource template.
func (o *Operation) SetAPIResourceTemplate(apiResourceTemplate string) {
	o.APIResourceTemplate = apiResourceTemplate
}
