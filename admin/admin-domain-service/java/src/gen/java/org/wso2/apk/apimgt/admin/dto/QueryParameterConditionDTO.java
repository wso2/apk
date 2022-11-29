package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class QueryParameterConditionDTO   {
  
  private String parameterName;

  private String parameterValue;


  /**
   * Name of the query parameter
   **/
  public QueryParameterConditionDTO parameterName(String parameterName) {
    this.parameterName = parameterName;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "Name of the query parameter")
  @JsonProperty("parameterName")
  @NotNull
  public String getParameterName() {
    return parameterName;
  }
  public void setParameterName(String parameterName) {
    this.parameterName = parameterName;
  }


  /**
   * Value of the query parameter to be matched
   **/
  public QueryParameterConditionDTO parameterValue(String parameterValue) {
    this.parameterValue = parameterValue;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "Value of the query parameter to be matched")
  @JsonProperty("parameterValue")
  @NotNull
  public String getParameterValue() {
    return parameterValue;
  }
  public void setParameterValue(String parameterValue) {
    this.parameterValue = parameterValue;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    QueryParameterConditionDTO queryParameterCondition = (QueryParameterConditionDTO) o;
    return Objects.equals(parameterName, queryParameterCondition.parameterName) &&
        Objects.equals(parameterValue, queryParameterCondition.parameterValue);
  }

  @Override
  public int hashCode() {
    return Objects.hash(parameterName, parameterValue);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class QueryParameterConditionDTO {\n");
    
    sb.append("    parameterName: ").append(toIndentedString(parameterName)).append("\n");
    sb.append("    parameterValue: ").append(toIndentedString(parameterValue)).append("\n");
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

