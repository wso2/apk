package org.wso2.apk.runtime.definitions;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.runtime.api.*;

import java.util.List;

public class AsyncApiParserUtil {

    private static final APIDefinition asyncApiParser = new AsyncApiParser();
    private static final Log log = LogFactory.getLog(AsyncApiParserUtil.class);
    private static final String PATH_SEPARATOR = "/";

    public static APIDefinitionValidationResponse validateAsyncAPISpecification(
            String schemaToBeValidated, boolean returnJSONContent) throws APIManagementException {

        APIDefinitionValidationResponse validationResponse = asyncApiParser.validateAPIDefinition(schemaToBeValidated, returnJSONContent);
        final String asyncAPIKeyNotFound = "#: required key [asyncapi] not found";

        if (!validationResponse.isValid()) {
            for (ErrorHandler errorItem : validationResponse.getErrorItems()) {
                if (asyncAPIKeyNotFound.equals(errorItem.getErrorMessage())) {    //change it other way
                    addErrorToValidationResponse(validationResponse, "#: attribute [asyncapi] should be present");
                    return validationResponse;
                }
            }
        }

        return validationResponse;
    }

    public static void updateValidationResponseAsSuccess(
            APIDefinitionValidationResponse validationResponse,
            String originalAPIDefinition,
            String asyncAPIVersion,
            String title,
            String version,
            String context,
            String description,
            List<String> endpoints
    ) {
        validationResponse.setValid(true);
        validationResponse.setContent(originalAPIDefinition);
        Info info = new Info();
        info.setOpenAPIVersion(asyncAPIVersion);
        info.setName(title);
        info.setVersion(version);
        info.setContext(context);
        info.setDescription(description);
        info.setEndpoints(endpoints);
        validationResponse.setInfo(info);
    }

    public static ErrorItem addErrorToValidationResponse(
            APIDefinitionValidationResponse validationResponse, String errMessage) {
        ErrorItem errorItem = new ErrorItem();
        errorItem.setMessage(errMessage);
        validationResponse.getErrorItems().add(errorItem);
        return errorItem;
    }
}
