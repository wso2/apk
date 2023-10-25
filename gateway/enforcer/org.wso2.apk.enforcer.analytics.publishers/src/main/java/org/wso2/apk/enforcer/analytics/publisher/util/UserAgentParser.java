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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import ua_parser.Client;
import ua_parser.Parser;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * User agent parser util class.
 */
public class UserAgentParser {
    private static final Logger log = LoggerFactory.getLogger(UserAgentParser.class);
    private static final UserAgentParser INSTANCE = new UserAgentParser();
    private boolean isInitialized = false;
    private Parser uaParser;
    private Map<String, Client> clientCache;

    private UserAgentParser() {
        uaParser = new Parser();
        isInitialized = true;
        clientCache = new LinkedHashMap<String, Client>(Constants.USER_AGENT_DEFAULT_CACHE_SIZE
                                                                + Constants.DEFAULT_WORKER_THREADS) {
            static final int MAX = Constants.USER_AGENT_DEFAULT_CACHE_SIZE;
            @Override
            protected boolean removeEldestEntry(Map.Entry eldest) {
                return size() > MAX;
            }
        };
    }

    public static UserAgentParser getInstance() {
        return INSTANCE;
    }

    public Client parseUserAgent(String userAgentHeader) {
        if (isInitialized) {
            Client client = clientCache.get(userAgentHeader);
            if (client == null) {
                client = uaParser.parse(userAgentHeader);
                clientCache.put(userAgentHeader, client);
                log.debug("user agent added to cache. Current cache size is " + clientCache.size());
                return client;
            } else {
                log.debug("User agent fetched from cache");
                return client;
            }
        } else {
            return null;
        }
    }
}
