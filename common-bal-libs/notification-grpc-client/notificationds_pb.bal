import ballerina/grpc;
import ballerina/protobuf;

public const string NOTIFICATIONDS_DESC = "0A146E6F74696669636174696F6E64732E70726F746F1218646973636F766572792E736572766963652E61706B6D677422CD030A0B4170706C69636174696F6E12180A076576656E74496418012001280952076576656E74496412120A046E616D6518022001280952046E616D6512120A047575696418032001280952047575696412140A056F776E657218042001280952056F776E657212160A06706F6C6963791805200128095206706F6C69637912550A0A6174747269627574657318062003280B32352E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E2E41747472696275746573456E747279520A61747472696275746573123D0A046B65797318072003280B32292E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E2E4B657952046B65797312220A0C6F7267616E697A6174696F6E180820012809520C6F7267616E697A6174696F6E121C0A0974696D655374616D70180920012809520974696D655374616D701A3D0A0F41747472696275746573456E74727912100A036B657918012001280952036B657912140A0576616C7565180220012809520576616C75653A0238011A370A034B657912100A036B657918012001280952036B6579121E0A0A6B65794D616E61676572180220012809520A6B65794D616E616765722298020A0C537562736372697074696F6E12180A076576656E74496418012001280952076576656E74496412260A0E6170706C69636174696F6E526566180220012809520E6170706C69636174696F6E52656612160A066170695265661803200128095206617069526566121A0A08706F6C69637949641804200128095208706F6C6963794964121C0A097375625374617475731805200128095209737562537461747573121E0A0A73756273637269626572180620012809520A7375627363726962657212120A0475756964180720012809520475756964121C0A0974696D655374616D70180820012809520974696D655374616D7012220A0C6F7267616E697A6174696F6E180920012809520C6F7267616E697A6174696F6E2294010A144E6F74696669636174696F6E526573706F6E7365124D0A04636F646518012001280E32392E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E73652E537461747573436F64655204636F6465222D0A0A537461747573436F6465120B0A07554E4B4E4F574E100012060A024F4B1001120A0A064641494C4544100232A3050A134E6F74696669636174696F6E53657276696365126A0A114372656174654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126A0A115570646174654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126A0A1144656C6574654170706C69636174696F6E12252E646973636F766572792E736572766963652E61706B6D67742E4170706C69636174696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A12437265617465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A12557064617465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E7365126C0A1244656C657465537562736372697074696F6E12262E646973636F766572792E736572766963652E61706B6D67742E537562736372697074696F6E1A2E2E646973636F766572792E736572766963652E61706B6D67742E4E6F74696669636174696F6E526573706F6E73654280010A2E6F72672E77736F322E61706B2E656E666F726365722E646973636F766572792E736572766963652E61706B6D677442136E6F74696669636174696F6E447350726F746F50005A346769746875622E636F6D2F77736F322F61706B2F616461707465722F646973636F766572792F736572766963652F61706B6D6774880101620670726F746F33";

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
    string applicationRef = "";
    string apiRef = "";
    string policyId = "";
    string subStatus = "";
    string subscriber = "";
    string uuid = "";
    string timeStamp = "";
    string organization = "";
|};

@protobuf:Descriptor {value: NOTIFICATIONDS_DESC}
public type Application record {|
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

@protobuf:Descriptor {value: NOTIFICATIONDS_DESC}
public type Application_Key record {|
    string key = "";
    string keyManager = "";
|};

