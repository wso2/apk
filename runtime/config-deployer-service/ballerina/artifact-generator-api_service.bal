import ballerina/http;
import wso2/apk_common_lib as commons;


isolated service /api/configurator on ep0 {
    # Create API configuration file from api specification.
    #
    # + request - parameter description 
    # + return - returns can be any of following types
    # OkAnydata (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/'generate\-configuration(http:Request request) returns OkAnydata|BadRequestError|InternalServerErrorError|commons:APKError {
        ConfigGeneratorClient apiclient = new ;
        return check apiclient.getGeneratedAPKConf(request);
    }
    # Generate K8s Resources
    #
    # + organization - **Organization ID** of the organization the API belongs to. 
    # + request - parameter description 
    # + return - returns can be any of following types
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/'generate\-k8s\-resources(string? organization, http:Request request) returns http:Response|BadRequestError|InternalServerErrorError|commons:APKError {
        ConfigGeneratorClient apiclient = new ;
        commons:Organization organizationObj  = {displayName: "default",
        name: "default",
        organizationClaimValue: "default",
        uuid: "",
        enabled: true};
        if (organization is string) {
            organizationObj.name = organization;
        }
        return check apiclient.getGeneratedK8sResources(request,organizationObj);
    }
}
