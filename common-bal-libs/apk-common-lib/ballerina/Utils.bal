import ballerina/lang.value;
import apk_common_lib.java.lang;
import ballerina/http;

public isolated function getAuthenticatedUserContext(http:RequestContext requestContext) returns UserContext|APKError {
    value:Cloneable|object {} userContextAttribute = requestContext.get(VALIDATED_USER_CONTEXT);
    UserContext|error userContext = userContextAttribute.ensureType(UserContext);
    if userContext is UserContext {
        return userContext;
    } else {
        APKError apkError = error("unauthorized Request", message = "Invalid Credentials", code = 900905, description = "Invalid Credentials", statusCode = 401);
        return apkError;
    }
}

public isolated function fromJsonStringToYaml(string jsonString) returns string?|error {
    YamlUtil yamlUtil = newYamlUtil1();
    string?|lang:Exception convertedString = yamlUtil.fromJsonStringToYaml(jsonString);
    if convertedString is string {
        return convertedString;
    } else if convertedString is () {
        return convertedString;
    } else {
        return convertedString.cause();
    }
}

public isolated function fromYamlStringToJson(string yamlString) returns json?|error {
    YamlUtil yamlUtil = newYamlUtil1();
    string?|lang:Exception convertedString = yamlUtil.fromYamlStringToJson(yamlString);
    if convertedString is string {
        return check value:fromJsonString(convertedString);
    } else if convertedString is () {
        return convertedString;
    } else {
        return convertedString.cause();
    }
}

