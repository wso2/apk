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

function getAPIList() returns string?|APIList|error {
    API[]|error? apis = db_getAPIsDAO();
    if apis is API[] {
        int count = apis.length();
        APIList apisList = {count: count, list: apis};
        return apisList;
    } else {
        return apis;
    }
}

function changeLifeCyleState(string action, string apiId, string organization) returns LifecycleState|error {
    string|error? lcState = db_changeLCState(action, apiId, organization);
    if lcState is string {
            LifecycleState lcStateCr = {state: lcState, availableTransitions: [
                {event: "", targetState: ""}
            ]};
            return lcStateCr;
    } else {
        return error("error while updating LC state");
    }
} 


function getLifeCyleState(string apiId, string organization) returns LifecycleState|error {
    string|error? currentLCState = db_getCurrentLCStatus(apiId, organization);
    if currentLCState is string {
            LifecycleState lcStateCr = {state: currentLCState, availableTransitions: [
                {event: "", targetState: ""}
            ]};
        return lcStateCr;
    } else {
        return error("error while updating LC state");
    }
}

function actionToLCState(any v) returns string {

    match v {
        "Demote to Created" => { return "Created"; }
        "Publish" => { return "Published"; }
        "Block" => { return "Blocked"; }
        "Deprecate" => { return "Deprecateed"; }
        "Retire" => { return "Retired"; }
        "Re-Publish" => { return "Published"; }
        _ => { return "any"; }
    }
}