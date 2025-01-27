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

// ExtendedAPI represents an extended API object with organization and environment details.
type ExtendedAPI struct {
	API
	OrganizationID string `json:"organizationId"`
	EnvironmentID  string `json:"environmentId"`
	APIContext     string `json:"apiContext"`
}

// GetOrganizationID returns the organization ID.
func (e *ExtendedAPI) GetOrganizationID() string {
	return e.OrganizationID
}

// SetOrganizationID sets the organization ID.
func (e *ExtendedAPI) SetOrganizationID(organizationID string) {
	e.OrganizationID = organizationID
}

// GetEnvironmentID returns the environment ID.
func (e *ExtendedAPI) GetEnvironmentID() string {
	return e.EnvironmentID
}

// SetEnvironmentID sets the environment ID.
func (e *ExtendedAPI) SetEnvironmentID(environmentID string) {
	e.EnvironmentID = environmentID
}

// GetAPIContext returns the API context.
func (e *ExtendedAPI) GetAPIContext() string {
	return e.APIContext
}

// SetAPIContext sets the API context.
func (e *ExtendedAPI) SetAPIContext(apiContext string) {
	e.APIContext = apiContext
}
