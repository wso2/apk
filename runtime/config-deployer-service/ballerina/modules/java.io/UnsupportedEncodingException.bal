// Ballerina error type for `java.io.UnsupportedEncodingException`.

public const UNSUPPORTEDENCODINGEXCEPTION = "UnsupportedEncodingException";

type UnsupportedEncodingExceptionData record {
    string message;
};

public type UnsupportedEncodingException distinct error<UnsupportedEncodingExceptionData>;

