/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package dto

//public
//type Organization record
//{|
//string uuid;
//string name;
//string displayName;
//string organizationClaimValue;
//boolean enabled;
//string[] serviceListingNamespaces = ["*"];
//OrganizationProperties[] properties = [];
//|}

// Organization represents the organization
type Organization struct {
	UUID                     string                   `json:"uuid" yaml:"uuid"`
	Name                     string                   `json:"name" yaml:"name"`
	DisplayName              string                   `json:"displayName" yaml:"displayName"`
	OrganizationClaimValue   string                   `json:"organizationClaimValue" yaml:"organizationClaimValue"`
	Enabled                  bool                     `json:"enabled" yaml:"enabled"`
	ServiceListingNamespaces []string                 `json:"serviceListingNamespaces" yaml:"serviceListingNamespaces"`
	Properties               []OrganizationProperties `json:"properties" yaml:"properties"`
}

// OrganizationProperties represents the organization Additional Properties
type OrganizationProperties struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// NewOrganization creates a new Organization with default values
func NewOrganization(uuid, name, displayName, organizationClaimValue string, enabled bool) *Organization {
	return &Organization{
		UUID:                     uuid,
		Name:                     name,
		DisplayName:              displayName,
		OrganizationClaimValue:   organizationClaimValue,
		Enabled:                  enabled,
		ServiceListingNamespaces: []string{"*"},
		Properties:               []OrganizationProperties{},
	}
}
