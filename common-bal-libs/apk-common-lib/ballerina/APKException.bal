// Ballerina error type for `org.wso2.apk.common.APKException`.

public const APKEXCEPTION = "APKException";

type APKExceptionData record {
    string message;
};

public type APKException distinct error<APKExceptionData>;

