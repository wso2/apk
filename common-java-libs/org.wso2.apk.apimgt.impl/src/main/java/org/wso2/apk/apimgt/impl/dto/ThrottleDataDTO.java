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


import java.util.Map;

/**
 * This class is used to hold throttling data before publish them to
 * global policy engine side. We decided to maintain this in impl as
 * this can be used by other components such as usage publisher.
 * In future we may consider adding all properties to this class.
 */
public class ThrottleDataDTO {
    String clientIP;
    int messageSizeInBytes;
    Map<String, String> transportHeaders;
    Map<String, String> queryParameters;

    public Map<String, String> getTransportHeaders() {
        return transportHeaders;
    }

    public void setTransportHeaders(Map<String, String> transportHeaders) {
        this.transportHeaders = transportHeaders;
    }

    public Map<String, String> getQueryParameters() {
        return queryParameters;
    }

    public void setQueryParameters(Map<String, String> queryParameters) {
        this.queryParameters = queryParameters;
    }

    public String getClientIP() {
        return clientIP;
    }

    public void setClientIP(String clientIP) {
        this.clientIP = clientIP;
    }

    public int getMessageSizeInBytes() {
        return messageSizeInBytes;
    }

    public void setMessageSizeInBytes(int messageSizeInBytes) {
        this.messageSizeInBytes = messageSizeInBytes;
    }

    /**
     * Cleaning DTO by setting null reference for all it instance variables.
     */
    public void cleanDTO(){
        this.clientIP = null;
        this.messageSizeInBytes = 0;
        this.transportHeaders =null;
        this.queryParameters =null;
    }

}
