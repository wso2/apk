/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.util;

import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.common.CacheProviderUtil;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.dto.JWTInfoDto;
import org.wso2.apk.enforcer.commons.exception.JWTGeneratorException;
import org.wso2.apk.enforcer.commons.jwtgenerator.APIMgtGatewayJWTGeneratorImpl;
import org.wso2.apk.enforcer.commons.jwtgenerator.APIMgtGatewayUrlSafeJWTGeneratorImpl;
import org.wso2.apk.enforcer.commons.jwtgenerator.AbstractAPIMgtGatewayJWTGenerator;
import org.wso2.apk.enforcer.commons.jwttransformer.JWTTransformer;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.JwtConstants;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.security.jwt.JwtTransformerAnnotation;
import org.wso2.apk.enforcer.security.jwt.validator.JWTConstants;

import java.lang.annotation.Annotation;
import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.util.Base64;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.ServiceLoader;

/**
 * Contains Util methods related to backend JWT generation.
 */
public class BackendJwtUtils {

    private static final Logger log = LogManager.getLogger(BackendJwtUtils.class);

    /**
     * Generates or gets the Cached Backend JWT token.
     *
     * @param jwtGenerator               the jwtGenerator instance to use if generating the token
     * @param tokenSignature             token signature to use in the cache key
     * @param jwtInfoDto                 information to include in the jwt
     * @param isGatewayTokenCacheEnabled whether gateway token cache is enabled
     * @return backend jwt token
     * @throws APISecurityException if an error occurs while generating the token
     */
    public static String generateAndRetrieveJWTToken(AbstractAPIMgtGatewayJWTGenerator jwtGenerator,
                                                     String tokenSignature, JWTInfoDto jwtInfoDto,
                                                     boolean isGatewayTokenCacheEnabled, String organization) throws APISecurityException {

        log.debug("Inside generateAndRetrieveJWTToken");
        String endUserToken = null;
        boolean valid = false;
        String jwtTokenCacheKey = jwtInfoDto.getApiContext().concat(":").concat(jwtInfoDto.getVersion()).concat(":")
                .concat(tokenSignature); // TODO: (suksw) Check if to add tenantName or label also

        if (jwtGenerator != null) {
            if (isGatewayTokenCacheEnabled) {
                try {
                    Object token =
                            CacheProviderUtil.getOrganizationCache(organization).getGatewayJWTTokenCache().get(jwtTokenCacheKey);
                    if (!JWTConstants.UNAVAILABLE.equals(token)) {
                        endUserToken = (String) token;
                        String[] splitToken = endUserToken.split("\\.");
                        org.json.JSONObject payload = new org.json.JSONObject(new String(Base64.getUrlDecoder().decode(splitToken[1])));
                        long exp = payload.getLong(JwtConstants.EXP);
                        valid = !JWTUtils.isExpired(exp);
                    }
                } catch (Exception e) {
                    log.error("Error while getting token from the cache", e);
                }

                if (StringUtils.isEmpty(endUserToken) || !valid) {
                    endUserToken = generateToken(jwtGenerator, jwtInfoDto, true, jwtTokenCacheKey, organization);
                }
            } else {
                endUserToken = generateToken(jwtGenerator, jwtInfoDto, false, jwtTokenCacheKey, organization);
            }
        } else {
            log.debug("Error while loading JWTGenerator");
        }
        return endUserToken;
    }

    private static String generateToken(AbstractAPIMgtGatewayJWTGenerator jwtGenerator, JWTInfoDto jwtInfoDto,
                                        boolean isGatewayTokenCacheEnabled, String jwtTokenCacheKey,
                                        String organization) throws APISecurityException {

        String endUserToken;
        JWTConfigurationDto jwtConfigurationDto = jwtGenerator.getJWTConfigurationDto();
        jwtGenerator.setJWTConfigurationDto(jwtConfigurationDto);
        try {
            endUserToken = jwtGenerator.generateToken(jwtInfoDto);
            if (isGatewayTokenCacheEnabled) {
                CacheProviderUtil.getOrganizationCache(organization).getGatewayJWTTokenCache().put(jwtTokenCacheKey,
                        endUserToken);
            }
        } catch (JWTGeneratorException e) {
            log.error("Error while Generating Backend JWT", e);
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_GENERAL_ERROR,
                    APISecurityConstants.API_AUTH_GENERAL_ERROR_MESSAGE, e);
        }
        return endUserToken;
    }

    /**
     * Load the specified backend JWT Generator.
     *
     * @param jwtConfigurationDtoFromAPI jwt configuration dto from api
     * @return an instance of the JWT Generator given in the config
     */
    public static AbstractAPIMgtGatewayJWTGenerator getApiMgtGatewayJWTGenerator(final JWTConfigurationDto jwtConfigurationDtoFromAPI) {

        JWTConfigurationDto jwtConfigurationDto = jwtConfigurationDtoFromAPI;
        AbstractAPIMgtGatewayJWTGenerator jwtGenerator = null;
        if ("Base64".equals(jwtConfigurationDto.getEncoding())){
            return new APIMgtGatewayJWTGeneratorImpl();
        }else if ("Base64Url".equals(jwtConfigurationDto.getEncoding())){
            return new APIMgtGatewayUrlSafeJWTGeneratorImpl();
        }
        return null;
    }

    /**
     * Load the specified JWT Transformers.
     *
     * @return a map of JWT Transformers
     */
    public static Map<String, JWTTransformer> loadJWTTransformers() {

        ServiceLoader<JWTTransformer> loader = ServiceLoader.load(JWTTransformer.class);
        Iterator<JWTTransformer> classIterator = loader.iterator();
        Map<String, JWTTransformer> jwtTransformersMap = new HashMap<>();

        if (!classIterator.hasNext()) {
            log.debug("No JWTTransformers found.");
            return jwtTransformersMap;
        }

        while (classIterator.hasNext()) {
            JWTTransformer transformer = classIterator.next();
            Annotation[] annotations = transformer.getClass().getAnnotations();
            if (annotations.length == 0) {
                log.debug("JWTTransformer is discarded as no annotations found : {}",
                        transformer.getClass().getCanonicalName());
                continue;
            }
            for (Annotation annotation : annotations) {
                if (annotation instanceof JwtTransformerAnnotation) {
                    JwtTransformerAnnotation jwtTransformerAnnotation =
                            (JwtTransformerAnnotation) annotation;
                    if (jwtTransformerAnnotation.enabled()) {
                        log.debug("JWTTransformer for the issuer : {} is enabled.",
                                jwtTransformerAnnotation.issuer());
                        jwtTransformersMap.put(jwtTransformerAnnotation.issuer(), transformer);
                    } else {
                        log.debug("JWTTransformer for the issuer : {} is disabled.",
                                jwtTransformerAnnotation.issuer());
                    }
                }
            }
        }
        return jwtTransformersMap;
    }
}
