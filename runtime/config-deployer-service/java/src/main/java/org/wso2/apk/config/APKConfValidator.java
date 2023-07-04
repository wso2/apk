package org.wso2.apk.config;

import org.everit.json.schema.Schema;
import org.everit.json.schema.ValidationException;
import org.everit.json.schema.loader.SchemaLoader;
import org.wso2.apk.config.api.ErrorItem;
import org.wso2.apk.config.api.APKConfValidationResponse;

import java.util.ArrayList;
import java.util.List;

/**
 * This class used to validate apk-conf schema.
 */
public class APKConfValidator {

    Schema schema;

    public APKConfValidator(String jsonSchema) throws Exception {

        org.json.JSONObject apkConfSchema = new org.json.JSONObject(jsonSchema);
        schema = SchemaLoader.load(apkConfSchema);
    }

    /**
     * This method use to validate apk-conf content with schema.
     *
     * @param apkConfString apk-conf as String.
     * @return apkConfValidationResponse based on result.
     * @throws Exception if error occurred at validation.
     */
    public APKConfValidationResponse validate(String apkConfString) throws Exception {

        org.json.JSONObject apkConfJson = new org.json.JSONObject(apkConfString);
        try {
            schema.validate(apkConfJson);
            return new APKConfValidationResponse(true);
        } catch (ValidationException e) {
            APKConfValidationResponse apkConfValidationResponse = new APKConfValidationResponse(false);
            // Data is invalid, handle the validation errors
            List<ErrorItem> errorItems = new ArrayList<>();
            for (ValidationException message : e.getCausingExceptions()) {
                ErrorItem errorItem = new ErrorItem();
                errorItem.setMessage(message.getErrorMessage());
                errorItem.setDescription(message.getLocalizedMessage());
                errorItems.add(errorItem);
            }
            if (e.getCausingExceptions() == null || e.getCausingExceptions().isEmpty()){
                ErrorItem errorItem = new ErrorItem();
                errorItem.setMessage(e.getErrorMessage());
                errorItem.setDescription(e.getLocalizedMessage());
                errorItems.add(errorItem);
            }
            apkConfValidationResponse.setErrorItems(errorItems.toArray(new ErrorItem[errorItems.size()]));
            return apkConfValidationResponse;
        }
    }
}
