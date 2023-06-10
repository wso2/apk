// Ballerina error type for `java.io.IOException`.

public const IOEXCEPTION = "IOException";

type IOExceptionData record {
    string message;
};

public type IOException distinct error<IOExceptionData>;

