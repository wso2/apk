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

public type BadRequestTokenErrorResponse record {|
    *http:BadRequest;
    TokenErrorResponse body;
|};

public type UnauthorizedTokenErrorResponse record {|
    *http:Unauthorized;
    TokenErrorResponse body;
|};

public type TokenErrorResponse record {
    # Error code classifying the type of preProcessingError.
    string preProcessingError;
};

public type TokenResponse record {
    # OAuth access tokn issues by authorization server.
    string access_token;
    # The type of the token issued.
    string token_type;
    # The lifetime in seconds of the access token.
    int expires_in?;
    # OPTIONAL.
    # The refresh token, which can be used to obtain new access tokens.
    string refresh_token?;
    # The scope of the access token requested.
    string scope?;
};

public type Token_body record {
    # Required OAuth grant type
    string grant_type;
    # Authorization code to be sent for authorization grant type
    string code?;
    # Clients redirection endpoint
    string redirect_uri?;
    # OAuth client identifier
    string client_id?;
    # OAuth client secret
    string client_secret?;
    # Refresh token issued to the client.
    string refresh_token?;
    # OAuth scopes
    string scope?;
    # username
    string username?;
    # password
    string password?;
    # Validity period of token
    int validity_period?;
};
