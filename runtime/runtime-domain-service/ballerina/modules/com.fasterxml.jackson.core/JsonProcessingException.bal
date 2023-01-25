// Ballerina error type for `com.fasterxml.jackson.core.JsonProcessingException`.

public const JSONPROCESSINGEXCEPTION = "JsonProcessingException";

type JsonProcessingExceptionData record {
    string message;
};

public type JsonProcessingException distinct error<JsonProcessingExceptionData>;

