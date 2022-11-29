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
package org.wso2.apk.apimgt.api.model;

import java.util.ArrayList;
import java.util.List;

public class EndpointConfig {


    private String url = null;


    private String timeout = null;


    private Boolean isPrimary = null;


    private List<EndpointConfigAttributes> attributes = new ArrayList<EndpointConfigAttributes>();

    public String getUrl() {

        return url;
    }

    public void setUrl(String url) {

        this.url = url;
    }

    public String getTimeout() {

        return timeout;
    }

    public void setTimeout(String timeout) {

        this.timeout = timeout;
    }

    public Boolean getPrimary() {

        return isPrimary;
    }

    public void setPrimary(Boolean primary) {

        isPrimary = primary;
    }

    public List<EndpointConfigAttributes> getAttributes() {

        return attributes;
    }

    public void setAttributes(List<EndpointConfigAttributes> attributes) {

        this.attributes = attributes;
    }

    @Override
    public String toString() {

        return "EndpointConfig{" +
                "url='" + url + '\'' +
                ", timeout='" + timeout + '\'' +
                ", isPrimary=" + isPrimary +
                ", attributes=" + attributes +
                '}';
    }
}
