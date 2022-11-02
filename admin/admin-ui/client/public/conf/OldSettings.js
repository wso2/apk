const OldSettings = {
  "app": {
    "context": "/admin",
    "customUrl": {
      "enabled": false,
      "forwardedHeader": "X-Forwarded-For"
    },
    "origin": {
      "host": "localhost"
    },
    "feedback": {
      "enable": false,
      "serviceURL": ""
    },
    "singleLogout": {
      "enabled": true,
      "timeout": 2000
    },
    "docUrl": "https://apim.docs.wso2.com/en/4.1.0/",
    "minScopesToLogin": [
      "apim:api_workflow_view",
      "apim:api_workflow_approve",
      "apim:tenantInfo",
      "apim:admin_settings"
    ]
  }
}

if (typeof module !== 'undefined') {
  module.exports = OldSettings; // For Jest unit tests
}
