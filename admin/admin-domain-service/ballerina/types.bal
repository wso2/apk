import ballerina/http;
import ballerina/constraint;

public type AcceptedWorkflowResponse record {|
    *http:Accepted;
    WorkflowResponse body;
|};

public type NotAcceptableError record {|
    *http:NotAcceptable;
    Error body;
|};

public type UnsupportedMediaTypeError record {|
    *http:UnsupportedMediaType;
    Error body;
|};

public type OkKeyManagerWellKnownResponse record {|
    *http:Ok;
    KeyManagerWellKnownResponse body;
|};

public type InternalServerErrorError record {|
    *http:InternalServerError;
    Error body;
|};

public type ForbiddenError record {|
    *http:Forbidden;
    Error body;
|};

public type ConflictError record {|
    *http:Conflict;
    Error body;
|};

public type OkWorkflowInfo record {|
    *http:Ok;
    WorkflowInfo body;
|};

public type NotFoundError record {|
    *http:NotFound;
    Error body;
|};

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type Policy record {
    # Id of plan
    string planId?;
    # Name of plan
    @constraint:String {maxLength: 60, minLength: 1}
    string planName;
    # Display name of the policy
    @constraint:String {maxLength: 512}
    string displayName?;
    # Description of the policy
    @constraint:String {maxLength: 1024}
    string description?;
    # Indicates whether the policy is deployed successfully or not.
    boolean isDeployed = false;
    # Indicates the type of throttle policy
    string 'type?;
};

public type EnvironmentList record {
    # Number of Environments returned.
    int count?;
    Environment[] list?;
};

# Blocking Conditions
public type BlockingCondition record {
    # Id of the blocking condition
    string policyId?;
    # Type of the blocking condition
    string conditionType;
    # Value of the blocking condition
    string conditionValue;
    # Status of the blocking condition
    boolean conditionStatus?;
};

public type WorkflowProperties record {
    string name?;
    boolean enable?;
    string[] properties?;
};

public type ApplicationRatePlan record {
    *Policy;
    ThrottleLimit defaultLimit;
};

public type Pagination record {
    int offset?;
    int 'limit?;
    int total?;
    # Link to the next subset of resources qualified.
    # Empty if no more resources are to be returned.
    # example: ""
    string next?;
    # Link to the previous subset of resources qualified.
    # Empty if current subset is the first subset returned.
    # example: ""
    string previous?;
};

public type EventCountLimit record {
    *ThrottleLimitBase;
    # Maximum number of events allowed
    int eventCount;
};

public type ThrottleLimitBase record {
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the throttling limit applies.
    int unitTime;
};

public type ClaimMappingEntry record {
    string remoteClaim?;
    string localClaim?;
};

public type BusinessPlan record {
    *Policy;
    *GraphQLQuery;
    ThrottleLimit defaultLimit;
    # Burst control request count
    int rateLimitCount?;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Number of subscriptions allowed
    int subscriberCount?;
    # Custom attributes added to the Subscription Throttling Policy
    CustomAttribute[] customAttributes?;
    BusinessPlanPermission permissions?;
};

public type PolicyDetails record {
    # Id of policy
    int policyId?;
    # UUId of policy
    string uuid?;
    # Name of policy
    @constraint:String {maxLength: 60, minLength: 1}
    string policyName;
    # Display name of the policy
    @constraint:String {maxLength: 512}
    string displayName?;
    # Description of the policy
    @constraint:String {maxLength: 1024}
    string description?;
    # Indicates whether the policy is deployed successfully or not.
    boolean isDeployed = false;
    # Indicates the type of throttle policy
    string 'type?;
};

public type KeyManagerWellKnownResponse record {
    boolean valid?;
    KeyManager value?;
};

public type KeyManager record {
    string id?;
    @constraint:String {maxLength: 100, minLength: 1}
    string name;
    # display name of Key Manager to  show in UI
    @constraint:String {maxLength: 100}
    string displayName?;
    @constraint:String {maxLength: 45, minLength: 1}
    string 'type;
    @constraint:String {maxLength: 256}
    string description?;
    # Well-Known Endpoint of Identity Provider.
    string wellKnownEndpoint?;
    KeyManagerEndpoint[] endpoints?;
    KeyManager_signingCertificate signingCertificate?;
    # PEM type certificate
    string tlsCertificate?;
    string issuer;
    string[] availableGrantTypes?;
    boolean enableTokenGeneration = true;
    boolean enableMapOAuthConsumerApps = false;
    boolean enableOauthAppValidation = true;
    boolean enableOAuthAppCreation = true;
    string consumerKeyClaim?;
    string scopesClaim?;
    boolean enabled = true;
    record {} additionalProperties?;
};

public type CustomUrlInfo_devPortal record {
    string url?;
};

public type Settings record {
    string[] scopes?;
    Settings_keyManagerConfiguration[] keyManagerConfiguration?;
    # To determine whether analytics is enabled or not
    boolean analyticsEnabled?;
};

public type KeyManagerConfiguration record {
    string name?;
    string label?;
    string 'type?;
    boolean required?;
    boolean mask?;
    boolean multiple?;
    string tooltip?;
    string default?;
    string[] values?;
};

public type BusinessPlanList record {
    # Number of Business Plans returned.
    int count?;
    BusinessPlan[] list?;
};

public type ApplicationRatePlanList record {
    # Number of Application Rate Plans returned.
    int count?;
    ApplicationRatePlan[] list?;
};

public type OrganizationList record {
    # Number of Organization returned.
    int count?;
    Organization[] list?;
};

public type ThrottleLimit record {
    # Type of the throttling limit. Allowed values are "REQUESTCOUNTLIMIT" and "BANDWIDTHLIMIT".
    # Please see schemas of "RequestCountLimit" and "BandwidthLimit" throttling limit types in
    # Definitions section.
    string 'type;
    RequestCountLimit requestCount?;
    BandwidthLimit bandwidth?;
    EventCountLimit eventCount?;
};

public type TokenValidation record {
    int id?;
    boolean enable?;
    string 'type?;
    record {} value?;
};

public type Keymanagers_discover_body record {
    # Well-Known Endpoint
    string url?;
    # Key Manager Type
    string 'type?;
};

public type Environment record {
    string id?;
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    @constraint:String {maxLength: 255, minLength: 1}
    string displayName?;
    string provider?;
    @constraint:String {maxLength: 1023}
    string description?;
    boolean isReadOnly?;
    @constraint:Array {minLength: 1}
    VHost[] vhosts;
    GatewayEnvironmentProtocolURI[] endpointURIs?;
    AdditionalProperty[] additionalProperties?;
};

public type KeyManager_signingCertificate record {
    string 'type?;
    string value?;
};

public type BusinessPlanPermission record {
    string permissionType;
    string[] roles;
};

public type Application record {
    string applicationId?;
    string name?;
    string throttlingPolicy?;
    string description?;
    # Type of the access token generated for this application.
    # **OAUTH:** A UUID based access token which is issued by default.
    # **JWT:** A self-contained, signed JWT based access token. **Note:** This can be only used in Microgateway environments.
    string tokenType?;
    string status?;
    string[] groups?;
    int subscriptionCount?;
    record {|string...;|} attributes?;
    ScopeInfo[] subscriptionScopes?;
    # Application created user
    string owner?;
};

public type VHost record {
    @constraint:String {maxLength: 255, minLength: 1}
    string host;
    @constraint:String {maxLength: 255}
    string httpContext?;
    int httpPort?;
    int httpsPort?;
    int wsPort?;
    int wssPort?;
};

public type Organization record {
    string id?;
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    @constraint:String {maxLength: 255, minLength: 1}
    string displayName;
    @constraint:String {maxLength: 255, minLength: 1}
    string organizationClaimValue?;
    boolean enabled = true;
    string[] serviceNamespaces = ["*"];
    WorkflowProperties[] workflows?;
    string[] production?;
    string[] sandbox?;
};

public type MonetizationUsagePublishInfo record {
    # State of usage publish job
    string state?;
    # Status of usage publish job
    string status?;
    # Timestamp of the started time of the Job
    string startedTime?;
    # Timestamp of the last published time
    string lastPublsihedTime?;
};

public type ErrorListItem record {
    # Error code
    string code;
    # Description about individual errors occurred
    string message;
};

public type BlockingConditionList record {
    # Number of Blocking Conditions returned.
    int count?;
    BlockingCondition[] list?;
};

public type CustomAttribute record {
    # Name of the custom attribute
    string name;
    # Value of the custom attribute
    string value;
};

public type APICategoryList record {
    # Number of API categories returned.
    int count?;
    APICategory[] list?;
};

public type ApplicationInfo record {
    string applicationId?;
    string name?;
    string owner?;
    string status?;
    string groupId?;
};

# The tenant information of the user
public type TenantInfo record {
    string username?;
    string tenantDomain?;
    int tenantId?;
};

public type KeyManagerList record {
    # Number of Key managers returned.
    int count?;
    KeyManagerInfo[] list?;
};

public type Policies_import_body record {
    # Json File
    record {byte[] fileContent; string fileName;} file;
};

# Blocking Conditions Status
public type BlockingConditionStatus record {
    # Id of the blocking condition
    string policyId?;
    # Status of the blocking condition
    boolean conditionStatus;
};

public type KeyManagerEndpoint record {
    string name;
    string value;
};

public type Settings_keyManagerConfiguration record {
    string 'type?;
    string displayName?;
    string defaultConsumerKeyClaim?;
    string defaultScopesClaim?;
    KeyManagerConfiguration[] configurations?;
    KeyManagerConfiguration[] endpointConfigurations?;
};

public type WorkflowList record {
    # Number of workflow processes returned.
    int count?;
    # Link to the next subset of resources qualified.
    # Empty if no more resources are to be returned.
    string next?;
    # Link to the previous subset of resources qualified.
    # Empty if current subset is the first subset returned.
    string previous?;
    WorkflowInfo[] list?;
};

public type ApplicationList record {
    # Number of applications returned.
    int count?;
    ApplicationInfo[] list?;
    Pagination pagination?;
};

public type PublishStatus record {
    # Status of the usage publish request
    string status?;
    # detailed message of the status
    string message?;
};

public type RequestCountLimit record {
    *ThrottleLimitBase;
    # Maximum number of requests allowed
    int requestCount;
};

public type KeyManagerInfo record {
    string id?;
    string name;
    string 'type;
    string description?;
    boolean enabled?;
};

public type WorkflowInfo record {
    # Type of the Workflow Request. It shows which type of request is it.
    string workflowType?;
    # Show the Status of the the workflow request whether it is approved or created.
    string workflowStatus?;
    # Time of the the workflow request created.
    string createdTime?;
    # Time of the the workflow request updated.
    string updatedTime?;
    # description is a message with basic details about the workflow request.
    string description?;
};

public type ExportPolicy record {
    string 'type?;
    string subtype?;
    string version?;
    record {} data?;
};

public type Error record {
    # Error code
    int code;
    # Error message.
    string message;
    # A detail description about the error message.
    string description?;
    # Preferably an url with more details about the error.
    string moreInfo?;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] 'error?;
};

public type ScopeInfo record {
    string 'key?;
    string name?;
    # Allowed roles for the scope
    string[] roles?;
    # Description of the scope
    string description?;
};

public type GatewayEnvironmentProtocolURI record {
    string protocol;
    string endpointURI;
};

public type GraphQLQuery record {
    # Maximum Complexity of the GraphQL query
    int graphQLMaxComplexity?;
    # Maximum Depth of the GraphQL query
    int graphQLMaxDepth?;
};

public type PolicyDetailsList record {
    # Number of Throttling Policies returned.
    int count?;
    PolicyDetails[] list?;
};

# The custom url information of the tenant domain
public type CustomUrlInfo record {
    string tenantDomain?;
    string tenantAdminUsername?;
    boolean enabled?;
    CustomUrlInfo_devPortal devPortal?;
};

public type WorkflowResponse record {
    # This attribute declares whether this workflow task is approved or rejected.
    string workflowStatus;
    # Attributes that returned after the workflow execution
    string jsonPayload?;
};

public type AdditionalProperty record {
    string 'key?;
    string value?;
};

public type APICategory record {
    string id?;
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    @constraint:String {maxLength: 1024}
    string description?;
    int numberOfAPIs?;
};

public type BandwidthLimit record {
    *ThrottleLimitBase;
    # Amount of data allowed to be transferred
    int dataAmount;
    # Unit of data allowed to be transferred. Allowed values are "KB", "MB" and "GB"
    string dataUnit;
};
