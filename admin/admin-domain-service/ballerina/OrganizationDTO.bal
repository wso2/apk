//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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

# Description
#
# + id - Org ID  
# + name - Org Name 
# + displayName - Org Display Name
# + enabled -   Org Enabled
# + claimKey -  Org Claim Key 
# + production -  Org Production
# + sandbox -  Org Sandbox
# + serviceNamespaces -  Org Service Namespaces
# + claimValue -    Org Claim Value
public type Organizations record {
    string id;
    string name;
    string displayName;
    boolean enabled;
    string[] serviceNamespaces;
    string claimKey;
    string production;
    string sandbox;
    string claimValue;
};


public type Internal_Organization record {
    string id;
    string name;
    string displayName;
    boolean enabled;
    string[] serviceNamespaces;
    string[] production?;
    string[] sandbox?;
    OrganizationClaim[] claimList;
};

public type OrganizationClaim record {
    string claimKey?;
    string claimValue?;
};