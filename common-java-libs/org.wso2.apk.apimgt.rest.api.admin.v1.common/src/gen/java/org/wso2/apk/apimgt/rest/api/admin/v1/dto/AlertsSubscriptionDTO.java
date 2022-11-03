package org.wso2.apk.apimgt.rest.api.admin.v1.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.wso2.apk.apimgt.rest.api.admin.v1.dto.AlertTypeDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class AlertsSubscriptionDTO   {
  
  private List<AlertTypeDTO> alerts = null;

  private List<String> emailList = null;


  /**
   **/
  public AlertsSubscriptionDTO alerts(List<AlertTypeDTO> alerts) {
    this.alerts = alerts;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("alerts")
  public List<AlertTypeDTO> getAlerts() {
    return alerts;
  }
  public void setAlerts(List<AlertTypeDTO> alerts) {
    this.alerts = alerts;
  }

  public AlertsSubscriptionDTO addAlertsItem(AlertTypeDTO alertsItem) {
    if (this.alerts == null) {
      this.alerts = new ArrayList<>();
    }
    this.alerts.add(alertsItem);
    return this;
  }


  /**
   **/
  public AlertsSubscriptionDTO emailList(List<String> emailList) {
    this.emailList = emailList;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("emailList")
  public List<String> getEmailList() {
    return emailList;
  }
  public void setEmailList(List<String> emailList) {
    this.emailList = emailList;
  }

  public AlertsSubscriptionDTO addEmailListItem(String emailListItem) {
    if (this.emailList == null) {
      this.emailList = new ArrayList<>();
    }
    this.emailList.add(emailListItem);
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
    AlertsSubscriptionDTO alertsSubscription = (AlertsSubscriptionDTO) o;
    return Objects.equals(alerts, alertsSubscription.alerts) &&
        Objects.equals(emailList, alertsSubscription.emailList);
  }

  @Override
  public int hashCode() {
    return Objects.hash(alerts, emailList);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class AlertsSubscriptionDTO {\n");
    
    sb.append("    alerts: ").append(toIndentedString(alerts)).append("\n");
    sb.append("    emailList: ").append(toIndentedString(emailList)).append("\n");
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

