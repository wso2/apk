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

// Application represents application attributes in an analytics event.
type Application struct {
	KeyType          string `json:"keyType"`
	ApplicationID    string `json:"applicationId"`
	ApplicationName  string `json:"applicationName"`
	ApplicationOwner string `json:"applicationOwner"`
}

// GetKeyType returns the key type.
func (a *Application) GetKeyType() string {
	return a.KeyType
}

// SetKeyType sets the key type.
func (a *Application) SetKeyType(keyType string) {
	a.KeyType = keyType
}

// GetApplicationID returns the application ID.
func (a *Application) GetApplicationID() string {
	return a.ApplicationID
}

// SetApplicationID sets the application ID.
func (a *Application) SetApplicationID(applicationID string) {
	a.ApplicationID = applicationID
}

// GetApplicationName returns the application name.
func (a *Application) GetApplicationName() string {
	return a.ApplicationName
}

// SetApplicationName sets the application name.
func (a *Application) SetApplicationName(applicationName string) {
	a.ApplicationName = applicationName
}

// GetApplicationOwner returns the application owner.
func (a *Application) GetApplicationOwner() string {
	return a.ApplicationOwner
}

// SetApplicationOwner sets the application owner.
func (a *Application) SetApplicationOwner(applicationOwner string) {
	a.ApplicationOwner = applicationOwner
}
