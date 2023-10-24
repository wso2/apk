/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.analytics.publisher.util;

import org.mockserver.integration.ClientAndServer;
import org.mockserver.model.JsonBody;
import org.testng.annotations.AfterClass;
import org.testng.annotations.BeforeClass;
import org.wso2.apk.enforcer.analytics.publisher.auth.TokenDetailsDTO;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

/**
 * Mocking auth-api service.
 */
public class AuthAPIMockService {

    protected static final String SAS_TOKEN = "SharedAccessSignature sr=sb://localhost/incoming-hub/publishers"
            + "/pub1&sig=signature&se=1641892957&skn=send-policy";
    private static final int TEST_PORT = 9191;
    private ClientAndServer mockServer;
    protected String authApiEndpoint;

    @BeforeClass
    public void startServer() {

        mockServer = ClientAndServer.startClientAndServer(TEST_PORT);
        authApiEndpoint = "http://localhost:" + TEST_PORT + "/auth-api";
    }

    @AfterClass
    public void stopServer() {

        mockServer.stop();
    }

    protected void mock(int responseCode, String token) {

        TokenDetailsDTO dto = new TokenDetailsDTO();
        dto.setToken(SAS_TOKEN);
        mockServer.when(
                        request()
                                .withMethod("GET")
                                .withPath("/auth-api/token")
                                .withHeader("Authorization", "Bearer " + token)
                )
                .respond(
                        response()
                                .withStatusCode(responseCode)
                                .withBody(JsonBody.json(dto))
                );
    }
}
