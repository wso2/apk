//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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

import ballerina/sql;

sql:ParameterizedQuery ADD_API_Suffix = `INSERT INTO api(api_uuid, api_name, api_version,context,api_provider,status,artifact) VALUES (`;
sql:ParameterizedQuery UPDATE_API_Suffix = `UPDATE api SET`;
sql:ParameterizedQuery DELETE_API_Suffix = `DELETE FROM api WHERE api_uuid = `;
