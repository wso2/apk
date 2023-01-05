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

import ballerina/jballerina.java;
import devportal_service.org.wso2.apk.devportal.sdk as sdk;
import devportal_service.java.util as javautil;
import devportal_service.java.lang as javalang;
import ballerina/http;

isolated function getAPIByAPIId(string apiId, string organization) returns string?|API|error {
    string?|API|error api = getAPIByIdDAO(apiId, organization);
    return api;
}

isolated function getAPIList(int 'limit, int  offset, string? query, string organization) returns string?|APIList|error {
    API[]|error? apis = getAPIsDAO(organization);
    if apis is API[] {
        int count = apis.length();
        ApplicationList apisList = {count: count, list: apis};
        return apisList;
    } else {
        return apis;
    }
}

isolated function getAPIDefinition(string apiId, string organization) returns APIDefinition|NotFoundError|error {
    APIDefinition|NotFoundError|error apiDefinition = getAPIDefinitionDAO(apiId,organization);
    return apiDefinition;
}

function generateSDK(string apiId, string language, string org) returns http:Response|sdk:APIClientGenerationException|string?|error {
    sdk:APIClientGenerationManager sdkClient = new sdk:APIClientGenerationManager(newSDKClient());
    string apiName;
    string apiVersion;
    string?|API|error api = getAPIByAPIId(apiId,org);
    if api is API {
        apiName = api.name;
        apiVersion = api.'version;
        APIDefinition|NotFoundError|error apiDefinition = getAPIDefinition(apiId,org);
        if apiDefinition is APIDefinition {
            string? schema = apiDefinition.schemaDefinition;
            if schema is string {
                javautil:Map|sdk:APIClientGenerationException sdkMap = sdkClient.generateSDK(language,apiName,apiVersion,schema);
                if sdkMap is javautil:Map {
                    string path = readMap(sdkMap,"zipFilePath");
                    string fileName = readMap(sdkMap,"zipFileName");
                    http:Response response = new;
                    response.setFileAsPayload(path);
                    response.addHeader("Content-Disposition","attachment; filename=" + fileName);
                    return response;
                } else {
                    return sdkMap;
                }
            } else {
                return error("Unable to retrieve Schema Definition");
            }
        }
    } else {
        return api;
    }
    return;
}

isolated function getSDKLanguages() returns string|json|error {
    sdk:APIClientGenerationManager sdkClient = new sdk:APIClientGenerationManager(newSDKClient());
    string? sdkLanguages = sdkClient.getSupportedSDKLanguages();
    return sdkLanguages;
}

isolated function newSDKClient() returns handle = @java:Constructor {
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager"
} external;

function readMap(javautil:Map sdkMap, string key) returns string{
    handle keyAsJavaStr = java:fromString(key);
    javalang:Object keyAsObj = new (keyAsJavaStr);
    javalang:Object value = sdkMap.get(keyAsObj);
    // Above simplified to one line
    // javalang:Object value = sdkMap.get(new (java:fromString("zipFilePath")));
    handle valueHandle = value.jObj;
    string? valueStr = java:toString(valueHandle);
    if valueStr is string {
        return valueStr;
    } else {
        return "";
    }
}
