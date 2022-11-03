/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.impl.dto;

import java.util.HashSet;
import java.util.Set;

public class GatewayCleanupSkipList {

    private Set<String> apis = new HashSet<>();
    private Set<String> endpoints = new HashSet();
    private Set<String> localEntries = new HashSet<>();
    private Set<String> sequences = new HashSet<>();

    public Set<String> getApis() {

        return apis;
    }

    public void setApis(Set<String> apis) {

        this.apis = apis;
    }

    public Set<String> getEndpoints() {

        return endpoints;
    }

    public void setEndpoints(Set<String> endpoints) {

        this.endpoints = endpoints;
    }

    public Set<String> getLocalEntries() {

        return localEntries;
    }

    public void setLocalEntries(Set<String> localEntries) {

        this.localEntries = localEntries;
    }

    public Set<String> getSequences() {

        return sequences;
    }

    public void setSequences(Set<String> sequences) {

        this.sequences = sequences;
    }
}
