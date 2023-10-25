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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Bounded concurrent queue wrapping {@link java.util.concurrent.ArrayBlockingQueue}.
 */
public class EventQueue {

    private static final Logger log = LoggerFactory.getLogger(EventQueue.class);
    private final BlockingQueue<MetricEventBuilder> eventQueue;
    private final ExecutorService publisherExecutorService;
    private final EventHubClient client;
    private final AtomicInteger failureCount;
    private final ScheduledExecutorService flushingExecutorService;

    public EventQueue(int queueSize, int workerThreadCount, EventHubClient client, int flushingDelay) {
        this.client = client;
        // Note : Using a fixed worker thread pool and a bounded queue to control the load on the server
        publisherExecutorService = Executors.newFixedThreadPool(workerThreadCount,
                                                                new DefaultAnalyticsThreadFactory("Queue-Worker"));
        flushingExecutorService = Executors.newScheduledThreadPool(workerThreadCount,
                                                                   new DefaultAnalyticsThreadFactory("Queue-Flusher"));
        eventQueue = new LinkedBlockingQueue<>(queueSize);
        failureCount = new AtomicInteger(0);
        for (int i = 0; i < workerThreadCount; i++) {
            if (i == 0) {
                publisherExecutorService.submit(new ParallelQueueWorker(eventQueue, client));
                flushingExecutorService.scheduleWithFixedDelay(new QueueFlusher(eventQueue, client), flushingDelay,
                                                               flushingDelay, TimeUnit.SECONDS);
            } else {
                EventHubClient clonedClient = client.clone();
                publisherExecutorService.submit(new ParallelQueueWorker(eventQueue, clonedClient));
                flushingExecutorService.scheduleWithFixedDelay(new QueueFlusher(eventQueue, clonedClient),
                                                               flushingDelay, flushingDelay, TimeUnit.SECONDS);
            }
        }
    }

    public void put(MetricEventBuilder builder) {
        try {
            if (!eventQueue.offer(builder)) {
                int count = failureCount.incrementAndGet();
                if (count == 1) {
                    log.error("Event queue is full. Starting to drop analytics events.");
                } else if (count % 1000 == 0) {
                    log.error("Event queue is full. " + count + " events dropped so far");
                }
            }
        } catch (RejectedExecutionException e) {
            log.warn("Task submission failed. Task queue might be full", e);
        }

    }

    @Override
    protected void finalize() throws Throwable {
        publisherExecutorService.shutdown();
        super.finalize();
    }

    protected EventHubClient getClient() {
        return this.client;
    }
}
