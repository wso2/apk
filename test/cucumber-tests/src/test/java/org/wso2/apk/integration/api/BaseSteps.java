/*
 * Copyright (c) 2024, WSO2 LLC (http://www.wso2.com).
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

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JOSEObjectType;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.jwk.source.JWKSourceBuilder;
import com.nimbusds.jose.proc.BadJOSEException;
import com.nimbusds.jose.proc.DefaultJOSEObjectTypeVerifier;
import com.nimbusds.jose.proc.JWSKeySelector;
import com.nimbusds.jose.proc.JWSVerificationKeySelector;
import com.nimbusds.jose.proc.SecurityContext;
import com.nimbusds.jose.util.Resource;
import com.nimbusds.jose.util.ResourceRetriever;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.proc.ConfigurableJWTProcessor;
import com.nimbusds.jwt.proc.DefaultJWTProcessor;
import io.cucumber.core.options.CurlOption;
import io.cucumber.datatable.DataTable;
import io.cucumber.java.Before;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.StatusLine;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleGRPCStudentClient;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentResponse;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.ContentType;

import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URL;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.text.ParseException;
import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * This class contains the common step definitions.
 */
public class BaseSteps {

    private static final Log logger = LogFactory.getLog(BaseSteps.class);
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

    @Then("the student response body should contain name: {string} age: {int}")
    public void theStudentResponseBodyShouldContainNameAndAge(String arg0, int arg1) {
        StudentResponse studentResponse = sharedContext.getStudentResponse();
        if (studentResponse == null) {
            Assert.fail("Student response is null.");
        }
        int age = studentResponse.getAge();
        String name = studentResponse.getName();
        Assert.assertEquals(name, arg0);
        Assert.assertEquals(age, arg1);
    }

    @Then("the response body should contain {string}")
    public void theResponseBodyShouldContain(String expectedText) throws IOException {
        Assert.assertTrue(sharedContext.getResponseBody().contains(expectedText), "Actual response body: " + sharedContext.getResponseBody());
    }
    @Then("the response body should not contain {string}")
    public void theResponseBodyShouldNotContain(String expectedText) throws IOException {
        Assert.assertFalse(sharedContext.getResponseBody().contains(expectedText), "Actual response body: " + sharedContext.getResponseBody());
    }

    @Then("the response body should contain")
    public void theResponseBodyShouldContain(DataTable dataTable) throws IOException {
        List<String> responseBodyLines = dataTable.asList(String.class);
        for (String line : responseBodyLines) {
            Assert.assertTrue(sharedContext.getResponseBody().contains(line), "Actual response body: " + sharedContext.getResponseBody());
        }
    }

    @Then("the response status code should be {int}")
    public void theResponseStatusCodeShouldBe(int expectedStatusCode) throws IOException {

        int actualStatusCode = sharedContext.getResponse().getStatusLine().getStatusCode();
        ((CloseableHttpResponse)sharedContext.getResponse()).close();
        Assert.assertEquals(actualStatusCode, expectedStatusCode);
    }

    @Then("the grpc error response status code should be {int}")
    public void theGrpcErrorResponseStatusCodeShouldBe(int expectedStatusCode) throws IOException {

        int actualStatusCode = sharedContext.getGrpcErrorCode();
        Assert.assertEquals(actualStatusCode, expectedStatusCode);
    }

    @Then("I send {string} request to {string} with body {string}")
    public void sendHttpRequest(String httpMethod, String url, String body) throws IOException {
        body = Utils.resolveVariables(body, sharedContext.getValueStore());
        if (sharedContext.getResponse() instanceof CloseableHttpResponse) {
            ((CloseableHttpResponse) sharedContext.getResponse()).close();
        }
        if (CurlOption.HttpMethod.GET.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doGet(url, sharedContext.getHeaders()));
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        } else if (CurlOption.HttpMethod.POST.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doPost(url, sharedContext.getHeaders(), body, null));
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        } else if (CurlOption.HttpMethod.PUT.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doPut(url, sharedContext.getHeaders(), body, null));
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        } else if (CurlOption.HttpMethod.DELETE.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doDelete(url, sharedContext.getHeaders()));
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        } else if (CurlOption.HttpMethod.OPTIONS.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
            sharedContext.setResponse(httpClient.doOptions(url, sharedContext.getHeaders(), null, null));
        }
    }

    @Then("I make grpc request GetStudent to {string} with port {int}")
    public void GetStudent(String arg0, int arg1) throws StatusRuntimeException {
        try {
            SimpleGRPCStudentClient grpcStudentClient = new SimpleGRPCStudentClient(arg0, arg1);
            sharedContext.setStudentResponse(grpcStudentClient.GetStudent(sharedContext.getHeaders()));
        } catch (StatusRuntimeException e) {
            if (e.getStatus().getCode()== Status.Code.PERMISSION_DENIED){
                sharedContext.setGrpcErrorCode(403);
            } else {
                logger.error(e.getMessage());
            }
        }
    }

    @Then("I make grpc request GetStudent default version to {string} with port {int}")
    public void GetStudentDefaultVersion(String arg0, int arg1) throws StatusRuntimeException {
        try {
            SimpleGRPCStudentClient grpcStudentClient = new SimpleGRPCStudentClient(arg0, arg1);
            sharedContext.setStudentResponse(grpcStudentClient.GetStudent(sharedContext.getHeaders()));
        } catch (StatusRuntimeException e) {
            if (e.getStatus().getCode()== Status.Code.PERMISSION_DENIED){
                sharedContext.setGrpcErrorCode(403);
            } else {
                logger.error(e.getMessage());
            }
        }
    }

    // It will send request using a new thread and forget about the response
    @Then("I send {string} async request to {string} with body {string}")
    public void sendAsyncHttpRequest(String httpMethod, String url, String body) throws IOException, NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        String finalBody = Utils.resolveVariables(body, sharedContext.getValueStore());
        if (sharedContext.getResponse() instanceof CloseableHttpResponse) {
            ((CloseableHttpResponse) sharedContext.getResponse()).close();
        }
        SimpleHTTPClient simpleHTTPClient = new SimpleHTTPClient();
        Thread thread = new Thread(() -> {
            try {
                if (CurlOption.HttpMethod.GET.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
                    simpleHTTPClient.doGet(url, sharedContext.getHeaders());
                } else if (CurlOption.HttpMethod.POST.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
                    simpleHTTPClient.doPost(url, sharedContext.getHeaders(), finalBody, null);
                } else if (CurlOption.HttpMethod.PUT.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
                    simpleHTTPClient.doPut(url, sharedContext.getHeaders(), finalBody, null);
                } else if (CurlOption.HttpMethod.DELETE.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
                    simpleHTTPClient.doPut(url, sharedContext.getHeaders(), finalBody, null);
                } else if (CurlOption.HttpMethod.OPTIONS.toString().toLowerCase().equals(httpMethod.toLowerCase())) {
                    simpleHTTPClient.doOptions(url, sharedContext.getHeaders(), null, null);
                }
            } catch (IOException e) {
                logger.warn("An async http request sending thread experienced an error: " + e);
            }
        });
        thread.start();
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

    @Then("I remove header {string}")
    public void setHeaders(String headerName) {
        String resolvedHeaderName = Utils.resolveVariables(headerName, sharedContext.getValueStore());
        sharedContext.removeHeader(resolvedHeaderName);
    }

    @Then("the response headers should contain")
    public void theResponseHeadersShouldContain(DataTable dataTable) {
        List<List<String>> rows = dataTable.asLists(String.class);
        for (List<String> columns : rows) {
            String key = columns.get(0);
            String value = columns.get(1);
            Header header = sharedContext.getResponse().getFirstHeader(key);
            Assert.assertNotNull(header);
            Assert.assertEquals(header.getValue(), value);
        }
    }

    @Then("the response headers should not contain")
    public void theResponseHeadersShouldNotContain(DataTable dataTable) {
        List<List<String>> rows = dataTable.asLists(String.class);
        for (List<String> columns : rows) {
            String key = columns.get(0);
            Header header = sharedContext.getResponse().getFirstHeader(key);
            Assert.assertNull(header);
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
        logger.info("Current time: " + LocalDateTime.now());
    }

    @Then("I wait for next minute strictly")
    public void waitForNextMinuteStrictly() throws InterruptedException {
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime nextMinute = now.plusMinutes(1).withSecond(0).withNano(0);
        long secondsToWait = now.until(nextMinute, ChronoUnit.SECONDS);
        Thread.sleep((secondsToWait+1) * 1000);
        logger.info("Current time: " + LocalDateTime.now());
    }

    @Then("I wait for {int} minute")
    public void waitForMinute(int minute) throws InterruptedException {
        Thread.sleep(minute * 1000);
    }

    @Then("I wait for {int} seconds")
    public void waitForSeconds(int seconds) throws InterruptedException {
        Thread.sleep(seconds * 1000);
    }

    @Then("the response headers contains key {string} and value {string}")
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
        Assert.assertEquals(value, actualValue,"Header with key found but value mismatched.");
    }
    @Then("the response headers not contains key {string}")
    public void notContainsHeader(String key) {
        key = Utils.resolveVariables(key, sharedContext.getValueStore());
        HttpResponse response = sharedContext.getResponse();
        if (response == null) {
            Assert.fail("Response is null.");
        }
        Header header = response.getFirstHeader(key);
        Assert.assertNull(header,"header contains in response headers");
    }

    @Then("the {string} jwt should validate from JWKS {string} and contain")
    public void decode_header_and_validate(String header,String jwksEndpoint, DataTable dataTable) throws MalformedURLException {
        List<Map<String, String>> claims = dataTable.asMaps(String.class, String.class);
        JsonObject jsonResponse = (JsonObject) JsonParser.parseString(sharedContext.getResponseBody());
        String headerValue = jsonResponse.get("headers").getAsJsonObject().get(header).getAsString();
        ConfigurableJWTProcessor<SecurityContext> jwtProcessor = new DefaultJWTProcessor<>();
        jwtProcessor.setJWSTypeVerifier(new DefaultJOSEObjectTypeVerifier<>(JOSEObjectType.JWT));
        ResourceRetriever retriever = url -> {
            try {
                HttpResponse httpResponse = new SimpleHTTPClient().doGet(url.toString(), Collections.emptyMap());
                StatusLine statusLine = httpResponse.getStatusLine();
                if (statusLine.getStatusCode() == 200) {
                    Header header1 = httpResponse.getFirstHeader("Content-Type");
                    try (InputStream content = httpResponse.getEntity().getContent()) {
                        return new Resource(IOUtils.toString(content), header1.getValue());
                    }
                } else {
                    throw new IOException("HTTP " + statusLine.getStatusCode() + ": " + statusLine.getReasonPhrase());
                }
            } catch (NoSuchAlgorithmException | KeyStoreException | KeyManagementException e) {
                throw new IOException(e);
            }
        };

        JWKSource<SecurityContext> keySource = JWKSourceBuilder.create(new URL(jwksEndpoint), retriever).build();
        JWSAlgorithm expectedJWSAlg = JWSAlgorithm.RS256;
        JWSKeySelector<SecurityContext> keySelector = new JWSVerificationKeySelector<>(expectedJWSAlg, keySource);
        jwtProcessor.setJWSKeySelector(keySelector);
        try {
            JWTClaimsSet claimsSet = jwtProcessor.process(headerValue, null);
            for (Map<String, String> claim : claims) {
                Object claim1 = claimsSet.getClaim(claim.get("claim"));
                Assert.assertNotNull(claim1, "Actual decoded JWT body: " + claimsSet);
                Assert.assertEquals(claim.get("value"), claim1.toString(), "Actual " +
                        "decoded JWT body: " + claimsSet);
            }
        } catch (BadJOSEException | JOSEException|ParseException e) {
            logger.error("JWT Signature verification fail", e);
            Assert.fail("JWT Signature verification fail");
        }
    }

    @Given("I have a valid subscription")
    public void iHaveValidSubscription() throws Exception {

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

        HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers, "grant_type=client_credentials&scope=" + Constants.API_CREATE_SCOPE,
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
        sharedContext.setAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());
        logger.info("Access Token: " + sharedContext.getAccessToken());
    }

    @Given("I have a valid subscription without api deploy permission")
    public void iHaveValidSubscriptionWithAPICreateScope() throws Exception {

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==");

        HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers, "grant_type=client_credentials",
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
        sharedContext.setAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("accessToken", sharedContext.getAccessToken());
    }

    @Given("I have a valid subscription with scopes")
    public void iHaveValidSubscriptionWithScope(DataTable dataTable) throws Exception {
        List<List<String>> rows = dataTable.asLists(String.class);
        String scopes = Constants.EMPTY_STRING;
        for (List<String> row : rows) {
            String scope = row.get(0);
            scopes += scope + Constants.SPACE_STRING;
        }
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, Constants.SUBSCRIPTION_BASIC_AUTH_TOKEN);

        HttpResponse httpResponse = httpClient.doPost(Utils.getTokenEndpointURL(), headers,
                                                      "grant_type=client_credentials&scope=" + scopes,
                                                      Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
        sharedContext.setAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue(Constants.ACCESS_TOKEN, sharedContext.getAccessToken());
    }

    @Then("I remove the header {string}")
    public void removeHeader(String key) {
        sharedContext.removeHeader(key);
    }

    @Given("I have a DCR application")
    public void iHaveADCRApplication() throws Exception {

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_APIM_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Basic YWRtaW46YWRtaW4=");

        HttpResponse httpResponse = httpClient.doPost(Utils.getDCREndpointURL(), headers, "{\n" +
                        "  \"callbackUrl\":\"www.google.lk\",\n" +
                        "  \"clientName\":\"rest_api_publisher\",\n" +
                        "  \"owner\":\"admin\",\n" +
                        "  \"grantType\":\"client_credentials password refresh_token\",\n" +
                        "  \"saasApp\":true\n" +
                        "  }",
                Constants.CONTENT_TYPES.APPLICATION_JSON);
        sharedContext.setBasicAuthToken(Utils.extractBasicToken(httpResponse));
        sharedContext.addStoreValue("publisherBasicAuthToken", sharedContext.getBasicAuthToken());
    }


    @Given("I have a valid Publisher access token")
    public void iHaveValidPublisherAccessToken() throws Exception {

        Map<String, String> headers = new HashMap<>();
        String basicAuthHeader = "Basic " + sharedContext.getBasicAuthToken();
        logger.info("Basic Auth Header: " + basicAuthHeader);
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_APIM_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, basicAuthHeader);

        HttpResponse httpResponse = httpClient.doPost(Utils.getAPIMTokenEndpointURL(), headers, "grant_type=password&username=admin&password=admin&scope=apim:api_view apim:api_create apim:api_publish apim:api_delete apim:api_manage apim:api_import_export apim:subscription_manage apim:client_certificates_add apim:client_certificates_update",
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);

        sharedContext.setPublisherAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("publisherAccessToken", sharedContext.getPublisherAccessToken());
    }

    @Given("I have a valid Devportal access token")
    public void iHaveValidDevportalAccessToken() throws Exception {
        logger.info("Basic Auth Header: " + sharedContext.getBasicAuthToken());

        Map<String, String> headers = new HashMap<>();
        String basicAuthHeader = "Basic " + sharedContext.getBasicAuthToken();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_APIM_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, basicAuthHeader);

        HttpResponse httpResponse = httpClient.doPost(Utils.getAPIMTokenEndpointURL(), headers, "grant_type=password&username=admin&password=admin&scope=apim:app_manage apim:sub_manage apim:subscribe",
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);

        sharedContext.setDevportalAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("devportalAccessToken", sharedContext.getDevportalAccessToken());
        logger.info("Devportal Access Token: " + sharedContext.getDevportalAccessToken());
    }

    @Given("I have a valid Adminportal access token")
    public void iHaveValidAdminportalAccessToken() throws Exception {
        logger.info("Basic Auth Header: " + sharedContext.getBasicAuthToken());

        Map<String, String> headers = new HashMap<>();
        String basicAuthHeader = "Basic " + sharedContext.getBasicAuthToken();
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_APIM_IDP_HOST);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, basicAuthHeader);

        HttpResponse httpResponse = httpClient.doPost(Utils.getAPIMTokenEndpointURL(), headers, "grant_type=password&username=admin&password=admin&scope=apim:app_manage apim:admin_tier_view apim:admin_tier_manage",
                Constants.CONTENT_TYPES.APPLICATION_X_WWW_FORM_URLENCODED);
        sharedContext.setAdminAccessToken(Utils.extractToken(httpResponse));
        sharedContext.addStoreValue("adminportalAccessToken", sharedContext.getAdminAccessToken());
        logger.info("Admin Access Token: " + sharedContext.getAdminAccessToken());
    }

    @Then("the response should be given as valid")
    public void theResponseShouldBeGivenAs() throws IOException {
        Boolean status = sharedContext.getDefinitionValidStatus();
        Assert.assertEquals(true, status,"Actual definition validation status: "+ status);
    }

    @Then("I set {string} as the new access token")
    public void set_invalid_access_token(String newToken) throws Exception {
            sharedContext.setApiAccessToken(newToken);
            sharedContext.addStoreValue("accessToken",sharedContext.getApiAccessToken());
    }
}
