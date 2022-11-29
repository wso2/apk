package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class ApplicationInfoDTO   {
  
  private String applicationId;

  private String name;

  private String throttlingPolicy;

  private String description;

  private String status = "";

  private List<String> groups = null;

  private Integer subscriptionCount;

  private Object attributes;

  private String owner;

  private String createdTime;

  private String updatedTime;


  /**
   **/
  public ApplicationInfoDTO applicationId(String applicationId) {
    this.applicationId = applicationId;
    return this;
  }

  
  @ApiModelProperty(example = "01234567-0123-0123-0123-012345678901", value = "")
  @JsonProperty("applicationId")
  public String getApplicationId() {
    return applicationId;
  }
  public void setApplicationId(String applicationId) {
    this.applicationId = applicationId;
  }


  /**
   **/
  public ApplicationInfoDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "CalculatorApp", value = "")
  @JsonProperty("name")
  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public ApplicationInfoDTO throttlingPolicy(String throttlingPolicy) {
    this.throttlingPolicy = throttlingPolicy;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", value = "")
  @JsonProperty("throttlingPolicy")
  public String getThrottlingPolicy() {
    return throttlingPolicy;
  }
  public void setThrottlingPolicy(String throttlingPolicy) {
    this.throttlingPolicy = throttlingPolicy;
  }


  /**
   **/
  public ApplicationInfoDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "Sample calculator application", value = "")
  @JsonProperty("description")
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   **/
  public ApplicationInfoDTO status(String status) {
    this.status = status;
    return this;
  }

  
  @ApiModelProperty(example = "APPROVED", value = "")
  @JsonProperty("status")
  public String getStatus() {
    return status;
  }
  public void setStatus(String status) {
    this.status = status;
  }


  /**
   **/
  public ApplicationInfoDTO groups(List<String> groups) {
    this.groups = groups;
    return this;
  }

  
  @ApiModelProperty(example = "", value = "")
  @JsonProperty("groups")
  public List<String> getGroups() {
    return groups;
  }
  public void setGroups(List<String> groups) {
    this.groups = groups;
  }

  public ApplicationInfoDTO addGroupsItem(String groupsItem) {
    if (this.groups == null) {
      this.groups = new ArrayList<>();
    }
    this.groups.add(groupsItem);
    return this;
  }


  /**
   **/
  public ApplicationInfoDTO subscriptionCount(Integer subscriptionCount) {
    this.subscriptionCount = subscriptionCount;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("subscriptionCount")
  public Integer getSubscriptionCount() {
    return subscriptionCount;
  }
  public void setSubscriptionCount(Integer subscriptionCount) {
    this.subscriptionCount = subscriptionCount;
  }


  /**
   **/
  public ApplicationInfoDTO attributes(Object attributes) {
    this.attributes = attributes;
    return this;
  }

  
  @ApiModelProperty(example = "External Reference ID, Billing Tier", value = "")
  @JsonProperty("attributes")
  public Object getAttributes() {
    return attributes;
  }
  public void setAttributes(Object attributes) {
    this.attributes = attributes;
  }


  /**
   **/
  public ApplicationInfoDTO owner(String owner) {
    this.owner = owner;
    return this;
  }

  
  @ApiModelProperty(example = "admin", value = "")
  @JsonProperty("owner")
  public String getOwner() {
    return owner;
  }
  public void setOwner(String owner) {
    this.owner = owner;
  }


  /**
   **/
  public ApplicationInfoDTO createdTime(String createdTime) {
    this.createdTime = createdTime;
    return this;
  }

  
  @ApiModelProperty(example = "1651555310208", value = "")
  @JsonProperty("createdTime")
  public String getCreatedTime() {
    return createdTime;
  }
  public void setCreatedTime(String createdTime) {
    this.createdTime = createdTime;
  }


  /**
   **/
  public ApplicationInfoDTO updatedTime(String updatedTime) {
    this.updatedTime = updatedTime;
    return this;
  }

  
  @ApiModelProperty(example = "1651555310208", value = "")
  @JsonProperty("updatedTime")
  public String getUpdatedTime() {
    return updatedTime;
  }
  public void setUpdatedTime(String updatedTime) {
    this.updatedTime = updatedTime;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ApplicationInfoDTO applicationInfo = (ApplicationInfoDTO) o;
    return Objects.equals(applicationId, applicationInfo.applicationId) &&
        Objects.equals(name, applicationInfo.name) &&
        Objects.equals(throttlingPolicy, applicationInfo.throttlingPolicy) &&
        Objects.equals(description, applicationInfo.description) &&
        Objects.equals(status, applicationInfo.status) &&
        Objects.equals(groups, applicationInfo.groups) &&
        Objects.equals(subscriptionCount, applicationInfo.subscriptionCount) &&
        Objects.equals(attributes, applicationInfo.attributes) &&
        Objects.equals(owner, applicationInfo.owner) &&
        Objects.equals(createdTime, applicationInfo.createdTime) &&
        Objects.equals(updatedTime, applicationInfo.updatedTime);
  }

  @Override
  public int hashCode() {
    return Objects.hash(applicationId, name, throttlingPolicy, description, status, groups, subscriptionCount, attributes, owner, createdTime, updatedTime);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ApplicationInfoDTO {\n");
    
    sb.append("    applicationId: ").append(toIndentedString(applicationId)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    throttlingPolicy: ").append(toIndentedString(throttlingPolicy)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
    sb.append("    status: ").append(toIndentedString(status)).append("\n");
    sb.append("    groups: ").append(toIndentedString(groups)).append("\n");
    sb.append("    subscriptionCount: ").append(toIndentedString(subscriptionCount)).append("\n");
    sb.append("    attributes: ").append(toIndentedString(attributes)).append("\n");
    sb.append("    owner: ").append(toIndentedString(owner)).append("\n");
    sb.append("    createdTime: ").append(toIndentedString(createdTime)).append("\n");
    sb.append("    updatedTime: ").append(toIndentedString(updatedTime)).append("\n");
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

