import ballerina/lang.value;
import ballerina/http;
public isolated function getAuthenticatedUserContext(http:RequestContext requestContext) returns UserContext|APKError {
    value:Cloneable|object {} userContextAttribute = requestContext.get(VALIDATED_USER_CONTEXT);
    UserContext|error userContext = userContextAttribute.ensureType(UserContext);
    if userContext is UserContext {
        return userContext;
    } else {
        APKError apkError = error("unauthorized Request", message = "Invalid Credentials",code = 900905,description = "Invalid Credentials",statusCode = 401);
        return apkError;
    }
}
    
