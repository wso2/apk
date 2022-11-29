package org.wso2.apk.apimgt.admin.dto;

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



public class RoleAliasDTO   {
  
  private String role;

  private List<String> aliases = null;


  /**
   * The original role
   **/
  public RoleAliasDTO role(String role) {
    this.role = role;
    return this;
  }

  
  @ApiModelProperty(example = "Internal/subscriber", value = "The original role")
  @JsonProperty("role")
  public String getRole() {
    return role;
  }
  public void setRole(String role) {
    this.role = role;
  }


  /**
   * The role mapping for role alias
   **/
  public RoleAliasDTO aliases(List<String> aliases) {
    this.aliases = aliases;
    return this;
  }

  
  @ApiModelProperty(example = "[\"Subscriber\",\"Internal/subscriber\"]", value = "The role mapping for role alias")
  @JsonProperty("aliases")
  public List<String> getAliases() {
    return aliases;
  }
  public void setAliases(List<String> aliases) {
    this.aliases = aliases;
  }

  public RoleAliasDTO addAliasesItem(String aliasesItem) {
    if (this.aliases == null) {
      this.aliases = new ArrayList<>();
    }
    this.aliases.add(aliasesItem);
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
    RoleAliasDTO roleAlias = (RoleAliasDTO) o;
    return Objects.equals(role, roleAlias.role) &&
        Objects.equals(aliases, roleAlias.aliases);
  }

  @Override
  public int hashCode() {
    return Objects.hash(role, aliases);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class RoleAliasDTO {\n");
    
    sb.append("    role: ").append(toIndentedString(role)).append("\n");
    sb.append("    aliases: ").append(toIndentedString(aliases)).append("\n");
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

