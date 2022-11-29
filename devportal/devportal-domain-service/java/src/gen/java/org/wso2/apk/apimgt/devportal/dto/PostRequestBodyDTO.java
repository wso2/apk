package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import java.util.Objects;



public class PostRequestBodyDTO   {
  
  private String content;

  private String category;


  /**
   * Content of the comment 
   **/
  public PostRequestBodyDTO content(String content) {
    this.content = content;
    return this;
  }

  
  @ApiModelProperty(example = "This is a comment", required = true, value = "Content of the comment ")
  @JsonProperty("content")
  @NotNull
 @Size(max=512)  public String getContent() {
    return content;
  }
  public void setContent(String content) {
    this.content = content;
  }


  /**
   * Category of the comment 
   **/
  public PostRequestBodyDTO category(String category) {
    this.category = category;
    return this;
  }

  
  @ApiModelProperty(example = "general", value = "Category of the comment ")
  @JsonProperty("category")
  public String getCategory() {
    return category;
  }
  public void setCategory(String category) {
    this.category = category;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    PostRequestBodyDTO postRequestBody = (PostRequestBodyDTO) o;
    return Objects.equals(content, postRequestBody.content) &&
        Objects.equals(category, postRequestBody.category);
  }

  @Override
  public int hashCode() {
    return Objects.hash(content, category);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class PostRequestBodyDTO {\n");
    
    sb.append("    content: ").append(toIndentedString(content)).append("\n");
    sb.append("    category: ").append(toIndentedString(category)).append("\n");
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

