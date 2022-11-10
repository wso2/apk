package org.wso2.apk.apimgt.rest.api.devportal.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.rest.api.devportal.v1.dto.GraphQLCustomComplexityInfoDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class GraphQLQueryComplexityInfoDTO   {
  
  private List<GraphQLCustomComplexityInfoDTO> _list = null;


  /**
   **/
  public GraphQLQueryComplexityInfoDTO _list(List<GraphQLCustomComplexityInfoDTO> _list) {
    this._list = _list;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("list")
  public List<GraphQLCustomComplexityInfoDTO> getList() {
    return _list;
  }
  public void setList(List<GraphQLCustomComplexityInfoDTO> _list) {
    this._list = _list;
  }

  public GraphQLQueryComplexityInfoDTO addListItem(GraphQLCustomComplexityInfoDTO _listItem) {
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
    GraphQLQueryComplexityInfoDTO graphQLQueryComplexityInfo = (GraphQLQueryComplexityInfoDTO) o;
    return Objects.equals(_list, graphQLQueryComplexityInfo._list);
  }

  @Override
  public int hashCode() {
    return Objects.hash(_list);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class GraphQLQueryComplexityInfoDTO {\n");
    
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

