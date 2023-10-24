/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package org.wso2.apk.enforcer.analytics.publisher.util;


import java.util.HashMap;
import java.util.Map;
import java.util.Set;

/**
 * Utility to filter only the required attributes.
 */
public class EventMapAttributeFilter {
    private static final EventMapAttributeFilter INSTANCE = new EventMapAttributeFilter();

    public static EventMapAttributeFilter getInstance() {
        return INSTANCE;
    }

    public Map<String, Object> filter(Map<String, Object> source, Map<String, Class> requiredAttributes) {

        Set<String> targetKeys = requiredAttributes.keySet();
        Map<String, Object> filteredEventMap = new HashMap<>();

        for (String key : targetKeys) {
            filteredEventMap.put(key, source.get(key));
        }

        return filteredEventMap;
    }


}
