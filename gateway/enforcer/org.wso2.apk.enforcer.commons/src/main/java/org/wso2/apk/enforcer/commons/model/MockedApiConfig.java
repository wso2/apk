/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.commons.model;
import java.util.Map;

/**
 * Defines JSON script structure used in the mock API implementations.
 */
public class MockedApiConfig {
    private Map<String, MockedResponseConfig> responses;  // status code => responseconfig

    public Map<String, MockedResponseConfig> getResponses() {
        return responses;
    }

    public void setResponses(Map<String, MockedResponseConfig> responses) {
        this.responses = responses;
    }
}
