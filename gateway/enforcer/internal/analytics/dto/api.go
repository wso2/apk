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

// API represents the API attributes in an analytics event.
type API struct {
	APIID                  string `json:"apiId"`
	APIType                string `json:"apiType"`
	APIName                string `json:"apiName"`
	APIVersion             string `json:"apiVersion"`
	APICreator             string `json:"apiCreator"`
	APICreatorTenantDomain string `json:"apiCreatorTenantDomain"`
}

// GetAPIID returns the API ID.
func (a *API) GetAPIID() string {
	return a.APIID
}

// SetAPIID sets the API ID.
func (a *API) SetAPIID(apiID string) {
	a.APIID = apiID
}

// GetAPIType returns the API type.
func (a *API) GetAPIType() string {
	return a.APIType
}

// SetAPIType sets the API type.
func (a *API) SetAPIType(apiType string) {
	a.APIType = apiType
}

// GetAPIName returns the API name.
func (a *API) GetAPIName() string {
	return a.APIName
}

// SetAPIName sets the API name.
func (a *API) SetAPIName(apiName string) {
	a.APIName = apiName
}

// GetAPIVersion returns the API version.
func (a *API) GetAPIVersion() string {
	return a.APIVersion
}

// SetAPIVersion sets the API version.
func (a *API) SetAPIVersion(apiVersion string) {
	a.APIVersion = apiVersion
}

// GetAPICreator returns the API creator.
func (a *API) GetAPICreator() string {
	return a.APICreator
}

// SetAPICreator sets the API creator.
func (a *API) SetAPICreator(apiCreator string) {
	a.APICreator = apiCreator
}

// GetAPICreatorTenantDomain returns the API creator's tenant domain.
func (a *API) GetAPICreatorTenantDomain() string {
	return a.APICreatorTenantDomain
}

// SetAPICreatorTenantDomain sets the API creator's tenant domain.
func (a *API) SetAPICreatorTenantDomain(apiCreatorTenantDomain string) {
	a.APICreatorTenantDomain = apiCreatorTenantDomain
}
