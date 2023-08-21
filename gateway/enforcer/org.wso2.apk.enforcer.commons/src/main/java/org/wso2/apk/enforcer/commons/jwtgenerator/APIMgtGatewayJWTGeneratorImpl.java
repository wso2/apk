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

package org.wso2.apk.enforcer.commons.jwtgenerator;

import org.apache.commons.lang3.StringUtils;
import org.wso2.apk.enforcer.commons.constants.JWTConstants;
import org.wso2.apk.enforcer.commons.dto.ClaimValueDTO;
import org.wso2.apk.enforcer.commons.dto.JWTInfoDto;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;


/**
 * Default implementation of backend jwt generation.
 */
public class APIMgtGatewayJWTGeneratorImpl extends AbstractAPIMgtGatewayJWTGenerator {

    @Override
    public Map<String, Object> populateStandardClaims(JWTInfoDto jwtInfoDto) {

        long currentTime = TimeUnit.MILLISECONDS.toSeconds(System.currentTimeMillis());
        long expireIn = currentTime + super.jwtConfigurationDto.getTTL();
        String dialect = getDialectURI();
        Map<String, Object> claims = new HashMap<>();
        claims.put("iss", API_GATEWAY_ID);
        claims.put("exp", String.valueOf(expireIn));
        claims.put("iat", String.valueOf(currentTime));
        // dialect is either empty or '/' do not append a backslash. otherwise append a backslash '/'
        if (!"".equals(dialect) && !"/".equals(dialect)) {
            dialect = dialect + "/";
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getSubscriber())) {
            claims.put(dialect + "subscriber", jwtInfoDto.getSubscriber());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getApplicationId())) {
            claims.put(dialect + "applicationid", jwtInfoDto.getApplicationId());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getApplicationName())) {
            claims.put(dialect + "applicationname", jwtInfoDto.getApplicationName());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getApplicationTier())) {
            claims.put(dialect + "applicationtier", jwtInfoDto.getApplicationTier());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getApiName())) {
            claims.put(dialect + "apiname", jwtInfoDto.getApiName());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getApiContext())) {
            claims.put(dialect + "apicontext", jwtInfoDto.getApiContext());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getVersion())) {
            claims.put(dialect + "version", jwtInfoDto.getVersion());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getSubscriptionTier())) {
            claims.put(dialect + "tier", jwtInfoDto.getSubscriptionTier());
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getKeyType())) {
            claims.put(dialect + "keytype", jwtInfoDto.getKeyType());
        } else {
            claims.put(dialect + "keytype", "PRODUCTION");
        }
        claims.put(dialect + "usertype", JWTConstants.AUTH_APPLICATION_USER_LEVEL_TOKEN);
        claims.put(dialect + "enduser", jwtInfoDto.getEndUser());
        claims.put(dialect + "enduserTenantId", String.valueOf(jwtInfoDto.getEndUserTenantId()));
        claims.put(dialect + "applicationUUId", jwtInfoDto.getApplicationUUId());
        Map<String, String> appAttributes = jwtInfoDto.getAppAttributes();
        if (appAttributes != null && !appAttributes.isEmpty()) {
            claims.put(dialect + "applicationAttributes", appAttributes);
        }
        if (StringUtils.isNotEmpty(jwtInfoDto.getSub())) {
            claims.put(JWTConstants.SUB, jwtInfoDto.getSub());
        }
        if (jwtInfoDto.getOrganizations() != null) {
            claims.put(JWTConstants.ORGANIZATIONS, jwtInfoDto.getOrganizations());
        }
        return claims;
    }

    @Override
    public Map<String, ClaimValueDTO> populateCustomClaims(JWTInfoDto jwtInfoDto) {

        String[] restrictedClaims = {"iss", "sub", "aud", "exp", "nbf", "iat", "jti", "application", "tierInfo",
                "subscribedAPIs", "aut"};
        Map<String, ClaimValueDTO> claims = new HashMap<>();
        Set<String> jwtExcludedClaims = jwtConfigurationDto.getJWTExcludedClaims();
        jwtExcludedClaims.addAll(Arrays.asList(restrictedClaims));
        Map<String, Object> jwtToken = jwtInfoDto.getJwtValidationInfo().getClaims();
        if (jwtToken != null) {
            for (Map.Entry<String, Object> jwtClaimEntry : jwtToken.entrySet()) {
                if (!jwtExcludedClaims.contains(jwtClaimEntry.getKey())) {
                    ClaimValueDTO claimValue = new ClaimValueDTO(jwtClaimEntry.getValue(), null);
                    claims.put(jwtClaimEntry.getKey(), claimValue);
                }
            }
        }
        Map<String, ClaimValueDTO> customClaimsAPI = jwtInfoDto.getClaims();
        if(customClaimsAPI != null) {
            for (Map.Entry<String, ClaimValueDTO> customClaimEntry : customClaimsAPI.entrySet()) {
                ClaimValueDTO claim = new ClaimValueDTO(customClaimEntry.getValue().getValue(), customClaimEntry.getValue().getType());
                if(claims.containsKey(customClaimEntry.getKey())) {
                    claims.replace(customClaimEntry.getKey(), claim);
                } else {
                    claims.put(customClaimEntry.getKey(), claim);
                }

            }
        }
        return claims;
    }
}

