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
import java.util.Base64;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;

public class Utils {

    public static String getConfigGeneratorURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_CONFIGURATOR + "apis/generate-configuration";
    }

    public static String getDCREndpointURL() {

        return "https://" + Constants.DEFAULT_IDP_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DCR_EP;
    }

    public static String getTokenEndpointURL() {

        return "https://" + Constants.DEFAULT_IDP_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_TOKEN_EP;
    }

    public static String getAPIDeployerURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/deploy";
    }

    public static String getImportAPIURL() {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/import-openapi";
    }

    public static String getAPIRevisionURL(String apiUUID) {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/" + apiUUID + "/revisions";
    }

    public static String getAPIChangeLifecycleURL(String apiUUID) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/change-lifecycle?action=Publish&apiId=" + apiUUID;
    }

    public static String getApplicationCreateURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL + "applications";
    }

    public static String getGenerateKeysURL(String applicationId) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL + "applications/" + applicationId + "/generate-keys";
    }

    public static String getOauthKeysURL(String applicationId) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL + "applications/" + applicationId + "/oauth-keys";
    }

    public static String getKeyManagerURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL+ "key-managers";
    }

    public static String getSubscriptionURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL + "subscriptions";
    }

    public static String getAccessTokenGenerationURL(String applicationId, String oauthKeyId) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_DEVPORTAL + "applications/" + applicationId + "/oauth-keys/" + oauthKeyId + "/generate-token";
    }

    public static String getAPIRevisionDeploymentURL(String apiUUID, String revisionId) {

        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/" + apiUUID + "/deploy-revision?revisionId=" + revisionId;
    }

    public static String getAPIUnDeployerURL(String apiID) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/" + apiID;
    }

    public static String getGQLSchemaValidatorURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/validate-graphql-schema";
    }

    public static String getGQLImportAPIURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/import-graphql-schema";
    }

    public static String getAPISearchEndpoint(String queryValue) {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "search?query=content:" + queryValue;
    }

    public static String getAPINewVersionCreationURL() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_API_DEPLOYER + "apis/copy-api";
    }

    public static String getAPIThrottlingConfigEndpoint() {
        return "https://" + Constants.DEFAULT_API_HOST + ":" + Constants.DEFAULT_GW_PORT + "/"
                + Constants.DEFAULT_ADMINPORTAL+ "throttling/policies/advanced";
    }

    public static String extractID(String payload) throws IOException {

        JSONParser parser = new JSONParser();
        try {
            // Parse the JSON string
            JSONObject jsonObject = (JSONObject) parser.parse(payload);

            // Get the value of the "id" attribute
            String idValue = (String) jsonObject.get("id");
            return idValue;
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
    }

    public static String extractApplicationID(String payload) throws IOException {

        JSONParser parser = new JSONParser();
        try {
            // Parse the JSON string
            JSONObject jsonObject = (JSONObject) parser.parse(payload);

            // Get the value of the "applicationId" attribute
            String idValue = (String) jsonObject.get("applicationId");
            return idValue;
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
    }

    public static String extractKeyManagerID(String payload) throws IOException {

        JSONParser parser = new JSONParser();
        try {
            // Parse the JSON string
            JSONObject jsonObject = (JSONObject) parser.parse(payload);

            // Get the value of the "id" attribute
            JSONArray idValue = (JSONArray)jsonObject.get("list");
            JSONObject keyManager = (JSONObject) idValue.get(0);
            String keyManagerId = (String) keyManager.get("id");
            return keyManagerId;
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
    }

    public static String extractOAuthMappingID(String payload, String keyMappingID) throws IOException {
        JSONParser parser = new JSONParser();
        try {
            JSONObject jsonObject = (JSONObject) parser.parse(payload);
            JSONArray list = (JSONArray) jsonObject.get("list");

            for (Object obj : list) {
                JSONObject keyManager = (JSONObject) obj;
                String currentKeyMappingId = (String) keyManager.get("keyMappingId");
                if (keyMappingID.equals(currentKeyMappingId)) {
                    return currentKeyMappingId;
                }
            }
            return null;
    
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
    }

    public static String extractKeys(String payload, String key) throws IOException {

        JSONParser parser = new JSONParser();
        try {
            // Parse the JSON string
            JSONObject jsonObject = (JSONObject) parser.parse(payload);

            // Get the value of the "applicationId" attribute
            String idValue = (String) jsonObject.get(key);
            return idValue;
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
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

    public static String extractBasicToken(HttpResponse response) throws IOException {

        int responseCode = response.getStatusLine().getStatusCode();
        String clientId = null;
        String clientSecret = null;

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
        if (jsonResponse.has("clientId")) {
            clientId = jsonResponse.get("clientId").getAsString();
        }
        if (jsonResponse.has("clientSecret")) {
            clientSecret = jsonResponse.get("clientSecret").getAsString();
        }
        if (clientId != null && clientSecret != null) {
            // base64 encode the clientId and clientSecret
            return Base64.getEncoder().encodeToString((clientId + ":" + clientSecret).getBytes());

        }
        throw new IOException("Missing key [access_token] in the response from the OAuth server");
    }

    public static String resolveVariables(String input, Map<String, Object> valueStore) {
        // Define the pattern to match variables like ${variableName}
        Pattern pattern = Pattern.compile("\\$\\{([^}]*)\\}");
        Matcher matcher = pattern.matcher(input);
        StringBuffer resolvedString = new StringBuffer();

        while (matcher.find()) {
            String variableName = matcher.group(1);
            String variableValue = valueStore.get(variableName).toString();

            // Replace the variable with its value from the value store if it exists
            // Otherwise, keep the variable placeholder as is in the string
            String replacement = (variableValue != null) ? variableValue : matcher.group();
            matcher.appendReplacement(resolvedString, Matcher.quoteReplacement(replacement));
        }

        matcher.appendTail(resolvedString);
        return resolvedString.toString();
    }

    public static Boolean extractValidStatus(String payload) throws IOException {
        JSONParser parser = new JSONParser();
        try {
            // Parse the JSON string
            JSONObject jsonObject = (JSONObject) parser.parse(payload);

            // Get the value of the "isValid" attribute
            Boolean validStatus = (Boolean) jsonObject.get("isValid");
            return validStatus;
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
    }

    public static String extractApplicationUUID(String payload) throws IOException {
        JSONParser parser = new JSONParser();
        try {
            JSONObject jsonObject = (JSONObject) parser.parse(payload);
            long count = (long) jsonObject.get("count");
            if (count == 1) {
                JSONArray list = (JSONArray) jsonObject.get("list");
                JSONObject applicationObj = (JSONObject) list.get(0);
                String applicationId = (String) applicationObj.get("applicationId");
                return applicationId;
            }
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
        return null; // Return null if count is not 1
    }

    public static String extractAPIUUID(String payload) throws IOException {
        JSONParser parser = new JSONParser();
        try {
            JSONObject jsonObject = (JSONObject) parser.parse(payload);
            long count = (long) jsonObject.get("count");
            if (count == 1) {
                JSONArray list = (JSONArray) jsonObject.get("list");
                JSONObject apiObj = (JSONObject) list.get(0);
                String apiId = (String) apiObj.get("id");
                return apiId;
            }
        } catch (ParseException e) {
            throw new IOException("Error while parsing the JSON payload: " + e.getMessage());
        }
        return null; // Return null if count is not 1
    }
}
