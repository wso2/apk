import ballerina/http;
import ballerinax/prometheus as _;

import wso2/apk_common_lib as commons;

isolated service http:InterceptableService /api/deployer on ep0 {
    # Deploy API
    #
    # + request - parameter description 
    # + return - returns can be any of following types
    # anydata (API deployed successfully)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/deploy(http:RequestContext requestContext, http:Request request) returns commons:APKError|http:Response {
        DeployerClient deployerClient = new;
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return check deployerClient.handleAPIDeployment(request, organization);
    }
    # Undeploy API
    #
    # + apiId - UUID of the K8s API Resource 
    # + return - returns can be any of following types
    # AcceptedString (API undeployed successfully)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/undeploy(http:RequestContext requestContext, string apiId) returns AcceptedString|BadRequestError|InternalServerErrorError|commons:APKError {
        DeployerClient deployerClient = new;
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return check deployerClient.handleAPIUndeployment(apiId, organization);
    }

    public function createInterceptors() returns http:Interceptor|http:Interceptor[] {
        http:Interceptor[] interceptors = [jwtValidationInterceptor, requestErrorInterceptor, responseErrorInterceptor];
        return interceptors;
    }
}
