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

package org.wso2.apk.enforcer.analytics.publisher.client;

import com.azure.core.credential.AccessToken;
import com.azure.core.credential.TokenCredential;
import com.azure.core.credential.TokenRequestContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.auth.AuthClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionRecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionUnrecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.util.BackoffRetryCounter;
import reactor.core.publisher.Mono;

import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.Arrays;
import java.util.Map;

/**
 * WSO2 SAS token refresh implementation for TokenCredential.
 */
class WSO2TokenCredential implements TokenCredential {
    private static final Logger log = LoggerFactory.getLogger(WSO2TokenCredential.class);
    private final String authEndpoint;
    private final String authToken;
    private final Map<String, String> properties;
    private BackoffRetryCounter backoffRetryCounter;

    public WSO2TokenCredential(String authEndpoint, String authToken, Map<String, String> properties) {
        this.authEndpoint = authEndpoint;
        this.authToken = authToken;
        this.properties = properties;
        backoffRetryCounter = new BackoffRetryCounter();
    }

    @Override
    public Mono<AccessToken> getToken(TokenRequestContext tokenRequestContext) {
        log.debug("Trying to retrieving a new SAS token.");
        try {
            String sasToken = AuthClient.getSASToken(this.authEndpoint, this.authToken, this.properties);
            backoffRetryCounter.reset();
            log.debug("New SAS token retrieved.");
            // Using lower duration than actual.
            OffsetDateTime time = getExpirationTime(sasToken);
            return Mono.fromCallable(() -> new AccessToken(sasToken, time));
        } catch (ConnectionRecoverableException e) {
            log.error("Error occurred when retrieving SAS token. Connection will be retried in "
                              + backoffRetryCounter.getTimeInterval().replaceAll("[\r\n]", ""), e);
            try {
                Thread.sleep(backoffRetryCounter.getTimeIntervalMillis());
            } catch (InterruptedException interruptedException) {
                Thread.currentThread().interrupt();
            }
            backoffRetryCounter.increment();
            return getToken(tokenRequestContext);
        } catch (ConnectionUnrecoverableException e) {
            //Do not pass the exception. Publishing threads will encounter authentication issue and then try to
            // reinitialize publisher.
            log.error("Error occurred when retrieving SAS token.", e);
            backoffRetryCounter.reset();
            return Mono.error(e);
        }
    }

    private OffsetDateTime getExpirationTime(String sharedAccessSignature) {
        String[] parts = sharedAccessSignature.split("&");
        return Arrays.stream(parts).map(part -> part.split("="))
                .filter(pair -> pair.length == 2 && pair[0].equalsIgnoreCase("se"))
                .findFirst().map(pair -> pair[1])
                .map((expirationTimeStr) -> {
                    try {
                        long epochSeconds = Long.parseLong(expirationTimeStr);
                        return Instant.ofEpochSecond(epochSeconds).atOffset(ZoneOffset.UTC);
                    } catch (NumberFormatException e) {
                        log.error("Invalid expiration time format in the SAS token.", e);
                        return OffsetDateTime.MAX;
                    }
                })
                .orElse(OffsetDateTime.MAX);
    }
}
