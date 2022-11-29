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



public class ScopeDTO   {
  
  private String tag;

  private String name;

  private String description;

  private List<String> roles = null;


  /**
   * Portal name. 
   **/
  public ScopeDTO tag(String tag) {
    this.tag = tag;
    return this;
  }

  
  @ApiModelProperty(example = "publisher", value = "Portal name. ")
  @JsonProperty("tag")
  public String getTag() {
    return tag;
  }
  public void setTag(String tag) {
    this.tag = tag;
  }


  /**
   * Scope name. 
   **/
  public ScopeDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "apim:api_publish", value = "Scope name. ")
  @JsonProperty("name")
  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   * About scope. 
   **/
  public ScopeDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "Publish API", value = "About scope. ")
  @JsonProperty("description")
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   * Roles for the particular scope. 
   **/
  public ScopeDTO roles(List<String> roles) {
    this.roles = roles;
    return this;
  }

  
  @ApiModelProperty(example = "[\"admin\",\"Internal/publisher\"]", value = "Roles for the particular scope. ")
  @JsonProperty("roles")
  public List<String> getRoles() {
    return roles;
  }
  public void setRoles(List<String> roles) {
    this.roles = roles;
  }

  public ScopeDTO addRolesItem(String rolesItem) {
    if (this.roles == null) {
      this.roles = new ArrayList<>();
    }
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
    ScopeDTO scope = (ScopeDTO) o;
    return Objects.equals(tag, scope.tag) &&
        Objects.equals(name, scope.name) &&
        Objects.equals(description, scope.description) &&
        Objects.equals(roles, scope.roles);
  }

  @Override
  public int hashCode() {
    return Objects.hash(tag, name, description, roles);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ScopeDTO {\n");
    
    sb.append("    tag: ").append(toIndentedString(tag)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
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

