/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 LLC. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.commons.dto;

import java.security.PrivateKey;
import java.security.cert.Certificate;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

/**
 * Holds configs related to jwt generation.
 */
public class JWTConfigurationDto {

    private boolean enabled = false;
    private String jwtHeader = "X-JWT-Assertion";
    private String consumerDialectUri = "http://wso2.org/claims";
    private String signatureAlgorithm = "SHA256withRSA";
    private String encoding = "Base64";
    private String gatewayJWTGeneratorImpl = "org.wso2.apk.enforcer.commons.jwtgenerator.APIMgtGatewayJWTGeneratorImpl";
    private Map<String, TokenIssuerDto> tokenIssuerDtoMap = new HashMap();
    private Set<String> jwtExcludedClaims = new HashSet<>();
    private Certificate publicCert;
    private PrivateKey privateKey;
    private long ttl;
    private Map<String, ClaimValueDTO> customClaims = new HashMap<>();

    private boolean useKid;

    public String getKidValue() {

        return kidValue;
    }

    public void setKidValue(String kidValue) {

        this.kidValue = kidValue;
    }

    private String kidValue;
    public JWTConfigurationDto() {

    }

    public JWTConfigurationDto(final boolean enabled, final String jwtHeader,
            final String signatureAlgorithm, final String encoding,
            final Certificate publicCert, final PrivateKey privateKey, final long ttl) {
        this.enabled = enabled;
        this.jwtHeader = jwtHeader;
        this.signatureAlgorithm = signatureAlgorithm;
        this.encoding = encoding;
        this.publicCert = publicCert;
        this.privateKey = privateKey;
        this.ttl = ttl;
    }

    public boolean useKid() {
        return useKid;
    }

    public void setUseKid(boolean useKid) {
        this.useKid = useKid;
    }

    public boolean isEnabled() {

        return enabled;
    }

    public void setEnabled(boolean enabled) {

        this.enabled = enabled;
    }

    public String getJwtHeader() {

        return jwtHeader;
    }

    public void setJwtHeader(String jwtHeader) {

        this.jwtHeader = jwtHeader;
    }

    public String getConsumerDialectUri() {

        return consumerDialectUri;
    }

    public void setConsumerDialectUri(String consumerDialectUri) {

        this.consumerDialectUri = consumerDialectUri;
    }

    public String getSignatureAlgorithm() {

        return signatureAlgorithm;
    }

    public String getEncoding() {

        return encoding;
    }

    public void setSignatureAlgorithm(String signatureAlgorithm) {

        this.signatureAlgorithm = signatureAlgorithm;
    }

    public String getGatewayJWTGeneratorImpl() {

        return gatewayJWTGeneratorImpl;
    }

    public void setGatewayJWTGeneratorImpl(String gatewayJWTGeneratorImpl) {

        this.gatewayJWTGeneratorImpl = gatewayJWTGeneratorImpl;
    }

    public void setEncoding(String encoding) {

        this.encoding = encoding;
    }

    public Set<String> getJWTExcludedClaims() {

        return jwtExcludedClaims;
    }

    public Certificate getPublicCert() {

        return publicCert;
    }

    public void setPublicCert(Certificate publicCert) {

        this.publicCert = publicCert;
    }

    public PrivateKey getPrivateKey() {

        return privateKey;
    }

    public void setPrivateKey(PrivateKey privateKey) {

        this.privateKey = privateKey;
    }

    public void setTtl(long ttl) {

        this.ttl = ttl;
    }

    public Map<String, ClaimValueDTO> getCustomClaims() {

        return customClaims;
    }

    public long getTTL() {

        return ttl;
    }

    public void populateConfigValues(final boolean enabled, final String jwtHeader,
            final String signatureAlgorithm, final String encoding,
            final Certificate publicCert, final PrivateKey privateKey, final long ttl, final Map<String, ClaimValueDTO> customClaims,boolean useKid,String kidValue) {
        this.enabled = enabled;
        this.jwtHeader = jwtHeader;
        this.signatureAlgorithm = signatureAlgorithm;
        this.encoding = encoding;
        this.publicCert = publicCert;
        this.privateKey = privateKey;
        this.ttl = ttl;
        this.customClaims = customClaims;
        this.useKid = useKid;
        this.kidValue = kidValue;
    }

}
