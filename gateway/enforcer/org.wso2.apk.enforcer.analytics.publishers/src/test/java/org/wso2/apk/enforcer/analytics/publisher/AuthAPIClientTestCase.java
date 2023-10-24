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

package org.wso2.apk.enforcer.analytics.publisher;

import org.testng.Assert;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.auth.AuthClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionRecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionUnrecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.util.AuthAPIMockService;

import java.util.HashMap;
import java.util.UUID;

public class AuthAPIClientTestCase extends AuthAPIMockService {

    @Test(expectedExceptions = {ConnectionUnrecoverableException.class},
            expectedExceptionsMessageRegExp = "Invalid/expired user token.*")
    public void testAuthClientWithAInvalidToken() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(401, authToken);

        AuthClient.getSASToken(authApiEndpoint, authToken, new HashMap<>());
    }

    @Test
    public void testAuthClientWithAValidToken() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(200, authToken);

        String sasToken = AuthClient.getSASToken(authApiEndpoint, authToken, new HashMap<>());
        Assert.assertEquals(sasToken, SAS_TOKEN);
    }

    @Test(expectedExceptions = {ConnectionRecoverableException.class})
    public void testAuthClientWithFor500Response() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(500, authToken);

        AuthClient.getSASToken(authApiEndpoint, authToken, new HashMap<>());
    }

    @Test(expectedExceptions = {ConnectionRecoverableException.class})
    public void testAuthClientWithFor400Response() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(400, authToken);

        AuthClient.getSASToken(authApiEndpoint, authToken, new HashMap<>());
    }

    @Test(expectedExceptions = {ConnectionRecoverableException.class},
            expectedExceptionsMessageRegExp = "Publisher has been temporarily revoked.")
    public void testAuthClientWithFor403Response() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(403, authToken);

        AuthClient.getSASToken(authApiEndpoint, authToken, new HashMap<>());
    }

    @Test(expectedExceptions = {ConnectionUnrecoverableException.class},
            expectedExceptionsMessageRegExp = "Invalid apim.analytics configurations provided.*")
    public void testAuthClientWithForInvalidAuthUrl() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(200, authToken);

        String authEndpoint = "invalid/host/auth-api";
        AuthClient.getSASToken(authEndpoint, authToken, new HashMap<>());
    }

    @Test(expectedExceptions = {ConnectionRecoverableException.class},
            expectedExceptionsMessageRegExp = "Provided authentication endpoint.*")
    public void testAuthClientWithForNonExistAuthUrl() throws Exception {

        String authToken = UUID.randomUUID().toString();
        mock(200, authToken);

        String authEndpoint = "https://no.such.host/auth-api";
        AuthClient.getSASToken(authEndpoint, authToken, new HashMap<>());
    }
}
