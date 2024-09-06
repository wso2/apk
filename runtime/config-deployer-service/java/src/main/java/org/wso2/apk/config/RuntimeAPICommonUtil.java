package org.wso2.apk.config;

import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.api.APIDefinitionValidationResponse;
import org.wso2.apk.config.api.APIManagementException;
import org.wso2.apk.config.api.ExceptionCodes;
import org.wso2.apk.config.definitions.GraphQLSchemaDefinition;
import org.wso2.apk.config.definitions.OASParserUtil;
import org.wso2.apk.config.definitions.ProtoParser;
import org.wso2.apk.config.definitions.ProtoParserUtil;
import org.wso2.apk.config.model.API;
import org.wso2.apk.config.model.URITemplate;

import graphql.schema.idl.SchemaParser;
import graphql.schema.idl.TypeDefinitionRegistry;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;

public class RuntimeAPICommonUtil {

    public static String generateDefinition(API api) throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getParser(api);
        return parser.generateAPIDefinition(api);
    }

    /**
     * @param inputByteArray API definition file
     * @param apiDefinition  API definition
     * @param fileName       Filename of the definition file
     * @param returnContent  Whether to return json or not
     * @return APIDefinitionValidationResponse
     * @throws APIManagementException when file parsing fails
     */
    public static APIDefinitionValidationResponse validateOpenAPIDefinition(String type, byte[] inputByteArray,
            String apiDefinition, String fileName, boolean returnContent) throws APIManagementException {

        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        if (APIConstants.ParserType.REST.name().equals(type)) {
            if (inputByteArray != null && inputByteArray.length > 0) {
                if (fileName != null) {
                    if (fileName.endsWith(".zip")) {
                        validationResponse = OASParserUtil.extractAndValidateOpenAPIArchive(inputByteArray,
                                returnContent);
                    } else {
                        String openAPIContent = new String(inputByteArray, StandardCharsets.UTF_8);
                        validationResponse = OASParserUtil.validateAPIDefinition(openAPIContent, returnContent);
                    }
                } else {
                    String openAPIContent = new String(inputByteArray, StandardCharsets.UTF_8);
                    validationResponse = OASParserUtil.validateAPIDefinition(openAPIContent, returnContent);
                }
            } else if (apiDefinition != null) {
                validationResponse = OASParserUtil.validateAPIDefinition(apiDefinition, returnContent);
            }
        } else if (APIConstants.ParserType.GRAPHQL.name().equals(type.toUpperCase())) {
            if (fileName.endsWith(".graphql") || fileName.endsWith(".txt") ||
                    fileName.endsWith(".sdl")) {
                validationResponse = OASParserUtil.validateGraphQLSchema(
                        new String(inputByteArray, StandardCharsets.UTF_8),
                        returnContent);
            } else {
                OASParserUtil.addErrorToValidationResponse(validationResponse,
                        "Invalid definition file type provided.");
            }
        } else if (APIConstants.ParserType.GRPC.name().equals(type.toUpperCase())) {
            if (inputByteArray != null && inputByteArray.length > 0) {
                if (fileName.endsWith(".zip") || fileName.endsWith(".proto")) {
                    validationResponse = ProtoParserUtil.validateGRPCAPIDefinition(inputByteArray);
                } else {
                    ProtoParserUtil.addErrorToValidationResponse(validationResponse,
                            "Invalid definition file type provided.");
                }
            } else {
                ProtoParserUtil.addErrorToValidationResponse(validationResponse,
                        "Invalid definition file type provided.");
            }
        }
        return validationResponse;
    }

    public static API getGRPCAPIFromProtoDefinition(byte[] definition, String fileName) {
        ProtoParser protoParser = new ProtoParser();
        try {
            return protoParser.getAPIFromProtoFile(definition, fileName);
        } catch (APIManagementException e) {
            throw new RuntimeException(e);
        }
    }

    public static Set<URITemplate> generateUriTemplatesFromAPIDefinition(String apiType, String content)
            throws APIManagementException {

        APIDefinition parser = DefinitionParserFactory.getParser(apiType);
        if (parser != null) {
            return parser.getURITemplates(content);
        } else {
            throw new APIManagementException("Couldn't find parser", ExceptionCodes.INTERNAL_ERROR);
        }
    }

    /**
     * This method used to retrieve API definition merged with custom properties.
     *
     * @param api        API object
     * @param definition user given definition.
     * @return definition
     * @throws APIManagementException
     */
    public static String generateDefinition(API api, String definition) throws APIManagementException {

        APIDefinition parser = DefinitionParserFactory.getParser(api);
        return parser.generateAPIDefinition(api, definition);
    }

    public static API getAPIFromDefinition(String definition, String apiType)
            throws APIManagementException {
        if (apiType.equalsIgnoreCase(APIConstants.GRAPHQL_API)) {
            return getGQLAPIFromDefinition(definition);
        } else {
            APIDefinition parser = DefinitionParserFactory.getParser(apiType);
            if (parser != null) {
                return parser.getAPIFromDefinition(definition);
            }
        }
        throw new APIManagementException("Definition parser not found");
    }

    private static API getGQLAPIFromDefinition(String definition) {
        SchemaParser schemaParser = new SchemaParser();
        TypeDefinitionRegistry registry = schemaParser.parse(definition);
        List<URITemplate> combinedUriTemplates = new ArrayList<>();

        // Directly add all URI templates for query, mutation, and subscription into a
        // combined list
        combinedUriTemplates
                .addAll(GraphQLSchemaDefinition.extractGraphQLOperationList(registry, APIConstants.GRAPHQL_QUERY));
        combinedUriTemplates
                .addAll(GraphQLSchemaDefinition.extractGraphQLOperationList(registry, APIConstants.GRAPHQL_MUTATION));
        combinedUriTemplates.addAll(
                GraphQLSchemaDefinition.extractGraphQLOperationList(registry, APIConstants.GRAPHQL_SUBSCRIPTION));

        API api = new API();
        api.setUriTemplates(combinedUriTemplates.toArray(new URITemplate[0]));
        api.setGraphQLSchema(definition);
        return api;
    }

    private RuntimeAPICommonUtil() {

    }

}
