package org.wso2.apk.config.definitions;

import java.io.*;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.config.api.*;
import org.wso2.apk.config.model.API;
import org.wso2.apk.config.model.proto.ProtoFile;
import org.wso2.apk.config.model.proto.Service;
import org.wso2.apk.config.model.SwaggerData;
import org.wso2.apk.config.model.URITemplate;

import io.swagger.v3.core.util.Json;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.Paths;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;

public class ProtoParser extends APIDefinition {
    private static final Log log = LogFactory.getLog(ProtoParser.class);

    public ProtoParser() {
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
    public APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition, boolean returnContent) {
        return new APIDefinitionValidationResponse();
    }

    @Override
    public API getAPIFromDefinition(String content) throws APIManagementException {
        throw new UnsupportedOperationException("Unimplemented method 'getAPIFromDefinition'");
    }

    @Override
    public String processOtherSchemeScopes(String resourceConfigsJSON) throws APIManagementException {
        return resourceConfigsJSON;
    }

    @Override
    public String getType() {
        return null;
    }

    @Override
    public boolean canHandleDefinition(String definition) {
        return true;
    }

    @Override
    public String generateAPIDefinition(API api, String swagger) throws APIManagementException {
        return null;
    }

    public API getAPIFromProtoFile(byte[] content, String fileName) throws APIManagementException {
        try {
            API api = new API();
            ProtoFile protoFile = new ProtoFile();
            List<URITemplate> uriTemplates = new ArrayList<>();
            if (fileName.endsWith(".zip")) {
                List<byte[]> protoContents = extractProtoFilesFromZip(content);
                for (byte[] protoContent : protoContents) {
                    uriTemplates.addAll(processProtoFile(protoContent, protoFile));
                }
            } else {
                uriTemplates = processProtoFile(content, protoFile);
            }
            api.setBasePath(protoFile.getBasePath());
            api.setProtoDefinition(new String(content, java.nio.charset.StandardCharsets.UTF_8));
            api.setVersion(protoFile.getVersion());
            api.setName(protoFile.getApiName());
            api.setUriTemplates(uriTemplates.toArray(new URITemplate[0]));
            return api;
        } catch (Exception e) {
            e.printStackTrace();
            throw new APIManagementException(e);
        }
    }

    private List<URITemplate> processProtoFile(byte[] definition, ProtoFile protoFile) throws APIManagementException {
        String content = new String(definition, java.nio.charset.StandardCharsets.UTF_8);
        String packageString = getPackageString(content);
        List<URITemplate> uriTemplates = new ArrayList<>();
        StringBuilder apiName = new StringBuilder().append(protoFile.getApiName());
        if (packageString == null && protoFile.getPackageName() != null) {
            throw new APIManagementException("Package string has not been defined in proto file");
        }
        String packageName = getPackageName(packageString);
        if (packageName == null) {
            packageName = packageString;
        }
        List<Service> services = new ArrayList<>();
        protoFile.setVersion(getVersion(packageString));
        protoFile.setBasePath(getBasePath(packageString));
        List<String> serviceBlocks = extractServiceBlocks(content);

        for (String serviceBlock : serviceBlocks) {
            String serviceName = getServiceName(serviceBlock);
            List<String> methodNames = extractMethodNames(serviceBlock);
            Service service = new Service(serviceName, methodNames);
            services.add(service);
            for (String method : service.getMethods()) {
                URITemplate uriTemplate = new URITemplate();
                uriTemplate.setUriTemplate(packageName + "." + service.getServiceName());
                uriTemplate.setVerb(method);
                uriTemplates.add(uriTemplate);
            }
        }

        for (Service service : services) {
            apiName.append(service.getServiceName()).append("-");
        }
        protoFile.setApiName(apiName.toString());
        return uriTemplates;
    }

    private List<byte[]> extractProtoFilesFromZip(byte[] zipContent) throws APIManagementException {
        List<byte[]> protoFiles = new ArrayList<>();
        try (ByteArrayInputStream bais = new ByteArrayInputStream(zipContent);
                ZipInputStream zis = new ZipInputStream(bais)) {

            ZipEntry zipEntry;
            while ((zipEntry = zis.getNextEntry()) != null) {
                if (zipEntry.getName().endsWith(".proto")) {
                    protoFiles.add(readProtoFileBytesFromZip(zis));
                }
            }
        } catch (IOException e) {
            e.printStackTrace();
            throw new APIManagementException("Failed to process zip file", e);
        }
        return protoFiles;
    }

    private String readProtoFileFromZip(InputStream zis) throws IOException {
        StringBuilder protoContent = new StringBuilder();
        BufferedReader reader = new BufferedReader(new InputStreamReader(zis));
        String line;
        while ((line = reader.readLine()) != null) {
            protoContent.append(line).append("\n");
        }
        return protoContent.toString();
    }

    private byte[] readProtoFileBytesFromZip(ZipInputStream zis) throws IOException {
        ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
        byte[] buffer = new byte[1024];
        int bytesRead;
        while ((bytesRead = zis.read(buffer)) != -1) {
            byteArrayOutputStream.write(buffer, 0, bytesRead);
        }
        return byteArrayOutputStream.toByteArray();
    }

    boolean validateProtoContent(byte[] fileContent, String fileName) {
        try {
            // ProtoFile protoFile = getProtoFileFromDefinition(fileContent, fileName);
            return true;
        } catch (Exception e) {
            log.error("Proto definition validation failed for " + fileName + ": " + e.getMessage());
            return false;
        }
    }

    public void validateGRPCAPIDefinition(byte[] inputByteArray, String fileName,
            APIDefinitionValidationResponse validationResponse, ArrayList<ErrorHandler> errors) {
        try {
            if (fileName.endsWith(".zip")) {
                try (ZipInputStream zis = new ZipInputStream(new ByteArrayInputStream(inputByteArray))) {
                    ZipEntry zipEntry;
                    while ((zipEntry = zis.getNextEntry()) != null) {
                        if (!zipEntry.isDirectory() && zipEntry.getName().endsWith(".proto")) {
                            byte[] protoFileContentBytes = zis.readAllBytes();
                            boolean validated = validateProtoContent(protoFileContentBytes, fileName);
                            if (!validated) {
                                throw new APIManagementException(
                                        "Invalid definition file provided. " + "Please provide a valid .zip or .proto file.");
                            }
                        }
                    }
                }
            } else if (fileName.endsWith(".proto")) {
                boolean validated = validateProtoContent(inputByteArray, fileName);
                validationResponse.setValid(validated);
                validationResponse.setProtoContent(inputByteArray);
                if (!validated) {
                    throw new APIManagementException(
                            "Invalid definition file provided. " + "Please provide a valid .zip or .proto file.");
                }
            } else {
                throw new APIManagementException(
                        "Invalid definition file provided. " + "Please provide a valid .zip or .desc file.");
            }
        } catch (Exception e) {
            ProtoParserUtil.addErrorToValidationResponse(validationResponse, e.getMessage());
            validationResponse.setValid(false);
            errors.add(new ErrorItem("API Definition Validation Error", "API Definition is invalid", 400, 400));
            validationResponse.setErrorItems(errors);
        }
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

    public String getPackageString(String content) {
        // package string has the format "package something"
        Pattern packagePattern = Pattern.compile("package\\s+([\\w\\.]+);");
        Matcher packageMatcher = packagePattern.matcher(content);
        if (packageMatcher.find()) {
            if (packageMatcher.group().length() > 1) {
                return packageMatcher.group(1);
            }
        }
        log.error("Package has not been defined in the proto file");
        return null;
    }

    public String getVersion(String packageString) {
        Pattern versionPattern = Pattern.compile("v\\d+(\\.\\d+)*");
        Matcher versionMatcher = versionPattern.matcher(packageString);
        if (versionMatcher.find()) {
            return versionMatcher.group(0);
        }
        log.error("Version not found in proto file");
        return null;
    }

    public String getPackageName(String packageString) {
        Pattern namePattern = Pattern.compile("v\\d+(\\.\\d+)*\\.(\\w+)$");
        Matcher nameMatcher = namePattern.matcher(packageString);
        if (nameMatcher.find()) {
            return nameMatcher.group(2);
        }
        log.error("Package name not found in proto file.");
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
            return "/" + basePath;
        }
        log.error("Base path not found in proto file");
        return null;
    }
}
