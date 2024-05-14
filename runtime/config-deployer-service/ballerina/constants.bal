const string DEFAULT_PARTITION = "default";
const string ALL_NAMESPACES = "*";
public final string[] & readonly HTTP_DEFAULT_METHODS = ["get", "put", "post", "delete", "patch"];
public final string[] SOAP_DEFAULT_METHODS = ["post"];
public final string[] SSE_DEFAULT_METHODS = ["get"];
public final string[] WS_DEFAULT_METHODS = ["post"];
public final string[] WEBSUB_DEFAULT_METHODS = ["post"];
public final string[] WEBSUB_SUPPORTED_METHODS = ["subscribe"];
public final string[] SSE_SUPPORTED_METHODS = ["subscribe"];
public final string[] WS_SUPPORTED_METHODS = ["subscribe", "publish"];

const string API_TYPE_REST = "REST";
const string API_TYPE_GRAPHQL = "GRAPHQL";
const string API_TYPE_GRPC = "GRPC";
const string API_TYPE_SOAP = "SOAP";
const string API_TYPE_SSE = "SSE";
const string API_TYPE_WS = "WS";
const string API_TYPE_WEBSUB = "WEBSUB";
const string APK_USER = "apkuser";
const string CURRENT_NAMESPACE = "CURRENT_NAME_SPACE";

const string SORT_BY_API_NAME = "apiName";
const string SORT_BY_POLICY_NAME = "policyName";
const string SORT_BY_ID = "id";
const string SORT_BY_CREATED_TIME = "createdTime";
const string SORT_ORDER_ASC = "asc";
const string SORT_ORDER_DESC = "desc";
const string SEARCH_CRITERIA_NAME = "name";
const string SEARCH_CRITERIA_TYPE = "type";
const string SEARCH_CRITERIA_NAMESPACE = "namespace";
const string SORT_BY_SERVICE_NAME = "serviceName";
const string SORT_BY_SERVICE_CREATED_TIME = "createdTime";

const string BASE_PATH_ALREADY_EXIST_K8s_VALIDATION_MESSAGE = "an API has been already created for the basePath";
const string PRODUCTION_TYPE = "production";
const string SANDBOX_TYPE = "sandbox";
const string URL = "url";
const string ENDPOINT_BASIC_USER_NAME = "username";
const string ENDPOINT_BASIC_PASSWORD = "password";
const string ENDPOINT_BASIC_SECRET_REF = "secretRef";

const string INTERCEPTOR_TYPE = "interceptor";
const string PRIMARY_ENDPOINT = "primary";
const string ENDPOINT_SECURITY_TYPE_BASIC = "basic";
const string ENDPOINT_SECURITY_TYPE_BASIC_CASE = "Basic";
const string DEFAULT_MODIFIED_ENDPOINT_PASSWORD = "*****"; //5 stars
const string ENDPOINT_SECURITY_FIELD = "endpoint_security";
const string ENDPOINT_SECURITY_TYPE = "type";
const string ENDPOINT_SECURITY_SECRET_REFERENCE_NAME = "secretRefName";
const string ENDPOINT_SECURITY_GENERATED_SECRET_REFERENCE_NAME = "generatedSecretRefName";
const string ENDPOINT_SECURITY_USERNAME = "username";
const string ENDPOINT_SECURITY_PASSWORD = "password";
const string ZIP_FILE_EXTENSTION = ".zip";
const string PROTOCOL_HTTP = "http";
const string PROTOCOL_HTTPS = "https";
final string[] & readonly ALLOWED_API_TYPES = [API_TYPE_REST, API_TYPE_GRAPHQL, API_TYPE_GRPC];

const string MEDIATION_POLICY_TYPE_REQUEST_HEADER_MODIFIER = "RequestHeaderModifier";
const string MEDIATION_POLICY_TYPE_RESPONSE_HEADER_MODIFIER = "ResponseHeaderModifier";
const string MEDIATION_POLICY_TYPE_INTERCEPTOR = "Interceptor";
const string POLICY_TYPE_BACKEND_JWT = "BackendJwt";
const string MEDIATION_POLICY_NAME_ADD_HEADER = "addHeader";
const string MEDIATION_POLICY_NAME_REMOVE_HEADER = "removeHeader";
const string MEDIATION_POLICY_TYPE_URL_REWRITE = "URLRewrite";
const string MEDIATION_POLICY_FLOW_REQUEST = "request";
const string MEDIATION_POLICY_FLOW_RESPONSE = "response";

const string API_NAME_HASH_LABEL = "api-name";
const string API_VERSION_HASH_LABEL = "api-version";
const string ORGANIZATION_HASH_LABEL = "organization";
const string CONFIG_TYPE_LABEL = "config-type";
const string MANAGED_BY_HASH_LABEL = "managed-by";
const string MANAGED_BY_HASH_LABEL_VALUE = "apk";

const string CERTIFICATE_VERSION_NUMBER = "wso2apk/certificate-version";
const string CERTIFICATE_SERIAL_NUMBER = "wso2apk/certificate-serial-number";
const string CERTIFICATE_ISSUER = "wso2apk/certificate-issuer";
const string CERTIFICATE_SUBJECT = "wso2apk/certificate-subject";
const string CERTIFICATE_NOT_BEFORE = "wso2apk/certificate-not-before";
const string CERTIFICATE_NOT_AFTER = "wso2apk/certificate-not-after";
const string CERTIFICATE_HOSTS = "wso2apk/certificate-host";
const string API_UUID_ANNOTATION = "wso2apk/api-uuid";
const string CERTIFICATE_KEY_CONFIG_MAP = "endoint.crt";
const string CONFIG_TYPE_LABEL_VALUE = "certificate";
const string CONFIGMAP_DEFINITION_KEY = "definition";

const string ORG_RESOLVER_CONTROL_PLANE = "controlPlane";
const string ORG_RESOLVER_K8s = "k8s";
const string ORG_RESOLVER_NONE = "none";

const string APPLICATION_JSON_MEDIA_TYPE = "application/json";
const string APPLICATION_YAML_MEDIA_TYPE = "application/yaml";
const string ALL_MEDIA_TYPE = "*/*";
const string AUTH_TYPE_JWT = "JWT";
