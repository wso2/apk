//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
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
