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
import wso2/apk_common_lib as commons;
import ballerina/uuid;
import ballerina/http;
import ballerina/constraint;
import ballerina/cache;

public isolated class ServiceBaseOrgResolver {
    *commons:OrganizationResolver;
    private final http:Client httpClient;
    private final map<string>? headers;
    private final cache:Cache orgCache;
    private final cache:Cache orgnizationClaimValueCache = new ({
        // The maximum size of the cache is 10.
        capacity: 1000,
        // The eviction factor is set to 0.2, which means at the
        // time of eviction 10*0.2=2 entries get removed from the cache.
        evictionFactor: 0.2,
        // The default max age of the cache entry is set to 2 seconds.
        defaultMaxAge: 600,
        // The cache cleanup task runs every 3 seconds and clears all
        // the expired entries.
        cleanupInterval: 60
    });
    public isolated function init(string serviceBaseURL, map<string>? headers, commons:KeyStore? certificate, boolean enableAuth) returns error? {
        self.orgCache = new ({
            // The maximum size of the cache is 10.
            capacity: 1000,
            // The eviction factor is set to 0.2, which means at the
            // time of eviction 10*0.2=2 entries get removed from the cache.
            evictionFactor: 0.2,
            // The default max age of the cache entry is set to 2 seconds.
            defaultMaxAge: 600,
            // The cache cleanup task runs every 3 seconds and clears all
            // the expired entries.
            cleanupInterval: 60
        });
        TokenIssuerConfiguration issuerConfiguration = runtimeConfiguration.tokenIssuerConfiguration;
        commons:KeyStore & readonly signingCert = runtimeConfiguration.keyStores.signing;
        self.headers = headers.cloneReadOnly();
        boolean secured = serviceBaseURL.startsWith("https:");
        self.httpClient = check new (serviceBaseURL,
        auth = enableAuth ? {
                username: "runtime-domain-service",
                issuer: issuerConfiguration.issuer,
                keyId: issuerConfiguration.keyId,
                jwtId: uuid:createType1AsString(),
                expTime: issuerConfiguration.expTime,
                signatureConfig: {
                    config: {
                        keyFile: <string>signingCert.keyFilePath
                    }
                }
            } : {},
        secureSocket = secured ? (certificate is commons:KeyStore ? {cert: certificate.certFilePath, verifyHostName: runtimeConfiguration.controlPlane.enableHostNameVerification} : ()) : {});
    }

    public isolated function retrieveOrganizationByName(string organizationName) returns commons:Organization|commons:APKError|() {
        do {

            if self.orgCache.hasKey(organizationName) {
                return check self.orgCache.get(organizationName).ensureType(commons:Organization);
            }
            lock {
                if self.orgCache.hasKey(organizationName) {
                    return check self.orgCache.get(organizationName).ensureType(commons:Organization);
                }
                OrganizationList|http:ClientError getOrgnizationByName = self.httpClient->get("/organizations?name=" + organizationName, targetType = OrganizationList, headers = self.headers);
                if getOrgnizationByName is OrganizationList {
                    Organization[]? organizations = getOrgnizationByName.list;
                    if organizations is Organization[] && organizations.length() > 0 {
                        Organization organization = organizations[0];
                        commons:Organization org = {
                            uuid: <string>organization.id,
                            name: organization.name,
                            displayName: organization.displayName,
                            organizationClaimValue: <string>organization.organizationClaimValue,
                            enabled: organization.enabled,
                            serviceListingNamespaces: organization.serviceNamespaces ?: ["*"]
                        };
                        check self.orgCache.put(organizationName, org);
                        check self.orgnizationClaimValueCache.put(org.organizationClaimValue, org.name);
                        return org.clone();
                    }
                    return;
                } else {
                    commons:APKError apkError = error("Error while retrieving organization by name", getOrgnizationByName, code = 900900, message = "Error while retrieving organization by name", description = "Error while retrieving organization by name", statusCode = 500);
                    return apkError;
                }
            }
        } on fail var e {
            commons:APKError apkError = error("Error while retrieving organization by name", e, code = 900900, message = "Error while retrieving organization by name", description = "Error while retrieving organization by name", statusCode = 500);
            return apkError;
        }
    }

    public isolated function retrieveOrganizationFromIDPClaimValue(map<anydata> claims, string organizationClaim) returns commons:Organization|commons:APKError|() {
        do {
            if self.orgnizationClaimValueCache.hasKey(organizationClaim) {
                string orgName = check self.orgnizationClaimValueCache.get(organizationClaim).ensureType(string);
                if self.orgCache.hasKey(orgName) {
                    return check self.orgCache.get(orgName).ensureType(commons:Organization);
                }
            }
            lock {
                if self.orgnizationClaimValueCache.hasKey(organizationClaim) {
                    string orgName = check self.orgnizationClaimValueCache.get(organizationClaim).ensureType(string);
                    if self.orgCache.hasKey(orgName) {
                        return check self.orgCache.get(orgName).ensureType(commons:Organization);
                    }
                }
                OrganizationList|http:ClientError getOrgnizationByName = self.httpClient->get("/organizations?organizationClaimValue=" + organizationClaim, targetType = OrganizationList, headers = self.headers);
                if getOrgnizationByName is OrganizationList {
                    Organization[]? organizations = getOrgnizationByName.list;
                    if organizations is Organization[] && organizations.length() > 0 {
                        Organization & readonly organization = organizations[0].cloneReadOnly();
                        return {
                            uuid: <string>organization.id,
                            name: organization.name,
                            displayName: organization.displayName,
                            organizationClaimValue: <string>organization.organizationClaimValue,
                            enabled: organization.enabled,
                            serviceListingNamespaces: organization.serviceNamespaces ?: ["*"]
                        };
                    }
                    return;
                } else {
                    commons:APKError apkError = error("Error while retrieving organization by name", getOrgnizationByName, code = 900900, message = "Error while retrieving organization by name", description = "Error while retrieving organization by name", statusCode = 500);
                    return apkError;
                }
            }
        } on fail var e {
            commons:APKError apkError = error("Error while retrieving organization by name", e, code = 900900, message = "Error while retrieving organization by name", description = "Error while retrieving organization by name", statusCode = 500);
            return apkError;
        }
    }
}

public type OrganizationList record {
    # Number of Organization returned.
    int count?;
    Organization[] list?;
};

public type Organization record {
    string id?;
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    @constraint:String {maxLength: 255, minLength: 1}
    string displayName;
    @constraint:String {maxLength: 255, minLength: 1}
    string organizationClaimValue?;
    boolean enabled = true;
    string[] serviceNamespaces?;
};
