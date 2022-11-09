package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class ExportThrottlePolicyDTO   {
  
  private String type;

  private String subtype;

  private String version;

  private Object data;


  /**
   **/
  public ExportThrottlePolicyDTO type(String type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("type")
  public String getType() {
    return type;
  }
  public void setType(String type) {
    this.type = type;
  }


  /**
   **/
  public ExportThrottlePolicyDTO subtype(String subtype) {
    this.subtype = subtype;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("subtype")
  public String getSubtype() {
    return subtype;
  }
  public void setSubtype(String subtype) {
    this.subtype = subtype;
  }


  /**
   **/
  public ExportThrottlePolicyDTO version(String version) {
    this.version = version;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("version")
  public String getVersion() {
    return version;
  }
  public void setVersion(String version) {
    this.version = version;
  }


  /**
   **/
  public ExportThrottlePolicyDTO data(Object data) {
    this.data = data;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("data")
  public Object getData() {
    return data;
  }
  public void setData(Object data) {
    this.data = data;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ExportThrottlePolicyDTO exportThrottlePolicy = (ExportThrottlePolicyDTO) o;
    return Objects.equals(type, exportThrottlePolicy.type) &&
        Objects.equals(subtype, exportThrottlePolicy.subtype) &&
        Objects.equals(version, exportThrottlePolicy.version) &&
        Objects.equals(data, exportThrottlePolicy.data);
  }

  @Override
  public int hashCode() {
    return Objects.hash(type, subtype, version, data);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ExportThrottlePolicyDTO {\n");
    
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    subtype: ").append(toIndentedString(subtype)).append("\n");
    sb.append("    version: ").append(toIndentedString(version)).append("\n");
    sb.append("    data: ").append(toIndentedString(data)).append("\n");
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

