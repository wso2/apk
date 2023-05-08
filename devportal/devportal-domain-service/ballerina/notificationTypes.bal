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

public type ApplicationGRPC record {|
    string eventId = "";
    string name = "";
    string uuid = "";
    string owner = "";
    string policy = "";
    Application_Key[] keys = [];
    string organization = "";
    string timeStamp = "";
    record {|string key; string value;|}[] attributes = [];
|};

public type Application_Key record {|
    string key = "";
    string keyManager = "";
|};

public type SubscriptionGRPC record {|
    string eventId = "";
    string applicationRef = "";
    string apiRef = "";
    string policyId = "";
    string subStatus = "";
    string subscriber = "";
    string uuid = "";
    string timeStamp = "";
    string organization = "";
|};

public type NotificationResponse record {|
    NotificationResponse_StatusCode code = UNKNOWN;
|};

public enum NotificationResponse_StatusCode {
    UNKNOWN, OK, FAILED
}
