public type Service record {
string apiVersion = "v1";
string kind = "Service";
Metadata metadata;
ServiceSpec spec;
};

public type ServiceSpec record {
string externalName?;
string 'type;
Port[] ports?;
};

public type Port record {
string name?;
int nodePort?;
int port;
string protocol?;
int targetPort?;
string appProtocol?;
};

public type ServiceList record {
string kind = "ServiceList";
string apiVersion = "v1";
ListMeta metadata;
Service[] items;
};