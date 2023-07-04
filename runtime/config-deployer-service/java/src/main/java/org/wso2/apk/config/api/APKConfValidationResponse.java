package org.wso2.apk.config.api;

import org.wso2.apk.config.api.ErrorHandler;

import java.util.ArrayList;
import java.util.List;

public class APKConfValidationResponse {

    boolean validated = false;
    private ErrorHandler[] errorItems;

    public APKConfValidationResponse(boolean validated) {

        this.validated = validated;
    }

    public boolean isValidated() {

        return validated;
    }

    public void setValidated(boolean validated) {

        this.validated = validated;
    }

    public ErrorHandler[] getErrorItems() {

        return errorItems;
    }

    public void setErrorItems(ErrorHandler[] errorItems) {

        this.errorItems = errorItems;
    }
}
