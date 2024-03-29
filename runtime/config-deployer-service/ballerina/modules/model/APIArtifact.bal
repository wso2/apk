public type APIArtifact record {|
    string name;
    string 'version;
    API api?;
    HTTPRoute[] productionHttpRoutes = [];
    HTTPRoute[] sandboxHttpRoutes = [];
    GQLRoute[] productionGqlRoutes = [];
    GQLRoute[] sandboxGqlRoutes = [];
    ConfigMap definition?;
    map<ConfigMap> endpointCertificates = {};
    map<string> certificateMap = {};
    map<Backend> backendServices = {};
    map<Authentication> authenticationMap = {};
    map<Scope> scopes = {};
    map<RateLimitPolicy> rateLimitPolicies = {};
    map<APIPolicy> apiPolicies = {};
    map<InterceptorService> interceptorServices = {};
    boolean sandboxEndpointAvailable = false;
    string productionUrl?;
    string sandboxUrl?;
    boolean productionEndpointAvailable = false;
    string uniqueId;
    map<K8sSecret> secrets = {};
    BackendJWT backendJwt?;
    string namespace?;
    string organization;
|};
