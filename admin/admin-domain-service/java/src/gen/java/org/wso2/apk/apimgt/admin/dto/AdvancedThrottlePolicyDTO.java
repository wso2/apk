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
import org.wso2.apk.apimgt.admin.dto.ThrottlePolicyDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class AdvancedThrottlePolicyDTO extends ThrottlePolicyDTO  {
  
  private ThrottleLimitDTO defaultLimit;

  private List<ConditionalGroupDTO> conditionalGroups = null;


  /**
   **/
  public AdvancedThrottlePolicyDTO defaultLimit(ThrottleLimitDTO defaultLimit) {
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
  public AdvancedThrottlePolicyDTO conditionalGroups(List<ConditionalGroupDTO> conditionalGroups) {
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

  public AdvancedThrottlePolicyDTO addConditionalGroupsItem(ConditionalGroupDTO conditionalGroupsItem) {
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
    AdvancedThrottlePolicyDTO advancedThrottlePolicy = (AdvancedThrottlePolicyDTO) o;
    return Objects.equals(defaultLimit, advancedThrottlePolicy.defaultLimit) &&
        Objects.equals(conditionalGroups, advancedThrottlePolicy.conditionalGroups) &&
        super.equals(o);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultLimit, conditionalGroups, super.hashCode());
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdvancedThrottlePolicyDTO {\n");
    sb.append("    ").append(toIndentedString(super.toString())).append("\n");
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

