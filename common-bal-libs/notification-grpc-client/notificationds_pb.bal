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

import ballerina/grpc;
import ballerina/protobuf;

public const string NOTIFICATIONDS_DESC = "0A146E6F74696669636174696F6E64732E70726F746F1218646973636F766572792E736572766963652E61706B6D677422A3010A0B4170706C69636174696F6E12180A076576656E74496418012001280952076576656E74496412240A0D6170706C69636174696F6E4964180220012809520D6170706C69636174696F6E496412120A0475756964180320012809520475756964121C0A0974696D655374616D70180420012809520974696D655374616D7012220A0C6F7267616E697A6174696F6E180520012809520C6F7267616E697A6174696F6E22A4010A0C537562736372697074696F6E12180A076576656E74496418012001280952076576656E74496412240A0D6170706C69636174696F6E4964180220012809520D6170706C69636174696F6E496412120A0475756964180320012809520475756964121C0A0974696D655374616D70180420012809520974696D655374616D7012220A0C6F7267616E697A6174696F6E180520012809520C6F7267616E697A6174696F6E2294010A144E6F74696669636174696F6E526573706F6E7365124D0A04636F646518012001280E32392E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E73652E537461747573436F64655204636F6465222D0A0A537461747573436F6465120B0A07554E4B4E4F574E100012060A024F4B1001120A0A064641494C4544100232A3050A134E6F74696669636174696F6E53657276696365126A0A114372656174654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126A0A115570646174654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126A0A1144656C6574654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A12437265617465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A12557064617465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A1244656C657465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E73654282010A306F72672E77736F322E63686F72656F2E636F6E6E6563742E646973636F766572792E736572766963652E61706B6D677442136E6F74696669636174696F6E447350726F746F50005A346769746875622E636F6D2F77736F322F61706B2F616461707465722F646973636F766572792F736572766963652F61706B6D6774880101620670726F746F33";

public isolated client class NotificationServiceClient {
    *grpc:AbstractClientEndpoint;

    private final grpc:Client grpcClient;

    public isolated function init(string url, *grpc:ClientConfiguration config) returns grpc:Error? {
        self.grpcClient = check new (url, config);
        check self.grpcClient.initStub(self, NOTIFICATIONDS_DESC);
    }

    isolated remote function CreateApplication(Application|ContextApplication req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/CreateApplication", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function CreateApplicationContext(Application|ContextApplication req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/CreateApplication", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }

    isolated remote function UpdateApplication(Application|ContextApplication req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/UpdateApplication", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function UpdateApplicationContext(Application|ContextApplication req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/UpdateApplication", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }

    isolated remote function DeleteApplication(Application|ContextApplication req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/DeleteApplication", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function DeleteApplicationContext(Application|ContextApplication req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Application message;
        if req is ContextApplication {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/DeleteApplication", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }

    isolated remote function CreateSubscription(Subscription|ContextSubscription req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/CreateSubscription", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function CreateSubscriptionContext(Subscription|ContextSubscription req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/CreateSubscription", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }

    isolated remote function UpdateSubscription(Subscription|ContextSubscription req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/UpdateSubscription", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function UpdateSubscriptionContext(Subscription|ContextSubscription req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/UpdateSubscription", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }

    isolated remote function DeleteSubscription(Subscription|ContextSubscription req) returns NotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/DeleteSubscription", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <NotificationResponse>result;
    }

    isolated remote function DeleteSubscriptionContext(Subscription|ContextSubscription req) returns ContextNotificationResponse|grpc:Error {
        map<string|string[]> headers = {};
        Subscription message;
        if req is ContextSubscription {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.NotificationService/DeleteSubscription", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <NotificationResponse>result, headers: respHeaders};
    }
}

public type ContextNotificationResponse record {|
    NotificationResponse content;
    map<string|string[]> headers;
|};

public type ContextSubscription record {|
    Subscription content;
    map<string|string[]> headers;
|};

public type ContextApplication record {|
    Application content;
    map<string|string[]> headers;
|};

@protobuf:Descriptor {value: NOTIFICATIONDS_DESC}
public type NotificationResponse record {|
    NotificationResponse_StatusCode code = UNKNOWN;
|};

public enum NotificationResponse_StatusCode {
    UNKNOWN, OK, FAILED
}

@protobuf:Descriptor {value: NOTIFICATIONDS_DESC}
public type Subscription record {|
    string eventId = "";
    string applicationId = "";
    string uuid = "";
    string timeStamp = "";
    string organization = "";
|};

@protobuf:Descriptor {value: NOTIFICATIONDS_DESC}
public type Application record {|
    string eventId = "";
    string applicationId = "";
    string uuid = "";
    string timeStamp = "";
    string organization = "";
|};

