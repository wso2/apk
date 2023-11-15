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

package org.wso2.apk.enforcer.analytics.publisher.util;

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;

import java.time.Clock;
import java.time.OffsetDateTime;
import java.util.List;

/**
 * Util class  containing helper methods for test
 */
public class TestUtils {

    public static void populateBuilder(MetricEventBuilder builder) throws MetricReportingException {

        String uaString = "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1_1 like Mac OS X) AppleWebKit/534.46 (KHTML, "
                + "like Gecko) Version/5.1 Mobile/9B206 Safari/7534.48.3";

        builder.addAttribute(Constants.REQUEST_TIMESTAMP, OffsetDateTime.now(Clock.systemUTC()).toString())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_TYPE, "HTTP")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CONTEXT, "/v1/")
                .addAttribute(Constants.API_RESOURCE_TEMPLATE, "/resource/{value}")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.ENVIRONMENT_ID, "Development")
                .addAttribute(Constants.DESTINATION, "localhost:8080")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.USER_AGENT_HEADER, uaString)
                .addAttribute(Constants.USER_NAME, "admin")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, 401)
                .addAttribute(Constants.RESPONSE_CACHE_HIT, true)
                .addAttribute(Constants.RESPONSE_LATENCY, 2000L)
                .addAttribute(Constants.BACKEND_LATENCY, 3000L)
                .addAttribute(Constants.REQUEST_MEDIATION_LATENCY, 1000L)
                .addAttribute(Constants.RESPONSE_MEDIATION_LATENCY, 1000L)
                .addAttribute(Constants.USER_IP, "127.0.0.1");
    }

    public static boolean isContains(List<String> messages, String message) {

        for (String log : messages) {
            if (log.contains(message)) {
                return true;
            }
        }
        return false;
    }
}
