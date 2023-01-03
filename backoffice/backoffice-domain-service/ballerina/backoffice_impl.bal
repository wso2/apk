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

# This function used to get API from database
#
# + return - Return Value string?|APIList|error
isolated function getAPIList() returns string?|APIList|error {
    API[]|error? apis = db_getAPIsDAO();
    if apis is API[] {
        int count = apis.length();
        APIList apisList = {count: count, list: apis};
        return apisList;
    } else {
        return apis;
    }
}

# This function used to change the lifecycle of API
#
# + targetState - lifecycle action
# + apiId - API Id
# + organization - organization
# + return - Return Value LifecycleState|error
isolated function changeLifeCyleState(string targetState, string apiId, string organization) returns LifecycleState|error {
    string prevLCState = check db_getCurrentLCStatus(apiId, organization);
    transaction {
        string|error lcState = db_changeLCState(targetState, apiId, organization);
        if lcState is string {
                string newvLCState = check db_getCurrentLCStatus(apiId, organization);
                string|error lcEvent = db_AddLCEvent(apiId, prevLCState, newvLCState, organization);
                if lcEvent is string {
                    check commit;
                    json lcPayload = check getTransitionsFromState(targetState);
                    LifecycleState lcCr = check lcPayload.cloneWithType(LifecycleState);
                    return lcCr;
                } else {
                    rollback;
                    return error("error while adding LC event" + lcEvent.message());
                }
        } else {
            rollback;
            return error("error while updating LC state" + lcState.message());
        }
    } 
}

# This function used to get current state of the API.
#
# + apiId - API Id parameter
# + organization - organization
# + return - Return Value LifecycleState|error
isolated function getLifeCyleState(string apiId, string organization) returns LifecycleState|error {
    string|error currentLCState = db_getCurrentLCStatus(apiId, organization);
    if currentLCState is string {
        json lcPayload = check getTransitionsFromState(currentLCState);
        LifecycleState lcGet = check lcPayload.cloneWithType(LifecycleState);
        return lcGet;
    } else {
        return error("error while Getting LC state" + currentLCState.message());
    }
}

# This function used to map user action to LC state
#
# + v - any parameter object
# + return - Return LC state
isolated function actionToLCState(any v) returns string {
    if(v.toString().equalsIgnoreCaseAscii("published")){
        return "PUBLISHED";
    } else if(v.toString().equalsIgnoreCaseAscii("created")){
        return "CREATED";
    } else if(v.toString().equalsIgnoreCaseAscii("blocked")){
        return "BLOCKED";
    } else if(v.toString().equalsIgnoreCaseAscii("deprecated")){
        return "DEPRECATED";
    } else if(v.toString().equalsIgnoreCaseAscii("prototyped")){
        return "PROTOTYPED";
    } else if(v.toString().equalsIgnoreCaseAscii("retired")){
        return "RETIRED";
    } else {
        return "any";   
    }
}

# This function used to get the availble event transitions from state
#
# + state - state parameter
# + return - Return Value jsons
isolated function getTransitionsFromState(string state) returns json|error {
    StatesList c =  check lifeCycleStateTransitions.cloneWithType(StatesList);
    foreach States x in c.States {
        if(state.equalsIgnoreCaseAscii(x.State)) {
            return x.toJson();
        }
    }
    
}

# This function used to connect API create service to database
#
# + apiId - API Id parameter
# + return - Return Value LifecycleHistory
isolated function getLcEventHistory(string apiId) returns LifecycleHistory|error? {
    LifecycleHistoryItem[]|error? lcHistory = db_getLCEventHistory(apiId);
    if lcHistory is LifecycleHistoryItem[] {
        int count = lcHistory.length();
        LifecycleHistory eventList = {count: count, list: lcHistory};
        return eventList;
    } else {
        return error("Error while retriving LC events", lcHistory);
    }
}



isolated function getSubscriptions(string? apiId) returns SubscriptionList|error {
    Subscription[]|error subcriptions;
        subcriptions = check db_getSubscriptionsForAPI(apiId.toString());
        if subcriptions is Subscription[] {
            int count = subcriptions.length();
            SubscriptionList subsList = {count: count, list: subcriptions};
            return subsList;
        } else {
            return error(subcriptions.message());
        } 
}


isolated function blockSubscription(string subscriptionId, string blockState) returns string|error {
    if ("blocked".equalsIgnoreCaseAscii(blockState) || "prod_only_blocked".equalsIgnoreCaseAscii(blockState)) {
        error|string blockSub = db_blockSubscription(subscriptionId, blockState);
        if blockSub is error {
            return error("Error while blocking subscription");
        } else {
            return  blockSub;
        }
    } else {
        return error("Invalid block state provided");
    }
}

isolated function unblockSubscription(string subscriptionId) returns string|error {
    error|string unblockSub = db_unblockSubscription(subscriptionId);
    if unblockSub is error {
        return error("Error while unblocking subscription");
    } else {
        return  unblockSub;
    }
}
