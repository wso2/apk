//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

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

public type CreatedSubscription record {|
    *http:Created;
    Subscription body;
|};

public type UnsupportedMediaTypeError record {|
    *http:UnsupportedMediaType;
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

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type UnauthorizedError record {|
    *http:Unauthorized;
    Error body;
|};

public type CreatedApplication record {|
    *http:Created;
    Application body;
|};

public type DocumentList record {
    # Number of Documents returned.
    int count?;
    Document[] list?;
    Pagination pagination?;
};

public type ApplicationKey record {
    # Key Manager Mapping UUID
    string keyMappingId?;
    # Key Manager Name
    string keyManager?;
    # Consumer key of the application
    string consumerKey?;
    # Consumer secret of the application
    string consumerSecret?;
    # The grant types that are supported by the application
    string[] supportedGrantTypes?;
    # Callback URL
    string callbackUrl?;
    # Describes the state of the key generation.
    string keyState?;
    # Describes to which endpoint the key belongs
    string keyType?;
    # Describe the which mode Application Mapped.
    string mode?;
    # Application group id (if any).
    string groupId?;
    ApplicationToken token?;
    # additionalProperties (if any).
    record {} additionalProperties?;
};

public type GraphQLSchemaTypeList record {
    GraphQLSchemaType[] typeList?;
};

public type GraphQLSchemaType record {
    # Type found within the GraphQL Schema
    string 'type?;
    # Array of fields under current type
    string[] fieldList?;
};

public type ApplicationKeyGenerateRequest record {
    string keyType;
    # key Manager to Generate Keys
    string keyManager?;
    # Grant types that should be supported by the application
    string[] grantTypesToBeSupported;
    # Callback URL
    string callbackUrl?;
    # Allowed scopes for the access token
    string[] scopes?;
    string validityTime?;
    # Client ID for generating access token.
    string clientId?;
    # Client secret for generating access token. This is given together with the client Id.
    string clientSecret?;
    # Additional properties needed.
    record {} additionalProperties?;
};

public type Document record {
    string documentId?;
    string name;
    string 'type;
    string summary?;
    string sourceType;
    string sourceUrl?;
    string otherTypeName?;
};

public type APIKeyRevokeRequest record {
    # API Key to revoke
    string apikey?;
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

public type ApiTiers record {
    string tierName?;
    string tierPlan?;
    ApiMonetizationattributes monetizationAttributes?;
};

public type AdditionalsubscriptioninfoSolacedeployedenvironments record {
    string environmentName?;
    string environmentDisplayName?;
    string organizationName?;
    AdditionalsubscriptioninfoSolaceurls[] solaceURLs?;
    AdditionalsubscriptioninfoSolacetopicsobject SolaceTopicsObject?;
};

public type RatingList record {
    # Average Rating of the API
    string avgRating?;
    # Rating given by the user
    int userRating?;
    # Number of Subscriber Ratings returned.
    int count?;
    Rating[] list?;
    Pagination pagination?;
};

public type APIOperations record {
    string id?;
    string target?;
    string verb?;
};

public type ScopeList record {
    # Number of results returned.
    int count?;
    ScopeInfo[] list?;
    Pagination pagination?;
};

public type APIList record {
    # Number of APIs returned.
    int count?;
    APIInfo[] list?;
    Pagination pagination?;
};

public type CurrentAndNewPasswords record {
    string currentPassword?;
    string newPassword?;
};

public type ApplicationToken record {
    # Access token
    string accessToken?;
    # Valid comma seperated scopes for the access token
    string[] tokenScopes?;
    # Maximum validity time for the access token
    int validityTime?;
};

public type APIMonetizationUsage record {
    # Map of custom properties related to monetization usage
    record {} properties?;
};

public type Subscription record {
    # The UUID of the subscription
    string subscriptionId?;
    # The UUID of the application
    string applicationId;
    # The unique identifier of the API.
    string apiId?;
    APIInfo apiInfo?;
    ApplicationInfo applicationInfo?;
    string throttlingPolicy;
    string requestedThrottlingPolicy?;
    string status?;
    # A url and other parameters the subscriber can be redirected.
    string redirectionParams?;
};

public type ApplicationsImportBody record {
    # Zip archive consisting of exported Application Configuration.
    string file;
};

public type Settings record {
    string[] grantTypes?;
    string[] scopes?;
    boolean applicationSharingEnabled?;
    boolean mapExistingAuthApps?;
    string apiGatewayEndpoint?;
    boolean monetizationEnabled?;
    boolean recommendationEnabled?;
    boolean IsUnlimitedTierPaid?;
    SettingsIdentityprovider identityProvider?;
    boolean IsAnonymousModeEnabled?;
    boolean IsPasswordChangeEnabled?;
    # The 'PasswordJavaRegEx' cofigured in the UserStoreManager
    string userStorePasswordPattern?;
    # The regex configured in the Password Policy property 'passwordPolicy.pattern'
    string passwordPolicyPattern?;
    # If Password Policy Feature is enabled, the property 'passwordPolicy.min.length' is returned as the 'passwordPolicyMinLength'. If password policy is not enabled, default value -1 will be returned. And it should be noted that the regex pattern(s) returned in 'passwordPolicyPattern' and 'userStorePasswordPattern' properties too will affect the minimum password length allowed and an intersection of all conditions will be considered finally to validate the password.
    int passwordPolicyMinLength?;
    # If Password Policy Feature is enabled, the property 'passwordPolicy.max.length' is returned as the 'passwordPolicyMaxLength'. If password policy is not enabled, default value -1 will be returned. And it should be noted that the regex pattern(s) returned in 'passwordPolicyPattern' and 'userStorePasswordPattern' properties too will affect the maximum password length allowed and an intersection of all conditions will be considered finally to validate the password.
    int passwordPolicyMaxLength?;
};

public type ApiinfoAdditionalproperties record {
    string name?;
    string value?;
    boolean display?;
};

public type ThrottlingPolicyPermissionInfo record {
    string 'type?;
    # roles for this permission
    string[] roles?;
};

public type ApiDefaultversionurls record {
    # HTTP environment default URL
    string http?;
    # HTTPS environment default URL
    string https?;
    # WS environment default URL
    string ws?;
    # WSS environment default URL
    string wss?;
};

public type APIMonetizationInfo record {
    # Flag to indicate the monetization status
    boolean enabled;
};

public type APIBusinessInformation record {
    string businessOwner?;
    string businessOwnerEmail?;
    string technicalOwner?;
    string technicalOwnerEmail?;
};

public type SolaceTopics record {
    string[] publishTopics?;
    string[] subscribeTopics?;
};

public type APISearchResult record {
    *SearchResult;
    # A brief description about the API
    string description?;
    # A string that represents the context of the user's request
    string context?;
    # The version of the API
    string 'version?;
    # If the provider value is notgiven, the user invoking the API will be used as the provider.
    string provider?;
    # This describes in which status of the lifecycle the API is
    string status?;
    string thumbnailUri?;
    APIBusinessInformation businessInformation?;
    # Average rating of the API
    string avgRating?;
};

public type ApplicationAttributeList record {
    # Number of application attributes returned.
    int count?;
    ApplicationAttribute[] list?;
};

public type SettingsIdentityprovider record {
    boolean 'external?;
};

public type ThrottlingPolicyList record {
    # Number of Throttling Policies returned.
    int count?;
    ThrottlingPolicy[] list?;
    Pagination pagination?;
};

public type ApiEndpointurls record {
    string environmentName?;
    string environmentDisplayName?;
    string environmentType?;
    ApiUrls URLs?;
    ApiDefaultversionurls defaultVersionURLs?;
};

public type AdditionalsubscriptioninfoSolacetopicsobject record {
    SolaceTopics defaultSyntax?;
    SolaceTopics mqttSyntax?;
};

public type ApplicationKeyList record {
    # Number of applications keys returned.
    int count?;
    ApplicationKey[] list?;
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

public type AdvertiseInfo record {
    boolean advertised?;
    string apiExternalProductionEndpoint?;
    string apiExternalSandboxEndpoint?;
    string originalDevPortalUrl?;
    string apiOwner?;
    string vendor?;
};

public type CommentList record {
    # Number of Comments returned.
    int count?;
    Comment[] list?;
    Pagination pagination?;
};

public type Application record {
    string applicationId?;
    int id?;
    @constraint:String {maxLength: 100, minLength: 1}
    string name;
    @constraint:String {minLength: 1}
    string throttlingPolicy;
    @constraint:String {maxLength: 512}
    string description?;
    # Type of the access token generated for this application.
    # 
    # **OAUTH:** A UUID based access token
    # **JWT:** A self-contained, signed JWT based access token which is issued by default.
    string tokenType = "JWT";
    string status = "";
    string[] groups?;
    int subscriptionCount?;
    ApplicationKey[] keys?;
    record {} attributes?;
    ScopeInfo[] subscriptionScopes?;
    # Application created user
    string owner?;
    boolean hashEnabled?;
    string createdTime?;
    string updatedTime?;
};

public type GraphQLCustomComplexityInfo record {
    # The type found within the schema of the API
    string 'type;
    # The field which is found under the type within the schema of the API
    string 'field;
    # The complexity value allocated for the associated field under the specified type
    int complexityValue;
};

public type AdditionalSubscriptionInfoList record {
    # Number of additional information sets of subscription returned.
    int count?;
    AdditionalSubscriptionInfo[] list?;
    Pagination pagination?;
};

public type ApiMonetizationattributes record {
    string fixedPrice?;
    string pricePerRequest?;
    string currencyType?;
    string billingCycle?;
};

public type User record {
    string username;
    string password;
    string firstName;
    string lastName;
    string email;
};

public type ErrorListItem record {
    string code;
    # Description about individual errors occurred
    string message;
};

public type Rating record {
    string ratingId?;
    string apiId?;
    @constraint:String {maxLength: 50}
    string ratedBy?;
    int rating;
};

public type APICategoryList record {
    # Number of API categories returned.
    int count?;
    APICategory[] list?;
};

public type ApplicationInfo record {
    string applicationId?;
    string name?;
    string throttlingPolicy?;
    string description?;
    string status?;
    string[] groups?;
    int subscriptionCount?;
    record {} attributes?;
    string owner?;
    string createdTime?;
    string updatedTime?;
};

public type ApplicationKeyReGenerateResponse record {
    # The consumer key associated with the application, used to indetify the client
    string consumerKey?;
    # The client secret that is used to authenticate the client with the authentication server
    string consumerSecret?;
};

public type KeyManagerList record {
    # Number of Key managers returned.
    int count?;
    KeyManagerInfo[] list?;
};

public type ApplicationAttribute record {
    # description of the application attribute
    string description?;
    # type of the input element to display
    string 'type?;
    # tooltop to display for the input element
    string tooltip?;
    # whether this is a required attribute
    string required?;
    # the name of the attribute
    string attribute?;
    # whether this is a hidden attribute
    string hidden?;
};

public type Recommendations record {
    # Number of APIs returned.
    int count?;
    RecommendedAPI[] list?;
};

public type SearchResultList record {
    # Number of results returned.
    int count?;
    record {}[] list?;
    Pagination pagination?;
};

public type Tenant record {
    # tenant domain
    string domain?;
    # current status of the tenant active/inactive
    string status?;
};

public type KeyManagerApplicationConfiguration record {
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

public type ApplicationList record {
    # Number of applications returned.
    int count?;
    ApplicationInfo[] list?;
    Pagination pagination?;
};

public type TopicList record {
    # Number of Topics returned.
    int count?;
    Topic[] list?;
    Pagination pagination?;
};

public type TagList record {
    # Number of Tags returned.
    int count?;
    Tag[] list?;
    Pagination pagination?;
};

public type TenantList record {
    # Number of tenants returned.
    int count?;
    Tenant[] list?;
    Pagination pagination?;
};

public type ApplicationKeyMappingRequest record {
    # Consumer key of the application
    string consumerKey;
    # Consumer secret of the application
    string consumerSecret?;
    # Key Manager Name
    string keyManager?;
    string keyType;
};

public type AdditionalsubscriptioninfoSolaceurls record {
    string protocol?;
    string endpointURL?;
};

public type Topic record {
    string apiId?;
    string name?;
    string 'type?;
};

public type WebhookSubscriptionList record {
    # Number of webhook subscriptions returned.
    int count?;
    WebhookSubscription[] list?;
    Pagination pagination?;
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

public type ApplicationTokenGenerateRequest record {
    # Consumer secret of the application
    string consumerSecret?;
    # Token validity period
    int validityPeriod?;
    # Allowed scopes (space seperated) for the access token
    string[] scopes?;
    # Token to be revoked, if any
    string revokeToken?;
    string grantType?;
    # Additional parameters if Authorization server needs any
    record {} additionalProperties?;
};

public type KeyManagerInfo record {
    string id?;
    string name;
    string 'type;
    # display name of Keymanager
    string displayName?;
    string description?;
    boolean enabled?;
    string[] availableGrantTypes?;
    string tokenEndpoint?;
    string revokeEndpoint?;
    string userInfoEndpoint?;
    boolean enableTokenGeneration?;
    boolean enableTokenEncryption = false;
    boolean enableTokenHashing = false;
    boolean enableOAuthAppCreation = true;
    boolean enableMapOAuthConsumerApps = false;
    KeyManagerApplicationConfiguration[] applicationConfiguration?;
    # The alias of Identity Provider.
    # If the tokenType is EXCHANGED, the alias value should be inclusive in the audience values of the JWT token
    string alias?;
    record {} additionalProperties?;
    # The type of the tokens to be used (exchanged or without exchanged). Accepted values are EXCHANGED, DIRECT and BOTH.
    string tokenType = "DIRECT";
};

public type WebhookSubscription record {
    string apiId?;
    string appId?;
    string topic?;
    string callBackUrl?;
    string deliveryTime?;
    int deliveryStatus?;
};

public type APIInfo record {
    string id?;
    string name?;
    string description?;
    string context?;
    string 'version?;
    string 'type?;
    string createdTime?;
    # If the provider value is not given, the user invoking the API will be used as the provider.
    string provider?;
    string lifeCycleStatus?;
    string thumbnailUri?;
    # Average rating of the API
    string avgRating?;
    # List of throttling policies of the API
    string[] throttlingPolicies?;
    AdvertiseInfo advertiseInfo?;
    APIBusinessInformation businessInformation?;
    boolean isSubscriptionAvailable?;
    string monetizationLabel?;
    string gatewayVendor?;
    # Custom(user defined) properties of API
    ApiinfoAdditionalproperties[] additionalProperties?;
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

public type APIKeyGenerateRequest record {
    # Token validity period
    int validityPeriod?;
    # Additional parameters if Authorization server needs any
    record {} additionalProperties?;
};

public type ScopeInfo record {
    string 'key?;
    string name?;
    # Allowed roles for the scope
    string[] roles?;
    # Description of the scope
    string description?;
};

public type SearchResult record {
    string id?;
    string name;
    string 'type?;
    # Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL
    string transportType?;
};

public type ThrottlingPolicy record {
    string name;
    string description?;
    string policyLevel?;
    # Custom attributes added to the throttling policy
    record {} attributes?;
    # Maximum number of requests which can be sent within a provided unit time
    int requestCount;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit?;
    int unitTime;
    string timeUnit?;
    # Burst control request count
    int rateLimitCount = 0;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Default quota limit type
    string quotaPolicyType?;
    # This attribute declares whether this tier is available under commercial or free
    string tierPlan;
    # If this attribute is set to false, you are capabale of sending requests
    # even if the request count exceeded within a unit time
    boolean stopOnQuotaReach;
    MonetizationInfo monetizationAttributes?;
    ThrottlingPolicyPermissionInfo throttlingPolicyPermissions?;
};

public type APIKey record {
    # API Key
    string apikey?;
    int validityTime?;
};

public type DocumentSearchResult record {
    *SearchResult;
    string docType?;
    string summary?;
    string sourceType?;
    string sourceUrl?;
    string otherTypeName?;
    string visibility?;
    # The name of the associated API
    string apiName?;
    # The version of the associated API
    string apiVersion?;
    string apiProvider?;
    string apiUUID?;
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

public type GraphQLQueryComplexityInfo record {
    GraphQLCustomComplexityInfo[] list?;
};

public type RecommendedAPI record {
    string id?;
    string name?;
    # Average rating of the API
    string avgRating?;
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
};

public type API record {
    # UUID of the api
    string id?;
    # ID of the api
    int apiId?;
    # Name of the API
    string name;
    # A brief description about the API
    string description?;
    # A string that represents thecontext of the user's request
    string context;
    # The version of the API
    string 'version;
    # If the provider value is not given user invoking the api will be used as the provider.
    string provider;
    # Swagger definition of the API which contains details about URI templates and scopes
    string apiDefinition?;
    # WSDL URL if the API is based on a WSDL endpoint
    string wsdlUri?;
    # This describes in which status of the lifecycle the API is.
    string lifeCycleStatus;
    boolean isDefaultVersion?;
    # This describes the transport type of the API
    string 'type?;
    string[] transport?;
    APIOperations[] operations?;
    # Name of the Authorization header used for invoking the API. If it is not set, Authorization header name specified
    # in tenant or system level will be used.
    string authorizationHeader?;
    # Types of API security, the current API secured with. It can be either OAuth2 or mutual SSL or both. If
    # it is not set OAuth2 will be set as the security for the current API.
    string[] securityScheme?;
    # Search keywords related to the API
    string[] tags?;
    # The subscription tiers selected for the particular API
    ApiTiers[] tiers?;
    boolean hasThumbnail = false;
    # Custom(user defined) properties of API
    ApiinfoAdditionalproperties[] additionalProperties?;
    APIMonetizationInfo monetization?;
    ApiEndpointurls[] endpointURLs?;
    APIBusinessInformation businessInformation?;
    # The environment list configured with non empty endpoint URLs for the particular API.
    string[] environmentList?;
    ScopeInfo[] scopes?;
    # The average rating of the API
    string avgRating?;
    AdvertiseInfo advertiseInfo?;
    boolean isSubscriptionAvailable?;
    # API categories
    string[] categories?;
    # API Key Managers
    record {} keyManagers?;
    string createdTime?;
    string lastUpdatedTime?;
    string gatewayVendor?;
    # Supported transports for the aync API.
    string[] asyncTransportProtocols?;
};

public type Tag record {
    string value?;
    int count?;
};

public type MonetizationInfo record {
    string billingType?;
    string billingCycle?;
    string fixedPrice?;
    string pricePerRequest?;
    string currencyType?;
};

public type AdditionalSubscriptionInfo record {
    # The UUID of the subscription
    string subscriptionId?;
    # The UUID of the application
    string applicationId?;
    # The name of the application
    string applicationName?;
    # The unique identifier of the API.
    string apiId?;
    boolean isSolaceAPI?;
    string solaceOrganization?;
    AdditionalsubscriptioninfoSolacedeployedenvironments[] solaceDeployedEnvironments?;
};

public type APIInfoList record {
    # Number of API Info objects returned.
    int count?;
    APIInfo[] list?;
};

public type APICategory record {
    string id?;
    string name;
    string description?;
};

public type ApiUrls record {
    # HTTP environment URL
    string http?;
    # HTTPS environment URL
    string https?;
    # WS environment URL
    string ws?;
    # WSS environment URL
    string wss?;
};
