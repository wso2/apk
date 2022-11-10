package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

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



public class SolaceTopicsDTO   {
  
  private List<String> publishTopics = null;

  private List<String> subscribeTopics = null;


  /**
   **/
  public SolaceTopicsDTO publishTopics(List<String> publishTopics) {
    this.publishTopics = publishTopics;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("publishTopics")
  public List<String> getPublishTopics() {
    return publishTopics;
  }
  public void setPublishTopics(List<String> publishTopics) {
    this.publishTopics = publishTopics;
  }

  public SolaceTopicsDTO addPublishTopicsItem(String publishTopicsItem) {
    if (this.publishTopics == null) {
      this.publishTopics = new ArrayList<>();
    }
    this.publishTopics.add(publishTopicsItem);
    return this;
  }


  /**
   **/
  public SolaceTopicsDTO subscribeTopics(List<String> subscribeTopics) {
    this.subscribeTopics = subscribeTopics;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("subscribeTopics")
  public List<String> getSubscribeTopics() {
    return subscribeTopics;
  }
  public void setSubscribeTopics(List<String> subscribeTopics) {
    this.subscribeTopics = subscribeTopics;
  }

  public SolaceTopicsDTO addSubscribeTopicsItem(String subscribeTopicsItem) {
    if (this.subscribeTopics == null) {
      this.subscribeTopics = new ArrayList<>();
    }
    this.subscribeTopics.add(subscribeTopicsItem);
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
    SolaceTopicsDTO solaceTopics = (SolaceTopicsDTO) o;
    return Objects.equals(publishTopics, solaceTopics.publishTopics) &&
        Objects.equals(subscribeTopics, solaceTopics.subscribeTopics);
  }

  @Override
  public int hashCode() {
    return Objects.hash(publishTopics, subscribeTopics);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SolaceTopicsDTO {\n");
    
    sb.append("    publishTopics: ").append(toIndentedString(publishTopics)).append("\n");
    sb.append("    subscribeTopics: ").append(toIndentedString(subscribeTopics)).append("\n");
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

