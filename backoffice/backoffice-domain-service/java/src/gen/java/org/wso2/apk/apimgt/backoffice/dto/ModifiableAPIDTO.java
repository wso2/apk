package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import javax.validation.constraints.*;


import java.util.Objects;



public class ModifiableAPIDTO   {
  
  private String id;

  private String name;

  private String context;

  private String description;

  private Boolean hasThumbnail;


public enum StateEnum {

    CREATED(String.valueOf("CREATED")), PUBLISHED(String.valueOf("PUBLISHED"));


    private String value;

    StateEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static StateEnum fromValue(String value) {
        for (StateEnum b : StateEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private StateEnum state = StateEnum.CREATED;

  private List<String> tags = null;

  private Map<String, APIAdditionalPropertiesValueDTO> additionalProperties = null;

  private APIMonetizationInfoDTO monetization;

  private APIBusinessInformationDTO businessInformation;

  private List<String> categories = null;


  /**
   * UUID of the API 
   **/
  public ModifiableAPIDTO id(String id) {
    this.id = id;
    return this;
  }

  
  @ApiModelProperty(example = "01234567-0123-0123-0123-012345678901", value = "UUID of the API ")
  @JsonProperty("id")
  public String getId() {
    return id;
  }
  public void setId(String id) {
    this.id = id;
  }


  /**
   * Name of the API
   **/
  public ModifiableAPIDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "PizzaShackAPI", required = true, value = "Name of the API")
  @JsonProperty("name")
  @NotNull
 @Size(min=1,max=50)  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public ModifiableAPIDTO context(String context) {
    this.context = context;
    return this;
  }

  
  @ApiModelProperty(example = "pizzaproduct", value = "")
  @JsonProperty("context")
 @Size(min=1,max=60)  public String getContext() {
    return context;
  }
  public void setContext(String context) {
    this.context = context;
  }


  /**
   * A brief description about the API
   **/
  public ModifiableAPIDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "This is a simple API for Pizza Shack online pizza delivery store", value = "A brief description about the API")
  @JsonProperty("description")
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   **/
  public ModifiableAPIDTO hasThumbnail(Boolean hasThumbnail) {
    this.hasThumbnail = hasThumbnail;
    return this;
  }

  
  @ApiModelProperty(example = "false", value = "")
  @JsonProperty("hasThumbnail")
  public Boolean getHasThumbnail() {
    return hasThumbnail;
  }
  public void setHasThumbnail(Boolean hasThumbnail) {
    this.hasThumbnail = hasThumbnail;
  }


  /**
   * State of the API. Only published APIs are visible on the Developer Portal 
   **/
  public ModifiableAPIDTO state(StateEnum state) {
    this.state = state;
    return this;
  }

  
  @ApiModelProperty(value = "State of the API. Only published APIs are visible on the Developer Portal ")
  @JsonProperty("state")
  public StateEnum getState() {
    return state;
  }
  public void setState(StateEnum state) {
    this.state = state;
  }


  /**
   **/
  public ModifiableAPIDTO tags(List<String> tags) {
    this.tags = tags;
    return this;
  }

  
  @ApiModelProperty(example = "[\"pizza\",\"food\"]", value = "")
  @JsonProperty("tags")
  public List<String> getTags() {
    return tags;
  }
  public void setTags(List<String> tags) {
    this.tags = tags;
  }

  public ModifiableAPIDTO addTagsItem(String tagsItem) {
    if (this.tags == null) {
      this.tags = new ArrayList<>();
    }
    this.tags.add(tagsItem);
    return this;
  }


  /**
   **/
  public ModifiableAPIDTO additionalProperties(Map<String, APIAdditionalPropertiesValueDTO> additionalProperties) {
    this.additionalProperties = additionalProperties;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("additionalProperties")
  public Map<String, APIAdditionalPropertiesValueDTO> getAdditionalProperties() {
    return additionalProperties;
  }
  public void setAdditionalProperties(Map<String, APIAdditionalPropertiesValueDTO> additionalProperties) {
    this.additionalProperties = additionalProperties;
  }


  public ModifiableAPIDTO putAdditionalPropertiesItem(String key, APIAdditionalPropertiesValueDTO additionalPropertiesItem) {
    if (this.additionalProperties == null) {
      this.additionalProperties = new HashMap<>();
    }
    this.additionalProperties.put(key, additionalPropertiesItem);
    return this;
  }

  /**
   **/
  public ModifiableAPIDTO monetization(APIMonetizationInfoDTO monetization) {
    this.monetization = monetization;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("monetization")
  public APIMonetizationInfoDTO getMonetization() {
    return monetization;
  }
  public void setMonetization(APIMonetizationInfoDTO monetization) {
    this.monetization = monetization;
  }


  /**
   **/
  public ModifiableAPIDTO businessInformation(APIBusinessInformationDTO businessInformation) {
    this.businessInformation = businessInformation;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("businessInformation")
  public APIBusinessInformationDTO getBusinessInformation() {
    return businessInformation;
  }
  public void setBusinessInformation(APIBusinessInformationDTO businessInformation) {
    this.businessInformation = businessInformation;
  }


  /**
   * API categories 
   **/
  public ModifiableAPIDTO categories(List<String> categories) {
    this.categories = categories;
    return this;
  }

  
  @ApiModelProperty(example = "[]", value = "API categories ")
  @JsonProperty("categories")
  public List<String> getCategories() {
    return categories;
  }
  public void setCategories(List<String> categories) {
    this.categories = categories;
  }

  public ModifiableAPIDTO addCategoriesItem(String categoriesItem) {
    if (this.categories == null) {
      this.categories = new ArrayList<>();
    }
    this.categories.add(categoriesItem);
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
    ModifiableAPIDTO modifiableAPI = (ModifiableAPIDTO) o;
    return Objects.equals(id, modifiableAPI.id) &&
        Objects.equals(name, modifiableAPI.name) &&
        Objects.equals(context, modifiableAPI.context) &&
        Objects.equals(description, modifiableAPI.description) &&
        Objects.equals(hasThumbnail, modifiableAPI.hasThumbnail) &&
        Objects.equals(state, modifiableAPI.state) &&
        Objects.equals(tags, modifiableAPI.tags) &&
        Objects.equals(additionalProperties, modifiableAPI.additionalProperties) &&
        Objects.equals(monetization, modifiableAPI.monetization) &&
        Objects.equals(businessInformation, modifiableAPI.businessInformation) &&
        Objects.equals(categories, modifiableAPI.categories);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, name, context, description, hasThumbnail, state, tags, additionalProperties, monetization, businessInformation, categories);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ModifiableAPIDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    context: ").append(toIndentedString(context)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
    sb.append("    hasThumbnail: ").append(toIndentedString(hasThumbnail)).append("\n");
    sb.append("    state: ").append(toIndentedString(state)).append("\n");
    sb.append("    tags: ").append(toIndentedString(tags)).append("\n");
    sb.append("    additionalProperties: ").append(toIndentedString(additionalProperties)).append("\n");
    sb.append("    monetization: ").append(toIndentedString(monetization)).append("\n");
    sb.append("    businessInformation: ").append(toIndentedString(businessInformation)).append("\n");
    sb.append("    categories: ").append(toIndentedString(categories)).append("\n");
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

