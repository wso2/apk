/*
 * Copyright (c) 2024, WSO2 LLC (http://www.wso2.com).
 *
 * WSO2 LLC licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.apk.integration.utils;

public class Constants {

    public static final String DEFAULT_IDP_HOST = "am.wso2.com";
    public static final String DEFAULT_API_HOST = "am.wso2.com";
    public static final String DEFAULT_GW_PORT = "9443";
    public static final String DEFAULT_TOKEN_EP = "oauth2/token";
    public static final String DEFAULT_DCR_EP = "client-registration/v0.17/register";
    public static final String DEFAULT_API_CONFIGURATOR = "api/configurator/1.0.0/";
    public static final String DEFAULT_API_DEPLOYER = "api/am/publisher/v4/";
    public static final String DEFAULT_DEVPORTAL = "api/am/devportal/v3/";
    public static final String ACCESS_TOKEN = "accessToken";
    public static final String EMPTY_STRING = "";
    public static final String API_CREATE_SCOPE = "apk:api_create";
    public static final String SPACE_STRING = " ";
    public static final String SUBSCRIPTION_BASIC_AUTH_TOKEN =
            "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==";

    public class REQUEST_HEADERS {

        public static final String HOST = "Host";
        public static final String AUTHORIZATION = "Authorization";
        public static final String CONTENT_TYPE = "Content-Type";
    }

    public class CONTENT_TYPES {

        public static final String APPLICATION_JSON = "application/json";
        public static final String APPLICATION_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded";

        public static final String MULTIPART_FORM_DATA = "multipart/form-data";

        public static final String APPLICATION_OCTET_STREAM = "application/octet-stream";

        public static final String APPLICATION_ZIP = "application/zip";

        public static final String TEXT_PLAIN = "text/plain";
    }
}
