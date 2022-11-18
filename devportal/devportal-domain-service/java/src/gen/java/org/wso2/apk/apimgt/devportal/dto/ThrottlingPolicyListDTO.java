package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class ThrottlingPolicyListDTO   {
  
  private Integer count;

  private List<ThrottlingPolicyDTO> _list = null;

  private PaginationDTO pagination;


  /**
   * Number of Throttling Policies returned. 
   **/
  public ThrottlingPolicyListDTO count(Integer count) {
    this.count = count;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "Number of Throttling Policies returned. ")
  @JsonProperty("count")
  public Integer getCount() {
    return count;
  }
  public void setCount(Integer count) {
    this.count = count;
  }


  /**
   **/
  public ThrottlingPolicyListDTO _list(List<ThrottlingPolicyDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<ThrottlingPolicyDTO> getList() {
    return _list;
  }
  public void setList(List<ThrottlingPolicyDTO> _list) {
    this._list = _list;
  }

  public ThrottlingPolicyListDTO addListItem(ThrottlingPolicyDTO _listItem) {
    if (this._list == null) {
      this._list = new ArrayList<>();
    }
    this._list.add(_listItem);
    return this;
  }


  /**
   **/
  public ThrottlingPolicyListDTO pagination(PaginationDTO pagination) {
    this.pagination = pagination;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("pagination")
  public PaginationDTO getPagination() {
    return pagination;
  }
  public void setPagination(PaginationDTO pagination) {
    this.pagination = pagination;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ThrottlingPolicyListDTO throttlingPolicyList = (ThrottlingPolicyListDTO) o;
    return Objects.equals(count, throttlingPolicyList.count) &&
        Objects.equals(_list, throttlingPolicyList._list) &&
        Objects.equals(pagination, throttlingPolicyList.pagination);
  }

  @Override
  public int hashCode() {
    return Objects.hash(count, _list, pagination);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ThrottlingPolicyListDTO {\n");
    
    sb.append("    count: ").append(toIndentedString(count)).append("\n");
    sb.append("    _list: ").append(toIndentedString(_list)).append("\n");
    sb.append("    pagination: ").append(toIndentedString(pagination)).append("\n");
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

