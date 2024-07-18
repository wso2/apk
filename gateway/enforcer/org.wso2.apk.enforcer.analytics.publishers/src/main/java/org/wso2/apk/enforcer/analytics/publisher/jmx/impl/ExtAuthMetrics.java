/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package org.wso2.apk.enforcer.analytics.publisher.jmx.impl;

import org.wso2.apk.enforcer.analytics.publisher.jmx.MBeanRegistrator;
import org.wso2.apk.enforcer.analytics.publisher.jmx.api.ExtAuthMetricsMXBean;
import org.wso2.apk.enforcer.analytics.publisher.reporter.prometheus.APIInvocationEvent;

import java.io.UnsupportedEncodingException;
import java.util.Timer;
import java.util.TimerTask;

/**
 * Singleton MBean for ExtAuth Service metrics.
 */
public class ExtAuthMetrics extends TimerTask implements ExtAuthMetricsMXBean {

    private static final long REQUEST_COUNT_INTERVAL_MILLIS = 5 * 60 * 1000;
    private static ExtAuthMetrics extAuthMetricsMBean = null;
    private long requestCountInLastFiveMinuteWindow = 0;
    private int tokenIssuerCount = 0;
    private long requestCountWindowStartTimeMillis = System.currentTimeMillis();
    private long totalRequestCount = 0;
    private int subscriptionCount = 0;
    private double averageResponseTimeMillis = 0;
    private double maxResponseTimeMillis = Double.MIN_VALUE;
    private double minResponseTimeMillis = Double.MAX_VALUE;

    private int apiMessages = 0;
    private int totalRequests = 0;
    private int postRequests = 0;
    private int getRequests = 0;
    private String resourcePaths = "";
    private String apiName = "sample1";
    private String applicationId = "";
    private int successRequests = 0;
    private int failureRequests = 0;

//    private ExtAuthMetrics() {
//        MBeanRegistrator.registerMBean(this);
//    }

    private ExtAuthMetrics(APIInvocationEvent event) {
        //getRequests++;
        //this.apiName = apiName;
        try {
            MBeanRegistrator.registerMBean(this, event);
        } catch (UnsupportedEncodingException e) {
            throw new RuntimeException(e);
        }
    }

/**
 * Getter for the Singleton ExtAuthMetrics instance.
 *
 * @return ExtAuthMetrics
 */
public static ExtAuthMetrics getInstance(APIInvocationEvent event) {
    if (extAuthMetricsMBean == null) {
        synchronized (ExtAuthMetrics.class) {
            if (extAuthMetricsMBean == null) {
                Timer timer = new Timer();
                extAuthMetricsMBean = new ExtAuthMetrics(event);
                extAuthMetricsMBean.requestCountWindowStartTimeMillis = System.currentTimeMillis();
                timer.schedule(extAuthMetricsMBean, REQUEST_COUNT_INTERVAL_MILLIS, REQUEST_COUNT_INTERVAL_MILLIS);
            }
        }
    }
    return extAuthMetricsMBean;
}

    @Override
    public long getTotalRequestCount() {
        return totalRequestCount;
    };

    @Override
    public double getAverageResponseTimeMillis() {
        return averageResponseTimeMillis;
    };

    @Override
    public double getMaxResponseTimeMillis() {
        return maxResponseTimeMillis;
    };

    @Override
    public double getMinResponseTimeMillis() {
        return minResponseTimeMillis;
    };

    public synchronized void recordMetric(long responseTimeMillis) {
        this.requestCountInLastFiveMinuteWindow += 1;
        this.totalRequestCount += 1;
        this.averageResponseTimeMillis = this.averageResponseTimeMillis +
                (responseTimeMillis - this.averageResponseTimeMillis) / totalRequestCount;
        this.minResponseTimeMillis = Math.min(this.minResponseTimeMillis, responseTimeMillis);
        this.maxResponseTimeMillis = Math.max(this.maxResponseTimeMillis, responseTimeMillis);
    }

    public synchronized void recordJWTIssuerMetrics(int jwtIssuers) {
        this.tokenIssuerCount = jwtIssuers;
    }

    public synchronized void recordSubscriptionMetrics(int subscriptionCount) {
        this.subscriptionCount = subscriptionCount;
    }

    @Override
    public synchronized void resetExtAuthMetrics() {
        this.totalRequestCount = 0;
        this.averageResponseTimeMillis = 0;
        this.maxResponseTimeMillis = Double.MIN_VALUE;
        this.minResponseTimeMillis = Double.MAX_VALUE;
    }

    @Override
    public synchronized void run() {
        requestCountWindowStartTimeMillis = System.currentTimeMillis();
        requestCountInLastFiveMinuteWindow = 0;
    }

    @Override
    public long getRequestCountInLastFiveMinuteWindow() {
        return requestCountInLastFiveMinuteWindow;
    }

    @Override
    public long getRequestCountWindowStartTimeMillis() {
        return requestCountWindowStartTimeMillis;
    }

    @Override
    public int getTokenIssuerCount() {
        return tokenIssuerCount;
    }

    @Override
    public int getSubscriptionCount() {
        return subscriptionCount;
    }

    public void recordApiMessages() {
        apiMessages++;;
    }

    @Override
    public int getApiMessages() {
        return apiMessages;
    }

    @Override
    public int getTotalRequests() {
        return totalRequests;
    }

    @Override
    public int getPostRequests() {
        return postRequests;
    }

    @Override
    public int getGetRequests() {
        return getRequests;
    }

    @Override
    public String getResourcePaths() {
        return resourcePaths;
    }

    public void incrementTotalRequests() {
        totalRequests++;
    }

    public void incrementPostRequests() {
        postRequests++;
    }

    public void incrementGetRequests() {
        getRequests++;
    }

    public void incrementResourcePath(String resourcePath) {
        resourcePaths = resourcePath;
    }

    public void incrementApiName(String apiName) {
        this.apiName = apiName;
    }

    @Override
    public String getApiName() {
        return apiName;
    }

    public void incrementApplicationId(String applicationId) {
        this.applicationId = applicationId;
    }

    @Override
    public String getApplicationId() {
        return applicationId;
    }

    public void incrementSuccessRequests() {
        successRequests++;
    }

    public void incrementFailureRequests() {
        failureRequests++;
    }

    @Override
    public int getSuccessRequests() {
        return successRequests;
    }

    @Override
    public int getFailureRequests() {
        return failureRequests;
    }
}