package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class ScopeListDTO   {
  
  private Integer count;

  private List<ScopeInfoDTO> _list = null;

  private PaginationDTO pagination;


  /**
   * Number of results returned. 
   **/
  public ScopeListDTO count(Integer count) {
    this.count = count;
    return this;
  }

  
  @ApiModelProperty(example = "1", value = "Number of results returned. ")
  @JsonProperty("count")
  public Integer getCount() {
    return count;
  }
  public void setCount(Integer count) {
    this.count = count;
  }


  /**
   **/
  public ScopeListDTO _list(List<ScopeInfoDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<ScopeInfoDTO> getList() {
    return _list;
  }
  public void setList(List<ScopeInfoDTO> _list) {
    this._list = _list;
  }

  public ScopeListDTO addListItem(ScopeInfoDTO _listItem) {
    if (this._list == null) {
      this._list = new ArrayList<>();
    }
    this._list.add(_listItem);
    return this;
  }


  /**
   **/
  public ScopeListDTO pagination(PaginationDTO pagination) {
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
    ScopeListDTO scopeList = (ScopeListDTO) o;
    return Objects.equals(count, scopeList.count) &&
        Objects.equals(_list, scopeList._list) &&
        Objects.equals(pagination, scopeList.pagination);
  }

  @Override
  public int hashCode() {
    return Objects.hash(count, _list, pagination);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ScopeListDTO {\n");
    
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

