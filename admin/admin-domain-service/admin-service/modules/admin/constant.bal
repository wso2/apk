Settings settingPayload = {
    "scopes": [
        "string"
    ],
    "keyManagerConfiguration": [
        {
            "type": "default",
            "displayName": "default",
            "defaultConsumerKeyClaim": "azp",
            "defaultScopesClaim": "scope",
            "configurations": [
                {
                    "name": "consumer_key",
                    "label": "Consumer Key",
                    "type": "select",
                    "required": true,
                    "mask": true,
                    "multiple": true,
                    "tooltip": "Enter username to connect to key manager",
                    "default": {},
                    "values": [
                        {}
                    ]
                }
            ],
            "endpointConfigurations": [
                {
                    "name": "consumer_key",
                    "label": "Consumer Key",
                    "type": "select",
                    "required": true,
                    "mask": true,
                    "multiple": true,
                    "tooltip": "Enter username to connect to key manager",
                    "default": {},
                    "values": [
                        {}
                    ]
                }
            ]
        }
    ],
    "analyticsEnabled": false
};

AdvancedThrottlePolicyList advancedPolicyList = {
    "count": 1,
    "list": [
        {
            "policyId": "0c6439fd-9b16-3c2e-be6e-1086e0b9aa93",
            "policyName": "30PerMin",
            "displayName": "30PerMin",
            "description": "Allows 30 request per minute",
            "isDeployed": false,
            "type": "string",
            "defaultLimit": {
                "type": "REQUESTCOUNTLIMIT",
                "requestCount": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "requestCount": 30
                },
                "bandwidth": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "dataAmount": 1000,
                    "dataUnit": "KB"
                },
                "eventCount": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "eventCount": 3000
                }
            }
        }
    ]
};

SubscriptionThrottlePolicyList subPolicyList = {
    "count": 1,
    "list": [
        {
            "policyId": "0c6439fd-9b16-3c2e-be6e-1086e0b9aa93",
            "policyName": "30PerMin",
            "displayName": "30PerMin",
            "description": "Allows 30 request per minute",
            "isDeployed": false,
            "type": "string",
            "graphQLMaxComplexity": 400,
            "graphQLMaxDepth": 10,
            "defaultLimit": {
                "type": "REQUESTCOUNTLIMIT",
                "requestCount": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "requestCount": 30
                },
                "bandwidth": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "dataAmount": 1000,
                    "dataUnit": "KB"
                },
                "eventCount": {
                    "timeUnit": "min",
                    "unitTime": 10,
                    "eventCount": 3000
                }
            },
            "monetization": {
                "monetizationPlan": "FixedRate",
                "properties": {
                    "property1": "string",
                    "property2": "string"
                }
            },
            "rateLimitCount": 10,
            "rateLimitTimeUnit": "min",
            "subscriberCount": 10,
            "customAttributes": [],
            "stopOnQuotaReach": false,
            "billingPlan": "FREE",
            "permissions": {
                "permissionType": "deny",
                "roles": [
                    "Internal/everyone"
                ]
            }
        }
    ]
};

AdvancedThrottlePolicy policyCreated = {
    "policyId": "4cf46441-a538-4f79-a499-ab81200c9bca",
    "policyName": "10KPerMin",
    "displayName": "10KPerMin",
    "description": "Allows 10000 requests per minute",
    "isDeployed": true,
    "defaultLimit": {
        "type": "REQUESTCOUNTLIMIT",
        "requestCount": {
            "timeUnit": "min",
            "unitTime": 1,
            "requestCount": 10000
        }
    },
    "conditionalGroups": []
};