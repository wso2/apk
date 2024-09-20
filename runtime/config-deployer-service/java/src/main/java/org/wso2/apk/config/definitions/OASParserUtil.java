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

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.databind.node.ObjectNode;

import graphql.schema.GraphQLSchema;
import graphql.schema.idl.SchemaParser;
import graphql.schema.idl.TypeDefinitionRegistry;
import graphql.schema.idl.UnExecutableSchemaGenerator;
import graphql.schema.validation.SchemaValidationError;
import graphql.schema.validation.SchemaValidator;
import io.swagger.models.RefModel;
import io.swagger.models.RefPath;
import io.swagger.models.RefResponse;
import io.swagger.models.Response;
import io.swagger.models.Swagger;
import io.swagger.models.parameters.RefParameter;
import io.swagger.models.properties.RefProperty;
import io.swagger.parser.SwaggerParser;
import io.swagger.parser.util.DeserializationUtils;
import io.swagger.parser.util.SwaggerDeserializationResult;
import io.swagger.v3.core.util.Json;
import io.swagger.v3.core.util.Yaml;
import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.Paths;
import io.swagger.v3.oas.models.media.ComposedSchema;
import io.swagger.v3.oas.models.media.ObjectSchema;
import io.swagger.v3.oas.models.media.Schema;
import io.swagger.v3.parser.ObjectMapperFactory;
import io.swagger.v3.parser.OpenAPIV3Parser;
import io.swagger.v3.parser.core.models.ParseOptions;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.json.JSONObject;
import org.wso2.apk.config.APIConstants;
import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.api.APIDefinitionValidationResponse;
import org.wso2.apk.config.api.APIManagementException;
import org.wso2.apk.config.api.ErrorHandler;
import org.wso2.apk.config.api.ErrorItem;
import org.wso2.apk.config.api.ExceptionCodes;
import org.wso2.apk.config.api.Info;
import org.wso2.apk.config.model.URITemplate;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.UUID;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

/**
 * Provide common functions related to OAS
 */
public class OASParserUtil {
    private static final Log log = LogFactory.getLog(OASParserUtil.class);
    private static final String OPENAPI_RESOURCE_KEY = "paths";
    private static final String[] UNSUPPORTED_RESOURCE_BLOCKS = new String[] { "servers" };
    private static final APIDefinition oas3Parser = new OAS3Parser();

    public enum SwaggerVersion {
        SWAGGER,
        OPEN_API,
    }

    private static final String REQUEST_BODIES = "requestBodies";
    private static final String SCHEMAS = "schemas";
    private static final String PARAMETERS = "parameters";
    private static final String RESPONSES = "responses";
    private static final String HEADERS = "headers";

    private static final String REF_PREFIX = "#/components/";

    static class SwaggerUpdateContext {
        private final Paths paths = new Paths();
        private final Map<String, Set<String>> referenceObjectMap = new HashMap<>();
        private final Set<Components> aggregatedComponents = new HashSet<>();

        SwaggerUpdateContext() {
            referenceObjectMap.put(REQUEST_BODIES, new HashSet<>());
            referenceObjectMap.put(SCHEMAS, new HashSet<>());
            referenceObjectMap.put(PARAMETERS, new HashSet<>());
            referenceObjectMap.put(RESPONSES, new HashSet<>());
            referenceObjectMap.put(HEADERS, new HashSet<>());
        }

        Paths getPaths() {
            return paths;
        }

        Map<String, Set<String>> getReferenceObjectMapping() {
            return referenceObjectMap;
        }

        public Set<Components> getAggregatedComponents() {
            return aggregatedComponents;
        }
    }

    public static SwaggerVersion getSwaggerVersion(String apiDefinition) throws APIManagementException {
        ObjectMapper mapper;
        if (apiDefinition.trim().startsWith("{")) {
            mapper = ObjectMapperFactory.createJson();
        } else {
            mapper = ObjectMapperFactory.createYaml();
        }
        JsonNode rootNode;
        try {
            rootNode = mapper.readTree(apiDefinition.getBytes());
        } catch (IOException e) {
            throw new APIManagementException("Error occurred while parsing OAS definition", e,
                    ExceptionCodes.OPENAPI_PARSE_EXCEPTION);
        }
        ObjectNode node = (ObjectNode) rootNode;
        JsonNode openapi = node.get("openapi");
        if (openapi != null && openapi.asText().startsWith("3.")) {
            return SwaggerVersion.OPEN_API;
        }
        JsonNode swagger = node.get("swagger");
        if (swagger != null) {
            return SwaggerVersion.SWAGGER;
        }

        throw new APIManagementException("Invalid OAS definition provided.",
                ExceptionCodes.MALFORMED_OPENAPI_DEFINITON);
    }

    private static List<String> addSchemaOfSchema(Schema schema) {
        List<String> references = new ArrayList<String>();
        ObjectSchema os = (ObjectSchema) schema;
        if (os.getProperties() != null) {
            for (String propertyName : os.getProperties().keySet()) {
                if (os.getProperties().get(propertyName) instanceof ComposedSchema) {
                    ComposedSchema cs = (ComposedSchema) os.getProperties().get(propertyName);
                    if (cs.getAllOf() != null) {
                        for (Schema sc : cs.getAllOf()) {
                            references.add(sc.get$ref());
                        }
                    } else if (cs.getAnyOf() != null) {
                        for (Schema sc : cs.getAnyOf()) {
                            references.add(sc.get$ref());
                        }
                    } else if (cs.getOneOf() != null) {
                        for (Schema sc : cs.getOneOf()) {
                            references.add(sc.get$ref());
                        }
                    } else {
                        log.error("Unidentified schema. The schema is not available in the API definition.");
                    }
                }
            }
        }
        return references;
    }

    private static void addToReferenceObjectMap(String ref, SwaggerUpdateContext context) {
        Map<String, Set<String>> referenceObjectMap = context.getReferenceObjectMapping();
        final String category = getComponentCategory(ref);
        if (referenceObjectMap.containsKey(category)) {
            Set<String> refObjects = referenceObjectMap.get(category);
            refObjects.add(getRefKey(ref));
        }
    }

    private static String getRefKey(String ref) {
        String[] split = ref.split("/");
        return split[split.length - 1];
    }

    private static String getComponentCategory(String ref) {
        String[] remainder = ref.split(REF_PREFIX);

        if (remainder.length == 2) {
            String[] split = remainder[1].split("/");

            if (split.length == 2) {
                return split[0];
            }
        }

        return "";
    }

    public static File checkMasterSwagger(File archiveDirectory) throws APIManagementException {
        File masterSwagger = null;
        if ((new File(archiveDirectory + "/" + APIConstants.OPENAPI_MASTER_JSON)).exists()) {
            masterSwagger = new File(archiveDirectory + "/" + APIConstants.OPENAPI_MASTER_JSON);
            return masterSwagger;
        } else if ((new File(archiveDirectory + "/" + APIConstants.OPENAPI_MASTER_YAML)).exists()) {
            masterSwagger = new File(archiveDirectory + "/" + APIConstants.OPENAPI_MASTER_YAML);
            return masterSwagger;
        } else {
            throw new APIManagementException("Could not find a master swagger file with the name of swagger.json " +
                    "/swagger.yaml", ExceptionCodes.INTERNAL_ERROR);
        }
    }

    /**
     * Update the APIDefinitionValidationResponse object with success state using
     * the values given
     *
     * @param validationResponse    APIDefinitionValidationResponse object to be
     *                              updated
     * @param originalAPIDefinition original API Definition
     * @param openAPIVersion        version of OpenAPI Spec (2.0 or 3.0.0)
     * @param title                 title of the OpenAPI Definition
     * @param version               version of the OpenAPI Definition
     * @param context               base path of the OpenAPI Definition
     * @param description           description of the OpenAPI Definition
     */
    public static void updateValidationResponseAsSuccess(APIDefinitionValidationResponse validationResponse,
            String originalAPIDefinition, String openAPIVersion, String title, String version, String context,
            String description, List<String> endpoints) {
        validationResponse.setValid(true);
        validationResponse.setContent(originalAPIDefinition);
        Info info = new Info();
        info.setOpenAPIVersion(openAPIVersion);
        info.setName(title);
        info.setVersion(version);
        info.setContext(context);
        info.setDescription(description);
        info.setEndpoints(endpoints);
        validationResponse.setInfo(info);
    }

    /**
     * Add error item with the provided message to the provided validation response
     * object
     *
     * @param validationResponse APIDefinitionValidationResponse object
     * @param errMessage         error message
     * @return added ErrorItem object
     */
    public static ErrorItem addErrorToValidationResponse(APIDefinitionValidationResponse validationResponse,
            String errMessage) {
        ErrorItem errorItem = new ErrorItem();
        errorItem.setErrorCode(ExceptionCodes.OPENAPI_PARSE_EXCEPTION.getErrorCode());
        errorItem.setMessage(ExceptionCodes.OPENAPI_PARSE_EXCEPTION.getErrorMessage());
        errorItem.setDescription(errMessage);
        validationResponse.getErrorItems().add(errorItem);
        return errorItem;
    }

    /**
     * Creates a json string using the swagger object.
     *
     * @param swaggerObj swagger object
     * @return json string using the swagger object
     * @throws APIManagementException error while creating swagger json
     */
    public static String getSwaggerJsonString(Swagger swaggerObj) throws APIManagementException {
        ObjectMapper mapper = new ObjectMapper();
        mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
        mapper.enable(SerializationFeature.INDENT_OUTPUT);

        // this is to ignore "originalRef" in schema objects
        mapper.addMixIn(RefModel.class, IgnoreOriginalRefMixin.class);
        mapper.addMixIn(RefProperty.class, IgnoreOriginalRefMixin.class);
        mapper.addMixIn(RefPath.class, IgnoreOriginalRefMixin.class);
        mapper.addMixIn(RefParameter.class, IgnoreOriginalRefMixin.class);
        mapper.addMixIn(RefResponse.class, IgnoreOriginalRefMixin.class);

        // this is to ignore "responseSchema" in response schema objects
        mapper.addMixIn(Response.class, ResponseSchemaMixin.class);
        try {
            return new String(mapper.writeValueAsBytes(swaggerObj));
        } catch (JsonProcessingException e) {
            throw new APIManagementException("Error while generating Swagger json from model", e);
        }
    }

    /**
     * Sets the scopes to the URL template object using the given list of scopes
     *
     * @param template       URL template
     * @param resourceScopes list of scopes of the resource
     * @param apiScopes      set of scopes defined for the API
     * @return URL template after setting the scopes
     */
    public static URITemplate setScopesToTemplate(URITemplate template, List<String> resourceScopes,
            String[] apiScopes) throws APIManagementException {

        for (String scopeName : resourceScopes) {
            if (StringUtils.isNotBlank(scopeName)) {
                String scope = ParserUtil.findScopeByKey(apiScopes, scopeName);
                if (scope == null) {
                    throw new APIManagementException("Resource Scope '" + scopeName + "' not found.",
                            ExceptionCodes.SCOPE_NOT_FOUND);
                }
                template.setScopes(scope);
            }
        }
        return template;
    }

    /**
     * Sort scopes by name.
     * This method was added to display scopes in publisher in a sorted manner.
     *
     * @param scopeSet
     * @return Scope set
     */
    static String[] sortScopes(Set<String> scopeSet) {

        String[] scopeArray = scopeSet.toArray(new String[scopeSet.size()]);
        Arrays.sort(scopeArray);
        return scopeArray;
    }

    public static void copyOperationVendorExtensions(Map<String, Object> existingExtensions,
            Map<String, Object> updatedVendorExtensions) {
        if (existingExtensions.get(APIConstants.SWAGGER_X_AUTH_TYPE) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_AUTH_TYPE, existingExtensions
                    .get(APIConstants.SWAGGER_X_AUTH_TYPE));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_THROTTLING_TIER) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_THROTTLING_TIER, existingExtensions
                    .get(APIConstants.SWAGGER_X_THROTTLING_TIER));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_THROTTLING_BANDWIDTH) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_THROTTLING_BANDWIDTH, existingExtensions
                    .get(APIConstants.SWAGGER_X_THROTTLING_BANDWIDTH));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_MEDIATION_SCRIPT) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_MEDIATION_SCRIPT, existingExtensions
                    .get(APIConstants.SWAGGER_X_MEDIATION_SCRIPT));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_WSO2_SECURITY) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_WSO2_SECURITY, existingExtensions
                    .get(APIConstants.SWAGGER_X_WSO2_SECURITY));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_SCOPE) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_SCOPE, existingExtensions
                    .get(APIConstants.SWAGGER_X_SCOPE));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_AMZN_RESOURCE_NAME) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_AMZN_RESOURCE_NAME, existingExtensions
                    .get(APIConstants.SWAGGER_X_AMZN_RESOURCE_NAME));
        }
        if (existingExtensions.get(APIConstants.SWAGGER_X_AMZN_RESOURCE_TIMEOUT) != null) {
            updatedVendorExtensions.put(APIConstants.SWAGGER_X_AMZN_RESOURCE_TIMEOUT, existingExtensions
                    .get(APIConstants.SWAGGER_X_AMZN_RESOURCE_TIMEOUT));
        }
    }

    /**
     * Extract the archive file and validates the openAPI definition
     *
     * @param inputByteArray file as input stream
     * @param returnContent  whether to return the content of the definition in the
     *                       response DTO
     * @return APIDefinitionValidationResponse
     * @throws APIManagementException if error occurred while parsing definition
     */
    public static APIDefinitionValidationResponse extractAndValidateOpenAPIArchive(byte[] inputByteArray,
            boolean returnContent) throws APIManagementException {
        String path = System.getProperty(APIConstants.JAVA_IO_TMPDIR) + File.separator +
                APIConstants.OPENAPI_ARCHIVES_TEMP_FOLDER + File.separator + UUID.randomUUID();
        String archivePath = path + File.separator + APIConstants.OPENAPI_ARCHIVE_ZIP_FILE;
        String extractedLocation = OASParserUtil.extractUploadedArchive(inputByteArray,
                APIConstants.OPENAPI_EXTRACTED_DIRECTORY, archivePath, path);
        File[] listOfFiles = new File(extractedLocation).listFiles();
        File archiveDirectory = null;
        if (listOfFiles != null) {
            if (listOfFiles.length > 1) {
                throw new APIManagementException("Swagger Definitions should be placed under one root folder.");
            }
            for (File file : listOfFiles) {
                if (file.isDirectory()) {
                    archiveDirectory = file.getAbsoluteFile();
                    break;
                }
            }
        }
        // Verify whether the zipped input is archive or file.
        // If it is a single swagger file without remote references it can be imported
        // directly, without zipping.
        if (archiveDirectory == null) {
            throw new APIManagementException("Could not find an archive in the given ZIP file.");
        }
        File masterSwagger = checkMasterSwagger(archiveDirectory);
        String content;
        try {
            InputStream masterInputStream = new FileInputStream(masterSwagger);
            content = IOUtils.toString(masterInputStream, StandardCharsets.UTF_8);
        } catch (IOException e) {
            throw new APIManagementException("Error reading master swagger file" + e);
        }
        String openAPIContent = "";
        SwaggerVersion version;
        version = getSwaggerVersion(content);
        String filePath = masterSwagger.getAbsolutePath();
        if (SwaggerVersion.OPEN_API.equals(version)) {
            OpenAPIV3Parser openAPIV3Parser = new OpenAPIV3Parser();
            ParseOptions options = new ParseOptions();
            options.setResolve(true);
            OpenAPI openAPI = openAPIV3Parser.read(filePath, null, options);
            openAPIContent = Json.pretty(openAPI);
        } else if (SwaggerVersion.SWAGGER.equals(version)) {
            SwaggerParser parser = new SwaggerParser();
            Swagger swagger = parser.read(filePath, null, true);
            try {
                openAPIContent = Yaml.pretty().writeValueAsString(swagger);
            } catch (IOException e) {
                throw new APIManagementException("Error in converting swagger to openAPI content. " + e);
            }
        }
        APIDefinitionValidationResponse apiDefinitionValidationResponse;
        apiDefinitionValidationResponse = OASParserUtil.validateAPIDefinition(openAPIContent, returnContent);
        return apiDefinitionValidationResponse;
    }

    public static String extractUploadedArchive(byte[] byteArray, String importedDirectoryName,
            String apiArchiveLocation, String extractLocation) throws APIManagementException {
        String archiveExtractLocation;

        try (ByteArrayInputStream uploadedApiArchiveInputStream = new ByteArrayInputStream(byteArray)) {
            // create api import directory structure
            createDirectory(extractLocation);
            // create archive
            createArchiveFromInputStream(uploadedApiArchiveInputStream, apiArchiveLocation);
            // extract the archive
            archiveExtractLocation = extractLocation + File.separator + importedDirectoryName;
            extractArchive(apiArchiveLocation, archiveExtractLocation);

        } catch (APIManagementException | IOException e) {
            deleteDirectory(extractLocation);
            String errorMsg = "Error in accessing uploaded API archive";
            log.error(errorMsg, e);
            throw new APIManagementException(errorMsg, e);
        }
        return archiveExtractLocation;
    }

    /**
     * Delete a given directory
     *
     * @param path Path to the directory to be deleted
     * @throws APIManagementException if unable to delete the directory
     */
    public static void deleteDirectory(String path) throws APIManagementException {
        try {
            FileUtils.deleteDirectory(new File(path));
        } catch (IOException e) {
            String errorMsg = "Error while deleting directory : " + path;
            log.error(errorMsg, e);
            throw new APIManagementException(errorMsg, e);
        }
    }

    /**
     * Extracts a a given zip archive
     *
     * @param archiveFilePath path of the zip archive
     * @param destination     extract location
     * @return name of the extracted zip archive
     * @throws APIManagementException if an error occurs while extracting the
     *                                archive
     */
    public static String extractArchive(String archiveFilePath, String destination)
            throws APIManagementException {
        int bufferSize = 512;
        long sizeLimit = 0x6400000; // Max size of unzipped data, 100MB
        int maxEntryCount = 1024;
        String archiveName = null;

        try {
            FileInputStream fis = new FileInputStream(archiveFilePath);
            ZipInputStream zis = new ZipInputStream(new BufferedInputStream(fis));
            ZipEntry entry;
            int entries = 0;
            long total = 0;

            // Process each entry
            while ((entry = zis.getNextEntry()) != null) {
                String currentEntry = entry.getName();
                int index = 0;
                // This index variable is used to get the extracted folder name; that is root
                // directory
                if (index == 0 && currentEntry.indexOf('/') != -1) {
                    archiveName = currentEntry.substring(0, currentEntry.indexOf('/'));
                    --index;
                }

                File destinationFile = new File(destination, currentEntry);
                File destinationParent = destinationFile.getParentFile();
                String canonicalizedDestinationFilePath = destinationFile.getCanonicalPath();

                if (!canonicalizedDestinationFilePath.startsWith(new File(destination).getCanonicalPath())) {
                    String errorMessage = "Attempt to upload invalid zip archive with file at " + currentEntry +
                            ". File path is outside target directory";
                    log.error(errorMessage);
                    throw new APIManagementException(errorMessage);
                }
                if (entry.isDirectory()) {
                    log.debug("Creating directory " + destinationFile.getAbsolutePath());
                    destinationFile.mkdir();
                    continue;
                }

                // create the parent directory structure
                if (destinationParent.mkdirs()) {
                    log.debug("Creation of folder is successful. Directory Name : " + destinationParent.getName());
                }

                int count;
                byte[] data = new byte[bufferSize];
                FileOutputStream fos = new FileOutputStream(destinationFile);
                BufferedOutputStream dest = new BufferedOutputStream(fos, bufferSize);
                while (total + bufferSize <= sizeLimit && (count = zis.read(data, 0, bufferSize)) != -1) {
                    dest.write(data, 0, count);
                    total += count;
                }
                dest.flush();
                dest.close();
                zis.closeEntry();
                entries++;
                if (entries > maxEntryCount) {
                    throw new APIManagementException("Too many files to unzip.");
                }
                if (total + bufferSize > sizeLimit) {
                    throw new APIManagementException("File being unzipped is too big.");
                }
            }
            return archiveName;
        } catch (IOException e) {
            String errorMsg = "Failed to extract archive file: " + archiveFilePath + " to destination: " + destination;
            log.error(errorMsg, e);
            throw new APIManagementException(errorMsg, e);
        }
    }

    /**
     * Creates a zip archive from the given {@link InputStream} inputStream
     *
     * @param inputStream {@link InputStream} instance
     * @param archivePath path to create the zip archive
     * @throws APIManagementException if an error occurs while creating the archive
     */
    public static void createArchiveFromInputStream(InputStream inputStream, String archivePath)
            throws APIManagementException {
        try (FileOutputStream outFileStream = new FileOutputStream(new File(archivePath))) {
            IOUtils.copy(inputStream, outFileStream);
        } catch (IOException e) {
            String errorMsg = "Error in Creating archive from data.";
            log.error(errorMsg, e);
            throw new APIManagementException(errorMsg, e);
        }
    }

    /**
     * Creates a directory
     *
     * @param path path of the directory to create
     * @throws APIManagementException if an error occurs while creating the
     *                                directory
     */
    public static void createDirectory(String path) throws APIManagementException {
        try {
            Files.createDirectories(java.nio.file.Paths.get(path));
        } catch (IOException e) {
            String msg = "Error in creating directory at: " + path;
            log.error(msg, e);
            throw new APIManagementException(msg, e);
        }
    }

    /**
     * Try to validate a give openAPI definition using OpenAPI 3 parser
     *
     * @param apiDefinition     definition
     * @param returnJsonContent whether to return definition as a json content
     * @return APIDefinitionValidationResponse
     * @throws APIManagementException if error occurred while parsing definition
     */
    public static APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition, boolean returnJsonContent)
            throws APIManagementException {
        String apiDefinitionProcessed = apiDefinition;
        if (!apiDefinition.trim().startsWith("{")) {
            try {
                JsonNode jsonNode = DeserializationUtils.readYamlTree(apiDefinition,
                        new SwaggerDeserializationResult());
                apiDefinitionProcessed = jsonNode.toString();
            } catch (IOException e) {
                throw new APIManagementException("Error while reading API definition yaml", e);
            }
        }
        apiDefinitionProcessed = removeUnsupportedBlocksFromResources(apiDefinitionProcessed);
        if (apiDefinitionProcessed != null) {
            apiDefinition = apiDefinitionProcessed;
        }
        APIDefinitionValidationResponse validationResponse = oas3Parser.validateAPIDefinition(apiDefinition,
                returnJsonContent);
        if (!validationResponse.isValid()) {
            // for (ErrorHandler handler : validationResponse.getErrorItems()) {
            // if (ExceptionCodes.INVALID_OAS3_FOUND.getErrorCode() ==
            // handler.getErrorCode()) {
            // return tryOAS2Validation(apiDefinition, returnJsonContent);
            // }
            // }
        }
        return validationResponse;
    }

    /**
     * Validate graphQL Schema
     * 
     * @return Validation response
     */
    public static APIDefinitionValidationResponse validateGraphQLSchema(String apiDefinition,
            boolean returnGraphQLSchemaContent) {
        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        ArrayList<ErrorHandler> errors = new ArrayList<>();

        if (apiDefinition.isBlank()) {
            validationResponse.setValid(false);
            errors.add(ExceptionCodes.GRAPHQL_SCHEMA_CANNOT_BE_NULL);
            validationResponse.setErrorItems(errors);
        }

        SchemaParser schemaParser = new SchemaParser();
        Set<SchemaValidationError> validationErrors = new HashSet<>();
        try {
            TypeDefinitionRegistry typeRegistry = schemaParser.parse(apiDefinition);
            GraphQLSchema graphQLSchema = UnExecutableSchemaGenerator.makeUnExecutableSchema(typeRegistry);
            SchemaValidator schemaValidation = new SchemaValidator();
            validationErrors = schemaValidation.validateSchema(graphQLSchema);
            if (validationErrors.toArray().length > 0) {
                validationResponse.setValid(false);
                errors.add(ExceptionCodes.API_NOT_GRAPHQL);
                validationResponse.setErrorItems(errors);
            } else {
                validationResponse.setValid(true);
                validationResponse.setContent(apiDefinition);
            }
        } catch (Exception e) {
            OASParserUtil.addErrorToValidationResponse(validationResponse, e.getMessage());
            validationResponse.setValid(false);
            errors.add(new ErrorItem("API Definition Validation Error", "API Definition is invalid", 400, 400));
            validationResponse.setErrorItems(errors);
        }

        return validationResponse;
    }

    /**
     * This method removes the unsupported json blocks from the given json string.
     *
     * @param jsonString Open api specification from which unsupported blocks must
     *                   be removed.
     * @return String open api specification without unsupported blocks. Null value
     *         if there is no unsupported blocks.
     */
    public static String removeUnsupportedBlocksFromResources(String jsonString) {
        JSONObject jsonObject = new JSONObject(jsonString);
        boolean definitionUpdated = false;
        if (jsonObject.has(OPENAPI_RESOURCE_KEY)) {
            JSONObject paths = jsonObject.optJSONObject(OPENAPI_RESOURCE_KEY);
            if (paths != null) {
                for (String unsupportedBlockKey : UNSUPPORTED_RESOURCE_BLOCKS) {
                    boolean result = removeBlocksRecursivelyFromJsonObject(unsupportedBlockKey, paths, false);
                    definitionUpdated = definitionUpdated || result;
                }
            }
        }
        if (definitionUpdated) {
            ObjectMapper om = new ObjectMapper();
            om.configure(SerializationFeature.ORDER_MAP_ENTRIES_BY_KEYS, true);
            try {
                Map<String, Object> map = om.readValue(jsonObject.toString(), HashMap.class);
                String json = om.writeValueAsString(map);
                return json;
            } catch (JsonProcessingException e) {
                return null;
            }
        } else {
            return null;
        }
    }

    /**
     * This method removes provided key from the json object recursively.
     *
     * @param keyToBeRemoved, Key to remove from open api spec.
     * @param jsonObject,     Open api spec as json object.
     */
    private static boolean removeBlocksRecursivelyFromJsonObject(String keyToBeRemoved, JSONObject jsonObject,
            boolean definitionUpdated) {
        if (jsonObject == null) {
            return definitionUpdated;
        }
        if (jsonObject.has(keyToBeRemoved)) {
            jsonObject.remove(keyToBeRemoved);
            definitionUpdated = true;
        }
        for (Object key : jsonObject.keySet()) {
            JSONObject subObj = jsonObject.optJSONObject(key.toString());
            if (subObj != null) {
                boolean result = removeBlocksRecursivelyFromJsonObject(keyToBeRemoved, subObj, definitionUpdated);
                definitionUpdated = definitionUpdated || result;
            }
        }
        return definitionUpdated;
    }
}
