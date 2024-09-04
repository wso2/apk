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

import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.dto.ClaimMappingDto;
import org.wso2.apk.enforcer.commons.dto.JWKSConfigurationDTO;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.dto.ExtendedTokenIssuerDto;
import org.wso2.apk.enforcer.constants.Constants;
import org.wso2.apk.enforcer.discovery.subscription.Certificate;
import org.wso2.apk.enforcer.discovery.subscription.JWTIssuer;
import org.wso2.apk.enforcer.models.Application;
import org.wso2.apk.enforcer.models.ApplicationKeyMapping;
import org.wso2.apk.enforcer.models.ApplicationMapping;
import org.wso2.apk.enforcer.models.SubscribedAPI;
import org.wso2.apk.enforcer.models.Subscription;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;
import org.wso2.apk.enforcer.util.TLSUtils;

import java.io.IOException;
import java.security.cert.CertificateException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Implementation of the subscription data store.
 */
public class SubscriptionDataStoreImpl implements SubscriptionDataStore {

    private static final Logger log = LogManager.getLogger(SubscriptionDataStoreImpl.class);
    private static final SubscriptionDataStoreImpl instance = new SubscriptionDataStoreImpl();

    public static final String DELEM_PERIOD = ":";

    // Maps for keeping Subscription related details.
    private Map<String, ApplicationKeyMapping> applicationKeyMappingMap = new ConcurrentHashMap<>();
    private Map<String, Map<String, ApplicationMapping>> applicationMappingMap = new ConcurrentHashMap<>();
    private Map<String, Application> applicationMap = new ConcurrentHashMap<>();
    private Map<String, Subscription> subscriptionMap = new ConcurrentHashMap<>();
    private Map<String, JWTValidator> jwtValidatorMap = new ConcurrentHashMap<>();

    SubscriptionDataStoreImpl() {

    }

    public static SubscriptionDataStoreImpl getInstance() {

        return instance;
    }

    @Override
    public Application getApplicationById(String appUUID) {

        return applicationMap.get(appUUID);
    }

    @Override
    public Subscription getSubscriptionById(String appId, String apiId) {

        return subscriptionMap.get(SubscriptionDataStoreUtil.getSubscriptionCacheKey(appId, apiId));
    }

    public void addSubscriptions(List<SubscriptionDto> subscriptionList) {

        Map<String, Subscription> newSubscriptionMap = new ConcurrentHashMap<>();

        for (SubscriptionDto subscription : subscriptionList) {
            SubscribedAPI subscribedAPI = new SubscribedAPI(subscription.getSubscribedApi());
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

        for (ApplicationDto application : applicationList) {
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

    public void addApplicationKeyMappings(List<ApplicationKeyMappingDTO> applicationKeyMappingList) {

        Map<String, ApplicationKeyMapping> newApplicationKeyMappingMap = new ConcurrentHashMap<>();

        for (ApplicationKeyMappingDTO applicationKeyMapping : applicationKeyMappingList) {
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

    public void addApplicationMappings(List<ApplicationMappingDto> applicationMappingList) {

        Map<String, Map<String, ApplicationMapping>> newApplicationMappingMap = new ConcurrentHashMap<>();
        for (ApplicationMappingDto applicationMapping : applicationMappingList) {
            ApplicationMapping appMapping = new ApplicationMapping();
            appMapping.setUuid(applicationMapping.getUuid());
            appMapping.setApplicationUUID(applicationMapping.getApplicationRef());
            appMapping.setSubscriptionUUID(applicationMapping.getSubscriptionRef());
            appMapping.setOrganization(applicationMapping.getOrganizationId());
            Map<String, ApplicationMapping> applicationMappings = new HashMap<>();
            if (newApplicationMappingMap.containsKey(appMapping.getApplicationUUID())) {
                applicationMappings = newApplicationMappingMap.get(appMapping.getApplicationUUID());
            }
            applicationMappings.put(appMapping.getCacheKey(), appMapping);
            newApplicationMappingMap.put(appMapping.getApplicationUUID(), applicationMappings);
        }
        if (log.isDebugEnabled()) {
            log.debug("Total Application Mappings in new cache: {}", newApplicationMappingMap.size());
        }
        this.applicationMappingMap = newApplicationMappingMap;
    }

    @Override
    public ApplicationKeyMapping getMatchingApplicationKeyMapping(String applicationIdentifier, String keyType,
                                                                  String securityScheme, String envType) {

        String cacheKey = SubscriptionDataStoreUtil.getApplicationKeyMappingCacheKey(applicationIdentifier, keyType,
                securityScheme, envType);
        return applicationKeyMappingMap.get(cacheKey);
    }

    @Override
    public Set<ApplicationMapping> getMatchingApplicationMappings(String uuid) {

        if (StringUtils.isNotEmpty(uuid)) {
            if (applicationMappingMap.containsKey(uuid)) {
                return new HashSet<>(applicationMappingMap.get(uuid).values());
            }
        }
        return new HashSet<>();
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

        if (StringUtils.isNotEmpty(uuid)) {
            return subscriptionMap.get(uuid);
        }
        return null;
    }

    @Override
    public void addJWTIssuers(List<JWTIssuer> jwtIssuers) {

        Map<String, JWTValidator> jwtValidatorMap = new ConcurrentHashMap<>();
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
                        java.security.cert.Certificate tlsCertificate =
                                TLSUtils.getCertificateFromContent(certificate.getJwks().getTls());
                        jwksConfigurationDTO.setCertificate(tlsCertificate);
                    }
                    jwksConfigurationDTO.setUrl(certificate.getJwks().getUrl());
                    jwksConfigurationDTO.setEnabled(true);
                    tokenIssuerDto.setJwksConfigurationDTO(jwksConfigurationDTO);
                }
                if (StringUtils.isNotEmpty(certificate.getCertificate())) {
                    java.security.cert.Certificate signingCertificate =
                            TLSUtils.getCertificateFromContent(certificate.getCertificate());
                    tokenIssuerDto.setCertificate(signingCertificate);
                }
                Map<String, String> claimMappingMap = jwtIssuer.getClaimMappingMap();
                Map<String, ClaimMappingDto> claimMappingDtos = new HashMap<>();
                claimMappingMap.forEach((remoteClaim, localClaim) -> claimMappingDtos.put(remoteClaim,
                        new ClaimMappingDto(remoteClaim, localClaim)));
                tokenIssuerDto.setClaimMappings(claimMappingDtos);
                JWTValidator jwtValidator = new JWTValidator(tokenIssuerDto);
                List<String> environments = getEnvironments(jwtIssuer);
                for (String environment : environments) {
                    String mapKey = getMapKey(environment, jwtIssuer.getIssuer());
                    jwtValidatorMap.put(mapKey, jwtValidator);
                }
                this.jwtValidatorMap = jwtValidatorMap;
            } catch (EnforcerException | CertificateException | IOException e) {
                log.error("Error occurred while configuring JWT Validator for issuer " + jwtIssuer.getIssuer(), e);
            }
        }
    }

    @Override
    public JWTValidator getJWTValidatorByIssuer(String issuer, String environment) {

        String mapKey = getMapKey(Constants.DEFAULT_ALL_ENVIRONMENTS_TOKEN_ISSUER, issuer);
        JWTValidator jwtValidator = jwtValidatorMap.get(mapKey);
        if (jwtValidator != null) {
            return jwtValidator;
        }
        mapKey = getMapKey(environment, issuer);
        return jwtValidatorMap.get(mapKey);
    }

    @Override
    public void addApplication(org.wso2.apk.enforcer.discovery.subscription.Application application) {

        Application resolvedApplication = new Application();
        resolvedApplication.setName(application.getName());
        resolvedApplication.setOwner(application.getOwner());
        resolvedApplication.setUUID(application.getUuid());
        resolvedApplication.setOrganization(application.getOrganization());
        resolvedApplication.setAttributes(application.getAttributesMap());
        if (applicationMap.containsKey(resolvedApplication.getUuid())) {
            applicationMap.replace(resolvedApplication.getUuid(), resolvedApplication);
        } else {
            applicationMap.put(resolvedApplication.getUuid(), resolvedApplication);
        }
    }

    @Override
    public void addSubscription(org.wso2.apk.enforcer.discovery.subscription.Subscription subscription) {

        Subscription resolvedSubscription = new Subscription();
        resolvedSubscription.setSubscriptionId(subscription.getUuid());
        resolvedSubscription.setSubscriptionStatus(subscription.getSubStatus());
        resolvedSubscription.setOrganization(subscription.getOrganization());
        resolvedSubscription.setSubscribedApi(new SubscribedAPI(subscription.getSubscribedApi()));
        resolvedSubscription.setRatelimitTier(subscription.getRatelimitTier());
        System.out.println(subscription.getRatelimitTier());
        if (subscriptionMap.containsKey(resolvedSubscription.getSubscriptionId())) {
            subscriptionMap.replace(resolvedSubscription.getSubscriptionId(), resolvedSubscription);
        } else {
            subscriptionMap.put(resolvedSubscription.getSubscriptionId(), resolvedSubscription);
        }
    }

    @Override
    public void addApplicationMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationMapping applicationMapping) {

        ApplicationMapping resolvedApplicationMapping = new ApplicationMapping();
        resolvedApplicationMapping.setUuid(applicationMapping.getUuid());
        resolvedApplicationMapping.setApplicationUUID(applicationMapping.getApplicationRef());
        resolvedApplicationMapping.setSubscriptionUUID(applicationMapping.getSubscriptionRef());
        resolvedApplicationMapping.setOrganization(applicationMapping.getOrganization());
        Map<String, ApplicationMapping> applicationMappings = new HashMap<>();
        if (applicationMappingMap.containsKey(resolvedApplicationMapping.getApplicationUUID())) {
            applicationMappings = applicationMappingMap.get(resolvedApplicationMapping.getApplicationUUID());
        }
        applicationMappings.put(resolvedApplicationMapping.getCacheKey(), resolvedApplicationMapping);
        applicationMappingMap.put(resolvedApplicationMapping.getApplicationUUID(), applicationMappings);
    }

    @Override
    public void addApplicationKeyMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping applicationKeyMapping) {

        ApplicationKeyMapping resolvedApplicationKeyMapping = new ApplicationKeyMapping();
        resolvedApplicationKeyMapping.setApplicationUUID(applicationKeyMapping.getApplicationUUID());
        resolvedApplicationKeyMapping.setSecurityScheme(applicationKeyMapping.getSecurityScheme());
        resolvedApplicationKeyMapping.setApplicationIdentifier(applicationKeyMapping.getApplicationIdentifier());
        resolvedApplicationKeyMapping.setKeyType(applicationKeyMapping.getKeyType());
        resolvedApplicationKeyMapping.setEnvId(applicationKeyMapping.getEnvID());
        Iterator<Map.Entry<String, ApplicationKeyMapping>> iterator = applicationKeyMappingMap.entrySet().iterator();
        while (iterator.hasNext()) {
            Map.Entry<String, ApplicationKeyMapping> cachedApplicationKeyMapping = iterator.next();
            ApplicationKeyMapping value = cachedApplicationKeyMapping.getValue();
            if (value.getApplicationIdentifier().equals(resolvedApplicationKeyMapping.getApplicationIdentifier()) && value.getSecurityScheme().equals(resolvedApplicationKeyMapping.getSecurityScheme()) && value.getKeyType().equals(resolvedApplicationKeyMapping.getKeyType()) && value.getEnvId().equals(resolvedApplicationKeyMapping.getEnvId()) && value.getApplicationUUID().equals(resolvedApplicationKeyMapping.getApplicationUUID())) {
                iterator.remove();
            }
        }
        applicationKeyMappingMap.put(resolvedApplicationKeyMapping.getCacheKey(), resolvedApplicationKeyMapping);
    }

    @Override
    public void removeApplicationMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationMapping applicationMapping) {

        ApplicationMapping resolvedApplicationMapping = new ApplicationMapping();
        resolvedApplicationMapping.setUuid(applicationMapping.getUuid());
        resolvedApplicationMapping.setApplicationUUID(applicationMapping.getApplicationRef());
        resolvedApplicationMapping.setSubscriptionUUID(applicationMapping.getSubscriptionRef());
        resolvedApplicationMapping.setOrganization(applicationMapping.getOrganization());
        Map<String, ApplicationMapping> applicationMappings =
                applicationMappingMap.get(applicationMapping.getApplicationRef());
        if (applicationMappings != null) {
            applicationMappings.remove(resolvedApplicationMapping.getCacheKey());
        }
    }

    @Override
    public void removeApplicationKeyMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping applicationKeyMapping) {

        ApplicationKeyMapping resolvedApplicationKeyMapping = new ApplicationKeyMapping();
        resolvedApplicationKeyMapping.setApplicationUUID(applicationKeyMapping.getApplicationUUID());
        resolvedApplicationKeyMapping.setSecurityScheme(applicationKeyMapping.getSecurityScheme());
        resolvedApplicationKeyMapping.setApplicationIdentifier(applicationKeyMapping.getApplicationIdentifier());
        resolvedApplicationKeyMapping.setKeyType(applicationKeyMapping.getKeyType());
        resolvedApplicationKeyMapping.setEnvId(applicationKeyMapping.getEnvID());
        Iterator<Map.Entry<String, ApplicationKeyMapping>> iterator = applicationKeyMappingMap.entrySet().iterator();
        while (iterator.hasNext()) {
            Map.Entry<String, ApplicationKeyMapping> cachedApplicationKeyMapping = iterator.next();
            ApplicationKeyMapping value = cachedApplicationKeyMapping.getValue();
            if (value.getApplicationIdentifier().equals(resolvedApplicationKeyMapping.getApplicationIdentifier()) && value.getSecurityScheme().equals(resolvedApplicationKeyMapping.getSecurityScheme()) && value.getKeyType().equals(resolvedApplicationKeyMapping.getKeyType()) && value.getEnvId().equals(resolvedApplicationKeyMapping.getEnvId()) && value.getApplicationUUID().equals(resolvedApplicationKeyMapping.getApplicationUUID())) {
                iterator.remove();
            }
        }
    }

    @Override
    public void removeSubscription(org.wso2.apk.enforcer.discovery.subscription.Subscription subscription) {

        subscriptionMap.remove(subscription.getUuid());

    }

    @Override
    public void removeApplication(org.wso2.apk.enforcer.discovery.subscription.Application application) {

        applicationMap.remove(application.getUuid());
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

    @Override
    public int getSubscriptionCount() {

        return subscriptionMap.size();
    }

    @Override
    public int getJWTIssuerCount() {

        return jwtValidatorMap.size();
    }

}
