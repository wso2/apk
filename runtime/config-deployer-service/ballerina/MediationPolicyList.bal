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

final MediationPolicy[] & readonly avilableMediationPolicyList = [
    {
        id: "1",
        'type: MEDIATION_POLICY_TYPE_REQUEST_HEADER_MODIFIER,
        name: MEDIATION_POLICY_NAME_ADD_HEADER,
        displayName: "Add Header",
        description: "This policy allows you to add a new header to the request",
        applicableFlows: [MEDIATION_POLICY_FLOW_REQUEST],
        supportedApiTypes: [API_TYPE_REST],
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
        name: MEDIATION_POLICY_NAME_REMOVE_HEADER,
        displayName: "Remove Header",
        description: "This policy allows you to remove a header from the request",
        applicableFlows: [MEDIATION_POLICY_FLOW_REQUEST],
        supportedApiTypes: [API_TYPE_REST],
        policyAttributes: [
            {
                name: "headerName",
                description: "Name of the header to be removed",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    },
    {
        id: "3",
        'type: MEDIATION_POLICY_TYPE_RESPONSE_HEADER_MODIFIER,
        name: MEDIATION_POLICY_NAME_ADD_HEADER,
        displayName: "Add Header",
        description: "This policy allows you to add a new header to the response",
        applicableFlows: [MEDIATION_POLICY_FLOW_RESPONSE],
        supportedApiTypes: [API_TYPE_REST],
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
        id: "4",
        'type: MEDIATION_POLICY_TYPE_RESPONSE_HEADER_MODIFIER,
        name: MEDIATION_POLICY_NAME_REMOVE_HEADER,
        displayName: "Remove Header",
        description: "This policy allows you to remove a header from the response",
        applicableFlows: [MEDIATION_POLICY_FLOW_RESPONSE],
        supportedApiTypes: [API_TYPE_REST],
        policyAttributes: [
            {
                name: "headerName",
                description: "Name of the header to be removed",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    },
    {
        id: "5",
        'type: MEDIATION_POLICY_TYPE_INTERCEPTOR,
        name: MEDIATION_POLICY_TYPE_INTERCEPTOR,
        displayName: "Interceptor",
        description: "This policy allows you to engage an interceptor service",
        applicableFlows: [MEDIATION_POLICY_FLOW_REQUEST, MEDIATION_POLICY_FLOW_RESPONSE],
        supportedApiTypes: [API_TYPE_REST],
        policyAttributes: [
            {
                name: "headersEnabled",
                description: "Indicates whether request/response header details should be sent to the interceptor service",
                'type: "boolean",
                required: false
            },
            {
                name: "bodyEnabled",
                description: "Indicates whether request/response body details should be sent to the interceptor service",
                'type: "boolean",
                required: false
            },
            {
                name: "contextEnabled",
                description: "Indicates whether context details should be sent to the interceptor service",
                'type: "boolean",
                required: false
            },
            {
                name: "trailersEnabled",
                description: "Indicates whether request/response trailer details should be sent to the interceptor service",
                'type: "boolean",
                required: false
            },
            {
                name: "backendUrl",
                description: "Backend URL of the interceptor service",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    },
    {
        id: "6",
        'type: POLICY_TYPE_BACKEND_JWT,
        name: POLICY_TYPE_BACKEND_JWT,
        displayName: "BackendJwt",
        description: "This policy allows you to add backend JWT",
        applicableFlows: [MEDIATION_POLICY_FLOW_REQUEST],
        supportedApiTypes: [API_TYPE_REST],
        policyAttributes: [
            {
                name: "enabled",
                description: "enabled holds the status of the policy",
                'type: "boolean",
                required: true
            },
            {
                name: "encoding",
                description: "Encoding holds the encoding type",
                'type: "String",
                required: false
            },
            {
                name: "signingAlgorithm",
                description: "signingAlgorithm holds the signing algorithm",
                'type: "String",
                required: false
            },
            {
                name: "header",
                description: "Header holds the header name",
                'type: "String",
                required: false
            },
            {
                name: "tokenTTL",
                description: "TokenTTL holds the token time to live in seconds",
                'type: "int",
                required: false
            },
            {
                name: "customClaims",
                description: "CustomClaim holds custom claim information",
                'type: "array",
                required: false
            }
        ]
    }
];



public type MediationPolicy record {|
    string id;
    string 'type;
    string name;
    string displayName?;
    string description?;
    string[] applicableFlows?;
    string[] supportedApiTypes?;
    MediationPolicySpecAttribute[] policyAttributes?;
|};

public type MediationPolicySpecAttribute record {|
    # Name of the attribute
    string name?;
    # Description of the attribute
    string description?;
    # Is this option mandatory for the policy
    boolean required?;
    # UI validation regex for the attribute
    string validationRegex?;
    # Type of the attribute
    string 'type?;
    # Default value for the attribute
    string defaultValue?;
|};

public type MediationPolicyList record {|
    # Number of mediation policies returned.
    int count?;
    MediationPolicy[] list?;
|};
