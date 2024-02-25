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
import org.wso2.apk.enforcer.common.CacheableEntity;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.discovery.JWTIssuerDiscoveryClient;
import org.wso2.apk.enforcer.discovery.subscription.Application;
import org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping;
import org.wso2.apk.enforcer.discovery.subscription.ApplicationMapping;
import org.wso2.apk.enforcer.discovery.subscription.Subscription;
import org.wso2.apk.enforcer.jmx.JMXUtils;
import org.wso2.apk.enforcer.metrics.jmx.impl.ExtAuthMetrics;
import org.wso2.apk.enforcer.util.ApacheFeignHttpClient;
import org.wso2.apk.enforcer.util.FilterUtils;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Utility methods related to subscription data store functionalities.
 */
public class SubscriptionDataStoreUtil {

    private static SubscriptionValidationDataRetrievalRestClient subscriptionValidationDataRetrievalRestClient;
    private static SubscriptionDataStoreUtil Instance;

    private SubscriptionDataStoreUtil() {

        String commonControllerHost = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHost();
        String commonControllerHostname = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHostname();
        int commonControllerRestPort =
                Integer.parseInt(ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerRestPort());
        subscriptionValidationDataRetrievalRestClient =
                Feign.builder().encoder(new GsonEncoder()).decoder(new GsonDecoder()).logger(new Slf4jLogger()).client(new ApacheFeignHttpClient(FilterUtils.getMutualSSLHttpClient("https", Arrays.asList(commonControllerHost, commonControllerHostname)))).target(SubscriptionValidationDataRetrievalRestClient.class, "https://" + commonControllerHost + ":" + commonControllerRestPort);
    }

    public static final String DELEM_PERIOD = ".";

    public static String getAPICacheKey(String context, String version) {

        return context + DELEM_PERIOD + version;
    }

    public static String getSubscriptionCacheKey(String appId, String apiId) {

        return appId + DELEM_PERIOD + apiId;
    }

    public static String getPolicyCacheKey(String tierName) {

        return tierName;
    }

    public static SubscriptionDataStoreUtil getInstance() {

        if (Instance == null) {
            synchronized (SubscriptionDataStoreUtil.class) {
                if (Instance == null) {
                    Instance = new SubscriptionDataStoreUtil();
                }
                return Instance;
            }
        }
        return Instance;
    }

    private static void loadApplicationKeyMappings() {

        new Thread(() -> {
            ApplicationKeyMappingDtoList applicationKeyMappings =
                    subscriptionValidationDataRetrievalRestClient.getAllApplicationKeyMappings();
            List<ApplicationKeyMappingDTO> list = applicationKeyMappings.getList();
            Map<String, List<ApplicationKeyMappingDTO>> orgWizeMAp = new HashMap<>();
            for (ApplicationKeyMappingDTO applicationKeyMappingDTO : list) {
                String organization = applicationKeyMappingDTO.getOrganizationId();
                List<ApplicationKeyMappingDTO> applicationKeyMappingDTOS = orgWizeMAp.computeIfAbsent(organization,
                        k -> new ArrayList<>());
                applicationKeyMappingDTOS.add(applicationKeyMappingDTO);
            }
            orgWizeMAp.forEach((k, v) -> {
                SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(k);
                subscriptionDataStore.addApplicationKeyMappings(v);
            });
        }).start();

    }

    private static void loadApplicationMappings() {

        new Thread(() -> {
            ApplicationMappingDtoList applicationMappings =
                    subscriptionValidationDataRetrievalRestClient.getAllApplicationMappings();
            List<ApplicationMappingDto> list = applicationMappings.getList();
            Map<String, List<ApplicationMappingDto>> orgWizeMAp = new HashMap<>();
            for (ApplicationMappingDto applicationMappingDto : list) {
                String organization = applicationMappingDto.getOrganizationId();
                List<ApplicationMappingDto> applicationMappingDtos = orgWizeMAp.computeIfAbsent(organization,
                        k -> new ArrayList<>());
                applicationMappingDtos.add(applicationMappingDto);
            }
            orgWizeMAp.forEach((k, v) -> {
                SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(k);
                subscriptionDataStore.addApplicationMappings(v);
            });
        }).start();

    }

    public static void initializeLoadingTasks() {

        JWTIssuerDiscoveryClient.getInstance().watchJWTIssuers();
        EventingGrpcClient.getInstance().watchEvents();
    }

    private static void loadApplications() {

        new Thread(() -> {
            ApplicationListDto applications = subscriptionValidationDataRetrievalRestClient.getAllApplications();
            List<ApplicationDto> list = applications.getList();
            Map<String, List<ApplicationDto>> orgWizeMAp = new HashMap<>();
            for (ApplicationDto applicationDto : list) {
                String organization = applicationDto.getOrganizationId();
                List<ApplicationDto> applicationDtos = orgWizeMAp.computeIfAbsent(organization, k -> new ArrayList<>());
                applicationDtos.add(applicationDto);
            }
            orgWizeMAp.forEach((k, v) -> {
                SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(k);
                subscriptionDataStore.addApplications(v);
            });
        }).start();
    }

    private static void loadSubscriptions() {

        new Thread(() -> {
            SubscriptionListDto subscriptions = subscriptionValidationDataRetrievalRestClient.getAllSubscriptions();
            List<SubscriptionDto> list = subscriptions.getList();
            if (JMXUtils.isJMXMetricsEnabled()) {
                ExtAuthMetrics.getInstance().recordSubscriptionMetrics(list.size());
            }
            Map<String, List<SubscriptionDto>> orgWizeMAp = new HashMap<>();
            for (SubscriptionDto subscriptionDto : list) {
                String organization = subscriptionDto.getOrganization();
                List<SubscriptionDto> subscriptionDtos = orgWizeMAp.computeIfAbsent(organization,
                        k -> new ArrayList<>());
                subscriptionDtos.add(subscriptionDto);
            }
            orgWizeMAp.forEach((k, v) -> {
                SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(k);
                subscriptionDataStore.addSubscriptions(v);
            });
        }).start();
    }

    public static String getApplicationKeyMappingCacheKey(String applicationIdentifier, String keyType,
                                                          String securityScheme, String envType) {

        return securityScheme + CacheableEntity.DELEM_PERIOD + envType + CacheableEntity.DELEM_PERIOD + keyType + CacheableEntity.DELEM_PERIOD + applicationIdentifier;
    }

    public static void addApplication(Application application) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(application.getOrganization());
        subscriptionDataStore.addApplication(application);

    }

    public static void addSubscription(Subscription subscription) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(subscription.getOrganization());
        subscriptionDataStore.addSubscription(subscription);

    }

    public static void addApplicationMapping(ApplicationMapping applicationMapping) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(applicationMapping.getOrganization());
        subscriptionDataStore.addApplicationMapping(applicationMapping);

    }

    public static void addApplicationKeyMapping(ApplicationKeyMapping applicationKeyMapping) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(applicationKeyMapping.getOrganization());
        subscriptionDataStore.addApplicationKeyMapping(applicationKeyMapping);

    }

    public static void removeApplicationMapping(ApplicationMapping applicationMapping) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(applicationMapping.getOrganization());
        subscriptionDataStore.removeApplicationMapping(applicationMapping);

    }

    private static SubscriptionDataStore getSubscriptionDataStore(String organization) {

        SubscriptionDataStore subscriptionDataStore =
                SubscriptionDataHolder.getInstance().getSubscriptionDataStore(organization);
        if (subscriptionDataStore == null) {
            synchronized (organization.concat("subscriptionDataStore").intern()) {
                subscriptionDataStore = SubscriptionDataHolder.getInstance().getSubscriptionDataStore(organization);
                if (subscriptionDataStore != null) {
                    return subscriptionDataStore;
                }
                subscriptionDataStore = SubscriptionDataHolder.getInstance().initializeSubscriptionDataStore(organization);
                return subscriptionDataStore;
            }
        }
        return subscriptionDataStore;
    }

    public static void removeApplicationKeyMapping(ApplicationKeyMapping applicationKeyMapping) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(applicationKeyMapping.getOrganization());

        subscriptionDataStore.removeApplicationKeyMapping(applicationKeyMapping);

    }

    public static void removeSubscription(Subscription subscription) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(subscription.getOrganization());

        subscriptionDataStore.removeSubscription(subscription);

    }

    public static void removeApplication(Application application) {

        SubscriptionDataStore subscriptionDataStore = getSubscriptionDataStore(application.getOrganization());
        subscriptionDataStore.removeApplication(application);

    }

    public void loadStartupArtifacts() {

        loadApplications();
        loadSubscriptions();
        loadApplicationMappings();
        loadApplicationKeyMappings();

    }
}

