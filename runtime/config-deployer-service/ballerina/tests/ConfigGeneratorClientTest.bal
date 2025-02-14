import ballerina/http;
import ballerina/test;
import wso2/apk_common_lib;
import ballerina/io;
import ballerina/mime;

@test:Config {}
public function testGetGeneratedAPKConf() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/api.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);
    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText(API_TYPE_REST);
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is OkAnydata {
        test:assertEquals(response.status.code, http:STATUS_OK, "Status code mismatched");
        string content = check io:fileReadString("./tests/resources/expectedAPK.apk-conf");
        test:assertEquals(response.body.toString(), content, "APK conf content mismatched");
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromUrl() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity urlPart = new;
    urlPart.setText("http://localhost:9090/test/definition");
    urlPart.setContentDisposition(getContentDispositionForFormData("url"));
    bodyParts.push(urlPart);
    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText(API_TYPE_REST);
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is OkAnydata {
        test:assertEquals(response.status.code, http:STATUS_OK, "Status code mismatched");
        string content = check io:fileReadString("./tests/resources/expectedAPK.apk-conf");
        test:assertEquals(response.body.toString(), content, "APK conf content mismatched.");
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromUrlAndDefinitionNegative() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity urlPart = new;
    urlPart.setText("http://localhost:9090/test/definition");
    urlPart.setContentDisposition(getContentDispositionForFormData("url"));
    bodyParts.push(urlPart);
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/api.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);

    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText(API_TYPE_REST);
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is BadRequestError {
        BadRequestError badRequest = {body: {code: 90091, message: "Specify either definition or url"}};
        test:assertEquals(response.status.code, http:STATUS_BAD_REQUEST, "Status code mismatched");
        test:assertEquals(response.body, badRequest.body);
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromUrlAndDefinitionNoType() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/api.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);
    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText("newType");
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);

    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is BadRequestError {
        BadRequestError badRequest = {body: {code: 90091, message: "Invalid API Type"}};
        test:assertEquals(response.status.code, http:STATUS_BAD_REQUEST, "Status code mismatched");
        test:assertEquals(response.body, badRequest.body);
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromUrlAndDefinitionNotDefinedType() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/api.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is OkAnydata {
        test:assertEquals(response.status.code, http:STATUS_OK, "Status code mismatched");
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfInvalidDefinition() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/invalidapi.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);
    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText(API_TYPE_REST);
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is BadRequestError {
        test:assertEquals(response.status.code, http:STATUS_BAD_REQUEST, "Status code mismatched");
        BadRequestError badRequest = {
            body: {
                code: 90091,
                message: "Invalid API Definition",
                'error: [
                    {
                        code: "900754",
                        message: "Error while parsing OpenAPI definition",
                        description: "attribute paths.'/status/{codes}'(post).parameters is not of type `array`"
                    },
                    {
                        code: "900754",
                        message: "Error while parsing OpenAPI definition",
                        description: "attribute paths.'/status/{codes}'(post).responses is missing"
                    },
                    {
                        code: "900754",
                        message: "Error while parsing OpenAPI definition",
                        description: "attribute paths.'/stream/{n}'(get).responses is missing"
                    },
                    {
                        code: "900754",
                        message: "Error while parsing OpenAPI definition",
                        description: "paths.'/status/{codes}'. Declared path parameter codes needs to be defined as a path parameter in path or operation level"
                    }
                ]
            }
        };
        test:assertEquals(response.body, badRequest.body);
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromInvalidUrl() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity urlPart = new;
    urlPart.setText("http://localhost:9090/test/definition1");
    urlPart.setContentDisposition(getContentDispositionForFormData("url"));
    bodyParts.push(urlPart);
    mime:Entity apiTypeBodyPart = new;
    apiTypeBodyPart.setText(API_TYPE_REST);
    apiTypeBodyPart.setContentDisposition(getContentDispositionForFormData("apiType"));
    bodyParts.push(apiTypeBodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is apk_common_lib:APKError {
        test:assertEquals(response.toBalString(),e909044().toBalString());
    } else {
        test:assertFail("Error occurred while generating APK conf");
    }
}

function getContentDispositionForFormData(string partName) returns (mime:ContentDisposition) {
    mime:ContentDisposition contentDisposition = new;
    contentDisposition.name = partName;
    contentDisposition.disposition = "form-data";
    return contentDisposition;
}
