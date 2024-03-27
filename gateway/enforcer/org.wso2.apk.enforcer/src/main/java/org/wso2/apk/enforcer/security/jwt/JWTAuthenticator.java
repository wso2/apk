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
package org.wso2.apk.enforcer.security.jwt;

import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import com.nimbusds.jwt.util.DateUtils;
import io.opentelemetry.context.Scope;
import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.ThreadContext;
import org.wso2.apk.enforcer.common.CacheProviderUtil;
import org.wso2.apk.enforcer.commons.dto.ClaimValueDTO;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.dto.JWTInfoDto;
import org.wso2.apk.enforcer.commons.dto.JWTValidationInfo;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.commons.jwtgenerator.AbstractAPIMgtGatewayJWTGenerator;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.dto.APIKeyValidationInfoDTO;
import org.wso2.apk.enforcer.security.Authenticator;
import org.wso2.apk.enforcer.security.KeyValidator;
import org.wso2.apk.enforcer.security.TokenValidationContext;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;
import org.wso2.apk.enforcer.security.jwt.validator.RevokedJWTDataHolder;
import org.wso2.apk.enforcer.server.RevokedTokenRedisClient;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;
import org.wso2.apk.enforcer.tracing.TracingConstants;
import org.wso2.apk.enforcer.tracing.TracingSpan;
import org.wso2.apk.enforcer.tracing.TracingTracer;
import org.wso2.apk.enforcer.tracing.Utils;
import org.wso2.apk.enforcer.util.BackendJwtUtils;
import org.wso2.apk.enforcer.util.FilterUtils;
import org.wso2.apk.enforcer.util.JWTUtils;

import java.text.ParseException;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Implements the authenticator interface to authenticate request using a JWT token.
 */
public class JWTAuthenticator implements Authenticator {

    private static final Logger log = LogManager.getLogger(JWTAuthenticator.class);
    private final boolean isGatewayTokenCacheEnabled;
    private AbstractAPIMgtGatewayJWTGenerator jwtGenerator;

    public JWTAuthenticator(final JWTConfigurationDto jwtConfigurationDto, final boolean isGatewayTokenCacheEnabled) {

        this.isGatewayTokenCacheEnabled = isGatewayTokenCacheEnabled;
        if (jwtConfigurationDto.isEnabled()) {
            this.jwtGenerator = BackendJwtUtils.getApiMgtGatewayJWTGenerator(jwtConfigurationDto);
            this.jwtGenerator.setJWTConfigurationDto(jwtConfigurationDto);
        }
    }

    @Override
    public boolean canAuthenticate(RequestContext requestContext) {
        String authHeader = getTokenHeader(requestContext.getMatchedResourcePaths());

        if (!StringUtils.equals(authHeader, "")) {
            String authHeaderValue = retrieveAuthHeaderValue(requestContext, authHeader);
            return authHeaderValue != null && !StringUtils.containsWhitespace(authHeaderValue) && authHeaderValue.split("\\.").length == 3;
        }
        return false;
    }

    @Override
    public AuthenticationContext authenticate(RequestContext requestContext) throws APISecurityException {

        TracingTracer tracer = null;
        TracingSpan jwtAuthenticatorInfoSpan = null;
        Scope jwtAuthenticatorInfoSpanScope = null;
        TracingSpan validateScopesSpan = null;

        try {
            if (Utils.tracingEnabled()) {
                tracer = Utils.getGlobalTracer();
                jwtAuthenticatorInfoSpan = Utils.startSpan(TracingConstants.JWT_AUTHENTICATOR_SPAN, tracer);
                jwtAuthenticatorInfoSpanScope = jwtAuthenticatorInfoSpan.getSpan().makeCurrent();
                Utils.setTag(jwtAuthenticatorInfoSpan, APIConstants.LOG_TRACE_ID,
                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
            }
            String authHeader = getTokenHeader(requestContext.getMatchedResourcePaths());
            String jwtToken = retrieveAuthHeaderValue(requestContext, authHeader);
            String context = requestContext.getMatchedAPI().getBasePath();
            String version = requestContext.getMatchedAPI().getVersion();
            String organization = requestContext.getMatchedAPI().getOrganizationId();
            String environment = requestContext.getMatchedAPI().getEnvironment();

            JWTValidationInfo validationInfo = getJwtValidationInfo(jwtToken, organization, environment);
            if (RevokedTokenRedisClient.getRevokedTokens().contains(validationInfo.getIdentifier())) {
                log.info("Revoked JWT token. ", validationInfo.getIdentifier());
                throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS,
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
            }
            if (validationInfo != null) {
                if (validationInfo.isValid()) {
                    List<String> audFromAPI = getAudience(requestContext.getMatchedResourcePaths());
                    List<String> audFromToken = validationInfo.getAudience();
                    if (!checkAnyExist(audFromAPI, audFromToken)) {
                        throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                                APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
                    }
                    Map<String, Object> claims = validationInfo.getClaims();
                    // Validate token type
                    Object keyType = claims.get("keytype");
                    if (keyType != null && !keyType.toString().equalsIgnoreCase(requestContext.getMatchedAPI().getEnvType())) {
                        throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                                APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, "Invalid key type.");
                    }

                    // Scope validation is only done for tokens that are not of type InternalKey
                    Object tokenType = claims.get(APIConstants.JwtTokenConstants.TOKEN_TYPE);
                    if (!isInternalKey(tokenType)) {
                        // Validate scopes
                        Scope validateScopesSpanScope = null;
                        try {
                            if (Utils.tracingEnabled()) {
                                validateScopesSpan = Utils.startSpan(TracingConstants.SCOPES_VALIDATION_SPAN, tracer);
                                validateScopesSpanScope = validateScopesSpan.getSpan().makeCurrent();
                                Utils.setTag(validateScopesSpan, APIConstants.LOG_TRACE_ID,
                                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
                            }
                            validateScopes(context, version, requestContext.getMatchedResourcePaths(), validationInfo,
                                    jwtToken);
                        } finally {
                            if (Utils.tracingEnabled()) {
                                validateScopesSpanScope.close();
                                Utils.finishSpan(validateScopesSpan);
                            }
                        }
                    }

                    log.debug("JWT authentication successful.");

                    // Generate or get backend JWT
                    String endUserToken = null;

                    // JWT generator is only set if the backend JWT is enabled
                    if (this.jwtGenerator != null) {
                        JWTConfigurationDto configurationDto = this.jwtGenerator.getJWTConfigurationDto();
                        Map<String, ClaimValueDTO> claimMap = new HashMap<>();
                        if (configurationDto != null) {
                            claimMap = configurationDto.getCustomClaims();
                        }
                        JWTInfoDto jwtInfoDto = FilterUtils.generateJWTInfoDto(null, validationInfo,
                                null, requestContext);

                        // set custom claims get from the CR
                        jwtInfoDto.setClaims(claimMap);
                        endUserToken = BackendJwtUtils.generateAndRetrieveJWTToken(this.jwtGenerator,
                                validationInfo.getIdentifier(), jwtInfoDto, isGatewayTokenCacheEnabled, organization);
                        // Set generated jwt token as a response header
                        // Change the backendJWTConfig to API level
                        requestContext.addOrModifyHeaders(this.jwtGenerator.getJWTConfigurationDto().getJwtHeader(),
                                endUserToken);
                    }

                    return FilterUtils.generateAuthenticationContext(requestContext, validationInfo.getIdentifier(),
                            validationInfo, null, endUserToken, jwtToken, true);
                } else {
                    throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                            validationInfo.getValidationCode(),
                            APISecurityConstants.getAuthenticationFailureMessage(validationInfo.getValidationCode()));
                }
            } else {
                throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                        APISecurityConstants.API_AUTH_GENERAL_ERROR,
                        APISecurityConstants.API_AUTH_GENERAL_ERROR_MESSAGE);
            }
        } finally {
            if (Utils.tracingEnabled()) {
                jwtAuthenticatorInfoSpanScope.close();
                Utils.finishSpan(jwtAuthenticatorInfoSpan);
            }
        }

    }

    public boolean isInternalKey(Object tokenType) {
        return tokenType != null && tokenType.toString().equalsIgnoreCase(APIConstants.JwtTokenConstants.INTERNAL_KEY_TOKEN_TYPE);
    }

    private String getTokenHeader(ArrayList<ResourceConfig> matchedResourceConfigs) {
        for (ResourceConfig resourceConfig : matchedResourceConfigs) {
            if (resourceConfig.getAuthenticationConfig() != null &&
                    resourceConfig.getAuthenticationConfig().getJwtAuthenticationConfig() != null) {
                return resourceConfig.getAuthenticationConfig().getJwtAuthenticationConfig().getHeader();
            }
        }
        return "";
    }

    private List<String> getAudience(ArrayList<ResourceConfig> matchedResourceConfigs) {
        List<String> audience = new ArrayList<>();
        for (ResourceConfig resourceConfig : matchedResourceConfigs) {
            if (resourceConfig.getAuthenticationConfig() != null &&
                    resourceConfig.getAuthenticationConfig().getJwtAuthenticationConfig() != null) {
                return resourceConfig.getAuthenticationConfig().getJwtAuthenticationConfig().getAudience();
            }
        }
        return audience;
    }

    @Override
    public String getChallengeString() {

        return "JWT realm=\"APK\"";
    }

    @Override
    public String getName() {

        return "JWT";
    }

    @Override
    public int getPriority() {

        return 10;
    }

    private String retrieveAuthHeaderValue(RequestContext requestContext, String header) {
        Map<String, String> headers = requestContext.getHeaders();
        return headers.get(header);
    }

    /**
     * Validate scopes bound to the resource of the API being invoked against the scopes specified
     * in the JWT token payload.
     *
     * @param apiContext        API Context
     * @param apiVersion        API Version
     * @param matchingResources Accessed API resources
     * @param jwtValidationInfo Validated JWT Information
     * @param jwtToken          JWT Token
     * @throws APISecurityException in case of scope validation failure
     */
    private void validateScopes(String apiContext, String apiVersion, ArrayList<ResourceConfig> matchingResources,
                                JWTValidationInfo jwtValidationInfo, String jwtToken) throws APISecurityException {

        APIKeyValidationInfoDTO apiKeyValidationInfoDTO = new APIKeyValidationInfoDTO();
        Set<String> scopeSet = new HashSet<>(jwtValidationInfo.getScopes());
        apiKeyValidationInfoDTO.setScopes(scopeSet);

        TokenValidationContext tokenValidationContext = new TokenValidationContext();
        tokenValidationContext.setValidationInfoDTO(apiKeyValidationInfoDTO);

        tokenValidationContext.setAccessToken(jwtToken);
        // since matching resources has same method for all, just get the first element's method is adequate.
        // i.e. graphQL matching resources has same operation type for a request.
        tokenValidationContext.setHttpVerb(matchingResources.get(0).getMethod().toString());
        tokenValidationContext.setMatchingResourceConfigs(matchingResources);
        tokenValidationContext.setContext(apiContext);
        tokenValidationContext.setVersion(apiVersion);

        boolean valid = KeyValidator.validateScopes(tokenValidationContext);
        if (valid) {
            log.debug("Scope validation was successful for the resource.");
        }
    }

    /**
     * Validate whether the user is subscribed to the invoked API using consumer key. If subscribed, validate
     * the API information embedded within the Subscription against the information from the request context.
     *
     * @param validationInfo Token validation related details. This will be populated based on the available data
     *                       during the subscription validation.
     * @param name           API name
     * @param version        API version
     * @param context        API context
     * @param consumerKey    Consumer key extracted from the jwt token claim set
     * @param envType        The environment type, i.e. PRODUCTION or SANDBOX
     * @param organization   Organization extracted from the request context
     * @param tokenPrefix     The split token
     * @param matchedAPI
     * @throws APISecurityException if the user is not subscribed to the API
     */
    private void validateSubscriptionUsingConsumerKey(APIKeyValidationInfoDTO validationInfo, String name,
                                                      String version, String context, String consumerKey,
                                                      String envType, String organization, String tokenPrefix,
                                                      APIConfig matchedAPI) throws APISecurityException {

        validationInfo.setApiName(name);
        validationInfo.setApiVersion(version);
        validationInfo.setApiContext(context);
        validationInfo.setConsumerKey(consumerKey);
        validationInfo.setType(matchedAPI.getApiType());
        validationInfo.setEnvType(envType);
        validationInfo.setEnvironment(matchedAPI.getEnvironment());
        validationInfo.setSecurityScheme(APIConstants.API_SECURITY_OAUTH2);
        validationInfo.setSubscriberOrganization(organization);
        validationInfo.setApiContext(matchedAPI.getBasePath());
        validationInfo.setApiVersion(matchedAPI.getVersion());
        validationInfo.setApiName(matchedAPI.getName());
        KeyValidator.validateSubscriptionUsingConsumerKey(validationInfo);

        if (validationInfo.isAuthorized()) {
            if (log.isDebugEnabled()) {
                log.debug("User is subscribed to the API: " + name + ", " + "version: " + version + ". Token:" + " " +
                        FilterUtils.getMaskedToken(tokenPrefix));
            }
        } else {
            if (log.isDebugEnabled()) {
                log.debug("User is not subscribed to access the API: " + name + ", version: " + version + ". " +
                        "Token: " + FilterUtils.getMaskedToken(tokenPrefix));
            }
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHORIZED.getCode(),
                    APISecurityConstants.API_AUTH_FORBIDDEN, APISecurityConstants.API_AUTH_FORBIDDEN_MESSAGE);
        }
    }

    /**
     * Validate whether the token is a valid JWT and generate the JWTValidationInfo object.
     *
     * @param jwtToken     The full JWT token
     * @param organization organization of the API
     * @param environment  environment of the API
     * @return
     * @throws APISecurityException
     */
    private JWTValidationInfo getJwtValidationInfo(String jwtToken, String organization, String environment) throws APISecurityException {

        if (isGatewayTokenCacheEnabled) {
            String[] jwtParts = jwtToken.split("\\.");
            String signature = jwtParts[2];
            Object validCacheToken = CacheProviderUtil.getOrganizationCache(organization).getGatewayKeyCache()
                    .getIfPresent(signature);
            if (validCacheToken != null) {
                JWTValidationInfo validationInfo = (JWTValidationInfo) validCacheToken;
                if (!isJWTExpired(validationInfo)) {
                    if (!StringUtils.equals(validationInfo.getToken(), jwtToken)) {
                        log.warn("Suspected tampered token; a JWT token with the same signature is " +
                                "already available in the cache. Tampered token: " + FilterUtils.getMaskedToken(jwtToken));
                        throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                                APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, "Invalid JWT token");
                    }
                    if (RevokedJWTDataHolder.isJWTTokenSignatureExistsInRevokedMap(validationInfo.getIdentifier())) {
                        log.debug("Token found in the revoked jwt token map.");
                        throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                                APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, "Invalid JWT token");
                    }
                    return validationInfo;
                } else {
                    CacheProviderUtil.getOrganizationCache(organization).getGatewayKeyCache().invalidate(signature);
                    CacheProviderUtil.getOrganizationCache(organization).getInvalidTokenCache().put(signature, true);
                    throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                            APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED,
                            APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED_MESSAGE);
                }
            } else if (CacheProviderUtil.getOrganizationCache(organization).getInvalidTokenCache()
                    .getIfPresent(signature) != null) {
                throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS,
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
            }
        }

        SignedJWT signedJWT;
        JWTClaimsSet jwtClaimsSet;
        SignedJWTInfo signedJWTInfo;
        Scope decodeTokenHeaderSpanScope = null;
        TracingSpan decodeTokenHeaderSpan = null;
        try {
            if (Utils.tracingEnabled()) {
                TracingTracer tracer = Utils.getGlobalTracer();
                decodeTokenHeaderSpan = Utils.startSpan(TracingConstants.DECODE_TOKEN_HEADER_SPAN, tracer);
                decodeTokenHeaderSpanScope = decodeTokenHeaderSpan.getSpan().makeCurrent();
                Utils.setTag(decodeTokenHeaderSpan, APIConstants.LOG_TRACE_ID,
                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
            }
            signedJWT = SignedJWT.parse(jwtToken);
            jwtClaimsSet = signedJWT.getJWTClaimsSet();
            signedJWTInfo = new SignedJWTInfo(jwtToken, signedJWT, jwtClaimsSet);
        } catch (ParseException | IllegalArgumentException e) {
            log.error("Failed to decode the token header. {}", e.getMessage());
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, "Not a JWT token. Failed to decode the " +
                    "token header", e);
        } finally {
            if (Utils.tracingEnabled()) {
                decodeTokenHeaderSpanScope.close();
                Utils.finishSpan(decodeTokenHeaderSpan);
            }
        }

        String signature = signedJWT.getSignature().toString();
        String jwtTokenIdentifier = StringUtils.isNotEmpty(jwtClaimsSet.getJWTID()) ? jwtClaimsSet.getJWTID() :
                signature;

        // check whether the token is revoked
        String jwtHeader = signedJWT.getHeader().toString();
        if (RevokedJWTDataHolder.isJWTTokenSignatureExistsInRevokedMap(jwtTokenIdentifier)) {
            log.debug("Token retrieved from the revoked jwt token map. Token: " +
                    FilterUtils.getMaskedToken(jwtHeader));
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_INVALID_CREDENTIALS, "Invalid JWT token");
        }

        try {
            // Get issuer
            String issuer = jwtClaimsSet.getIssuer();
            SubscriptionDataStore subscriptionDataStore = SubscriptionDataHolder.getInstance()
                    .getSubscriptionDataStore(organization);
            if (subscriptionDataStore == null) {
                throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS,
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
            }
            JWTValidator jwtValidator = subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment);
            // If no validator found for the issuer, we are not caching the token.
            if (jwtValidator == null) {
                throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS,
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
            }

            JWTValidationInfo jwtValidationInfo = jwtValidator.validateToken(jwtToken, signedJWTInfo);
            if (isGatewayTokenCacheEnabled) {
                // Add token to tenant token cache
                if (jwtValidationInfo.isValid()) {
                    CacheProviderUtil.getOrganizationCache(organization).getGatewayKeyCache().put(signature,
                            jwtValidationInfo);
                } else {
                    CacheProviderUtil.getOrganizationCache(organization).getInvalidTokenCache().put(signature, true);
                }
            }
            return jwtValidationInfo;
        } catch (EnforcerException e) {
            log.error("JWT Validation failed", e);
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_GENERAL_ERROR,
                    APISecurityConstants.API_AUTH_GENERAL_ERROR_MESSAGE);
        }
    }

    /**
     * Check whether the jwt token is expired or not.
     *
     * @param payload The payload of the JWT token
     * @return boolean true if the token is not expired, false otherwise
     */
    private Boolean isJWTExpired(JWTValidationInfo payload) {

        long timestampSkew = FilterUtils.getTimeStampSkewInSeconds();
        Date now = new Date();
        Date exp = new Date(payload.getExpiryTime());
        return !DateUtils.isAfter(exp, now, timestampSkew);
    }

    /**
     * Checks if at least one element from list1 exists in list2.
     *
     * @param list1 The first list to check.
     * @param list2 The second list to check against.
     * @return true if list1 is empty else if at least one element from list1 exists in list2, otherwise false.
     */
    public static boolean checkAnyExist(List<String> list1, List<String> list2) {
        if (list1.size() == 0) {
            return true;
        }
        return list1.stream().anyMatch(list2::contains);
    }
}
