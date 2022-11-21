package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class APIAdditionalPropertiesValueDTO   {
  
  private String name;

  private String value;

  private Boolean display = false;


  /**
   **/
  public APIAdditionalPropertiesValueDTO name(String name) {
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
  public APIAdditionalPropertiesValueDTO value(String value) {
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
  public APIAdditionalPropertiesValueDTO display(Boolean display) {
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
    APIAdditionalPropertiesValueDTO apIAdditionalPropertiesValue = (APIAdditionalPropertiesValueDTO) o;
    return Objects.equals(name, apIAdditionalPropertiesValue.name) &&
        Objects.equals(value, apIAdditionalPropertiesValue.value) &&
        Objects.equals(display, apIAdditionalPropertiesValue.display);
  }

  @Override
  public int hashCode() {
    return Objects.hash(name, value, display);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIAdditionalPropertiesValueDTO {\n");
    
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

