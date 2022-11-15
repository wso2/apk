package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;


import java.util.Objects;



public class GraphQLSchemaTypeListDTO   {
  
  private List<GraphQLSchemaTypeDTO> typeList = null;


  /**
   **/
  public GraphQLSchemaTypeListDTO typeList(List<GraphQLSchemaTypeDTO> typeList) {
    this.typeList = typeList;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("typeList")
  public List<GraphQLSchemaTypeDTO> getTypeList() {
    return typeList;
  }
  public void setTypeList(List<GraphQLSchemaTypeDTO> typeList) {
    this.typeList = typeList;
  }

  public GraphQLSchemaTypeListDTO addTypeListItem(GraphQLSchemaTypeDTO typeListItem) {
    if (this.typeList == null) {
      this.typeList = new ArrayList<>();
    }
    this.typeList.add(typeListItem);
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
    GraphQLSchemaTypeListDTO graphQLSchemaTypeList = (GraphQLSchemaTypeListDTO) o;
    return Objects.equals(typeList, graphQLSchemaTypeList.typeList);
  }

  @Override
  public int hashCode() {
    return Objects.hash(typeList);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class GraphQLSchemaTypeListDTO {\n");
    
    sb.append("    typeList: ").append(toIndentedString(typeList)).append("\n");
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

