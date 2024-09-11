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

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.protobuf.Struct;
import com.google.protobuf.Value;
import io.envoyproxy.envoy.config.core.v3.Metadata;
import io.envoyproxy.envoy.service.ext_proc.v3.BodyMutation;
import io.envoyproxy.envoy.service.ext_proc.v3.BodyResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.CommonResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.ExternalProcessorGrpc;
import io.envoyproxy.envoy.service.ext_proc.v3.HeadersResponse;
import io.envoyproxy.envoy.service.ext_proc.v3.ProcessingRequest;
import io.envoyproxy.envoy.service.ext_proc.v3.ProcessingResponse;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.grpc.client.RatelimitClient;
import java.util.ArrayList;
import java.util.List;
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
    RatelimitClient ratelimitClient = new RatelimitClient();
    @Override
    public StreamObserver<ProcessingRequest> process(
            final StreamObserver<ProcessingResponse> responseObserver) {
        FilterMetadata filterMetadata = new FilterMetadata();
        System.out.println("process ....");
        return new StreamObserver<ProcessingRequest>() {

            @Override
            public void onNext(ProcessingRequest request) {
                System.out.println("on next ....");
                ProcessingRequest.RequestCase r = request.getRequestCase();
                System.out.println("case: " + r.name());
                switch (r) {
                    case REQUEST_HEADERS:
                        if (!request.getAttributesMap().isEmpty() && request.getAttributesMap().get("envoy.filters.http.ext_proc") != null && request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata") != null){
                            Value value = request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata");
                            FilterMetadata metadata = convertStringToFilterMetadata(value.getStringValue());
                            System.out.println("Metadata generated: "+ metadata);
                            filterMetadata.backendBasedAIRatelimitDescriptorValue = metadata.backendBasedAIRatelimitDescriptorValue;
                            filterMetadata.enableBackendBasedAIRatelimit = metadata.enableBackendBasedAIRatelimit;
                            filterMetadata.enableSubscriptionBasedAIRatelimit = metadata.enableSubscriptionBasedAIRatelimit;
                        }
                        responseObserver.onNext(ProcessingResponse.newBuilder().build());
                    case RESPONSE_BODY:
                        if (!request.getAttributesMap().isEmpty() && request.getAttributesMap().get("envoy.filters.http.ext_proc") != null && request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata") != null){
                            Value value = request.getAttributesMap().get("envoy.filters.http.ext_proc").getFieldsMap().get("xds.route_metadata");
                            FilterMetadata metadata = convertStringToFilterMetadata(value.getStringValue());
                            System.out.println("Metadata generated: "+ metadata);
                            filterMetadata.backendBasedAIRatelimitDescriptorValue = metadata.backendBasedAIRatelimitDescriptorValue;
                            filterMetadata.enableBackendBasedAIRatelimit = metadata.enableBackendBasedAIRatelimit;
                            filterMetadata.enableSubscriptionBasedAIRatelimit = metadata.enableSubscriptionBasedAIRatelimit;
                        }
                        System.out.println("In the response flow metadata descirtor:" + filterMetadata.backendBasedAIRatelimitDescriptorValue);
                        if (request.hasResponseBody()) {
                            String body = request.getResponseBody().getBody().toStringUtf8();
                            Struct filterMetadataFromAuthZ = request.getMetadataContext().getFilterMetadataOrDefault("envoy.filters.http.ext_authz", null);
                            if (filterMetadataFromAuthZ != null) {
                                String extractTokenFrom = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_EXTRACT_TOKEN_FROM).getStringValue();
                                System.out.println("Extract Token From: " + extractTokenFrom);

                                String promptTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_PROMPT_TOKEN_ID).getStringValue();
                                System.out.println("Prompt Token ID: " + promptTokenID);

                                String completionTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_COMPLETION_TOKEN_ID).getStringValue();
                                System.out.println("Completion Token ID: " + completionTokenID);

                                String totalTokenID = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_TOTAL_TOKEN_ID).getStringValue();
                                System.out.println("Total Token ID: " + totalTokenID);

                                Usage usage = extractUsageFromBody(body, completionTokenID, promptTokenID, totalTokenID);
                                if (usage == null) {
                                    logger.error("Usage details not found..");
                                    System.out.println("Usage details not found..");
                                    responseObserver.onCompleted();
                                    return;
                                }
                                System.out.println("body: " +request.getResponseBody().getBody().toStringUtf8());
                                List<RatelimitClient.KeyValueHitsAddend> configs = new ArrayList<>();
                                if (filterMetadata.enableBackendBasedAIRatelimit) {
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_REQUEST_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getPrompt_tokens()));
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_RESPONSE_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getCompletion_tokens()));
                                    configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_TOTAL_TOKEN_COUNT, filterMetadata.backendBasedAIRatelimitDescriptorValue, usage.getTotal_tokens()));
                                }
                                if (filterMetadata.enableSubscriptionBasedAIRatelimit) {
                                    if (request.hasMetadataContext()) {
                                        if (filterMetadataFromAuthZ != null) {
                                            String orgAndAIRLPolicyValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_ORGANIZATION_AND_AIRL_POLICY).getStringValue();
                                            String aiRLSubsValue = filterMetadataFromAuthZ.getFieldsMap().get(DYNAMIC_METADATA_KEY_FOR_SUBSCRIPTION).getStringValue();
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_REQUEST_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getPrompt_tokens())));
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_RESPONSE_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getCompletion_tokens())));
                                            configs.add(new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_SUBSCRIPTION_BASED_AI_TOTAL_TOKEN_COUNT, orgAndAIRLPolicyValue, new RatelimitClient.KeyValueHitsAddend(DESCRIPTOR_KEY_FOR_AI_SUBSCRIPTION, aiRLSubsValue, usage.getTotal_tokens())));
                                        }
                                    }
                                }
                                ratelimitClient.shouldRatelimit(configs);
                            }
                            responseObserver.onCompleted();
                        } else {
                            System.out.println("Request does not have response body");
                            responseObserver.onCompleted();
                        }

                }
            }

            @Override
            public void onError(Throwable err) {
                System.out.println("on error ...."+ err.getLocalizedMessage() + " " + err.getMessage() + " " + err.toString()+ " ****");
            }

            @Override
            public void onCompleted() {
                System.out.println("on completed ....");
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
                                .setBodyMutation(BodyMutation.newBuilder().build())
                                .build())
                .build();
    }

    // The FilterMetadata class as per your request
    private static class FilterMetadata {
        boolean enableSubscriptionBasedAIRatelimit;
        boolean enableBackendBasedAIRatelimit;
        String backendBasedAIRatelimitDescriptorValue;
        @Override
        public String toString() {
            return "FilterMetadata{" +
                    "enableSubscriptionBasedAIRatelimit=" + enableSubscriptionBasedAIRatelimit +
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
        String enableSubscriptionPattern = "key: \"EnableSubscriptionBasedAIRatelimit\".*?string_value: \"(.*?)\"";

        // Extract and assign to the FilterMetadata object
        metadata.backendBasedAIRatelimitDescriptorValue = extractValue(input, backendValuePattern);
        metadata.enableBackendBasedAIRatelimit = Boolean.parseBoolean(extractValue(input, enableBackendPattern));
        metadata.enableSubscriptionBasedAIRatelimit = Boolean.parseBoolean(extractValue(input, enableSubscriptionPattern));

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


    public static void main(String[] args) {
        String input = "filter_metadata { key: \"envoy.filters.http.ext_proc\" value { fields { key: \"BackendBasedAIRatelimitDescriptorValue\" value { string_value: \"default-apk-backend-ratelimit-xxx-apk-backend-93fa16a00dbec9ff438aa40e6358d91dcdc22f48-api\" } } fields { key: \"EnableBackendBasedAIRatelimit\" value { string_value: \"true\" } } fields { key: \"EnableSubscriptionBasedAIRatelimit\" value { string_value: \"false\" } } } }";
        input = "{ \"choices\":[ { \"content_filter_results\":{ \"hate\":{ \"filtered\":false, \"severity\":\"safe\" }, \"self_harm\":{ \"filtered\":false, \"severity\":\"safe\" }, \"sexual\":{ \"filtered\":false, \"severity\":\"safe\" }, \"violence\":{ \"filtered\":false, \"severity\":\"safe\" } }, \"finish_reason\":\"stop\", \"index\":0, \"logprobs\":null, \"message\":{ \"content\":\"Arr, matey! Ye be askin' a great question. Care for a parrot be a task that requires careful attention and love. Here be some tips to keep yer feathered friend happy and healthy: 1. Provide a proper cage: A parrot needs an adequate-sized cage to freely stretch its wings. Make sure it has enough space for perching, playing, and spreading those beautiful feathers. Also, ensure the bars are close enough together to prevent escape. 2. Nourishing grub: A parrot's diet be important. Offer a balanced diet of high-quality parrot pellets, fresh fruits, vegetables, and some seeds. Avoid avacados, chocolate, caffeine, and anythin' toxic to a bird's delicate system. 3. Fresh water: Change the water in yer parrot's bowl daily, matey. Keeps it clean and fresh. Parrots love to dunk their beaks, so ensure they have ample water for sippin' and splish-splashin'. 4. Feathered entertainment: Parrots be social creatures and need mental stimulation. Provide 'em with plenty of toys, such as ropes, bells, and puzzle toys, to keep 'em entertained. Rotate the toys frequently to avoid boredom. 5. Avast, matey! Give 'em attention: Parrots be fond of interaction with their human companions. Spend time talkin' to 'em, singin' shanties, and makin' 'em feel loved. They may even learn a few words or phrases! 6. Exercise be important: Encourage yer parrot to exercise its wings, me hearty. Free-flyin' in a safe, enclosed area be ideal, or let 'em out of their cage for supervised playtime. 7. Regular vet visits: Aye, take yer parrot to a qualified avian vet for regular check-ups. They'll make sure yer feathered friend's health be in shipshape and suggest any necessary vaccinations or treatments. 8. Aye, watch for signs of illness: Keep a keen eye on yer parrot for any signs of illness, such as changes in appetite, behavior, or feather condition. If ye spot any concerns, consult an avian vet right quick. Remember, matey, each parrot be unique, so get to know yer bird and pay attention to its specific needs. Aye, with proper care, ye and yer parrot will forge a bond that be stronger than the mightiest of pirate ships. Fair winds and happy parrot keepin'!\", \"role\":\"assistant\" } } ], \"created\":1724232516, \"id\":\"chatcmpl-9ybvo8dQte9Hb0IkD2NCaOME9Q1LH\", \"model\":\"gpt-35-turbo\", \"object\":\"chat.completion\", \"prompt_filter_results\":[ { \"prompt_index\":0, \"content_filter_results\":{ \"hate\":{ \"filtered\":false, \"severity\":\"safe\" }, \"self_harm\":{ \"filtered\":false, \"severity\":\"safe\" }, \"sexual\":{ \"filtered\":false, \"severity\":\"safe\" }, \"violence\":{ \"filtered\":false, \"severity\":\"safe\" } } } ], \"system_fingerprint\":null, \"usage\":{ \"completion_tokens\":514, \"prompt_tokens\":33, \"total_tokens\":547 } }";
        Usage usage = extractUsageFromBody(input, "usage.completion_tokens", "usage.prompt_tokens", "usage.total_tokens");
        System.out.println(usage.completion_tokens);
//        FilterMetadata metadata = convertStringToFilterMetadata(input);
//        System.out.println(metadata.backendBasedAIRatelimitDescriptorValue);  // Printing the FilterMetadata object
    }

    public static String sanitize(String input) {
        // Replace all newline characters and tabs with a space
        return input.replaceAll("[\\t\\n\\r]+", " ").trim();
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
            System.out.println("Usage extracted: "+ usage);
            return usage;

        } catch (Exception e) {
            System.out.println(String.format("Unexpected error while extracting usage from the body: %s", body) + " \n" + e);
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

}
