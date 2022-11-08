package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.KeyManagerDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class KeyManagerWellKnownResponseDTO   {
  
  private Boolean valid = false;

  private KeyManagerDTO value;


  /**
   **/
  public KeyManagerWellKnownResponseDTO valid(Boolean valid) {
    this.valid = valid;
    return this;
  }

  
  @ApiModelProperty(example = "true", value = "")
  @JsonProperty("valid")
  public Boolean getValid() {
    return valid;
  }
  public void setValid(Boolean valid) {
    this.valid = valid;
  }


  /**
   **/
  public KeyManagerWellKnownResponseDTO value(KeyManagerDTO value) {
    this.value = value;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("value")
  public KeyManagerDTO getValue() {
    return value;
  }
  public void setValue(KeyManagerDTO value) {
    this.value = value;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    KeyManagerWellKnownResponseDTO keyManagerWellKnownResponse = (KeyManagerWellKnownResponseDTO) o;
    return Objects.equals(valid, keyManagerWellKnownResponse.valid) &&
        Objects.equals(value, keyManagerWellKnownResponse.value);
  }

  @Override
  public int hashCode() {
    return Objects.hash(valid, value);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class KeyManagerWellKnownResponseDTO {\n");
    
    sb.append("    valid: ").append(toIndentedString(valid)).append("\n");
    sb.append("    value: ").append(toIndentedString(value)).append("\n");
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

