import ballerina/http;
import ballerina/constraint;

public type CreatedAdvancedThrottlePolicy record {|
    *http:Created;
    AdvancedThrottlePolicy body;
|};

public type AcceptedWorkflowResponse record {|
    *http:Accepted;
    WorkflowResponse body;
|};

public type PayloadTooLargeError record {|
    *http:PayloadTooLarge;
    Error body;
|};

public type AcceptedPublishStatus record {|
    *http:Accepted;
    PublishStatus body;
|};

public type InternalServerErrorError record {|
    *http:InternalServerError;
    Error body;
|};

public type ConflictError record {|
    *http:Conflict;
    Error body;
|};

public type CreatedEnvironment record {|
    *http:Created;
    Environment body;
|};

public type NotFoundError record {|
    *http:NotFound;
    Error body;
|};

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type CreatedAPICategory record {|
    *http:Created;
    APICategory body;
|};

public type NotAcceptableError record {|
    *http:NotAcceptable;
    Error body;
|};

public type CreatedBlockingCondition record {|
    *http:Created;
    BlockingCondition body;
|};

public type UnsupportedMediaTypeError record {|
    *http:UnsupportedMediaType;
    Error body;
|};

public type CreatedKeyManager record {|
    *http:Created;
    KeyManager body;
|};

public type ForbiddenError record {|
    *http:Forbidden;
    Error body;
|};

public type CreatedSubscriptionThrottlePolicy record {|
    *http:Created;
    SubscriptionThrottlePolicy body;
|};

public type CreatedApplicationThrottlePolicy record {|
    *http:Created;
    ApplicationThrottlePolicy body;
|};

public type SubscriptionThrottlePolicyList record {
    # Number of Subscription Throttling Policies returned.
    int count?;
    SubscriptionThrottlePolicy[] list?;
};

public type EnvironmentList record {
    # Number of Environments returned.
    int count?;
    Environment[] list?;
};

# Blocking Conditions
public type BlockingCondition record {
    # Id of the blocking condition
    string conditionId?;
    # Type of the blocking condition
    string conditionType;
    # Value of the blocking condition
    record {} conditionValue;
    # Status of the blocking condition
    boolean conditionStatus?;
};

public type ApplicationThrottlePolicyList record {
    # Number of Application Throttling Policies returned.
    int count?;
    ApplicationThrottlePolicy[] list?;
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

public type FileInfo record {
    # relative location of the file (excluding the base context and host of the Admin API)
    string relativePath?;
    # media-type of the file
    string mediaType?;
};

public type RoleAlias record {
    # The original role
    string role?;
    # The role mapping for role alias
    string[] aliases?;
};

public type ThrottlePolicyList record {
    # Number of Throttling Policies returned.
    int count?;
    ThrottlePolicy[] list?;
    Pagination pagination?;
};

public type ThrottleLimitBase record {
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the throttling limit applies.
    int unitTime;
};

public type ScopeList record {
    # Number of scopes available for tenant.
    int count?;
    Scope[] list?;
};

public type ClaimMappingEntry record {
    string remoteClaim?;
    string localClaim?;
};

public type KeyManagerWellKnownResponse record {
    boolean valid?;
    KeyManager value?;
};

public type SettingsKeymanagerconfiguration record {
    string 'type?;
    string displayName?;
    string defaultConsumerKeyClaim?;
    string defaultScopesClaim?;
    KeyManagerConfiguration[] configurations?;
    KeyManagerConfiguration[] endpointConfigurations?;
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
    string introspectionEndpoint?;
    string clientRegistrationEndpoint?;
    string tokenEndpoint?;
    string displayTokenEndpoint?;
    string revokeEndpoint?;
    string displayRevokeEndpoint?;
    string userInfoEndpoint?;
    string authorizeEndpoint?;
    KeyManagerEndpoint[] endpoints?;
    KeymanagerCertificates certificates?;
    string issuer?;
    # The alias of Identity Provider.
    # If the tokenType is EXCHANGED, the alias value should be inclusive in the audience values of the JWT token
    string alias?;
    string scopeManagementEndpoint?;
    string[] availableGrantTypes?;
    boolean enableTokenGeneration?;
    boolean enableTokenEncryption = false;
    boolean enableTokenHashing = false;
    boolean enableMapOAuthConsumerApps = false;
    boolean enableOAuthAppCreation = false;
    boolean enableSelfValidationJWT = true;
    ClaimMappingEntry[] claimMapping?;
    string consumerKeyClaim?;
    string scopesClaim?;
    TokenValidation[] tokenValidation?;
    boolean enabled?;
    record {} additionalProperties?;
    # The type of the tokens to be used (exchanged or without exchanged). Accepted values are EXCHANGED, DIRECT and BOTH.
    string tokenType = "DIRECT";
};

public type Settings record {
    string[] scopes?;
    SettingsKeymanagerconfiguration[] keyManagerConfiguration?;
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
    record {} default?;
    record {}[] values?;
};

public type KeymanagersDiscoverBody record {
    # Well-Known Endpoint
    string url?;
    # Key Manager Type
    string 'type?;
};

public type ScopeSettings record {
    string name?;
};

public type ThrottlePolicyDetailsList record {
    # Number of Throttling Policies returned.
    int count?;
    ThrottlePolicyDetails[] list?;
};

public type HeaderCondition record {
    # Name of the header
    string headerName;
    # Value of the header
    string headerValue;
};

public type TenantthemeBody record {
    # Zip archive consisting of tenant theme configuration
    string file;
};

public type Workflow record {
    # This attribute declares whether this workflow task is approved or rejected.
    string status;
    # Custom attributes to complete the workflow task
    record {} attributes?;
    string description?;
};

public type ThrottlePolicy record {
    # Id of policy
    string policyId?;
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

public type BotDetectionData record {
    # The time of detection
    int recordedTime?;
    # The message ID
    string messageID?;
    # The api method
    string apiMethod?;
    # The header set
    string headerSet?;
    # The content of the message body
    string messageBody?;
    # The IP of the client
    string clientIp?;
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

public type SubscriptionThrottlePolicy record {
    *ThrottlePolicy;
    *GraphQLQuery;
    ThrottleLimit defaultLimit;
    MonetizationInfo monetization?;
    # Burst control request count
    int rateLimitCount?;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Number of subscriptions allowed
    int subscriberCount?;
    # Custom attributes added to the Subscription Throttling Policy
    CustomAttribute[] customAttributes?;
    # This indicates the action to be taken when a user goes beyond the allocated quota. If checked, the user's requests will be dropped. If unchecked, the requests will be allowed to pass through.
    boolean stopOnQuotaReach = false;
    # define whether this is Paid or a Free plan. Allowed values are FREE or COMMERCIAL.
    string billingPlan?;
    SubscriptionThrottlePolicyPermission permissions?;
};

public type ApplicationThrottlePolicy record {
    *ThrottlePolicy;
    ThrottleLimit defaultLimit;
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

public type ThrottlePolicyDetails record {
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

# Conditions used for Throttling
public type ThrottleCondition record {
    # Type of the throttling condition. Allowed values are "HEADERCONDITION", "IPCONDITION", "JWTCLAIMSCONDITION"
    # and "QUERYPARAMETERCONDITION".
    string 'type;
    # Specifies whether inversion of the condition to be matched against the request.
    #
    # **Note:** When you add conditional groups for advanced throttling policies, this paramater should have the
    # same value ('true' or 'false') for the same type of conditional group.
    boolean invertCondition = false;
    HeaderCondition headerCondition?;
    IPCondition ipCondition?;
    JWTClaimsCondition jwtClaimsCondition?;
    QueryParameterCondition queryParameterCondition?;
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
    record {} attributes?;
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

public type JWTClaimsCondition record {
    # JWT claim URL
    string claimUrl;
    # Attribute to be matched
    string attribute;
};

public type ConditionalGroup record {
    # Description of the Conditional Group
    string description?;
    # Individual throttling conditions. They can be defined as either HeaderCondition, IPCondition, JWTClaimsCondition, QueryParameterCondition
    # Please see schemas of each of those throttling condition in Definitions section.
    ThrottleCondition[] conditions;
    ThrottleLimit 'limit;
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

public type CustomurlinfoDevportal record {
    string url?;
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

public type QueryParameterCondition record {
    # Name of the query parameter
    string parameterName;
    # Value of the query parameter to be matched
    string parameterValue;
};

# Blocking Conditions Status
public type BlockingConditionStatus record {
    # Id of the blocking condition
    string conditionId?;
    # Status of the blocking condition
    boolean conditionStatus;
};

public type SubscriptionThrottlePolicyPermission record {
    string permissionType;
    string[] roles;
};

public type KeyManagerEndpoint record {
    string name;
    string value;
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

public type RoleAliasList record {
    # The number of role aliases
    int count?;
    RoleAlias[] list?;
};

public type ApplicationList record {
    # Number of applications returned.
    int count?;
    ApplicationInfo[] list?;
    Pagination pagination?;
};

public type BotDetectionDataList record {
    # Number of Bot Detection Data returned.
    int count?;
    BotDetectionData[] list?;
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
    # The type of the tokens to be used (exchanged or without exchanged). Accepted values are EXCHANGED, DIRECT and BOTH.
    string tokenType = "DIRECT";
};

public type KeymanagerCertificates record {
    string 'type?;
    string value?;
};

public type PoliciesImportBody record {
    # Json File
    string file;
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
    # Workflow external reference is used to identify the workflow requests uniquely.
    string referenceId?;
    record {} properties?;
    # description is a message with basic details about the workflow request.
    string description?;
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

public type AdvancedThrottlePolicy record {
    *ThrottlePolicy;
    ThrottleLimit defaultLimit;
    # Group of conditions which allow adding different parameter conditions to the throttling limit.
    ConditionalGroup[] conditionalGroups?;
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

public type IPCondition record {
    # Type of the IP condition. Allowed values are "IPRANGE" and "IPSPECIFIC"
    string ipConditionType?;
    # Specific IP when "IPSPECIFIC" is used as the ipConditionType
    string specificIP?;
    # Staring IP when "IPRANGE" is used as the ipConditionType
    string startingIP?;
    # Ending IP when "IPRANGE" is used as the ipConditionType
    string endingIP?;
};

public type AdvancedThrottlePolicyInfo record {
    *ThrottlePolicy;
    ThrottleLimit defaultLimit?;
};

public type Scope record {
    # Portal name.
    string tag?;
    # Scope name.
    string name?;
    # About scope.
    string description?;
    # Roles for the particular scope.
    string[] roles?;
};

# The custom url information of the tenant domain
public type CustomUrlInfo record {
    string tenantDomain?;
    string tenantAdminUsername?;
    boolean enabled?;
    CustomurlinfoDevportal devPortal?;
};

public type WorkflowResponse record {
    # This attribute declares whether this workflow task is approved or rejected.
    string workflowStatus;
    # Attributes that returned after the workflow execution
    string jsonPayload?;
};

public type MonetizationInfo record {
    # Flag to indicate the monetization plan
    string monetizationPlan?;
    # Map of custom properties related to each monetization plan
    record {} properties;
};

public type AdvancedThrottlePolicyList record {
    # Number of Advanced Throttling Policies returned.
    int count?;
    AdvancedThrottlePolicyInfo[] list?;
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

public type ExportThrottlePolicy record {
    string 'type?;
    string subtype?;
    string 'version?;
    record {} data?;
};

public type BandwidthLimit record {
    *ThrottleLimitBase;
    # Amount of data allowed to be transfered
    int dataAmount;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit;
};
