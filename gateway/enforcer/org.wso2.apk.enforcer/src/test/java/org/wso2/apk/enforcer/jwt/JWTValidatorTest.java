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
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
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
import org.wso2.apk.enforcer.commons.model.JWTAuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.config.dto.ExtendedTokenIssuerDto;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.security.KeyValidator;
import org.wso2.apk.enforcer.security.jwt.JWTAuthenticator;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStoreImpl;

public class JWTValidatorTest {

    @Test
    public void testJWTValidator() throws APISecurityException, EnforcerException {
        String organization = "org1";
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
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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

        try (MockedStatic<LogManager> logManagerDummy = Mockito.mockStatic(LogManager.class);
             MockedStatic<LogFactory> logFactoryDummy = Mockito.mockStatic(LogFactory.class);
                MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
                MockedStatic<ConfigHolder> configHolderDummy = Mockito.mockStatic(ConfigHolder.class);
                MockedStatic<SubscriptionDataStoreImpl> subscriptionDataStoreImplDummy =
                     Mockito.mockStatic(SubscriptionDataStoreImpl.class);
                MockedStatic<KeyValidator> keyValidaterDummy = Mockito.mockStatic(KeyValidator.class)) {
            Logger logger = Mockito.mock(Logger.class);
            logManagerDummy.when(() -> LogManager.getLogger(JWTAuthenticator.class)).thenReturn(logger);
            Log logger2 = Mockito.mock(Log.class);
            logFactoryDummy.when(() -> LogFactory.getLog(AbstractAPIMgtGatewayJWTGenerator.class)).thenReturn(logger2);
            ////
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

            SubscriptionDataStoreImpl subscriptionDataStore = Mockito.mock(SubscriptionDataStoreImpl.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer,
                    organization)).thenReturn(jwtValidator);
            subscriptionDataStoreImplDummy.when(SubscriptionDataStoreImpl::getInstance).thenReturn(subscriptionDataStore);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidaterDummy.when(()->KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = jwtAuthenticator.authenticate(requestContext);
            Assert.assertNotNull(authenticate);
            Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).put(signature, jwtValidationInfo);
        }
    }

    @Test
    public void testCachedJWTValidator() throws APISecurityException, EnforcerException {
        String organization = "org1";
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
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<SubscriptionDataStoreImpl> subscriptionDataStoreImplDummy = Mockito.mockStatic(SubscriptionDataStoreImpl.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class)
        ) {
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            SubscriptionDataStoreImpl subscriptionDataStore = Mockito.mock(SubscriptionDataStoreImpl.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer,
                    organization)).thenReturn(jwtValidator);
            subscriptionDataStoreImplDummy.when(SubscriptionDataStoreImpl::getInstance).thenReturn(subscriptionDataStore);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = jwtAuthenticator.authenticate(requestContext);
            Assert.assertNotNull(authenticate);
            Assert.assertEquals(authenticate.getConsumerKey(), jwtValidationInfo.getConsumerKey());
            Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).getIfPresent(signature);
        }
    }

    @Test
    public void testNonJTIJWTValidator() throws APISecurityException, EnforcerException {
        String organization = "org1";
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
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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
        Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class);
             MockedStatic<SubscriptionDataStoreImpl> subscriptionDataStoreImplDummy = Mockito.mockStatic(SubscriptionDataStoreImpl.class);) {
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            SubscriptionDataStoreImpl subscriptionDataStore = Mockito.mock(SubscriptionDataStoreImpl.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer,
                    organization)).thenReturn(jwtValidator);
            subscriptionDataStoreImplDummy.when(SubscriptionDataStoreImpl::getInstance).thenReturn(subscriptionDataStore);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            AuthenticationContext authenticate = jwtAuthenticator.authenticate(requestContext);
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
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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
                jwtAuthenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for expired tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), e.getMessage(), APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED_MESSAGE);
                Mockito.verify(gatewayKeyCache, Mockito.atLeast(1)).getIfPresent(signature);
                Mockito.verify(invalidTokenCache, Mockito.atLeast(1)).put(signature, true);
            }
        }
    }

    @Test
    public void testNoCacheExpiredJWTValidator() throws EnforcerException {
        String organization = "org1";
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
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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
        Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(null);
        Mockito.when(invalidTokenCache.getIfPresent(signature)).thenReturn(null);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<SubscriptionDataStoreImpl> subscriptionDataStoreImplDummy =
                     Mockito.mockStatic(SubscriptionDataStoreImpl.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class)
        ) {
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).
                    thenReturn(cacheProvider);
            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            SubscriptionDataStoreImpl subscriptionDataStore = Mockito.mock(SubscriptionDataStoreImpl.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer, organization)).thenReturn(jwtValidator);
            subscriptionDataStoreImplDummy.when(SubscriptionDataStoreImpl::getInstance).thenReturn(subscriptionDataStore);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);

            try {
                jwtAuthenticator.authenticate(requestContext);
                Assert.fail("Authentication should fail for expired tokens");
            } catch (APISecurityException e) {
                Assert.assertEquals(e.getMessage(), APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED_MESSAGE);
                Mockito.verify(invalidTokenCache, Mockito.atLeast(1)).put(signature, true);
            }
        }

    }

    @Test
    public void testTamperedPayloadJWTValidator() throws EnforcerException {
        String organization = "org1";
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
        String tamperedJWT = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5UZG1aak00WkRrM05qWTBZemM1T" +
                "W1abU9EZ3dNVEUzTVdZd05ERTVNV1JsWkRnNE56YzRaQT09In0" +
                ".ewogICJhdWQiOiAiaHR0cDovL29yZy53c28yLmFwaW1ndC9nYXRld2F5IiwKICAic3ViIjogImFkbWluQGNhcm" +
                "Jvbi5zdXBlciIsCiAgImFwcGxpY2F0aW9uIjogewogICAgIm93bmVyIjogImFkbWluIiwKICAgICJ0aWVyUXVvd" +
                "GFUeXBlIjogInJlcXVlc3RDb3VudCIsCiAgICAidGllciI6ICJVbmxpbWl0ZWQiLAogICAgIm5hbWUiOiAiRGVm" +
                "YXVsdEFwcGxpY2F0aW9uMiIsCiAgICAiaWQiOiAyLAogICAgInV1aWQiOiBudWxsCiAgfSwKICAic2NvcGUiOiA" +
                "iYW1fYXBwbGljYXRpb25fc2NvcGUgZGVmYXVsdCIsCiAgImlzcyI6ICJodHRwczovL2xvY2FsaG9zdDo5NDQzL2" +
                "9hdXRoMi90b2tlbiIsCiAgInRpZXJJbmZvIjoge30sCiAgImtleXR5cGUiOiAiUFJPRFVDVElPTiIsCiAgInN1Y" +
                "nNjcmliZWRBUElzIjogW10sCiAgImNvbnN1bWVyS2V5IjogIlhnTzM5NklIRks3ZUZZeWRycVFlNEhLR3oxa2Ei" +
                "LAogICJleHAiOiAxNTkwMzQyMzEzLAogICJpYXQiOiAxNTkwMzM4NzEzLAogICJqdGkiOiAiYjg5Mzg3NjgtMjN" +
                "mZC00ZGVjLThiNzAtYmVkNDVlYjdjMzNkIgp9." + signature;

        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        jwtValidationInfo.setValid(false);
        jwtValidationInfo.setExpiryTime(System.currentTimeMillis() - 100);
        jwtValidationInfo.setConsumerKey(UUID.randomUUID().toString());
        jwtValidationInfo.setUser("user1");
        jwtValidationInfo.setKeyManager("Default");
        jwtValidationInfo.setToken(jwt);
        jwtValidationInfo.setIdentifier(signature);

        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        JWTAuthenticator jwtAuthenticator = new JWTAuthenticator(jwtConfigurationDto, true);
        RequestContext requestContext = Mockito.mock(RequestContext.class);

        ArrayList<ResourceConfig> resourceConfigs = new ArrayList<>();
        ResourceConfig resourceConfig = Mockito.mock(ResourceConfig.class);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader("Authorization");
        authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
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
        Mockito.when(apiConfig.getOrganizationId()).thenReturn(organization);
        Mockito.when(requestContext.getMatchedAPI()).thenReturn(apiConfig);
        try (MockedStatic<CacheProviderUtil> cacheProviderUtilDummy = Mockito.mockStatic(CacheProviderUtil.class);
             MockedStatic<SubscriptionDataStoreImpl> subscriptionDataStoreImplDummy =
                     Mockito.mockStatic(SubscriptionDataStoreImpl.class);
             MockedStatic<KeyValidator> keyValidatorDummy = Mockito.mockStatic(KeyValidator.class)
        ) {
            CacheProvider cacheProvider = Mockito.mock(CacheProvider.class);
            LoadingCache gatewayKeyCache = Mockito.mock(LoadingCache.class);
            LoadingCache invalidTokenCache = Mockito.mock(LoadingCache.class);
            Mockito.when(gatewayKeyCache.getIfPresent(signature)).thenReturn(jwtValidationInfo);
            Mockito.when(invalidTokenCache.getIfPresent(signature)).thenReturn(null);
            Mockito.when(cacheProvider.getGatewayKeyCache()).thenReturn(gatewayKeyCache);
            Mockito.when(cacheProvider.getInvalidTokenCache()).thenReturn(invalidTokenCache);
            cacheProviderUtilDummy.when(() -> CacheProviderUtil.getOrganizationCache(organization)).thenReturn(cacheProvider);

            JWTValidator jwtValidator = Mockito.mock(JWTValidator.class);
            SubscriptionDataStoreImpl subscriptionDataStore = Mockito.mock(SubscriptionDataStoreImpl.class);
            Mockito.when(subscriptionDataStore.getJWTValidatorByIssuer(issuer,
                    organization)).thenReturn(jwtValidator);

            subscriptionDataStoreImplDummy.when(SubscriptionDataStoreImpl::getInstance).thenReturn(subscriptionDataStore);
            Mockito.when(jwtValidator.validateToken(Mockito.eq(jwt), Mockito.any())).thenReturn(jwtValidationInfo);
            keyValidatorDummy.when(() -> KeyValidator.validateScopes(Mockito.any())).thenReturn(true);
            try {
                jwtAuthenticator.authenticate(requestContext);
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
