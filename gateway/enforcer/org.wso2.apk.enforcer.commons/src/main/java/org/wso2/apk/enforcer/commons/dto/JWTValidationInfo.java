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

import com.nimbusds.jwt.JWTClaimsSet;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Holds jwt validation related information.
 */
public class JWTValidationInfo implements Serializable {
    private static final long serialVersionUID = 1L;

    private String user;
    private long expiryTime;
    private String consumerKey;
    private boolean valid;
    private List<String> scopes = new ArrayList<>();
    private Map<String, Object> claims = new HashMap<>();
    private int validationCode;
    private String keyManager;
    private String identifier;
    private JWTClaimsSet jwtClaimsSet;
    private String token;
    private List<String> audience = new ArrayList<>();

    public JWTValidationInfo() {

    }

    public JWTValidationInfo(JWTValidationInfo jwtValidationInfo) {

        this.user = jwtValidationInfo.getUser();
        this.expiryTime = jwtValidationInfo.getExpiryTime();
        this.consumerKey = jwtValidationInfo.getConsumerKey();
        this.valid = jwtValidationInfo.isValid();
        this.scopes = jwtValidationInfo.getScopes();
        this.claims = jwtValidationInfo.getClaims();
        this.validationCode = jwtValidationInfo.getValidationCode();
        this.keyManager = jwtValidationInfo.getKeyManager();
        this.audience = jwtValidationInfo.audience;
    }

    public List<String> getAudience() {
        return audience;
    }

    public void setAudience(List<String> audience) {
        this.audience = audience;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public String getIdentifier() {
        return identifier;
    }

    public void setIdentifier(String identifier) {
        this.identifier = identifier;
    }

    public JWTClaimsSet getJwtClaimsSet() {
        return jwtClaimsSet;
    }

    public void setJwtClaimsSet(JWTClaimsSet jwtClaimsSet) {
        this.jwtClaimsSet = jwtClaimsSet;
    }

    public String getUser() {

        return user;
    }

    public void setUser(String user) {

        this.user = user;
    }

    public long getExpiryTime() {

        return expiryTime;
    }

    public void setExpiryTime(long expiryTime) {

        this.expiryTime = expiryTime;
    }

    public boolean isValid() {

        return valid;
    }

    public void setValid(boolean valid) {

        this.valid = valid;
    }

    public List<String> getScopes() {

        return scopes;
    }

    public void setScopes(List<String> scopes) {

        this.scopes = scopes;
    }

    public Map<String, Object> getClaims() {

        return claims;
    }

    public void setClaims(Map<String, Object> claims) {

        this.claims = claims;
    }

    public String getConsumerKey() {

        return consumerKey;
    }

    public void setConsumerKey(String consumerKey) {

        this.consumerKey = consumerKey;
    }

    public int getValidationCode() {

        return validationCode;
    }

    public void setValidationCode(int validationCode) {

        this.validationCode = validationCode;
    }

    public String getKeyManager() {

        return keyManager;
    }

    public void setKeyManager(String keyManager) {

        this.keyManager = keyManager;
    }
}
