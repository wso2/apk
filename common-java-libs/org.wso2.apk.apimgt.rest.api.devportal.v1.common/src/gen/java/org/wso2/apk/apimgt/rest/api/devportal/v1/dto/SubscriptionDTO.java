package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APIInfoDTO;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationInfoDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class SubscriptionDTO   {
  
  private String subscriptionId;

  private String applicationId;

  private String apiId;

  private APIInfoDTO apiInfo;

  private ApplicationInfoDTO applicationInfo;

  private String throttlingPolicy;

  private String requestedThrottlingPolicy;


public enum StatusEnum {

    BLOCKED(String.valueOf("BLOCKED")), PROD_ONLY_BLOCKED(String.valueOf("PROD_ONLY_BLOCKED")), UNBLOCKED(String.valueOf("UNBLOCKED")), ON_HOLD(String.valueOf("ON_HOLD")), REJECTED(String.valueOf("REJECTED")), TIER_UPDATE_PENDING(String.valueOf("TIER_UPDATE_PENDING")), DELETE_PENDING(String.valueOf("DELETE_PENDING"));


    private String value;

    StatusEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static StatusEnum fromValue(String value) {
        for (StatusEnum b : StatusEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private StatusEnum status;

  private String redirectionParams;


  /**
   * The UUID of the subscription
   **/
  public SubscriptionDTO subscriptionId(String subscriptionId) {
    this.subscriptionId = subscriptionId;
    return this;
  }

  
  @ApiModelProperty(example = "faae5fcc-cbae-40c4-bf43-89931630d313", value = "The UUID of the subscription")
  @JsonProperty("subscriptionId")
  public String getSubscriptionId() {
    return subscriptionId;
  }
  public void setSubscriptionId(String subscriptionId) {
    this.subscriptionId = subscriptionId;
  }


  /**
   * The UUID of the application
   **/
  public SubscriptionDTO applicationId(String applicationId) {
    this.applicationId = applicationId;
    return this;
  }

  
  @ApiModelProperty(example = "b3ade481-30b0-4b38-9a67-498a40873a6d", required = true, value = "The UUID of the application")
  @JsonProperty("applicationId")
  @NotNull
  public String getApplicationId() {
    return applicationId;
  }
  public void setApplicationId(String applicationId) {
    this.applicationId = applicationId;
  }


  /**
   * The unique identifier of the API.
   **/
  public SubscriptionDTO apiId(String apiId) {
    this.apiId = apiId;
    return this;
  }

  
  @ApiModelProperty(example = "2962f3bb-8330-438e-baee-0ee1d6434ba4", value = "The unique identifier of the API.")
  @JsonProperty("apiId")
  public String getApiId() {
    return apiId;
  }
  public void setApiId(String apiId) {
    this.apiId = apiId;
  }


  /**
   **/
  public SubscriptionDTO apiInfo(APIInfoDTO apiInfo) {
    this.apiInfo = apiInfo;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("apiInfo")
  public APIInfoDTO getApiInfo() {
    return apiInfo;
  }
  public void setApiInfo(APIInfoDTO apiInfo) {
    this.apiInfo = apiInfo;
  }


  /**
   **/
  public SubscriptionDTO applicationInfo(ApplicationInfoDTO applicationInfo) {
    this.applicationInfo = applicationInfo;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("applicationInfo")
  public ApplicationInfoDTO getApplicationInfo() {
    return applicationInfo;
  }
  public void setApplicationInfo(ApplicationInfoDTO applicationInfo) {
    this.applicationInfo = applicationInfo;
  }


  /**
   **/
  public SubscriptionDTO throttlingPolicy(String throttlingPolicy) {
    this.throttlingPolicy = throttlingPolicy;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", required = true, value = "")
  @JsonProperty("throttlingPolicy")
  @NotNull
  public String getThrottlingPolicy() {
    return throttlingPolicy;
  }
  public void setThrottlingPolicy(String throttlingPolicy) {
    this.throttlingPolicy = throttlingPolicy;
  }


  /**
   **/
  public SubscriptionDTO requestedThrottlingPolicy(String requestedThrottlingPolicy) {
    this.requestedThrottlingPolicy = requestedThrottlingPolicy;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", value = "")
  @JsonProperty("requestedThrottlingPolicy")
  public String getRequestedThrottlingPolicy() {
    return requestedThrottlingPolicy;
  }
  public void setRequestedThrottlingPolicy(String requestedThrottlingPolicy) {
    this.requestedThrottlingPolicy = requestedThrottlingPolicy;
  }


  /**
   **/
  public SubscriptionDTO status(StatusEnum status) {
    this.status = status;
    return this;
  }

  
  @ApiModelProperty(example = "UNBLOCKED", value = "")
  @JsonProperty("status")
  public StatusEnum getStatus() {
    return status;
  }
  public void setStatus(StatusEnum status) {
    this.status = status;
  }


  /**
   * A url and other parameters the subscriber can be redirected.
   **/
  public SubscriptionDTO redirectionParams(String redirectionParams) {
    this.redirectionParams = redirectionParams;
    return this;
  }

  
  @ApiModelProperty(example = "", value = "A url and other parameters the subscriber can be redirected.")
  @JsonProperty("redirectionParams")
  public String getRedirectionParams() {
    return redirectionParams;
  }
  public void setRedirectionParams(String redirectionParams) {
    this.redirectionParams = redirectionParams;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    SubscriptionDTO subscription = (SubscriptionDTO) o;
    return Objects.equals(subscriptionId, subscription.subscriptionId) &&
        Objects.equals(applicationId, subscription.applicationId) &&
        Objects.equals(apiId, subscription.apiId) &&
        Objects.equals(apiInfo, subscription.apiInfo) &&
        Objects.equals(applicationInfo, subscription.applicationInfo) &&
        Objects.equals(throttlingPolicy, subscription.throttlingPolicy) &&
        Objects.equals(requestedThrottlingPolicy, subscription.requestedThrottlingPolicy) &&
        Objects.equals(status, subscription.status) &&
        Objects.equals(redirectionParams, subscription.redirectionParams);
  }

  @Override
  public int hashCode() {
    return Objects.hash(subscriptionId, applicationId, apiId, apiInfo, applicationInfo, throttlingPolicy, requestedThrottlingPolicy, status, redirectionParams);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SubscriptionDTO {\n");
    
    sb.append("    subscriptionId: ").append(toIndentedString(subscriptionId)).append("\n");
    sb.append("    applicationId: ").append(toIndentedString(applicationId)).append("\n");
    sb.append("    apiId: ").append(toIndentedString(apiId)).append("\n");
    sb.append("    apiInfo: ").append(toIndentedString(apiInfo)).append("\n");
    sb.append("    applicationInfo: ").append(toIndentedString(applicationInfo)).append("\n");
    sb.append("    throttlingPolicy: ").append(toIndentedString(throttlingPolicy)).append("\n");
    sb.append("    requestedThrottlingPolicy: ").append(toIndentedString(requestedThrottlingPolicy)).append("\n");
    sb.append("    status: ").append(toIndentedString(status)).append("\n");
    sb.append("    redirectionParams: ").append(toIndentedString(redirectionParams)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }
}

