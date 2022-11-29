package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class APIExternalStoreDTO   {
  
  private String id;

  private String lastUpdatedTime;


  /**
   * The external store identifier, which is a unique value. 
   **/
  public APIExternalStoreDTO id(String id) {
    this.id = id;
    return this;
  }

  
  @ApiModelProperty(example = "Store123#", value = "The external store identifier, which is a unique value. ")
  @JsonProperty("id")
  public String getId() {
    return id;
  }
  public void setId(String id) {
    this.id = id;
  }


  /**
   * The recent timestamp which a given API is updated in the external store. 
   **/
  public APIExternalStoreDTO lastUpdatedTime(String lastUpdatedTime) {
    this.lastUpdatedTime = lastUpdatedTime;
    return this;
  }

  
  @ApiModelProperty(example = "2019-09-09T13:57:16.229", value = "The recent timestamp which a given API is updated in the external store. ")
  @JsonProperty("lastUpdatedTime")
  public String getLastUpdatedTime() {
    return lastUpdatedTime;
  }
  public void setLastUpdatedTime(String lastUpdatedTime) {
    this.lastUpdatedTime = lastUpdatedTime;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    APIExternalStoreDTO apIExternalStore = (APIExternalStoreDTO) o;
    return Objects.equals(id, apIExternalStore.id) &&
        Objects.equals(lastUpdatedTime, apIExternalStore.lastUpdatedTime);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, lastUpdatedTime);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIExternalStoreDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    lastUpdatedTime: ").append(toIndentedString(lastUpdatedTime)).append("\n");
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

