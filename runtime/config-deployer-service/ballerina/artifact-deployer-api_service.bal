import ballerina/http;
import wso2/apk_common_lib as commons;

isolated service /api/deployer on ep0 {
    # Deploy API
    #
    # + request - parameter description 
    # + return - returns can be any of following types
    # anydata (API deployed successfully)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/deploy(http:Request request) returns commons:APKError|http:Response {
        DeployerClient deployerClient = new;
        return check deployerClient.handleAPIDeployment(request);
    }
    # Undeploy API
    #
    # + apiId - UUID of the K8s API Resource 
    # + return - returns can be any of following types
    # AcceptedString (API undeployed successfully)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/undeploy(string apiId) returns AcceptedString|BadRequestError|InternalServerErrorError|commons:APKError|error {
        DeployerClient deployerClient = new;
        return check deployerClient.handleAPIUndeployment(apiId);
    }
}
