// Ballerina error type for `org.wso2.apk.config.api.APIManagementException`.

public const APIMANAGEMENTEXCEPTION = "APIManagementException";

type APIManagementExceptionData record {
    string message;
};

public type APIManagementException distinct error<APIManagementExceptionData>;

