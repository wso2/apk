## KGW Config Deployer Service

KGW Config Deployer Service.

# Functionalities.

1. Generate KGW configuration (api.apk-conf) from given OAS definition.
2. Generate K8s artifacts from given definition and KGW configuration file.
3. Deploy API into Gateway getting from KGW configuration and definition.
4. Undeploy API from Gateway.



"typemeta:\n  kind: Backend\n  apiversion: gateway.envoyproxy.io/v1alpha1\nobjectmeta:\n  name: Sample-API7603c81c1b-1-1-6261636b\n  generatename: \"\"\n  namespace: \"\"\n  selflink: \"\"\n  uid: \"\"\n  resourceversion: \"\"\n  generation: 0\n  creationtimestamp: \"0001-01-01T00:00:00Z\"\n  deletiontimestamp: null\n  deletiongraceperiodseconds: null\n  labels: {}\n  annotations: {}\n  ownerreferences: []\n  finalizers: []\n  managedfields: []\nspec:\n  type: null\n  endpoints:\n    - hostname: null\n      fqdn:\n        hostname: dev-tools.w...+148 more"