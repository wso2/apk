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

import ballerina/log;

public isolated function createApplication(Application createApplicationRequest, string endpoint) returns error|NotificationResponse {
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse createApplicationResponse = check ep->CreateApplication(createApplicationRequest);
    log:printDebug(createApplicationResponse.toString());
    return createApplicationResponse;
}

public isolated function updateApplication(Application updateApplicationRequest, string endpoint) returns error|NotificationResponse {
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse updateApplicationResponse = check ep->UpdateApplication(updateApplicationRequest);
    log:printDebug(updateApplicationResponse.toString());
    return updateApplicationResponse;
}

public isolated function deleteApplication(Application deleteApplicationRequest, string endpoint) returns error|NotificationResponse {
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse deleteApplicationResponse = check ep->DeleteApplication(deleteApplicationRequest);
    log:printDebug(deleteApplicationResponse.toString());
    return deleteApplicationResponse;
}

public isolated function createSubscription(Subscription createSubscriptionRequest, string endpoint) returns error|NotificationResponse {
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse createSubscriptionResponse = check ep->CreateSubscription(createSubscriptionRequest);
    log:printDebug(createSubscriptionResponse.toString());
    return createSubscriptionResponse;
}

public isolated function updateSubscription(Subscription updateSubscriptionRequest, string endpoint) returns error|NotificationResponse {
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse updateSubscriptionResponse = check ep->UpdateSubscription(updateSubscriptionRequest);
    log:printDebug(updateSubscriptionResponse.toString());
    return updateSubscriptionResponse;
}

public isolated function deleteSubscription(Subscription deleteSubscriptionRequest, string endpoint) returns error|NotificationResponse{
    NotificationServiceClient ep = check new (endpoint);
    NotificationResponse deleteSubscriptionResponse = check ep->DeleteSubscription(deleteSubscriptionRequest);
    log:printDebug(deleteSubscriptionResponse.toString());
    return deleteSubscriptionResponse;
}
