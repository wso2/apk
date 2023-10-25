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

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;

import java.util.Map;

/**
 * Main interface class for Metric Event Builders. Metric Event Builders are responsible of collecting metrics,
 * validating them and later returning them as a Map&lt;String, Object&gt;. Default builders will be implemented and
 * for any custom message building new builders have to be introduced
 */
public interface MetricEventBuilder {

    /**
     * Validates the provided attributes and build a flat {@link Map}. Any validation failures will cause
     * {@link MetricReportingException}.
     *
     * @return Map containing all attributes related to Metric Event
     * @throws MetricReportingException if validation failed
     */
    public Map<String, Object> build() throws MetricReportingException;

    /**
     * Checks the validity of the added attributes. If all required attributes are present true will be returned.
     * Else {@link MetricReportingException} will be thrown.
     *
     * @return Validity state of the added data
     * @throws MetricReportingException if validation failed
     */
    public boolean validate() throws MetricReportingException;

    /**
     * Method to add any attribute to the builder. Each concrete implementation can implement validations
     * based on the key.
     *
     * @param key    Key of the attribute
     * @param number Value of the attribute
     * @return Returns itself to support chaining
     * @throws MetricReportingException if validation failed
     */
    public MetricEventBuilder addAttribute(String key, Object number) throws MetricReportingException;
}
