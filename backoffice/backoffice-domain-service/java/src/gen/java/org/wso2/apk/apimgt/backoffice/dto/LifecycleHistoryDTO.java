package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class LifecycleHistoryDTO   {
  
  private Integer count;

  private List<LifecycleHistoryItemDTO> _list = null;


  /**
   **/
  public LifecycleHistoryDTO count(Integer count) {
    this.count = count;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "")
  @JsonProperty("count")
  public Integer getCount() {
    return count;
  }
  public void setCount(Integer count) {
    this.count = count;
  }


  /**
   **/
  public LifecycleHistoryDTO _list(List<LifecycleHistoryItemDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<LifecycleHistoryItemDTO> getList() {
    return _list;
  }
  public void setList(List<LifecycleHistoryItemDTO> _list) {
    this._list = _list;
  }

  public LifecycleHistoryDTO addListItem(LifecycleHistoryItemDTO _listItem) {
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
    LifecycleHistoryDTO lifecycleHistory = (LifecycleHistoryDTO) o;
    return Objects.equals(count, lifecycleHistory.count) &&
        Objects.equals(_list, lifecycleHistory._list);
  }

  @Override
  public int hashCode() {
    return Objects.hash(count, _list);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class LifecycleHistoryDTO {\n");
    
    sb.append("    count: ").append(toIndentedString(count)).append("\n");
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

