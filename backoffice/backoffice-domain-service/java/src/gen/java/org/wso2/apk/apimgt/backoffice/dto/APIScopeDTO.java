package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.*;


import java.util.Objects;



public class APIScopeDTO   {
  
  private ScopeDTO scope;


  /**
   **/
  public APIScopeDTO scope(ScopeDTO scope) {
    this.scope = scope;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "")
  @JsonProperty("scope")
  @NotNull
  public ScopeDTO getScope() {
    return scope;
  }
  public void setScope(ScopeDTO scope) {
    this.scope = scope;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIScopeDTO apIScope = (APIScopeDTO) o;
    return Objects.equals(scope, apIScope.scope);
  }

  @Override
  public int hashCode() {
    return Objects.hash(scope);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIScopeDTO {\n");
    
    sb.append("    scope: ").append(toIndentedString(scope)).append("\n");
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

