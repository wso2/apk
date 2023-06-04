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
import runtime_domain_service.model;
import wso2/apk_common_lib as commons;
import ballerina/http;
import ballerina/log;

isolated function convertK8sAPItoAPI(model:API api, boolean lightWeight) returns API|commons:APKError {
    API convertedModel = {
        id: getAPIUUIDFromAPI(api),
        name: api.spec.apiDisplayName,
        context: api.spec.context,
        'version: api.spec.apiVersion,
        'type: api.spec.apiType,
        createdTime: api.metadata.creationTimestamp
    };
    model:APIStatus? status = api.status;
    if status is model:APIStatus {
        convertedModel.lastUpdatedTime = status.transitionTime;
    }
    if !lightWeight {
        model:RuntimeAPI|http:ClientError internalAPI = getInternalAPI(api.metadata.name, <string>api.metadata.namespace);
        if internalAPI is model:RuntimeAPI {
            record {|anydata...;|}? endpointConfig = internalAPI.spec.endpointConfig.clone();
            if endpointConfig is record {} {
                convertedModel.endpointConfig = maskEndpointSecurityPassword(endpointConfig);
            }
            model:OperationPolicies? apiPolicies = internalAPI.spec.apiPolicies;
            if apiPolicies is model:OperationPolicies {
                convertedModel.apiPolicies = convertOperationPolicies(apiPolicies);
            }
            APIOperations[] apiOperations = [];
            model:Operations[]? operations = internalAPI.spec.operations;
            if operations is model:Operations[] {
                foreach model:Operations operation in operations {
                    apiOperations.push({
                        verb: operation.verb,
                        target: operation.target,
                        authTypeEnabled: operation.authTypeEnabled,
                        scopes: operation.scopes,
                        endpointConfig: maskEndpointSecurityPassword(operation.endpointConfig),
                        operationPolicies: convertOperationPolicies(operation.operationPolicies),
                        operationRateLimit: operation.operationRateLimit
                    });
                }
            }
            convertedModel.operations = apiOperations;
            model:RateLimit? apiRateLimit = internalAPI.spec.apiRateLimit;
            if apiRateLimit is model:RateLimit {
                convertedModel.apiRateLimit = {requestsPerUnit: apiRateLimit.requestsPerUnit, unit: apiRateLimit.unit};
            }
            string[]? securitySchemes = internalAPI.spec.securityScheme;
            if securitySchemes is string[] {
                log:printDebug("securitySchemes: " + securitySchemes.toString());
                convertedModel.securityScheme = securitySchemes;
            }
            model:ServiceInfo? serviceInfo = internalAPI.spec.serviceInfo;
            if serviceInfo is model:ServiceInfo {
                if serviceInfo.endpointSecurity is map<anydata> {
                    convertedModel.serviceInfo = {name: serviceInfo.name, namespace: serviceInfo.namespace, endpoint_security: serviceInfo.endpointSecurity};
                } else {
                    convertedModel.serviceInfo = {name: serviceInfo.name, namespace: serviceInfo.namespace};
                }
            }
        } else if internalAPI is http:ApplicationResponseError {
            if internalAPI.detail().statusCode != 404 {
                return error("Error while converting k8s API to API", internalAPI, code = 900900, message = "Internal Server Error", statusCode = 500, description = "Internal Server Error");
            }
        } else {
            return error("Error while converting k8s API to API", internalAPI, code = 900900, message = "Internal Server Error", statusCode = 500, description = "Internal Server Error");
        }
    }
    return convertedModel;
}

isolated function maskEndpointSecurityPassword(record {}? endpointConfig) returns record {}? {
    if endpointConfig is record {} {
        record {} clonedEndpointconfig = endpointConfig.clone();
        anydata|error endpointSecurity = trap clonedEndpointconfig.get("endpoint_security");
        if endpointSecurity is map<anydata> {
            if endpointSecurity.hasKey(SANDBOX_TYPE) {
                map<anydata> sandboxEndpointSecurity = <map<anydata>>endpointSecurity.get(SANDBOX_TYPE);
                anydata|error endpointSecurityType = trap sandboxEndpointSecurity.get(ENDPOINT_SECURITY_TYPE);
                if endpointSecurityType is string && endpointSecurityType == ENDPOINT_SECURITY_TYPE_BASIC_CASE {
                    if sandboxEndpointSecurity.hasKey(ENDPOINT_SECURITY_PASSWORD) {
                        sandboxEndpointSecurity[ENDPOINT_SECURITY_PASSWORD] = DEFAULT_MODIFIED_ENDPOINT_PASSWORD;
                    }

                }
            }
            if endpointSecurity.hasKey(PRODUCTION_TYPE) {
                map<anydata> productionEndpoinySecurity = <map<anydata>>endpointSecurity.get(PRODUCTION_TYPE);
                anydata|error endpointSecurityType = trap productionEndpoinySecurity.get(ENDPOINT_SECURITY_TYPE);
                if endpointSecurityType is string && endpointSecurityType == ENDPOINT_SECURITY_TYPE_BASIC_CASE {
                    if productionEndpoinySecurity.hasKey(ENDPOINT_SECURITY_PASSWORD) {
                        productionEndpoinySecurity[ENDPOINT_SECURITY_PASSWORD] = DEFAULT_MODIFIED_ENDPOINT_PASSWORD;
                    }

                }
            }
        }
        return clonedEndpointconfig;
    }
    return;
}

isolated function convertOperationPolicies(model:OperationPolicies? operation) returns APIOperationPolicies|() {
    if operation is model:OperationPolicies {
        OperationPolicy[] requestPolicies = [];
        OperationPolicy[] responsePolicies = [];
        foreach model:OperationPolicy requestPolicy in operation.request {
            OperationPolicy policy = {...requestPolicy};
            requestPolicies.push(policy);
        }
        foreach model:OperationPolicy responsePolicy in operation.response {
            OperationPolicy policy = {...responsePolicy};
            responsePolicies.push(policy);
        }
        return {request: requestPolicies, response: responsePolicies};
    } else {
        return ();
    }
}

isolated function convertPolicyModeltoPolicy(model:MediationPolicy mediationPolicy) returns MediationPolicy {
    MediationPolicy mediationPolicyData = {
        id: mediationPolicy.id,
        'type: mediationPolicy.'type,
        name: mediationPolicy.name,
        displayName: mediationPolicy.displayName,
        description: mediationPolicy.description,
        applicableFlows: mediationPolicy.applicableFlows,
        supportedApiTypes: mediationPolicy.supportedApiTypes,
        policyAttributes: mediationPolicy.policyAttributes

    };
    return mediationPolicyData;
}

public isolated function convertAPIListToAPIInfoList(API[] apiList) returns APIInfo[] {
    APIInfo[] apiInfoList = [];
    foreach API api in apiList {
        APIInfo apiInfo = {
            id: api.id,
            name: api.name,
            context: api.context,
            'version: api.'version,
            'type: api.'type,
            createdTime: api.createdTime
        };
        apiInfoList.push(apiInfo);
    }
    return apiInfoList;
}

public isolated function fromAPIToAPKConf(API api) returns APKConf|error {
    APKConf apkConf = {
        id: api.id,
        name: api.name,
        context: api.context,
        version: api.'version,
        'type: api.'type,
        apiPolicies: api.apiPolicies,
        apiRateLimit: api.apiRateLimit,
        securityScheme: api.securityScheme,
        serviceInfo: convertApiServiceInfo(api.serviceInfo)
    };
    if api.operations is APIOperations[] {
        apkConf.operations = check convetAPIOperations(api.operations);
    }
    if api.endpointConfig is record {} {
        apkConf.endpointConfig = check convertEndpointConfig(api.endpointConfig);
    }
    return apkConf;
}

isolated function convertApiServiceInfo(API_serviceInfo? serviceInfo) returns APKServiceInfo? {
    if serviceInfo is API_serviceInfo {
        APKServiceInfo apkServiceInfo = {
            name: serviceInfo.name,
            namespace: serviceInfo.namespace
        };
        if serviceInfo.endpoint_security is record {} {
            apkServiceInfo.endpointSecurity = convertEndpointSecurity(<record {}>serviceInfo.endpoint_security, PRODUCTION_TYPE);
        }
        return apkServiceInfo;
    }
    return;
}

isolated function convetAPIOperations(APIOperations[]? apiOperations) returns APKOperations[]|error {
    APKOperations[] apkOperations = [];
    if apiOperations is APIOperations[] {
        foreach APIOperations apiOperation in apiOperations {
            APKOperations apkOperation = {
                authTypeEnabled: apiOperation.authTypeEnabled,
                target: apiOperation.target,
                verb: apiOperation.verb,
                scopes: apiOperation.scopes,
                operationPolicies: apiOperation.operationPolicies
            };
            if apiOperation.endpointConfig is record {} {
                apkOperation.endpointConfig = check convertEndpointConfig(apiOperation.endpointConfig);
            }
            if apiOperation.operationRateLimit is APIRateLimit {
                apkOperation.operationRateLimit = {
                    requestsPerUnit: <int>apiOperation.operationRateLimit?.requestsPerUnit,
                    unit: <string>apiOperation.operationRateLimit?.unit
                };
            }
            apkOperations.push(apkOperation);
        }
    }
    return apkOperations;
}

isolated function convertEndpointConfig(record {}? apiEndpointConfig) returns EndpointConfig|error {
    EndpointConfig endpointConfig = {};
    Endpoint? production = ();
    Endpoint? sandbox = ();
    if apiEndpointConfig is record {} {
        anydata|error sandboxEndpointConfig = trap apiEndpointConfig.get("sandbox_endpoints");
        anydata|error productionEndpointConfig = trap apiEndpointConfig.get("production_endpoints");
        anydata|error endpoint_security = trap apiEndpointConfig.get("endpoint_security");
        if sandboxEndpointConfig is map<anydata> {
            if sandboxEndpointConfig.hasKey("url") {
                anydata url = sandboxEndpointConfig.get("url");
                EndpointSecurity? endpointSecurity = ();
                if endpoint_security is record {} {
                    endpointSecurity = convertEndpointSecurity(endpoint_security, SANDBOX_TYPE);
                }
                sandbox = {
                    endpointURL: <string>url,
                    endpointSecurity: endpointSecurity
                };
            } else {
                return e909013();
            }
        }
        if productionEndpointConfig is map<anydata> {
            if productionEndpointConfig.hasKey("url") {
                anydata url = productionEndpointConfig.get("url");
                EndpointSecurity? endpointSecurity = ();
                if endpoint_security is record {} {
                    endpointSecurity = convertEndpointSecurity(endpoint_security, PRODUCTION_TYPE);
                }
                production = {
                    endpointURL: <string>url,
                    endpointSecurity: endpointSecurity
                };
            } else {
                return e909014();
            }
        }
        endpointConfig = {
            sandbox: sandbox,
            production: production
        };
    }
    return endpointConfig;
}

isolated function convertEndpointSecurity(record {} endpoint_security, string endpointType) returns EndpointSecurity? {
    if endpoint_security.hasKey(endpointType) {
        EndpointSecurity? endpointSecurity = ();
        map<anydata> endpoinySecurityData = <map<anydata>>endpoint_security.get(endpointType);
        map<string> securityProperties = {
            [ENDPOINT_BASIC_USER_NAME] : endpoinySecurityData.hasKey(ENDPOINT_BASIC_USER_NAME) ? <string>endpoinySecurityData.get(ENDPOINT_BASIC_USER_NAME) : "",
            [ENDPOINT_BASIC_PASSWORD] : endpoinySecurityData.hasKey(ENDPOINT_BASIC_PASSWORD) ? <string>endpoinySecurityData.get(ENDPOINT_BASIC_PASSWORD) : ""
        };
        if endpoinySecurityData.hasKey(ENDPOINT_BASIC_SECRET_REF) {
            securityProperties[ENDPOINT_BASIC_SECRET_REF] = <string>endpoinySecurityData.get(ENDPOINT_BASIC_SECRET_REF);
        }

        endpointSecurity = {
            enable: endpoinySecurityData.hasKey("enabled") ? <boolean>endpoinySecurityData.get("enabled") : false,
            securityType: endpoinySecurityData.hasKey("type") ? <string>endpoinySecurityData.get("type") : "",
            securityProperties: securityProperties
        };
        return endpointSecurity;
    }
    return;
}
