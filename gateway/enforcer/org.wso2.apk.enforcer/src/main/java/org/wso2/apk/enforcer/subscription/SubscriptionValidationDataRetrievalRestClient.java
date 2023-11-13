package org.wso2.apk.enforcer.subscription;
import feign.Headers;
import feign.Param;
import feign.RequestLine;

public interface SubscriptionValidationDataRetrievalRestClient {

    @RequestLine("GET /applications")
    @Headers("Content-Type: application/json")
    ApplicationListDto getAllApplications();

    @RequestLine("GET /subscriptions")
    @Headers("Content-Type: application/json")
    SubscriptionListDto getAllSubscriptions();

    @RequestLine("GET /applicationmappings")
    @Headers("Content-Type: application/json")
    ApplicationMappingDtoList getAllApplicationMappings();

    @RequestLine("GET /applicationkeymappings")
    @Headers("Content-Type: application/json")
    ApplicationKeyMappingDtoList getAllApplicationKeyMappings();
}
