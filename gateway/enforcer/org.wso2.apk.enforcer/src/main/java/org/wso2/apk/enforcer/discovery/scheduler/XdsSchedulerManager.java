/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.discovery.scheduler;

import org.wso2.apk.enforcer.config.EnvVarConfig;
import org.wso2.apk.enforcer.discovery.ApiDiscoveryClient;
import org.wso2.apk.enforcer.discovery.ConfigDiscoveryClient;
import org.wso2.apk.enforcer.discovery.JWTIssuerDiscoveryClient;
import org.wso2.apk.enforcer.subscription.EventingGrpcClient;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

/**
 * Manages all the scheduling tasks that runs for retrying discovery requests.
 */
public class XdsSchedulerManager {

    private static int retryPeriod;
    private static volatile XdsSchedulerManager instance;
    private static ScheduledExecutorService discoveryClientScheduler;
    private static ScheduledExecutorService eventingScheduler;
    private ScheduledFuture<?> apiDiscoveryScheduledFuture;
    private ScheduledFuture<?> jwtIssuerDiscoveryScheduledFuture;

    private ScheduledFuture<?> eventingScheduledFuture;
    private ScheduledFuture<?> configDiscoveryScheduledFuture;

    public static XdsSchedulerManager getInstance() {

        if (instance == null) {
            synchronized (XdsSchedulerManager.class) {
                if (instance == null) {
                    instance = new XdsSchedulerManager();
                    discoveryClientScheduler = Executors.newSingleThreadScheduledExecutor();
                    eventingScheduler = Executors.newSingleThreadScheduledExecutor();
                    retryPeriod = Integer.parseInt(EnvVarConfig.getInstance().getXdsRetryPeriod());
                }
            }
        }
        return instance;
    }

    public synchronized void startAPIDiscoveryScheduling() {

        if (apiDiscoveryScheduledFuture == null || apiDiscoveryScheduledFuture.isDone()) {
            apiDiscoveryScheduledFuture = discoveryClientScheduler
                    .scheduleWithFixedDelay(ApiDiscoveryClient.getInstance(), 1, retryPeriod, TimeUnit.SECONDS);
        }
    }

    public synchronized void stopAPIDiscoveryScheduling() {

        if (apiDiscoveryScheduledFuture != null && !apiDiscoveryScheduledFuture.isDone()) {
            apiDiscoveryScheduledFuture.cancel(false);
        }
    }

    public synchronized void startJWTIssuerDiscoveryScheduling() {

        if (jwtIssuerDiscoveryScheduledFuture == null || jwtIssuerDiscoveryScheduledFuture.isDone()) {
            jwtIssuerDiscoveryScheduledFuture = discoveryClientScheduler
                    .scheduleWithFixedDelay(JWTIssuerDiscoveryClient.getInstance(), 1, retryPeriod, TimeUnit.SECONDS);
        }
    }

    public synchronized void startEventScheduling() {

        if (eventingScheduledFuture == null || eventingScheduledFuture.isDone()) {
            eventingScheduledFuture = eventingScheduler
                    .scheduleWithFixedDelay(EventingGrpcClient.getInstance(), 1, retryPeriod, TimeUnit.SECONDS);
        }
    }

    public synchronized void stopJWTIssuerDiscoveryScheduling() {

        if (jwtIssuerDiscoveryScheduledFuture != null && !jwtIssuerDiscoveryScheduledFuture.isDone()) {
            jwtIssuerDiscoveryScheduledFuture.cancel(false);
        }
    }

    public synchronized void stopEventStreamScheduling() {

        if (eventingScheduledFuture != null && !eventingScheduledFuture.isDone()) {
            eventingScheduledFuture.cancel(false);
        }
    }

    public synchronized void startConfigDiscoveryScheduling() {

        if (configDiscoveryScheduledFuture == null || configDiscoveryScheduledFuture.isDone()) {
            configDiscoveryScheduledFuture = discoveryClientScheduler
                    .scheduleWithFixedDelay(ConfigDiscoveryClient.getInstance(), 1, retryPeriod,
                            TimeUnit.SECONDS);
        }
    }

    public synchronized void stopConfigDiscoveryScheduling() {

        if (configDiscoveryScheduledFuture != null && !configDiscoveryScheduledFuture.isDone()) {
            configDiscoveryScheduledFuture.cancel(false);
        }
    }

}
