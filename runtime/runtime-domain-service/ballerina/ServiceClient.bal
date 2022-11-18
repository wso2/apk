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

# This returns services in a namsepace.
#
# + namespace - namespace value
# + return - list of services in namespace.
function getServicesListInNamespace(string namespace) returns ServiceList|error {
    Service[] servicesList = getServicesList();
    Service[] filteredList = [];
    foreach Service item in servicesList {
        if item.namespace == namespace {
            filteredList.push(item);
        }
    }
    return {list: filteredList, pagination: {total: filteredList.length()}};
}

# This returns list of services in all namespaces.
# + return - list of services in namespaces.
function getServicesListFromK8s() returns ServiceList|error {
    return {list: getServicesList(), pagination: {total: getServicesList().length()}};
}

# This retrieve specific service from name space.
#
# + name - name of service.
# + namespace - namespace of service.
# + return - service in namespace.
function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
    Service? serviceResult = getService(name, namespace);
    if serviceResult is null {
        return {list: []};
    } else {
        return {list: [serviceResult]};
    }
}
