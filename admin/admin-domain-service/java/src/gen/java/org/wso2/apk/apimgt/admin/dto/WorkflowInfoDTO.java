package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class WorkflowInfoDTO   {
  

public enum WorkflowTypeEnum {

    APPLICATION_CREATION(String.valueOf("APPLICATION_CREATION")), SUBSCRIPTION_CREATION(String.valueOf("SUBSCRIPTION_CREATION")), USER_SIGNUP(String.valueOf("USER_SIGNUP")), APPLICATION_REGISTRATION_PRODUCTION(String.valueOf("APPLICATION_REGISTRATION_PRODUCTION")), APPLICATION_REGISTRATION_SANDBOX(String.valueOf("APPLICATION_REGISTRATION_SANDBOX")), APPLICATION_DELETION(String.valueOf("APPLICATION_DELETION")), API_STATE(String.valueOf("API_STATE")), API_PRODUCT_STATE(String.valueOf("API_PRODUCT_STATE")), SUBSCRIPTION_DELETION(String.valueOf("SUBSCRIPTION_DELETION")), SUBSCRIPTION_UPDATE(String.valueOf("SUBSCRIPTION_UPDATE"));


    private String value;

    WorkflowTypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static WorkflowTypeEnum fromValue(String value) {
        for (WorkflowTypeEnum b : WorkflowTypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private WorkflowTypeEnum workflowType;


public enum WorkflowStatusEnum {

    APPROVED(String.valueOf("APPROVED")), CREATED(String.valueOf("CREATED"));


    private String value;

    WorkflowStatusEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static WorkflowStatusEnum fromValue(String value) {
        for (WorkflowStatusEnum b : WorkflowStatusEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private WorkflowStatusEnum workflowStatus;

  private String createdTime;

  private String updatedTime;

  private String referenceId;

  private Object properties;

  private String description;


  /**
   * Type of the Workflow Request. It shows which type of request is it. 
   **/
  public WorkflowInfoDTO workflowType(WorkflowTypeEnum workflowType) {
    this.workflowType = workflowType;
    return this;
  }

  
  @ApiModelProperty(example = "APPLICATION_CREATION", value = "Type of the Workflow Request. It shows which type of request is it. ")
  @JsonProperty("workflowType")
  public WorkflowTypeEnum getWorkflowType() {
    return workflowType;
  }
  public void setWorkflowType(WorkflowTypeEnum workflowType) {
    this.workflowType = workflowType;
  }


  /**
   * Show the Status of the the workflow request whether it is approved or created. 
   **/
  public WorkflowInfoDTO workflowStatus(WorkflowStatusEnum workflowStatus) {
    this.workflowStatus = workflowStatus;
    return this;
  }

  
  @ApiModelProperty(example = "APPROVED", value = "Show the Status of the the workflow request whether it is approved or created. ")
  @JsonProperty("workflowStatus")
  public WorkflowStatusEnum getWorkflowStatus() {
    return workflowStatus;
  }
  public void setWorkflowStatus(WorkflowStatusEnum workflowStatus) {
    this.workflowStatus = workflowStatus;
  }


  /**
   * Time of the the workflow request created. 
   **/
  public WorkflowInfoDTO createdTime(String createdTime) {
    this.createdTime = createdTime;
    return this;
  }

  
  @ApiModelProperty(example = "2020-02-10 10:10:19.704", value = "Time of the the workflow request created. ")
  @JsonProperty("createdTime")
  public String getCreatedTime() {
    return createdTime;
  }
  public void setCreatedTime(String createdTime) {
    this.createdTime = createdTime;
  }


  /**
   * Time of the the workflow request updated. 
   **/
  public WorkflowInfoDTO updatedTime(String updatedTime) {
    this.updatedTime = updatedTime;
    return this;
  }

  
  @ApiModelProperty(example = "2020-02-10 10:10:19.704", value = "Time of the the workflow request updated. ")
  @JsonProperty("updatedTime")
  public String getUpdatedTime() {
    return updatedTime;
  }
  public void setUpdatedTime(String updatedTime) {
    this.updatedTime = updatedTime;
  }


  /**
   * Workflow external reference is used to identify the workflow requests uniquely. 
   **/
  public WorkflowInfoDTO referenceId(String referenceId) {
    this.referenceId = referenceId;
    return this;
  }

  
  @ApiModelProperty(example = "5871244b-d6f3-466e-8995-8accd1e64303", value = "Workflow external reference is used to identify the workflow requests uniquely. ")
  @JsonProperty("referenceId")
  public String getReferenceId() {
    return referenceId;
  }
  public void setReferenceId(String referenceId) {
    this.referenceId = referenceId;
  }


  /**
   **/
  public WorkflowInfoDTO properties(Object properties) {
    this.properties = properties;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("properties")
  public Object getProperties() {
    return properties;
  }
  public void setProperties(Object properties) {
    this.properties = properties;
  }


  /**
   * description is a message with basic details about the workflow request. 
   **/
  public WorkflowInfoDTO description(String description) {
    this.description = description;
    return this;
  }

  
  @ApiModelProperty(example = "Approve application [APP1] creation request from application creator - admin with throttling tier - 10MinPer", value = "description is a message with basic details about the workflow request. ")
  @JsonProperty("description")
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    WorkflowInfoDTO workflowInfo = (WorkflowInfoDTO) o;
    return Objects.equals(workflowType, workflowInfo.workflowType) &&
        Objects.equals(workflowStatus, workflowInfo.workflowStatus) &&
        Objects.equals(createdTime, workflowInfo.createdTime) &&
        Objects.equals(updatedTime, workflowInfo.updatedTime) &&
        Objects.equals(referenceId, workflowInfo.referenceId) &&
        Objects.equals(properties, workflowInfo.properties) &&
        Objects.equals(description, workflowInfo.description);
  }

  @Override
  public int hashCode() {
    return Objects.hash(workflowType, workflowStatus, createdTime, updatedTime, referenceId, properties, description);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class WorkflowInfoDTO {\n");
    
    sb.append("    workflowType: ").append(toIndentedString(workflowType)).append("\n");
    sb.append("    workflowStatus: ").append(toIndentedString(workflowStatus)).append("\n");
    sb.append("    createdTime: ").append(toIndentedString(createdTime)).append("\n");
    sb.append("    updatedTime: ").append(toIndentedString(updatedTime)).append("\n");
    sb.append("    referenceId: ").append(toIndentedString(referenceId)).append("\n");
    sb.append("    properties: ").append(toIndentedString(properties)).append("\n");
    sb.append("    description: ").append(toIndentedString(description)).append("\n");
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

