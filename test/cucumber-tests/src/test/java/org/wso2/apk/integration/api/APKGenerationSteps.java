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
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import java.io.File;
import java.net.URL;
import java.nio.charset.StandardCharsets;

import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

/**
 * This class contains the step definitions for APK generation.
 */
public class APKGenerationSteps {

    private final SharedContext sharedContext;
    private static final Log logger = LogFactory.getLog(BaseSteps.class);
    private File definitionFile;

    public APKGenerationSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @When("I use the definition file {string} in resources")
    public void i_use_the_definition_file_in_resources(String definitionFilePath) {

        URL url = Resources.getResource(definitionFilePath);
        definitionFile = new File(url.getPath());
    }

    @When("generate the APK conf file for a {string} API")
    public void generate_the_apk_conf_file(String apiType) throws Exception {

        // Create a MultipartEntityBuilder to build the request entity
        MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                .addTextBody("apiType", apiType)
                .addPart("definition", new FileBody(definitionFile));

        HttpEntity multipartEntity = builder.build();
                HttpResponse httpResponse = sharedContext.getHttpClient().doPostWithMultipart(Utils.getConfigGeneratorURL(),
                multipartEntity);
        sharedContext.setResponse(httpResponse);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    }

    @Then("the response body should be {string} in resources")
    public void the_response_body_should_be_in_resources(String expectedAPKConfFilePath) throws Exception {

        URL url = Resources.getResource(expectedAPKConfFilePath);
        String text = Resources.toString(url, StandardCharsets.UTF_8);
        Assert.assertEquals(sharedContext.getResponseBody(), text);
    }
}
