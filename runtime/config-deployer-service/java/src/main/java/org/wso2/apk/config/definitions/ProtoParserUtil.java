package org.wso2.apk.config.definitions;

import java.util.*;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
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
    public static APIDefinitionValidationResponse validateGRPCAPIDefinition(byte[] inputByteArray, String fileName,
            boolean returnGRPCContent) {
        APIDefinitionValidationResponse validationResponse = new APIDefinitionValidationResponse();
        ArrayList<ErrorHandler> errors = new ArrayList<>();

        if (inputByteArray.length == 0) {
            validationResponse.setValid(false);
            errors.add(ExceptionCodes.GRPC_PROTO_DEFINTION_CANNOT_BE_NULL);
            validationResponse.setErrorItems(errors);
        } else {
            protoParser.validateGRPCAPIDefinition(inputByteArray, fileName, validationResponse, errors);
        }
        return validationResponse;
    }

    /**
     * Add error item with the provided message to the provided validation response object
     *
     * @param validationResponse APIDefinitionValidationResponse object
     * @param errMessage         error message
     */
    public static void addErrorToValidationResponse(APIDefinitionValidationResponse validationResponse,
            String errMessage) {
        ErrorItem errorItem = new ErrorItem();
        errorItem.setErrorCode(ExceptionCodes.PROTO_DEFINITION_PARSE_EXCEPTION.getErrorCode());
        errorItem.setMessage(ExceptionCodes.PROTO_DEFINITION_PARSE_EXCEPTION.getErrorMessage());
        errorItem.setDescription(errMessage);
        validationResponse.getErrorItems().add(errorItem);
    }

}
