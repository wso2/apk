package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class AdditionalSubscriptionInfoDTO   {
  
  private String subscriptionId;

  private String applicationId;

  private String applicationName;

  private String apiId;

  private Boolean isSolaceAPI;

  private String solaceOrganization;

  private List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO> solaceDeployedEnvironments = null;


  /**
   * The UUID of the subscription
   **/
  public AdditionalSubscriptionInfoDTO subscriptionId(String subscriptionId) {
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
  public AdditionalSubscriptionInfoDTO applicationId(String applicationId) {
    this.applicationId = applicationId;
    return this;
  }

  
  @ApiModelProperty(example = "b3ade481-30b0-4b38-9a67-498a40873a6d", value = "The UUID of the application")
  @JsonProperty("applicationId")
  public String getApplicationId() {
    return applicationId;
  }
  public void setApplicationId(String applicationId) {
    this.applicationId = applicationId;
  }


  /**
   * The name of the application
   **/
  public AdditionalSubscriptionInfoDTO applicationName(String applicationName) {
    this.applicationName = applicationName;
    return this;
  }

  
  @ApiModelProperty(example = "Sample Application", value = "The name of the application")
  @JsonProperty("applicationName")
  public String getApplicationName() {
    return applicationName;
  }
  public void setApplicationName(String applicationName) {
    this.applicationName = applicationName;
  }


  /**
   * The unique identifier of the API.
   **/
  public AdditionalSubscriptionInfoDTO apiId(String apiId) {
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
  public AdditionalSubscriptionInfoDTO isSolaceAPI(Boolean isSolaceAPI) {
    this.isSolaceAPI = isSolaceAPI;
    return this;
  }

  
  @ApiModelProperty(example = "false", value = "")
  @JsonProperty("isSolaceAPI")
  public Boolean getIsSolaceAPI() {
    return isSolaceAPI;
  }
  public void setIsSolaceAPI(Boolean isSolaceAPI) {
    this.isSolaceAPI = isSolaceAPI;
  }


  /**
   **/
  public AdditionalSubscriptionInfoDTO solaceOrganization(String solaceOrganization) {
    this.solaceOrganization = solaceOrganization;
    return this;
  }

  
  @ApiModelProperty(example = "SolaceWso2", value = "")
  @JsonProperty("solaceOrganization")
  public String getSolaceOrganization() {
    return solaceOrganization;
  }
  public void setSolaceOrganization(String solaceOrganization) {
    this.solaceOrganization = solaceOrganization;
  }


  /**
   **/
  public AdditionalSubscriptionInfoDTO solaceDeployedEnvironments(List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO> solaceDeployedEnvironments) {
    this.solaceDeployedEnvironments = solaceDeployedEnvironments;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("solaceDeployedEnvironments")
  public List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO> getSolaceDeployedEnvironments() {
    return solaceDeployedEnvironments;
  }
  public void setSolaceDeployedEnvironments(List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO> solaceDeployedEnvironments) {
    this.solaceDeployedEnvironments = solaceDeployedEnvironments;
  }

  public AdditionalSubscriptionInfoDTO addSolaceDeployedEnvironmentsItem(AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO solaceDeployedEnvironmentsItem) {
    if (this.solaceDeployedEnvironments == null) {
      this.solaceDeployedEnvironments = new ArrayList<>();
    }
    this.solaceDeployedEnvironments.add(solaceDeployedEnvironmentsItem);
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
    AdditionalSubscriptionInfoDTO additionalSubscriptionInfo = (AdditionalSubscriptionInfoDTO) o;
    return Objects.equals(subscriptionId, additionalSubscriptionInfo.subscriptionId) &&
        Objects.equals(applicationId, additionalSubscriptionInfo.applicationId) &&
        Objects.equals(applicationName, additionalSubscriptionInfo.applicationName) &&
        Objects.equals(apiId, additionalSubscriptionInfo.apiId) &&
        Objects.equals(isSolaceAPI, additionalSubscriptionInfo.isSolaceAPI) &&
        Objects.equals(solaceOrganization, additionalSubscriptionInfo.solaceOrganization) &&
        Objects.equals(solaceDeployedEnvironments, additionalSubscriptionInfo.solaceDeployedEnvironments);
  }

  @Override
  public int hashCode() {
    return Objects.hash(subscriptionId, applicationId, applicationName, apiId, isSolaceAPI, solaceOrganization, solaceDeployedEnvironments);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdditionalSubscriptionInfoDTO {\n");
    
    sb.append("    subscriptionId: ").append(toIndentedString(subscriptionId)).append("\n");
    sb.append("    applicationId: ").append(toIndentedString(applicationId)).append("\n");
    sb.append("    applicationName: ").append(toIndentedString(applicationName)).append("\n");
    sb.append("    apiId: ").append(toIndentedString(apiId)).append("\n");
    sb.append("    isSolaceAPI: ").append(toIndentedString(isSolaceAPI)).append("\n");
    sb.append("    solaceOrganization: ").append(toIndentedString(solaceOrganization)).append("\n");
    sb.append("    solaceDeployedEnvironments: ").append(toIndentedString(solaceDeployedEnvironments)).append("\n");
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

