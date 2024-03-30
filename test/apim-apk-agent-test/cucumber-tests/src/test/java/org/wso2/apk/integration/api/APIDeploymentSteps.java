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

import com.google.common.io.Resources;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import io.cucumber.java.en.Given;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.message.BasicNameValuePair;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;
import java.nio.file.Files;
import java.nio.charset.StandardCharsets;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.net.URL;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * This class contains the step definitions for API Deployment.
 */
public class APIDeploymentSteps {

    private final SharedContext sharedContext;
    private File payloadFile;
    private File definitionFile;

    private String OASURL;

    private static final Log logger = LogFactory.getLog(APIDeploymentSteps.class);

    public APIDeploymentSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @When("I use the Payload file {string}")
    public void iHaveTheAPKConf(String payloadFileName) throws IOException {

        URL url = Resources.getResource(payloadFileName);
        payloadFile = new File(url.getPath());
    }

    @When("I use the OAS URL {string}")
    public void iHaveTheOASURL(String pOASURL) throws IOException {
        OASURL = pOASURL;
    }

    @When("the definition file {string}")
    public void iHaveTheDefinitionFile(String definitionFileName) throws IOException {

        URL url = Resources.getResource(definitionFileName);
        definitionFile = new File(url.getPath());
    }

    @When("make the import API Creation request using OAS {string}")
    public void make_import_api_creation_request(String definitionType) throws Exception {
        MultipartEntityBuilder builder = null;
        if(definitionType.equals("URL")){
                logger.info("OAS URL: " + OASURL);
                builder = MultipartEntityBuilder.create()
                        .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                        .addTextBody("url", OASURL, ContentType.TEXT_PLAIN)
                        .addPart("additionalProperties", new FileBody(payloadFile));

                logger.info("Payload File: "+ new FileBody(payloadFile));
        }
        if(definitionType.equals("File")){
                logger.info("OAS File: " + definitionFile.getName());
                builder = MultipartEntityBuilder.create()
                        .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                        .addPart("file", new FileBody(definitionFile))
                        .addPart("additionalProperties", new FileBody(payloadFile));

                logger.info("Payload File: "+ new FileBody(payloadFile));
        }


        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpEntity multipartEntity = builder.build();

        HttpResponse response = sharedContext.getHttpClient().doPostWithMultipart(Utils.getImportAPIURL(),
                multipartEntity, headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setApiUUID(Utils.extractID(sharedContext.getResponseBody()));
        Thread.sleep(3000);
    }

    @When("make the API Revision Deployment request")
    public void make_a_api_revision_deployment_request() throws Exception {
        String apiUUID = sharedContext.getApiUUID();
        String payload = "{\"description\":\"Initial Revision\"}";

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getAPIRevisionURL(apiUUID),
                 headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setRevisionUUID(Utils.extractID(SimpleHTTPClient.responseEntityBodyToString(response)));

        Thread.sleep(3000);

        String payload2 = "[{\"name\": \"Default\", \"vhost\": \"default.gw.wso2.com\", \"displayOnDevportal\": true}]";

        HttpResponse response2 = sharedContext.getHttpClient().doPost(Utils.getAPIRevisionDeploymentURL(apiUUID, sharedContext.getRevisionUUID()),
                headers, payload2, Constants.CONTENT_TYPES.APPLICATION_JSON);

        logger.info("Response: "+ response2);

        sharedContext.setResponse(response2);
        Thread.sleep(3000);
    }

    @When("make the Change Lifecycle request")
    public void make_a_change_lifecycle_request() throws Exception {
        String apiUUID = sharedContext.getApiUUID();
        String payload = "";

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getAPIChangeLifecycleURL(apiUUID),
                headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        Thread.sleep(3000);
    }

    @When("make the Application Creation request with the name {string}")
    public void make_application_creation_request(String applicationName) throws Exception {
        logger.info("Creating an application");
        String payload = "{\"name\":\"" + applicationName + "\",\"throttlingPolicy\":\"10PerMin\",\"description\":\"test app\",\"tokenType\":\"JWT\",\"groups\":null,\"attributes\":{}}";
    
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);
    
        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getApplicationCreateURL(),
                headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);
    
        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        logger.info("Response: " + sharedContext.getResponseBody());
        sharedContext.setApplicationUUID(Utils.extractApplicationID(sharedContext.getResponseBody()));
        Thread.sleep(3000);
    }
    

    @When("I have a KeyManager")
    public void i_have_a_key_manager() throws Exception {
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doGet(Utils.getKeyManagerURL(),
                headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setKeyManagerUUID(Utils.extractKeyManagerID(sharedContext.getResponseBody()));
        Thread.sleep(3000);
    }

    @When("make the Generate Keys request")
    public void make_generate_keys_request() throws Exception {
        String applicationUUID = sharedContext.getApplicationUUID();
        String keyManagerUUID = sharedContext.getKeyManagerUUID();
        logger.info("Key Manager UUID: " + keyManagerUUID);
        logger.info("Application UUID: " + applicationUUID);
        String payloadForProdKeys = "{\"keyType\":\"PRODUCTION\",\"grantTypesToBeSupported\":[\"password\",\"client_credentials\"]," +
                "\"callbackUrl\":\"\",\"additionalProperties\":{\"application_access_token_expiry_time\":\"N/A\"," +
                "\"user_access_token_expiry_time\":\"N/A\",\"refresh_token_expiry_time\":\"N/A\"," +
                "\"id_token_expiry_time\":\"N/A\",\"pkceMandatory\":\"false\",\"pkceSupportPlain\":\"false\"," +
                "\"bypassClientCredentials\":\"false\"},\"keyManager\":\"" + keyManagerUUID +"\"," +
                "\"validityTime\":3600,\"scopes\":[\"default\"]}";

        
        String payloadForSandboxKeys = "{\"keyType\":\"SANDBOX\",\"grantTypesToBeSupported\":[\"password\",\"client_credentials\"]," +
                "\"callbackUrl\":\"\",\"additionalProperties\":{\"application_access_token_expiry_time\":\"N/A\"," +
                "\"user_access_token_expiry_time\":\"N/A\",\"refresh_token_expiry_time\":\"N/A\"," +
                "\"id_token_expiry_time\":\"N/A\",\"pkceMandatory\":\"false\",\"pkceSupportPlain\":\"false\"," +
                "\"bypassClientCredentials\":\"false\"},\"keyManager\":\"" + keyManagerUUID +"\"," +
                "\"validityTime\":3600,\"scopes\":[\"default\"]}";

       
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getGenerateKeysURL(applicationUUID),
                headers, payloadForProdKeys, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setConsumerSecret(Utils.extractKeys(sharedContext.getResponseBody(), "consumerSecret"), "production");
        sharedContext.setConsumerKey(Utils.extractKeys(sharedContext.getResponseBody(), "consumerKey"), "production");
        sharedContext.setKeyMappingID(Utils.extractKeys(sharedContext.getResponseBody(), "keyMappingId"), "production");
        Thread.sleep(3000);

        HttpResponse response2 = sharedContext.getHttpClient().doPost(Utils.getGenerateKeysURL(applicationUUID),
                headers, payloadForSandboxKeys, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response2);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setConsumerSecret(Utils.extractKeys(sharedContext.getResponseBody(), "consumerSecret"), "sandbox");
        sharedContext.setConsumerKey(Utils.extractKeys(sharedContext.getResponseBody(), "consumerKey"), "sandbox");
        sharedContext.setKeyMappingID(Utils.extractKeys(sharedContext.getResponseBody(), "keyMappingId"), "sandbox");
        Thread.sleep(3000);
    }

    @When("make the Subscription request")
    public void make_subscription_request() throws Exception {
        String applicationUUID = sharedContext.getApplicationUUID();
        String apiUUID = sharedContext.getApiUUID();
        logger.info("API UUID: " + apiUUID);
        logger.info("Application UUID: " + applicationUUID);
        String payload = "{\"apiId\":\"" + apiUUID + "\",\"applicationId\":\"" + applicationUUID + "\",\"throttlingPolicy\":\"Unlimited\"}";

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getSubscriptionURL(),
                headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        Thread.sleep(3000);
    }

    @When("I get {string} oauth keys for application")
    public void get_oauth_keys_for_application(String type) throws Exception {
        String applicationUUID = sharedContext.getApplicationUUID();
        String keyType = (type.equals("production")) ? "production" : "sandbox";
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doGet(Utils.getOauthKeysURL(applicationUUID),
                headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setOauthKeyUUID(Utils.extractOAuthMappingID(sharedContext.getResponseBody(), sharedContext.getKeyMappingID(keyType)));
        Thread.sleep(3000);
    }

    @When("make the Access Token Generation request for {string}")
    public void make_access_token_generation_request(String type) throws Exception {
        String applicationUUID = sharedContext.getApplicationUUID();
        String oauthKeyUUID = sharedContext.getOauthKeyUUID();
        String keyType = (type.equals("production")) ? "production" : "sandbox";
        logger.info("Generating keys for: " + keyType);
        String consumerKey = sharedContext.getConsumerKey(keyType);
        String consumerSecret = sharedContext.getConsumerSecret(keyType);
        
        logger.info("Application UUID: " + applicationUUID);
        logger.info("Oauth Key UUID: " + oauthKeyUUID);
        
        String payload = "{\"consumerSecret\":\"" + consumerSecret + "\",\"validityPeriod\":3600,\"revokeToken\":null," +
                "\"scopes\":[\"write:pets\",\"read:pets\",\"query:hero\"],\"additionalProperties\":{\"id_token_expiry_time\":3600," +
                "\"application_access_token_expiry_time\":3600,\"user_access_token_expiry_time\":3600,\"bypassClientCredentials\":false," +
                "\"pkceMandatory\":false,\"pkceSupportPlain\":false,\"refresh_token_expiry_time\":86400}}";

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getAccessTokenGenerationURL(applicationUUID, oauthKeyUUID),
                headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setApiAccessToken(Utils.extractKeys(sharedContext.getResponseBody(), "accessToken"));
        logger.info("Access Token: " + sharedContext.getApiAccessToken());
        Thread.sleep(3000);
    }

    @When("I make Access Token Generation request without scopes")
    public void make_access_token_generation_request_without_scopes() throws Exception {
        String applicationUUID = sharedContext.getApplicationUUID();
        String oauthKeyUUID = sharedContext.getOauthKeyUUID();
        String keyType = "production"; //Use the same ternary logic above if both sandbox and production routes need to be tested
        String consumerKey = sharedContext.getConsumerKey(keyType);
        String consumerSecret = sharedContext.getConsumerSecret(keyType);
        
        logger.info("Application UUID: " + applicationUUID);
        logger.info("Oauth Key UUID: " + oauthKeyUUID);

        String payload = "{\"consumerSecret\":\"" + consumerSecret + "\",\"validityPeriod\":3600,\"revokeToken\":null," +
                "\"scopes\":[],\"additionalProperties\":{\"id_token_expiry_time\":3600," +
                "\"application_access_token_expiry_time\":3600,\"user_access_token_expiry_time\":3600,\"bypassClientCredentials\":false," +
                "\"pkceMandatory\":false,\"pkceSupportPlain\":false,\"refresh_token_expiry_time\":86400}}";

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(Utils.getAccessTokenGenerationURL(applicationUUID, oauthKeyUUID),
                headers, payload, Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        sharedContext.setApiAccessToken(Utils.extractKeys(sharedContext.getResponseBody(), "accessToken"));
        logger.info("Access Token without scopes: " + sharedContext.getApiAccessToken());
        Thread.sleep(3000);
    }

    @When("make the API Deployment request")
    public void make_a_api_deployment_request() throws Exception {

        // Create a MultipartEntityBuilder to build the request entity
        MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                .addPart("url", new FileBody(definitionFile))
                .addPart("apkConfiguration", new FileBody(payloadFile));

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpEntity multipartEntity = builder.build();

        HttpResponse response = sharedContext.getHttpClient().doPostWithMultipart(Utils.getAPIDeployerURL(),
                multipartEntity, headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        Thread.sleep(3000);
    }

    @When("make the API deployment request for organization {string}")
    public void makeAPIDeploymentFromOrganization(String organization) throws Exception {

        // Create a MultipartEntityBuilder to build the request entity
        MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                .addPart("definitionFile", new FileBody(definitionFile))
                .addPart("apkConfiguration", new FileBody(payloadFile));

        Map<String, String> headers = new HashMap<>();
        Object accessToken = sharedContext.getStoreValue(organization);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + accessToken);
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpEntity multipartEntity = builder.build();

        HttpResponse response = sharedContext.getHttpClient().doPostWithMultipart(Utils.getAPIDeployerURL(),
                multipartEntity, headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
        Thread.sleep(3000);
    }

    // @When("I undeploy the API whose ID is {string}")
    // public void i_undeploy_the_api_whose_id_is(String apiID) throws Exception {

    //     // Create query parameters
    //     List<NameValuePair> queryParams = new ArrayList<>();
    //     queryParams.add(new BasicNameValuePair("apiId", apiID));

    //     URI uri = new URIBuilder(Utils.getAPIUnDeployerURL()).addParameters(queryParams).build();

    //     Map<String, String> headers = new HashMap<>();
    //     headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
    //     headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

    //     HttpResponse response = sharedContext.getHttpClient().doPost(uri.toString(), headers, "",
    //             Constants.CONTENT_TYPES.APPLICATION_JSON);

    //     sharedContext.setResponse(response);
    //     sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    // }

    // @When("I undeploy the API whose ID is {string} and organization {string}")
    // public void undeployAPIByIdAndOrganization(String apiID,String organization) throws Exception {

    //     // Create query parameters
    //     List<NameValuePair> queryParams = new ArrayList<>();
    //     queryParams.add(new BasicNameValuePair("apiId", apiID));

    //     URI uri = new URIBuilder(Utils.getAPIUnDeployerURL()).addParameters(queryParams).build();

    //     Map<String, String> headers = new HashMap<>();
    //     Object header = sharedContext.getStoreValue(organization);
    //     headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + header);
    //     headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

    //     HttpResponse response = sharedContext.getHttpClient().doPost(uri.toString(), headers, "",
    //             Constants.CONTENT_TYPES.APPLICATION_JSON);

    //     sharedContext.setResponse(response);
    //     sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    // }

        @Given("a valid graphql definition file")
        public void iHaveValidGraphQLDefinition() throws Exception {
    
            // Create a MultipartEntityBuilder to build the request entity
            MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                    .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                    .addPart("file", new FileBody(definitionFile));
    
            logger.info("Definition File: "+ new FileBody(definitionFile));
    
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_APIM_HOST);
    
            HttpEntity multipartEntity = builder.build();
    
            HttpResponse response = sharedContext.getHttpClient().doPostWithMultipart(Utils.getGQLSchemaValidatorURL(),
                    multipartEntity, headers);
    
            sharedContext.setResponse(response);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setAPIDefinitionValidStatus(Utils.extractValidStatus(sharedContext.getResponseBody()));
            Thread.sleep(3000);
        }
    
        @Then("I make the import GraphQLAPI Creation request")
        public void make_import_gqlapi_creation_request() throws Exception {
    
            // Create a MultipartEntityBuilder to build the request entity
            MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                    .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                    .addPart("additionalProperties", new FileBody(payloadFile))
                    .addPart("file", new FileBody(definitionFile));
    
    
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);
    
            HttpEntity multipartEntity = builder.build();
    
            HttpResponse response = sharedContext.getHttpClient().doPostWithMultipart(Utils.getGQLImportAPIURL(),
                    multipartEntity, headers);
    
            sharedContext.setResponse(response);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setApiUUID(Utils.extractID(sharedContext.getResponseBody()));
            Thread.sleep(3000);
        }

        @Then("I update the GQL API settings")
        public void make_update_gql_request() throws Exception {
            String fileContent = new String(Files.readAllBytes(payloadFile.toPath()), StandardCharsets.UTF_8);
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);
    
            HttpResponse response = sharedContext.getHttpClient().doPut(Utils.getAPIUnDeployerURL(sharedContext.getApiUUID()), headers, fileContent ,Constants.CONTENT_TYPES.APPLICATION_JSON);
    
            sharedContext.setResponse(response);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setApiUUID(Utils.extractID(sharedContext.getResponseBody()));
            Thread.sleep(3000);
        }

        @Then("I delete the application {string} from devportal")
        public void make_application_deletion_request(String applicationName) throws Exception {
            logger.info("Fetching the applications");
        
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getDevportalAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

            List<NameValuePair> queryParams = new ArrayList<>();
            queryParams.add(new BasicNameValuePair("query", applicationName));

            URI uri = new URIBuilder(Utils.getApplicationCreateURL()).addParameters(queryParams).build();
            HttpResponse appSearchResponse = sharedContext.getHttpClient().doGet(uri.toString(), headers);
        
            sharedContext.setResponse(appSearchResponse);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setApplicationUUID(Utils.extractApplicationUUID(sharedContext.getResponseBody()));
            HttpResponse deleteResponse = sharedContext.getHttpClient().doDelete(Utils.getApplicationCreateURL() + "/" + sharedContext.getApplicationUUID(), headers);
        
            sharedContext.setResponse(deleteResponse);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            Thread.sleep(3000);
        }

        @Then("I find the apiUUID of the API created with the name {string}")
        public void find_api_uuid_using_name(String apiName) throws Exception {
            logger.info("Fetching the APIs");
        
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);
            
            HttpResponse appSearchResponse = sharedContext.getHttpClient().doGet(Utils.getAPISearchEndpoint(apiName), headers);
        
            sharedContext.setResponse(appSearchResponse);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setApiUUID(Utils.extractAPIUUID(sharedContext.getResponseBody()));
            Thread.sleep(3000);
        }

        @When("I undeploy the selected API")
        public void i_undeploy_the_api() throws Exception {
            logger.info("API UUID to be deleted: " + sharedContext.getApiUUID());
            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

            HttpResponse response = sharedContext.getHttpClient().doDelete(Utils.getAPIUnDeployerURL(sharedContext.getApiUUID()), headers);
    
            sharedContext.setResponse(response);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            Thread.sleep(3000);
        }

        @When("I create the new version {string} of the same API with default version set to {string}")
        public void create_new_version_of_the_api(String newVersion, String isDefaultVersion) throws Exception {
            String apiUUID = sharedContext.getApiUUID();

            Map<String, String> headers = new HashMap<>();
            headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getPublisherAccessToken());
            headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

                    // Create query parameters
            List<NameValuePair> queryParams = new ArrayList<>();
            queryParams.add(new BasicNameValuePair("newVersion", newVersion));
            queryParams.add(new BasicNameValuePair("defaultVersion", isDefaultVersion));
            queryParams.add(new BasicNameValuePair("apiId", apiUUID));

            URI uri = new URIBuilder(Utils.getAPINewVersionCreationURL()).addParameters(queryParams).build();
    
            HttpResponse response = sharedContext.getHttpClient().doPost(uri.toString(), headers,"",Constants.CONTENT_TYPES.APPLICATION_JSON);
    
            sharedContext.setResponse(response);
            sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
            sharedContext.setApiUUID(Utils.extractID(sharedContext.getResponseBody()));
            Thread.sleep(3000);
        }
}
