package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APITiersInnerMonetizationAttributesDTO   {
  
  private String fixedPrice;

  private String pricePerRequest;

  private String currencyType;

  private String billingCycle;


  /**
   **/
  public APITiersInnerMonetizationAttributesDTO fixedPrice(String fixedPrice) {
    this.fixedPrice = fixedPrice;
    return this;
  }

  
  @ApiModelProperty(example = "10", value = "")
  @JsonProperty("fixedPrice")
  public String getFixedPrice() {
    return fixedPrice;
  }
  public void setFixedPrice(String fixedPrice) {
    this.fixedPrice = fixedPrice;
  }


  /**
   **/
  public APITiersInnerMonetizationAttributesDTO pricePerRequest(String pricePerRequest) {
    this.pricePerRequest = pricePerRequest;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "")
  @JsonProperty("pricePerRequest")
  public String getPricePerRequest() {
    return pricePerRequest;
  }
  public void setPricePerRequest(String pricePerRequest) {
    this.pricePerRequest = pricePerRequest;
  }


  /**
   **/
  public APITiersInnerMonetizationAttributesDTO currencyType(String currencyType) {
    this.currencyType = currencyType;
    return this;
  }

  
  @ApiModelProperty(example = "USD", value = "")
  @JsonProperty("currencyType")
  public String getCurrencyType() {
    return currencyType;
  }
  public void setCurrencyType(String currencyType) {
    this.currencyType = currencyType;
  }


  /**
   **/
  public APITiersInnerMonetizationAttributesDTO billingCycle(String billingCycle) {
    this.billingCycle = billingCycle;
    return this;
  }

  
  @ApiModelProperty(example = "month", value = "")
  @JsonProperty("billingCycle")
  public String getBillingCycle() {
    return billingCycle;
  }
  public void setBillingCycle(String billingCycle) {
    this.billingCycle = billingCycle;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APITiersInnerMonetizationAttributesDTO apITiersInnerMonetizationAttributes = (APITiersInnerMonetizationAttributesDTO) o;
    return Objects.equals(fixedPrice, apITiersInnerMonetizationAttributes.fixedPrice) &&
        Objects.equals(pricePerRequest, apITiersInnerMonetizationAttributes.pricePerRequest) &&
        Objects.equals(currencyType, apITiersInnerMonetizationAttributes.currencyType) &&
        Objects.equals(billingCycle, apITiersInnerMonetizationAttributes.billingCycle);
  }

  @Override
  public int hashCode() {
    return Objects.hash(fixedPrice, pricePerRequest, currencyType, billingCycle);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APITiersInnerMonetizationAttributesDTO {\n");
    
    sb.append("    fixedPrice: ").append(toIndentedString(fixedPrice)).append("\n");
    sb.append("    pricePerRequest: ").append(toIndentedString(pricePerRequest)).append("\n");
    sb.append("    currencyType: ").append(toIndentedString(currencyType)).append("\n");
    sb.append("    billingCycle: ").append(toIndentedString(billingCycle)).append("\n");
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

