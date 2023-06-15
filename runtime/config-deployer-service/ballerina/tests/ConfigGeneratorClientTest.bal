import ballerina/http;
import ballerina/test;
import apk_common_lib;
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
        test:assertEquals(response.body.toString(), content, "APK conf content mismatched");
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
        BadRequestError badRequest = {body: {code: 90091, message: "API Type need to specified"}};
        test:assertEquals(response.status.code, http:STATUS_BAD_REQUEST, "Status code mismatched");
        test:assertEquals(response.body, badRequest.body);
    }
}

@test:Config {}
public function testGetGeneratedAPKConfFromUrlAndDefinitionInvalidType() returns error? {
    ConfigGeneratorClient configGeneratorClient = new ();
    http:Request request = new;
    mime:Entity[] bodyParts = [];
    mime:Entity bodyPart = new;
    bodyPart.setFileAsEntityBody("./tests/resources/api.yaml", APPLICATION_YAML_MEDIA_TYPE);
    bodyPart.setContentDisposition(getContentDispositionForFormData("definition"));
    bodyParts.push(bodyPart);
    request.setBodyParts(bodyParts, contentType = mime:MULTIPART_FORM_DATA);

    OkAnydata|apk_common_lib:APKError|BadRequestError response = configGeneratorClient.getGeneratedAPKConf(request);
    if response is BadRequestError {
        BadRequestError badRequest = {body: {code: 90091, message: "API Type need to specified"}};
        test:assertEquals(response.status.code, http:STATUS_BAD_REQUEST, "Status code mismatched");
        test:assertEquals(response.body, badRequest.body);
    }
}

function getContentDispositionForFormData(string partName) returns (mime:ContentDisposition) {
    mime:ContentDisposition contentDisposition = new;
    contentDisposition.name = partName;
    contentDisposition.disposition = "form-data";
    return contentDisposition;
}
