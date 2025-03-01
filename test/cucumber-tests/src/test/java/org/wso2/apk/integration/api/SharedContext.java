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
import org.wso2.apk.integration.utils.clients.student_service.StudentResponse;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.ArrayList;
import java.util.List;

public class SharedContext {

    private SimpleHTTPClient httpClient;
    private String accessToken;
    private HttpResponse response;
    private String responseBody;
    private String publisherAccessToken;
    private String devportalAccessToken;
    private String adminportalAccessToken;
    private String basicAuthToken;
    private String apiUUID;
    private String revisionUUID;
    private StudentResponse studentResponse;
    private String applicationUUID;
    private String keyManagerUUID;
    private String oauthKeyUUID;
    private String consumerSecret;
    private String consumerKey;
    private String sandboxConsumerSecret;
    private String sandboxConsumerKey;
    private String prodKeyMappingID;
    private String sandboxKeyMappingID;
    private String apiAccessToken;
    private Boolean definitionValidStatus;
    private String subscriptionID;
    private String internalKey;
    private static String policyID;
    private HashMap<String, Object> valueStore = new HashMap<>();
    private HashMap<String, String> headers = new HashMap<>();
    private int grpcStatusCode;
    private int grpcErrorCode;
    private List<String> responses = new ArrayList<>();


    public SimpleHTTPClient getHttpClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        if (httpClient == null) {
            httpClient = new SimpleHTTPClient();
        }
        return httpClient;
    }

    public void addResponse(String response) {
        responses.add(response);
    }

    public List<String> getResponses() {
        return responses;
    }

    public int getGrpcStatusCode() {
        return grpcStatusCode;
    }

    public void setGrpcStatusCode(int grpcStatusCode) {
        this.grpcStatusCode = grpcStatusCode;
    }

    public StudentResponse getStudentResponse() {

        return studentResponse;
    }

    public void setStudentResponse(StudentResponse studentResponse) {

        this.studentResponse = studentResponse;
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

    public void removeHeader(String key) {
        headers.remove(key);
    }

    public String getResponseBody() {

        return responseBody;
    }

    public void setResponseBody(String responseBody) {

        this.responseBody = responseBody;
    }

    public void clearResposes() {
        responses.clear();
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

    public String getAdminAccessToken() {

        return adminportalAccessToken;
    }

    public void setAdminAccessToken(String accessToken) {

        this.adminportalAccessToken = accessToken;
    }

    public String getBasicAuthToken() {

        return basicAuthToken;
    }

    public void setBasicAuthToken(String basicAuthToken) {

        this.basicAuthToken = basicAuthToken;
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

    public String getApplicationUUID() {

        return applicationUUID;
    }

    public void setApplicationUUID(String applicationUUID) {

        this.applicationUUID = applicationUUID;
    }

    public String getKeyManagerUUID() {

        return keyManagerUUID;
    }

    public void setKeyManagerUUID(String keyManagerUUID) {

        this.keyManagerUUID = keyManagerUUID;
    }

    public String getOauthKeyUUID() {

        return oauthKeyUUID;
    }

    public void setOauthKeyUUID(String oauthKeyUUID) {

        this.oauthKeyUUID = oauthKeyUUID;
    }

    public void setAPIInternalKey(String internalKey) {
        this.internalKey = internalKey;
    }

    public String getAPIInternalKey() {
        return internalKey;
    }

    public String getConsumerSecret(String keyType) {
        if ("production".equals(keyType))
            return consumerSecret;
        else if ("sandbox".equals(keyType))
            return sandboxConsumerSecret;
        return "";
    }

    public void setConsumerSecret(String consumerSecret, String keyType) {
        if ("production".equals(keyType))
            this.consumerSecret = consumerSecret;
        else if ("sandbox".equals(keyType))
            this.sandboxConsumerSecret = consumerSecret;
    }

    public String getConsumerKey(String keyType) {
        if ("production".equals(keyType))
            return consumerKey;
        else if ("sandbox".equals(keyType))
            return sandboxConsumerKey;
        return "";
    }

    public void setConsumerKey(String consumerKey, String keyType) {
        if ("production".equals(keyType))
            this.consumerKey = consumerKey;
        else if ("sandbox".equals(keyType))
            this.sandboxConsumerKey = consumerKey;
    }

    public void setKeyMappingID(String keyMappingID, String keyType) {
        if ("production".equals(keyType))
            this.prodKeyMappingID = keyMappingID;
        else if ("sandbox".equals(keyType))
            this.sandboxKeyMappingID = keyMappingID;
    }

    public String getKeyMappingID(String keyType) {
        if ("production".equals(keyType))
            return prodKeyMappingID;
        else if ("sandbox".equals(keyType))
            return sandboxKeyMappingID;
        return "";
    }

    public String getApiAccessToken() {

        return apiAccessToken;
    }

    public void setApiAccessToken(String apiAccessToken) {

        this.apiAccessToken = apiAccessToken;
    }

    public void setAPIDefinitionValidStatus(Boolean definitionValidStatus) {
        this.definitionValidStatus = definitionValidStatus;
    }

    public Boolean getDefinitionValidStatus() {
        return definitionValidStatus;
    }

    public String getSubscriptionID() {

        return subscriptionID;
    }

    public void setSubscriptionID(String subID) {

        this.subscriptionID = subID;
    }

    public String getPolicyID() {

        return policyID;
    }

    public void setPolicyID(String policyId) {

        this.policyID = policyId;
    }
}
