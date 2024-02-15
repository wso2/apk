/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.apk.enforcer.jwt;

import com.google.common.cache.LoadingCache;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.mockito.MockedStatic;
import org.mockito.Mockito;
import org.wso2.apk.enforcer.common.CacheProvider;
import org.wso2.apk.enforcer.common.CacheProviderUtil;
import org.wso2.apk.enforcer.commons.dto.JWKSConfigurationDTO;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.dto.JWTValidationInfo;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.commons.jwtgenerator.AbstractAPIMgtGatewayJWTGenerator;
import org.wso2.apk.enforcer.commons.jwttransformer.DefaultJWTTransformer;
import org.wso2.apk.enforcer.commons.jwttransformer.JWTTransformer;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.AuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.Oauth2AuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.config.dto.ExtendedTokenIssuerDto;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.security.KeyValidator;
import org.wso2.apk.enforcer.security.jwt.Oauth2Authenticator;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;
import org.wso2.apk.enforcer.server.RevokedTokenRedisClient;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.UUID;

public class Oauth2ValidatorTest {

    @Before
    public void setup() {

        RevokedTokenRedisClient.setRevokedTokens(new HashSet<>());
    }

    @Test
    public void testJWTValidator() throws APISecurityException, EnforcerException {

        String organization = "org1";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfFCEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJUMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV" +
                "3YXkiLCJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZ" +
                "SI6InJlcXVlc3RDb3VudCIsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoiRGVmYXVsdEFwcGxpY2F0aW9uIiwiaWQiOjEsInV1aWQ" +
                "iOm51bGx9LCJzY29wZSI6ImFtX2FwcGxpY2F0aW9uX3Njb3BlIGRlZmF1bHQiLCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0" +
                "NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6e30sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOltdL" +
                "CJjb25zdW1lcktleSI6IlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLCJleHAiOjE1OTAzNDIzMTMsImlhdCI6MTU5MDMzO" +
                "DcxMywianRpIjoiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIn0." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() + 5000L);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);
        try (MockedStatic<LogManager> logManagerDummy = Mockito.mockStatic(LogManager.class);
             MockedStatic<LogFactory> logFactoryDummy = Mockito.mockStatic(LogFactory.class);
             MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<ConfigHolder> configHolderDummy = Mockito.mockStatic(ConfigHolder.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);
             MockedStatic<KeyValidator> keyValidaterDummy = Mockito.mockStatic(KeyValidator.class)) {
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            Logger logger = Mockito.mock(Logger.class);
            logManagerDummy.when(() -> LogManager.getLogger(Oauth2Authenticator.class)).thenReturn(logger);
            Log logger2 = Mockito.mock(Log.class);
            logFactoryDummy.when(() -> LogFactory.getLog(AbstractAPIMgtGatewayJWTGenerator.class)).thenReturn(logger2);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization))
                    .thenReturn(cacheProvider);
            ExtendedTokenIssuerDto tokenIssuerDto = Mockito.mock(ExtendedTokenIssuerDto.class);
            Mockito.when(tokenIssuerDto.getIssuer()).thenReturn(issuer);
            JWKSConfigurationDTO jwksConfigurationDTO = new JWKSConfigurationDTO();
            jwksConfigurationDTO.setEnabled(true);
            Mockito.when(tokenIssuerDto.getJwksConfigurationDTO()).thenReturn(jwksConfigurationDTO);

            EnforcerConfig enforcerConfig = Mockito.mock(EnforcerConfig.class);
            ConfigHolder configHolder = Mockito.mock(ConfigHolder.class);
            configHolderDummy.when(ConfigHolder::getInstance).thenReturn(configHolder);
            Mockito.when(configHolder.getConfig()).thenReturn(enforcerConfig);
            JWTTransformer jwtTransformer = new DefaultJWTTransformer();
            Mockito.when(enforcerConfig.getJwtTransformer(issuer)).thenReturn(jwtTransformer);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);

            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment)).thenReturn(jwtValidator);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidaterDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = oauth2Authenticator.authenticate(requestContext);
            Assert.assertNotNull(authenticate);
            Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).put(signature, jwtValidationInfo);
        }
    }

    @Test
    public void testRevokedToken() throws APISecurityException, EnforcerException {
        HashSet<String> revokedTokens = new HashSet<>();
        String revokedTokenJTI = "b8938768-23fd-4dec-8b70-bed45eb7c33d";
        revokedTokens.add(revokedTokenJTI);
        RevokedTokenRedisClient.setRevokedTokens(revokedTokens);
        String organization = "org1";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfFCEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJUMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV" +
                "3YXkiLCJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZ" +
                "SI6InJlcXVlc3RDb3VudCIsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoiRGVmYXVsdEFwcGxpY2F0aW9uIiwiaWQiOjEsInV1aWQ" +
                "iOm51bGx9LCJzY29wZSI6ImFtX2FwcGxpY2F0aW9uX3Njb3BlIGRlZmF1bHQiLCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0" +
                "NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6e30sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOltdL" +
                "CJjb25zdW1lcktleSI6IlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLCJleHAiOjE1OTAzNDIzMTMsImlhdCI6MTU5MDMzO" +
                "DcxMywianRpIjoiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIn0." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() + 5000L);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setIdentifier(revokedTokenJTI);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);

        try (MockedStatic<LogManager> logManagerDummy = Mockito.mockStatic(LogManager.class);
             MockedStatic<LogFactory> logFactoryDummy = Mockito.mockStatic(LogFactory.class);
             MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<ConfigHolder> configHolderDummy = Mockito.mockStatic(ConfigHolder.class);
             MockedStatic<KeyValidator> keyValidaterDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);) {
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            Logger logger = Mockito.mock(Logger.class);
            logManagerDummy.when(() -> LogManager.getLogger(Oauth2Authenticator.class)).thenReturn(logger);
            Log logger2 = Mockito.mock(Log.class);
            logFactoryDummy.when(() -> LogFactory.getLog(AbstractAPIMgtGatewayJWTGenerator.class)).thenReturn(logger2);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization))
                    .thenReturn(cacheProvider);
            ExtendedTokenIssuerDto tokenIssuerDto = Mockito.mock(ExtendedTokenIssuerDto.class);
            Mockito.when(tokenIssuerDto.getIssuer()).thenReturn(issuer);
            JWKSConfigurationDTO jwksConfigurationDTO = new JWKSConfigurationDTO();
            jwksConfigurationDTO.setEnabled(true);
            Mockito.when(tokenIssuerDto.getJwksConfigurationDTO()).thenReturn(jwksConfigurationDTO);

            EnforcerConfig enforcerConfig = Mockito.mock(EnforcerConfig.class);
            ConfigHolder configHolder = Mockito.mock(ConfigHolder.class);
            configHolderDummy.when(ConfigHolder::getInstance).thenReturn(configHolder);
            Mockito.when(configHolder.getConfig()).thenReturn(enforcerConfig);
            JWTTransformer jwtTransformer = new DefaultJWTTransformer();
            Mockito.when(enforcerConfig.getJwtTransformer(issuer)).thenReturn(jwtTransformer);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidaterDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            try {
                oauth2Authenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for revoked tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), e.getMessage(),
                        APISecurityConstants.API_AUTH_INVALID_CREDENTIALS_MESSAGE);
            }
        } finally {
            RevokedTokenRedisClient.setRevokedTokens(new HashSet<>());
        }
    }

    @Test
    public void testCachedJWTValidator() throws APISecurityException, EnforcerException {

        String organization = "org1";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfF" +
                "CEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJ" +
                "UMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_" +
                "-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV" +
                "3YXkiLCJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZ" +
                "SI6InJlcXVlc3RDb3VudCIsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoiRGVmYXVsdEFwcGxpY2F0aW9uIiwiaWQiOjEsInV1aWQ" +
                "iOm51bGx9LCJzY29wZSI6ImFtX2FwcGxpY2F0aW9uX3Njb3BlIGRlZmF1bHQiLCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0" +
                "NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6e30sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOltdL" +
                "CJjb25zdW1lcktleSI6IlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLCJleHAiOjE1OTAzNDIzMTMsImlhdCI6MTU5MDMzO" +
                "DcxMywianRpIjoiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIn0." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() + 5000L);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);) {
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment)).thenReturn(jwtValidator);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = oauth2Authenticator.authenticate(requestContext);
            Assert.assertNotNull(authenticate);
            Assert.assertEquals(authenticate.getConsumerKey(), jwtValidationInfo.getConsumerKey());
            Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).getIfPresent(signature);
        }
    }

    @Test
    public void testNonJTIJWTValidator() throws APISecurityException, EnforcerException {

        String organization = "org1";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "SSQyg_VTxF5drIogztn2SyEK2wRE07wG6OW3tufD3vo";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9" +
                ".eyJpc3MiOiJodHRwczovL2xvY2FsaG9zdCIsImlhdCI6MTU5OTU0ODE3NCwiZXhwIjoxNjMxMDg0MTc0LC" +
                "JhdWQiOiJ3d3cuZXhhbXBsZS5jb20iLCJzdWIiOiJqcm9ja2V0QGV4YW1wbGUuY29tIiwiR2l2ZW5OYW1l" +
                "IjoiSm9obm55IiwiU3VybmFtZSI6IlJvY2tldCIsIkVtYWlsIjoianJvY2tldEBleGFtcGxlLmNvbSIsIl" +
                "JvbGUiOlsiTWFuYWdlciIsIlByb2plY3QgQWRtaW5pc3RyYXRvciJdfQ." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() + 5000L);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);) {
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment))
                    .thenReturn(jwtValidator);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = oauth2Authenticator.authenticate(requestContext);
            Assert.assertNotNull(authenticate);
            Assert.assertEquals(authenticate.getConsumerKey(), jwtValidationInfo.getConsumerKey());
            Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).getIfPresent(signature);
        }
    }

    @Test
    public void testExpiredJWTValidator() {

        String organization = "org1";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfFCEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJUMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV" +
                "3YXkiLCJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZ" +
                "SI6InJlcXVlc3RDb3VudCIsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoiRGVmYXVsdEFwcGxpY2F0aW9uIiwiaWQiOjEsInV1aWQ" +
                "iOm51bGx9LCJzY29wZSI6ImFtX2FwcGxpY2F0aW9uX3Njb3BlIGRlZmF1bHQiLCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0" +
                "NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6e30sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOltdL" +
                "CJjb25zdW1lcktleSI6IlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLCJleHAiOjE1OTAzNDIzMTMsImlhdCI6MTU5MDMzO" +
                "DcxMywianRpIjoiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIn0." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() - 5000L);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
        Mockito.when(invalidTokenCache.getIfPresent(signature)).thenReturn(null);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class)) {
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            try {
                oauth2Authenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for expired tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), e.getMessage(),
                        APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED_MESSAGE);
                Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).getIfPresent(signature);
                Mockito.verify(invalidTokenCache, Mockito.atLeast(1)).put(signature, true);
            }
        }
    }

    @Test
    public void testNoCacheExpiredJWTValidator() throws EnforcerException {

        String organization = "org1";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfFCEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJUMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV" +
                "3YXkiLCJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZ" +
                "SI6InJlcXVlc3RDb3VudCIsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoiRGVmYXVsdEFwcGxpY2F0aW9uIiwiaWQiOjEsInV1aWQ" +
                "iOm51bGx9LCJzY29wZSI6ImFtX2FwcGxpY2F0aW9uX3Njb3BlIGRlZmF1bHQiLCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0" +
                "NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6e30sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOltdL" +
                "CJjb25zdW1lcktleSI6IlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLCJleHAiOjE1OTAzNDIzMTMsImlhdCI6MTU5MDMzO" +
                "DcxMywianRpIjoiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIn0." + signature;
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(false);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() - 100);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setValidationCode(APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED);
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + jwt);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
        LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
        LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(null);
        Mockito.when(invalidTokenCache.getIfPresent(signature)).thenReturn(null);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);
        ) {
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).
                    thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment)).thenReturn(jwtValidator);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);

            try {
                oauth2Authenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for expired tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED_MESSAGE);
                Mockito.verify(invalidTokenCache, Mockito.atLeast(1)).put(signature, true);
            }
        }

    }

    @Test
    public void testTamperedPayloadJWTValidator() throws EnforcerException {

        String organization = "org2";
        String environment = "development";
        String issuer = "https://localhost:9443/oauth2/token";
        String signature = "sBgeoqJn0log5EZflj_G7ADvm6B3KQ9bdfFCEFVQS1U3oY9" +
                "-cqPwAPyOLLh95pdfjYjakkf1UtjPZjeIupwXnzg0SffIc704RoVlZocAx9Ns2XihjU6Imx2MbXq9ARmQxQkyGVkJUMTwZ8" +
                "-SfOnprfrhX2cMQQS8m2Lp7hcsvWFRGKxAKIeyUrbY4ihRIA5vOUrMBWYUx9Di1N7qdKA4S3e8O4KQX2VaZPBzN594c9TG" +
                "riiH8AuuqnrftfvidSnlRLaFJmko8-QZo8jDepwacaFhtcaPVVJFG4uYP-_-N6sqfxLw3haazPN0_xU0T1zJLPRLC5HPfZMJDMGp" +
                "EuSe9w";
        String jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".ewogICJhdWQiOiAiaHR0cDovL29yZy53c28yLmFwaW1ndC9nYXRld2F5IiwKICAic3ViIjogImFkbWluQGNhcmJ" +
                "vbi5zdXBlciIsCiAgImFwcGxpY2F0aW9uIjogewogICAgIm93bmVyIjogImFkbWluIiwKICAgICJ0aWVyUXVvdGFU" +
                "eXBlIjogInJlcXVlc3RDb3VudCIsCiAgICAidGllciI6ICJVbmxpbWl0ZWQiLAogICAgIm5hbWUiOiAiRGVmYXVsd" +
                "EFwcGxpY2F0aW9uIiwKICAgICJpZCI6IDEsCiAgICAidXVpZCI6IG51bGwKICB9LAogICJzY29wZSI6ICJhbV9hcHB" +
                "saWNhdGlvbl9zY29wZSBkZWZhdWx0IiwKICAiaXNzIjogImh0dHBzOi8vbG9jYWxob3N0Ojk0NDMvb2F1dGgyL3Rva" +
                "2VuIiwKICAidGllckluZm8iOiB7fSwKICAia2V5dHlwZSI6ICJQUk9EVUNUSU9OIiwKICAic3Vic2NyaWJlZEFQSXM" +
                "iOiBbXSwKICAiY29uc3VtZXJLZXkiOiAiWGdPMzk2SUhGSzdlRll5ZHJxUWU0SEtHejFrYSIsCiAgImV4cCI6IDQxMz" +
                "IzODM0NzcsCiAgImlhdCI6IDE1OTAzMzg3MTMsCiAgImp0aSI6ICJiODkzODc2OC0yM2ZkLTRkZWMtOGI3MC1iZWQ0N" +
                "WViN2MzM2QiCn0=." + signature;
        String tamperedJWT = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".ewogICJhdWQiOiAiaHR0cDovL29yZy53c28yLmFwaW1ndC9nYXRld2F5IiwKICAic3ViIjogImFkbWluQGNhcmJvbi5" +
                "zdXBlciIsCiAgImFwcGxpY2F0aW9uIjogewogICAgIm93bmVyIjogImFkbWluIiwKICAgICJ0aWVyUXVvdGFUeXBlIjo" +
                "gInJlcXVlc3RDb3VudCIsCiAgICAidGllciI6ICJVbmxpbWl0ZWQiLAogICAgIm5hbWUiOiAiRGVmYXVsdEFwcGxpY2F" +
                "0aW9uMiIsCiAgICAiaWQiOiAyLAogICAgInV1aWQiOiBudWxsCiAgfSwKICAic2NvcGUiOiAiYW1fYXBwbGljYXRpb25" +
                "fc2NvcGUgZGVmYXVsdCIsCiAgImlzcyI6ICJodHRwczovL2xvY2FsaG9zdDo5NDQzL29hdXRoMi90b2tlbiIsCiAgInR" +
                "pZXJJbmZvIjoge30sCiAgImtleXR5cGUiOiAiUFJPRFVDVElPTiIsCiAgInN1YnNjcmliZWRBUElzIjogW10sCiAgImN" +
                "vbnN1bWVyS2V5IjogIlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2EiLAogICJleHAiOiA0MTMyMzgzNDc3LAogICJ" +
                "pYXQiOiAxNTkwMzM4NzEzLAogICJqdGkiOiAiYjg5Mzg3NjgtMjNmZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIgp9." +
                signature;

        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(false);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() + 120000);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        Oauth2Authenticator oauth2Authenticator = new Oauth2Authenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        Oauth2AuthenticationConfig oauth2AuthenticationConfig = new Oauth2AuthenticationConfig();
        oauth2AuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setOauth2AuthenticationConfig(oauth2AuthenticationConfig);
        Mockito.when(resourceConfig.getAuthenticationConfig()).thenReturn(authenticationConfig);
        Mockito.when(resourceConfig.getMethod()).thenReturn(ResourceConfig.HttpMethods.GET);
        resourceConfigs.add(resourceConfig);
        Mockito.when(requestContext.getMatchedResourcePaths()).thenReturn(resourceConfigs);
        Map<String, String> headers = new HashMap<>();
        headers.put("Authorization", "Bearer " + tamperedJWT);
        Mockito.when(requestContext.getHeaders()).thenReturn(headers);
        Mockito.when(requestContext.getAuthenticationContext()).thenReturn(new AuthenticationContext());
        APIConfig apiConfig = Mockito.mock(APIConfig.class);
        Mockito.when(apiConfig.getName()).thenReturn("api1");
        Mockito.when(apiConfig.getEnvironment()).thenReturn(environment);
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataHolder> subscriptionDataHolderMockedStatic =
                     Mockito.mockStatic(SubscriptionDataHolder.class);
        ) {
            CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
            SubscriptionDataStore subscriptionDataStore = Mockito.mock(SubscriptionDataStore.class);
            SubscriptionDataHolder subscriptionDataHolder = Mockito.mock(SubscriptionDataHolder.class);
            subscriptionDataHolderMockedStatic.when(SubscriptionDataHolder::getInstance).thenReturn(subscriptionDataHolder);
            Mockito.when(subscriptionDataHolder.getSubscriptionDataStore(organization)).thenReturn(subscriptionDataStore);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).
                    thenReturn(cacheProvider);
            LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
            LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
            Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
            Mockito.when(invalidTokenCache.getIfPresent(signature)).thenReturn(null);
            Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
            Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization))
                    .thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, environment)).thenReturn(jwtValidator);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            try {
                oauth2Authenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for tampered tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), "Invalid JWT token");
                Mockito.verify(invalidTokenCache, Mockito.never()).put(signature, true);
                Mockito.verify(gatewayKeyCache, Mockito.never()).put(signature, true);
                Mockito.verify(gatewayKeyCache, Mockito.never()).invalidate(signature);
            }
        }
    }
}
