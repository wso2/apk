package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.*;


import java.util.Objects;



public class SubscriptionDTO   {
  
  private String subscriptionId;

  private ApplicationInfoDTO applicationInfo;

  private String usagePlan;


public enum SubscriptionStatusEnum {

    BLOCKED(String.valueOf("BLOCKED")), PROD_ONLY_BLOCKED(String.valueOf("PROD_ONLY_BLOCKED")), UNBLOCKED(String.valueOf("UNBLOCKED")), ON_HOLD(String.valueOf("ON_HOLD")), REJECTED(String.valueOf("REJECTED")), TIER_UPDATE_PENDING(String.valueOf("TIER_UPDATE_PENDING")), DELETE_PENDING(String.valueOf("DELETE_PENDING"));


    private String value;

    SubscriptionStatusEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static SubscriptionStatusEnum fromValue(String value) {
        for (SubscriptionStatusEnum b : SubscriptionStatusEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private SubscriptionStatusEnum subscriptionStatus;


  /**
   **/
  public SubscriptionDTO subscriptionId(String subscriptionId) {
    this.subscriptionId = subscriptionId;
    return this;
  }

  
  @ApiModelProperty(example = "01234567-0123-0123-0123-012345678901", required = true, value = "")
  @JsonProperty("subscriptionId")
  @NotNull
  public String getSubscriptionId() {
    return subscriptionId;
  }
  public void setSubscriptionId(String subscriptionId) {
    this.subscriptionId = subscriptionId;
  }


  /**
   **/
  public SubscriptionDTO applicationInfo(ApplicationInfoDTO applicationInfo) {
    this.applicationInfo = applicationInfo;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "")
  @JsonProperty("applicationInfo")
  @NotNull
  public ApplicationInfoDTO getApplicationInfo() {
    return applicationInfo;
  }
  public void setApplicationInfo(ApplicationInfoDTO applicationInfo) {
    this.applicationInfo = applicationInfo;
  }


  /**
   **/
  public SubscriptionDTO usagePlan(String usagePlan) {
    this.usagePlan = usagePlan;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", required = true, value = "")
  @JsonProperty("usagePlan")
  @NotNull
  public String getUsagePlan() {
    return usagePlan;
  }
  public void setUsagePlan(String usagePlan) {
    this.usagePlan = usagePlan;
  }


  /**
   **/
  public SubscriptionDTO subscriptionStatus(SubscriptionStatusEnum subscriptionStatus) {
    this.subscriptionStatus = subscriptionStatus;
    return this;
  }

  
  @ApiModelProperty(example = "BLOCKED", required = true, value = "")
  @JsonProperty("subscriptionStatus")
  @NotNull
  public SubscriptionStatusEnum getSubscriptionStatus() {
    return subscriptionStatus;
  }
  public void setSubscriptionStatus(SubscriptionStatusEnum subscriptionStatus) {
    this.subscriptionStatus = subscriptionStatus;
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
        Objects.equals(applicationInfo, subscription.applicationInfo) &&
        Objects.equals(usagePlan, subscription.usagePlan) &&
        Objects.equals(subscriptionStatus, subscription.subscriptionStatus);
  }

  @Override
  public int hashCode() {
    return Objects.hash(subscriptionId, applicationInfo, usagePlan, subscriptionStatus);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SubscriptionDTO {\n");
    
    sb.append("    subscriptionId: ").append(toIndentedString(subscriptionId)).append("\n");
    sb.append("    applicationInfo: ").append(toIndentedString(applicationInfo)).append("\n");
    sb.append("    usagePlan: ").append(toIndentedString(usagePlan)).append("\n");
    sb.append("    subscriptionStatus: ").append(toIndentedString(subscriptionStatus)).append("\n");
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

