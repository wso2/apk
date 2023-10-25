/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *
 */

package org.wso2.apk.enforcer.commons.analytics.publishers.impl;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import java.util.List;


/**
 * Fault event publisher implementation.
 */
public class FaultyRequestDataPublisher extends AbstractRequestDataPublisher {

    private static final Log log = LogFactory.getLog(FaultyRequestDataPublisher.class);

    @Override
    public CounterMetric getCounterMetric() {
        return null;
    }

    @Override
    public List<CounterMetric> getMultipleCounterMetrics() {
        try {
            return AnalyticsDataPublisher.getInstance().getFaultyMetricReporters();
        } catch (MetricCreationException e) {
            log.error("Unable to get faulty counter metrics", e);
            return null;
        }
    }
}
