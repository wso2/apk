import ballerina/grpc;
import ballerina/protobuf;

const string APIDS_DESC = "0A0B61706964732E70726F746F1218646973636F766572792E736572766963652E61706B6D67742281040A0341504912120A0475756964180120012809520475756964121A0A0870726F7669646572180220012809520870726F766964657212180A0776657273696F6E180320012809520776657273696F6E12120A046E616D6518042001280952046E616D6512180A07636F6E746578741805200128095207636F6E7465787412120A047479706518062001280952047479706512260A0E6F7267616E697A6174696F6E4964180720012809520E6F7267616E697A6174696F6E4964121C0A09637265617465644279180820012809520963726561746564427912200A0B6372656174656454696D65180920012809520B6372656174656454696D65121C0A09757064617465644279180A20012809520975706461746564427912200A0B7570646174656454696D65180B20012809520B7570646174656454696D65121E0A0A646566696E6974696F6E180C20012809520A646566696E6974696F6E121E0A0A7472616E73706F727473180D20032809520A7472616E73706F72747312400A097265736F757263657318102003280B32222E646973636F766572792E736572766963652E61706B6D67742E5265736F7572636552097265736F757263657312440A0A636F7273436F6E66696718112001280B32242E646973636F766572792E736572766963652E61706B6D67742E436F7273436F6E666967520A636F7273436F6E66696722C8020A0A436F7273436F6E666967123A0A18636F7273436F6E66696775726174696F6E456E61626C65641801200128085218636F7273436F6E66696775726174696F6E456E61626C6564123C0A19616363657373436F6E74726F6C416C6C6F774F726967696E731802200328095219616363657373436F6E74726F6C416C6C6F774F726967696E7312440A1D616363657373436F6E74726F6C416C6C6F7743726564656E7469616C73180320012808521D616363657373436F6E74726F6C416C6C6F7743726564656E7469616C73123C0A19616363657373436F6E74726F6C416C6C6F77486561646572731804200328095219616363657373436F6E74726F6C416C6C6F7748656164657273123C0A19616363657373436F6E74726F6C416C6C6F774D6574686F64731805200328095219616363657373436F6E74726F6C416C6C6F774D6574686F647322FE020A085265736F7572636512120A047061746818012001280952047061746812120A047665726218022001280952047665726212520A0F61757468656E7469636174696F6E7318032003280B32282E646973636F766572792E736572766963652E61706B6D67742E41757468656E7469636174696F6E520F61757468656E7469636174696F6E7312370A0673636F70657318042003280B321F2E646973636F766572792E736572766963652E61706B6D67742E53636F7065520673636F70657312590A116F7065726174696F6E506F6C696369657318062001280B322B2E646973636F766572792E736572766963652E61706B6D67742E4F7065726174696F6E506F6C696369657352116F7065726174696F6E506F6C696369657312460A0B7175657279506172616D7318072003280B32242E646973636F766572792E736572766963652E61706B6D67742E5175657279506172616D520B7175657279506172616D73121A0A08686F73746E616D651808200328095208686F73746E616D6522130A114F7065726174696F6E506F6C696369657322360A0A5175657279506172616D12120A046E616D6518012001280952046E616D6512140A0576616C7565180220012809520576616C7565227B0A0553636F706512120A046E616D6518012001280952046E616D6512200A0B646973706C61794E616D65180220012809520B646973706C61794E616D6512200A0B6465736372697074696F6E180320012809520B6465736372697074696F6E121A0A0862696E64696E6773180420032809520862696E64696E677322B0010A0E41757468656E7469636174696F6E12120A047479706518012001280952047479706512100A03697373180220012809520369737312100A03617564180320012809520361756412180A076A776B7355726918042001280952076A776B73557269124C0A0E63726564656E7469616C4C69737418052003280B32242E646973636F766572792E736572766963652E61706B6D67742E43726564656E7469616C520E63726564656E7469616C4C69737422440A0A43726564656E7469616C121A0A08757365726E616D651801200128095208757365726E616D65121A0A0870617373776F7264180220012809520870617373776F726422220A08526573706F6E736512160A06726573756C741801200128085206726573756C7432FC010A0A41504953657276696365124E0A09637265617465415049121D2E646973636F766572792E736572766963652E61706B6D67742E4150491A222E646973636F766572792E736572766963652E61706B6D67742E526573706F6E7365124E0A09757064617465415049121D2E646973636F766572792E736572766963652E61706B6D67742E4150491A222E646973636F766572792E736572766963652E61706B6D67742E526573706F6E7365124E0A0964656C657465415049121D2E646973636F766572792E736572766963652E61706B6D67742E4150491A222E646973636F766572792E736572766963652E61706B6D67742E526573706F6E736542365A346769746875622E636F6D2F77736F322F61706B2F616461707465722F646973636F766572792F736572766963652F61706B6D6774620670726F746F33";

public isolated client class APIServiceClient {
    *grpc:AbstractClientEndpoint;

    private final grpc:Client grpcClient;

    public isolated function init(string url, *grpc:ClientConfiguration config) returns grpc:Error? {
        self.grpcClient = check new (url, config);
        check self.grpcClient.initStub(self, APIDS_DESC);
    }

    isolated remote function createAPI(API|ContextAPI req) returns Response|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/createAPI", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Response>result;
    }

    isolated remote function createAPIContext(API|ContextAPI req) returns ContextResponse|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/createAPI", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Response>result, headers: respHeaders};
    }

    isolated remote function updateAPI(API|ContextAPI req) returns Response|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/updateAPI", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Response>result;
    }

    isolated remote function updateAPIContext(API|ContextAPI req) returns ContextResponse|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/updateAPI", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Response>result, headers: respHeaders};
    }

    isolated remote function deleteAPI(API|ContextAPI req) returns Response|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/deleteAPI", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Response>result;
    }

    isolated remote function deleteAPIContext(API|ContextAPI req) returns ContextResponse|grpc:Error {
        map<string|string[]> headers = {};
        API message;
        if req is ContextAPI {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("discovery.service.apkmgt.APIService/deleteAPI", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Response>result, headers: respHeaders};
    }
}

public client class APIServiceResponseCaller {
    private grpc:Caller caller;

    public isolated function init(grpc:Caller caller) {
        self.caller = caller;
    }

    public isolated function getId() returns int {
        return self.caller.getId();
    }

    isolated remote function sendResponse(Response response) returns grpc:Error? {
        return self.caller->send(response);
    }

    isolated remote function sendContextResponse(ContextResponse response) returns grpc:Error? {
        return self.caller->send(response);
    }

    isolated remote function sendError(grpc:Error response) returns grpc:Error? {
        return self.caller->sendError(response);
    }

    isolated remote function complete() returns grpc:Error? {
        return self.caller->complete();
    }

    public isolated function isCancelled() returns boolean {
        return self.caller.isCancelled();
    }
}

public type ContextResponse record {|
    Response content;
    map<string|string[]> headers;
|};

public type ContextAPI record {|
    API content;
    map<string|string[]> headers;
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type Response record {|
    boolean result = false;
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type OperationPolicies record {|
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type Credential record {|
    string username = "";
    string password = "";
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type CorsConfig record {|
    boolean corsConfigurationEnabled = false;
    string[] accessControlAllowOrigins = [];
    boolean accessControlAllowCredentials = false;
    string[] accessControlAllowHeaders = [];
    string[] accessControlAllowMethods = [];
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type Scope record {|
    string name = "";
    string displayName = "";
    string description = "";
    string[] bindings = [];
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type QueryParam record {|
    string name = "";
    string value = "";
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type Authentication record {|
    string 'type = "";
    string iss = "";
    string aud = "";
    string jwksUri = "";
    Credential[] credentialList = [];
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type Resource record {|
    string path = "";
    string verb = "";
    Authentication[] authentications = [];
    Scope[] scopes = [];
    OperationPolicies operationPolicies = {};
    QueryParam[] queryParams = [];
    string[] hostname = [];
|};

@protobuf:Descriptor {value: APIDS_DESC}
public type API record {|
    string uuid = "";
    string provider = "";
    string 'version = "";
    string name = "";
    string context = "";
    string 'type = "";
    string organizationId = "";
    string createdBy = "";
    string createdTime = "";
    string updatedBy = "";
    string updatedTime = "";
    string definition = "";
    string[] transports = [];
    Resource[] resources = [];
    CorsConfig corsConfig = {};
|};

