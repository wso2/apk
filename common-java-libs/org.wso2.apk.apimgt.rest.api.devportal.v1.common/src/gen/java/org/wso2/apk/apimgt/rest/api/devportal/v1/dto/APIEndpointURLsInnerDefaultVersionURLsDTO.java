package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APIEndpointURLsInnerDefaultVersionURLsDTO   {
  
  private String http;

  private String https;

  private String ws;

  private String wss;


  /**
   * HTTP environment default URL
   **/
  public APIEndpointURLsInnerDefaultVersionURLsDTO http(String http) {
    this.http = http;
    return this;
  }

  
  @ApiModelProperty(example = "http://localhost:8280/phoneverify/", value = "HTTP environment default URL")
  @JsonProperty("http")
  public String getHttp() {
    return http;
  }
  public void setHttp(String http) {
    this.http = http;
  }


  /**
   * HTTPS environment default URL
   **/
  public APIEndpointURLsInnerDefaultVersionURLsDTO https(String https) {
    this.https = https;
    return this;
  }

  
  @ApiModelProperty(example = "https://localhost:8243/phoneverify/", value = "HTTPS environment default URL")
  @JsonProperty("https")
  public String getHttps() {
    return https;
  }
  public void setHttps(String https) {
    this.https = https;
  }


  /**
   * WS environment default URL
   **/
  public APIEndpointURLsInnerDefaultVersionURLsDTO ws(String ws) {
    this.ws = ws;
    return this;
  }

  
  @ApiModelProperty(example = "ws://localhost:9099/phoneverify/", value = "WS environment default URL")
  @JsonProperty("ws")
  public String getWs() {
    return ws;
  }
  public void setWs(String ws) {
    this.ws = ws;
  }


  /**
   * WSS environment default URL
   **/
  public APIEndpointURLsInnerDefaultVersionURLsDTO wss(String wss) {
    this.wss = wss;
    return this;
  }

  
  @ApiModelProperty(example = "wss://localhost:9099/phoneverify/", value = "WSS environment default URL")
  @JsonProperty("wss")
  public String getWss() {
    return wss;
  }
  public void setWss(String wss) {
    this.wss = wss;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIEndpointURLsInnerDefaultVersionURLsDTO apIEndpointURLsInnerDefaultVersionURLs = (APIEndpointURLsInnerDefaultVersionURLsDTO) o;
    return Objects.equals(http, apIEndpointURLsInnerDefaultVersionURLs.http) &&
        Objects.equals(https, apIEndpointURLsInnerDefaultVersionURLs.https) &&
        Objects.equals(ws, apIEndpointURLsInnerDefaultVersionURLs.ws) &&
        Objects.equals(wss, apIEndpointURLsInnerDefaultVersionURLs.wss);
  }

  @Override
  public int hashCode() {
    return Objects.hash(http, https, ws, wss);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIEndpointURLsInnerDefaultVersionURLsDTO {\n");
    
    sb.append("    http: ").append(toIndentedString(http)).append("\n");
    sb.append("    https: ").append(toIndentedString(https)).append("\n");
    sb.append("    ws: ").append(toIndentedString(ws)).append("\n");
    sb.append("    wss: ").append(toIndentedString(wss)).append("\n");
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

