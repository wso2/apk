import wso2/apk_common_lib as commons;

public isolated function e909100(string id) returns commons:APKError {
    return error commons:APKError( id + " not found",
        code = 909100,
        message = id + " not found",
        statusCode = 404,
        description = id + " not found"
    ); 
}

public isolated function e90911(string apiName) returns commons:APKError {
    return error commons:APKError( "API Name - " + apiName + " already exist",
        code = 90911,
        message = "API Name - " + apiName + " already exist",
        statusCode = 409,
        description = "API Name - " + apiName + " already exist"
    ); 
}

public isolated function e90912(string apiContext) returns commons:APKError {
    return error commons:APKError( "API Context - " + apiContext + " already exist",
        code = 90912,
        message = "API Context - " + apiContext + " already exist",
        statusCode = 409,
        description = "API Context - " + apiContext + " already exist"
    ); 
}

public isolated function e90913() returns commons:APKError {
    return error commons:APKError( "Sandbox endpoint not specified",
        code = 90913,
        message = "Sandbox endpoint not specified",
        statusCode = 400,
        description = "Sandbox endpoint not specified"
    ); 
}

public isolated function e90914() returns commons:APKError {
    return error commons:APKError( "Production endpoint not specified",
        code = 90914,
        message = "Production endpoint not specified",
        statusCode = 400,
        description = "Production endpoint not specified"
    ); 
}

public isolated function e90915(string apiContext) returns commons:APKError {
    return error commons:APKError("API context " + apiContext + " invalid",
        code = 90915,
        message = "API context " + apiContext + " invalid",
        statusCode = 400,
        description = "API context " + apiContext + " invalid"
    ); 
}

public isolated function e90916(string apiName) returns commons:APKError {
    return error commons:APKError("API name " + apiName + " invalid",
        code = 90916,
        message = "API name " + apiName + " invalid",
        statusCode = 400,
        description = "API name " + apiName + " invalid"
    ); 
}

public isolated function e90917() returns commons:APKError {
    return error commons:APKError("Invalid API request",
        code = 90917,
        message = "Invalid API request",
        statusCode = 400,
        description = "Invalid API request"
    ); 
}

public isolated function e90918() returns commons:APKError {
    return error commons:APKError("Error while generating token",
        code = 90918,
        message = "Error while generating token",
        statusCode = 400,
        description = "Error while generating token"
    ); 
}