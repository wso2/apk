public type ErrorHandler record {|
    int code;
    string message;
    string statusCode;
    string description;
    map<string> moreInfo = {};
|};
public type APKError distinct (error<ErrorHandler>);

