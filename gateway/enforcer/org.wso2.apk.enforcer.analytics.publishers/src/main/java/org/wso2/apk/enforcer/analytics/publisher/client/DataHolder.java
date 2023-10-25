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
package org.wso2.apk.enforcer.analytics.publisher.client;

/**
 * Data holder class to generate mock analytics events.
 */
public class DataHolder {
    static final String[] NODE_ID = new String[]{"1", "2", "3", "4", "5"};
    static final String[] DEPLOYMENT_ID =
            new String[]{"prod", "prod", "prod", "prod", "prod"};
    static final String[] API_UUID = new String[]{"apiUUID1", "apiUUID2", "apiUUID3", "apiUUID4", "apiUUID5"};
    static final String[] REGION_ID = new String[]{"region1", "region2", "region3", "region4", "region5"};
    static final String[] GATEWAY_TYPE = new String[]{"type1", "type2", "type3", "type4", "type5"};
    static final String[] DESTINATION =
            new String[]{"destination1", "destination2", "destination3", "destination4", "destination5"};
    static final String[] REQUEST_MED_LATENCY = new String[]{"100", "200", "300", "400", "500"};
    static final String[] RESPONSE_MED_LATENCY = new String[]{"500", "400", "300", "200", "100"};
    static final String[] RESPONSE_LATENCY = new String[]{"100", "200", "300", "400", "500"};
    static final String[] RESPONSE_CODE = new String[]{"100", "200", "300", "400", "500"};
    static final String[] RESPONSE_SIZE = new String[]{"100", "200", "300", "400", "500"};
    static final String[]
            API_CREATOR = new String[]{"creator1", "creator2", "creator3", "creator4", "creator5"};
    static final String[] API_METHOD = new String[]{"POST", "GET", "PUT", "DELETE", "PATCH"};
    static final String[] API_RESOURCE_TEMPLATE = new String[]{"/{id}", "/{name}", "/{age}", "/{gender}", "/{country}"};
    static final String[] API_VERSION = new String[]{"1.0.0", "2.0.0", "3.0.0", "4.0.0", "5.0.0"};
    static final String[] API_NAME = new String[]{"api1", "api2", "api3", "api4", "api5"};
    static final String[]
            API_CONTEXT = new String[]{"/context1", "/context2", "/context3", "/context4", "/context5"};
    static final String[] APPLICATION_NAME = new String[]{"app1", "app2", "app3", "app4", "app5"};
    static final String[] KEY_TYPE = new String[]{"production", "sandbox"};
    static final String[]
            API_CREATOR_TENANT_DOMAIN = new String[]{"carbon.super", "carbon.super", "carbon.super", "carbon.super",
                                                     "carbon.super"};
    static final String[] APPLICATION_CONSUMER_KEY = new String[]{"key1", "key2", "key3", "key4", "key5"};
    static final String[] APPLICATION_OWNER = new String[]{"owner1", "owner2", "owner3", "owner4", "owner5"};
    static final String[] USER_AGENT = new String[]{"agent1", "agent2", "agent3", "agent4", "agent5"};
    static final String[] EVENT_TYPE = new String[]{"response", "response", "response", "response", "response"};
}
