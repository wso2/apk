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
import wso2/apk_common_lib as commons;
import ballerina/regex;
import ballerina/log;

isolated function getAPIByAPIId(string apiId) returns API|NotFoundError|commons:APKError {
    API|commons:APKError|NotFoundError api = getAPIByIdDAO(apiId);
    return api;
}

isolated function getAPIList(int 'limit, int offset, string? query, commons:Organization organization, anydata? groups) returns APIList|commons:APKError {
    string[] groupsArray = getUserGroups(groups);

    if query !is string {
        APIInfo[]|commons:APKError apis = getAPIsDAO(organization.uuid, groupsArray);
        if apis is APIInfo[] {
            APIInfo[] limitSet = [];
            if apis.length() > offset {
                foreach int i in offset ... (apis.length() - 1) {
                    if limitSet.length() < 'limit {
                        limitSet.push(apis[i]);
                    }
                }
            }
            APIList apisList = {count: limitSet.length(), list: limitSet, pagination: {total: apis.length(), 'limit: 'limit, offset: offset}};
            return apisList;
        } else {
            return apis;
        }
    } else {
        boolean hasPrefix = query.startsWith("content");
        if hasPrefix {
            int? index = query.indexOf(":");
            if index is int {
                string modifiedQuery = "%" + query.substring(index + 1) + "%";
                APIInfo[]|commons:APKError apis = getAPIsByQueryDAO(modifiedQuery, organization.uuid, groupsArray);
                if apis is APIInfo[] {
                    APIInfo[] limitSet = [];
                    if apis.length() > offset {
                        foreach int i in offset ... (apis.length() - 1) {
                            if limitSet.length() < 'limit {
                                limitSet.push(apis[i]);
                            }
                        }
                    }
                    APIList apisList = {count: limitSet.length(), list: limitSet, pagination: {total: apis.length(), 'limit: 'limit, offset: offset}};
                    return apisList;
                } else {
                    return apis;
                }
            } else {
                string message = "Invalid Content Search Text Provided. Missing :";
                commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 400);
                return e;
            }
        } else {
            string message = "Invalid Content Search Text Provided. Missing content keyword";
            commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 400);
            return e;
        }
    }
}

isolated function getAPIDefinition(string apiId) returns APIDefinition|NotFoundError|commons:APKError {
    APIDefinition|NotFoundError|commons:APKError apiDefinition = getAPIDefinitionDAO(apiId);
    return apiDefinition;
}

isolated function generateSDKImpl(string apiId, string language) returns http:Response|NotFoundError|commons:APKError {
    sdk:APIClientGenerationManager sdkClient = new sdk:APIClientGenerationManager(newSDKClient());
    string apiName;
    string apiVersion;
    API|NotFoundError api = check getAPIByAPIId(apiId);
    if api is API {
        apiName = api.name;
        apiVersion = api.'version;
        APIDefinition|NotFoundError|commons:APKError apiDefinition = getAPIDefinition(apiId);
        if apiDefinition is APIDefinition {
            string? schema = apiDefinition.schemaDefinition;
            if schema is string {
                javautil:Map|sdk:APIClientGenerationException sdkMap = sdkClient.generateSDK(language, apiName, apiVersion, schema,
                sdkConfig.groupId, sdkConfig.artifactId, sdkConfig.modelPackage, sdkConfig.apiPackage);
                if sdkMap is javautil:Map {
                    string path = readMap(sdkMap, "zipFilePath");
                    string fileName = readMap(sdkMap, "zipFileName");
                    http:Response response = new;
                    response.setFileAsPayload(path);
                    response.addHeader("Content-Disposition", "attachment; filename=" + fileName);
                    return response;
                } else {
                    commons:APKError e = error("Unable to generate SDK", message = "Unable to generate SDK", description = "Unable to generate SDK", code = 90911, statusCode = 500);
                    return e;
                }
            } else {
                string message = "Unable to retrieve schema mediation";
                commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 500);
                return e;
            }
        }
    } else if api is NotFoundError|commons:APKError {
        return api;
    }
    string message = "Unable to generate SDK";
    commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 500);
    return e;
}

isolated function getSDKLanguages() returns string|json|commons:APKError {
    sdk:APIClientGenerationManager sdkClient = new sdk:APIClientGenerationManager(newSDKClient());
    string? sdkLanguages = sdkClient.getSupportedSDKLanguages();
    return sdkLanguages;
}

isolated function newSDKClient() returns handle = @java:Constructor {
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager"
} external;

isolated function readMap(javautil:Map sdkMap, string key) returns string {
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

isolated function getThumbnail(string apiId) returns http:Response|NotFoundError|commons:APKError {
    API|NotFoundError|commons:APKError api = getAPIByAPIId(apiId);
    if api is API {
        int|commons:APKError thumbnailCategoryId = getResourceCategoryIdByCategoryTypeDAO(RESOURCE_TYPE_THUMBNAIL);
        if thumbnailCategoryId is int {
            Resource|NotFoundError|commons:APKError thumbnail = getResourceByResourceCategoryDAO(apiId, thumbnailCategoryId);
            if thumbnail is Resource {
                http:Response outResponse = new;
                outResponse.setBinaryPayload(thumbnail.resourceBinaryValue, thumbnail.dataType);
                return outResponse;
            } else {
                return thumbnail;
            }
        }
        return thumbnailCategoryId;
    } else if api is NotFoundError|commons:APKError {
        return api;
    }
    string message = "Unable to get the thumbnail";
    commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 500);
    return e;
}

isolated function getDocumentMetaData(string apiId, string documentId) returns Document|NotFoundError|commons:APKError {
    API|NotFoundError|commons:APKError api = getAPIByAPIId(apiId);
    if api is API {
        DocumentMetaData|NotFoundError|commons:APKError getDocumentMetaData = getDocumentByDocumentIdDAO(documentId, apiId);
        if getDocumentMetaData is DocumentMetaData {
            // Convert documentMetadata object to Document object
            Document document = {
                documentId: getDocumentMetaData.documentId,
                name: getDocumentMetaData.name,
                summary: getDocumentMetaData.summary,
                sourceType: <"INLINE"|"MARKDOWN"|"URL"|"FILE">getDocumentMetaData.sourceType,
                sourceUrl: getDocumentMetaData.sourceUrl,
                documentType: <"HOWTO"|"SAMPLES"|"PUBLIC_FORUM"|"SUPPORT_FORUM"|"API_MESSAGE_FORMAT"|"SWAGGER_DOC"|"OTHER">getDocumentMetaData.documentType,
                otherTypeName: getDocumentMetaData.otherTypeName
            };
            return document;
        } else {
            return getDocumentMetaData;
        }
    } else if api is NotFoundError|commons:APKError {
        return api;
    }
    string message = "Unable to get the Document meta data";
    commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 500);
    return e;
}


isolated function getDocumentContent(string apiId, string documentId) returns http:Response|NotFoundError|commons:APKError {
    API|NotFoundError|commons:APKError api = getAPIByAPIId(apiId);
    if api is API {
        DocumentMetaData|NotFoundError|commons:APKError getDocumentMetaData = getDocumentByDocumentIdDAO(documentId, apiId);
        if getDocumentMetaData is DocumentMetaData {
            Resource|commons:APKError getDocumentResource = getResourceByResourceIdDAO(<string>getDocumentMetaData.resourceId);
            if getDocumentResource is Resource {
                    http:Response outResponse = new;
                    outResponse.setBinaryPayload(<byte[]>getDocumentResource.resourceBinaryValue, getDocumentResource.dataType);
                    return outResponse;
            } else {
                return getDocumentResource;
            }
        } else {
            return getDocumentMetaData;
        }
    } else if api is NotFoundError|commons:APKError {
        return api;
    }
    string message = "Unable to get the Document content";
    commons:APKError e = error(message, message = message, description = message, code = 90911, statusCode = 500);
    return e;
}

isolated function getDocumentList(string apiId, int 'limit, int offset) returns DocumentList|NotFoundError|commons:APKError {
    API|NotFoundError|commons:APKError api = getAPIByAPIId(apiId);
    if api is API {
        Document[]|commons:APKError documents = getDocumentsDAO(apiId);
        if documents is Document[] {
            Document[] limitSet = [];
            if documents.length() > offset {
                foreach int i in offset ... (documents.length() - 1) {
                    if limitSet.length() < 'limit {
                        limitSet.push(documents[i]);
                    }
                }
            }
            DocumentList documentList = {count: limitSet.length(), list: limitSet, pagination: {total: documents.length(), 'limit: 'limit, offset: offset}};
            return documentList;
        } else {
            return documents;
        }
    } else if api is NotFoundError  {
        return api;
    } else {
        string message = "Internal Error occured while retrieving API for Docuements retieval";
        return error(message, message = message, description = message, code = 909001, statusCode = 500);
    }
}

isolated function getUserGroups(anydata groups) returns string[] {
    string[] groupsArray = [];
    if (groups is json[]) {
        json[] groupsArr = <json[]>groups;
        foreach json group in groupsArr {
            groupsArray.push(group.toString());
        }
    } else if (groups is string) {
        string groupsStr = <string>groups;
        string[] tmp = regex:split(groupsStr, ",");
        foreach string group in tmp {
            string trimmedGroup = group.trim();
            if (trimmedGroup != "") {
                groupsArray.push(trimmedGroup);
            }
        }
    } else {
        log:printDebug("No user groups found");
    }
    return groupsArray;
}
