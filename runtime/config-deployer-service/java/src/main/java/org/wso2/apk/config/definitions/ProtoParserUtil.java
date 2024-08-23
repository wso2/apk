package org.wso2.apk.config.definitions;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.config.api.APIDefinition;
import org.wso2.apk.config.api.APIDefinitionValidationResponse;
import org.wso2.apk.config.api.ErrorHandler;
import org.wso2.apk.config.api.ErrorItem;
import org.wso2.apk.config.api.ExceptionCodes;

public class ProtoParserUtil {
    /**
     * Provide common functions related to OAS
     */
    private static final Log log = LogFactory.getLog(ProtoParserUtil.class);
    private static final ProtoParser protoParser = new ProtoParser();

    /**
     * Validate graphQL Schema
     * 
     * @return Validation response
     */
    public static APIDefinitionValidationResponse validateGRPCAPIDefinition(String apiDefinition,
            boolean returnGRPCContent) {
        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        ArrayList<ErrorHandler> errors = new ArrayList<>();

        if (apiDefinition.isBlank()) {
            validationResponse.setValid(false);
            errors.add(ExceptionCodes.GRPC_PROTO_DEFINTION_CANNOT_BE_NULL);
            validationResponse.setErrorItems(errors);
        }

        try {
            boolean validated = protoParser.validateProtoFile(apiDefinition);
            validationResponse.setValid(validated);
            validationResponse.setContent(apiDefinition);
        } catch (Exception e) {
            ProtoParserUtil.addErrorToValidationResponse(validationResponse, e.getMessage());
            validationResponse.setValid(false);
            errors.add(new ErrorItem("API Definition Validation Error", "API Definition is invalid", 400, 400));
            validationResponse.setErrorItems(errors);
        }

        return validationResponse;
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
        errorItem.setErrorCode(ExceptionCodes.PROTO_DEFINITION_PARSE_EXCEPTION.getErrorCode());
        errorItem.setMessage(ExceptionCodes.PROTO_DEFINITION_PARSE_EXCEPTION.getErrorMessage());
        errorItem.setDescription(errMessage);
        validationResponse.getErrorItems().add(errorItem);
        return errorItem;
    }

}
