package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.SettingsKeyManagerConfigurationInnerDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class SettingsDTO   {
  
  private List<String> scopes = null;

  private List<SettingsKeyManagerConfigurationInnerDTO> keyManagerConfiguration = null;

  private Boolean analyticsEnabled;


  /**
   **/
  public SettingsDTO scopes(List<String> scopes) {
    this.scopes = scopes;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("scopes")
  public List<String> getScopes() {
    return scopes;
  }
  public void setScopes(List<String> scopes) {
    this.scopes = scopes;
  }

  public SettingsDTO addScopesItem(String scopesItem) {
    if (this.scopes == null) {
      this.scopes = new ArrayList<>();
    }
    this.scopes.add(scopesItem);
    return this;
  }


  /**
   **/
  public SettingsDTO keyManagerConfiguration(List<SettingsKeyManagerConfigurationInnerDTO> keyManagerConfiguration) {
    this.keyManagerConfiguration = keyManagerConfiguration;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("keyManagerConfiguration")
  public List<SettingsKeyManagerConfigurationInnerDTO> getKeyManagerConfiguration() {
    return keyManagerConfiguration;
  }
  public void setKeyManagerConfiguration(List<SettingsKeyManagerConfigurationInnerDTO> keyManagerConfiguration) {
    this.keyManagerConfiguration = keyManagerConfiguration;
  }

  public SettingsDTO addKeyManagerConfigurationItem(SettingsKeyManagerConfigurationInnerDTO keyManagerConfigurationItem) {
    if (this.keyManagerConfiguration == null) {
      this.keyManagerConfiguration = new ArrayList<>();
    }
    this.keyManagerConfiguration.add(keyManagerConfigurationItem);
    return this;
  }


  /**
   * To determine whether analytics is enabled or not
   **/
  public SettingsDTO analyticsEnabled(Boolean analyticsEnabled) {
    this.analyticsEnabled = analyticsEnabled;
    return this;
  }

  
  @ApiModelProperty(example = "false", value = "To determine whether analytics is enabled or not")
  @JsonProperty("analyticsEnabled")
  public Boolean getAnalyticsEnabled() {
    return analyticsEnabled;
  }
  public void setAnalyticsEnabled(Boolean analyticsEnabled) {
    this.analyticsEnabled = analyticsEnabled;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    SettingsDTO settings = (SettingsDTO) o;
    return Objects.equals(scopes, settings.scopes) &&
        Objects.equals(keyManagerConfiguration, settings.keyManagerConfiguration) &&
        Objects.equals(analyticsEnabled, settings.analyticsEnabled);
  }

  @Override
  public int hashCode() {
    return Objects.hash(scopes, keyManagerConfiguration, analyticsEnabled);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SettingsDTO {\n");
    
    sb.append("    scopes: ").append(toIndentedString(scopes)).append("\n");
    sb.append("    keyManagerConfiguration: ").append(toIndentedString(keyManagerConfiguration)).append("\n");
    sb.append("    analyticsEnabled: ").append(toIndentedString(analyticsEnabled)).append("\n");
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

