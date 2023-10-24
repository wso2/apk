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
package org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.util;

/**
 * Class for constants related to external Moesif microservice.
 */
public class MoesifMicroserviceConstants {
    public static final String MOESIF_PROTOCOL_WITH_FQDN_KEY = "moesifProtocolWithFQDN";
    public static final String MOESIF_EP_COMMON_PATH = "moesif/moesif_key";
    public static final String MOESIF_MS_VERSIONING_KEY = "moesifMSVersioning";
    public static final String MS_USERNAME_CONFIG_KEY = "moesifMSAuthUsername";
    public static final String MS_PWD_CONFIG_KEY = "moesifMSAuthPwd";
    public static final String CONTENT_TYPE = "application/json";
    public static final String QUERY_PARAM = "org_id";
    public static final int NUM_RETRY_ATTEMPTS = 3;
    public static final long TIME_TO_WAIT = 10000;
    public static final int NUM_RETRY_ATTEMPTS_PUBLISH = 3;
    public static final long TIME_TO_WAIT_PUBLISH = 10000;
    public static final int REQUEST_READ_TIMEOUT = 10000;
    public static final long PERIODIC_CALL_DELAY = 300000;

}
