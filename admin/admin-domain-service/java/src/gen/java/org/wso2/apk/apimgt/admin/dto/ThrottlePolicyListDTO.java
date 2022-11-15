package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.admin.dto.PaginationDTO;
import org.wso2.apk.apimgt.admin.dto.ThrottlePolicyDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class ThrottlePolicyListDTO   {
  
  private Integer count;

  private List<ThrottlePolicyDTO> _list = null;

  private PaginationDTO pagination;


  /**
   * Number of Throttling Policies returned. 
   **/
  public ThrottlePolicyListDTO count(Integer count) {
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
  public ThrottlePolicyListDTO _list(List<ThrottlePolicyDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<ThrottlePolicyDTO> getList() {
    return _list;
  }
  public void setList(List<ThrottlePolicyDTO> _list) {
    this._list = _list;
  }

  public ThrottlePolicyListDTO addListItem(ThrottlePolicyDTO _listItem) {
    if (this._list == null) {
      this._list = new ArrayList<>();
    }
    this._list.add(_listItem);
    return this;
  }


  /**
   **/
  public ThrottlePolicyListDTO pagination(PaginationDTO pagination) {
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
    ThrottlePolicyListDTO throttlePolicyList = (ThrottlePolicyListDTO) o;
    return Objects.equals(count, throttlePolicyList.count) &&
        Objects.equals(_list, throttlePolicyList._list) &&
        Objects.equals(pagination, throttlePolicyList.pagination);
  }

  @Override
  public int hashCode() {
    return Objects.hash(count, _list, pagination);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ThrottlePolicyListDTO {\n");
    
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

