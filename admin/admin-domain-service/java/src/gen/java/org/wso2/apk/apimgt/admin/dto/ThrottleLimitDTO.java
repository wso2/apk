package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.admin.dto.BandwidthLimitDTO;
import org.wso2.apk.apimgt.admin.dto.EventCountLimitDTO;
import org.wso2.apk.apimgt.admin.dto.RequestCountLimitDTO;
import javax.validation.constraints.*;


import io.swagger.annotations.*;
import java.util.Objects;



public class ThrottleLimitDTO   {
  

public enum TypeEnum {

    REQUESTCOUNTLIMIT(String.valueOf("REQUESTCOUNTLIMIT")), BANDWIDTHLIMIT(String.valueOf("BANDWIDTHLIMIT")), EVENTCOUNTLIMIT(String.valueOf("EVENTCOUNTLIMIT"));


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

  private RequestCountLimitDTO requestCount;

  private BandwidthLimitDTO bandwidth;

  private EventCountLimitDTO eventCount;


  /**
   * Type of the throttling limit. Allowed values are \&quot;REQUESTCOUNTLIMIT\&quot; and \&quot;BANDWIDTHLIMIT\&quot;. Please see schemas of \&quot;RequestCountLimit\&quot; and \&quot;BandwidthLimit\&quot; throttling limit types in Definitions section. 
   **/
  public ThrottleLimitDTO type(TypeEnum type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(example = "REQUESTCOUNTLIMIT", required = true, value = "Type of the throttling limit. Allowed values are \"REQUESTCOUNTLIMIT\" and \"BANDWIDTHLIMIT\". Please see schemas of \"RequestCountLimit\" and \"BandwidthLimit\" throttling limit types in Definitions section. ")
  @JsonProperty("type")
  @NotNull
  public TypeEnum getType() {
    return type;
  }
  public void setType(TypeEnum type) {
    this.type = type;
  }


  /**
   **/
  public ThrottleLimitDTO requestCount(RequestCountLimitDTO requestCount) {
    this.requestCount = requestCount;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("requestCount")
  public RequestCountLimitDTO getRequestCount() {
    return requestCount;
  }
  public void setRequestCount(RequestCountLimitDTO requestCount) {
    this.requestCount = requestCount;
  }


  /**
   **/
  public ThrottleLimitDTO bandwidth(BandwidthLimitDTO bandwidth) {
    this.bandwidth = bandwidth;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("bandwidth")
  public BandwidthLimitDTO getBandwidth() {
    return bandwidth;
  }
  public void setBandwidth(BandwidthLimitDTO bandwidth) {
    this.bandwidth = bandwidth;
  }


  /**
   **/
  public ThrottleLimitDTO eventCount(EventCountLimitDTO eventCount) {
    this.eventCount = eventCount;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("eventCount")
  public EventCountLimitDTO getEventCount() {
    return eventCount;
  }
  public void setEventCount(EventCountLimitDTO eventCount) {
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
    ThrottleLimitDTO throttleLimit = (ThrottleLimitDTO) o;
    return Objects.equals(type, throttleLimit.type) &&
        Objects.equals(requestCount, throttleLimit.requestCount) &&
        Objects.equals(bandwidth, throttleLimit.bandwidth) &&
        Objects.equals(eventCount, throttleLimit.eventCount);
  }

  @Override
  public int hashCode() {
    return Objects.hash(type, requestCount, bandwidth, eventCount);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ThrottleLimitDTO {\n");
    
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    requestCount: ").append(toIndentedString(requestCount)).append("\n");
    sb.append("    bandwidth: ").append(toIndentedString(bandwidth)).append("\n");
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

