package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO   {
  
  private String environmentName;

  private String environmentDisplayName;

  private String organizationName;

  private List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceURLsInnerDTO> solaceURLs = null;

  private AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO solaceTopicsObject;


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO environmentName(String environmentName) {
    this.environmentName = environmentName;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("environmentName")
  public String getEnvironmentName() {
    return environmentName;
  }
  public void setEnvironmentName(String environmentName) {
    this.environmentName = environmentName;
  }


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO environmentDisplayName(String environmentDisplayName) {
    this.environmentDisplayName = environmentDisplayName;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("environmentDisplayName")
  public String getEnvironmentDisplayName() {
    return environmentDisplayName;
  }
  public void setEnvironmentDisplayName(String environmentDisplayName) {
    this.environmentDisplayName = environmentDisplayName;
  }


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO organizationName(String organizationName) {
    this.organizationName = organizationName;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("organizationName")
  public String getOrganizationName() {
    return organizationName;
  }
  public void setOrganizationName(String organizationName) {
    this.organizationName = organizationName;
  }


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO solaceURLs(List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceURLsInnerDTO> solaceURLs) {
    this.solaceURLs = solaceURLs;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("solaceURLs")
  public List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceURLsInnerDTO> getSolaceURLs() {
    return solaceURLs;
  }
  public void setSolaceURLs(List<AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceURLsInnerDTO> solaceURLs) {
    this.solaceURLs = solaceURLs;
  }

  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO addSolaceURLsItem(AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceURLsInnerDTO solaceURLsItem) {
    if (this.solaceURLs == null) {
      this.solaceURLs = new ArrayList<>();
    }
    this.solaceURLs.add(solaceURLsItem);
    return this;
  }


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO solaceTopicsObject(AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO solaceTopicsObject) {
    this.solaceTopicsObject = solaceTopicsObject;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("SolaceTopicsObject")
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO getSolaceTopicsObject() {
    return solaceTopicsObject;
  }
  public void setSolaceTopicsObject(AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO solaceTopicsObject) {
    this.solaceTopicsObject = solaceTopicsObject;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO additionalSubscriptionInfoSolaceDeployedEnvironmentsInner = (AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO) o;
    return Objects.equals(environmentName, additionalSubscriptionInfoSolaceDeployedEnvironmentsInner.environmentName) &&
        Objects.equals(environmentDisplayName, additionalSubscriptionInfoSolaceDeployedEnvironmentsInner.environmentDisplayName) &&
        Objects.equals(organizationName, additionalSubscriptionInfoSolaceDeployedEnvironmentsInner.organizationName) &&
        Objects.equals(solaceURLs, additionalSubscriptionInfoSolaceDeployedEnvironmentsInner.solaceURLs) &&
        Objects.equals(solaceTopicsObject, additionalSubscriptionInfoSolaceDeployedEnvironmentsInner.solaceTopicsObject);
  }

  @Override
  public int hashCode() {
    return Objects.hash(environmentName, environmentDisplayName, organizationName, solaceURLs, solaceTopicsObject);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerDTO {\n");
    
    sb.append("    environmentName: ").append(toIndentedString(environmentName)).append("\n");
    sb.append("    environmentDisplayName: ").append(toIndentedString(environmentDisplayName)).append("\n");
    sb.append("    organizationName: ").append(toIndentedString(organizationName)).append("\n");
    sb.append("    solaceURLs: ").append(toIndentedString(solaceURLs)).append("\n");
    sb.append("    solaceTopicsObject: ").append(toIndentedString(solaceTopicsObject)).append("\n");
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

