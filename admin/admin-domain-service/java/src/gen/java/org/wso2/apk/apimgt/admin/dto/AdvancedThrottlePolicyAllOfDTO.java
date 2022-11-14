package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.admin.dto.ConditionalGroupDTO;
import org.wso2.apk.apimgt.admin.dto.ThrottleLimitDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class AdvancedThrottlePolicyAllOfDTO   {
  
  private ThrottleLimitDTO defaultLimit;

  private List<ConditionalGroupDTO> conditionalGroups = null;


  /**
   **/
  public AdvancedThrottlePolicyAllOfDTO defaultLimit(ThrottleLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "")
  @JsonProperty("defaultLimit")
  @NotNull
  public ThrottleLimitDTO getDefaultLimit() {
    return defaultLimit;
  }
  public void setDefaultLimit(ThrottleLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
  }


  /**
   * Group of conditions which allow adding different parameter conditions to the throttling limit. 
   **/
  public AdvancedThrottlePolicyAllOfDTO conditionalGroups(List<ConditionalGroupDTO> conditionalGroups) {
    this.conditionalGroups = conditionalGroups;
    return this;
  }

  
  @ApiModelProperty(value = "Group of conditions which allow adding different parameter conditions to the throttling limit. ")
  @JsonProperty("conditionalGroups")
  public List<ConditionalGroupDTO> getConditionalGroups() {
    return conditionalGroups;
  }
  public void setConditionalGroups(List<ConditionalGroupDTO> conditionalGroups) {
    this.conditionalGroups = conditionalGroups;
  }

  public AdvancedThrottlePolicyAllOfDTO addConditionalGroupsItem(ConditionalGroupDTO conditionalGroupsItem) {
    if (this.conditionalGroups == null) {
      this.conditionalGroups = new ArrayList<>();
    }
    this.conditionalGroups.add(conditionalGroupsItem);
    return this;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    AdvancedThrottlePolicyAllOfDTO advancedThrottlePolicyAllOf = (AdvancedThrottlePolicyAllOfDTO) o;
    return Objects.equals(defaultLimit, advancedThrottlePolicyAllOf.defaultLimit) &&
        Objects.equals(conditionalGroups, advancedThrottlePolicyAllOf.conditionalGroups);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultLimit, conditionalGroups);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdvancedThrottlePolicyAllOfDTO {\n");
    
    sb.append("    defaultLimit: ").append(toIndentedString(defaultLimit)).append("\n");
    sb.append("    conditionalGroups: ").append(toIndentedString(conditionalGroups)).append("\n");
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

