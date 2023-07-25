/*
 * Copyright (c) 2023, WSO2 LLC (http://www.wso2.com).
 *
 * WSO2 LLC licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.apk.integration.api;

import io.cucumber.core.options.CurlOption;
import io.cucumber.datatable.DataTable;
import io.cucumber.java.Before;
import io.cucumber.java.en.And;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import org.apache.http.Header;
import org.apache.http.HttpResponse;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import java.io.IOException;
import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * This class contains the common step definitions.
 */
public class BaseSteps {

    private final SharedContext sharedContext;
    private SimpleHTTPClient httpClient;
    private static final int MAX_WAIT_FOR_NEXT_MINUTE_IN_SECONDS = 10;

    public BaseSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @Before
    public void setup() throws Exception {

        httpClient = sharedContext.getHttpClient();
    }

    @Given("The system is ready")
    public void systemIsReady() {

    }

    @Then("the response body should contain {string}")
    public void theResponseBodyShouldContain(String expectedText) throws IOException {

        Assert.assertTrue(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()).contains(expectedText));
    }

    @Then("the response status code should be {int}")
    public void theResponseStatusCodeShouldBe(int expectedStatusCode) {

        int actualStatusCode = sharedContext.getResponse().getStatusLine().getStatusCode();
        Assert.assertEquals(actualStatusCode, expectedStatusCode);
    }

    @Then("I send {string} request to {string} with body {string}")
    public void sendHttpRequest(String httpMethod, String url, String body) throws IOException {
        if (sharedContext.getResponse() instanceof CloseableHttpResponse) {
            ((CloseableHttpResponse) sharedContext.getResponse()).close();
        }
        if (CurlOption.HttpMethod.GET.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doGet(url, sharedContext.getHeaders()));
        } else if (CurlOption.HttpMethod.POST.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doPost(url, sharedContext.getHeaders(), body, null));
        } else if (CurlOption.HttpMethod.PUT.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doPut(url, sharedContext.getHeaders(), body, null));
        } else if (CurlOption.HttpMethod.DELETE.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doPut(url, sharedContext.getHeaders(), body, null));
        }
    }

    @Then("I set headers")
    public void setHeaders(DataTable dataTable) {
        List<List<String>> rows = dataTable.asLists(String.class);
        for (List<String> columns : rows) {
            String key = columns.get(0);
            String value = columns.get(1);
            key = Utils.resolveVariables(key, sharedContext.getValueStore());
            value = Utils.resolveVariables(value, sharedContext.getValueStore());
            sharedContext.addHeader(key, value);
        }
    }

    @Then("I eventually receive {int} response code, not accepting")
    public void eventualSuccess(int statusCode, DataTable dataTable) throws IOException, InterruptedException {
        List<Integer> nonAcceptableCodes = dataTable.asList(Integer.class);
        if (sharedContext.getResponse().getStatusLine().getStatusCode() == statusCode) {
            Assert.assertTrue(true);
        } else {
            HttpResponse httpResponse = httpClient.executeLastRequestForEventualConsistentResponse(statusCode,
                    nonAcceptableCodes);
            sharedContext.setResponse(httpResponse);
            Assert.assertEquals(httpResponse.getStatusLine().getStatusCode(), statusCode);
        }
    }

    @Then("I wait for next minute")
    public void waitForNextMinute() throws InterruptedException {
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime nextMinute = now.plusMinutes(1).withSecond(0).withNano(0);
        long secondsToWait = now.until(nextMinute, ChronoUnit.SECONDS);
        if (secondsToWait > MAX_WAIT_FOR_NEXT_MINUTE_IN_SECONDS) {
            return;
        }
        Thread.sleep((secondsToWait+1) * 1000);
    }

    @Then("the response headers contains key {string} and value {string} ")
    public void containsHeader(String key, String value) {
        key = Utils.resolveVariables(key, sharedContext.getValueStore());
        value = Utils.resolveVariables(value, sharedContext.getValueStore());
        HttpResponse response = sharedContext.getResponse();
        if (response == null) {
            Assert.fail("Response is null.");
        }
        Header header = response.getFirstHeader(key);
        if (header == null) {
            Assert.fail("Could not find a header with the given key: " + key);
        }
        if ("*".equals(value)) {
            return; // Any value is acceptable
        }
        String actualValue = header.getValue();
        Assert.assertEquals("Header with key found but value mismatched.", value, actualValue);
    }


    @Given("I have a valid subscription")
    public void iHaveValidSubscription() throws Exception {

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

        HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers, "grant_type=client_credentials",
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
        sharedContext.setAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());
    }
}
