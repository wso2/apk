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

package org.wso2.apk.integration.api;

import org.apache.http.HttpResponse;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class SharedContext {

    private SimpleHTTPClient httpClient;
    private String publisherAccessToken;
    private String devportalAccessToken;
    private String publisherBasicAuthToken;
    private String devportalBasicAuthToken;
    private HttpResponse response;
    private String responseBody;
    private String apiUUID;
    private String revisionUUID;
    private HashMap<String, Object> valueStore = new HashMap<>();
    private HashMap<String, String> headers = new HashMap<>();

    public SimpleHTTPClient getHttpClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        if (httpClient == null) {
            httpClient = new SimpleHTTPClient();
        }
        return httpClient;
    }

    public String getPublisherAccessToken() {

        return publisherAccessToken;
    }

    public void setPublisherAccessToken(String accessToken) {

        this.publisherAccessToken = accessToken;
    }

    public String getDevportalAccessToken() {

        return devportalAccessToken;
    }

    public void setDevportalAccessToken(String accessToken) {

        this.devportalAccessToken = accessToken;
    }

    public String getPublisherBasicAuthToken() {

        return publisherBasicAuthToken;
    }

    public void setPublisherBasicAuthToken(String basicAuthToken) {

        this.publisherBasicAuthToken = basicAuthToken;
    }

    public String getDevportalBasicAuthToken() {

        return devportalBasicAuthToken;
    }

    public void setDevportalBasicAuthToken(String basicAuthToken) {

        this.devportalBasicAuthToken = basicAuthToken;
    }

    public HttpResponse getResponse() {

        return response;
    }

    public void setResponse(HttpResponse response) {

        this.response = response;
    }

    public Object getStoreValue(String key) {
        return valueStore.get(key);
    }

    public void addStoreValue(String key, Object value) {
        valueStore.put(key, value);
    }

    public Map<String, Object> getValueStore() {
        return Collections.unmodifiableMap(valueStore);
    }

    public Map<String, String> getHeaders() {
        return Collections.unmodifiableMap(headers);
    }

    public void addHeader(String key, String value) {
        headers.put(key, value);
    }

    public String getResponseBody() {

        return responseBody;
    }

    public void setResponseBody(String responseBody) {

        this.responseBody = responseBody;
    }

    public String getApiUUID() {

        return apiUUID;
    }

    public void setApiUUID(String apiUUID) {

        this.apiUUID = apiUUID;
    }

    public String getRevisionUUID() {

        return revisionUUID;
    }

    public void setRevisionUUID(String revisionUUID) {

        this.revisionUUID = revisionUUID;
    }
}
