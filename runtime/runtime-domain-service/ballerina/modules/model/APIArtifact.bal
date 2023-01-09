public type APIArtifact record {|
API api?;
Httproute productionRoute?;
Httproute sandboxRoute?;
ConfigMap definition?;
K8sServiceMapping[] serviceMapping = [];
Service[] backendServices = [];
boolean sandboxEndpointAvailable = false;
string productionUrl?;
string sandboxUrl?;
boolean productionEndpointAvailable = false;
string uniqueId;

|};