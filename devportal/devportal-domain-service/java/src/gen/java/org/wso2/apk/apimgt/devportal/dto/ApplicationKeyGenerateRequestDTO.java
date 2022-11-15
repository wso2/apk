package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import javax.validation.constraints.*;


import java.util.Objects;



public class ApplicationKeyGenerateRequestDTO   {
  

public enum KeyTypeEnum {

    PRODUCTION(String.valueOf("PRODUCTION")), SANDBOX(String.valueOf("SANDBOX"));


    private String value;

    KeyTypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static KeyTypeEnum fromValue(String value) {
        for (KeyTypeEnum b : KeyTypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private KeyTypeEnum keyType;

  private String keyManager;

  private List<String> grantTypesToBeSupported = new ArrayList<>();

  private String callbackUrl;

  private List<String> scopes = null;

  private String validityTime;

  private String clientId;

  private String clientSecret;

  private Object additionalProperties;


  /**
   **/
  public ApplicationKeyGenerateRequestDTO keyType(KeyTypeEnum keyType) {
    this.keyType = keyType;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "")
  @JsonProperty("keyType")
  @NotNull
  public KeyTypeEnum getKeyType() {
    return keyType;
  }
  public void setKeyType(KeyTypeEnum keyType) {
    this.keyType = keyType;
  }


  /**
   * key Manager to Generate Keys
   **/
  public ApplicationKeyGenerateRequestDTO keyManager(String keyManager) {
    this.keyManager = keyManager;
    return this;
  }

  
  @ApiModelProperty(example = "Resident Key Manager", value = "key Manager to Generate Keys")
  @JsonProperty("keyManager")
  public String getKeyManager() {
    return keyManager;
  }
  public void setKeyManager(String keyManager) {
    this.keyManager = keyManager;
  }


  /**
   * Grant types that should be supported by the application
   **/
  public ApplicationKeyGenerateRequestDTO grantTypesToBeSupported(List<String> grantTypesToBeSupported) {
    this.grantTypesToBeSupported = grantTypesToBeSupported;
    return this;
  }

  
  @ApiModelProperty(example = "[\"password\",\"client_credentials\"]", required = true, value = "Grant types that should be supported by the application")
  @JsonProperty("grantTypesToBeSupported")
  @NotNull
  public List<String> getGrantTypesToBeSupported() {
    return grantTypesToBeSupported;
  }
  public void setGrantTypesToBeSupported(List<String> grantTypesToBeSupported) {
    this.grantTypesToBeSupported = grantTypesToBeSupported;
  }

  public ApplicationKeyGenerateRequestDTO addGrantTypesToBeSupportedItem(String grantTypesToBeSupportedItem) {
    this.grantTypesToBeSupported.add(grantTypesToBeSupportedItem);
    return this;
  }


  /**
   * Callback URL
   **/
  public ApplicationKeyGenerateRequestDTO callbackUrl(String callbackUrl) {
    this.callbackUrl = callbackUrl;
    return this;
  }

  
  @ApiModelProperty(example = "http://sample.com/callback/url", value = "Callback URL")
  @JsonProperty("callbackUrl")
  public String getCallbackUrl() {
    return callbackUrl;
  }
  public void setCallbackUrl(String callbackUrl) {
    this.callbackUrl = callbackUrl;
  }


  /**
   * Allowed scopes for the access token
   **/
  public ApplicationKeyGenerateRequestDTO scopes(List<String> scopes) {
    this.scopes = scopes;
    return this;
  }

  
  @ApiModelProperty(example = "[\"am_application_scope\",\"default\"]", value = "Allowed scopes for the access token")
  @JsonProperty("scopes")
  public List<String> getScopes() {
    return scopes;
  }
  public void setScopes(List<String> scopes) {
    this.scopes = scopes;
  }

  public ApplicationKeyGenerateRequestDTO addScopesItem(String scopesItem) {
    if (this.scopes == null) {
      this.scopes = new ArrayList<>();
    }
    this.scopes.add(scopesItem);
    return this;
  }


  /**
   **/
  public ApplicationKeyGenerateRequestDTO validityTime(String validityTime) {
    this.validityTime = validityTime;
    return this;
  }

  
  @ApiModelProperty(example = "3600", value = "")
  @JsonProperty("validityTime")
  public String getValidityTime() {
    return validityTime;
  }
  public void setValidityTime(String validityTime) {
    this.validityTime = validityTime;
  }


  /**
   * Client ID for generating access token.
   **/
  public ApplicationKeyGenerateRequestDTO clientId(String clientId) {
    this.clientId = clientId;
    return this;
  }

  
  @ApiModelProperty(example = "sZzoeSCI_vL2cjSXZQmsmV8JEyga", value = "Client ID for generating access token.")
  @JsonProperty("clientId")
  public String getClientId() {
    return clientId;
  }
  public void setClientId(String clientId) {
    this.clientId = clientId;
  }


  /**
   * Client secret for generating access token. This is given together with the client Id.
   **/
  public ApplicationKeyGenerateRequestDTO clientSecret(String clientSecret) {
    this.clientSecret = clientSecret;
    return this;
  }

  
  @ApiModelProperty(example = "nrs3YAP4htxnz_DqpvGhf9Um04oa", value = "Client secret for generating access token. This is given together with the client Id.")
  @JsonProperty("clientSecret")
  public String getClientSecret() {
    return clientSecret;
  }
  public void setClientSecret(String clientSecret) {
    this.clientSecret = clientSecret;
  }


  /**
   * Additional properties needed.
   **/
  public ApplicationKeyGenerateRequestDTO additionalProperties(Object additionalProperties) {
    this.additionalProperties = additionalProperties;
    return this;
  }

  
  @ApiModelProperty(example = "{}", value = "Additional properties needed.")
  @JsonProperty("additionalProperties")
  public Object getAdditionalProperties() {
    return additionalProperties;
  }
  public void setAdditionalProperties(Object additionalProperties) {
    this.additionalProperties = additionalProperties;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ApplicationKeyGenerateRequestDTO applicationKeyGenerateRequest = (ApplicationKeyGenerateRequestDTO) o;
    return Objects.equals(keyType, applicationKeyGenerateRequest.keyType) &&
        Objects.equals(keyManager, applicationKeyGenerateRequest.keyManager) &&
        Objects.equals(grantTypesToBeSupported, applicationKeyGenerateRequest.grantTypesToBeSupported) &&
        Objects.equals(callbackUrl, applicationKeyGenerateRequest.callbackUrl) &&
        Objects.equals(scopes, applicationKeyGenerateRequest.scopes) &&
        Objects.equals(validityTime, applicationKeyGenerateRequest.validityTime) &&
        Objects.equals(clientId, applicationKeyGenerateRequest.clientId) &&
        Objects.equals(clientSecret, applicationKeyGenerateRequest.clientSecret) &&
        Objects.equals(additionalProperties, applicationKeyGenerateRequest.additionalProperties);
  }

  @Override
  public int hashCode() {
    return Objects.hash(keyType, keyManager, grantTypesToBeSupported, callbackUrl, scopes, validityTime, clientId, clientSecret, additionalProperties);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ApplicationKeyGenerateRequestDTO {\n");
    
    sb.append("    keyType: ").append(toIndentedString(keyType)).append("\n");
    sb.append("    keyManager: ").append(toIndentedString(keyManager)).append("\n");
    sb.append("    grantTypesToBeSupported: ").append(toIndentedString(grantTypesToBeSupported)).append("\n");
    sb.append("    callbackUrl: ").append(toIndentedString(callbackUrl)).append("\n");
    sb.append("    scopes: ").append(toIndentedString(scopes)).append("\n");
    sb.append("    validityTime: ").append(toIndentedString(validityTime)).append("\n");
    sb.append("    clientId: ").append(toIndentedString(clientId)).append("\n");
    sb.append("    clientSecret: ").append(toIndentedString(clientSecret)).append("\n");
    sb.append("    additionalProperties: ").append(toIndentedString(additionalProperties)).append("\n");
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

