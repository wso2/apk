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

import com.azure.core.amqp.AmqpRetryOptions;
import com.azure.core.amqp.exception.AmqpErrorCondition;
import com.azure.core.amqp.exception.AmqpException;
import com.azure.messaging.eventhubs.EventData;
import com.azure.messaging.eventhubs.EventDataBatch;
import com.azure.messaging.eventhubs.EventHubProducerClient;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.core.config.Configuration;
import org.junit.Ignore;
import org.mockito.MockedStatic;
import org.mockito.Mockito;
import org.testng.Assert;
import org.testng.annotations.AfterClass;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Optional;
import org.testng.annotations.Parameters;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubProducerClientFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.ConnectionUnrecoverableException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.cloud.DefaultAnalyticsMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.util.AuthAPIMockService;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;
import org.wso2.apk.enforcer.analytics.publisher.util.TestUtils;
import org.wso2.apk.enforcer.analytics.publisher.util.UnitTestAppender;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.TimeoutException;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyMap;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.timeout;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

/**
 * Unit tests related to publisher client and AMQP producer.
 */
public class EventHubClientTestCase extends AuthAPIMockService {

    private EventHubProducerClient client;
    private Map<String, String> configs;
    private UnitTestAppender appender;
    private List<String> messages;
    private MockedStatic<EventHubProducerClientFactory> clientFactoryMocked;

    @BeforeClass
    public void init() {

        clientFactoryMocked = Mockito.mockStatic(EventHubProducerClientFactory.class);
    }

    @AfterClass
    public void finalized() {

        clientFactoryMocked.close();
    }

    @Parameters({"proxyConfigEnabled"})
    @BeforeMethod
    public void setup(@Optional("false") String proxyConfigEnabled) {

        client = Mockito.mock(EventHubProducerClient.class);
        clientFactoryMocked.when(() -> EventHubProducerClientFactory
                        .create(anyString(), anyString(), any(AmqpRetryOptions.class), anyMap()))
                .thenReturn(client);

        String authToken = UUID.randomUUID().toString();
        mock(200, authToken);

        configs = new HashMap<>();
        configs.put(Constants.AUTH_API_URL, authApiEndpoint);
        configs.put(Constants.AUTH_API_TOKEN, authToken);
        if (proxyConfigEnabled.equals("true")) {
            configs.put(Constants.PROXY_PORT, String.valueOf(3128));
            configs.put(Constants.PROXY_HOST, "localhost");
            configs.put(Constants.PROXY_USERNAME, "admin");
            configs.put(Constants.PROXY_PASSWORD, "admin");
        }

        LoggerContext context = LoggerContext.getContext(false);
        Configuration config = context.getConfiguration();
        appender = config.getAppender("UnitTestAppender");
    }

    @Test
    public void testEventFlushing() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        EventDataBatch newEventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch).thenReturn(newEventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1);
        when(newEventDataBatch.getCount()).thenReturn(0);

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting to worker thread adding event to the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));
    }

    @Test
    public void testEventFlushingWithAMQPAuthException() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1);

        doThrow(new AmqpException(false, AmqpErrorCondition.UNAUTHORIZED_ACCESS, "", null))
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter1", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting to adding event to the queue
        verify(eventDataBatch, timeout(10000).times(1)).tryAdd(any(EventData.class));

        // waiting to flushing thread try to send
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // verify flushing thread identified the auth error when try to send via AMQP
        Thread.sleep(1000);
        List<String> appenderMessages = appender.getMessages();
        Assert.assertTrue(TestUtils
                .isContains(appenderMessages, "Marked client status as FLUSHING_FAILED due to AMQP authentication failure."));

        // Try to publish another event
        metric.incrementCount(builder);

        // verify it is also trying to add event queue, after identified state as FLUSHING_FAILED
        verify(eventDataBatch, timeout(10000).times(2)).tryAdd(any(EventData.class));

        // verify worker thread has already identified the FLUSHING_FAILED state
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, "client status is FLUSHING_FAILED. Producer client "
                + "will be re-initialized retaining the Event Data Batch"));
    }

    @Test
    public void testEventFlushingWithConnectionUnrecoverableException() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1);

        doThrow(new RuntimeException(new ConnectionUnrecoverableException("ConnectionUnrecoverableException")))
                .when(client).send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter2", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting to adding event to the queue
        verify(eventDataBatch, timeout(10000).times(1)).tryAdd(any(EventData.class));

        // waiting to flushing thread try to send
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // verify worker thread has already identified the Unrecoverable error
        String msg = "Unrecoverable error occurred when event flushing. Analytics event flushing will be disabled "
                + "until issue is rectified. Reason: org.wso2.apk.enforcer.analytics.publisher.exception."
                + "ConnectionUnrecoverableException: ConnectionUnrecoverableException";
        List<String> appenderMessages = new ArrayList<>(appender.getMessages());
        Assert.assertTrue(TestUtils.isContains(appenderMessages, msg));

        // try to publish another event
        metric.incrementCount(builder);
        // waiting to confirm that the event is not added to the queue
        verify(eventDataBatch, timeout(10000).times(1)).tryAdd(any(EventData.class));
    }

    @Test
    public void testEventPublishingInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        EventDataBatch newEventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch).thenReturn(newEventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false);
        when(newEventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1);
        when(newEventDataBatch.getCount()).thenReturn(1);

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));
        // waiting for worker thread, adding event to the queue
        verify(newEventDataBatch, timeout(10000).times(1)).tryAdd(any(EventData.class));
    }

    @Test
    public void testEventPublishingAndAuthExceptionInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1).thenReturn(0);

        doThrow(new AmqpException(false, AmqpErrorCondition.UNAUTHORIZED_ACCESS, "", null))
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));
        // waiting for worker thread, and verify adding event to the queue get succeed second time
        verify(eventDataBatch, timeout(10000).times(2)).tryAdd(any(EventData.class));

        // verify worker thread has already identified the Unrecoverable error
        String msg = "Authentication issue happened. Producer client will be re-initialized retaining the Event Data "
                + "Batch";
        List<String> appenderMessages = new ArrayList<>(appender.getMessages());
        Assert.assertTrue(TestUtils.isContains(appenderMessages, msg));
    }

    @Test
    public void testEventPublishingAndResourceLimitExceededInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1).thenReturn(0);

        doThrow(new AmqpException(false, AmqpErrorCondition.RESOURCE_LIMIT_EXCEEDED, "", null))
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // waiting for worker thread, and verify adding event to the queue get succeed second time
        verify(eventDataBatch, timeout(100000).times(2)).tryAdd(any(EventData.class));

        // verify worker thread has already identified the Resource limit exceeded error
        String msg = "Resource limit exceeded when publishing Event Data Batch. Operation will be retried after "
                + "constant delay";
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, msg));
    }

    @Ignore
    public void testEventPublishingAndAnyAMQPExceptionInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        EventDataBatch newEventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch).thenReturn(newEventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false);
        when(newEventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1).thenReturn(0);
        when(newEventDataBatch.getCount()).thenReturn(1).thenReturn(0);

        doThrow(new AmqpException(false, AmqpErrorCondition.ARGUMENT_ERROR, "", null))
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // verify existing producer client is close, before create new producer and batch
        verify(client, timeout(10000).times(1)).close();

        // verify worker thread has already identified the amqp exception
        String msg = "AMQP error occurred while publishing Event Data Batch. Producer client will be re-initialized. "
                + "Events may be lost in the process.";
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, msg));
    }

    @Test
    public void testEventPublishingAndTimeOutExceptionInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1).thenReturn(0);

        doThrow(new RuntimeException(new TimeoutException()))
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // waiting for worker thread, and verify adding event to the queue get succeed second time
        verify(eventDataBatch, timeout(20000).times(2)).tryAdd(any(EventData.class));

        // verify worker thread has already identified the Timeout exception
        String msg = "Timeout occurred after retrying 2 times with an timeout of 30 seconds while trying to publish "
                + "Event Data Batch. Next retry cycle will begin shortly.";
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, msg));
    }

    @Test
    public void testEventPublishingAndAnyOtherExceptionInWorkerThread() throws Exception {

        EventDataBatch eventDataBatch = Mockito.mock(EventDataBatch.class);
        EventDataBatch newEventDataBatch = Mockito.mock(EventDataBatch.class);
        when(client.createBatch()).thenReturn(eventDataBatch).thenReturn(newEventDataBatch);
        when(eventDataBatch.tryAdd(any(EventData.class))).thenReturn(false);
        when(newEventDataBatch.tryAdd(any(EventData.class))).thenReturn(true);
        when(eventDataBatch.getCount()).thenReturn(1).thenReturn(0);
        when(newEventDataBatch.getCount()).thenReturn(1).thenReturn(0);

        doThrow(new RuntimeException())
                .when(client)
                .send(any(EventDataBatch.class));

        MetricReporter metricReporter = new DefaultAnalyticsMetricReporter(configs);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);

        // try publishing an event
        metric.incrementCount(builder);

        // waiting for worker thread flush the queue
        verify(client, timeout(20000).times(1)).send(any(EventDataBatch.class));

        // verify existing producer client is close, before create new producer and batch
        verify(client, timeout(10000).times(1)).close();

        // verify worker thread has already identified the runtime exception
        String msg = "Unknown error occurred while publishing Event Data Batch. Producer client will be re-initialized."
                + " Events may be lost in the process.";
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, msg));
    }
}
