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

package org.wso2.apk.integration.utils;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.HttpStatus;
import org.apache.http.entity.ContentType;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

public class Utils {

    public static String getConfigGeneratorURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_CONFIGURATOR + "apis/generate-configuration";
    }

    public static String getTokenEndpointURL() {

        return "https://" + Constants.DEFAULT_IDP_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_TOKEN_EP;
    }

    public static String getAPIDeployerURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/deploy";
    }

    public static String getAPIUnDeployerURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/undeploy";
    }

    public static String extractToken(HttpResponse response) throws IOException {

        int responseCode = response.getStatusLine().getStatusCode();

        HttpEntity entity = response.getEntity();
        Charset charset = ContentType.getOrDefault(entity).getCharset();
        if (charset == null) {
            charset = StandardCharsets.UTF_8;
        }

        BufferedReader reader = new BufferedReader(new InputStreamReader(entity.getContent(), charset));
        String inputLine;
        StringBuilder stringBuilder = new StringBuilder();

        while ((inputLine = reader.readLine()) != null) {
            stringBuilder.append(inputLine);
        }

        if (responseCode != HttpStatus.SC_OK) {
            throw new IOException("Error while accessing the Token URL. "
                    + response.getStatusLine());
        }

        JsonParser parser = new JsonParser();
        JsonObject jsonResponse = (JsonObject) parser.parse(stringBuilder.toString());
        if (jsonResponse.has("access_token")) {
            return jsonResponse.get("access_token").getAsString();
        }
        throw new IOException("Missing key [access_token] in the response from the OAuth server");
    }
}
