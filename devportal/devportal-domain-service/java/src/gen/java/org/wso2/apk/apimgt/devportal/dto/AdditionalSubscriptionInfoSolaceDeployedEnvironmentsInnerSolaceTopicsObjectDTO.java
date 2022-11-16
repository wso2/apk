package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO   {
  
  private SolaceTopicsDTO defaultSyntax;

  private SolaceTopicsDTO mqttSyntax;


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO defaultSyntax(SolaceTopicsDTO defaultSyntax) {
    this.defaultSyntax = defaultSyntax;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("defaultSyntax")
  public SolaceTopicsDTO getDefaultSyntax() {
    return defaultSyntax;
  }
  public void setDefaultSyntax(SolaceTopicsDTO defaultSyntax) {
    this.defaultSyntax = defaultSyntax;
  }


  /**
   **/
  public AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO mqttSyntax(SolaceTopicsDTO mqttSyntax) {
    this.mqttSyntax = mqttSyntax;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("mqttSyntax")
  public SolaceTopicsDTO getMqttSyntax() {
    return mqttSyntax;
  }
  public void setMqttSyntax(SolaceTopicsDTO mqttSyntax) {
    this.mqttSyntax = mqttSyntax;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO additionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObject = (AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO) o;
    return Objects.equals(defaultSyntax, additionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObject.defaultSyntax) &&
        Objects.equals(mqttSyntax, additionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObject.mqttSyntax);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultSyntax, mqttSyntax);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdditionalSubscriptionInfoSolaceDeployedEnvironmentsInnerSolaceTopicsObjectDTO {\n");
    
    sb.append("    defaultSyntax: ").append(toIndentedString(defaultSyntax)).append("\n");
    sb.append("    mqttSyntax: ").append(toIndentedString(mqttSyntax)).append("\n");
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

