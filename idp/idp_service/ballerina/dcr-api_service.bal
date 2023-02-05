import ballerina/http;

service /dcr on ep0{
    resource function post register(@http:Payload RegistrationRequest payload) returns CreatedApplication|BadRequestClientRegistrationError|ConflictClientRegistrationError|InternalServerErrorClientRegistrationError {
        DCRMClient dcrmClient = new;
        return dcrmClient.createDCRApplication(payload);

    }
    resource function get register/[string client_id]() returns Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError {
        DCRMClient dcrmClient = new;
        return dcrmClient.getApplication(client_id);

    }
    resource function put register/[string client_id](@http:Payload UpdateRequest payload) returns Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError|BadRequestClientRegistrationError {
        DCRMClient dcrmClient = new;
        return dcrmClient.updateDCRApplication(client_id, payload);

    }
    resource function delete register/[string client_id]() returns http:NoContent|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError {
        DCRMClient dcrmClient = new;
        return dcrmClient.deleteApplication(client_id);

    }
}
