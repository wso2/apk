/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.constants;

/**
 * Interceptor Constants class holds constants shared between router's interceptor lua script and enforcer.
 * If a certain value is changed, the change should be added to the lua script implementation
 * as well.
 */
public class InterceptorConstants {
    /**
     * Key in a Key-Value pair of authentication context in request body to interceptor service.
     */
    public static class AuthContextFields {
        public static final String TOKEN_TYPE = "tokenType";
        public static final String TOKEN = "token";
        public static final String KEY_TYPE = "keyType"; // PRODUCTION or SANDBOX
    }

    /**
     * Data names for api metadata fileds
     */
    public static class APIMetadataFields {
        public static final String API_BASE_PATH = "basePath";
        public static final String API_NAME = "apiName";
        public static final String API_VERSION = "apiVersion";
        public static final String API_VHOST = "vhost";
        public static final String API_ORGANIZATION_ID = "organizationId";
        public static final String ENVIRONMENT = "environment";

    }
}
