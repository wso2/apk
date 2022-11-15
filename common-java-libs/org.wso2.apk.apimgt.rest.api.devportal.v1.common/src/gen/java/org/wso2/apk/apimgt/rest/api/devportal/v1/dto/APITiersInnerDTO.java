package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APITiersInnerMonetizationAttributesDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APITiersInnerDTO   {
  
  private String tierName;

  private String tierPlan;

  private APITiersInnerMonetizationAttributesDTO monetizationAttributes;


  /**
   **/
  public APITiersInnerDTO tierName(String tierName) {
    this.tierName = tierName;
    return this;
  }

  
  @ApiModelProperty(example = "Gold", value = "")
  @JsonProperty("tierName")
  public String getTierName() {
    return tierName;
  }
  public void setTierName(String tierName) {
    this.tierName = tierName;
  }


  /**
   **/
  public APITiersInnerDTO tierPlan(String tierPlan) {
    this.tierPlan = tierPlan;
    return this;
  }

  
  @ApiModelProperty(example = "COMMERCIAL", value = "")
  @JsonProperty("tierPlan")
  public String getTierPlan() {
    return tierPlan;
  }
  public void setTierPlan(String tierPlan) {
    this.tierPlan = tierPlan;
  }


  /**
   **/
  public APITiersInnerDTO monetizationAttributes(APITiersInnerMonetizationAttributesDTO monetizationAttributes) {
    this.monetizationAttributes = monetizationAttributes;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("monetizationAttributes")
  public APITiersInnerMonetizationAttributesDTO getMonetizationAttributes() {
    return monetizationAttributes;
  }
  public void setMonetizationAttributes(APITiersInnerMonetizationAttributesDTO monetizationAttributes) {
    this.monetizationAttributes = monetizationAttributes;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APITiersInnerDTO apITiersInner = (APITiersInnerDTO) o;
    return Objects.equals(tierName, apITiersInner.tierName) &&
        Objects.equals(tierPlan, apITiersInner.tierPlan) &&
        Objects.equals(monetizationAttributes, apITiersInner.monetizationAttributes);
  }

  @Override
  public int hashCode() {
    return Objects.hash(tierName, tierPlan, monetizationAttributes);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APITiersInnerDTO {\n");
    
    sb.append("    tierName: ").append(toIndentedString(tierName)).append("\n");
    sb.append("    tierPlan: ").append(toIndentedString(tierPlan)).append("\n");
    sb.append("    monetizationAttributes: ").append(toIndentedString(monetizationAttributes)).append("\n");
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

