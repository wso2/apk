package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APIEndpointURLsInnerDefaultVersionURLsDTO;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APIEndpointURLsInnerURLsDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APIEndpointURLsInnerDTO   {
  
  private String environmentName;

  private String environmentDisplayName;

  private String environmentType;

  private APIEndpointURLsInnerURLsDTO urLs;

  private APIEndpointURLsInnerDefaultVersionURLsDTO defaultVersionURLs;


  /**
   **/
  public APIEndpointURLsInnerDTO environmentName(String environmentName) {
    this.environmentName = environmentName;
    return this;
  }

  
  @ApiModelProperty(example = "Default", value = "")
  @JsonProperty("environmentName")
  public String getEnvironmentName() {
    return environmentName;
  }
  public void setEnvironmentName(String environmentName) {
    this.environmentName = environmentName;
  }


  /**
   **/
  public APIEndpointURLsInnerDTO environmentDisplayName(String environmentDisplayName) {
    this.environmentDisplayName = environmentDisplayName;
    return this;
  }

  
  @ApiModelProperty(example = "Default", value = "")
  @JsonProperty("environmentDisplayName")
  public String getEnvironmentDisplayName() {
    return environmentDisplayName;
  }
  public void setEnvironmentDisplayName(String environmentDisplayName) {
    this.environmentDisplayName = environmentDisplayName;
  }


  /**
   **/
  public APIEndpointURLsInnerDTO environmentType(String environmentType) {
    this.environmentType = environmentType;
    return this;
  }

  
  @ApiModelProperty(example = "hybrid", value = "")
  @JsonProperty("environmentType")
  public String getEnvironmentType() {
    return environmentType;
  }
  public void setEnvironmentType(String environmentType) {
    this.environmentType = environmentType;
  }


  /**
   **/
  public APIEndpointURLsInnerDTO urLs(APIEndpointURLsInnerURLsDTO urLs) {
    this.urLs = urLs;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("URLs")
  public APIEndpointURLsInnerURLsDTO getUrLs() {
    return urLs;
  }
  public void setUrLs(APIEndpointURLsInnerURLsDTO urLs) {
    this.urLs = urLs;
  }


  /**
   **/
  public APIEndpointURLsInnerDTO defaultVersionURLs(APIEndpointURLsInnerDefaultVersionURLsDTO defaultVersionURLs) {
    this.defaultVersionURLs = defaultVersionURLs;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("defaultVersionURLs")
  public APIEndpointURLsInnerDefaultVersionURLsDTO getDefaultVersionURLs() {
    return defaultVersionURLs;
  }
  public void setDefaultVersionURLs(APIEndpointURLsInnerDefaultVersionURLsDTO defaultVersionURLs) {
    this.defaultVersionURLs = defaultVersionURLs;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIEndpointURLsInnerDTO apIEndpointURLsInner = (APIEndpointURLsInnerDTO) o;
    return Objects.equals(environmentName, apIEndpointURLsInner.environmentName) &&
        Objects.equals(environmentDisplayName, apIEndpointURLsInner.environmentDisplayName) &&
        Objects.equals(environmentType, apIEndpointURLsInner.environmentType) &&
        Objects.equals(urLs, apIEndpointURLsInner.urLs) &&
        Objects.equals(defaultVersionURLs, apIEndpointURLsInner.defaultVersionURLs);
  }

  @Override
  public int hashCode() {
    return Objects.hash(environmentName, environmentDisplayName, environmentType, urLs, defaultVersionURLs);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIEndpointURLsInnerDTO {\n");
    
    sb.append("    environmentName: ").append(toIndentedString(environmentName)).append("\n");
    sb.append("    environmentDisplayName: ").append(toIndentedString(environmentDisplayName)).append("\n");
    sb.append("    environmentType: ").append(toIndentedString(environmentType)).append("\n");
    sb.append("    urLs: ").append(toIndentedString(urLs)).append("\n");
    sb.append("    defaultVersionURLs: ").append(toIndentedString(defaultVersionURLs)).append("\n");
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

