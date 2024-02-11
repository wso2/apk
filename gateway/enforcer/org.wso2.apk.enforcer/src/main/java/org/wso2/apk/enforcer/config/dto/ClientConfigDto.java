/*
 * Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.enforcer.config.dto;

/**
 * This contains HttpClient properties.
 */
public class ClientConfigDto {

    private boolean enableSslVerification;
    private String hostnameVerifier;
    private int connectionTimeout;
    private int socketTimeout;
    private int maxConnections;
    private int maxConnectionsPerRoute;

    public boolean isEnableSslVerification() {

        return enableSslVerification;
    }

    public void setEnableSslVerification(boolean enableSslVerification) {

        this.enableSslVerification = enableSslVerification;
    }

    public String getHostnameVerifier() {

        return hostnameVerifier;
    }

    public void setHostnameVerifier(String hostnameVerifier) {

        this.hostnameVerifier = hostnameVerifier;
    }

    public int getConnectionTimeout() {

        return connectionTimeout;
    }

    public void setConnectionTimeout(int connectionTimeout) {

        this.connectionTimeout = connectionTimeout;
    }

    public int getSocketTimeout() {

        return socketTimeout;
    }

    public void setSocketTimeout(int socketTimeout) {

        this.socketTimeout = socketTimeout;
    }

    public int getMaxConnections() {

        return maxConnections;
    }

    public void setMaxConnections(int maxConnections) {

        this.maxConnections = maxConnections;
    }

    public int getMaxConnectionsPerRoute() {

        return maxConnectionsPerRoute;
    }

    public void setMaxConnectionsPerRoute(int maxConnectionsPerRoute) {

        this.maxConnectionsPerRoute = maxConnectionsPerRoute;
    }
}
