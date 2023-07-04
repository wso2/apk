// Ballerina error type for `java.lang.Exception`.

public const EXCEPTION = "Exception";

type ExceptionData record {
    string message;
};

public type Exception distinct error<ExceptionData>;

