/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package org.wso2.apk.enforcer.analytics.publisher.reporter.moesif;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.TimerMetric;
import org.wso2.apk.enforcer.analytics.publisher.retriever.MoesifKeyRetriever;
import org.wso2.apk.enforcer.analytics.publisher.retriever.MoesifKeyRetrieverFactory;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.util.Map;

/**
 * Moesif Metric Reporter Implementation. This implementation is responsible for sending analytics data into Moesif
 * dashboard in a secure and reliable way.
 */
public class MoesifReporter extends AbstractMetricReporter {
    private static final Logger log = LoggerFactory.getLogger(MoesifReporter.class);
    private final EventQueue eventQueue;

    public MoesifReporter(Map<String, String> properties) throws MetricCreationException {
        super(properties);
        MoesifKeyRetriever keyRetriever  = MoesifKeyRetrieverFactory.getMoesifKeyRetriever(properties);
        int queueSize = Constants.DEFAULT_QUEUE_SIZE;
        int workerThreads = Constants.DEFAULT_WORKER_THREADS;
        if (properties.get(Constants.QUEUE_SIZE) != null) {
            queueSize = Integer.parseInt(properties.get(Constants.QUEUE_SIZE));
        }
        if (properties.get(Constants.WORKER_THREAD_COUNT) != null) {
            workerThreads = Integer.parseInt(properties.get(Constants.WORKER_THREAD_COUNT));
        }
        this.eventQueue = new EventQueue(queueSize, workerThreads, keyRetriever);

    }

    @Override
    protected void validateConfigProperties(Map<String, String> map) throws MetricCreationException {

    }

    @Override
    public CounterMetric createCounter(String name, MetricSchema metricSchema) throws MetricCreationException {

        return new MoesifCounterMetric(name, eventQueue, metricSchema);
    }

    @Override
    protected TimerMetric createTimer(String s) {
        return null;
    }
}

