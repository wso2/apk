public type APIArtifact record {|
API api?;
Httproute[] productionRoute = [];
Httproute[] sandboxRoute = [];
ConfigMap definition?;
K8sServiceMapping[] serviceMapping = [];
RuntimeAPI runtimeAPI?;
map<Backend> backendServices = {};
map<Authentication> authenticationMap = {};
map<Scope> scopes = {};
map<RateLimitPolicy> rateLimitPolicies = {};
boolean sandboxEndpointAvailable = false;
string productionUrl?;
string sandboxUrl?;
boolean productionEndpointAvailable = false;
string uniqueId;
|};