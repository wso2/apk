//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

# Organization List CR definition.
#
# + kind - Field Description  
# + apiVersion - Field Description  
# + metadata - Field Description  
# + items - Field Description
public type OrganizationList record {
    string kind = "OrganizationList";
    string apiVersion = "cp.wso2.com/v1alpha1";
    ListMeta metadata;
    Organization[] items;

};

# Organization CR.
#
# + apiVersion - Field Description  
# + kind - Field Description  
# + metadata - Field Description  
# + spec - Field Description
public type Organization record {|
    string apiVersion = "cp.wso2.com/v1alpha1";
    string kind = "Organization";
    Metadata metadata;
    OrganizationSpec spec;
|};

# Organization definition.
#
# + uuid - uuid in control-plane  
# + name - name of organization.  
# + displayName - displayname organization.  
# + organizationClaimValue - organization claim value came from jwt.  
# + serviceListingNamespaces - service listing namespaces.
# + enabled - Organization enabld in system. 
# + properties - additional properties.
public type OrganizationSpec record {
    string uuid;
    string name;
    string displayName;
    string organizationClaimValue;
    boolean enabled = true;
    string[] serviceListingNamespaces = ["*"];
    OrganizationProperties[] properties = [];
};

# Organization AdditonalProperties Definition.
#
# + key - property key.  
# + value - property value.
public type OrganizationProperties record {
    string key;
    string value;
};