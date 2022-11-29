package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class APIOperationsDTO   {
  
  private String id;

  private String target;

  private String verb;

  private String usagePlan;


  /**
   **/
  public APIOperationsDTO id(String id) {
    this.id = id;
    return this;
  }

  
  @ApiModelProperty(example = "postapiresource", value = "")
  @JsonProperty("id")
  public String getId() {
    return id;
  }
  public void setId(String id) {
    this.id = id;
  }


  /**
   **/
  public APIOperationsDTO target(String target) {
    this.target = target;
    return this;
  }

  
  @ApiModelProperty(example = "/order/{orderId}", value = "")
  @JsonProperty("target")
  public String getTarget() {
    return target;
  }
  public void setTarget(String target) {
    this.target = target;
  }


  /**
   **/
  public APIOperationsDTO verb(String verb) {
    this.verb = verb;
    return this;
  }

  
  @ApiModelProperty(example = "POST", value = "")
  @JsonProperty("verb")
  public String getVerb() {
    return verb;
  }
  public void setVerb(String verb) {
    this.verb = verb;
  }


  /**
   **/
  public APIOperationsDTO usagePlan(String usagePlan) {
    this.usagePlan = usagePlan;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", value = "")
  @JsonProperty("usagePlan")
  public String getUsagePlan() {
    return usagePlan;
  }
  public void setUsagePlan(String usagePlan) {
    this.usagePlan = usagePlan;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIOperationsDTO apIOperations = (APIOperationsDTO) o;
    return Objects.equals(id, apIOperations.id) &&
        Objects.equals(target, apIOperations.target) &&
        Objects.equals(verb, apIOperations.verb) &&
        Objects.equals(usagePlan, apIOperations.usagePlan);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, target, verb, usagePlan);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIOperationsDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    target: ").append(toIndentedString(target)).append("\n");
    sb.append("    verb: ").append(toIndentedString(verb)).append("\n");
    sb.append("    usagePlan: ").append(toIndentedString(usagePlan)).append("\n");
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

