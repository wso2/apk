package org.wso2.apk.enforcer.analytics.publisher.reporter.prometheus;

import java.util.Map;

public class APIInvocationEvent {
    private String apiName;
    private int proxyResponseCode;
    private String destination;
    private String apiCreatorTenantDomain;
    private String platform;
    private String organizationId;
    private String apiMethod;
    private String apiVersion;
    private String gatewayType;
    private String environmentId;
    private String apiCreator;
    private boolean responseCacheHit;
    private int backendLatency;
    private String correlationId;
    private int requestMediationLatency;
    private String keyType;
    private String apiId;
    private String applicationName;
    private int targetResponseCode;
    private String requestTimestamp;
    private String applicationOwner;
    private String userAgent;
    private String userName;
    private String apiResourceTemplate;
    private String regionId;
    private int responseLatency;
    private int responseMediationLatency;
    private String userIp;
    private String apiContext;
    private String applicationId;
    private String apiType;
    private Map<String, String> properties;

    // Getters and setters for all fields
    public String getApiName() {
        return apiName;
    }

    public void setApiName(String apiName) {
        this.apiName = apiName;
    }

    public int getProxyResponseCode() {
        return proxyResponseCode;
    }

    public void setProxyResponseCode(int proxyResponseCode) {
        this.proxyResponseCode = proxyResponseCode;
    }

    public String getDestination() {
        return destination;
    }

    public void setDestination(String destination) {
        this.destination = destination;
    }

    public String getApiCreatorTenantDomain() {
        return apiCreatorTenantDomain;
    }

    public void setApiCreatorTenantDomain(String apiCreatorTenantDomain) {
        this.apiCreatorTenantDomain = apiCreatorTenantDomain;
    }

    public String getPlatform() {
        return platform;
    }

    public void setPlatform(String platform) {
        this.platform = platform;
    }

    public String getOrganizationId() {
        return organizationId;
    }

    public void setOrganizationId(String organizationId) {
        this.organizationId = organizationId;
    }

    public String getApiMethod() {
        return apiMethod;
    }

    public void setApiMethod(String apiMethod) {
        this.apiMethod = apiMethod;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getGatewayType() {
        return gatewayType;
    }

    public void setGatewayType(String gatewayType) {
        this.gatewayType = gatewayType;
    }

    public String getEnvironmentId() {
        return environmentId;
    }

    public void setEnvironmentId(String environmentId) {
        this.environmentId = environmentId;
    }

    public String getApiCreator() {
        return apiCreator;
    }

    public void setApiCreator(String apiCreator) {
        this.apiCreator = apiCreator;
    }

    public boolean isResponseCacheHit() {
        return responseCacheHit;
    }

    public void setResponseCacheHit(boolean responseCacheHit) {
        this.responseCacheHit = responseCacheHit;
    }

    public int getBackendLatency() {
        return backendLatency;
    }

    public void setBackendLatency(int backendLatency) {
        this.backendLatency = backendLatency;
    }

    public String getCorrelationId() {
        return correlationId;
    }

    public void setCorrelationId(String correlationId) {
        this.correlationId = correlationId;
    }

    public int getRequestMediationLatency() {
        return requestMediationLatency;
    }

    public void setRequestMediationLatency(int requestMediationLatency) {
        this.requestMediationLatency = requestMediationLatency;
    }

    public String getKeyType() {
        return keyType;
    }

    public void setKeyType(String keyType) {
        this.keyType = keyType;
    }

    public String getApiId() {
        return apiId;
    }

    public void setApiId(String apiId) {
        this.apiId = apiId;
    }

    public String getApplicationName() {
        return applicationName;
    }

    public void setApplicationName(String applicationName) {
        this.applicationName = applicationName;
    }

    public int getTargetResponseCode() {
        return targetResponseCode;
    }

    public void setTargetResponseCode(int targetResponseCode) {
        this.targetResponseCode = targetResponseCode;
    }

    public String getRequestTimestamp() {
        return requestTimestamp;
    }

    public void setRequestTimestamp(String requestTimestamp) {
        this.requestTimestamp = requestTimestamp;
    }

    public String getApplicationOwner() {
        return applicationOwner;
    }

    public void setApplicationOwner(String applicationOwner) {
        this.applicationOwner = applicationOwner;
    }

    public String getUserAgent() {
        return userAgent;
    }

    public void setUserAgent(String userAgent) {
        this.userAgent = userAgent;
    }

    public String getUserName() {
        return userName;
    }

    public void setUserName(String userName) {
        this.userName = userName;
    }

    public String getApiResourceTemplate() {
        return apiResourceTemplate;
    }

    public void setApiResourceTemplate(String apiResourceTemplate) {
        this.apiResourceTemplate = apiResourceTemplate;
    }

    public String getRegionId() {
        return regionId;
    }

    public void setRegionId(String regionId) {
        this.regionId = regionId;
    }

    public int getResponseLatency() {
        return responseLatency;
    }

    public void setResponseLatency(int responseLatency) {
        this.responseLatency = responseLatency;
    }

    public int getResponseMediationLatency() {
        return responseMediationLatency;
    }

    public void setResponseMediationLatency(int responseMediationLatency) {
        this.responseMediationLatency = responseMediationLatency;
    }

    public String getUserIp() {
        return userIp;
    }

    public void setUserIp(String userIp) {
        this.userIp = userIp;
    }

    public String getApiContext() {
        return apiContext;
    }

    public void setApiContext(String apiContext) {
        this.apiContext = apiContext;
    }

    public String getApplicationId() {
        return applicationId;
    }

    public void setApplicationId(String applicationId) {
        this.applicationId = applicationId;
    }

    public String getApiType() {
        return apiType;
    }

    public void setApiType(String apiType) {
        this.apiType = apiType;
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    public APIInvocationEvent() {
    }
}
