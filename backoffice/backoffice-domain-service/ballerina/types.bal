import ballerina/http;
import ballerina/constraint;

public type NotAcceptableError record {|
    *http:NotAcceptable;
    Error body;
|};

public type UnsupportedMediaTypeError record {|
    *http:UnsupportedMediaType;
    Error body;
|};

public type ForbiddenError record {|
    *http:Forbidden;
    Error body;
|};

public type InternalServerErrorError record {|
    *http:InternalServerError;
    Error body;
|};

public type ConflictError record {|
    *http:Conflict;
    Error body;
|};

public type PreconditionFailedError record {|
    *http:PreconditionFailed;
    Error body;
|};

public type CreatedComment record {|
    *http:Created;
    Comment body;
|};

public type NotFoundError record {|
    *http:NotFound;
    Error body;
|};

public type CreatedDocument record {|
    *http:Created;
    Document body;
|};

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type UnauthorizedError record {|
    *http:Unauthorized;
    Error body;
|};

public type DocumentList record {
    # Number of Documents returned.
    int count?;
    Document[] list?;
    Pagination pagination?;
};

public type UsageLimitBase record {
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the usage limit applies.
    int unitTime;
};

public type APIScope record {
    Scope scope;
};

public type Document record {
    string documentId?;
    @constraint:String {maxLength: 60, minLength: 1}
    string name;
    string 'type;
    @constraint:String {maxLength: 32766, minLength: 1}
    string summary?;
    string sourceType;
    string sourceUrl?;
    string fileName?;
    string inlineContent?;
    string otherTypeName?;
    string visibility;
    string createdTime?;
    string createdBy?;
    string lastUpdatedTime?;
    string lastUpdatedBy?;
};

public type ExternalStore record {
    # The external store identifier, which is a unique value.
    string id?;
    # The name of the external API Store that is displayed in the Publisher UI.
    string displayName?;
    # The type of the Store. This can be a WSO2-specific API Store or an external one.
    string 'type?;
    # The endpoint URL of the external store
    string endpoint?;
};

public type Pagination record {
    int offset?;
    int 'limit?;
    int total?;
    # Link to the next subset of resources qualified.
    # Empty if no more resources are to be returned.
    string next?;
    # Link to the previous subset of resources qualified.
    # Empty if current subset is the first subset returned.
    string previous?;
};

public type EventCountLimit record {
    *UsageLimitBase;
    # Maximum number of events allowed
    int eventCount;
};

public type ResourcePath record {
    int id;
    string resourcePath?;
    string httpVerb?;
};

public type API_additionalProperties record {
    string name?;
    string value?;
    boolean display?;
};

public type FileInfo record {
    # relative location of the file (excluding the base context and host of the Publisher API)
    string relativePath?;
    # media-type of the file
    string mediaType?;
};

public type APIOperations record {
    string id?;
    string target?;
    string verb?;
    string usagePlan?;
};

public type APIList record {
    # Number of APIs returned.
    int count?;
    APIInfo[] list?;
    Pagination pagination?;
};

public type APIMonetizationUsage record {
    # Map of custom properties related to monetization usage
    record {} properties?;
};

public type Subscription record {
    string subscriptionId;
    ApplicationInfo applicationInfo;
    string usagePlan;
    string subscriptionStatus;
};

public type Settings record {
    # The Developer Portal URL
    string devportalUrl?;
    Environment[] environment?;
    string[] scopes?;
    MonetizationAttribute[] monetizationAttributes?;
    # Is Document Visibility configuration enabled
    boolean docVisibilityEnabled?;
    # Authorization Header
    string authorizationHeader?;
};

public type APIRevision record {
    string displayName?;
    string id?;
    @constraint:String {maxLength: 255}
    string description?;
    string createdTime?;
};

public type APIExternalStoreList record {
    # Number of external stores returned.
    int count?;
    APIExternalStore[] list?;
};

public type ResourcePathList record {
    # Number of API Resource Paths returned.
    int count?;
    ResourcePath[] list?;
    Pagination pagination?;
};

public type APIBusinessInformation record {
    @constraint:String {maxLength: 120}
    string businessOwner?;
    string businessOwnerEmail?;
    @constraint:String {maxLength: 120}
    string technicalOwner?;
    string technicalOwnerEmail?;
};

public type APIMonetizationInfo record {
    # Flag to indicate the monetization status
    boolean enabled;
    # Map of custom properties related to monetization
    record {} properties?;
};

public type LifecycleHistoryItem record {
    string previousState?;
    string postState?;
    string user?;
    string updatedTime?;
};

public type UsagePlanList record {
    # Number of Usage Plans returned.
    int count?;
    # Array of Usage Policies
    UsagePlan[] list?;
    Pagination pagination?;
};

public type APIDeployment record {
    @constraint:String {maxLength: 255, minLength: 1}
    string name?;
    string deployedTime?;
};

public type MonetizationAttribute record {
    # Is attribute required
    boolean required?;
    # Name of the attribute
    string name?;
    # Display name of the attribute
    string displayName?;
    # Description of the attribute
    string description?;
    # Is attribute hidden
    boolean hidden?;
    # Default value of the attribute
    string default?;
};

public type APIDefinition record {
    string 'type;
    string schemaDefinition?;
};

public type CommenterInfo record {
    string firstName?;
    string lastName?;
    string fullName?;
};

public type Environment record {
    string id;
    string name;
    string displayName?;
    string 'type;
    string serverUrl;
    string provider?;
    boolean showInApiConsole;
    GatewayEnvironmentProtocolURI[] endpointURIs?;
    AdditionalProperty[] additionalProperties?;
};

public type CommentList record {
    # Number of Comments returned.
    int count?;
    Comment[] list?;
    Pagination pagination?;
};

public type CustomAttribute record {
    # Name of the custom attribute
    string name;
    # Value of the custom attribute
    string value;
};

public type ErrorListItem record {
    string code;
    # Description about individual errors occurred
    string message;
    # A detail description about the error message.
    string description?;
};

public type APIRevenue record {
    # Map of custom properties related to API revenue
    record {} properties?;
};

public type APICategoryList record {
    # Number of API categories returned.
    int count?;
    APICategory[] list?;
};

public type DocumentId_content_body record {
    # Document to upload
    string file?;
    # Inline content of the document
    string inlineContent?;
};

public type ApplicationInfo record {
    string applicationId?;
    string name?;
    string subscriber?;
    string description?;
    int subscriptionCount?;
};

public type ModifiableAPI record {
    # UUID of the API
    string id?;
    # Name of the API
    @constraint:String {maxLength: 50, minLength: 1}
    string name;
    @constraint:String {maxLength: 60, minLength: 1}
    string context?;
    # A brief description about the API
    string description?;
    boolean hasThumbnail?;
    # State of the API. Only published APIs are visible on the Developer Portal
    string state = "CREATED";
    string[] tags?;
    record {} additionalProperties?;
    APIMonetizationInfo monetization?;
    APIBusinessInformation businessInformation?;
    # API categories
    string[] categories?;
};

public type LifecycleState record {
    string state?;
    LifecycleState_availableTransitions[] availableTransitions?;
};

public type ApiId_thumbnail_body record {
    # Image to upload
    string file;
};

public type SubscriptionThrottlePolicyPermission record {
    string permissionType;
    string[] roles;
};

public type SearchResultList record {
    # Number of results returned.
    int count?;
    record {}[] list?;
    Pagination pagination?;
};

public type LifecycleHistory record {
    int count?;
    LifecycleHistoryItem[] list?;
};

public type LifecycleState_availableTransitions record {
    string event?;
    string targetState?;
};

public type ExternalStoreList record {
    # Number of external stores returned.
    int count?;
    ExternalStore[] list?;
};

public type UsagePlan record {
    # Id of policy
    int policyId?;
    # policy uuid
    string uuid?;
    # Name of policy
    @constraint:String {maxLength: 60, minLength: 1}
    string policyName?;
    # Display name of the policy
    @constraint:String {maxLength: 512}
    string displayName?;
    # Description of the policy
    @constraint:String {maxLength: 1024}
    string description?;
    # Usage policy organization
    string organization?;
    UsageLimit defaultLimit;
    # Burst control request count
    int rateLimitCount?;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Number of subscriptions allowed
    int subscriberCount?;
    # Custom attributes added to the Usage plan
    CustomAttribute[] customAttributes?;
    # This indicates the action to be taken when a user goes beyond the allocated quota. If checked, the user's requests will be dropped. If unchecked, the requests will be allowed to pass through.
    boolean stopOnQuotaReach = false;
    # define whether this is Paid or a Free plan. Allowed values are FREE or COMMERCIAL.
    string billingPlan?;
    SubscriptionThrottlePolicyPermission permissions?;
};

public type Comment record {
    string id?;
    @constraint:String {maxLength: 512}
    string content;
    string createdTime?;
    string createdBy?;
    string updatedTime?;
    string category = "general";
    string parentCommentId?;
    string entryPoint?;
    CommenterInfo commenterInfo?;
    CommentList replies?;
};

public type RequestCountLimit record {
    *UsageLimitBase;
    # Maximum number of requests allowed
    int requestCount;
};

public type ThreatProtectionPolicy record {
    # Policy ID
    string uuid?;
    # Name of the policy
    string name;
    # Type of the policy
    string 'type;
    # policy as a json string
    string policy;
};

public type APIInfo record {
    string id?;
    string name?;
    string description?;
    string context?;
    string 'version?;
    string 'type?;
    string createdTime?;
    string updatedTime?;
    boolean hasThumbnail?;
    # State of the API. Only published APIs are visible on the Developer Portal
    string state?;
};

public type APIExternalStore record {
    # The external store identifier, which is a unique value.
    string id?;
    # The recent timestamp which a given API is updated in the external store.
    string lastUpdatedTime?;
};

public type Error record {
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

public type SearchResult record {
    string id?;
    string name;
    string 'type?;
    # Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL
    string transportType?;
};

public type GatewayEnvironmentProtocolURI record {
    string protocol;
    string endpointURI;
};

public type UsageLimit record {
    # Type of the usage limit. Allowed values are "REQUESTCOUNTLIMIT" and "BANDWIDTHLIMIT".
    # Please see schemas of "RequestCountLimit" and "BandwidthLimit" usage limit types in
    # Definitions section.
    string 'type;
    RequestCountLimit requestCount?;
    BandwidthLimit bandwidth?;
    EventCountLimit eventCount?;
};

public type Scope record {
    # UUID of the Scope. Valid only for shared scopes.
    string id?;
    # name of Scope
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    # display name of Scope
    @constraint:String {maxLength: 255}
    string displayName?;
    # description of Scope
    @constraint:String {maxLength: 512}
    string description?;
    # role bindings list of the Scope
    string[] bindings?;
    # usage count of Scope
    int usageCount?;
};

public type PatchRequestBody record {
    # Content of the comment
    @constraint:String {maxLength: 512}
    string content?;
    # Category of the comment
    string category?;
};

public type PostRequestBody record {
    # Content of the comment
    @constraint:String {maxLength: 512}
    string content;
    # Category of the comment
    string category?;
};

public type SubscriptionList record {
    # Number of Subscriptions returned.
    int count?;
    Subscription[] list?;
    Pagination pagination?;
};

public type WorkflowResponse record {
    # This attribute declares whether this workflow task is approved or rejected.
    string workflowStatus;
    # Attributes that returned after the workflow execution
    string jsonPayload?;
    LifecycleState lifecycleState?;
};

public type API record {
    # UUID of the API
    string id?;
    @constraint:String {maxLength: 60, minLength: 1}
    string name;
    @constraint:String {maxLength: 32766}
    string description?;
    @constraint:String {maxLength: 232, minLength: 1}
    string context;
    @constraint:String {maxLength: 30, minLength: 1}
    string 'version;
    # The api creation type to be used. Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL, WEBSUB, SSE, WEBHOOK, ASYNC
    string 'type = "HTTP";
    # Supported transports for the API (http and/or https).
    string[] transport?;
    boolean hasThumbnail?;
    # State of the API. Only published APIs are visible on the Developer Portal
    string state = "CREATED";
    string[] tags?;
    # API categories
    string[] categories?;
    record {} additionalProperties?;
    string createdTime?;
    string lastUpdatedTime?;
    APIOperations[] operations?;
    # The API level usage policy selected for the particular Runtime API
    string apiUsagePolicy?;
    APIMonetizationInfo monetization?;
    APIBusinessInformation businessInformation?;
    APIRevision revision?;
    APIDeployment[] deployments?;
};

public type AdditionalProperty record {
    string 'key?;
    string value?;
};

public type APICategory record {
    string id?;
    string name;
    string description?;
};

public type SubscriberInfo record {
    string name?;
};

public type BandwidthLimit record {
    *UsageLimitBase;
    # Amount of data allowed to be transfered
    int dataAmount;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit;
};
