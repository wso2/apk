// Ballerina error type for `org.wso2.apk.devportal.sdk.APIClientGenerationException`.

public const APICLIENTGENERATIONEXCEPTION = "APIClientGenerationException";

type APIClientGenerationExceptionData record {
    string message;
};

public type APIClientGenerationException distinct error<APIClientGenerationExceptionData>;

