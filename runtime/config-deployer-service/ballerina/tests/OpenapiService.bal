import ballerina/io;
import ballerina/http;

service /test on new http:Listener(9090) {

    # Create API configuration file from api specification.
    #
    # + request - parameter description 
    # + return - returns can be any of following types
    # OkAnydata (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get definition(http:Request request) returns http:Response|error {
        string fileReadString = check io:fileReadString("./tests/resources/api.yaml");
        http:Response response = new;
        response.setPayload(fileReadString);
        response.setHeader("Content-Type", "application/yaml");
        return response;
    }
}