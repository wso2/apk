package org.wso2.apk.apimgt.backoffice.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;
import javax.validation.constraints.*;


import java.util.Objects;



public class EventCountLimitDTO   {
  
  private String timeUnit;

  private Integer unitTime;

  private Long eventCount;


  /**
   * Unit of the time. Allowed values are \&quot;sec\&quot;, \&quot;min\&quot;, \&quot;hour\&quot;, \&quot;day\&quot;
   **/
  public EventCountLimitDTO timeUnit(String timeUnit) {
    this.timeUnit = timeUnit;
    return this;
  }

  
  @ApiModelProperty(example = "min", required = true, value = "Unit of the time. Allowed values are \"sec\", \"min\", \"hour\", \"day\"")
  @JsonProperty("timeUnit")
  @NotNull
  public String getTimeUnit() {
    return timeUnit;
  }
  public void setTimeUnit(String timeUnit) {
    this.timeUnit = timeUnit;
  }


  /**
   * Time limit that the usage limit applies.
   **/
  public EventCountLimitDTO unitTime(Integer unitTime) {
    this.unitTime = unitTime;
    return this;
  }

  
  @ApiModelProperty(example = "10", required = true, value = "Time limit that the usage limit applies.")
  @JsonProperty("unitTime")
  @NotNull
  public Integer getUnitTime() {
    return unitTime;
  }
  public void setUnitTime(Integer unitTime) {
    this.unitTime = unitTime;
  }


  /**
   * Maximum number of events allowed
   **/
  public EventCountLimitDTO eventCount(Long eventCount) {
    this.eventCount = eventCount;
    return this;
  }

  
  @ApiModelProperty(example = "3000", required = true, value = "Maximum number of events allowed")
  @JsonProperty("eventCount")
  @NotNull
  public Long getEventCount() {
    return eventCount;
  }
  public void setEventCount(Long eventCount) {
    this.eventCount = eventCount;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    EventCountLimitDTO eventCountLimit = (EventCountLimitDTO) o;
    return Objects.equals(timeUnit, eventCountLimit.timeUnit) &&
        Objects.equals(unitTime, eventCountLimit.unitTime) &&
        Objects.equals(eventCount, eventCountLimit.eventCount);
  }

  @Override
  public int hashCode() {
    return Objects.hash(timeUnit, unitTime, eventCount);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class EventCountLimitDTO {\n");
    
    sb.append("    timeUnit: ").append(toIndentedString(timeUnit)).append("\n");
    sb.append("    unitTime: ").append(toIndentedString(unitTime)).append("\n");
    sb.append("    eventCount: ").append(toIndentedString(eventCount)).append("\n");
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

