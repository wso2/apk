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

import runtime_domain_service.model;
import wso2/apk_common_lib as commons;
import ballerina/io;
import ballerina/crypto;
import ballerina/time;
import ballerina/uuid;
import ballerina/regex;
import ballerina/file;

public isolated function getConfigMapById(string configMapId, model:API api, commons:Organization organization) returns model:ConfigMap|commons:APKError {
    do {
        string apiNameHash = crypto:hashSha1(api.spec.apiDisplayName.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(api.spec.apiVersion.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
        lock {
            map<model:ConfigMap> & readonly readOnlyconfigMapList = configMapList.cloneReadOnly();
            foreach model:ConfigMap & readonly item in readOnlyconfigMapList {
                (map<string> & readonly) labels = item.metadata.labels ?: {};
                if (labels.hasKey(ORGANIZATION_HASH_LABEL) && labels.get(ORGANIZATION_HASH_LABEL) == organizationHash) {
                    if (labels.hasKey(API_NAME_HASH_LABEL) && labels.get(API_NAME_HASH_LABEL) == apiNameHash) {
                        if (labels.hasKey(API_VERSION_HASH_LABEL) && labels.get(API_VERSION_HASH_LABEL) == apiVersionHash) {
                            if (labels.hasKey(CONFIG_TYPE_LABEL) && labels.get(CONFIG_TYPE_LABEL) == CONFIG_TYPE_LABEL_VALUE) {
                                if (item.metadata.uid == configMapId) {
                                    return item;
                                }
                            }
                        }
                    }
                }
            }
        }
        return error("ConfigMap not found", message = string:'join("ConfigMap not found for configMap id: ", configMapId), code = 900200, description = string:'join("ConfigMap not found for configMap id: ", configMapId), statusCode = 404);
    } on fail var e {
        return error("Error while getting configMap", e, message = string:'join("Error while getting configMap for configMap id: ", configMapId), code = 900200, description = string:'join("Error while getting configMap for configMap id: ", configMapId), statusCode = 404);
    }
}

public isolated function getCertificateById(string certificateId, model:API api, commons:Organization organization) returns model:Certificate|commons:APKError {
    do {
        string apiNameHash = crypto:hashSha1(api.spec.apiDisplayName.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(api.spec.apiVersion.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
        lock {
            map<model:ConfigMap> & readonly readOnlyconfigMapList = configMapList.cloneReadOnly();
            foreach model:ConfigMap & readonly item in readOnlyconfigMapList {
                (map<string> & readonly) labels = item.metadata.labels ?: {};
                if (labels.hasKey(ORGANIZATION_HASH_LABEL) && labels.get(ORGANIZATION_HASH_LABEL) == organizationHash) {
                    if (labels.hasKey(API_NAME_HASH_LABEL) && labels.get(API_NAME_HASH_LABEL) == apiNameHash) {
                        if (labels.hasKey(API_VERSION_HASH_LABEL) && labels.get(API_VERSION_HASH_LABEL) == apiVersionHash) {
                            if (labels.hasKey(CONFIG_TYPE_LABEL) && labels.get(CONFIG_TYPE_LABEL) == CONFIG_TYPE_LABEL_VALUE) {
                                if (item.metadata.uid == certificateId) {
                                    return check convertConfigMapToCertificate(item);
                                }
                            }
                        }
                    }
                }
            }
        }
        return error("Certificate not found", message = string:'join("Certificate not found for certificate id: ", certificateId), code = 900200, description = string:'join("Certificate not found for certificate id: ", certificateId), statusCode = 404);
    } on fail var e {
        return error("Error while getting certificate", e, message = string:'join("Error while getting certificate for certificate id: ", certificateId), code = 900200, description = string:'join("Error while getting certificate for certificate id: ", certificateId), statusCode = 404);
    }
}

public isolated function getCertificatesForAPIId(model:API api, commons:Organization organization) returns model:Certificate[]|commons:APKError {
    do {
        string apiNameHash = crypto:hashSha1(api.spec.apiDisplayName.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(api.spec.apiVersion.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
        lock {
            model:Certificate[] certificates = [];

            map<model:ConfigMap> & readonly readOnlyconfigMapList = configMapList.cloneReadOnly();
            foreach model:ConfigMap & readonly item in readOnlyconfigMapList {
                (map<string> & readonly) labels = item.metadata.labels ?: {};
                if (labels.hasKey(ORGANIZATION_HASH_LABEL) && labels.get(ORGANIZATION_HASH_LABEL) == organizationHash) {
                    if (labels.hasKey(API_NAME_HASH_LABEL) && labels.get(API_NAME_HASH_LABEL) == apiNameHash) {
                        if (labels.hasKey(API_VERSION_HASH_LABEL) && labels.get(API_VERSION_HASH_LABEL) == apiVersionHash) {
                            if (labels.hasKey(CONFIG_TYPE_LABEL) && labels.get(CONFIG_TYPE_LABEL) == CONFIG_TYPE_LABEL_VALUE) {
                                certificates.push(check convertConfigMapToCertificate(item));
                            }
                        }
                    }
                }
            }
            return certificates.clone();
        }
    } on fail var e {
        return error("Error while getting certificates", e, message = "Error while getting certificates", code = 900200, description = "Error while getting certificates", statusCode = 404);
    }
}

isolated function convertConfigMapToCertificate(model:ConfigMap configMap) returns model:Certificate|error {
    map<string> annotations = configMap.metadata.annotations ?: {};
    string hosts = annotations.hasKey(CERTIFICATE_HOSTS) ? annotations.get(CERTIFICATE_HOSTS) : "";
    string serialNumber = annotations.hasKey(CERTIFICATE_SERIAL_NUMBER) ? annotations.get(CERTIFICATE_SERIAL_NUMBER) : "";
    string issuer = annotations.hasKey(CERTIFICATE_ISSUER) ? annotations.get(CERTIFICATE_ISSUER) : "";
    string subject = annotations.hasKey(CERTIFICATE_SUBJECT) ? annotations.get(CERTIFICATE_SUBJECT) : "";
    string notBefore = annotations.hasKey(CERTIFICATE_NOT_BEFORE) ? annotations.get(CERTIFICATE_NOT_BEFORE) : "";
    string notAfter = annotations.hasKey(CERTIFICATE_NOT_AFTER) ? annotations.get(CERTIFICATE_NOT_AFTER) : "";
    string certificateVersion = annotations.hasKey(CERTIFICATE_VERSION_NUMBER) ? annotations.get(CERTIFICATE_VERSION_NUMBER) : "";
    map<string> data = configMap.data ?: {};
    string certificateContent = data.hasKey(CERTIFICATE_KEY_CONFIG_MAP) ? data.get(CERTIFICATE_KEY_CONFIG_MAP) : "";
    boolean active = false;
    [int, decimal] currentTime = time:utcNow();
    if currentTime[0] >= check int:fromString(notBefore) && currentTime[0] <= check int:fromString(notAfter) {
        active = true;
    }
    return {
        'version: certificateVersion,
        certificateId: <string>configMap.metadata.uid,
        notAfter: notAfter,
        certificateContent: certificateContent,
        srtialNumber: serialNumber,
        subject: subject,
        issuer: issuer,
        notBefore: notBefore,
        hostname: hosts,
        active: active
    };
}

isolated function validateCertificateExpiry(EndpointCertificateRequest endpointCertificateRequest) returns [crypto:Certificate?, boolean]|error {
    string tmpDir = check file:createTempDir();
    string certPath = tmpDir + file:pathSeparator + endpointCertificateRequest.fileName;
    _ = check io:fileWriteBytes(certPath, endpointCertificateRequest.certificateFileContent);
    crypto:PublicKey decodedCertFile = check crypto:decodeRsaPublicKeyFromCertFile(certPath);
    crypto:Certificate? certificate = decodedCertFile.certificate;
    if certificate is crypto:Certificate {
        [int, decimal] & readonly notAfter = certificate.notAfter;
        [int, decimal] & readonly notBefore = certificate.notBefore;
        [int, decimal] currentTime = time:utcNow();
        if currentTime[0] >= notBefore[0] && currentTime[0] <= notAfter[0] {
            return [certificate, true];
        }
    }
    return [certificate, false];
}

isolated function createCertificateConfigMapEntry(string apiName, string apiVersion, EndpointCertificateRequest endpointCertificateRequest, crypto:Certificate certificate, commons:Organization organization) returns model:ConfigMap|error {
    byte[] certificateFileContent = endpointCertificateRequest.certificateFileContent;
    string content = check string:fromBytes(certificateFileContent);
    model:ConfigMap configMap = {
        metadata: {
            name: uuid:createType1AsString(),
            namespace: currentNameSpace,
            labels: getLabelsForCertificates(apiName, apiVersion, organization),
            annotations: {
                [CERTIFICATE_SERIAL_NUMBER] : certificate.serial.toString(),
                [CERTIFICATE_ISSUER] : certificate.issuer,
                [CERTIFICATE_SUBJECT] : certificate.subject,
                [CERTIFICATE_NOT_BEFORE] : certificate.notBefore[0].toString(),
                [CERTIFICATE_NOT_AFTER] : certificate.notAfter[0].toString(),
                [CERTIFICATE_VERSION_NUMBER] : certificate.version0.toString()
            }
        },
        data: {
            [CERTIFICATE_KEY_CONFIG_MAP] : content
        }
    };
    if endpointCertificateRequest.host is string {
        configMap.metadata.annotations[CERTIFICATE_HOSTS] = <string>endpointCertificateRequest.host;
    }
    return configMap;
}

public isolated function getLabelsForCertificates(string apiName, string apiVersion, commons:Organization organization) returns map<string> {
    string apiNameHash = crypto:hashSha1(apiName.toBytes()).toBase16();
    string apiVersionHash = crypto:hashSha1(apiVersion.toBytes()).toBase16();
    string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
    map<string> labels = {
        [API_NAME_HASH_LABEL] : apiNameHash,
        [API_VERSION_HASH_LABEL] : apiVersionHash,
        [ORGANIZATION_HASH_LABEL] : organizationHash,
        [MANAGED_BY_HASH_LABEL] : MANAGED_BY_HASH_LABEL_VALUE,
        [CONFIG_TYPE_LABEL] : CONFIG_TYPE_LABEL_VALUE
    };
    return labels;
}

public isolated function getLabels(APKConf apkConf, commons:Organization organization) returns map<string> {
    string apiNameHash = crypto:hashSha1(apkConf.name.toBytes()).toBase16();
    string apiVersionHash = crypto:hashSha1(apkConf.'version.toBytes()).toBase16();
    string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
    map<string> labels = {
        [API_NAME_HASH_LABEL] : apiNameHash,
        [API_VERSION_HASH_LABEL] : apiVersionHash,
        [ORGANIZATION_HASH_LABEL] : organizationHash,
        [MANAGED_BY_HASH_LABEL] : MANAGED_BY_HASH_LABEL_VALUE
    };
    return labels;
}

public isolated function getConfigMapNameByHostname(model:APIArtifact apiArtifact, APKConf apkConf, commons:Organization organization, Endpoint endpointConfig) returns string|commons:APKError? {
    do {
        string apiNameHash = crypto:hashSha1(apkConf.name.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(apkConf.'version.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
        lock {
            map<model:ConfigMap> endpointCertificates = apiArtifact.endpointCertificates;
            if endpointConfig.certification is string {
                if apiArtifact.certificateMap.hasKey(<string>endpointConfig.certification) {
                    return apiArtifact.certificateMap[<string>endpointConfig.certification];
                }
            }
            foreach model:ConfigMap item in endpointCertificates {
                map<string> labels = item.metadata.labels ?: {};
                if (labels.hasKey(ORGANIZATION_HASH_LABEL) && labels.get(ORGANIZATION_HASH_LABEL) == organizationHash) {
                    if (labels.hasKey(API_NAME_HASH_LABEL) && labels.get(API_NAME_HASH_LABEL) == apiNameHash) {
                        if (labels.hasKey(API_VERSION_HASH_LABEL) && labels.get(API_VERSION_HASH_LABEL) == apiVersionHash) {
                            if (labels.hasKey(CONFIG_TYPE_LABEL) && labels.get(CONFIG_TYPE_LABEL) == CONFIG_TYPE_LABEL_VALUE) {
                                map<string> annotations = item.metadata.annotations ?: {};
                                string hosts = annotations.hasKey(CERTIFICATE_HOSTS) ? annotations.get(CERTIFICATE_HOSTS) : "";
                                if regex:matches(endpointConfig.endpointURL, hosts) {
                                    return item.metadata.name;
                                }
                            }
                        }
                    }
                }
            }
        }
    } on fail var e {
        return error("Error while getting certificates", e, message = "Error while getting certificates", code = 900200, description = "Error while getting certificates", statusCode = 404);
    }
    return;
}

public isolated function getConfigMapsForAPICertificate(string apiName, string apiVersion, commons:Organization organization) returns model:ConfigMap[]|commons:APKError {
    do {
        string apiNameHash = crypto:hashSha1(apiName.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(apiVersion.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.uuid.toBytes()).toBase16();
        lock {
            model:ConfigMap[] configMaps = [];
            map<model:ConfigMap> & readonly readOnlyconfigMapList = configMapList.cloneReadOnly();
            foreach model:ConfigMap & readonly item in readOnlyconfigMapList {
                (map<string> & readonly) labels = item.metadata.labels ?: {};
                if (labels.hasKey(ORGANIZATION_HASH_LABEL) && labels.get(ORGANIZATION_HASH_LABEL) == organizationHash) {
                    if (labels.hasKey(API_NAME_HASH_LABEL) && labels.get(API_NAME_HASH_LABEL) == apiNameHash) {
                        if (labels.hasKey(API_VERSION_HASH_LABEL) && labels.get(API_VERSION_HASH_LABEL) == apiVersionHash) {
                            if (labels.hasKey(CONFIG_TYPE_LABEL) && labels.get(CONFIG_TYPE_LABEL) == CONFIG_TYPE_LABEL_VALUE) {
                                configMaps.push(item);
                            }
                        }
                    }
                }
            }
            return configMaps.cloneReadOnly();
        }
    } on fail var e {
        return error("Error while getting certificates", e, message = "Error while getting certificates", code = 900200, description = "Error while getting certificates", statusCode = 404);
    }
}
