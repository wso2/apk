/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package org.wso2.apk.enforcer.server;

import com.google.protobuf.ByteString;
import io.envoyproxy.envoy.service.auth.v3.CheckRequest;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.ThreadContext;
import org.wso2.apk.enforcer.commons.constants.GraphQLConstants;
import org.wso2.apk.enforcer.api.API;
import org.wso2.apk.enforcer.api.APIFactory;
import org.wso2.apk.enforcer.api.ResponseObject;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.commons.logging.ErrorDetails;
import org.wso2.apk.enforcer.commons.logging.LoggingConstants;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.AdapterConstants;
import org.wso2.apk.enforcer.constants.HttpConstants;
import org.wso2.apk.enforcer.graphql.GraphQLPayloadUtils;
import org.wso2.apk.enforcer.util.FilterUtils;

import java.util.ArrayList;
import java.util.Map;

/**
 * This class handles the request coming via the external auth gRPC service.
 */
public class HttpRequestHandler implements RequestHandler<CheckRequest, ResponseObject> {
    private static final Logger logger = LogManager.getLogger(RequestHandler.class);

    public ResponseObject process(CheckRequest request) {

        // Handle the '/ready' request separately.
        if (APIConstants.ReadinessCheck.ENDPOINT.equals(request.getAttributes().getRequest().getHttp().getPath())) {
            return processReadyRequest();
        }

        API matchedAPI = APIFactory.getInstance().getMatchedAPI(request);
        if (matchedAPI == null) {
            ResponseObject responseObject = new ResponseObject();
            responseObject.setStatusCode(APIConstants.StatusCodes.NOTFOUND.getCode());
            responseObject.setErrorCode(APIConstants.StatusCodes.NOTFOUND.getValue());
            responseObject.setDirectResponse(true);
            responseObject.setErrorMessage(APIConstants.NOT_FOUND_MESSAGE);
            responseObject.setErrorDescription(APIConstants.NOT_FOUND_DESCRIPTION);
            return responseObject;
        }
        APIConfig api = matchedAPI.getAPIConfig();
        ThreadContext.put(APIConstants.API_UUID, api.getUuid());
        logger.debug("API {}/{} found in the cache", api.getBasePath(), api.getVersion());

        // putting API details into ThreadContext for logging purposes
        ThreadContext.push(api.getName());
        ThreadContext.push(api.getOrganizationId());
        ThreadContext.push(api.getBasePath());

        RequestContext requestContext = buildRequestContext(matchedAPI, request);
        ResponseObject responseObject = matchedAPI.process(requestContext);

        // to clear the ThreadContext's stack used for logging
        ThreadContext.removeStack();
        return responseObject;
    }

    private RequestContext buildRequestContext(API api, CheckRequest request) {
        String requestPath = request.getAttributes().getRequest().getHttp().getPath();
        String method = request.getAttributes().getRequest().getHttp().getMethod();
        String certificate = request.getAttributes().getSource().getCertificate();
        Map<String, String> headers = request.getAttributes().getRequest().getHttp().getHeadersMap();
        String pathTemplate = request.getAttributes().getContextExtensionsMap().get(APIConstants.GW_RES_PATH_PARAM);
        String routeName = request.getAttributes().getContextExtensionsMap().get(APIConstants.ROUTE_NAME_PARAM);
        String cluster = request.getAttributes().getContextExtensionsMap()
                .get(AdapterConstants.CLUSTER_HEADER_KEY);
        long requestTimeInMillis = request.getAttributes().getRequest().getTime().getSeconds() * 1000 +
                request.getAttributes().getRequest().getTime().getNanos() / 1000000;
        String requestID = request.getAttributes().getRequest().getHttp().
                getHeadersOrDefault(HttpConstants.X_REQUEST_ID_HEADER,
                        request.getAttributes().getRequest().getHttp().getId());
        String address = "";
        if (request.getAttributes().getSource().hasAddress() &&
                request.getAttributes().getSource().getAddress().hasSocketAddress()) {
            address = request.getAttributes().getSource().getAddress().getSocketAddress().getAddress();
        }
        address = FilterUtils.getClientIp(headers, address);
        String requestPayload = null;
        if (!request.getAttributes().getRequest().getHttp().getRawBody().isEmpty()) {
            ByteString byteString = request.getAttributes().getRequest().getHttp().getRawBody();
            if (byteString.isValidUtf8()) {
                requestPayload = byteString.toStringUtf8();
            }
        }
        if (!request.getAttributes().getRequest().getHttp().getBody().isEmpty()) {
            ByteString byteString = request.getAttributes().getRequest().getHttp().getBodyBytes();
            if (byteString.isValidUtf8()) {
                requestPayload = byteString.toStringUtf8();
            }
        }
        ResourceConfig resourceConfig = null;
        ArrayList<ResourceConfig> resourceConfigs = null;
        boolean isGraphQLAPI = api.getAPIConfig().getApiType().equals(APIConstants.ApiType.GRAPHQL);
        EnforcerConfig enforcerConfig = ConfigHolder.getInstance().getConfig();
        if (isGraphQLAPI && !HttpConstants.OPTIONS.equals(method)) {
            // need to decode the payload if request is graphql and a non option call.
            try {
                requestPayload = GraphQLPayloadUtils.getGQLRequestPayload(requestPayload, headers);
                resourceConfigs = GraphQLPayloadUtils.buildGQLRequestContext(api, requestPayload, routeName);
            } catch (EnforcerException exception) {
                logger.error("Error while processing the graphql api request for {}",
                        api.getAPIConfig().getName(),
                        ErrorDetails.errorLog(LoggingConstants.Severity.MINOR, 6704), exception);
                RequestContext requestContext = new RequestContext.Builder(requestPath).requestMethod(method)
                        .matchedAPI(api.getAPIConfig()).headers(headers).requestID(requestID).address(address)
                        .clusterHeader(cluster).certificate(certificate)
                        .requestTimeStamp(requestTimeInMillis).requestPayload(requestPayload).build();
                requestContext.getProperties().put(APIConstants.MessageFormat.STATUS_CODE,
                        APIConstants.StatusCodes.BAD_REQUEST_ERROR.getCode());
                requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_CODE,
                        GraphQLConstants.GRAPHQL_INVALID_QUERY);
                requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_MESSAGE,
                        GraphQLConstants.GRAPHQL_INVALID_QUERY_MESSAGE);
                requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_DESCRIPTION,
                        exception.getMessage());
                return requestContext;
            }
        } else if (!isGraphQLAPI && !enforcerConfig.getEnableGatewayClassController()) {
            resourceConfig = APIFactory.getInstance().getMatchedResource(api, pathTemplate, method);
            if (resourceConfig != null) {
                resourceConfigs = new ArrayList<>();
                resourceConfigs.add(resourceConfig);
            }
        } else if (!isGraphQLAPI && enforcerConfig.getEnableGatewayClassController()) {
            resourceConfig = APIFactory.getInstance().getMatchedResource(api,
                    request.getAttributes().getContextExtensionsMap().get(APIConstants.ROUTE_NAME_PARAM));
            if (resourceConfig != null) {
                resourceConfigs = new ArrayList<>();
                resourceConfigs.add(resourceConfig);
            }
        }
        return new RequestContext.Builder(requestPath).matchedResourceConfigs(resourceConfigs).requestMethod(method)
                .certificate(certificate).matchedAPI(api.getAPIConfig()).headers(headers).requestID(requestID)
                .address(address).clusterHeader(cluster)
                .requestTimeStamp(requestTimeInMillis).pathTemplate(pathTemplate).requestPayload(requestPayload)
                .build();
    }

    /**
     * This method handles the processing of '/ready' request sent by router.
     * 
     * @return ResponseObject with status 200
     */
    private ResponseObject processReadyRequest() {
        ResponseObject responseObject = new ResponseObject();
        responseObject.setStatusCode(APIConstants.StatusCodes.OK.getCode());
        responseObject.setRequestPath(APIConstants.ReadinessCheck.ENDPOINT);
        responseObject.setDirectResponse(true);
        return responseObject;
    }
}
