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
# Description
#
# + issuer - issuer of the JWT token  
# + jwksUrl - URL of the JWKS endpoint  
# + organizationClaim - organization claim of the JWT token  
# + userClaim - user claim of the JWT token
# + authorizationHeader - authorization header of the JWT token  
# + publicKey - public key of the JWT token
public type IDPConfiguration record {|
string issuer = "wso2.org/products/am";
string jwksUrl? ;
string organizationClaim ="x-wso2-organization";
string userClaim = "sub";
string authorizationHeader = "X-JWT-Assertion";
KeyStore publicKey;
|};

