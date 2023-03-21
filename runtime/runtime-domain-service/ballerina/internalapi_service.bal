import wso2/apk_common_lib as commons;
import ballerina/http;

http:Service internalRuntimeService = service object {
    isolated resource function get apis/[string apiId]/definition(http:Request request) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        string|http:HeaderNotFoundError header = request.getHeader("X-WSO2-Organization");
        if header is string {
            APIClient apiClient = new ();
            commons:Organization organization = {
                displayName: header,
                name: header,
                organizationClaimValue: header,
                uuid: header,
                enabled: true
            };
            return apiClient.getAPIDefinitionByID(apiId, organization);
        }else {
            PreconditionFailedError preconditionFailedError = {body: {code: 900901, message: "X-WSO2-Organization is missing in request"}};
            return preconditionFailedError;
        }
    }
};
