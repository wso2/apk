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

package org.wso2.apk.enforcer.analytics.publisher.reporter.cloud;

import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;

import java.util.concurrent.BlockingQueue;

/**
 * Class to periodically flush event batch.
 */
public class QueueFlusher implements Runnable {

    private final BlockingQueue<MetricEventBuilder> queue;
    private final EventHubClient client;

    public QueueFlusher(BlockingQueue<MetricEventBuilder> queue, EventHubClient client) {
        this.queue = queue;
        this.client = client;
    }

    @Override public void run() {
        if (queue.isEmpty()) {
            //For scenarios where no API invocation is happening additional check in the EventHubClient will stop
            // empty batch flushing
            client.flushEvents();
        }
    }
}
