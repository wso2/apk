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

// AIMetadata represents AI metadata in an analytics event.
type AIMetadata struct {
	Model         string `json:"model"`
	VendorName    string `json:"vendor_name"`
	VendorVersion string `json:"vendor_version"`
}

// GetModel returns the model.
func (a *AIMetadata) GetModel() string {
	return a.Model
}

// SetModel sets the model.
func (a *AIMetadata) SetModel(model string) {
	a.Model = model
}

// GetVendorName returns the vendor name.
func (a *AIMetadata) GetVendorName() string {
	return a.VendorName
}

// SetVendorName sets the vendor name.
func (a *AIMetadata) SetVendorName(vendorName string) {
	a.VendorName = vendorName
}

// GetVendorVersion returns the vendor version.
func (a *AIMetadata) GetVendorVersion() string {
	return a.VendorVersion
}

// SetVendorVersion sets the vendor version.
func (a *AIMetadata) SetVendorVersion(vendorVersion string) {
	a.VendorVersion = vendorVersion
}
