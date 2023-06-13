/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.config.definitions;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import io.swagger.v3.core.util.Json;
import io.swagger.v3.core.util.Yaml;
import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.Operation;
import io.swagger.v3.oas.models.PathItem;
import io.swagger.v3.oas.models.Paths;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.media.Content;
import io.swagger.v3.oas.models.media.MediaType;
import io.swagger.v3.oas.models.media.Schema;
import io.swagger.v3.oas.models.parameters.Parameter;
import io.swagger.v3.oas.models.parameters.RequestBody;
import io.swagger.v3.oas.models.responses.ApiResponse;
import io.swagger.v3.oas.models.responses.ApiResponses;
import io.swagger.v3.oas.models.security.OAuthFlow;
import io.swagger.v3.oas.models.security.OAuthFlows;
import io.swagger.v3.oas.models.security.Scopes;
import io.swagger.v3.oas.models.security.SecurityRequirement;
import io.swagger.v3.oas.models.security.SecurityScheme;
import io.swagger.v3.oas.models.servers.Server;
import io.swagger.v3.parser.OpenAPIV3Parser;
import io.swagger.v3.parser.core.models.ParseOptions;
import io.swagger.v3.parser.core.models.SwaggerParseResult;
import io.swagger.v3.parser.util.DeserializationUtils;
import org.apache.commons.collections.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.config.APIConstants;
import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.api.APIDefinitionValidationResponse;
import org.wso2.apk.config.api.APIManagementException;
import org.wso2.apk.config.api.ErrorItem;
import org.wso2.apk.config.api.ExceptionCodes;
import org.wso2.apk.config.model.API;
import org.wso2.apk.config.model.SwaggerData;
import org.wso2.apk.config.model.URITemplate;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Models API definition using OAS (OpenAPI 3.0) parser
 */
public class OAS3Parser extends APIDefinition {

    private static final Log log = LogFactory.getLog(OAS3Parser.class);
    static final String OPENAPI_SECURITY_SCHEMA_KEY = "default";
    static final String OPENAPI_DEFAULT_AUTHORIZATION_URL = "https://test.com";
    private List<String> otherSchemes;

    private List<String> getOtherSchemes() {

        return otherSchemes;
    }

    private void setOtherSchemes(List<String> otherSchemes) {

        this.otherSchemes = otherSchemes;
    }

    /**
     * This method returns URI templates according to the given swagger file
     *
     * @param openAPI OpenAPI Model.
     * @return URI Templates
     * @throws APIManagementException
     */
    private Set<URITemplate> getURITemplates(OpenAPI openAPI) throws APIManagementException {

        Set<URITemplate> urlTemplates = new LinkedHashSet<>();
        String[] scopes = getScopes(openAPI);

        for (String pathKey : openAPI.getPaths().keySet()) {
            PathItem pathItem = openAPI.getPaths().get(pathKey);
            for (Map.Entry<PathItem.HttpMethod, Operation> entry : pathItem.readOperationsMap().entrySet()) {
                Operation operation = entry.getValue();
                URITemplate template = new URITemplate();
                if (APIConstants.SUPPORTED_METHODS.contains(entry.getKey().name().toLowerCase())) {
                    template.setHTTPVerb(entry.getKey().name().toUpperCase());
                    template.setUriTemplate(pathKey);
                    List<String> opScopes = getScopeOfOperations(OPENAPI_SECURITY_SCHEMA_KEY, operation);
                    if (!opScopes.isEmpty()) {
                        template = OASParserUtil.setScopesToTemplate(template, opScopes, scopes);

                    } else if (!getScopeOfOperations("OAuth2Security", operation).isEmpty()) {
                        opScopes = getScopeOfOperations("OAuth2Security", operation);
                        template = OASParserUtil.setScopesToTemplate(template, opScopes, scopes);
                    }
                    if (operation.getServers() != null) {
                        template.setEndpoint(operation.getServers().get(0).getUrl());
                    }
                    urlTemplates.add(template);
                }
            }
        }
        return urlTemplates;
    }

    /**
     * This method returns URI templates according to the given swagger file
     *
     * @param resourceConfigsJSON swaggerJSON
     * @return URI Templates
     * @throws APIManagementException
     */
    @Override
    public Set<URITemplate> getURITemplates(String resourceConfigsJSON) throws APIManagementException {

        OpenAPI openAPI = getOpenAPI(resourceConfigsJSON);
        return getURITemplates(openAPI);
    }

    /**
     * This method returns the oauth scopes according to the given swagger
     *
     * @param resourceConfigsJSON resource json
     * @return scope set
     * @throws APIManagementException
     */
    @Override
    public String[] getScopes(String resourceConfigsJSON) throws APIManagementException {

        OpenAPI openAPI = getOpenAPI(resourceConfigsJSON);
        return getScopes(openAPI);
    }

    private String[] getScopes(OpenAPI openAPI) {

        Map<String, SecurityScheme> securitySchemes;
        SecurityScheme securityScheme;
        OAuthFlows oAuthFlows;
        OAuthFlow oAuthFlow;
        Scopes scopes;
        Set<String> scopeSet = new HashSet<>();
        if (openAPI.getComponents() != null && (securitySchemes = openAPI.getComponents().getSecuritySchemes()) != null) {
            if ((securityScheme = securitySchemes.get(OPENAPI_SECURITY_SCHEMA_KEY)) != null && (oAuthFlows =
                    securityScheme.getFlows()) != null && (oAuthFlow = oAuthFlows.getImplicit()) != null && (scopes =
                    oAuthFlow.getScopes()) != null) {
                for (Map.Entry<String, String> entry : scopes.entrySet()) {
                    scopeSet.add(entry.getKey());
                }
            } else if ((securityScheme = securitySchemes.get("OAuth2Security")) != null && (oAuthFlows =
                    securityScheme.getFlows()) != null && (oAuthFlow = oAuthFlows.getPassword()) != null && (scopes =
                    oAuthFlow.getScopes()) != null) {
                for (Map.Entry<String, String> entry : scopes.entrySet()) {
                    scopeSet.add(entry.getKey());
                }
            }
        }
        return OASParserUtil.sortScopes(scopeSet);
    }

    @Override
    public String generateAPIDefinition(API api) throws APIManagementException {

        SwaggerData swaggerData = new SwaggerData(api);
        return generateAPIDefinition(swaggerData);
    }

    /**
     * This method generates API definition to the given api
     *
     * @param swaggerData api
     * @return API definition in string format
     * @throws APIManagementException
     */
    private String generateAPIDefinition(SwaggerData swaggerData) {

        OpenAPI openAPI = new OpenAPI();

        // create path if null
        if (openAPI.getPaths() == null) {
            openAPI.setPaths(new Paths());
        }

        //Create info object
        Info info = new Info();
        info.setTitle(swaggerData.getTitle());
        if (swaggerData.getDescription() != null) {
            info.setDescription(swaggerData.getDescription());
        }

        Contact contact = new Contact();
        //Create contact object and map business owner info
        if (swaggerData.getContactName() != null) {
            contact.setName(swaggerData.getContactName());
        }
        if (swaggerData.getContactEmail() != null) {
            contact.setEmail(swaggerData.getContactEmail());
        }
        if (swaggerData.getContactName() != null || swaggerData.getContactEmail() != null) {
            //put contact object to info object
            info.setContact(contact);
        }

        info.setVersion(swaggerData.getVersion());
        openAPI.setInfo(info);
        updateSwaggerSecurityDefinition(openAPI, swaggerData, OPENAPI_DEFAULT_AUTHORIZATION_URL);
        if (APIConstants.GRAPHQL_API.equals(swaggerData.getTransportType())) {
            modifyGraphQLSwagger(openAPI);
        } else {
            for (SwaggerData.Resource resource : swaggerData.getResources()) {
                addOrUpdatePathToSwagger(openAPI, resource);
            }
        }
        return Json.pretty(openAPI);
    }

    /**
     * This method generates API definition using the given api's URI templates and the swagger.
     * It will alter the provided swagger definition based on the URI templates. For example: if there is a new
     * URI template which is not included in the swagger, it will be added to the swagger as a basic resource. Any
     * additional resources inside the swagger will be removed from the swagger. Changes to scopes, throtting policies,
     * on the resource will be updated on the swagger
     *
     * @param api
     * @param swagger swagger definition
     * @return API definition in string format
     * @throws APIManagementException if error occurred when generating API Definition
     */
    @Override
    public String generateAPIDefinition(API api, String swagger) throws APIManagementException {

        OpenAPI openAPI = getOpenAPI(swagger);
        return generateAPIDefinition(api, openAPI);
    }

    @Override
    public APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition, boolean returnJsonContent) throws APIManagementException {

        return validateAPIDefinition(apiDefinition, "", returnJsonContent);
    }

    @Override
    public API getAPIFromDefinition(String content) throws APIManagementException {

        OpenAPI openAPI = getOpenAPI(content);
        List<Server> servers = openAPI.getServers();
        API api = new API();
        Info info = openAPI.getInfo();
        api.setName(info.getTitle());
        api.setVersion(info.getVersion());
        if (servers != null && !servers.isEmpty()) {
            // set 1st server as endpoint
            api.setEndpoint(servers.get(0).getUrl());
        }
        Set<URITemplate> uriTemplates = getURITemplates(openAPI);
        api.setUriTemplates(uriTemplates.toArray(new URITemplate[uriTemplates.size()]));
        return api;
    }

    /**
     * This method generates API definition using the given api's URI templates and the swagger.
     * It will alter the provided swagger definition based on the URI templates. For example: if there is a new
     * URI template which is not included in the swagger, it will be added to the swagger as a basic resource. Any
     * additional resources inside the swagger will be removed from the swagger. Changes to scopes, throtting policies,
     * on the resource will be updated on the swagger
     *
     * @param api
     * @param openAPI OpenAPI
     * @return API definition in string format
     * @throws APIManagementException if error occurred when generating API Definition
     */
    private String generateAPIDefinition(API api, OpenAPI openAPI) {

        SwaggerData swaggerData = new SwaggerData(api);
        Set<SwaggerData.Resource> copy = new HashSet<>(swaggerData.getResources());

        Iterator<Map.Entry<String, PathItem>> itr = openAPI.getPaths().entrySet().iterator();
        while (itr.hasNext()) {
            Map.Entry<String, PathItem> pathEntry = itr.next();
            String pathKey = pathEntry.getKey();
            PathItem pathItem = pathEntry.getValue();
            for (Map.Entry<PathItem.HttpMethod, Operation> entry : pathItem.readOperationsMap().entrySet()) {
                Operation operation = entry.getValue();
                boolean operationFound = false;
                for (SwaggerData.Resource resource : swaggerData.getResources()) {
                    if (pathKey.equalsIgnoreCase(resource.getPath()) && entry.getKey().name().equalsIgnoreCase(resource.getVerb())) {
                        //update operations in definition
                        operationFound = true;
                        copy.remove(resource);
                        updateOperationManagedInfo(resource, operation);
                        break;
                    }
                }
                // remove operation from definition
                if (!operationFound) {
                    pathItem.operation(entry.getKey(), null);
                }
            }
            if (pathItem.readOperations().isEmpty()) {
                itr.remove();
            }
        }
        if (APIConstants.GRAPHQL_API.equals(swaggerData.getTransportType())) {
            modifyGraphQLSwagger(openAPI);
        } else {
            //adding new operations to the definition
            for (SwaggerData.Resource resource : copy) {
                addOrUpdatePathToSwagger(openAPI, resource);
            }
        }
        updateSwaggerSecurityDefinition(openAPI, swaggerData, OPENAPI_DEFAULT_AUTHORIZATION_URL);

        openAPI.getInfo().setTitle(swaggerData.getTitle());

        if (StringUtils.isEmpty(openAPI.getInfo().getVersion())) {
            openAPI.getInfo().setVersion(swaggerData.getVersion());
        }
        if (!APIConstants.GRAPHQL_API.equals(swaggerData.getTransportType())) {
            preserveResourcePathOrderFromAPI(swaggerData, openAPI);
        }
        return Json.pretty(openAPI);
    }

    /**
     * Preserve and rearrange the OpenAPI definition according to the resource path order of the updating API payload.
     *
     * @param swaggerData Updating API swagger data
     * @param openAPI     Updated OpenAPI definition
     */
    private void preserveResourcePathOrderFromAPI(SwaggerData swaggerData, OpenAPI openAPI) {

        Set<String> orderedResourcePaths = new LinkedHashSet<>();
        Paths orderedOpenAPIPaths = new Paths();
        // Iterate the URI template order given in the updating API payload (Swagger Data) and rearrange resource paths
        // order in OpenAPI with relevance to the first matching resource path item from the swagger data path list.
        for (SwaggerData.Resource resource : swaggerData.getResources()) {
            String path = resource.getPath();
            if (!orderedResourcePaths.contains(path)) {
                orderedResourcePaths.add(path);
                // Get the resource path item for the path from existing OpenAPI
                PathItem resourcePathItem = openAPI.getPaths().get(path);
                orderedOpenAPIPaths.addPathItem(path, resourcePathItem);
            }
        }
        openAPI.setPaths(orderedOpenAPIPaths);
    }

    /**
     * Construct path parameters to the Operation
     *
     * @param operation OpenAPI operation
     * @param pathName  pathname
     */
    private void populatePathParameters(Operation operation, String pathName) {

        List<String> pathParams = getPathParamNames(pathName);
        Parameter parameter;
        if (pathParams.size() > 0) {
            for (String pathParam : pathParams) {
                parameter = new Parameter();
                parameter.setName(pathParam);
                parameter.setRequired(true);
                parameter.setIn("path");
                Schema schema = new Schema();
                schema.setType("string");
                parameter.setSchema(schema);
                operation.addParametersItem(parameter);
            }
        }
    }

    /**
     * This method validates the given OpenAPI definition by content
     *
     * @param apiDefinition     OpenAPI Definition content
     * @param host              OpenAPI Definition url
     * @param returnJsonContent whether to return the converted json form of the OpenAPI definition
     * @return APIDefinitionValidationResponse object with validation information
     */
    private APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition, String host,
                                                                  boolean returnJsonContent) throws APIManagementException {

        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        OpenAPIV3Parser openAPIV3Parser = new OpenAPIV3Parser();
        ParseOptions options = new ParseOptions();
        options.setResolve(true);
        SwaggerParseResult parseAttemptForV3 = openAPIV3Parser.readContents(apiDefinition, null, options);
        if (CollectionUtils.isNotEmpty(parseAttemptForV3.getMessages())) {
            validationResponse.setValid(false);
            for (String message : parseAttemptForV3.getMessages()) {
                OASParserUtil.addErrorToValidationResponse(validationResponse, message);
                if (message.contains(APIConstants.OPENAPI_IS_MISSING_MSG)) {
                    ErrorItem errorItem = new ErrorItem();
                    errorItem.setErrorCode(ExceptionCodes.INVALID_OAS3_FOUND.getErrorCode());
                    errorItem.setMessage(ExceptionCodes.INVALID_OAS3_FOUND.getErrorMessage());
                    errorItem.setDescription(ExceptionCodes.INVALID_OAS3_FOUND.getErrorMessage());
                    validationResponse.getErrorItems().add(errorItem);
                }
            }
        } else {
            OpenAPI openAPI = parseAttemptForV3.getOpenAPI();
            Info info = openAPI.getInfo();
            List<String> endpoints;
            String endpointWithHost = "";
            if (openAPI.getServers() == null || openAPI.getServers().isEmpty()) {
                endpoints = null;
            } else {
                endpoints = openAPI.getServers().stream().map(url -> url.getUrl()).collect(Collectors.toList());
                for (String endpoint : endpoints) {
                    if (endpoint.startsWith("/")) {
                        if (StringUtils.isEmpty(host)) {
                            endpointWithHost = "http://api.yourdomain.com" + endpoint;
                        } else {
                            endpointWithHost = host + endpoint;
                        }
                        endpoints.set(endpoints.indexOf(endpoint), endpointWithHost);
                    }
                }
            }
            String title = null;
            String context = null;
            if (!StringUtils.isBlank(info.getTitle())) {
                title = info.getTitle();
                context = info.getTitle().replaceAll("\\s", "").toLowerCase();
            }
            OASParserUtil.updateValidationResponseAsSuccess(validationResponse, apiDefinition, openAPI.getOpenapi(),
                    title, info.getVersion(), context, info.getDescription(), endpoints);
            validationResponse.setParser(this);
            if (returnJsonContent) {
                if (!apiDefinition.trim().startsWith("{")) { // not a json (it is yaml)
                    JsonNode jsonNode = DeserializationUtils.readYamlTree(apiDefinition);
                    validationResponse.setJsonContent(jsonNode.toString());
                } else {
                    validationResponse.setJsonContent(apiDefinition);
                }
            }
        }
        return validationResponse;
    }

    /**
     * Gets a list of scopes using the security requirements
     *
     * @param oauth2SchemeKey OAuth2 security element key
     * @param operation       Swagger path operation
     * @return list of scopes using the security requirements
     */
    private List<String> getScopeOfOperations(String oauth2SchemeKey, Operation operation) {

        List<SecurityRequirement> security = operation.getSecurity();
        if (security != null) {
            for (Map<String, List<String>> requirement : security) {
                if (requirement.get(oauth2SchemeKey) != null) {
                    return requirement.get(oauth2SchemeKey);
                }
            }
        }
        return getScopeOfOperationsFromExtensions(operation);
    }

    /**
     * Get scope of operation
     *
     * @param operation
     * @return
     */
    private List<String> getScopeOfOperationsFromExtensions(Operation operation) {

        Map<String, Object> extensions = operation.getExtensions();
        if (extensions != null && extensions.containsKey(APIConstants.SWAGGER_X_SCOPE)) {
            String scopeKey = (String) extensions.get(APIConstants.SWAGGER_X_SCOPE);
            return Stream.of(scopeKey.split(",")).collect(Collectors.toList());
        }
        return Collections.emptyList();
    }

    /**
     * Include Scope details to the definition
     *
     * @param openAPI     openapi definition
     * @param swaggerData Swagger related API data
     */
    private void updateSwaggerSecurityDefinition(OpenAPI openAPI, SwaggerData swaggerData, String authUrl) {

        if (openAPI.getComponents() == null) {
            openAPI.setComponents(new Components());
        }
        Map<String, SecurityScheme> securitySchemes = openAPI.getComponents().getSecuritySchemes();
        if (securitySchemes == null) {
            securitySchemes = new HashMap<>();
            openAPI.getComponents().setSecuritySchemes(securitySchemes);
        }
        SecurityScheme securityScheme = securitySchemes.get(OPENAPI_SECURITY_SCHEMA_KEY);
        if (securityScheme == null) {
            securityScheme = new SecurityScheme();
            securityScheme.setType(SecurityScheme.Type.OAUTH2);
            securitySchemes.put(OPENAPI_SECURITY_SCHEMA_KEY, securityScheme);
            List<SecurityRequirement> security = new ArrayList<SecurityRequirement>();
            SecurityRequirement secReq = new SecurityRequirement();
            secReq.addList(OPENAPI_SECURITY_SCHEMA_KEY, new ArrayList<String>());
            security.add(secReq);
            openAPI.setSecurity(security);
        }
        if (securityScheme.getFlows() == null) {
            securityScheme.setFlows(new OAuthFlows());
        }
        OAuthFlow oAuthFlow = securityScheme.getFlows().getImplicit();
        if (oAuthFlow == null) {
            oAuthFlow = new OAuthFlow();
            securityScheme.getFlows().setImplicit(oAuthFlow);
        }
        oAuthFlow.setAuthorizationUrl(authUrl);
        Scopes oas3Scopes = new Scopes();
        Set<String> scopes = swaggerData.getScopes();
        if (scopes != null && !scopes.isEmpty()) {
            for (String scope : scopes) {
                oas3Scopes.put(scope, scope);
            }
        }
        oAuthFlow.setScopes(oas3Scopes);
    }

    /**
     * Add a new path based on the provided URI template to swagger if it does not exists. If it exists,
     * adds the respective operation to the existing path
     *
     * @param openAPI  swagger object
     * @param resource API resource data
     */
    private void addOrUpdatePathToSwagger(OpenAPI openAPI, SwaggerData.Resource resource) {

        PathItem path;
        if (openAPI.getPaths() == null) {
            openAPI.setPaths(new Paths());
        }
        if (openAPI.getPaths().get(resource.getPath()) != null) {
            path = openAPI.getPaths().get(resource.getPath());
        } else {
            path = new PathItem();
        }

        Operation operation = createOperation(resource);
        PathItem.HttpMethod httpMethod = PathItem.HttpMethod.valueOf(resource.getVerb());
        path.operation(httpMethod, operation);
        openAPI.getPaths().addPathItem(resource.getPath(), path);
    }

    /**
     * Creates a new operation object using the URI template object
     *
     * @param resource API resource data
     * @return a new operation object using the URI template object
     */
    private Operation createOperation(SwaggerData.Resource resource) {

        Operation operation = new Operation();
        populatePathParameters(operation, resource.getPath());
        updateOperationManagedInfo(resource, operation);

        ApiResponses apiResponses = new ApiResponses();
        ApiResponse apiResponse = new ApiResponse();
        apiResponse.description("OK");
        apiResponses.addApiResponse(APIConstants.SWAGGER_RESPONSE_200, apiResponse);
        operation.setResponses(apiResponses);
        return operation;
    }

    /**
     * Updates managed info of a provided operation such as auth type and throttling
     *
     * @param resource  API resource data
     * @param operation swagger operation
     */
    private void updateOperationManagedInfo(SwaggerData.Resource resource, Operation operation) {

        List<SecurityRequirement> security = operation.getSecurity();
        if (security == null) {
            security = new ArrayList<>();
            operation.setSecurity(security);
        }
        for (Map<String, List<String>> requirement : security) {
            if (requirement.get(OPENAPI_SECURITY_SCHEMA_KEY) != null) {

                if (resource.getScopes() == null || resource.getScopes().length==0) {
                    requirement.put(OPENAPI_SECURITY_SCHEMA_KEY, Collections.EMPTY_LIST);
                } else {
                    requirement.put(OPENAPI_SECURITY_SCHEMA_KEY, Arrays.asList(resource.getScopes()));
                }
                return;
            }
        }
        // if oauth2SchemeKey not present, add a new
        SecurityRequirement defaultRequirement = new SecurityRequirement();
        if (resource.getScopes() != null || resource.getScopes().length==0) {
            defaultRequirement.put(OPENAPI_SECURITY_SCHEMA_KEY, Collections.EMPTY_LIST);
        } else {
            defaultRequirement.put(OPENAPI_SECURITY_SCHEMA_KEY,Arrays.asList(resource.getScopes()));
        }
        security.add(defaultRequirement);
    }

    /**
     * Get parsed OpenAPI object
     *
     * @param oasDefinition OAS definition
     * @return OpenAPI
     */
    OpenAPI getOpenAPI(String oasDefinition) {

        OpenAPIV3Parser openAPIV3Parser = new OpenAPIV3Parser();
        SwaggerParseResult parseAttemptForV3 = openAPIV3Parser.readContents(oasDefinition, null, null);
        if (CollectionUtils.isNotEmpty(parseAttemptForV3.getMessages())) {
            log.debug("Errors found when parsing OAS definition");
        }
        return parseAttemptForV3.getOpenAPI();
    }

    /**
     * Construct openAPI definition for graphQL. Add get and post operations
     *
     * @param openAPI OpenAPI
     * @return modified openAPI for GraphQL
     */
    private void modifyGraphQLSwagger(OpenAPI openAPI) {

        SwaggerData.Resource resource = new SwaggerData.Resource();
        resource.setAuthType(true);
        resource.setPolicy(APIConstants.DEFAULT_SUB_POLICY_UNLIMITED);
        resource.setPath("/");
        resource.setVerb(APIConstants.HTTP_POST);
        Operation postOperation = createOperation(resource);

        //post operation
        RequestBody requestBody = new RequestBody();
        requestBody.setDescription("Query or mutation to be passed to graphQL API");
        requestBody.setRequired(true);

        JsonObject typeOfPayload = new JsonObject();
        JsonObject payload = new JsonObject();
        typeOfPayload.addProperty(APIConstants.TYPE, APIConstants.STRING);
        payload.add(APIConstants.OperationParameter.PAYLOAD_PARAM_NAME, typeOfPayload);

        Schema postSchema = new Schema();
        postSchema.setType(APIConstants.OBJECT);
        postSchema.setProperties(new Gson().fromJson(payload, Map.class));

        MediaType mediaType = new MediaType();
        mediaType.setSchema(postSchema);

        Content content = new Content();
        content.addMediaType(APIConstants.APPLICATION_JSON_MEDIA_TYPE, mediaType);
        requestBody.setContent(content);
        postOperation.setRequestBody(requestBody);

        //add post and get operations to path /*
        PathItem pathItem = new PathItem();
        pathItem.setPost(postOperation);
        Paths paths = new Paths();
        paths.put("/", pathItem);

        openAPI.setPaths(paths);
    }

    /**
     * This method returns the boolean value which checks whether the swagger is included default security scheme or not
     *
     * @param swaggerContent resource json
     * @return boolean
     */
    private boolean isDefaultGiven(String swaggerContent) {

        OpenAPI openAPI = getOpenAPI(swaggerContent);

        Components components = openAPI.getComponents();
        if (components == null) {
            return false;
        }
        Map<String, SecurityScheme> securitySchemes = components.getSecuritySchemes();
        if (securitySchemes == null) {
            return false;
        }
        SecurityScheme checkDefault = openAPI.getComponents().getSecuritySchemes().get(OPENAPI_SECURITY_SCHEMA_KEY);
        return checkDefault != null;
    }

    /**
     * This method will inject scopes of other schemes to the swagger definition
     *
     * @param swaggerContent resource json
     * @return String
     * @throws APIManagementException
     */
    @Override
    public String processOtherSchemeScopes(String swaggerContent) throws APIManagementException {

        if (!isDefaultGiven(swaggerContent)) {
            OpenAPI openAPI = getOpenAPI(swaggerContent);
            openAPI = injectOtherScopesToDefaultScheme(openAPI);
            openAPI = injectOtherResourceScopesToDefaultScheme(openAPI);
            return Json.pretty(openAPI);
        }
        return swaggerContent;
    }

    /**
     * This method returns the oauth scopes according to the given swagger(version 3)
     *
     * @param openAPI - OpenApi object
     * @return OpenAPI
     * @throws APIManagementException
     */
    private OpenAPI injectOtherScopesToDefaultScheme(OpenAPI openAPI) throws APIManagementException {

        Map<String, SecurityScheme> securitySchemes;
        Components component = openAPI.getComponents();
        List<String> otherSetOfSchemes = new ArrayList<>();

        if (openAPI.getComponents() != null && (securitySchemes = openAPI.getComponents().getSecuritySchemes()) != null) {
            //If there is no default type schemes set a one
            SecurityScheme defaultScheme = securitySchemes.get(OPENAPI_SECURITY_SCHEMA_KEY);
            if (defaultScheme == null) {
                SecurityScheme newDefault = new SecurityScheme();
                newDefault.setType(SecurityScheme.Type.OAUTH2);
                //Populating the default security scheme with default values
                OAuthFlows newDefaultFlows = new OAuthFlows();
                OAuthFlow newDefaultFlow = new OAuthFlow();
                newDefaultFlow.setAuthorizationUrl(OPENAPI_DEFAULT_AUTHORIZATION_URL);
                Scopes newDefaultScopes = new Scopes();
                newDefaultFlow.setScopes(newDefaultScopes);
                newDefaultFlows.setImplicit(newDefaultFlow);
                newDefault.setFlows(newDefaultFlows);

                securitySchemes.put(OPENAPI_SECURITY_SCHEMA_KEY, newDefault);
            }
            for (Map.Entry<String, SecurityScheme> entry : securitySchemes.entrySet()) {
                if (!OPENAPI_SECURITY_SCHEMA_KEY.equals(entry.getKey()) && "oauth2".equals(entry.getValue().getType().toString())) {
                    otherSetOfSchemes.add(entry.getKey());
                    //Check for default one
                    SecurityScheme defaultType = securitySchemes.get(OPENAPI_SECURITY_SCHEMA_KEY);
                    OAuthFlows defaultTypeFlows = defaultType.getFlows();
                    if (defaultTypeFlows == null) {
                        defaultTypeFlows = new OAuthFlows();
                    }
                    OAuthFlow defaultTypeFlow = defaultTypeFlows.getImplicit();
                    if (defaultTypeFlow == null) {
                        defaultTypeFlow = new OAuthFlow();
                    }

                    SecurityScheme noneDefaultType = entry.getValue();
                    OAuthFlows noneDefaultTypeFlows = noneDefaultType.getFlows();
                    //Get Implicit Flows
                    OAuthFlow noneDefaultTypeFlowImplicit = noneDefaultTypeFlows.getImplicit();
                    if (noneDefaultTypeFlowImplicit != null) {
                        defaultTypeFlow = extractAndInjectScopesFromFlow(noneDefaultTypeFlowImplicit, defaultTypeFlow);
                        defaultTypeFlows.setImplicit(defaultTypeFlow);
                    }
                    //Get AuthorizationCode Flow
                    OAuthFlow noneDefaultTypeFlowAuthorizationCode = noneDefaultTypeFlows.getAuthorizationCode();
                    if (noneDefaultTypeFlowAuthorizationCode != null) {
                        defaultTypeFlow = extractAndInjectScopesFromFlow(noneDefaultTypeFlowAuthorizationCode,
                                defaultTypeFlow);
                        defaultTypeFlows.setImplicit(defaultTypeFlow);
                    }
                    //Get ClientCredentials Flow
                    OAuthFlow noneDefaultTypeFlowClientCredentials = noneDefaultTypeFlows.getClientCredentials();
                    if (noneDefaultTypeFlowClientCredentials != null) {
                        defaultTypeFlow = extractAndInjectScopesFromFlow(noneDefaultTypeFlowClientCredentials,
                                defaultTypeFlow);
                        defaultTypeFlows.setImplicit(defaultTypeFlow);
                    }
                    //Get Password Flow
                    OAuthFlow noneDefaultTypeFlowPassword = noneDefaultTypeFlows.getPassword();
                    if (noneDefaultTypeFlowPassword != null) {
                        defaultTypeFlow = extractAndInjectScopesFromFlow(noneDefaultTypeFlowPassword, defaultTypeFlow);
                        defaultTypeFlows.setImplicit(defaultTypeFlow);
                    }

                    defaultType.setFlows(defaultTypeFlows);
                }
            }
            component.setSecuritySchemes(securitySchemes);
            openAPI.setComponents(component);
        }
        setOtherSchemes(otherSetOfSchemes);
        return openAPI;
    }

    /**
     * This method returns the oauth scopes of Oauthflows according to the given swagger(version 3)
     *
     * @param noneDefaultTypeFlow , OAuthflow
     * @param defaultTypeFlow,    OAuthflow
     * @return OAuthFlow
     */
    private OAuthFlow extractAndInjectScopesFromFlow(OAuthFlow noneDefaultTypeFlow, OAuthFlow defaultTypeFlow) {

        Scopes noneDefaultFlowScopes = noneDefaultTypeFlow.getScopes();
        Scopes defaultFlowScopes = defaultTypeFlow.getScopes();
        Map<String, String> defaultScopeBindings = null;
        if (defaultFlowScopes == null) {
            defaultFlowScopes = new Scopes();
        }
        if (noneDefaultFlowScopes != null) {
            for (Map.Entry<String, String> input : noneDefaultFlowScopes.entrySet()) {
                //Inject scopes set into default scheme
                defaultFlowScopes.addString(input.getKey(), input.getValue());
            }
        }
        defaultTypeFlow.setScopes(defaultFlowScopes);
        //Check X-Scope Bindings
        Map<String, String> noneDefaultScopeBindings = null;
        Map<String, Object> defaultTypeExtension = defaultTypeFlow.getExtensions();
        if (defaultTypeExtension == null) {
            defaultTypeExtension = new HashMap<>();
        }
        if (noneDefaultTypeFlow.getExtensions() != null && (noneDefaultScopeBindings =
                (Map<String, String>) noneDefaultTypeFlow.getExtensions().get(APIConstants.SWAGGER_X_SCOPES_BINDINGS)) != null) {
            defaultScopeBindings =
                    (Map<String, String>) defaultTypeExtension.get(APIConstants.SWAGGER_X_SCOPES_BINDINGS);
            if (defaultScopeBindings == null) {
                defaultScopeBindings = new HashMap<>();
            }
            for (Map.Entry<String, String> roleInUse : noneDefaultScopeBindings.entrySet()) {
                defaultScopeBindings.put(roleInUse.getKey(), roleInUse.getValue());
            }
        }
        defaultTypeExtension.put(APIConstants.SWAGGER_X_SCOPES_BINDINGS, defaultScopeBindings);
        defaultTypeFlow.setExtensions(defaultTypeExtension);
        return defaultTypeFlow;
    }

    /**
     * This method returns URI templates according to the given swagger file(Swagger version 3)
     *
     * @param openAPI OpenAPI
     * @return OpenAPI
     * @throws APIManagementException
     */
    private OpenAPI injectOtherResourceScopesToDefaultScheme(OpenAPI openAPI) throws APIManagementException {

        List<String> schemes = getOtherSchemes();

        Paths paths = openAPI.getPaths();
        for (String pathKey : paths.keySet()) {
            PathItem pathItem = paths.get(pathKey);
            Map<PathItem.HttpMethod, Operation> operationsMap = pathItem.readOperationsMap();
            for (Map.Entry<PathItem.HttpMethod, Operation> entry : operationsMap.entrySet()) {
                SecurityRequirement updatedDefaultSecurityRequirement = new SecurityRequirement();
                PathItem.HttpMethod httpMethod = entry.getKey();
                Operation operation = entry.getValue();
                List<SecurityRequirement> securityRequirements = operation.getSecurity();
                if (securityRequirements == null) {
                    securityRequirements = new ArrayList<>();
                }
                if (APIConstants.SUPPORTED_METHODS.contains(httpMethod.name().toLowerCase())) {
                    List<String> opScopesDefault = new ArrayList<>();
                    List<String> opScopesDefaultInstance = getScopeOfOperations(OPENAPI_SECURITY_SCHEMA_KEY, operation);
                    if (opScopesDefaultInstance != null) {
                        opScopesDefault.addAll(opScopesDefaultInstance);
                    }
                    updatedDefaultSecurityRequirement.put(OPENAPI_SECURITY_SCHEMA_KEY, opScopesDefault);
                    for (Map<String, List<String>> input : securityRequirements) {
                        for (String scheme : schemes) {
                            if (!OPENAPI_SECURITY_SCHEMA_KEY.equals(scheme)) {
                                List<String> opScopesOthers = getScopeOfOperations(scheme, operation);
                                if (opScopesOthers != null) {
                                    for (String scope : opScopesOthers) {
                                        if (!opScopesDefault.contains(scope)) {
                                            opScopesDefault.add(scope);
                                        }
                                    }
                                }
                            }
                            updatedDefaultSecurityRequirement.put(OPENAPI_SECURITY_SCHEMA_KEY, opScopesDefault);
                        }
                    }
                    securityRequirements.add(updatedDefaultSecurityRequirement);
                }
                operation.setSecurity(securityRequirements);
                entry.setValue(operation);
                operationsMap.put(httpMethod, operation);
            }
            paths.put(pathKey, pathItem);
        }
        openAPI.setPaths(paths);
        return openAPI;
    }

    /**
     * Remove x-examples from all the paths from the OpenAPI definition.
     *
     * @param apiDefinition OpenAPI definition as String
     */
    public static String removeExamplesFromOpenAPI(String apiDefinition) throws APIManagementException {

        try {
            OpenAPIV3Parser openAPIV3Parser = new OpenAPIV3Parser();
            SwaggerParseResult parseAttemptForV3 = openAPIV3Parser.readContents(apiDefinition, null, null);
            if (CollectionUtils.isNotEmpty(parseAttemptForV3.getMessages())) {
                log.debug("Errors found when parsing OAS definition");
            }
            OpenAPI openAPI = parseAttemptForV3.getOpenAPI();
            for (Map.Entry<String, PathItem> entry : openAPI.getPaths().entrySet()) {
                String path = entry.getKey();
                List<Operation> operations = openAPI.getPaths().get(path).readOperations();
                for (Operation operation : operations) {
                    if (operation.getExtensions() != null && operation.getExtensions().containsKey(APIConstants.SWAGGER_X_EXAMPLES)) {
                        operation.getExtensions().remove(APIConstants.SWAGGER_X_EXAMPLES);
                    }
                }
            }
            return Yaml.pretty().writeValueAsString(openAPI);
        } catch (JsonProcessingException e) {
            throw new APIManagementException("Error while removing examples from OpenAPI definition", e,
                    ExceptionCodes.ERROR_REMOVING_EXAMPLES);
        }
    }

    @Override
    public String getType() {

        return null;
    }

    @Override
    public boolean canHandleDefinition(String definition) {

        return (StringUtils.isNotEmpty(definition) && definition.contains("openapi"));
    }
}
