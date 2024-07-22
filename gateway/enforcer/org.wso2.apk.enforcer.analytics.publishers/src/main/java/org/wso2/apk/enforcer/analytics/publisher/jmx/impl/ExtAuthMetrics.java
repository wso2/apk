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
    private long requestCountWindowStartTimeMillis = System.currentTimeMillis();

    private int apiMessages = 0;

    private ExtAuthMetrics(APIInvocationEvent event) {
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
    Timer timer = new Timer();
    extAuthMetricsMBean = new ExtAuthMetrics(event);
    extAuthMetricsMBean.requestCountWindowStartTimeMillis = System.currentTimeMillis();
    timer.schedule(extAuthMetricsMBean, REQUEST_COUNT_INTERVAL_MILLIS, REQUEST_COUNT_INTERVAL_MILLIS);
    return extAuthMetricsMBean;
}


    @Override
    public synchronized void run() {
        requestCountWindowStartTimeMillis = System.currentTimeMillis();
        requestCountInLastFiveMinuteWindow = 0;
    }


    public void recordApiMessages() {
        apiMessages++;;
    }

    @Override
    public int getApiMessages() {
        return apiMessages;
    }

}