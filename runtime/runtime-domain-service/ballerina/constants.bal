const string ALL_NAMESPACES = "*";
public final string[] & readonly HTTP_DEFAULT_METHODS = ["get", "put", "post", "delete", "patch"];
public final string[] SOAP_DEFAULT_METHODS = ["post"];
public final string[] SSE_DEFAULT_METHODS = ["get"];
public final string[] WS_DEFAULT_METHODS = ["post"];
public final string[] WEBSUB_DEFAULT_METHODS = ["post"];
public final string[] WEBSUB_SUPPORTED_METHODS = ["subscribe"];
public final string[] SSE_SUPPORTED_METHODS = ["subscribe"];
public final string[] WS_SUPPORTED_METHODS = ["subscribe", "publish"];

const string API_TYPE_HTTP = "HTTP";
const string API_TYPE_SOAP = "SOAP";
const string API_TYPE_SSE = "SSE";
const string API_TYPE_WS = "WS";
const string API_TYPE_WEBSUB = "WEBSUB";
const string APK_USER = "apkuser";
const string CURRENT_NAMESPACE = "CURRENT_NAME_SPACE";

const string SORT_BY_API_NAME = "apiName";
const string SORT_BY_CREATED_TIME = "createdTime";
const string SORT_ORDER_ASC = "asc";
const string SORT_ORDER_DESC = "desc";
const string SEARCH_CRITERIA_NAME = "name";
const string SEARCH_CRITERIA_TYPE = "type";
const string SORT_BY_SERVICE_NAME = "serviceName";
const string SORT_BY_SERVICE_CREATED_TIME = "createdTime";

const string CONTEXT_ALREADY_EXIST_K8s_VALIDATION_MESSAGE = "an API has been already created for the context";
