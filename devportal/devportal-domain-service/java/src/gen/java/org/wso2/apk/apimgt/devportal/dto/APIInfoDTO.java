package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class APIInfoDTO   {
  
  private String id;

  private String name;

  private String description;

  private String context;

  private String version;

  private String type;

  private String createdTime;

  private String provider;

  private String lifeCycleStatus;

  private String thumbnailUri;

  private String avgRating;

  private List<String> throttlingPolicies = null;

  private AdvertiseInfoDTO advertiseInfo;

  private APIBusinessInformationDTO businessInformation;

  private Boolean isSubscriptionAvailable;

  private String monetizationLabel;

  private String gatewayVendor;

  private List<APIInfoAdditionalPropertiesInnerDTO> additionalProperties = null;


  /**
   **/
  public APIInfoDTO id(String id) {
    this.id = id;
    return this;
  }

  
  @ApiModelProperty(example = "01234567-0123-0123-0123-012345678901", value = "")
  @JsonProperty("id")
  public String getId() {
    return id;
  }
  public void setId(String id) {
    this.id = id;
  }


  /**
   **/
  public APIInfoDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "CalculatorAPI", value = "")
  @JsonProperty("name")
  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public APIInfoDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "A calculator API that supports basic operations", value = "")
  @JsonProperty("description")
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   **/
  public APIInfoDTO context(String context) {
    this.context = context;
    return this;
  }

  
  @ApiModelProperty(example = "CalculatorAPI", value = "")
  @JsonProperty("context")
  public String getContext() {
    return context;
  }
  public void setContext(String context) {
    this.context = context;
  }


  /**
   **/
  public APIInfoDTO version(String version) {
    this.version = version;
    return this;
  }

  
  @ApiModelProperty(example = "1.0.0", value = "")
  @JsonProperty("version")
  public String getVersion() {
    return version;
  }
  public void setVersion(String version) {
    this.version = version;
  }


  /**
   **/
  public APIInfoDTO type(String type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(example = "WS", value = "")
  @JsonProperty("type")
  public String getType() {
    return type;
  }
  public void setType(String type) {
    this.type = type;
  }


  /**
   **/
  public APIInfoDTO createdTime(String createdTime) {
    this.createdTime = createdTime;
    return this;
  }

  
  @ApiModelProperty(example = "1614020559444", value = "")
  @JsonProperty("createdTime")
  public String getCreatedTime() {
    return createdTime;
  }
  public void setCreatedTime(String createdTime) {
    this.createdTime = createdTime;
  }


  /**
   * If the provider value is not given, the user invoking the API will be used as the provider. 
   **/
  public APIInfoDTO provider(String provider) {
    this.provider = provider;
    return this;
  }

  
  @ApiModelProperty(example = "admin", value = "If the provider value is not given, the user invoking the API will be used as the provider. ")
  @JsonProperty("provider")
  public String getProvider() {
    return provider;
  }
  public void setProvider(String provider) {
    this.provider = provider;
  }


  /**
   **/
  public APIInfoDTO lifeCycleStatus(String lifeCycleStatus) {
    this.lifeCycleStatus = lifeCycleStatus;
    return this;
  }

  
  @ApiModelProperty(example = "PUBLISHED", value = "")
  @JsonProperty("lifeCycleStatus")
  public String getLifeCycleStatus() {
    return lifeCycleStatus;
  }
  public void setLifeCycleStatus(String lifeCycleStatus) {
    this.lifeCycleStatus = lifeCycleStatus;
  }


  /**
   **/
  public APIInfoDTO thumbnailUri(String thumbnailUri) {
    this.thumbnailUri = thumbnailUri;
    return this;
  }

  
  @ApiModelProperty(example = "/apis/01234567-0123-0123-0123-012345678901/thumbnail", value = "")
  @JsonProperty("thumbnailUri")
  public String getThumbnailUri() {
    return thumbnailUri;
  }
  public void setThumbnailUri(String thumbnailUri) {
    this.thumbnailUri = thumbnailUri;
  }


  /**
   * Average rating of the API
   **/
  public APIInfoDTO avgRating(String avgRating) {
    this.avgRating = avgRating;
    return this;
  }

  
  @ApiModelProperty(example = "4.5", value = "Average rating of the API")
  @JsonProperty("avgRating")
  public String getAvgRating() {
    return avgRating;
  }
  public void setAvgRating(String avgRating) {
    this.avgRating = avgRating;
  }


  /**
   * List of throttling policies of the API
   **/
  public APIInfoDTO throttlingPolicies(List<String> throttlingPolicies) {
    this.throttlingPolicies = throttlingPolicies;
    return this;
  }

  
  @ApiModelProperty(example = "[\"Unlimited\",\"Bronze\"]", value = "List of throttling policies of the API")
  @JsonProperty("throttlingPolicies")
  public List<String> getThrottlingPolicies() {
    return throttlingPolicies;
  }
  public void setThrottlingPolicies(List<String> throttlingPolicies) {
    this.throttlingPolicies = throttlingPolicies;
  }

  public APIInfoDTO addThrottlingPoliciesItem(String throttlingPoliciesItem) {
    if (this.throttlingPolicies == null) {
      this.throttlingPolicies = new ArrayList<>();
    }
    this.throttlingPolicies.add(throttlingPoliciesItem);
    return this;
  }


  /**
   **/
  public APIInfoDTO advertiseInfo(AdvertiseInfoDTO advertiseInfo) {
    this.advertiseInfo = advertiseInfo;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("advertiseInfo")
  public AdvertiseInfoDTO getAdvertiseInfo() {
    return advertiseInfo;
  }
  public void setAdvertiseInfo(AdvertiseInfoDTO advertiseInfo) {
    this.advertiseInfo = advertiseInfo;
  }


  /**
   **/
  public APIInfoDTO businessInformation(APIBusinessInformationDTO businessInformation) {
    this.businessInformation = businessInformation;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("businessInformation")
  public APIBusinessInformationDTO getBusinessInformation() {
    return businessInformation;
  }
  public void setBusinessInformation(APIBusinessInformationDTO businessInformation) {
    this.businessInformation = businessInformation;
  }


  /**
   **/
  public APIInfoDTO isSubscriptionAvailable(Boolean isSubscriptionAvailable) {
    this.isSubscriptionAvailable = isSubscriptionAvailable;
    return this;
  }

  
  @ApiModelProperty(example = "false", value = "")
  @JsonProperty("isSubscriptionAvailable")
  public Boolean getIsSubscriptionAvailable() {
    return isSubscriptionAvailable;
  }
  public void setIsSubscriptionAvailable(Boolean isSubscriptionAvailable) {
    this.isSubscriptionAvailable = isSubscriptionAvailable;
  }


  /**
   **/
  public APIInfoDTO monetizationLabel(String monetizationLabel) {
    this.monetizationLabel = monetizationLabel;
    return this;
  }

  
  @ApiModelProperty(example = "Free", value = "")
  @JsonProperty("monetizationLabel")
  public String getMonetizationLabel() {
    return monetizationLabel;
  }
  public void setMonetizationLabel(String monetizationLabel) {
    this.monetizationLabel = monetizationLabel;
  }


  /**
   **/
  public APIInfoDTO gatewayVendor(String gatewayVendor) {
    this.gatewayVendor = gatewayVendor;
    return this;
  }

  
  @ApiModelProperty(example = "WSO2", value = "")
  @JsonProperty("gatewayVendor")
  public String getGatewayVendor() {
    return gatewayVendor;
  }
  public void setGatewayVendor(String gatewayVendor) {
    this.gatewayVendor = gatewayVendor;
  }


  /**
   * Custom(user defined) properties of API 
   **/
  public APIInfoDTO additionalProperties(List<APIInfoAdditionalPropertiesInnerDTO> additionalProperties) {
    this.additionalProperties = additionalProperties;
    return this;
  }

  
  @ApiModelProperty(example = "{}", value = "Custom(user defined) properties of API ")
  @JsonProperty("additionalProperties")
  public List<APIInfoAdditionalPropertiesInnerDTO> getAdditionalProperties() {
    return additionalProperties;
  }
  public void setAdditionalProperties(List<APIInfoAdditionalPropertiesInnerDTO> additionalProperties) {
    this.additionalProperties = additionalProperties;
  }

  public APIInfoDTO addAdditionalPropertiesItem(APIInfoAdditionalPropertiesInnerDTO additionalPropertiesItem) {
    if (this.additionalProperties == null) {
      this.additionalProperties = new ArrayList<>();
    }
    this.additionalProperties.add(additionalPropertiesItem);
    return this;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIInfoDTO apIInfo = (APIInfoDTO) o;
    return Objects.equals(id, apIInfo.id) &&
        Objects.equals(name, apIInfo.name) &&
        Objects.equals(description, apIInfo.description) &&
        Objects.equals(context, apIInfo.context) &&
        Objects.equals(version, apIInfo.version) &&
        Objects.equals(type, apIInfo.type) &&
        Objects.equals(createdTime, apIInfo.createdTime) &&
        Objects.equals(provider, apIInfo.provider) &&
        Objects.equals(lifeCycleStatus, apIInfo.lifeCycleStatus) &&
        Objects.equals(thumbnailUri, apIInfo.thumbnailUri) &&
        Objects.equals(avgRating, apIInfo.avgRating) &&
        Objects.equals(throttlingPolicies, apIInfo.throttlingPolicies) &&
        Objects.equals(advertiseInfo, apIInfo.advertiseInfo) &&
        Objects.equals(businessInformation, apIInfo.businessInformation) &&
        Objects.equals(isSubscriptionAvailable, apIInfo.isSubscriptionAvailable) &&
        Objects.equals(monetizationLabel, apIInfo.monetizationLabel) &&
        Objects.equals(gatewayVendor, apIInfo.gatewayVendor) &&
        Objects.equals(additionalProperties, apIInfo.additionalProperties);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, name, description, context, version, type, createdTime, provider, lifeCycleStatus, thumbnailUri, avgRating, throttlingPolicies, advertiseInfo, businessInformation, isSubscriptionAvailable, monetizationLabel, gatewayVendor, additionalProperties);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIInfoDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
    sb.append("    context: ").append(toIndentedString(context)).append("\n");
    sb.append("    version: ").append(toIndentedString(version)).append("\n");
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    createdTime: ").append(toIndentedString(createdTime)).append("\n");
    sb.append("    provider: ").append(toIndentedString(provider)).append("\n");
    sb.append("    lifeCycleStatus: ").append(toIndentedString(lifeCycleStatus)).append("\n");
    sb.append("    thumbnailUri: ").append(toIndentedString(thumbnailUri)).append("\n");
    sb.append("    avgRating: ").append(toIndentedString(avgRating)).append("\n");
    sb.append("    throttlingPolicies: ").append(toIndentedString(throttlingPolicies)).append("\n");
    sb.append("    advertiseInfo: ").append(toIndentedString(advertiseInfo)).append("\n");
    sb.append("    businessInformation: ").append(toIndentedString(businessInformation)).append("\n");
    sb.append("    isSubscriptionAvailable: ").append(toIndentedString(isSubscriptionAvailable)).append("\n");
    sb.append("    monetizationLabel: ").append(toIndentedString(monetizationLabel)).append("\n");
    sb.append("    gatewayVendor: ").append(toIndentedString(gatewayVendor)).append("\n");
    sb.append("    additionalProperties: ").append(toIndentedString(additionalProperties)).append("\n");
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

