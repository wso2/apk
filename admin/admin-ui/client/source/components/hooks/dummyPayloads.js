export const settings = {
    "analyticsEnabled": false,
    "keyManagerConfiguration": [
      {
        "defaultConsumerKeyClaim": "azp",
        "endpointConfigurations": [
          null,
          null
        ],
        "displayName": "default",
        "configurations": [
          {
            "default": "admin",
            "values": [
              {},
              {}
            ],
            "name": "consumer_key",
            "multiple": true,
            "tooltip": "Enter username to connect to key manager",
            "label": "Consumer Key",
            "type": "select",
            "required": true,
            "mask": true
          },
          {
            "default": "admin",
            "values": [
              {},
              {}
            ],
            "name": "consumer_key",
            "multiple": true,
            "tooltip": "Enter username to connect to key manager",
            "label": "Consumer Key",
            "type": "select",
            "required": true,
            "mask": true
          }
        ],
        "defaultScopesClaim": "scope",
        "type": "default"
      },
      {
        "defaultConsumerKeyClaim": "azp",
        "endpointConfigurations": [
          null,
          null
        ],
        "displayName": "default",
        "configurations": [
          {
            "default": "admin",
            "values": [
              {},
              {}
            ],
            "name": "consumer_key",
            "multiple": true,
            "tooltip": "Enter username to connect to key manager",
            "label": "Consumer Key",
            "type": "select",
            "required": true,
            "mask": true
          },
          {
            "default": "admin",
            "values": [
              {},
              {}
            ],
            "name": "consumer_key",
            "multiple": true,
            "tooltip": "Enter username to connect to key manager",
            "label": "Consumer Key",
            "type": "select",
            "required": true,
            "mask": true
          }
        ],
        "defaultScopesClaim": "scope",
        "type": "default"
      }
    ],
    "scopes": [
      "scopes",
      "scopes"
    ]
  }

  export const tenant = {
    "tenantId": -1234,
    "tenantDomain": "carbon.super",
    "username": "john"
  }

  export const apiCategories = {
    "count": 1,
    "list": [
      {
        "numberOfAPIs": 1,
        "name": "Finance",
        "description": "Finance related APIs",
        "id": "01234567-0123-0123-0123-012345678901"
      },
      {
        "numberOfAPIs": 1,
        "name": "Finance",
        "description": "Finance related APIs",
        "id": "01234567-0123-0123-0123-012345678901"
      }
    ]
  };

  export const applicationThrottlePolicies = {
    "count": 2,
    "list": [
        {
            "defaultLimit": {
                "type": "REQUESTCOUNTLIMIT",
                "requestCount": {
                    "requestCount": 41,
                    "timeUnit": "min",
                    "unitTime": 1
                }
            },
            "policyId": "cbee719f-ea93-4578-91a5-df94df99a008",
            "policyName": "42PerMin",
            "displayName": "41PerMin",
            "description": "Allows 30 request per minute",
            "isDeployed": false,
            "type": "ApplicationThrottlePolicy"
        },
        {
            "defaultLimit": {
                "type": "REQUESTCOUNTLIMIT",
                "requestCount": {
                    "requestCount": 32,
                    "timeUnit": "min",
                    "unitTime": 1
                }
            },
            "policyId": "d94c27e8-3867-482d-8218-9dc69e027ebe",
            "policyName": "32PerMin",
            "displayName": "32PerMin",
            "description": "Allows 32 request per minute",
            "isDeployed": false,
            "type": "ApplicationThrottlePolicy"
        }
    ]
}