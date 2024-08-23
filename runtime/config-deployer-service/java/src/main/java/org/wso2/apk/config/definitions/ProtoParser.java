package org.wso2.apk.config.definitions;

import java.io.*;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.Iterator;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.api.APIDefinitionValidationResponse;
import org.wso2.apk.config.api.APIManagementException;
import org.wso2.apk.config.api.ErrorHandler;
import org.wso2.apk.config.api.ErrorItem;
import org.wso2.apk.config.api.ExceptionCodes;
import org.wso2.apk.config.model.API;
import org.wso2.apk.config.model.SwaggerData;
import org.wso2.apk.config.model.URITemplate;

import io.swagger.v3.core.util.Json;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.Operation;
import io.swagger.v3.oas.models.PathItem;
import io.swagger.v3.oas.models.Paths;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.parser.OpenAPIV3Parser;
import io.swagger.v3.parser.core.models.SwaggerParseResult;

import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.TextFormat;
import com.google.protobuf.InvalidProtocolBufferException;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Path;

public class ProtoParser extends APIDefinition {

    public ProtoFile protoFile;

    public ProtoParser() {
    }

    public void setContent(String content) {
        protoFile = parseProtoContent(content);
    }

    public String getPackageString(String content) {
        Pattern packagePattern = Pattern.compile("package\\s+([\\w\\.]+);");
        Matcher packageMatcher = packagePattern.matcher(content);
        if (packageMatcher.find()) {
            return packageMatcher.group(1);
        }
        return null;
    }

    public String getVersion(String packageString) {
        Pattern versionPattern = Pattern.compile("v\\d+(\\.\\d+)*");
        Matcher versionMatcher = versionPattern.matcher(packageString);
        if (versionMatcher.find()) {
            return versionMatcher.group(0);
        }
        System.out.println("No version found");
        return null;
    }

    public String getPackageName(String packageString) {
        Pattern namePattern = Pattern.compile("v\\d+(\\.\\d+)*\\.(\\w+)");
        Matcher nameMatcher = namePattern.matcher(packageString);
        if (nameMatcher.find()) {
            return nameMatcher.group(2);
        }
        System.out.println("No name found");
        return null;
    }

    public String getBasePath(String packageString) {
        Pattern basePathPattern = Pattern.compile("^(.*?)v\\d");

        Matcher basePathMatcher = basePathPattern.matcher(packageString);
        if (basePathMatcher.find()) {
            String basePath = basePathMatcher.group(1);
            if (basePath.charAt(basePath.length() - 1) == '.') {
                basePath = basePath.substring(0, basePath.length() - 1);
            }
            return basePath;
        }
        System.out.println("No base path found");
        return null;
    }

    // Method to extract service blocks from a given text
    public List<String> extractServiceBlocks(String text) {
        // Regular expression pattern to match the service blocks
        String patternString = "service\\s+\\w+\\s*\\{[^{}]*(?:\\{[^{}]*\\}[^{}]*)*\\}";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString, Pattern.DOTALL);
        Matcher matcher = pattern.matcher(text);

        // Find all matches and append them to the result
        List<String> result = new ArrayList<>();
        while (matcher.find()) {
            result.add(matcher.group());
        }
        return result;
    }

    public List<String> extractMethodNames(String serviceBlock) {
        // Regular expression pattern to match the method names
        String patternString = "(?<=rpc\\s)\\w+";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString);
        Matcher matcher = pattern.matcher(serviceBlock);

        // Find all matches and append them to the result
        List<String> result = new ArrayList<>();
        while (matcher.find()) {
            result.add(matcher.group());
        }
        return result;
    }

    public String getServiceName(String serviceBlock) {
        // Regular expression pattern to match the service name
        String patternString = "(?<=service\\s)\\w+";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString);
        Matcher matcher = pattern.matcher(serviceBlock);

        // Find the first match and return it
        if (matcher.find()) {
            return matcher.group();
        }
        return null;
    }

    public ProtoFile parseProtoContent(String content) {
        ProtoFile protoFile = new ProtoFile();
        protoFile.services = new ArrayList<>();

        List<String> serviceBlocks = extractServiceBlocks(content);
        for (String serviceBlock : serviceBlocks) {
            Service service = new Service();
            service.name = getServiceName(serviceBlock);
            service.methods = new ArrayList<>();
            service.methods.addAll(extractMethodNames(serviceBlock));
            protoFile.services.add(service);
        }

        // Extract package name
        String packageName = getPackageString(content);
        protoFile.packageName = getPackageName(packageName);
        protoFile.version = getVersion(packageName);
        protoFile.basePath = getBasePath(packageName);

        // System.out.println(protoFile);

        return protoFile;
    }

    public List<String> getMethods(Service Service) {
        return Service.methods;
    }

    public List<Service> getServices() {
        return this.protoFile.services;
    }

    public class ProtoFile {
        public String packageName;
        public String basePath;
        public String version;
        public List<Service> services;

        @Override
        public String toString() {
            return "ProtoFile{" +
                    "packageName='" + packageName + '\'' +
                    ", basePath='" + basePath + '\'' +
                    ", version='" + version + '\'' +
                    ", services=" + services +
                    '}';
        }

    }

    public class Service {
        public String name;
        public List<String> methods;

        @Override
        public String toString() {
            return " Service{" +
                    "name='" + name + '\'' +
                    ", methods=" + methods +
                    '}';
        }
    }

    @Override
    public Set<URITemplate> getURITemplates(String resourceConfigsJSON) throws APIManagementException {
        // TODO Auto-generated method stub
        throw new UnsupportedOperationException("Unimplemented method 'getURITemplates'");
    }

    @Override
    public String[] getScopes(String resourceConfigsJSON) throws APIManagementException {
        // TODO Auto-generated method stub
        throw new UnsupportedOperationException("Unimplemented method 'getScopes'");
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

        // Create info object
        Info info = new Info();
        info.setTitle(swaggerData.getTitle());
        if (swaggerData.getDescription() != null) {
            info.setDescription(swaggerData.getDescription());
        }

        Contact contact = new Contact();
        // Create contact object and map business owner info
        if (swaggerData.getContactName() != null) {
            contact.setName(swaggerData.getContactName());
        }
        if (swaggerData.getContactEmail() != null) {
            contact.setEmail(swaggerData.getContactEmail());
        }
        if (swaggerData.getContactName() != null || swaggerData.getContactEmail() != null) {
            // put contact object to info object
            info.setContact(contact);
        }

        info.setVersion(swaggerData.getVersion());
        openAPI.setInfo(info);
        return Json.pretty(openAPI);
    }

    /**
     * Validate gRPC proto definition
     *
     * @return Validation response
     */
    public APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition,
            boolean returnContent) {
        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        ArrayList<ErrorHandler> errors = new ArrayList<>();
        try {
            if (apiDefinition.isBlank()) {
                validationResponse.setValid(false);
                errors.add(ExceptionCodes.GRPC_PROTO_DEFINTION_CANNOT_BE_NULL);
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

    @Override
    public API getAPIFromDefinition(String content) throws APIManagementException {
        // TODO Auto-generated method stub
        throw new UnsupportedOperationException("Unimplemented method 'getAPIFromDefinition'");
    }

    @Override
    public String processOtherSchemeScopes(String resourceConfigsJSON) throws APIManagementException {
        // TODO Auto-generated method stub
        throw new UnsupportedOperationException("Unimplemented method 'processOtherSchemeScopes'");
    }

    @Override
    public String getType() {
        // TODO Auto-generated method stub
        throw new UnsupportedOperationException("Unimplemented method 'getType'");
    }

    @Override
    public boolean canHandleDefinition(String definition) {
        return true;
    }

    @Override
    public String generateAPIDefinition(API api, String swagger) throws APIManagementException {
        return null;
    }

    public boolean validateProtoFile(String protoContent) {
        try {
            DescriptorProtos.FileDescriptorProto.Builder builder = DescriptorProtos.FileDescriptorProto.newBuilder();
            TextFormat.getParser().merge(protoContent, builder);
            // If parsing succeeds, return true
            return true;
        } catch (IOException e) {
            // If an exception occurs, the proto file is invalid
            System.err.println("Validation failed: " + e.getMessage());
            return false;
        }
    }

}