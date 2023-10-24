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

import com.moesif.api.MoesifAPIClient;
import com.moesif.api.controllers.APIController;
import com.moesif.api.http.client.APICallBack;
import com.moesif.api.http.client.HttpContext;
import com.moesif.api.http.response.HttpResponse;
import com.moesif.api.models.EventModel;
import com.moesif.api.models.EventRequestBuilder;
import com.moesif.api.models.EventRequestModel;
import com.moesif.api.models.EventResponseBuilder;
import com.moesif.api.models.EventResponseModel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.util.MoesifMicroserviceConstants;
import org.wso2.apk.enforcer.analytics.publisher.retriever.MoesifKeyRetriever;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.io.IOException;
import java.time.OffsetDateTime;
import java.util.Date;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Moesif Client is responsible for sending events to
 * Moesif Analytics Dashboard.
 */
public class MoesifClient {

    private final Logger log = LoggerFactory.getLogger(MoesifClient.class);
    private final MoesifKeyRetriever keyRetriever;

    public MoesifClient(MoesifKeyRetriever keyRetriever) {

        this.keyRetriever = keyRetriever;
    }

    private void doRetry(String orgId, MetricEventBuilder builder) {

        Integer currentAttempt = MoesifClientContextHolder.PUBLISH_ATTEMPTS.get();

        if (currentAttempt > 0) {
            currentAttempt -= 1;
            MoesifClientContextHolder.PUBLISH_ATTEMPTS.set(currentAttempt);
            try {
                Thread.sleep(MoesifMicroserviceConstants.TIME_TO_WAIT_PUBLISH);
                publish(builder);
            } catch (MetricReportingException e) {
                log.error("Failing retry attempt at Moesif client", e);
            } catch (InterruptedException e) {
                log.error("Failing retry attempt at Moesif client", e);
            }
        } else if (currentAttempt == 0) {
            log.error("Failed all retrying attempts. Event will be dropped for organization {}",
                    orgId.replaceAll("[\r\n]", ""));
        }
    }

    /**
     * publish method is responsible for checking the availability of relevant moesif key
     * and initiating moesif client sdk.
     */
    public void publish(MetricEventBuilder builder) throws MetricReportingException {

        Map<String, Object> event = builder.build();

        String orgId = (String) event.get(Constants.ORGANIZATION_ID);
        String eventEnvironment = (String) event.getOrDefault(Constants.ENVIRONMENT_ID, Constants.DEFAULT_ENVIRONMENT);
        String moesifKey = keyRetriever.getKey(orgId, eventEnvironment);

        // init moesif api client
        MoesifAPIClient client = new MoesifAPIClient(moesifKey);

        // api object is a singleton which will make calls to
        // moesif endpoints with the latest MoesifAPI client being provided.
        // Hence avoid maintaining a map of MoesifAPIClient against moesif keys.
        APIController api = client.getAPI();

        APICallBack<HttpResponse> callBack = new APICallBack<HttpResponse>() {
            public void onSuccess(HttpContext context, HttpResponse response) {

                int statusCode = context.getResponse().getStatusCode();
                if (statusCode == 200 || statusCode == 201 || statusCode == 202 || statusCode == 204) {
                    log.debug("Event successfully published.");
                } else if (statusCode >= 400 && statusCode < 500) {
                    log.error("Event publishing failed for organization: {}. Moesif returned {}.",
                            orgId.replaceAll("[\r\n]", ""), String.valueOf(statusCode).replaceAll("[\r\n]", ""));
                } else {
                    log.error("Event publishing failed for organization: {}. Retrying.",
                            orgId.replaceAll("[\r\n]", ""));
                    doRetry(orgId, builder);
                }
            }

            public void onFailure(HttpContext context, Throwable error) {

                int statusCode = context.getResponse().getStatusCode();

                if (statusCode >= 400 && statusCode < 500) {
                    log.error("Event publishing failed for organization: {}. Moesif returned {}.",
                            orgId.replaceAll("[\r\n]", ""), String.valueOf(statusCode).replaceAll("[\r\n]", ""));
                } else if (error != null) {
                    log.error("Event publishing failed for organization: {}. Event publishing failed.",
                            orgId.replaceAll("[\r\n]", ""),
                            error.getMessage().replaceAll("[\r\n]", ""));
                } else {
                    log.error("Event publishing failed for organization: {}. Retrying.",
                            orgId.replaceAll("[\r\n]", ""));
                    doRetry(orgId, builder);
                }
            }
        };
        try {
            api.createEventAsync(buildEventResponse(event), callBack);
        } catch (IOException e) {
            log.error("Analytics event sending failed. Event will be dropped", e);
        }
    }

    private EventModel buildEventResponse(Map<String, Object> data) throws IOException, MetricReportingException {
        //      Preprocessing data
        final String userIP = (String) data.get(Constants.USER_IP);
        final String userName = (String) data.get(Constants.USER_NAME);
        final String apiContext = (String) data.get(Constants.API_CONTEXT);
        final String apiResourceTemplate = (String) data.get(Constants.API_RESOURCE_TEMPLATE);
        final long responseLatency = (long) data.get(Constants.RESPONSE_LATENCY);
        final String requestTimeStamp = (String) data.get(Constants.REQUEST_TIMESTAMP);
        Map<String, String> reqHeaders = new HashMap<String, String>();

        reqHeaders.put(Constants.MOESIF_USER_AGENT_KEY,
                (String) data.getOrDefault(Constants.USER_AGENT_HEADER, Constants.UNKNOWN_VALUE));
        reqHeaders.put(Constants.MOESIF_CONTENT_TYPE_KEY, Constants.MOESIF_CONTENT_TYPE_HEADER);

        Map<String, String> rspHeaders = new HashMap<String, String>();

        rspHeaders.put("Vary", "Accept-Encoding");
        rspHeaders.put("Pragma", "no-cache");
        rspHeaders.put("Expires", "-1");
        rspHeaders.put(Constants.MOESIF_CONTENT_TYPE_KEY, "application/json; charset=utf-8");
        rspHeaders.put("Cache-Control", "no-cache");

        LinkedHashMap properties = (LinkedHashMap) data.get(Constants.PROPERTIES);
        String gwURL = (String) properties.get(Constants.GATEWAY_URL);
        String uri = apiContext + apiResourceTemplate;
        if (gwURL != null) {
            uri = gwURL;
        }
        Date requestDate = getDate(requestTimeStamp);
        Date responseDate = new Date(requestDate.getTime() + responseLatency);
        EventRequestModel eventReq = new EventRequestBuilder()
                .time(requestDate)
                .uri(uri)
                .verb((String) data.get(Constants.API_METHOD))
                .apiVersion((String) data.get(Constants.API_VERSION))
                .ipAddress(userIP)
                .headers(reqHeaders)
                .build();

        EventResponseModel eventRsp = new EventResponseBuilder()
                .time(responseDate)
                .status((int) data.get(Constants.PROXY_RESPONSE_CODE))
                .headers(rspHeaders)
                .build();

        EventModel eventModel = new EventModel();

        eventModel.setRequest(eventReq);
        eventModel.setResponse(eventRsp);
        eventModel.setUserId(userName);
        eventModel.setCompanyId(null);

        return eventModel;
    }

    private static Date getDate(String time) {

        OffsetDateTime offsetDateTime = OffsetDateTime.parse(time);
        long milliSecondsTime = offsetDateTime.toInstant().toEpochMilli();
        return new Date(milliSecondsTime);
    }
}
