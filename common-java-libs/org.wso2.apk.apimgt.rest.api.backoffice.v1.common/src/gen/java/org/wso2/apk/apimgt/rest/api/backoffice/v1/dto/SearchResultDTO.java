package org.wso2.apk.apimgt.rest.api.backoffice.v1.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;


@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "name", visible = true)
@JsonSubTypes({
})

public class SearchResultDTO   {
  
  private String id;

  private String name;


public enum TypeEnum {

    DOC(String.valueOf("DOC")), API(String.valueOf("API"));


    private String value;

    TypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static TypeEnum fromValue(String value) {
        for (TypeEnum b : TypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private TypeEnum type;

  private String transportType;


  /**
   **/
  public SearchResultDTO id(String id) {
    this.id = id;
    return this;
  }

  
  @ApiModelProperty(example = "01234567-0123-0123-0123-012345678901", value = "")
  @JsonProperty("id")
  public String getId() {
    return id;
  }
  public void setId(String id) {
    this.id = id;
  }


  /**
   **/
  public SearchResultDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "TestAPI", required = true, value = "")
  @JsonProperty("name")
  @NotNull
  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public SearchResultDTO type(TypeEnum type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(example = "API", value = "")
  @JsonProperty("type")
  public TypeEnum getType() {
    return type;
  }
  public void setType(TypeEnum type) {
    this.type = type;
  }


  /**
   * Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL
   **/
  public SearchResultDTO transportType(String transportType) {
    this.transportType = transportType;
    return this;
  }

  
  @ApiModelProperty(value = "Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL")
  @JsonProperty("transportType")
  public String getTransportType() {
    return transportType;
  }
  public void setTransportType(String transportType) {
    this.transportType = transportType;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    SearchResultDTO searchResult = (SearchResultDTO) o;
    return Objects.equals(id, searchResult.id) &&
        Objects.equals(name, searchResult.name) &&
        Objects.equals(type, searchResult.type) &&
        Objects.equals(transportType, searchResult.transportType);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, name, type, transportType);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class SearchResultDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    transportType: ").append(toIndentedString(transportType)).append("\n");
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

