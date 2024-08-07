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

package org.wso2.apk.enforcer.analytics.publisher.jmx;

/**
 * JMX Utilities
 */
public class JMXUtils {

    private static final String APK_JMX_METRICS_ENABLE = "apk.jmx.metrics.enabled";

    /**
     * Returns true if jmx metrics enabled as a system property, otherwise false.
     * 
     * @return boolean
     */
    public static boolean isJMXMetricsEnabled() {
        return Boolean.getBoolean(APK_JMX_METRICS_ENABLE);
    }
}
