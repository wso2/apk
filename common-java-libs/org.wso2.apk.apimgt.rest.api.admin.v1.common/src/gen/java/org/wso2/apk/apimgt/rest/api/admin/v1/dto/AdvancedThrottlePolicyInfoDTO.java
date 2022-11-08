package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.ThrottleLimitDTO;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.ThrottlePolicyDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class AdvancedThrottlePolicyInfoDTO extends ThrottlePolicyDTO  {
  
  private ThrottleLimitDTO defaultLimit;


  /**
   **/
  public AdvancedThrottlePolicyInfoDTO defaultLimit(ThrottleLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("defaultLimit")
  public ThrottleLimitDTO getDefaultLimit() {
    return defaultLimit;
  }
  public void setDefaultLimit(ThrottleLimitDTO defaultLimit) {
    this.defaultLimit = defaultLimit;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    AdvancedThrottlePolicyInfoDTO advancedThrottlePolicyInfo = (AdvancedThrottlePolicyInfoDTO) o;
    return Objects.equals(defaultLimit, advancedThrottlePolicyInfo.defaultLimit) &&
        super.equals(o);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultLimit, super.hashCode());
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AdvancedThrottlePolicyInfoDTO {\n");
    sb.append("    ").append(toIndentedString(super.toString())).append("\n");
    sb.append("    defaultLimit: ").append(toIndentedString(defaultLimit)).append("\n");
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

