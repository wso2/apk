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

import com.google.common.io.Resources;

import io.cucumber.java.Before;
import io.cucumber.java.en.Then;

import java.net.URL;
import java.nio.charset.StandardCharsets;
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

        @Then("I have a valid token with a client certificate {string}")
        public void getValidClientCertificateForMTLS(String clientCertificatePath) throws Exception {

                Map<String, String> headers = new HashMap<>();
                headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
                headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION,
                                "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

                HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers,
                                "grant_type=client_credentials&scope=" + Constants.API_CREATE_SCOPE,
                                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
                sharedContext.setAccessToken(Utils.extractToken(httpResponse));
                sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());

                URL url = Resources.getResource("artifacts/certificates/" + clientCertificatePath);
                String clientCertificate = Resources.toString(url, StandardCharsets.UTF_8);
                sharedContext.addStoreValue("clientCertificate", clientCertificate);

        }
}