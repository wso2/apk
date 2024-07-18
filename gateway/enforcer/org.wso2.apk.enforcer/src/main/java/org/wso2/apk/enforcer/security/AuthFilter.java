/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package org.wso2.apk.enforcer.security;

import java.util.stream.Collectors;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.Filter;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.commons.logging.ErrorDetails;
import org.wso2.apk.enforcer.commons.logging.LoggingConstants;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.APIKeyAuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.JWTAuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.Oauth2AuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.config.dto.MutualSSLDto;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.InterceptorConstants;
import org.wso2.apk.enforcer.security.jwt.APIKeyAuthenticator;
import org.wso2.apk.enforcer.security.jwt.JWTAuthenticator;
import org.wso2.apk.enforcer.security.jwt.Oauth2Authenticator;
import org.wso2.apk.enforcer.security.jwt.UnsecuredAPIAuthenticator;
import org.wso2.apk.enforcer.security.mtls.MTLSAuthenticator;
import org.wso2.apk.enforcer.util.FilterUtils;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * This is the filter handling the authentication for the requests flowing
 * through the gateway.
 */
public class AuthFilter implements Filter {
    private List<Authenticator> authenticators = new ArrayList<>();
    private static final Logger log = LogManager.getLogger(AuthFilter.class);
    private boolean isMutualSSLMandatory;
    private boolean isOAuth2Mandatory;
    private boolean isAPIKeyMandatory;

    @Override
    public void init(APIConfig apiConfig, Map<String, String> configProperties) {
        initializeAuthenticators(apiConfig);
    }

    private void initializeAuthenticators(APIConfig apiConfig) {
        // TODO: Check security schema and add relevant authenticators.
        boolean isMutualSSLProtected = false;
        isMutualSSLMandatory = false;

        // Set security conditions
        Map<String, Boolean> securityMaps = apiConfig.getApplicationSecurity();
        isOAuth2Mandatory = securityMaps.getOrDefault("OAuth2", true);
        isAPIKeyMandatory = securityMaps.getOrDefault("APIKey", false);

        if (!Objects.isNull(apiConfig.getMutualSSL())) {
            if (apiConfig.isTransportSecurity()) {
                if (apiConfig.getMutualSSL().equalsIgnoreCase(APIConstants.Optionality.MANDATORY)) {
                    isMutualSSLProtected = true;
                    isMutualSSLMandatory = true;
                } else if (apiConfig.getMutualSSL().equalsIgnoreCase(APIConstants.Optionality.OPTIONAL)) {
                    isMutualSSLProtected = true;
                }
            } else {
                isMutualSSLProtected = false;
            }
        }

        if (isMutualSSLProtected) {
            Authenticator mtlsAuthenticator = new MTLSAuthenticator();
            authenticators.add(mtlsAuthenticator);
        }

        // check whether the backend JWT token is enabled
        EnforcerConfig enforcerConfig = ConfigHolder.getInstance().getConfig();
        boolean isGatewayTokenCacheEnabled = enforcerConfig.getCacheDto().isEnabled();
        JWTConfigurationDto jwtConfigurationDto = apiConfig.getJwtConfigurationDto();

        Authenticator oauthAuthenticator = new Oauth2Authenticator(jwtConfigurationDto, isGatewayTokenCacheEnabled);
        authenticators.add(oauthAuthenticator);

        APIKeyAuthenticator apiKeyAuthenticator = new APIKeyAuthenticator(jwtConfigurationDto);
        authenticators.add(apiKeyAuthenticator);

        Authenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, isGatewayTokenCacheEnabled);
        authenticators.add(jwtAuthenticator);

        Authenticator unsecuredAPIAuthenticator = new UnsecuredAPIAuthenticator();
        authenticators.add(unsecuredAPIAuthenticator);

        authenticators.sort(Comparator.comparingInt(Authenticator::getPriority));
    }

    @Override
    public boolean handleRequest(RequestContext requestContext) {
        populateRemoveAndProtectedAuthHeaders(requestContext);
        // Set API metadata for interceptors
        setInterceptorAPIMetadata(requestContext);

        // It is required to skip the auth Filter if the lifecycle status is prototype
        if (APIConstants.PROTOTYPED_LIFE_CYCLE_STATUS.equals(
                requestContext.getMatchedAPI().getApiLifeCycleState()) &&
                !requestContext.getMatchedAPI().isMockedApi()) {
            return true;
        }

        // Authentication status of the request
        boolean authenticated = false;
        // Any auth token has been provided for application-level security or not
        boolean canAuthenticated = false;
        for (Authenticator authenticator : authenticators) {
            if (authenticator.canAuthenticate(requestContext)) {
                // For transport level securities (mTLS), canAuthenticated will not be applied
                if (!authenticator.getName().contains(APIConstants.API_SECURITY_MUTUAL_SSL_NAME)) {
                    canAuthenticated = true;
                }
                AuthenticationResponse authenticateResponse = authenticate(authenticator, requestContext);

                authenticated = authenticateResponse.isAuthenticated();
                setInterceptorAuthContextMetadata(authenticator, requestContext);

                if (!authenticateResponse.isContinueToNextAuthenticator()) {
                    break;
                }
            } else {
                // Check if the failed authentication is mandatory mTLS
                if (isMutualSSLMandatory && authenticator.getName()
                        .contains(APIConstants.API_SECURITY_MUTUAL_SSL_NAME)) {
                    authenticated = false;
                    log.debug("mTLS authentication was failed for the request: {} , API: {}:{} APIUUID: {} ",
                            requestContext.getMatchedResourcePaths().get(0).getPath(),
                            requestContext.getMatchedAPI().getName(), requestContext.getMatchedAPI().getVersion(),
                            requestContext.getMatchedAPI().getUuid());
                    break;
                }
                // Check if the failed authentication is a mandatory application level security
                if (!authenticator.getName()
                        .contains(APIConstants.API_SECURITY_MUTUAL_SSL_NAME) && (isAPIKeyMandatory||isOAuth2Mandatory)){
                    authenticated = false;
                }

            }
        }
        if (authenticated) {
            return true;
        }
        if (!canAuthenticated) {
            FilterUtils.setUnauthenticatedErrorToContext(requestContext);
        }
        log.debug("None of the authenticators were able to authenticate the request: {}",
                requestContext.getRequestPathTemplate(),
                ErrorDetails.errorLog(LoggingConstants.Severity.MINOR, 6600));
        // set WWW_AUTHENTICATE header to error response
        requestContext.addOrModifyHeaders(APIConstants.WWW_AUTHENTICATE, getAuthenticatorsChallengeString() +
                ", error=\"invalid_token\"" +
                ", error_description=\"The provided token is invalid\"");
        return false;
    }

    private AuthenticationResponse authenticate(Authenticator authenticator, RequestContext requestContext) {
        try {
            AuthenticationContext authenticate = authenticator.authenticate(requestContext);
            requestContext.setAuthenticationContext(authenticate);
            if (authenticator.getName().contains(APIConstants.API_SECURITY_MUTUAL_SSL_NAME)) {
                // This section is for mTLS authentication
                if (authenticate.isAuthenticated()) {
                    log.debug("mTLS authentication was passed for the request: {} , API: {}:{}, APIUUID: {} ",
                            requestContext.getMatchedResourcePaths().get(0).getPath(),
                            requestContext.getMatchedAPI().getName(), requestContext.getMatchedAPI().getVersion(),
                            requestContext.getMatchedAPI().getUuid());

                    boolean isApplicationSecurityDisabled = requestContext.getMatchedResourcePaths().get(0)
                            .getAuthenticationConfig().isDisabled();
                    // proceed to the next authenticator only if application security is enabled and
                    // is mandatory
                    return new AuthenticationResponse(true, isMutualSSLMandatory,
                            !isApplicationSecurityDisabled && (isOAuth2Mandatory || isAPIKeyMandatory));
                } else {
                    if (isMutualSSLMandatory) {
                        log.debug("Mandatory mTLS authentication was failed for the request: {} , API: {}:{}, " +
                                        "APIUUID: {} ",
                                requestContext.getMatchedResourcePaths().get(0).getPath(),
                                requestContext.getMatchedAPI().getName(), requestContext.getMatchedAPI().getVersion(),
                                requestContext.getMatchedAPI().getUuid());
                    } else {
                        log.debug("Optional mTLS authentication was failed for the request: {} , API: {}:{}, " +
                                        "APIUUID: {} ",
                                requestContext.getMatchedResourcePaths().get(0).getPath(),
                                requestContext.getMatchedAPI().getName(), requestContext.getMatchedAPI().getVersion(),
                                requestContext.getMatchedAPI().getUuid());
                    }
                    return new AuthenticationResponse(false, isMutualSSLMandatory, false);

                }
                // for all authenticators other than mTLS
            } else if (authenticate.isAuthenticated()) {
                return new AuthenticationResponse(true, isOAuth2Mandatory && isAPIKeyMandatory, false);
            }
        } catch (APISecurityException e) {
            // TODO: (VirajSalaka) provide the error code properly based on exception (401,
            // 403, 429 etc)
            FilterUtils.setErrorToContext(requestContext, e);
        }
        boolean continueToNextAuth = true;
        if (authenticator.getName().contains(APIConstants.API_SECURITY_MUTUAL_SSL_NAME)) {
            continueToNextAuth = false;
        }
        return new AuthenticationResponse(false,
                isOAuth2Mandatory || isMutualSSLMandatory, continueToNextAuth);
    }

    private String getAuthenticatorsChallengeString() {
        StringBuilder challengeString = new StringBuilder();
        if (authenticators != null) {
            for (Authenticator authenticator : authenticators) {
                challengeString.append(authenticator.getChallengeString()).append(" ");
            }
        }
        return challengeString.toString().trim();
    }

    private void setInterceptorAuthContextMetadata(Authenticator authenticator, RequestContext requestContext) {
        // add auth context to metadata, lua script will add it to the auth context of
        // the interceptor
        AuthenticationContext authContext = requestContext.getAuthenticationContext();
        String tokenType = authenticator.getName();
        authContext.setTokenType(tokenType);
        requestContext.addMetadataToMap(InterceptorConstants.AuthContextFields.TOKEN_TYPE,
                Objects.toString(tokenType, ""));
        requestContext.addMetadataToMap(InterceptorConstants.AuthContextFields.TOKEN,
                Objects.toString(authContext.getRawToken(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.AuthContextFields.KEY_TYPE,
                Objects.toString(requestContext.getMatchedAPI().getEnvType(), ""));
    }

    private void setInterceptorAPIMetadata(RequestContext requestContext) {
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.API_BASE_PATH,
                Objects.toString(requestContext.getMatchedAPI().getBasePath(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.API_VERSION,
                Objects.toString(requestContext.getMatchedAPI().getVersion(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.API_NAME,
                Objects.toString(requestContext.getMatchedAPI().getName(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.API_VHOST,
                Objects.toString(requestContext.getMatchedAPI().getVhost(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.API_ORGANIZATION_ID,
                Objects.toString(requestContext.getMatchedAPI().getOrganizationId(), ""));
        requestContext.addMetadataToMap(InterceptorConstants.APIMetadataFields.ENVIRONMENT,
                Objects.toString(requestContext.getMatchedAPI().getEnvironment(), ""));
    }

    private void populateRemoveAndProtectedAuthHeaders(RequestContext requestContext) {
        requestContext.getMatchedResourcePaths().forEach(resourcePath -> {
            Oauth2AuthenticationConfig oauth2AuthenticationConfig = resourcePath.getAuthenticationConfig()
                    .getOauth2AuthenticationConfig();
            JWTAuthenticationConfig jwtAuthenticationConfig = resourcePath.getAuthenticationConfig()
                    .getJwtAuthenticationConfig();
            List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfig = resourcePath.getAuthenticationConfig()
                    .getApiKeyAuthenticationConfigs();
            if (oauth2AuthenticationConfig != null && !oauth2AuthenticationConfig.isSendTokenToUpstream()) {
                requestContext.getRemoveHeaders().add(oauth2AuthenticationConfig.getHeader());
            }
            if (jwtAuthenticationConfig != null && !jwtAuthenticationConfig.isSendTokenToUpstream()) {
                requestContext.getRemoveHeaders().add(jwtAuthenticationConfig.getHeader());
            }
            if (apiKeyAuthenticationConfig != null && !apiKeyAuthenticationConfig.isEmpty()) {
                requestContext.getQueryParamsToRemove().addAll(apiKeyAuthenticationConfig.stream()
                        .filter(apiKeyAuthenticationConfig1 -> !apiKeyAuthenticationConfig1.isSendTokenToUpstream()
                                && Objects.equals(apiKeyAuthenticationConfig1.getIn(), "In"))
                        .map(APIKeyAuthenticationConfig::getName).collect(Collectors.toList()));
                List<String> apikeyHeadersToRemove = apiKeyAuthenticationConfig.stream()
                        .filter(apiKeyAuthenticationConfig1 -> !apiKeyAuthenticationConfig1.isSendTokenToUpstream()
                                && Objects.equals(apiKeyAuthenticationConfig1.getIn(), "Header"))
                        .map(APIKeyAuthenticationConfig::getName).collect(Collectors.toList());
                requestContext.getRemoveHeaders().addAll(apikeyHeadersToRemove);
            }
        });

        // Remove mTLS certificate header
        MutualSSLDto mtlsInfo = ConfigHolder.getInstance().getConfig().getMtlsInfo();
        String certificateHeaderName = FilterUtils.getCertificateHeaderName();
        if (!mtlsInfo.isEnableOutboundCertificateHeader()) {
            requestContext.getRemoveHeaders().add(certificateHeaderName);
        }
    }
}
