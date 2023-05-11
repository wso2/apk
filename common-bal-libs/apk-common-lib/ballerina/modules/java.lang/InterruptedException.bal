// Ballerina error type for `java.lang.InterruptedException`.

public const INTERRUPTEDEXCEPTION = "InterruptedException";

type InterruptedExceptionData record {
    string message;
};

public type InterruptedException distinct error<InterruptedExceptionData>;

