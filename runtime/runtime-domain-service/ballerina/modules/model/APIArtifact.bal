public type APIArtifact record {|
API api?;
Httproute productionRoute?;
Httproute sandboxRoute?;
ConfigMap definition?;
K8sServiceMapping[] serviceMapping = [];
RuntimeAPI runtimeAPI?;
map<Service> backendServices = {};
map<BackendPolicy> backendPolicies = {};
map<Authentication> authenticationMap = {};
boolean sandboxEndpointAvailable = false;
string productionUrl?;
string sandboxUrl?;
boolean productionEndpointAvailable = false;
string uniqueId;
|};