package org.wso2.apk.common;

import java.io.IOException;

public class APKException extends Exception {

    public APKException(String message) {

        super(message);
    }

    public APKException(String message, IOException e) {

        super(message, e);

    }
}
