package org.wso2.apk.apimgt.rest.api.backoffice.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.Valid;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;



public class LifecycleStateDTO   {
  
    private String state = null;
    private List<LifecycleStateAvailableTransitionsDTO> availableTransitions = new ArrayList<LifecycleStateAvailableTransitionsDTO>();

  /**
   **/
  public LifecycleStateDTO state(String state) {
    this.state = state;
    return this;
  }

  
  @ApiModelProperty(example = "Created", value = "")
  @JsonProperty("state")
  public String getState() {
    return state;
  }
  public void setState(String state) {
    this.state = state;
  }

  /**
   **/
  public LifecycleStateDTO availableTransitions(List<LifecycleStateAvailableTransitionsDTO> availableTransitions) {
    this.availableTransitions = availableTransitions;
    return this;
  }

  
  @ApiModelProperty(value = "")
      @Valid
  @JsonProperty("availableTransitions")
  public List<LifecycleStateAvailableTransitionsDTO> getAvailableTransitions() {
    return availableTransitions;
  }
  public void setAvailableTransitions(List<LifecycleStateAvailableTransitionsDTO> availableTransitions) {
    this.availableTransitions = availableTransitions;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    LifecycleStateDTO lifecycleState = (LifecycleStateDTO) o;
    return Objects.equals(state, lifecycleState.state) &&
        Objects.equals(availableTransitions, lifecycleState.availableTransitions);
  }

  @Override
  public int hashCode() {
    return Objects.hash(state, availableTransitions);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class LifecycleStateDTO {\n");
    
    sb.append("    state: ").append(toIndentedString(state)).append("\n");
    sb.append("    availableTransitions: ").append(toIndentedString(availableTransitions)).append("\n");
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

