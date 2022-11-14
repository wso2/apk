package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class GatewayEnvironmentProtocolURIDTO   {
  
  private String protocol;

  private String endpointURI;


  /**
   **/
  public GatewayEnvironmentProtocolURIDTO protocol(String protocol) {
    this.protocol = protocol;
    return this;
  }

  
  @ApiModelProperty(example = "default", required = true, value = "")
  @JsonProperty("protocol")
  @NotNull
  public String getProtocol() {
    return protocol;
  }
  public void setProtocol(String protocol) {
    this.protocol = protocol;
  }


  /**
   **/
  public GatewayEnvironmentProtocolURIDTO endpointURI(String endpointURI) {
    this.endpointURI = endpointURI;
    return this;
  }

  
  @ApiModelProperty(example = "default", required = true, value = "")
  @JsonProperty("endpointURI")
  @NotNull
  public String getEndpointURI() {
    return endpointURI;
  }
  public void setEndpointURI(String endpointURI) {
    this.endpointURI = endpointURI;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    GatewayEnvironmentProtocolURIDTO gatewayEnvironmentProtocolURI = (GatewayEnvironmentProtocolURIDTO) o;
    return Objects.equals(protocol, gatewayEnvironmentProtocolURI.protocol) &&
        Objects.equals(endpointURI, gatewayEnvironmentProtocolURI.endpointURI);
  }

  @Override
  public int hashCode() {
    return Objects.hash(protocol, endpointURI);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class GatewayEnvironmentProtocolURIDTO {\n");
    
    sb.append("    protocol: ").append(toIndentedString(protocol)).append("\n");
    sb.append("    endpointURI: ").append(toIndentedString(endpointURI)).append("\n");
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

