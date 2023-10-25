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

import com.azure.core.amqp.AmqpRetryOptions;
import com.azure.core.amqp.exception.AmqpErrorCondition;
import com.azure.core.amqp.exception.AmqpException;
import com.azure.messaging.eventhubs.EventData;
import com.azure.messaging.eventhubs.EventDataBatch;
import com.azure.messaging.eventhubs.EventHubProducerClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionRecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionUnrecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.cloud.DefaultAnalyticsThreadFactory;
import org.wso2.apk.enforcer.analytics.publisher.util.BackoffRetryCounter;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import javax.xml.bind.DatatypeConverter;

/**
 * Event Hub client is responsible for sending events to
 * Azure Event Hub.
 */
public class EventHubClient implements Cloneable {
    private static final Logger log = LoggerFactory.getLogger(EventHubClient.class);
    private static final String TOKEN_HASH_USER_PROP = "token-hash";
    private final String authEndpoint;
    private final String authToken;
    private final String authTokenHash;
    private final Lock publishingLock;
    private final BackoffRetryCounter producerRetryCounter;
    private final BackoffRetryCounter eventBatchRetryCounter;
    private final Lock threadBarrier;
    private final AmqpRetryOptions retryOptions;
    private final Condition waitCondition;
    private final ScheduledExecutorService scheduledExecutorService;
    private EventHubProducerClient producer;
    private EventDataBatch batch;
    private ClientStatus clientStatus;
    private Map<String, String> properties = new HashMap<>();

    public EventHubClient(String authEndpoint, String authToken, AmqpRetryOptions retryOptions,
                          Map<String, String> properties) {
        threadBarrier = new ReentrantLock();
        waitCondition = threadBarrier.newCondition();
        publishingLock = new ReentrantLock();
        scheduledExecutorService = Executors.newScheduledThreadPool(2, new DefaultAnalyticsThreadFactory(
                "Reconnection-Service"));
        producerRetryCounter = new BackoffRetryCounter();
        eventBatchRetryCounter = new BackoffRetryCounter();
        this.authEndpoint = authEndpoint;
        this.authToken = authToken;
        this.authTokenHash = toHash(authToken);
        this.retryOptions = retryOptions;
        this.clientStatus = ClientStatus.NOT_CONNECTED;
        this.properties = properties;
        createProducerWithRetry(authEndpoint, authToken, retryOptions, true, properties);
    }

    private void retryWithBackoff(String authEndpoint, String authToken,
                                  AmqpRetryOptions retryOptions, boolean createBatch) {
        scheduledExecutorService.schedule(new Runnable() {
            @Override
            public void run() {
                createProducerWithRetry(authEndpoint, authToken, retryOptions, createBatch, properties);
            }
        }, producerRetryCounter.getTimeIntervalMillis(), TimeUnit.MILLISECONDS);
        producerRetryCounter.increment();
    }

    private void createProducerWithRetry(String authEndpoint, String authToken, AmqpRetryOptions retryOptions,
                                         boolean createBatch, Map<String, String> properties) {
        log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                          + "- Creating Eventhub client instance.");
        try {
            if (producer != null) {
                producer.close();
            }
            producer = EventHubProducerClientFactory.create(authEndpoint, authToken, retryOptions, properties);
            try {
                if (createBatch) {
                    batch = producer.createBatch();
                }
            } catch (IllegalStateException e) {
                throw new ConnectionRecoverableException("Event batch creation failed. " + e.getMessage()
                        .replaceAll("[\r\n]", ""));
            } catch (AmqpException e) {
                throw new ConnectionRecoverableException("Event batch creation failed. " + e.getMessage()
                        .replaceAll("[\r\n]", ""));
            } catch (Exception e) {
                throw new ConnectionUnrecoverableException("Event batch creation failed. " + e.getMessage()
                        .replaceAll("[\r\n]", ""));
            }
            clientStatus = ClientStatus.CONNECTED;
            log.info("[" + Thread.currentThread().getName().replaceAll("[\r\n]", "") + "] "
                             + "- Eventhub client successfully connected.");
            producerRetryCounter.reset();
            try {
                threadBarrier.lock();
                waitCondition.signalAll();
            } finally {
                threadBarrier.unlock();
            }
        } catch (ConnectionRecoverableException e) {
            clientStatus = ClientStatus.RETRYING;
            log.error("Recoverable error occurred when creating Eventhub Client. Retry attempts will be made in "
                              + producerRetryCounter.getTimeInterval().replaceAll("[\r\n]", "") + ". Reason :"
                              + e.getMessage().replaceAll("[\r\n]", ""));
            if (log.isDebugEnabled()) {
                log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] - "
                                  + "Recoverable error occurred when creating Eventhub Client using following "
                                  + "attributes. Auth endpoint: " + authEndpoint.replaceAll("[\r\n]", "")
                                  + ". Retry attempts will be made. Reason : "
                                  + e.getMessage().replaceAll("[\r\n]", ""), e);
            }
            retryWithBackoff(authEndpoint, authToken, retryOptions, createBatch);
        } catch (ConnectionUnrecoverableException e) {
            clientStatus = ClientStatus.NOT_CONNECTED;
            log.error("Unrecoverable error occurred when creating Eventhub Client. Analytics event publishing will be"
                              + " disabled until issue is rectified. Reason: "
                              + e.getMessage().replaceAll("[\r\n]", ""));
            if (log.isDebugEnabled()) {
                log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                  + "- Unrecoverable error occurred when creating Eventhub Client using following "
                                  + "attributes. Auth endpoint: " + authEndpoint.replaceAll("[\r\n]", "") + ". "
                                  + "Analytics event publishing will be disabled until issue is rectified. Reason: "
                                  + e.getMessage().replaceAll("[\r\n]", ""), e);
            }
        }
    }

    public void sendEvent(String event) {
        if (clientStatus == ClientStatus.CONNECTED) {
            EventData eventData = new EventData(event);
            eventData.getProperties().put(TOKEN_HASH_USER_PROP, this.authTokenHash);
            try {
                publishingLock.lock();
                boolean isAdded = batch.tryAdd(eventData);
                if (!isAdded) {
                    try {
                        int size = 0;
                        if (log.isDebugEnabled()) {
                            size = batch.getCount();
                        }
                        producer.send(batch);
                        batch = createBatchWithRetry();
                        isAdded = batch.tryAdd(eventData);
                        log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                          + "- Published " + size + " events to Analytics cluster.");

                    } catch (AmqpException e) {
                        if (isAuthenticationFailure(e)) {
                            //if authentication error try to reinitialize publisher. Retrying will deal with any
                            // network or revocation failures.
                            log.error("Authentication issue happened. Producer client will be re-initialized "
                                              + "retaining the Event Data Batch");
                            this.clientStatus = ClientStatus.RETRYING;
                            createProducerWithRetry(authEndpoint, authToken, retryOptions, false, properties);
                            sendEvent(event);
                        } else if (e.getErrorCondition() == AmqpErrorCondition.RESOURCE_LIMIT_EXCEEDED) {
                            //If resource limit is exceeded we will retry after a constant delay
                            log.error("Resource limit exceeded when publishing Event Data Batch. Operation will be "
                                              + "retried after constant delay");
                            try {
                                Thread.sleep(1000 * 60);
                            } catch (InterruptedException interruptedException) {
                                Thread.currentThread().interrupt();
                            }
                            sendEvent(event);
                        } else {
                            //For any other exception
                            log.error("AMQP error occurred while publishing Event Data Batch. Producer client will "
                                              + "be re-initialized. Events may be lost in the process.");
                            log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                              + "- AMQP error occurred while "
                                              + "publishing Event Data Batch. Producer client will "
                                              + "be re-initialized. Events may be lost in the process.", e);
                            this.clientStatus = ClientStatus.RETRYING;
                            createProducerWithRetry(authEndpoint, authToken, retryOptions, true, properties);
                            sendEvent(event);
                        }
                    } catch (Exception e) {
                        if (e.getCause() instanceof TimeoutException) {
                            log.error("Timeout occurred after retrying " + retryOptions.getMaxRetries() + " "
                                              + "times with an timeout of " + retryOptions.getTryTimeout().getSeconds()
                                              + " seconds while trying to publish Event Data Batch. Next retry cycle "
                                              + "will begin shortly.");
                            sendEvent(event);
                        } else {
                            //For any other exception
                            log.error("Unknown error occurred while publishing Event Data Batch. Producer client will "
                                              + "be re-initialized. Events may be lost in the process.");
                            log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                              + "- Unknown error occurred while publishing Event Data Batch. "
                                              + "Producer client will "
                                              + "be re-initialized. Events may be lost in the process.", e);
                            this.clientStatus = ClientStatus.RETRYING;
                            createProducerWithRetry(authEndpoint, authToken, retryOptions, true, properties);
                            sendEvent(event);
                        }
                    }
                }
                if (isAdded) {
                    if (log.isTraceEnabled()) {
                        log.trace("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                          + "- Adding event: " + event.replaceAll("[\r\n]", ""));
                    }
                } else {
                    if (log.isTraceEnabled()) {
                        log.trace("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                          + "- Failed to add event: " + event.replaceAll("[\r\n]", ""));
                    }
                }
            } finally {
                publishingLock.unlock();
            }
        } else if (clientStatus == ClientStatus.FLUSHING_FAILED) {
            log.debug("client status is FLUSHING_FAILED. Producer client will be re-initialized "
                    + "retaining the Event Data Batch");
            this.clientStatus = ClientStatus.RETRYING;
            createProducerWithRetry(authEndpoint, authToken, retryOptions, false, properties);
            sendEvent(event);
        } else {
            try {
                threadBarrier.lock();
                if (log.isDebugEnabled()) {
                    log.debug(Thread.currentThread().getName().replaceAll("[\r\n]", "") + " will be parked as EventHub "
                                      + "Client is inactive.");
                }
                waitCondition.await();
                if (log.isDebugEnabled()) {
                    log.debug(Thread.currentThread().getName().replaceAll("[\r\n]", "") + " will be resumes as "
                                      + "EventHub Client is active.");
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                threadBarrier.unlock();
            }
            sendEvent(event);
        }
    }

    private boolean isAuthenticationFailure(AmqpException exception) {
        AmqpErrorCondition condition = exception.getErrorCondition();
        return (condition == AmqpErrorCondition.UNAUTHORIZED_ACCESS ||
                condition == AmqpErrorCondition.PUBLISHER_REVOKED_ERROR);
    }

    private EventDataBatch createBatchWithRetry() {
        log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "")
                          + " }] Creating Event Data Batch");
        try {
            EventDataBatch batch = producer.createBatch();
            eventBatchRetryCounter.reset();
            return batch;
        } catch (IllegalStateException e) {
            log.error("Error in creating Event Data Batch. Operation will be retried in "
                              + eventBatchRetryCounter.getTimeInterval().replaceAll("[\r\n]", ""));
            try {
                Thread.sleep(eventBatchRetryCounter.getTimeIntervalMillis());
            } catch (InterruptedException interruptedException) {
                Thread.currentThread().interrupt();
            }
            eventBatchRetryCounter.increment();
            return createBatchWithRetry();
        }
    }

    public void flushEvents() {
        if (this.clientStatus == ClientStatus.CONNECTED && batch.getCount() > 0) {
                if (publishingLock.tryLock()) {
                    try {
                        int size = batch.getCount();
                        producer.send(batch);
                        batch = createBatchWithRetry();
                        log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                          + "Flushed " + size + " events to Analytics cluster.");
                    } catch (Exception e) {
                        if (e instanceof AmqpException && isAuthenticationFailure((AmqpException) e)) {
                            log.error("Marked client status as FLUSHING_FAILED due to AMQP authentication failure.");
                            this.clientStatus = ClientStatus.FLUSHING_FAILED;
                        } else if (e.getCause() instanceof ConnectionUnrecoverableException) {
                            this.clientStatus = ClientStatus.NOT_CONNECTED;
                            log.error(
                                    "Unrecoverable error occurred when event flushing. Analytics event flushing will be"
                                            + " disabled until issue is rectified. Reason: " + e.getMessage()
                                            .replaceAll("[\r\n]", ""));
                            return;
                        }
                        log.error("Event flushing operation failed. Will be retried again according to the configured "
                                          + "client.flushing.delay. Error will be handled by publishing threads once "
                                          + "Event Data Batch is filled.");
                        log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] "
                                          + "Event flushing operation failed. Will be retried again according to the "
                                          + "configured client.flushing.delay. Error will be handled by publishing "
                                          + "threads once Event Data Batch is filled.", e);
                        //Dont do anything for any exception. If it is recoverable exception next run will succeed.
                        //If not recoverable then next run will be filtered by if condition
                    } finally {
                        publishingLock.unlock();
                    }
                } else {
                    log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] Event "
                               + "flushing operation aborted as publisher threads are trying to send events");
                }
        } else {
            if (log.isDebugEnabled()) {
                log.debug("[{ " + Thread.currentThread().getName().replaceAll("[\r\n]", "") + " }] Event flushing "
                                  + "is aborted as Event Data Batch is empty or connection to Event Hub is not made");
            }
        }
    }

    public ClientStatus getStatus() {
        return clientStatus;
    }

    /**
     * Returns a clone of this eventhub client.
     *
     * @return New clone of current object
     */
    public EventHubClient clone() {
        return new EventHubClient(this.authEndpoint, this.authToken, this.retryOptions, this.properties);
    }

    private String toHash(String text) {
        if (text == null) {
            log.debug("The text trying to hash is empty.");
            return null;
        }
        final MessageDigest messageDigest;
        try {
            messageDigest = MessageDigest.getInstance("SHA-256");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("Error occurred when getting hash algorithm.", e);
        }
        byte[] digestBytes = messageDigest.digest(text.getBytes(StandardCharsets.UTF_8));
        return DatatypeConverter.printHexBinary(digestBytes).toUpperCase(Locale.ENGLISH);
    }
}
