package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.WorkflowInfoDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class WorkflowListDTO   {
  
  private Integer count;

  private String next;

  private String previous;

  private List<WorkflowInfoDTO> _list = null;


  /**
   * Number of workflow processes returned. 
   **/
  public WorkflowListDTO count(Integer count) {
    this.count = count;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "Number of workflow processes returned. ")
  @JsonProperty("count")
  public Integer getCount() {
    return count;
  }
  public void setCount(Integer count) {
    this.count = count;
  }


  /**
   * Link to the next subset of resources qualified. Empty if no more resources are to be returned. 
   **/
  public WorkflowListDTO next(String next) {
    this.next = next;
    return this;
  }

  
  @ApiModelProperty(example = "/workflows?limit=1&offset=2&user=", value = "Link to the next subset of resources qualified. Empty if no more resources are to be returned. ")
  @JsonProperty("next")
  public String getNext() {
    return next;
  }
  public void setNext(String next) {
    this.next = next;
  }


  /**
   * Link to the previous subset of resources qualified. Empty if current subset is the first subset returned. 
   **/
  public WorkflowListDTO previous(String previous) {
    this.previous = previous;
    return this;
  }

  
  @ApiModelProperty(example = "/workflows?limit=1&offset=0&user=", value = "Link to the previous subset of resources qualified. Empty if current subset is the first subset returned. ")
  @JsonProperty("previous")
  public String getPrevious() {
    return previous;
  }
  public void setPrevious(String previous) {
    this.previous = previous;
  }


  /**
   **/
  public WorkflowListDTO _list(List<WorkflowInfoDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<WorkflowInfoDTO> getList() {
    return _list;
  }
  public void setList(List<WorkflowInfoDTO> _list) {
    this._list = _list;
  }

  public WorkflowListDTO addListItem(WorkflowInfoDTO _listItem) {
    if (this._list == null) {
      this._list = new ArrayList<>();
    }
    this._list.add(_listItem);
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
    WorkflowListDTO workflowList = (WorkflowListDTO) o;
    return Objects.equals(count, workflowList.count) &&
        Objects.equals(next, workflowList.next) &&
        Objects.equals(previous, workflowList.previous) &&
        Objects.equals(_list, workflowList._list);
  }

  @Override
  public int hashCode() {
    return Objects.hash(count, next, previous, _list);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class WorkflowListDTO {\n");
    
    sb.append("    count: ").append(toIndentedString(count)).append("\n");
    sb.append("    next: ").append(toIndentedString(next)).append("\n");
    sb.append("    previous: ").append(toIndentedString(previous)).append("\n");
    sb.append("    _list: ").append(toIndentedString(_list)).append("\n");
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

