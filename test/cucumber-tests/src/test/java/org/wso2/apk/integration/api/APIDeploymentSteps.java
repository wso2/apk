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

import com.google.common.io.Resources;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.utils.URIBuilder;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.message.BasicNameValuePair;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

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
    private File apkConfFile;
    private File definitionFile;

    public APIDeploymentSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @When("I use the APK Conf file {string}")
    public void iHaveTheAPKConf(String apkConfFileName) throws IOException {

        URL url = Resources.getResource(apkConfFileName);
        apkConfFile = new File(url.getPath());
    }

    @When("the definition file {string}")
    public void iHaveTheDefinitionFile(String definitionFileName) throws IOException {

        URL url = Resources.getResource(definitionFileName);
        definitionFile = new File(url.getPath());
    }

    @When("make the API deployment request")
    public void make_a_deployment_request() throws Exception {

        // Create a MultipartEntityBuilder to build the request entity
        MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                .addPart("definitionFile", new FileBody(definitionFile))
                .addPart("apkConfiguration", new FileBody(apkConfFile));

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getAccessToken());
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
                .addPart("apkConfiguration", new FileBody(apkConfFile));

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



    @When("I undeploy the API whose ID is {string}")
    public void i_undeploy_the_api_whose_id_is(String apiID) throws Exception {

        // Create query parameters
        List<NameValuePair> queryParams = new ArrayList<>();
        queryParams.add(new BasicNameValuePair("apiId", apiID));

        URI uri = new URIBuilder(Utils.getAPIUnDeployerURL()).addParameters(queryParams).build();

        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(uri.toString(), headers, "",
                Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    }

    @When("I undeploy the API whose ID is {string} and organization {string}")
    public void undeployAPIByIdAndOrganization(String apiID,String organization) throws Exception {

        // Create query parameters
        List<NameValuePair> queryParams = new ArrayList<>();
        queryParams.add(new BasicNameValuePair("apiId", apiID));

        URI uri = new URIBuilder(Utils.getAPIUnDeployerURL()).addParameters(queryParams).build();

        Map<String, String> headers = new HashMap<>();
        Object header = sharedContext.getStoreValue(organization);
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + header);
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doPost(uri.toString(), headers, "",
                Constants.CONTENT_TYPES.APPLICATION_JSON);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    }
}
