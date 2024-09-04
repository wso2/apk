package org.wso2.apk.config.definitions;

import java.io.ByteArrayInputStream;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

import com.google.protobuf.*;
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

    private static Descriptors.FileDescriptor resolveWellKnownType(String descriptorName)
            throws ClassNotFoundException, NoSuchMethodException, IllegalAccessException, InvocationTargetException {
        // Extract the proto file base name (e.g., "timestamp.proto" -> "Timestamp")
        String baseName = descriptorName.substring(descriptorName.lastIndexOf('/') + 1,
                descriptorName.lastIndexOf('.'));
        // Convert to CamelCase (e.g., "timestamp" -> "Timestamp")
        String className = baseName.substring(0, 1).toUpperCase() + baseName.substring(1);
        // Find the corresponding class in the com.google.protobuf package
        Class<?> clazz = Class.forName("com.google.protobuf." + className);
        // Use reflection to get the descriptor
        Method getDescriptorMethod = clazz.getMethod("getDescriptor");
        Descriptors.Descriptor descriptor = (Descriptors.Descriptor) getDescriptorMethod.invoke(null);
        return descriptor.getFile();
    }

    private static void populateProtoMap(Map<String, DescriptorProtos.FileDescriptorProto> protoMap,
            DescriptorProtos.FileDescriptorSet fileDescriptorSet) {
        for (DescriptorProtos.FileDescriptorProto fileDescriptorProto : fileDescriptorSet.getFileList()) {
            protoMap.put(fileDescriptorProto.getName(), fileDescriptorProto);
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

    public API getAPIFromProtoFile(byte[] definition, String fileName) throws APIManagementException {
        API api = new API();
        ProtoFile protoFile = getProtoFileFromDefinition(definition, fileName);
        StringBuilder apiName = new StringBuilder();
        List<String> sortedServices = new ArrayList<>();
        List<URITemplate> uriTemplates = new ArrayList<>();

        if (protoFile == null) {
            throw new APIManagementException("Error in validating API definition");
        }

        System.out.println("Printing services");
        for (Service service : protoFile.getServices()) {
            sortedServices.add(service.getServiceName());
            for (String method : service.getMethods()) {
                System.out.println("Method " + method + "of service " + service.getServiceName());
                URITemplate uriTemplate = new URITemplate();
                uriTemplate.setUriTemplate(protoFile.getPackageName() + "." + service.getServiceName());
                uriTemplate.setVerb(method);
                uriTemplates.add(uriTemplate);
            }
        }
        sortedServices.sort(String::compareTo);
        for (String service : sortedServices) {
            apiName.append(service).append("-");
        }

        api.setBasePath(protoFile.getBasePath());
        System.out.println("Basepath " + protoFile.getBasePath());
        api.setProtoDefinition(definition);
        api.setVersion(protoFile.getVersion());
        api.setName(apiName.toString());
        api.setUriTemplates(uriTemplates.toArray(new URITemplate[0]));
        return api;
    }

    ProtoFile getProtoFileFromDefinition(byte[] fileContent, String fileName) {
        Map<String, DescriptorProtos.FileDescriptorProto> protoMap = new HashMap<>();
        Map<String, Descriptors.FileDescriptor> descriptorMap = new HashMap<>();
        ArrayList<Descriptors.ServiceDescriptor> services = new ArrayList<>();
        String packageName = "";
        ProtoFile tempProtoFile = new ProtoFile();
        Map<String, Descriptors.FileDescriptor> wellKnownTypesMap = new HashMap<>();
        try {
            DescriptorProtos.FileDescriptorSet fileDescriptorSet = DescriptorProtos.FileDescriptorSet.parseFrom(
                    fileContent);

            populateProtoMap(protoMap, fileDescriptorSet);

            for (DescriptorProtos.FileDescriptorProto fileDescriptorProto : fileDescriptorSet.getFileList()) {
                packageName = processFileDescriptor(fileName, descriptorMap, protoMap, services, wellKnownTypesMap,
                        fileDescriptorProto);
            }

            tempProtoFile.setServices(convertServiceDescriptorsToServices(services, packageName));
            String[] info = packageName.split("\\.");
            if (info.length < 3) {
                throw new APIManagementException(
                        "Invalid package name: specify in the format of basepath.version.packageName");
            }
            tempProtoFile.setVersion(info[info.length - 2]);
            tempProtoFile.setPackageName(info[info.length - 1]);
            StringBuilder basePath = new StringBuilder("/").append(info[0]);
            for (int i = 1; i < info.length - 2; i++) {
                System.out.println(info[i]);
                basePath.append(".").append(info[i]);
            }
            tempProtoFile.setBasePath(basePath.toString());
            return tempProtoFile;
        } catch (Exception e) {
            e.printStackTrace();
            log.error("Proto definition validation failed for " + fileName + ": " + e.getMessage());
            return null;
        }
    }

    /**
     * @param fileName            - The name of the .desc file provided as input for the config generator
     * @param descriptorMap
     * @param protoMap
     * @param services
     * @param wellKnownTypesMap
     * @param fileDescriptorProto
     * @return
     * @throws Descriptors.DescriptorValidationException
     */
    private String processFileDescriptor(String fileName, Map<String, Descriptors.FileDescriptor> descriptorMap,
            Map<String, DescriptorProtos.FileDescriptorProto> protoMap,
            ArrayList<Descriptors.ServiceDescriptor> services,
            Map<String, Descriptors.FileDescriptor> wellKnownTypesMap,
            DescriptorProtos.FileDescriptorProto fileDescriptorProto) throws Descriptors.DescriptorValidationException {

        String packageName = fileDescriptorProto.getPackage();

        // Process and resolve dependencies for a given file descriptor
        Descriptors.FileDescriptor[] dependencies = fileDescriptorProto.getDependencyList().stream()
                .map(descriptorMap::get).toArray(Descriptors.FileDescriptor[]::new);

        // Build the file descriptor based on the proto and its dependencies
        Descriptors.FileDescriptor fileDescriptor = Descriptors.FileDescriptor.buildFrom(fileDescriptorProto,
                dependencies);
        services.addAll(fileDescriptor.getServices());
        descriptorMap.put(fileDescriptor.getName(), fileDescriptor);
        return packageName;
    }

    private Descriptors.FileDescriptor resolveDependency(Map<String, Descriptors.FileDescriptor> descriptorMap,
            Map<String, DescriptorProtos.FileDescriptorProto> protoMap,
            Map<String, Descriptors.FileDescriptor> wellKnownTypesMap, String descriptorName) {
        Descriptors.FileDescriptor dependency = descriptorMap.get(descriptorName);
        // Dependency has not been resolved yet
        if (dependency == null) {
            try {
                // if the dependency is a well known type
                if (descriptorName.startsWith("com.google.protobuf.")) {
                    dependency = resolveWellKnownType(descriptorName);
                    wellKnownTypesMap.put(descriptorName, dependency);
                } else {
                    // if the dependency is on another file that was imported, we resolve it and add it to the
                    // descriptor map
                    dependency = buildAndCacheDescriptor(descriptorName, protoMap, descriptorMap, wellKnownTypesMap);
                }
            } catch (Exception e) {
                System.err.println("Error loading well-known type: " + descriptorName + " - " + e.getMessage());
            }
        }
        if (dependency == null) {
            System.err.println("Missing dependency for " + descriptorName);
        }
        return dependency;
    }

    private Descriptors.FileDescriptor buildAndCacheDescriptor(String descriptorName,
            Map<String, DescriptorProtos.FileDescriptorProto> protoMap,
            Map<String, Descriptors.FileDescriptor> descriptorMap,
            Map<String, Descriptors.FileDescriptor> wellKnownTypesMap) {
        // this scenario is when you have an import in your proto file but that file hasnt been built yet
        // in that scenario, it needs to have its dependencies resolved as well
        DescriptorProtos.FileDescriptorProto dependencyProto = protoMap.get(descriptorName);
        if (dependencyProto != null) {
            //            Descriptors.FileDescriptor dependency = resolveDependency(descriptorMap, protoMap, wellKnownTypesMap,
            //                    descriptorName);
            //            descriptorMap.put(dependency.getName(), dependency);
            //            return dependency;
        }
        return null;
    }

    boolean validateProtoContent(byte[] fileContent, String fileName) {
        try {
            ProtoFile protoFile = getProtoFileFromDefinition(fileContent, fileName);
            return protoFile != null;
        } catch (Exception e) {
            log.error("Proto definition validation failed for " + fileName + ": " + e.getMessage());
            return false;
        }
    }

    public ArrayList<Service> convertServiceDescriptorsToServices(
            ArrayList<Descriptors.ServiceDescriptor> serviceDescriptors, String packageName) {
        ArrayList<Service> services = new ArrayList<>();
        for (Descriptors.ServiceDescriptor serviceDescriptor : serviceDescriptors) {
            List<Descriptors.MethodDescriptor> methodDescriptors = serviceDescriptor.getMethods();
            ArrayList<String> methods = new ArrayList<>();
            for (Descriptors.MethodDescriptor methodDescriptor : methodDescriptors) {
                methods.add(methodDescriptor.getName());
            }
            services.add(new Service(serviceDescriptor.getName(), methods, packageName));
        }
        return services;
    }

    public void validateGRPCAPIDefinition(byte[] inputByteArray, String fileName,
            APIDefinitionValidationResponse validationResponse, ArrayList<ErrorHandler> errors) {
        try {
            //            if (fileName.endsWith(".zip")) {
            //                try (ZipInputStream zis = new ZipInputStream(new ByteArrayInputStream(inputByteArray))) {
            //                    ZipEntry zipEntry;
            //                    while ((zipEntry = zis.getNextEntry()) != null) {
            //                        if (!zipEntry.isDirectory() && zipEntry.getName().endsWith(".desc")) {
            //                            byte[] protoFileContentBytes = zis.readAllBytes();
            //                            boolean validated = validateProtoContent(protoFileContentBytes, fileName);
            //                            if (!validated) {
            //                                throw new APIManagementException(
            //                                        "Invalid definition file provided. " + "Please provide a valid .zip or .proto file.");
            //                            }
            //                        }
            //                    }
            //                }
            //            } else
            if (fileName.endsWith(".desc")) {
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
}
