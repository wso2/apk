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

import com.google.gson.Gson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;

import java.util.Map;
import java.util.concurrent.BlockingQueue;

/**
 * Will removes the events from queues and send then to the endpoints.
 */
public class ParallelQueueWorker implements Runnable {

    private static final Logger log = LoggerFactory.getLogger(ParallelQueueWorker.class);
    private BlockingQueue<MetricEventBuilder> eventQueue;
    private EventHubClient client;

    public ParallelQueueWorker(BlockingQueue<MetricEventBuilder> queue, EventHubClient client) {
        this.client = client;
        this.eventQueue = queue;
    }

    public void run() {
            if (log.isDebugEnabled()) {
                log.debug(eventQueue.size() + " messages in queue before " +
                                  Thread.currentThread().getName().replaceAll("[\r\n]", "")
                                  + " worker has polled queue");
            }
            while (true) {
                MetricEventBuilder eventBuilder;
                String event;
                try {
                    eventBuilder = eventQueue.take();
                    if (eventBuilder != null) {
                        Map<String, Object> eventMap = eventBuilder.build();
                        event = new Gson().toJson(eventMap);
                        client.sendEvent(event);
                    }
                } catch (MetricReportingException e) {
                    log.error("Builder instance is not duly filled. Event building failed", e);
                    continue;
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } catch (Exception e) {
                    log.error("Analytics event sending failed. Event will be dropped", e);
                }

                if (log.isDebugEnabled()) {
                    log.debug(eventQueue.size() + " messages in queue after " +
                                      Thread.currentThread().getName().replaceAll("[\r\n]", "")
                                      + " worker has finished work");
                }
            }
    }
}
