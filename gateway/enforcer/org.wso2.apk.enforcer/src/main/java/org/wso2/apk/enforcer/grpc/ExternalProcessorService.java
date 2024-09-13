/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.grpc;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.protobuf.Struct;
import com.google.protobuf.Value;
import io.envoyproxy.envoy.config.core.v3.HeaderValue;
import io.envoyproxy.envoy.service.ext_proc.v3.BodyMutation;
import io.envoyproxy.envoy.service.ext_proc.v3.BodyResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.CommonResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.ExternalProcessorGrpc;
import io.envoyproxy.envoy.service.ext_proc.v3.HeaderMutation;
import io.envoyproxy.envoy.service.ext_proc.v3.HeadersResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.HttpHeaders;
import io.envoyproxy.envoy.service.ext_proc.v3.ProcessingRequest;
import io.envoyproxy.envoy.service.ext_proc.v3.ProcessingResponse;
import io.grpc.stub.StreamObserver;
import org.apache.commons.compress.compressors.CompressorStreamFactory;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.grpc.client.RatelimitClient;

import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * This is the gRPC server written to match with the envoy ext-authz filter proto file. Envoy proxy call this service.
 * This is the entry point to the filter chain process for a request.
 */
public class ExternalProcessorService extends ExternalProcessorGrpc.ExternalProcessorImplBase {
    private static final Logger logger = LogManager.getLogger(ExternalProcessorService.class);
    private static final String DESCRIPTOR_KEY_FOR_AI_REQUEST_TOKEN_COUNT  = "airequesttokencount";
    private static final String DESCRIPTOR_KEY_FOR_AI_RESPONSE_TOKEN_COUNT = "airesponsetokencount";
    private static final String DESCRIPTOR_KEY_FOR_AI_TOTAL_TOKEN_COUNT    = "aitotaltokencount";
    private static final String DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_REQUEST_TOKEN_COUNT  = "airequesttokencountsubs";
    private static final String DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_RESPONSE_TOKEN_COUNT = "airesponsetokencountsubs";
    private static final String DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_TOTAL_TOKEN_COUNT    = "aitotaltokencountsubs";
    private static final String DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION = "subscription";
    private static final String DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY = "ratelimit:organization-and-rlpolicy";
    private static final String DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION = "ratelimit:subscription";
    private static final String DYNAMIC_METADATA_KEY_FOR_EXTRACT_TOKEN_FROM = "aitoken:extracttokenfrom";
    private static final String DYNAMIC_METADATA_KEY_FOR_PROMPT_TOKEN_ID = "aitoken:prompttokenid";
    private static final String DYNAMIC_METADATA_KEY_FOR_COMPLETION_TOKEN_ID = "aitoken:completiontokenid";
    private static final String DYNAMIC_METADATA_KEY_FOR_TOTAL_TOKEN_ID = "aitoken:totaltokenid";
    private final ExecutorService executorService = Executors.newFixedThreadPool(10);;
    RatelimitClient ratelimitClient = new RatelimitClient();
    @Override
    public StreamObserver<ProcessingRequest> process(
            final StreamObserver<ProcessingResponse> responseObserver) {
        FilterMetadata filterMetadata = new FilterMetadata();
        return new StreamObserver<ProcessingRequest>() {

            @Override
            public void onNext(ProcessingRequest request) {
                ProcessingRequest.RequestCase r = request.getRequestCase();
                switch (r) {
                    case RESPONSE_HEADERS:
                        if (!request.getAttributesMap().isEmpty() && request.getAttributesMap().get("envoy.filters.http.ext_proc") != null && request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata") != null){
                            Value value = request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata");
                            FilterMetadata metadata = convertStringToFilterMetadata(value.getStringValue());
                            filterMetadata.backendBasedAIRatelimitDescriptorValue = metadata.backendBasedAIRatelimitDescriptorValue;
                            filterMetadata.enableBackendBasedAIRatelimit = metadata.enableBackendBasedAIRatelimit;
                        }
                        executorService.submit(() -> {
                            Struct filterMetadataFromAuthZ = request.getMetadataContext().getFilterMetadataOrDefault("envoy.filters.http.ext_authz", null);
                            if (filterMetadataFromAuthZ != null) {
                                String extractTokenFrom = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_EXTRACT_TOKEN_FROM).getStringValue();
                                String promptTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_PROMPT_TOKEN_ID).getStringValue();
                                String completionTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_COMPLETION_TOKEN_ID).getStringValue();
                                String totalTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_TOTAL_TOKEN_ID).getStringValue();

                                Usage usage = extractUsageFromHeaders(request.getResponseHeaders(), completionTokenID, promptTokenID, totalTokenID);
                                if (usage == null) {
                                    logger.error("Usage details not found..");
                                    responseObserver.onCompleted();
                                    return;
                                }
                                List<RatelimitClient.KeyValueHitsAddend> configs = new ArrayList<>();
                                if (filterMetadata.enableBackendBasedAIRatelimit) {
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_REQUEST_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getPrompt_tokens() - 1));
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_RESPONSE_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getCompletion_tokens() - 1));
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_TOTAL_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getTotal_tokens() - 1));
                                }
                                if (request.hasMetadataContext()) {
                                    if (filterMetadataFromAuthZ != null) {
                                        if (filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY) != null && filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION) != null) {
                                            String orgAndAIRLPolicyValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY).getStringValue();
                                            String aiRLSubsValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION).getStringValue();
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_REQUEST_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getPrompt_tokens() - 1)));
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_RESPONSE_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getCompletion_tokens() - 1)));
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_TOTAL_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getTotal_tokens() - 1)));
                                        }
                                    }
                                }
                                ratelimitClient.shouldRatelimit(configs);
                            }
                        });
                        responseObserver.onCompleted();
                    case RESPONSE_BODY:
                        if (!request.getAttributesMap().isEmpty() && request.getAttributesMap().get("envoy.filters.http.ext_proc") != null && request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata") != null){
                            Value value = request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata");
                            FilterMetadata metadata = convertStringToFilterMetadata(value.getStringValue());
                            filterMetadata.backendBasedAIRatelimitDescriptorValue = metadata.backendBasedAIRatelimitDescriptorValue;
                            filterMetadata.enableBackendBasedAIRatelimit = metadata.enableBackendBasedAIRatelimit;
                        }
                        if (request.hasResponseBody()) {
                            final byte[] bodyFromResponse = request.getResponseBody().getBody().toByteArray();
                            executorService.submit(() -> {
                                String body;
                                try {
                                    body = decompress(bodyFromResponse);
                                } catch (Exception e) {
                                    throw new RuntimeException(e);
                                }

                                Struct filterMetadataFromAuthZ = request.getMetadataContext().getFilterMetadataOrDefault("envoy.filters.http.ext_authz", null);
                                if (filterMetadataFromAuthZ != null) {
                                    String extractTokenFrom = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_EXTRACT_TOKEN_FROM).getStringValue();
                                    String promptTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_PROMPT_TOKEN_ID).getStringValue();
                                    String completionTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_COMPLETION_TOKEN_ID).getStringValue();
                                    String totalTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_TOTAL_TOKEN_ID).getStringValue();

                                    Usage usage = extractUsageFromBody(body, completionTokenID, promptTokenID, totalTokenID);
                                    if (usage == null) {
                                        logger.error("Usage details not found..");
                                        responseObserver.onCompleted();
                                        return;
                                    }
                                    List<RatelimitClient.KeyValueHitsAddend> configs = new ArrayList<>();
                                    if (filterMetadata.enableBackendBasedAIRatelimit) {
                                        configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_REQUEST_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getPrompt_tokens() - 1));
                                        configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_RESPONSE_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getCompletion_tokens() - 1));
                                        configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_TOTAL_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getTotal_tokens() - 1));
                                    }
                                    if (request.hasMetadataContext()) {
                                        if (filterMetadataFromAuthZ != null) {
                                            if (filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY) != null && filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION) != null) {
                                                String orgAndAIRLPolicyValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY).getStringValue();
                                                String aiRLSubsValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION).getStringValue();
                                                configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_REQUEST_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getPrompt_tokens() - 1)));
                                                configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_RESPONSE_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getCompletion_tokens() - 1)));
                                                configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_TOTAL_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getTotal_tokens() - 1)));
                                            }
                                        }
                                    }
                                    ratelimitClient.shouldRatelimit(configs);
                                }
                            });
                            responseObserver.onCompleted();
                        } else {
                            responseObserver.onCompleted();
                        }

                }
            }

            @Override
            public void onError(Throwable err) {
                logger.error("Error initiated from envoy in the external processing session. Error: " + err);
            }

            @Override
            public void onCompleted() {
                responseObserver.onCompleted();
            }
        };
    }

    protected BodyResponse prepareBodyResponse() {
        return BodyResponse.newBuilder()
                .setResponse(
                        CommonResponse.newBuilder()
                                .setStatus(CommonResponse.ResponseStatus.CONTINUE)
                                .setBodyMutation(BodyMutation.newBuilder().build())
                                .build())
                .build();
    }

    protected HeadersResponse prepareHeadersResponse() {
        return HeadersResponse.newBuilder()
                .setResponse(
                        CommonResponse.newBuilder()
                                .setStatus(CommonResponse.ResponseStatus.CONTINUE)
                                .setHeaderMutation(
                                        HeaderMutation.newBuilder()
                                                .build())
                                .setBodyMutation(BodyMutation.newBuilder().build())
                                .build())
                .build();
    }

    // The FilterMetadata class as per your request
    private static class FilterMetadata {
        boolean enableBackendBasedAIRatelimit;
        String backendBasedAIRatelimitDescriptorValue;
        @Override
        public String toString() {
            return "FilterMetadata{" +
                    ", enableBackendBasedAIRatelimit=" + enableBackendBasedAIRatelimit +
                    ", backendBasedAIRatelimitDescriptorValue='" + backendBasedAIRatelimitDescriptorValue + '\'' +
                    '}';
        }
    }

    // Method to parse the string and create FilterMetadata object
    public static FilterMetadata convertStringToFilterMetadata(String input) {
        FilterMetadata metadata = new FilterMetadata();
        // Regex patterns to extract specific fields
        String backendValuePattern = "key: \"BackendBasedAIRatelimitDescriptorValue\".*?string_value: \"(.*?)\"";
        String enableBackendPattern = "key: \"EnableBackendBasedAIRatelimit\".*?string_value: \"(.*?)\"";

        // Extract and assign to the FilterMetadata object
        metadata.backendBasedAIRatelimitDescriptorValue = extractValue(input, backendValuePattern);
        metadata.enableBackendBasedAIRatelimit = Boolean.parseBoolean(extractValue(input, enableBackendPattern));

        return metadata;
    }

    // Helper method to extract value based on a regex pattern
    private static String extractValue(String input, String pattern) {
        Pattern p = Pattern.compile(pattern);
        Matcher m = p.matcher(input);
        if (m.find()) {
            return m.group(1);
        }
        return null;
    }

    public static String sanitize(String input) {
        // Replace all newline characters and tabs with a space
        return input.replaceAll("[\\t\\n\\r]+", " ").trim();
    }

    private static Usage extractUsageFromHeaders(HttpHeaders headers, String completionTokenPath, String promptTokenPath, String totalTokenPath) {
        try {
            Usage usage = new Usage();
            for (HeaderValue headerValue : headers.getHeaders().getHeadersList()) {
                if (headerValue.getKey().equals(completionTokenPath)) {
                    usage.completion_tokens = Integer.parseInt(headerValue.getValue());
                }
                if (headerValue.getKey().equals(promptTokenPath)) {
                    usage.prompt_tokens = Integer.parseInt(headerValue.getValue());
                }
                if (headerValue.getKey().equals(totalTokenPath)) {
                    usage.total_tokens = Integer.parseInt(headerValue.getValue());
                }
            }
            return usage;
        } catch (Exception e) {
            logger.error("Error occured while getting yusage info from headers" + e);
            return null;
        }
    }

    private static Usage extractUsageFromBody(String body, String completionTokenPath, String promptTokenPath, String totalTokenPath) {
        body = sanitize(body);
        ObjectMapper mapper = new ObjectMapper();
        try {
            Usage usage = new Usage();
            // Parse the JSON string
            JsonNode rootNode = mapper.readTree(body);
            // Extract prompt token count
            String[] keysForPromtTokens = promptTokenPath.split("\\.");
            JsonNode currentNodeForPromtToken = null;
            if (rootNode.has(keysForPromtTokens[0])) {
                currentNodeForPromtToken = rootNode.get(keysForPromtTokens[0]);
            } else {
                return null;
            }
            for (int i = 1; i < keysForPromtTokens.length; i++) {
                if (currentNodeForPromtToken.has(keysForPromtTokens[i])) {
                    currentNodeForPromtToken = currentNodeForPromtToken.get(keysForPromtTokens[i]);
                } else {
                    return null;
                }
            }
            usage.setPrompt_tokens(currentNodeForPromtToken.asInt());

            // Extract completion token count
            String[] keysForCompletionTokens = completionTokenPath.split("\\.");
            JsonNode currentNodeForCompletionToken = null;
            if (rootNode.has(keysForCompletionTokens[0])) {
                currentNodeForCompletionToken = rootNode.get(keysForCompletionTokens[0]);
            } else {
                return null;
            }
            for (int i = 1; i < keysForCompletionTokens.length; i++) {
                if (currentNodeForCompletionToken.has(keysForCompletionTokens[i])) {
                    currentNodeForCompletionToken = currentNodeForCompletionToken.get(keysForCompletionTokens[i]);
                } else {
                    return null;
                }
            }
            usage.setCompletion_tokens(currentNodeForCompletionToken.asInt());

            // Extract total token count
            String[] keysForTotalTokens = totalTokenPath.split("\\.");
            JsonNode currentNodeForTotalToken = null;
            if (rootNode.has(keysForTotalTokens[0])) {
                currentNodeForTotalToken = rootNode.get(keysForTotalTokens[0]);
            } else {
                return null;
            }
            for (int i = 1; i < keysForTotalTokens.length; i++) {
                if (currentNodeForTotalToken.has(keysForTotalTokens[i])) {
                    currentNodeForTotalToken = currentNodeForTotalToken.get(keysForTotalTokens[i]);
                } else {
                    return null;
                }
            }
            usage.setTotal_tokens(currentNodeForTotalToken.asInt());
            return usage;

        } catch (Exception e) {
            logger.error(String.format("Unexpected error while extracting usage from the body: %s", body) + " \n" + e);
            return null;
        }
    }

    public static class Usage {
        private int completion_tokens;
        private int prompt_tokens;
        private int total_tokens;

        // Getters and Setters
        public int getCompletion_tokens() {
            return completion_tokens;
        }

        public void setCompletion_tokens(int completion_tokens) {
            this.completion_tokens = completion_tokens;
        }

        public int getPrompt_tokens() {
            return prompt_tokens;
        }

        public void setPrompt_tokens(int prompt_tokens) {
            this.prompt_tokens = prompt_tokens;
        }

        public int getTotal_tokens() {
            return total_tokens;
        }

        public void setTotal_tokens(int total_tokens) {
            this.total_tokens = total_tokens;
        }

        @Override
        public String toString() {
            return String.format("%s_%s_%s", prompt_tokens, completion_tokens, total_tokens);
        }
    }

    public static String decompress(byte[] compressed) throws Exception {
        String body = new String(compressed, StandardCharsets.UTF_8);
        if (isValidJson(body)) {
            return body;
        }
        try (InputStream is = new CompressorStreamFactory()
                .createCompressorInputStream(new ByteArrayInputStream(compressed));
             BufferedReader reader = new BufferedReader(new InputStreamReader(is, StandardCharsets.UTF_8))) {
            StringBuilder outStr = new StringBuilder();
            String line;
            while ((line = reader.readLine()) != null) {
                outStr.append(line);
            }
            if (isValidJson(outStr.toString())) {
                return outStr.toString();
            } else {
                throw new RuntimeException("Could not decompress response body");
            }
        }
    }

    public static boolean isValidJson(String json) {
        try {
            ObjectMapper mapper = new ObjectMapper();
            mapper.readTree(json);
            return true;
        } catch (Exception e) {
            return false;
        }
    }
}
