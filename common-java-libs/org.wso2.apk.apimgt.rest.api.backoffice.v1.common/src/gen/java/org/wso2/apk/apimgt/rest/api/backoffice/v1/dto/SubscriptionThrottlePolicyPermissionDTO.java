package org.wso2.apk.apimgt.rest.api.backoffice.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class SubscriptionThrottlePolicyPermissionDTO   {
  

public enum PermissionTypeEnum {

    ALLOW(String.valueOf("ALLOW")), DENY(String.valueOf("DENY"));


    private String value;

    PermissionTypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static PermissionTypeEnum fromValue(String value) {
        for (PermissionTypeEnum b : PermissionTypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private PermissionTypeEnum permissionType;

  private List<String> roles = new ArrayList<>();


  /**
   **/
  public SubscriptionThrottlePolicyPermissionDTO permissionType(PermissionTypeEnum permissionType) {
    this.permissionType = permissionType;
    return this;
  }

  
  @ApiModelProperty(example = "deny", required = true, value = "")
  @JsonProperty("permissionType")
  @NotNull
  public PermissionTypeEnum getPermissionType() {
    return permissionType;
  }
  public void setPermissionType(PermissionTypeEnum permissionType) {
    this.permissionType = permissionType;
  }


  /**
   **/
  public SubscriptionThrottlePolicyPermissionDTO roles(List<String> roles) {
    this.roles = roles;
    return this;
  }

  
  @ApiModelProperty(example = "[\"Internal/everyone\"]", required = true, value = "")
  @JsonProperty("roles")
  @NotNull
  public List<String> getRoles() {
    return roles;
  }
  public void setRoles(List<String> roles) {
    this.roles = roles;
  }

  public SubscriptionThrottlePolicyPermissionDTO addRolesItem(String rolesItem) {
    this.roles.add(rolesItem);
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
    SubscriptionThrottlePolicyPermissionDTO subscriptionThrottlePolicyPermission = (SubscriptionThrottlePolicyPermissionDTO) o;
    return Objects.equals(permissionType, subscriptionThrottlePolicyPermission.permissionType) &&
        Objects.equals(roles, subscriptionThrottlePolicyPermission.roles);
  }

  @Override
  public int hashCode() {
    return Objects.hash(permissionType, roles);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SubscriptionThrottlePolicyPermissionDTO {\n");
    
    sb.append("    permissionType: ").append(toIndentedString(permissionType)).append("\n");
    sb.append("    roles: ").append(toIndentedString(roles)).append("\n");
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

