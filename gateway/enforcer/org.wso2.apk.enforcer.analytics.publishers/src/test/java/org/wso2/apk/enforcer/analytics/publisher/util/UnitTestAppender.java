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

import org.apache.logging.log4j.core.Filter;
import org.apache.logging.log4j.core.LogEvent;
import org.apache.logging.log4j.core.appender.AbstractAppender;
import org.apache.logging.log4j.core.config.plugins.Plugin;
import org.apache.logging.log4j.core.config.plugins.PluginAttribute;
import org.apache.logging.log4j.core.config.plugins.PluginElement;
import org.apache.logging.log4j.core.config.plugins.PluginFactory;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

@Plugin(name = "UnitTestAppender", category = "Core", elementType = "appender", printObject = true)
public class UnitTestAppender extends AbstractAppender {

    private List<String> messages = Collections.synchronizedList(new ArrayList<>());

    protected UnitTestAppender(String name, Filter filter) {

        super(name, filter, null, false, null);
    }

    @PluginFactory
    public static UnitTestAppender createAppender(
            @PluginAttribute("name") String name,
            @PluginElement("Filter") Filter filter) {

        if (name == null) {
            LOGGER.error("No name provided for UnitTestAppender");
            return null;
        }
        return new UnitTestAppender(name, filter);
    }

    public List<String> getMessages() {

        return messages;
    }

    public boolean checkContains(String message) {

        for (String log : messages) {
            if (log.contains(message)) {
                return true;
            }
        }
        return false;
    }

    @Override
    public void append(LogEvent event) {

        messages.add(event.getMessage().getFormattedMessage());
    }
}
