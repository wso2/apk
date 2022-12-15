package org.wso2.apk.runtime;

import org.wso2.apk.runtime.api.APIDefinition;
import org.wso2.apk.runtime.api.APIDefinitionValidationResponse;
import org.wso2.apk.runtime.api.APIManagementException;
import org.wso2.apk.runtime.api.ExceptionCodes;
import org.wso2.apk.runtime.definitions.OASParserUtil;
import org.wso2.apk.runtime.model.API;
import org.wso2.apk.runtime.model.URITemplate;

import java.nio.charset.StandardCharsets;
import java.util.Set;


public class RuntimeAPICommonUtil {

    public static String generateDefinition(API api) throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getParser(api);
        return parser.generateAPIDefinition(api);
    }

    /**
     * @param inputByteArray OpenAPI definition file
     * @param apiDefinition  OpenAPI definition
     * @param fileName       Filename of the definition file
     * @param returnContent  Whether to return json or not
     * @return APIDefinitionValidationResponse
     * @throws APIManagementException when file parsing fails
     */
    public static APIDefinitionValidationResponse validateOpenAPIDefinition(String type, byte[] inputByteArray,
                                                                            String apiDefinition, String fileName,
                                                                            boolean returnContent)
            throws APIManagementException {
        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        if (APIConstants.ParserType.OAS3.name().equals(type)){
            if (inputByteArray != null && inputByteArray.length>0) {
                if (fileName != null) {
                    if (fileName.endsWith(".zip")) {
                        validationResponse =
                                OASParserUtil.extractAndValidateOpenAPIArchive(inputByteArray, returnContent);
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
        }
        return validationResponse;
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

    private RuntimeAPICommonUtil() {
    }

}
