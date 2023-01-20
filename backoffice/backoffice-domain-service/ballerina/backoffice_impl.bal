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
isolated function getAPIList() returns APIList|APKError {
    API[]|APKError apis = db_getAPIsDAO();
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
        json lcPayload =  check getTransitionsFromState(currentLCState);
        LifecycleState|error lcGet =  lcPayload.cloneWithType(LifecycleState);
        if lcGet is error {
            string message = "Error while retrieving connection";
            return error(message, message = message, description = message, code = 909000, statusCode = "500");
        }
        return lcGet;
    } else {
        return currentLCState;
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
isolated function getLcEventHistory(string apiId) returns LifecycleHistory|APKError {
    LifecycleHistoryItem[]|APKError lcHistory = db_getLCEventHistory(apiId);
    if lcHistory is LifecycleHistoryItem[] {
        int count = lcHistory.length();
        LifecycleHistory eventList = {count: count, list: lcHistory};
        return eventList;
    } else {
        return lcHistory;
    }
}



isolated function getSubscriptions(string? apiId) returns SubscriptionList|APKError {
    Subscription[]|APKError subcriptions;
        subcriptions = check db_getSubscriptionsForAPI(apiId.toString());
        if subcriptions is Subscription[] {
            int count = subcriptions.length();
            SubscriptionList subsList = {count: count, list: subcriptions};
            return subsList;
        } else {
            return subcriptions;
        } 
}


isolated function blockSubscription(string subscriptionId, string blockState) returns string|APKError {
    if ("blocked".equalsIgnoreCaseAscii(blockState) || "prod_only_blocked".equalsIgnoreCaseAscii(blockState)) {
        APKError|string blockSub = db_blockSubscription(subscriptionId, blockState);
        return blockSub;
    } else {
        string message = "Invalid blockState provided";
        return error(message, message = message, description = message, code = 909002, statusCode = "400");    
    }
}

isolated function unblockSubscription(string subscriptionId) returns string|APKError {
    APKError|string unblockSub = db_unblockSubscription(subscriptionId);
    return  unblockSub;
}

isolated function getAPI(string apiId) returns API|APKError {
    API|APKError getAPI = check db_getAPI(apiId);
    return  getAPI;
}

isolated function getAPIDefinition(string apiId) returns APIDefinition|NotFoundError|APKError {
    APIDefinition|NotFoundError|APKError apiDefinition = db_getAPIDefinition(apiId);
    return apiDefinition;
}


isolated function updateAPI(string apiId, ModifiableAPI payload, string organization) returns API|APKError {
    API|APKError api = db_updateAPI(apiId, payload, organization);
    return api;
}

isolated function handleAPKError(APKError errorDetail) returns InternalServerErrorError|BadRequestError{
            ErrorHandler & readonly detail = errorDetail.detail();
        if detail.statusCode == "400" {
            BadRequestError badRequest = {body: {code: detail.code, message: detail.message}};
            return badRequest;
        }
        InternalServerErrorError internalServerError = {body: {code: detail.code, message: detail.message}};
        return internalServerError;
}
