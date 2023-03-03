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
import runtime_domain_service.model;



enum PolicyFlows {
    request,
    response
}


isolated model:MediationPolicy[] avilableMediationPolicyList = [
    {
        id: "1",
        'type: MEDIATION_POLICY_TYPE_REQUEST_HEADER_MODIFIER,
        name: "addHeader",
        displayName: "Add Header",
        description: "This policy allows you to add a new header to the request",
        applicableFlows: [request],
        supportedApiTypes: [API_TYPE_REST],
        isApplicableforAPILevel: true,
        isApplicableforOperationLevel: true,
        policyAttributes: [
            {
                name: "headerName",
                description: "Name of the header to be added",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            },
            {
                name: "headerValue",
                description: "Value of the header",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    },
        {
        id: "2",
        'type: MEDIATION_POLICY_TYPE_REQUEST_HEADER_MODIFIER,
        name: "removeHeader",
        displayName: "Remove Header",
        description: "This policy allows you to remove a header from the request",
        applicableFlows: [request],
        supportedApiTypes: [API_TYPE_REST],
        isApplicableforAPILevel: true,
        isApplicableforOperationLevel: true,
        policyAttributes: [
            {
                name: "headerName",
                description: "Name of the header to be removed",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    }
];
