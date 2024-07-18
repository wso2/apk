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

package org.wso2.apk.enforcer.analytics.publisher.reporter;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.cloud.DefaultAnalyticsMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.elk.ELKMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.MoesifReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.prometheus.PrometheusMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.util.HashMap;
import java.util.Map;

/**
 * Factory class to create {@link MetricReporter}. Based on the passed argument relevant concrete implementation will
 * be created and returned. Factory will behave in Singleton manner and if same type of instance is requested again
 * Factory will return earlier requested instance.
 */
public class MetricReporterFactory {

    private static final Logger log = LoggerFactory.getLogger(MetricReporterFactory.class);
    private static final MetricReporterFactory instance = new MetricReporterFactory();
    private static Map<String, MetricReporter> reporterRegistry = new HashMap<>();

    private MetricReporterFactory() {
        //private constructor
    }

    public MetricReporter createMetricReporter(Map<String, String> properties)
            throws MetricCreationException {

        if (reporterRegistry.get(Constants.DEFAULT_REPORTER) == null) {
            synchronized (this) {
                if (reporterRegistry.get(Constants.DEFAULT_REPORTER) == null) {
                    MetricReporter reporterInstance = new DefaultAnalyticsMetricReporter(properties);
                    reporterRegistry.put(Constants.DEFAULT_REPORTER, reporterInstance);
                    return reporterInstance;
                }
            }
        }
        MetricReporter reporterInstance = reporterRegistry.get(Constants.DEFAULT_REPORTER);
        log.info("Metric Reporter of type " + reporterInstance.getClass().toString().replaceAll("[\r\n]", "") +
                " is already created. Hence returning same instance");
        return reporterInstance;
    }

    public MetricReporter createLogMetricReporter(Map<String, String> properties) throws MetricCreationException {

        if (reporterRegistry.get(Constants.ELK_REPORTER) == null) {
            synchronized (this) {
                if (reporterRegistry.get(Constants.ELK_REPORTER) == null) {
                    MetricReporter reporterInstance = new ELKMetricReporter(properties);
                    reporterRegistry.put(Constants.ELK_REPORTER, reporterInstance);
                    return reporterInstance;
                }
            }
        }

        MetricReporter reporterInstance = reporterRegistry.get(Constants.ELK_REPORTER);
        log.info("Metric Reporter of type " + reporterInstance.getClass().toString().replaceAll("[\r\n]", "") +
                " is already created. Hence returning same instance");
        return reporterInstance;
    }

    public MetricReporter createMetricReporter(String fullyQualifiedClassName, Map<String, String> properties)
            throws MetricCreationException {

        if (reporterRegistry.get(fullyQualifiedClassName) == null) {
            synchronized (this) {
                if (reporterRegistry.get(fullyQualifiedClassName) == null) {
                    if (fullyQualifiedClassName != null && !fullyQualifiedClassName.isEmpty()) {
                        try {
                            Class<MetricReporter> clazz =
                                    (Class<MetricReporter>) Class.forName(fullyQualifiedClassName);
                            Constructor<MetricReporter> constructor =
                                    clazz.getConstructor(Map.class);
                            MetricReporter reporterInstance = constructor.newInstance(properties);
                            reporterRegistry.put(fullyQualifiedClassName, reporterInstance);
                            return reporterInstance;
                        } catch (InstantiationException | IllegalAccessException | ClassNotFoundException
                                 | NoSuchMethodException | InvocationTargetException e) {
                            throw new MetricCreationException("Error occurred while creating a Metric Reporter of type"
                                    + " " + fullyQualifiedClassName, e);
                        }
                    } else {
                        throw new MetricCreationException("Provided class name is either empty or null. Hence cannot "
                                + "create the Reporter.");
                    }
                }
            }
        }
        MetricReporter reporterInstance = reporterRegistry.get(fullyQualifiedClassName);
        log.info("Metric Reporter of type " + reporterInstance.getClass().toString().replaceAll("[\r\n]", "") +
                " is already created. Hence returning same instance");
        return reporterInstance;
    }

    public MetricReporter createMetricReporterFromType(String type, Map<String, String> properties)
            throws MetricCreationException {

        String fullyQualifiedClassName = DefaultAnalyticsMetricReporter.class.getName();
        if (Constants.ELK_REPORTER.equals(type)) {
            fullyQualifiedClassName = ELKMetricReporter.class.getName();
        } else if (Constants.DEFAULT_REPORTER.equals(type)) {
            fullyQualifiedClassName = DefaultAnalyticsMetricReporter.class.getName();
        } else if (Constants.MOESIF_REPORTER.equals(type)) {
            fullyQualifiedClassName = MoesifReporter.class.getName();
        } else if (Constants.PROMETHEUS_REPORTER.equals(type)) {
            fullyQualifiedClassName = PrometheusMetricReporter.class.getName();
        }
        return createMetricReporter(fullyQualifiedClassName , properties);
    }

    /**
     * Reset the MetricReporterFactory registry. Only intended to be used in testing
     */
    public void reset() {

        reporterRegistry.clear();
    }

    public static MetricReporterFactory getInstance() {

        return instance;
    }
}
