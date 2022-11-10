package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APIInfoAdditionalPropertiesInnerDTO   {
  
  private String name;

  private String value;

  private Boolean display;


  /**
   **/
  public APIInfoAdditionalPropertiesInnerDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("name")
  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public APIInfoAdditionalPropertiesInnerDTO value(String value) {
    this.value = value;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("value")
  public String getValue() {
    return value;
  }
  public void setValue(String value) {
    this.value = value;
  }


  /**
   **/
  public APIInfoAdditionalPropertiesInnerDTO display(Boolean display) {
    this.display = display;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("display")
  public Boolean getDisplay() {
    return display;
  }
  public void setDisplay(Boolean display) {
    this.display = display;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIInfoAdditionalPropertiesInnerDTO apIInfoAdditionalPropertiesInner = (APIInfoAdditionalPropertiesInnerDTO) o;
    return Objects.equals(name, apIInfoAdditionalPropertiesInner.name) &&
        Objects.equals(value, apIInfoAdditionalPropertiesInner.value) &&
        Objects.equals(display, apIInfoAdditionalPropertiesInner.display);
  }

  @Override
  public int hashCode() {
    return Objects.hash(name, value, display);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIInfoAdditionalPropertiesInnerDTO {\n");
    
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    value: ").append(toIndentedString(value)).append("\n");
    sb.append("    display: ").append(toIndentedString(display)).append("\n");
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

