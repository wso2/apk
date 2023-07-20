/*
 * Copyright (c) 2023, WSO2 LLC (http://www.wso2.com).
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

package org.wso2.apk.integration.api;

import org.apache.http.HttpResponse;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

public class SharedContext {

    private SimpleHTTPClient httpClient;
    private String accessToken;
    private HttpResponse response;
    public SimpleHTTPClient getHttpClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        if (httpClient == null) {
            httpClient = new SimpleHTTPClient();
        }
        return httpClient;
    }

    public String getAccessToken() {

        return accessToken;
    }

    public void setAccessToken(String accessToken) {

        this.accessToken = accessToken;
    }

    public HttpResponse getResponse() {

        return response;
    }

    public void setResponse(HttpResponse response) {

        this.response = response;
    }
}
