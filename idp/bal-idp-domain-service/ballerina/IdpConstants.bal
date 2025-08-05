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
const string CLIENT_NAME_EMPTY_ERROR = "100150";
const string GRANT_TYPES_EMPTY_ERROR = "100151";
const string INTERNAL_ERROR = "100100";
const string INVALID_USERNAME_OR_PASSWORD = "100171";
const string INVALID_PERMISSION = "100172";
const string INVALID_ORGANIZATION = "100173";
const string CLIENT_ID_NOT_FOUND_ERROR = "100152";
const string UNSUPPORTED_GRANT_TYPE_ERROR = "100153";
const string CLIENT_CREDENTIALS_GRANT_TYPE = "client_credentials";
const string AUTHORIZATION_CODE_GRANT_TYPE = "authorization_code";
const string ACCESS_TOKEN_TYPE = "access_token";
const string REFRESH_TOKEN_TYPE = "refresh_token";
const string REFRESH_TOKEN_GRANT_TYPE = "refresh_token";
const string AUTHORIZATION_CODE_TYPE = "authorization_code";
const string SESSION_KEY_TYPE = "session_key";
const string ORGANIZATION_CLAIM = "organization";
const string TOKEN_TYPE_CLAIM = "token_type";
const string TOKEN_TYPE_BEARER = "Bearer";
const string REDIRECT_URI_CLAIM = "redirectUri";
const string SCOPES_CLAIM = "scope";
const string SESSION_KEY_PREFIX = "session-";
const string STATE_KEY_QUERY_PARAM = "stateKey";
const string LOCATION_HEADER = "Location";
const string CLIENT_ID_CLAIM  = "clientId";
const string AUTHORIZATION_CODE_QUERY_PARAM = "code";
const string TYPE_CLAIM = "type";
isolated  string[] ALLOWED_GRANT_TYPES = [CLIENT_CREDENTIALS_GRANT_TYPE,REFRESH_TOKEN_GRANT_TYPE,AUTHORIZATION_CODE_GRANT_TYPE];