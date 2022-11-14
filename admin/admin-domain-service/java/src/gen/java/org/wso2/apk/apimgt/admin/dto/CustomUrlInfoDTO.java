package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.admin.dto.CustomUrlInfoDevPortalDTO;
import javax.validation.constraints.*;

/**
 * The custom url information of the tenant domain
 **/

import io.swagger.annotations.*;
import java.util.Objects;


@ApiModel(description = "The custom url information of the tenant domain")
public class CustomUrlInfoDTO   {
  
  private String tenantDomain;

  private String tenantAdminUsername;

  private Boolean enabled;

  private CustomUrlInfoDevPortalDTO devPortal;


  /**
   **/
  public CustomUrlInfoDTO tenantDomain(String tenantDomain) {
    this.tenantDomain = tenantDomain;
    return this;
  }

  
  @ApiModelProperty(example = "carbon.super", value = "")
  @JsonProperty("tenantDomain")
  public String getTenantDomain() {
    return tenantDomain;
  }
  public void setTenantDomain(String tenantDomain) {
    this.tenantDomain = tenantDomain;
  }


  /**
   **/
  public CustomUrlInfoDTO tenantAdminUsername(String tenantAdminUsername) {
    this.tenantAdminUsername = tenantAdminUsername;
    return this;
  }

  
  @ApiModelProperty(example = "john@foo.com", value = "")
  @JsonProperty("tenantAdminUsername")
  public String getTenantAdminUsername() {
    return tenantAdminUsername;
  }
  public void setTenantAdminUsername(String tenantAdminUsername) {
    this.tenantAdminUsername = tenantAdminUsername;
  }


  /**
   **/
  public CustomUrlInfoDTO enabled(Boolean enabled) {
    this.enabled = enabled;
    return this;
  }

  
  @ApiModelProperty(example = "true", value = "")
  @JsonProperty("enabled")
  public Boolean getEnabled() {
    return enabled;
  }
  public void setEnabled(Boolean enabled) {
    this.enabled = enabled;
  }


  /**
   **/
  public CustomUrlInfoDTO devPortal(CustomUrlInfoDevPortalDTO devPortal) {
    this.devPortal = devPortal;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("devPortal")
  public CustomUrlInfoDevPortalDTO getDevPortal() {
    return devPortal;
  }
  public void setDevPortal(CustomUrlInfoDevPortalDTO devPortal) {
    this.devPortal = devPortal;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    CustomUrlInfoDTO customUrlInfo = (CustomUrlInfoDTO) o;
    return Objects.equals(tenantDomain, customUrlInfo.tenantDomain) &&
        Objects.equals(tenantAdminUsername, customUrlInfo.tenantAdminUsername) &&
        Objects.equals(enabled, customUrlInfo.enabled) &&
        Objects.equals(devPortal, customUrlInfo.devPortal);
  }

  @Override
  public int hashCode() {
    return Objects.hash(tenantDomain, tenantAdminUsername, enabled, devPortal);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class CustomUrlInfoDTO {\n");
    
    sb.append("    tenantDomain: ").append(toIndentedString(tenantDomain)).append("\n");
    sb.append("    tenantAdminUsername: ").append(toIndentedString(tenantAdminUsername)).append("\n");
    sb.append("    enabled: ").append(toIndentedString(enabled)).append("\n");
    sb.append("    devPortal: ").append(toIndentedString(devPortal)).append("\n");
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

