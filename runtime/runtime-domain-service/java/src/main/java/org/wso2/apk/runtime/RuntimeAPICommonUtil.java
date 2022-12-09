package org.wso2.apk.runtime;

import org.wso2.apk.runtime.api.APIDefinition;
import org.wso2.apk.runtime.api.APIDefinitionValidationResponse;
import org.wso2.apk.runtime.api.APIManagementException;
import org.wso2.apk.runtime.model.API;
import org.wso2.apk.runtime.model.URITemplate;

import java.util.Set;


public class RuntimeAPICommonUtil {

    public static String generateDefinition(API api) throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getParser(api);
        return parser.generateAPIDefinition(api);
    }

    public static APIDefinitionValidationResponse validateDefinition(String content, boolean returnContent)
            throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getValidatedParser(content);
        return parser.validateAPIDefinition(content, returnContent);
    }

    public static Set<URITemplate> generateUriTemplatesFromAPIDefinition(String apiType, String content)
            throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getParser(apiType);
        return parser.getURITemplates(content);
    }

    private RuntimeAPICommonUtil() {
    }

}
