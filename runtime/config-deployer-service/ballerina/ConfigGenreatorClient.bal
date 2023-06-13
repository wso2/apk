import apk_common_lib as commons;
import ballerina/http;
import config_deployer_service.org.wso2.apk.config.api as runtimeapi;
import config_deployer_service.org.wso2.apk.config.model as runtimeModels;
import config_deployer_service.java.util as utilapis;
import config_deployer_service.org.wso2.apk.config as runtimeUtil;
import ballerina/mime;
import ballerina/jballerina.java;
import config_deployer_service.model;
import ballerina/log;
import ballerina/io;
import ballerina/file;
import ballerina/uuid;

public class ConfigGeneratorClient {

    public isolated function getGeneratedAPKConf(http:Request request) returns OkAnydata|commons:APKError|BadRequestError {
        do {
            DefinitionBody definitionBody = check self.prepareDefinitionBodyFromRequest(request);
            runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException? validateAndRetrieveDefinitionResult = ();
            if definitionBody.url is string {
                validateAndRetrieveDefinitionResult = check self.validateAndRetrieveDefinition(<string>definitionBody.'apiType, definitionBody.url, (), ());
            } else if definitionBody.definition is record {|byte[] fileContent; string fileName; anydata...;|} {
                record {|byte[] fileContent; string fileName; anydata...;|} definition = <record {|byte[] fileContent; string fileName; anydata...;|}>definitionBody.definition;
                validateAndRetrieveDefinitionResult = check self.validateAndRetrieveDefinition(<string>definitionBody.'apiType, (), <byte[]>definition.fileContent, <string>definition.fileName);
            }
            if validateAndRetrieveDefinitionResult is runtimeapi:APIDefinitionValidationResponse {
                if validateAndRetrieveDefinitionResult.isValid() {
                    runtimeapi:APIDefinition parser = validateAndRetrieveDefinitionResult.getParser();
                    runtimeModels:API apiFromDefinition = check parser.getAPIFromDefinition(validateAndRetrieveDefinitionResult.getContent());
                    APIClient apiclient = new ();
                    APKConf generatedAPKConf =check apiclient.fromAPIModelToAPKConf(apiFromDefinition);
                    string|() apkConfYaml = check commons:newYamlUtil1().fromJsonStringToYaml(generatedAPKConf.toJsonString());
                    OkAnydata response = {body: apkConfYaml, mediaType: "application/yaml"};
                    return response;
                } else {
                    // Error definition.
                    ErrorListItem[] errorItems = [];
                    utilapis:ArrayList errorItemsResult = validateAndRetrieveDefinitionResult.getErrorItems();
                    foreach int i in 0 ... errorItemsResult.size() - 1 {
                        runtimeapi:ErrorItem errorItem = check java:cast(errorItemsResult.get(i));
                        ErrorListItem errorListItem = {code: errorItem.getErrorCode().toString(), message: <string>errorItem.getErrorMessage(), description: errorItem.getErrorDescription()};
                        errorItems.push(errorListItem);
                    }
                    BadRequestError badRequest = {body: {code: 90091, message: "Invalid API Definition", 'error: errorItems}};
                    return badRequest;
                }
            } else if validateAndRetrieveDefinitionResult is runtimeapi:APIManagementException {
                return e909022("Error occured while validating the definition", validateAndRetrieveDefinitionResult.cause());
            } else {
                return e909022("Error occured while validating the definition", ());
            }
        } on fail var e {
            return e909022("Internal error occured while creating APK conf", e);
        }
    }
    private isolated function prepareDefinitionBodyFromRequest(http:Request request) returns DefinitionBody|error {
        DefinitionBody definitionBody = {};
        mime:Entity[] payloadParts = check request.getBodyParts();
        foreach mime:Entity payloadPart in payloadParts {
            mime:ContentDisposition contentDisposition = payloadPart.getContentDisposition();
            string fieldName = contentDisposition.name;
            if fieldName == "definition" {
                definitionBody.definition = {fileName: contentDisposition.fileName, fileContent: check payloadPart.getByteArray()};
            }
            if fieldName == "url" {
                definitionBody.url = check payloadPart.getText();
            }
            if fieldName == "apiType" {
                definitionBody.apiType = check payloadPart.getText();
            }
        }
        return definitionBody;
    }
    private isolated function validateAndRetrieveDefinition(string 'type, string? url, byte[]? content, string? fileName) returns runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|error|commons:APKError {
        runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|error validationResponse;
        boolean typeAvailable = 'type.length() > 0;
        string[] ALLOWED_API_DEFINITION_TYPES = ["REST", "GRAPHQL", "ASYNC"];
        if !typeAvailable {
            return e909005("type");
        }
        if (ALLOWED_API_DEFINITION_TYPES.indexOf('type) is ()) {
            return e909006();
        }
        if url is string {
            string|error retrieveDefinitionFromUrlResult = self.retrieveDefinitionFromUrl(url);
            if retrieveDefinitionFromUrlResult is string {
                validationResponse = runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition('type, [], retrieveDefinitionFromUrlResult, fileName ?: "", true);
            } else {
                log:printError("Error occured while retrieving definition from url", retrieveDefinitionFromUrlResult);
                return e909044();
            }
        } else if fileName is string && content is byte[] {
            validationResponse = runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition('type, <byte[]>content, "", <string>fileName, true);
        } else {
            return e909008();
        }
        return validationResponse;
    }
    private isolated function retrieveDefinitionFromUrl(string url) returns string|error {
        string domain = getDomain(url);
        string path = getPath(url);
        if domain.length() > 0 {
            http:Client httpClient = check new (domain);
            http:Response response = check httpClient->get(path, targetType = http:Response);
            return response.getTextPayload();
        } else {
            return error("invalid url " + url);
        }
    }
    public isolated function getGeneratedK8sResources(http:Request request) returns http:Response|BadRequestError|InternalServerErrorError|commons:APKError {
        GenerateK8sResourcesBody body = {};
        do {
            mime:Entity[] payload = check request.getBodyParts();
            foreach mime:Entity payLoadPart in payload {
                mime:ContentDisposition contentDisposition = payLoadPart.getContentDisposition();
                string fieldName = contentDisposition.name;
                if fieldName == "apkConfiguration" {
                    body.apkConfiguration = {fileName: contentDisposition.fileName, fileContent: check payLoadPart.getByteArray()};
                }
                if fieldName == "definitionFile" {
                    body.definitionFile = {fileName: contentDisposition.fileName, fileContent: check payLoadPart.getByteArray()};
                }
                if fieldName == "apiType" {
                    body.apiType = check payLoadPart.getText();
                }
            }
            APIClient apiclient = new ();
            model:APIArtifact apiArtifact = check apiclient.prepareArtifact(body.apkConfiguration, body.definitionFile);
            [string, string] zipName = check self.zipAPIArtifact(apiArtifact.uniqueId, apiArtifact);
            http:Response response = new;
            response.setFileAsPayload(zipName[1]);
            response.addHeader("Content-Disposition", "attachment; filename=" + zipName[0]);
            return response;
        } on fail var e {
            return e909052(e);
        }
    }
    private isolated function zipAPIArtifact(string apiId, model:APIArtifact apiArtifact) returns [string, string]|error {
        string zipDir = check file:createTempDir(uuid:createType1AsString());
        model:API? k8sAPI = apiArtifact.api;
        if k8sAPI is model:API {
            string yamlString = check self.convertJsonToYaml(k8sAPI.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(k8sAPI.metadata.name, ".yaml"), zipDir);
        }
        model:ConfigMap? definition = apiArtifact.definition;
        if definition is model:ConfigMap {
            string yamlString = check self.convertJsonToYaml(definition.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(definition.metadata.name, ".yaml"), zipDir);
        }
        foreach model:Authentication authenticationCr in apiArtifact.authenticationMap {
            string yamlString = check self.convertJsonToYaml(authenticationCr.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(authenticationCr.metadata.name, ".yaml"), zipDir);
        }
        foreach model:Httproute httpRoute in apiArtifact.productionRoute {
            string yamlString = check self.convertJsonToYaml(httpRoute.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(httpRoute.metadata.name, ".yaml"), zipDir);
        }
        foreach model:Httproute httpRoute in apiArtifact.sandboxRoute {
            string yamlString = check self.convertJsonToYaml(httpRoute.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(httpRoute.metadata.name, ".yaml"), zipDir);
        }
        foreach model:Backend backend in apiArtifact.backendServices {
            string yamlString = check self.convertJsonToYaml(backend.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(backend.metadata.name, ".yaml"), zipDir);
        }
        foreach model:Scope scope in apiArtifact.scopes {
            string yamlString = check self.convertJsonToYaml(scope.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(scope.metadata.name, ".yaml"), zipDir);
        }
        foreach model:RateLimitPolicy rateLimitPolicy in apiArtifact.rateLimitPolicies {
            string yamlString = check self.convertJsonToYaml(rateLimitPolicy.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(rateLimitPolicy.metadata.name, ".yaml"), zipDir);
        }
        foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
            string yamlString = check self.convertJsonToYaml(apiPolicy.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(apiPolicy.metadata.name, ".yaml"), zipDir);
        }
        foreach model:InterceptorService interceptorService in apiArtifact.interceptorServices {
            string yamlString = check self.convertJsonToYaml(interceptorService.toJsonString());
            _ = check self.storeFile(yamlString, string:concat(interceptorService.metadata.name, ".yaml"), zipDir);
        }
        string zipfileName = string:concat(apiArtifact.name, "-", apiArtifact.'version);
        [string, string] zipName = check self.zipDirectory(zipfileName, zipDir);
        return zipName;
    }

    private isolated function convertJsonToYaml(string jsonString) returns string|error {
        commons:YamlUtil yamlUtil = commons:newYamlUtil1();
        string|() convertedYaml = check yamlUtil.fromJsonStringToYaml(jsonString);
        if convertedYaml is string {
            return convertedYaml;
        }
        return e909022("Error while converting json to yaml", convertedYaml);
    }
    private isolated function storeFile(string jsonString, string fileName, string? directroy = ()) returns error? {
        commons:YamlUtil yamlUtil = commons:newYamlUtil1();
        string|() convertedYaml = check yamlUtil.fromJsonStringToYaml(jsonString);
        string fullPath = directroy ?: "";
        if convertedYaml is string {
            fullPath = fullPath + file:pathSeparator + fileName + ".yaml";
            _ = check io:fileWriteString(fullPath, convertedYaml);
        }
    }

    private isolated function zipDirectory(string zipfileName, string directoryPath) returns [string, string]|error {
        string zipName = zipfileName + ZIP_FILE_EXTENSTION;
        string zipPath = directoryPath + ZIP_FILE_EXTENSTION;
        _ = check commons:ZIPUtils_zipDir(directoryPath, zipPath);
        return [zipName, zipPath];
    }

}
