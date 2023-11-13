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

package org.wso2.apk.enforcer.subscription;

import feign.Feign;
import feign.gson.GsonDecoder;
import feign.gson.GsonEncoder;
import feign.slf4j.Slf4jLogger;
import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.dto.ClaimMappingDto;
import org.wso2.apk.enforcer.commons.dto.JWKSConfigurationDTO;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.dto.ExtendedTokenIssuerDto;
import org.wso2.apk.enforcer.constants.Constants;
import org.wso2.apk.enforcer.discovery.ApiListDiscoveryClient;
import org.wso2.apk.enforcer.discovery.JWTIssuerDiscoveryClient;
import org.wso2.apk.enforcer.discovery.subscription.APIs;
import org.wso2.apk.enforcer.discovery.subscription.Certificate;
import org.wso2.apk.enforcer.discovery.subscription.JWTIssuer;
import org.wso2.apk.enforcer.models.API;
import org.wso2.apk.enforcer.models.Application;
import org.wso2.apk.enforcer.models.ApplicationKeyMapping;
import org.wso2.apk.enforcer.models.ApplicationMapping;
import org.wso2.apk.enforcer.models.SubscribedAPI;
import org.wso2.apk.enforcer.models.Subscription;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;
import org.wso2.apk.enforcer.util.ApacheFeignHttpClient;
import org.wso2.apk.enforcer.util.FilterUtils;
import org.wso2.apk.enforcer.util.TLSUtils;

import java.io.IOException;
import java.security.cert.CertificateException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Implementation of the subscription data store.
 */
public class SubscriptionDataStoreImpl implements SubscriptionDataStore {

    private static final Logger log = LogManager.getLogger(SubscriptionDataStoreImpl.class);
    private static final SubscriptionDataStoreImpl instance = new SubscriptionDataStoreImpl();

    public static final String DELEM_PERIOD = ":";

    // Maps for keeping Subscription related details.
    private Map<String, ApplicationKeyMapping> applicationKeyMappingMap;
    private Map<String, ApplicationMapping> applicationMappingMap;
    private Map<String, Application> applicationMap;
    private Map<String, API> apiMap;
    private Map<String, Subscription> subscriptionMap;

    private Map<String, Map<String, JWTValidator>> jwtValidatorMap;
    SubscriptionValidationDataRetrievalRestClient subscriptionValidationDataRetrievalRestClient;

    SubscriptionDataStoreImpl() {

    }

    public static SubscriptionDataStoreImpl getInstance() {

        return instance;
    }

    public void initializeStore() {

        String commonControllerHost = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHost();
        String commonControllerHostname = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHostname();
        int commonControllerRestPort =
                Integer.parseInt(ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerRestPort());
        subscriptionValidationDataRetrievalRestClient = Feign.builder()
                .encoder(new GsonEncoder())
                .decoder(new GsonDecoder())
                .logger(new Slf4jLogger())
                .client(new ApacheFeignHttpClient(FilterUtils.getMutualSSLHttpClient("https",
                        Arrays.asList(commonControllerHost, commonControllerHostname))))
                .target(SubscriptionValidationDataRetrievalRestClient.class,
                        "https://" + commonControllerHost + ":" + commonControllerRestPort);
        this.applicationKeyMappingMap = new ConcurrentHashMap<>();
        this.applicationMap = new ConcurrentHashMap<>();
        this.apiMap = new ConcurrentHashMap<>();
        this.subscriptionMap = new ConcurrentHashMap<>();
        this.applicationMappingMap = new ConcurrentHashMap<>();
        this.jwtValidatorMap = new ConcurrentHashMap<>();
        initializeLoadingTasks();
    }

    @Override
    public Application getApplicationById(String appUUID) {

        return applicationMap.get(appUUID);
    }

    @Override
    public API getApiByContextAndVersion(String uuid) {

        return apiMap.get(uuid);
    }

    @Override
    public Subscription getSubscriptionById(String appId, String apiId) {

        return subscriptionMap.get(SubscriptionDataStoreUtil.getSubscriptionCacheKey(appId, apiId));
    }

    private void initializeLoadingTasks() {
        loadSubscriptions();
        loadApplications();
        loadApplicationMappings();
        loadApplicationKeyMappings();
        ApiListDiscoveryClient.getInstance().watchApiList();
        JWTIssuerDiscoveryClient.getInstance().watchJWTIssuers();
    }

    private void loadApplicationKeyMappings() {
        new Thread(() -> {
            ApplicationKeyMappingDtoList applicationKeyMappings =
                    subscriptionValidationDataRetrievalRestClient.getAllApplicationKeyMappings();
            addApplicationKeyMappings(applicationKeyMappings.getList());
        }).start();

    }

    private void loadApplicationMappings() {
        new Thread(() -> {
            ApplicationMappingDtoList applicationMappings = subscriptionValidationDataRetrievalRestClient
                    .getAllApplicationMappings();
            addApplicationMappings(applicationMappings.getList());
        }).start();

    }

    private void loadApplications(){
        new Thread(() -> {
            ApplicationListDto applications = subscriptionValidationDataRetrievalRestClient.getAllApplications();
            addApplications(applications.getList());
        }).start();
    }
    private void loadSubscriptions() {

        new Thread(() -> {
            SubscriptionListDto subscriptions = subscriptionValidationDataRetrievalRestClient.getAllSubscriptions();
            addSubscriptions(subscriptions.getList());
        }).start();
    }

    public void addSubscriptions(List<SubscriptionDto> subscriptionList) {

        Map<String, Subscription> newSubscriptionMap = new ConcurrentHashMap<>();

        for (SubscriptionDto subscription : subscriptionList) {
            SubscribedAPI subscribedAPI = new SubscribedAPI();
            subscribedAPI.setName(subscription.getSubscribedApi().getName());
            subscribedAPI.setVersion(subscription.getSubscribedApi().getVersion());

            Subscription newSubscription = new Subscription();
            newSubscription.setSubscriptionId(subscription.getUuid());
            newSubscription.setSubscriptionStatus(subscription.getSubStatus());
            newSubscription.setOrganization(subscription.getOrganization());
            newSubscription.setSubscribedApi(subscribedAPI);
            newSubscriptionMap.put(newSubscription.getCacheKey(), newSubscription);
        }

        if (log.isDebugEnabled()) {
            log.debug("Total Subscriptions in new cache: {}", newSubscriptionMap.size());
        }
        this.subscriptionMap = newSubscriptionMap;
    }

    public void addApplications(List<ApplicationDto> applicationList) {

        Map<String, Application> newApplicationMap = new ConcurrentHashMap<>();

        for (ApplicationDto application: applicationList) {
            Application newApplication = new Application();
            newApplication.setUUID(application.getUuid());
            newApplication.setName(application.getName());
            newApplication.setOwner(application.getOwner());
            application.getAttributes().forEach(newApplication::addAttribute);

            newApplicationMap.put(newApplication.getCacheKey(), newApplication);
        }
        if (log.isDebugEnabled()) {
            log.debug("Total Applications in new cache: {}", newApplicationMap.size());
        }
        this.applicationMap = newApplicationMap;
    }

    public void addApis(List<APIs> apisList) {

        Map<String, API> newApiMap = new ConcurrentHashMap<>();

        for (APIs api : apisList) {
            API newApi = new API();
            // newApi.setApiId(Integer.parseInt(api.getApiId()));
            newApi.setApiName(api.getName());
            newApi.setApiProvider(api.getProvider());
            newApi.setApiType(api.getApiType());
            newApi.setApiVersion(api.getVersion());
            newApi.setContext(api.getBasePath());
            newApi.setApiTier(api.getPolicy());
            newApi.setApiUUID(api.getUuid());
            newApi.setLcState(api.getLcState());
            newApiMap.put(newApi.getCacheKey(), newApi);
        }
        if (log.isDebugEnabled()) {
            log.debug("Total Apis in new cache: {}", newApiMap.size());
        }
        this.apiMap = newApiMap;
    }

    public void addApplicationKeyMappings(
            List<ApplicationKeyMappingDTO> applicationKeyMappingList) {

        Map<String, ApplicationKeyMapping> newApplicationKeyMappingMap = new ConcurrentHashMap<>();

        for (ApplicationKeyMappingDTO applicationKeyMapping :
                applicationKeyMappingList) {
            ApplicationKeyMapping mapping = new ApplicationKeyMapping();
            mapping.setApplicationUUID(applicationKeyMapping.getApplicationUUID());
            mapping.setSecurityScheme(applicationKeyMapping.getSecurityScheme());
            mapping.setApplicationIdentifier(applicationKeyMapping.getApplicationIdentifier());
            mapping.setKeyType(applicationKeyMapping.getKeyType());
            mapping.setEnvId(applicationKeyMapping.getEnvID());
            newApplicationKeyMappingMap.put(mapping.getCacheKey(), mapping);
        }
        if (log.isDebugEnabled()) {
            log.debug("Total Application Key Mappings in new cache: {}", newApplicationKeyMappingMap.size());
        }
        this.applicationKeyMappingMap = newApplicationKeyMappingMap;
    }

    public void addApplicationMappings(
            List<ApplicationMappingDto> applicationMappingList) {

        Map<String, ApplicationMapping> newApplicationMappingMap = new ConcurrentHashMap<>();

        for (ApplicationMappingDto applicationMapping :
                applicationMappingList) {
            ApplicationMapping appMapping = new ApplicationMapping();
            appMapping.setUuid(applicationMapping.getUuid());
            appMapping.setApplicationRef(applicationMapping.getApplicationRef());
            appMapping.setSubscriptionRef(applicationMapping.getSubscriptionRef());
            newApplicationMappingMap.put(appMapping.getCacheKey(), appMapping);
        }
        if (log.isDebugEnabled()) {
            log.debug("Total Application Mappings in new cache: {}", newApplicationMappingMap.size());
        }
        this.applicationMappingMap = newApplicationMappingMap;
    }

    @Override
    public API getMatchingAPI(String context, String version) {

        for (API api : apiMap.values()) {
            if (StringUtils.isNotEmpty(context) && StringUtils.isNotEmpty(version)) {
                if (api.getContext().equals(context) && api.getApiVersion().equals(version)) {
                    return api;
                }
            }
        }
        return null;
    }

    @Override
    public ApplicationKeyMapping getMatchingApplicationKeyMapping(String applicationIdentifier, String keyType,
                                                                  String securityScheme) {

        for (ApplicationKeyMapping applicationKeyMapping : applicationKeyMappingMap.values()) {
            boolean isApplicationIdentifierMatching = false;
            boolean isSecuritySchemeMatching = false;
            boolean isKeyTypeMatching = false;

            if (StringUtils.isNotEmpty(applicationIdentifier)) {
                if (applicationKeyMapping.getApplicationIdentifier().equals(applicationIdentifier)) {
                    isApplicationIdentifierMatching = true;
                }
            }
            if (StringUtils.isNotEmpty(securityScheme)) {
                if (applicationKeyMapping.getSecurityScheme().equals(securityScheme)) {
                    isSecuritySchemeMatching = true;
                }
            }
            if (StringUtils.isNotEmpty(keyType)) {
                if (applicationKeyMapping.getKeyType().equals(keyType)) {
                    isKeyTypeMatching = true;
                }
            }

            if (isApplicationIdentifierMatching && isSecuritySchemeMatching && isKeyTypeMatching) {
                return applicationKeyMapping;
            }
        }
        return null;
    }

    @Override
    public ApplicationMapping getMatchingApplicationMapping(String uuid) {

        for (ApplicationMapping applicationMapping : applicationMappingMap.values()) {
            if (StringUtils.isNotEmpty(uuid)) {
                if (applicationMapping.getApplicationRef().equals(uuid)) {
                    return applicationMapping;
                }
            }
        }
        return null;
    }

    @Override
    public Application getMatchingApplication(String uuid) {

        for (Application application : applicationMap.values()) {
            if (StringUtils.isNotEmpty(uuid)) {
                if (application.getUUID().equals(uuid)) {
                    return application;
                }
            }
        }
        return null;
    }

    @Override
    public Subscription getMatchingSubscription(String uuid) {

        for (Subscription subscription : subscriptionMap.values()) {
            if (StringUtils.isNotEmpty(uuid)) {
                if (subscription.getSubscriptionId().equals(uuid)) {
                    return subscription;
                }
            }
        }
        return null;
    }

    @Override
    public void addJWTIssuers(List<JWTIssuer> jwtIssuers) {

        Map<String, Map<String, JWTValidator>> jwtValidatorMap = new ConcurrentHashMap<>();
        for (JWTIssuer jwtIssuer : jwtIssuers) {
            try {
                ExtendedTokenIssuerDto tokenIssuerDto = new ExtendedTokenIssuerDto(jwtIssuer.getIssuer());
                tokenIssuerDto.setName(jwtIssuer.getName());
                tokenIssuerDto.setConsumerKeyClaim(jwtIssuer.getConsumerKeyClaim());
                tokenIssuerDto.setScopesClaim(jwtIssuer.getScopesClaim());
                Certificate certificate = jwtIssuer.getCertificate();
                if (StringUtils.isNotEmpty(certificate.getJwks().getUrl())) {
                    JWKSConfigurationDTO jwksConfigurationDTO = new JWKSConfigurationDTO();
                    if (StringUtils.isNotEmpty(certificate.getJwks().getTls())) {
                        java.security.cert.Certificate tlsCertificate = TLSUtils
                                .getCertificateFromContent(certificate.getJwks().getTls());
                        jwksConfigurationDTO.setCertificate(tlsCertificate);
                    }
                    jwksConfigurationDTO.setUrl(certificate.getJwks().getUrl());
                    jwksConfigurationDTO.setEnabled(true);
                    tokenIssuerDto.setJwksConfigurationDTO(jwksConfigurationDTO);
                }
                if (StringUtils.isNotEmpty(certificate.getCertificate())) {
                    java.security.cert.Certificate signingCertificate = TLSUtils
                            .getCertificateFromContent(certificate.getCertificate());
                    tokenIssuerDto.setCertificate(signingCertificate);
                }
                Map<String, String> claimMappingMap = jwtIssuer.getClaimMappingMap();
                Map<String, ClaimMappingDto> claimMappingDtos = new HashMap<>();
                claimMappingMap.forEach((remoteClaim, localClaim) -> {
                    claimMappingDtos.put(remoteClaim, new ClaimMappingDto(remoteClaim, localClaim));
                });
                tokenIssuerDto.setClaimMappings(claimMappingDtos);
                JWTValidator jwtValidator = new JWTValidator(tokenIssuerDto);
                Map<String, JWTValidator> orgBasedJWTValidatorMap = new ConcurrentHashMap<>();
                if (jwtValidatorMap.containsKey(jwtIssuer.getOrganization())) {
                    orgBasedJWTValidatorMap = jwtValidatorMap.get(jwtIssuer.getOrganization());
                }

                List<String> environments = getEnvironments(jwtIssuer);
                for (String environment : environments) {
                    String mapKey = getMapKey(environment, jwtIssuer.getIssuer());
                    orgBasedJWTValidatorMap.put(mapKey, jwtValidator);
                }

                jwtValidatorMap.put(jwtIssuer.getOrganization(), orgBasedJWTValidatorMap);
                this.jwtValidatorMap = jwtValidatorMap;
            } catch (EnforcerException | CertificateException | IOException e) {
                log.error("Error occurred while configuring JWT Validator for issuer " + jwtIssuer.getIssuer(), e);
            }
        }
    }

    @Override
    public JWTValidator getJWTValidatorByIssuer(String issuer, String organization, String environment) {

        Map<String, JWTValidator> orgBaseJWTValidators = jwtValidatorMap.get(organization);

        if (orgBaseJWTValidators != null) {

            String mapKey = getMapKey(Constants.DEFAULT_ALL_ENVIRONMENTS_TOKEN_ISSUER, issuer);
            JWTValidator jwtValidator = orgBaseJWTValidators.get(mapKey);
            if (jwtValidator != null) {
                return jwtValidator;
            }

            mapKey = getMapKey(environment, issuer);
            return orgBaseJWTValidators.get(mapKey);
        }

        return null;
    }

    private List<String> getEnvironments(JWTIssuer jwtIssuer) {

        List<String> environmentsList = new ArrayList<>();
        int environmentCount = jwtIssuer.getEnvironmentsCount();

        if (environmentCount > 0) {
            for (int i = 0; i < environmentCount; i++) {
                environmentsList.add(jwtIssuer.getEnvironments(i));
            }
        } else {
            environmentsList.add(Constants.DEFAULT_ALL_ENVIRONMENTS_TOKEN_ISSUER);
        }
        return environmentsList;
    }

    private String getMapKey(String environment, String issuer) {

        return environment + DELEM_PERIOD + issuer;
    }

}
