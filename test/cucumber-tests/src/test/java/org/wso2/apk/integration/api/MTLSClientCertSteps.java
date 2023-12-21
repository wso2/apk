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
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import io.cucumber.java.Before;
import io.cucumber.java.en.Then;

import java.util.HashMap;
import java.util.Map;

/**
 * This class contains the common step definitions.
 */
public class MTLSClientCertSteps {

        private final SharedContext sharedContext;
        private SimpleHTTPClient httpClient;

        public MTLSClientCertSteps(SharedContext sharedContext) {

                this.sharedContext = sharedContext;
        }

        @Before
        public void setup() throws Exception {

                httpClient = sharedContext.getHttpClient();
        }

        @Then("I have a valid subscription with a valid client certificate")
        public void getValidClientCertificateForMTLS() throws Exception {

                Map<String, String> headers = new HashMap<>();
                headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
                headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION,
                                "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

                HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers,
                                "grant_type=client_credentials&scope=" + Constants.API_CREATE_SCOPE,
                                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
                sharedContext.setAccessToken(Utils.extractToken(httpResponse));
                sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());
                sharedContext.addStoreValue("clientCertificate",
                                "-----BEGIN CERTIFICATE-----MIIDGTCCAgECFANIkLQBkd76qiTXzSXjBS2scPJsMA0GCSqGSIb3DQEBCwUAME0xCzAJBgNVBAYTAkxLMRMwEQYDVQQIDApTb21lLVN0YXRlMQ0wCwYDVQQKDAR3c28yMQwwCgYDVQQLDANhcGsxDDAKBgNVBAMMA2FwazAeFw0yMzEyMDYxMDEyNDhaFw0yNTA0MTkxMDEyNDhaMEUxCzAJBgNVBAYTAkxLMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCdG90W/Tlk4u9awHPteD5zpVcThUKwMLvAKw9ivVQBC0AG6GzPbakol5gKVm+kBUDFzzzF6eayEXKWbyaZDty66A2+7HLLcKBop5M/a57Q9XtU3lRYvotgutLWuHcI7mLCScZDrjA3rnb/KjjbhZ602ZS1pp5jtyUz6DwLm7w4wQ/RProqCdBj8QqoAvnDDLSPeDfsx14J5VeNJVGJV2wax65jWRjRkj6wE7z2qzWAlP5vDeED6bogYYVDpC8DtgayQ+vKAQLi1uj+I9Yqb/nPUrdUh9IlxudlqiFQQxyvsXMJEzbWWmlbD0kXYkHmHzetJNPK9ayOS/fJcAcfAb01AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAFmUc7+cI8d0Dl4wTdq+gfyWdqjQb7AYVO9DvJi3XGxdc5Kp1nCSsKzKUz9gvxXHeaYKrBNYf4SSU+Pkdf/BWePqi7UX/SIxNXby2da8zWg+W6UhxZfKlLYGMp3mCjueZpZTJ7SKOOGFA8IIgEzjJD9Ln1gl3ywMaCwlNrG9RpiD1McTCOKvyWNKnSRVr/RvCklLVrAMTJr50kce2czcdFl/xF4Hm66vp7cP/bYJKWAL8hBGzUa9aQBKncOoAO+zQ/SGy7uJxTDUF8SverDsmjOc6AU6IhBGVUyX/JQbYyJfZinBYlviYxVzIm6IaNJHx4sihw4U1/jMFWRXT470zcQ=-----END CERTIFICATE-----");
        }

        @Then("I have a valid subscription with an invalid client certificate")
        public void getInvalidClientCertificateForMTLS() throws Exception {

                Map<String, String> headers = new HashMap<>();
                headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
                headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION,
                                "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

                HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers,
                                "grant_type=client_credentials&scope=" + Constants.API_CREATE_SCOPE,
                                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
                sharedContext.setAccessToken(Utils.extractToken(httpResponse));
                sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());
                sharedContext.addStoreValue("clientCertificate",
                                "-----BEGIN CERTIFICATE-----MIIDJDCfeXw==-----END CERTIFICATE-----");
        }
}