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

import ballerina/log;
import ballerina/uuid;

function addApplication(Application application, string org, string user) returns string?|Application|error {
    string applicationId = uuid:createType1AsString();
    application.applicationId = applicationId;
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        log:printError("Invalid Policy");
        return error("Invalid Policy");
    }
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        string?|Application|error createdApp = addApplicationDAO(application, subscriberId, org);
        return createdApp;
    } else {
        return subscriberId;
    }
}

function validateApplicationUsagePolicy(string policyName, string org) returns string?|error {
    string?|error policy = getApplicationUsagePlanByNameDAO(policyName,org);
    return policy;
}

function getApplicationById(string appId, string org) returns string?|Application|error {
    string?|Application|error application = getApplicationByIdDAO(appId, org);
    return application;
}

function getApplicationList(string? sortBy, string? groupId, string? query, string? sortOrder, int 'limit, int offset, string org) returns string?|ApplicationList|error {
    Application[]|error? applications = getApplicationsDAO(org);
    if applications is Application[] {
        int count = applications.length();
        ApplicationList applicationsList = {count: count, list: applications};
        return applicationsList;
    } else {
        return applications;
    }
}

function updateApplication(string appId, Application application, string org, string user) returns string?|Application|NotFoundError|error {
    string?|Application|error existingApp = getApplicationByIdDAO(appId, org);
    if existingApp is Application {
        application.applicationId = appId;
    } else {
        Error err = {code:9010101, message:"Application Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    }
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        log:printError("Invalid Policy");
        return error("Invalid Policy");
    }
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        string?|Application|error updatedApp = updateApplicationDAO(application, subscriberId, org);
        return updatedApp;
    } else {
        return subscriberId;
    }
}

function deleteApplication(string appId, string organization) returns string|error? {
    error?|string status = deleteApplicationDAO(appId,organization);
    return status;
}
