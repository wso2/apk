package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import javax.validation.constraints.*;


import java.util.Objects;



public class UsagePlanDTO   {
  
  private Integer policyId;

  private String uuid;

  private String policyName;

  private String displayName;

  private String description;

  private String organization;

  private UsageLimitDTO defaultLimit;

  private Integer rateLimitCount;

  private String rateLimitTimeUnit;

  private Integer subscriberCount;

  private List<CustomAttributeDTO> customAttributes = null;

  private Boolean stopOnQuotaReach = false;

  private String billingPlan;

  private SubscriptionThrottlePolicyPermissionDTO permissions;


  /**
   * Id of policy
   **/
  public UsagePlanDTO policyId(Integer policyId) {
    this.policyId = policyId;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "Id of policy")
  @JsonProperty("policyId")
  public Integer getPolicyId() {
    return policyId;
  }
  public void setPolicyId(Integer policyId) {
    this.policyId = policyId;
  }


  /**
   * policy uuid
   **/
  public UsagePlanDTO uuid(String uuid) {
    this.uuid = uuid;
    return this;
  }

  
  @ApiModelProperty(example = "0c6439fd-9b16-3c2e-be6e-1086e0b9aa93", value = "policy uuid")
  @JsonProperty("uuid")
  public String getUuid() {
    return uuid;
  }
  public void setUuid(String uuid) {
    this.uuid = uuid;
  }


  /**
   * Name of policy
   **/
  public UsagePlanDTO policyName(String policyName) {
    this.policyName = policyName;
    return this;
  }

  
  @ApiModelProperty(example = "30PerMin", value = "Name of policy")
  @JsonProperty("policyName")
 @Size(min=1,max=60)  public String getPolicyName() {
    return policyName;
  }
  public void setPolicyName(String policyName) {
    this.policyName = policyName;
  }


  /**
   * Display name of the policy
   **/
  public UsagePlanDTO displayName(String displayName) {
    this.displayName = displayName;
    return this;
  }

  
  @ApiModelProperty(example = "30PerMin", value = "Display name of the policy")
  @JsonProperty("displayName")
 @Size(max=512)  public String getDisplayName() {
    return displayName;
  }
  public void setDisplayName(String displayName) {
    this.displayName = displayName;
  }


  /**
   * Description of the policy
   **/
  public UsagePlanDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "Allows 30 request per minute", value = "Description of the policy")
  @JsonProperty("description")
 @Size(max=1024)  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   * Usage policy organization
   **/
  public UsagePlanDTO organization(String organization) {
    this.organization = organization;
    return this;
  }

  
  @ApiModelProperty(example = "wso2", value = "Usage policy organization")
  @JsonProperty("organization")
  public String getOrganization() {
    return organization;
  }
  public void setOrganization(String organization) {
    this.organization = organization;
  }


  /**
   **/
  public UsagePlanDTO defaultLimit(UsageLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "")
  @JsonProperty("defaultLimit")
  @NotNull
  public UsageLimitDTO getDefaultLimit() {
    return defaultLimit;
  }
  public void setDefaultLimit(UsageLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
  }


  /**
   * Burst control request count
   **/
  public UsagePlanDTO rateLimitCount(Integer rateLimitCount) {
    this.rateLimitCount = rateLimitCount;
    return this;
  }

  
  @ApiModelProperty(example = "10", value = "Burst control request count")
  @JsonProperty("rateLimitCount")
  public Integer getRateLimitCount() {
    return rateLimitCount;
  }
  public void setRateLimitCount(Integer rateLimitCount) {
    this.rateLimitCount = rateLimitCount;
  }


  /**
   * Burst control time unit
   **/
  public UsagePlanDTO rateLimitTimeUnit(String rateLimitTimeUnit) {
    this.rateLimitTimeUnit = rateLimitTimeUnit;
    return this;
  }

  
  @ApiModelProperty(example = "min", value = "Burst control time unit")
  @JsonProperty("rateLimitTimeUnit")
  public String getRateLimitTimeUnit() {
    return rateLimitTimeUnit;
  }
  public void setRateLimitTimeUnit(String rateLimitTimeUnit) {
    this.rateLimitTimeUnit = rateLimitTimeUnit;
  }


  /**
   * Number of subscriptions allowed
   **/
  public UsagePlanDTO subscriberCount(Integer subscriberCount) {
    this.subscriberCount = subscriberCount;
    return this;
  }

  
  @ApiModelProperty(example = "10", value = "Number of subscriptions allowed")
  @JsonProperty("subscriberCount")
  public Integer getSubscriberCount() {
    return subscriberCount;
  }
  public void setSubscriberCount(Integer subscriberCount) {
    this.subscriberCount = subscriberCount;
  }


  /**
   * Custom attributes added to the Usage plan 
   **/
  public UsagePlanDTO customAttributes(List<CustomAttributeDTO> customAttributes) {
    this.customAttributes = customAttributes;
    return this;
  }

  
  @ApiModelProperty(example = "[]", value = "Custom attributes added to the Usage plan ")
  @JsonProperty("customAttributes")
  public List<CustomAttributeDTO> getCustomAttributes() {
    return customAttributes;
  }
  public void setCustomAttributes(List<CustomAttributeDTO> customAttributes) {
    this.customAttributes = customAttributes;
  }

  public UsagePlanDTO addCustomAttributesItem(CustomAttributeDTO customAttributesItem) {
    if (this.customAttributes == null) {
      this.customAttributes = new ArrayList<>();
    }
    this.customAttributes.add(customAttributesItem);
    return this;
  }


  /**
   * This indicates the action to be taken when a user goes beyond the allocated quota. If checked, the user&#39;s requests will be dropped. If unchecked, the requests will be allowed to pass through. 
   **/
  public UsagePlanDTO stopOnQuotaReach(Boolean stopOnQuotaReach) {
    this.stopOnQuotaReach = stopOnQuotaReach;
    return this;
  }

  
  @ApiModelProperty(value = "This indicates the action to be taken when a user goes beyond the allocated quota. If checked, the user's requests will be dropped. If unchecked, the requests will be allowed to pass through. ")
  @JsonProperty("stopOnQuotaReach")
  public Boolean getStopOnQuotaReach() {
    return stopOnQuotaReach;
  }
  public void setStopOnQuotaReach(Boolean stopOnQuotaReach) {
    this.stopOnQuotaReach = stopOnQuotaReach;
  }


  /**
   * define whether this is Paid or a Free plan. Allowed values are FREE or COMMERCIAL. 
   **/
  public UsagePlanDTO billingPlan(String billingPlan) {
    this.billingPlan = billingPlan;
    return this;
  }

  
  @ApiModelProperty(example = "FREE", value = "define whether this is Paid or a Free plan. Allowed values are FREE or COMMERCIAL. ")
  @JsonProperty("billingPlan")
  public String getBillingPlan() {
    return billingPlan;
  }
  public void setBillingPlan(String billingPlan) {
    this.billingPlan = billingPlan;
  }


  /**
   **/
  public UsagePlanDTO permissions(SubscriptionThrottlePolicyPermissionDTO permissions) {
    this.permissions = permissions;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("permissions")
  public SubscriptionThrottlePolicyPermissionDTO getPermissions() {
    return permissions;
  }
  public void setPermissions(SubscriptionThrottlePolicyPermissionDTO permissions) {
    this.permissions = permissions;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    UsagePlanDTO usagePlan = (UsagePlanDTO) o;
    return Objects.equals(policyId, usagePlan.policyId) &&
        Objects.equals(uuid, usagePlan.uuid) &&
        Objects.equals(policyName, usagePlan.policyName) &&
        Objects.equals(displayName, usagePlan.displayName) &&
        Objects.equals(description, usagePlan.description) &&
        Objects.equals(organization, usagePlan.organization) &&
        Objects.equals(defaultLimit, usagePlan.defaultLimit) &&
        Objects.equals(rateLimitCount, usagePlan.rateLimitCount) &&
        Objects.equals(rateLimitTimeUnit, usagePlan.rateLimitTimeUnit) &&
        Objects.equals(subscriberCount, usagePlan.subscriberCount) &&
        Objects.equals(customAttributes, usagePlan.customAttributes) &&
        Objects.equals(stopOnQuotaReach, usagePlan.stopOnQuotaReach) &&
        Objects.equals(billingPlan, usagePlan.billingPlan) &&
        Objects.equals(permissions, usagePlan.permissions);
  }

  @Override
  public int hashCode() {
    return Objects.hash(policyId, uuid, policyName, displayName, description, organization, defaultLimit, rateLimitCount, rateLimitTimeUnit, subscriberCount, customAttributes, stopOnQuotaReach, billingPlan, permissions);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class UsagePlanDTO {\n");
    
    sb.append("    policyId: ").append(toIndentedString(policyId)).append("\n");
    sb.append("    uuid: ").append(toIndentedString(uuid)).append("\n");
    sb.append("    policyName: ").append(toIndentedString(policyName)).append("\n");
    sb.append("    displayName: ").append(toIndentedString(displayName)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
    sb.append("    organization: ").append(toIndentedString(organization)).append("\n");
    sb.append("    defaultLimit: ").append(toIndentedString(defaultLimit)).append("\n");
    sb.append("    rateLimitCount: ").append(toIndentedString(rateLimitCount)).append("\n");
    sb.append("    rateLimitTimeUnit: ").append(toIndentedString(rateLimitTimeUnit)).append("\n");
    sb.append("    subscriberCount: ").append(toIndentedString(subscriberCount)).append("\n");
    sb.append("    customAttributes: ").append(toIndentedString(customAttributes)).append("\n");
    sb.append("    stopOnQuotaReach: ").append(toIndentedString(stopOnQuotaReach)).append("\n");
    sb.append("    billingPlan: ").append(toIndentedString(billingPlan)).append("\n");
    sb.append("    permissions: ").append(toIndentedString(permissions)).append("\n");
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

