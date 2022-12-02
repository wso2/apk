package org.wso2.apk.runtime;

import org.wso2.apk.apimgt.api.APIDefinition;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.APIIdentifier;
import org.wso2.apk.runtime.model.API;

public class RuntimeAPICommonUtil {

    public static String generateDefinition(API api) throws APIManagementException {
        APIDefinition parser = DefinitionParserFactory.getParser(api);
        org.wso2.apk.apimgt.api.model.API apiModel = retrieveAPIModel(api);
        return parser.generateAPIDefinition(apiModel);
    }

    private static org.wso2.apk.apimgt.api.model.API retrieveAPIModel(API api) {
        org.wso2.apk.apimgt.api.model.API model = new org.wso2.apk.apimgt.api.model.API(new APIIdentifier("apkuser", api.getName(), api.getVersion()));
        model.setType(api.getType());
        model.setUriTemplates(api.getUriTemplates());
        return model;
    }
}
