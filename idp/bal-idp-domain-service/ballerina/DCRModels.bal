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

public type NotFoundClientRegistrationError record {|
    *http:NotFound;
    ClientRegistrationError body;
|};

public type InternalServerErrorClientRegistrationError record {|
    *http:InternalServerError;
    ClientRegistrationError body;
|};

public type ConflictClientRegistrationError record {|
    *http:Conflict;
    ClientRegistrationError body;
|};

public type BadRequestClientRegistrationError record {|
    *http:BadRequest;
    ClientRegistrationError body;
|};

public type CreatedApplication record {|
    *http:Created;
    Application body;
|};

public type ClientRegistrationError record {
    string 'error?;
    string error_description?;
};

public type UpdateRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};



public type RegistrationRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};

public type Application record {
    string client_id?;
    string client_secret?;
    string[] redirect_uris?;
    string[] grant_types?;
    string client_name?;
    int client_secret_expires_at?;
};
