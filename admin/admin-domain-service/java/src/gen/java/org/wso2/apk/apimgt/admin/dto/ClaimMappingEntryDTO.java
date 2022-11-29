package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class ClaimMappingEntryDTO   {
  
  private String remoteClaim;

  private String localClaim;


  /**
   **/
  public ClaimMappingEntryDTO remoteClaim(String remoteClaim) {
    this.remoteClaim = remoteClaim;
    return this;
  }

  
  @ApiModelProperty(example = "http://idp.org/username", value = "")
  @JsonProperty("remoteClaim")
  public String getRemoteClaim() {
    return remoteClaim;
  }
  public void setRemoteClaim(String remoteClaim) {
    this.remoteClaim = remoteClaim;
  }


  /**
   **/
  public ClaimMappingEntryDTO localClaim(String localClaim) {
    this.localClaim = localClaim;
    return this;
  }

  
  @ApiModelProperty(example = "http://wso2.org/username", value = "")
  @JsonProperty("localClaim")
  public String getLocalClaim() {
    return localClaim;
  }
  public void setLocalClaim(String localClaim) {
    this.localClaim = localClaim;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ClaimMappingEntryDTO claimMappingEntry = (ClaimMappingEntryDTO) o;
    return Objects.equals(remoteClaim, claimMappingEntry.remoteClaim) &&
        Objects.equals(localClaim, claimMappingEntry.localClaim);
  }

  @Override
  public int hashCode() {
    return Objects.hash(remoteClaim, localClaim);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ClaimMappingEntryDTO {\n");
    
    sb.append("    remoteClaim: ").append(toIndentedString(remoteClaim)).append("\n");
    sb.append("    localClaim: ").append(toIndentedString(localClaim)).append("\n");
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

