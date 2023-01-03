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

sql:ParameterizedQuery GET_API = `SELECT API_UUID AS ID, API_ID as APIID,
        API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, ARTIFACT as ARTIFACT
        FROM API `;
sql:ParameterizedQuery GET_API_LifeCycle_Prefix = `SELECT status from api where api_uuid = `;
sql:ParameterizedQuery UPDATE_API_LifeCycle_Prefix = `UPDATE api SET status = `;
sql:ParameterizedQuery ADD_LC_EVENT_Prefix = `INSERT INTO api_lc_event (api_id,previous_state,new_state,user_id,organization,event_date) VALUES (`;
sql:ParameterizedQuery SUBSCRIPTION_BLOCK_Prefix = `UPDATE subscription set sub_status = `; 
sql:ParameterizedQuery SUBSCRIPTION_UNBLOCK_Prefix = `UPDATE subscription set sub_status = 'UNBLOCKED'`; 
