package org.wso2.apk.integration.utils;

public class Constants {

    public static final String DEFAULT_IDP_HOST = "idp.am.wso2.com";
    public static final String DEFAULT_API_HOST = "api.am.wso2.com";
    public static final String DEFAULT_GW_PORT = "9095";
    public static final String DEFAULT_TOKEN_EP = "oauth2/token";
    public static final String DEFAULT_API_CONFIGURATOR = "api/configurator/1.0.0/";
    public static final String DEFAULT_API_DEPLOYER = "api/deployer/1.0.0/";

    public class REQUEST_HEADERS {
        public static final String HOST = "Host";
        public static final String AUTHORIZATION = "Authorization";
        public static final String CONTENT_TYPE = "Content-Type";
    }

    public class CONTENT_TYPES {
        public static final String APPLICATION_JSON = "application/json";
        public static final String APPLICATION_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded";
    }
}
