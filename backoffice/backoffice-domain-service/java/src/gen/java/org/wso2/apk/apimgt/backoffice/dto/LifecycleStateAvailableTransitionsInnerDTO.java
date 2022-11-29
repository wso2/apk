package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class LifecycleStateAvailableTransitionsInnerDTO   {
  
  private String event;

  private String targetState;


  /**
   **/
  public LifecycleStateAvailableTransitionsInnerDTO event(String event) {
    this.event = event;
    return this;
  }

  
  @ApiModelProperty(example = "Publish", value = "")
  @JsonProperty("event")
  public String getEvent() {
    return event;
  }
  public void setEvent(String event) {
    this.event = event;
  }


  /**
   **/
  public LifecycleStateAvailableTransitionsInnerDTO targetState(String targetState) {
    this.targetState = targetState;
    return this;
  }

  
  @ApiModelProperty(example = "Published", value = "")
  @JsonProperty("targetState")
  public String getTargetState() {
    return targetState;
  }
  public void setTargetState(String targetState) {
    this.targetState = targetState;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    LifecycleStateAvailableTransitionsInnerDTO lifecycleStateAvailableTransitionsInner = (LifecycleStateAvailableTransitionsInnerDTO) o;
    return Objects.equals(event, lifecycleStateAvailableTransitionsInner.event) &&
        Objects.equals(targetState, lifecycleStateAvailableTransitionsInner.targetState);
  }

  @Override
  public int hashCode() {
    return Objects.hash(event, targetState);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class LifecycleStateAvailableTransitionsInnerDTO {\n");
    
    sb.append("    event: ").append(toIndentedString(event)).append("\n");
    sb.append("    targetState: ").append(toIndentedString(targetState)).append("\n");
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

