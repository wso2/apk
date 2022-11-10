package org.wso2.apk.apimgt.rest.api.backoffice.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIAdditionalPropertiesValueDTO;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIBusinessInformationDTO;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIDeploymentDTO;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIMonetizationInfoDTO;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIOperationsDTO;
import org.wso2.apk.apimgt.rest.api.backoffice.v1.dto.APIRevisionDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class APIDTO   {
  
  private String id;

  private String name;

  private String description;

  private String context;

  private String version;


public enum TypeEnum {

    HTTP(String.valueOf("HTTP")), WS(String.valueOf("WS")), SOAPTOREST(String.valueOf("SOAPTOREST")), SOAP(String.valueOf("SOAP")), GRAPHQL(String.valueOf("GRAPHQL")), WEBSUB(String.valueOf("WEBSUB")), SSE(String.valueOf("SSE")), WEBHOOK(String.valueOf("WEBHOOK")), ASYNC(String.valueOf("ASYNC"));


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

  private TypeEnum type = TypeEnum.HTTP;

  private List<String> transport = null;

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

  private List<String> categories = null;

  private Map<String, APIAdditionalPropertiesValueDTO> additionalProperties = null;

  private String createdTime;

  private String lastUpdatedTime;

  private List<APIOperationsDTO> operations = null;

  private String apiUsagePolicy;

  private APIMonetizationInfoDTO monetization;

  private APIBusinessInformationDTO businessInformation;

  private APIRevisionDTO revision;

  private List<APIDeploymentDTO> deployments = null;


  /**
   * UUID of the API 
   **/
  public APIDTO id(String id) {
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
   **/
  public APIDTO name(String name) {
    this.name = name;
    return this;
  }

  
  @ApiModelProperty(example = "PizzaShackAPI", required = true, value = "")
  @JsonProperty("name")
  @NotNull
 @Pattern(regexp="(^[^~!@#;:%^*()+={}|\\\\<>\"',&$\\[\\]/]*$)") @Size(min=1,max=60)  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
  }


  /**
   **/
  public APIDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "This is a simple API for Pizza Shack online pizza delivery store.", value = "")
  @JsonProperty("description")
 @Size(max=32766)  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }


  /**
   **/
  public APIDTO context(String context) {
    this.context = context;
    return this;
  }

  
  @ApiModelProperty(example = "pizza", required = true, value = "")
  @JsonProperty("context")
  @NotNull
 @Size(min=1,max=232)  public String getContext() {
    return context;
  }
  public void setContext(String context) {
    this.context = context;
  }


  /**
   **/
  public APIDTO version(String version) {
    this.version = version;
    return this;
  }

  
  @ApiModelProperty(example = "1.0.0", required = true, value = "")
  @JsonProperty("version")
  @NotNull
 @Pattern(regexp="^[^~!@#;:%^*()+={}|\\\\<>\"',&/$\\[\\]\\s+/]+$") @Size(min=1,max=30)  public String getVersion() {
    return version;
  }
  public void setVersion(String version) {
    this.version = version;
  }


  /**
   * The api creation type to be used. Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL, WEBSUB, SSE, WEBHOOK, ASYNC
   **/
  public APIDTO type(TypeEnum type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(example = "HTTP", value = "The api creation type to be used. Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL, WEBSUB, SSE, WEBHOOK, ASYNC")
  @JsonProperty("type")
  public TypeEnum getType() {
    return type;
  }
  public void setType(TypeEnum type) {
    this.type = type;
  }


  /**
   * Supported transports for the API (http and/or https). 
   **/
  public APIDTO transport(List<String> transport) {
    this.transport = transport;
    return this;
  }

  
  @ApiModelProperty(example = "[\"http\",\"https\"]", value = "Supported transports for the API (http and/or https). ")
  @JsonProperty("transport")
  public List<String> getTransport() {
    return transport;
  }
  public void setTransport(List<String> transport) {
    this.transport = transport;
  }

  public APIDTO addTransportItem(String transportItem) {
    if (this.transport == null) {
      this.transport = new ArrayList<>();
    }
    this.transport.add(transportItem);
    return this;
  }


  /**
   **/
  public APIDTO hasThumbnail(Boolean hasThumbnail) {
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
  public APIDTO state(StateEnum state) {
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
  public APIDTO tags(List<String> tags) {
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

  public APIDTO addTagsItem(String tagsItem) {
    if (this.tags == null) {
      this.tags = new ArrayList<>();
    }
    this.tags.add(tagsItem);
    return this;
  }


  /**
   * API categories 
   **/
  public APIDTO categories(List<String> categories) {
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

  public APIDTO addCategoriesItem(String categoriesItem) {
    if (this.categories == null) {
      this.categories = new ArrayList<>();
    }
    this.categories.add(categoriesItem);
    return this;
  }


  /**
   **/
  public APIDTO additionalProperties(Map<String, APIAdditionalPropertiesValueDTO> additionalProperties) {
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


  public APIDTO putAdditionalPropertiesItem(String key, APIAdditionalPropertiesValueDTO additionalPropertiesItem) {
    if (this.additionalProperties == null) {
      this.additionalProperties = new HashMap<>();
    }
    this.additionalProperties.put(key, additionalPropertiesItem);
    return this;
  }

  /**
   **/
  public APIDTO createdTime(String createdTime) {
    this.createdTime = createdTime;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("createdTime")
  public String getCreatedTime() {
    return createdTime;
  }
  public void setCreatedTime(String createdTime) {
    this.createdTime = createdTime;
  }


  /**
   **/
  public APIDTO lastUpdatedTime(String lastUpdatedTime) {
    this.lastUpdatedTime = lastUpdatedTime;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("lastUpdatedTime")
  public String getLastUpdatedTime() {
    return lastUpdatedTime;
  }
  public void setLastUpdatedTime(String lastUpdatedTime) {
    this.lastUpdatedTime = lastUpdatedTime;
  }


  /**
   **/
  public APIDTO operations(List<APIOperationsDTO> operations) {
    this.operations = operations;
    return this;
  }

  
  @ApiModelProperty(example = "[{\"target\":\"/order/{orderId}\",\"verb\":\"POST\",\"usagePlan\":\"Unlimited\"},{\"target\":\"/menu\",\"verb\":\"GET\",\"usagePlan\":\"Unlimited\"}]", value = "")
  @JsonProperty("operations")
  public List<APIOperationsDTO> getOperations() {
    return operations;
  }
  public void setOperations(List<APIOperationsDTO> operations) {
    this.operations = operations;
  }

  public APIDTO addOperationsItem(APIOperationsDTO operationsItem) {
    if (this.operations == null) {
      this.operations = new ArrayList<>();
    }
    this.operations.add(operationsItem);
    return this;
  }


  /**
   * The API level usage policy selected for the particular Runtime API
   **/
  public APIDTO apiUsagePolicy(String apiUsagePolicy) {
    this.apiUsagePolicy = apiUsagePolicy;
    return this;
  }

  
  @ApiModelProperty(example = "Unlimited", value = "The API level usage policy selected for the particular Runtime API")
  @JsonProperty("apiUsagePolicy")
  public String getApiUsagePolicy() {
    return apiUsagePolicy;
  }
  public void setApiUsagePolicy(String apiUsagePolicy) {
    this.apiUsagePolicy = apiUsagePolicy;
  }


  /**
   **/
  public APIDTO monetization(APIMonetizationInfoDTO monetization) {
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
  public APIDTO businessInformation(APIBusinessInformationDTO businessInformation) {
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
   **/
  public APIDTO revision(APIRevisionDTO revision) {
    this.revision = revision;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("revision")
  public APIRevisionDTO getRevision() {
    return revision;
  }
  public void setRevision(APIRevisionDTO revision) {
    this.revision = revision;
  }


  /**
   **/
  public APIDTO deployments(List<APIDeploymentDTO> deployments) {
    this.deployments = deployments;
    return this;
  }

  
  @ApiModelProperty(example = "[{\"name\":\"US\",\"deployedTime\":\"2022-10-28T06:13:35.024Z\"},{\"name\":\"Europe\",\"deployedTime\":\"2022-10-28T06:13:35.024Z\"}]", value = "")
  @JsonProperty("deployments")
  public List<APIDeploymentDTO> getDeployments() {
    return deployments;
  }
  public void setDeployments(List<APIDeploymentDTO> deployments) {
    this.deployments = deployments;
  }

  public APIDTO addDeploymentsItem(APIDeploymentDTO deploymentsItem) {
    if (this.deployments == null) {
      this.deployments = new ArrayList<>();
    }
    this.deployments.add(deploymentsItem);
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
    APIDTO API = (APIDTO) o;
    return Objects.equals(id, API.id) &&
        Objects.equals(name, API.name) &&
        Objects.equals(description, API.description) &&
        Objects.equals(context, API.context) &&
        Objects.equals(version, API.version) &&
        Objects.equals(type, API.type) &&
        Objects.equals(transport, API.transport) &&
        Objects.equals(hasThumbnail, API.hasThumbnail) &&
        Objects.equals(state, API.state) &&
        Objects.equals(tags, API.tags) &&
        Objects.equals(categories, API.categories) &&
        Objects.equals(additionalProperties, API.additionalProperties) &&
        Objects.equals(createdTime, API.createdTime) &&
        Objects.equals(lastUpdatedTime, API.lastUpdatedTime) &&
        Objects.equals(operations, API.operations) &&
        Objects.equals(apiUsagePolicy, API.apiUsagePolicy) &&
        Objects.equals(monetization, API.monetization) &&
        Objects.equals(businessInformation, API.businessInformation) &&
        Objects.equals(revision, API.revision) &&
        Objects.equals(deployments, API.deployments);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, name, description, context, version, type, transport, hasThumbnail, state, tags, categories, additionalProperties, createdTime, lastUpdatedTime, operations, apiUsagePolicy, monetization, businessInformation, revision, deployments);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class APIDTO {\n");
    
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
    sb.append("    context: ").append(toIndentedString(context)).append("\n");
    sb.append("    version: ").append(toIndentedString(version)).append("\n");
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    transport: ").append(toIndentedString(transport)).append("\n");
    sb.append("    hasThumbnail: ").append(toIndentedString(hasThumbnail)).append("\n");
    sb.append("    state: ").append(toIndentedString(state)).append("\n");
    sb.append("    tags: ").append(toIndentedString(tags)).append("\n");
    sb.append("    categories: ").append(toIndentedString(categories)).append("\n");
    sb.append("    additionalProperties: ").append(toIndentedString(additionalProperties)).append("\n");
    sb.append("    createdTime: ").append(toIndentedString(createdTime)).append("\n");
    sb.append("    lastUpdatedTime: ").append(toIndentedString(lastUpdatedTime)).append("\n");
    sb.append("    operations: ").append(toIndentedString(operations)).append("\n");
    sb.append("    apiUsagePolicy: ").append(toIndentedString(apiUsagePolicy)).append("\n");
    sb.append("    monetization: ").append(toIndentedString(monetization)).append("\n");
    sb.append("    businessInformation: ").append(toIndentedString(businessInformation)).append("\n");
    sb.append("    revision: ").append(toIndentedString(revision)).append("\n");
    sb.append("    deployments: ").append(toIndentedString(deployments)).append("\n");
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

