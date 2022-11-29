package org.wso2.apk.apimgt.admin.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.wso2.apk.apimgt.admin.dto.HeaderConditionDTO;
import org.wso2.apk.apimgt.admin.dto.IPConditionDTO;
import org.wso2.apk.apimgt.admin.dto.JWTClaimsConditionDTO;
import org.wso2.apk.apimgt.admin.dto.QueryParameterConditionDTO;
import javax.validation.constraints.*;

/**
 * Conditions used for Throttling
 **/

import io.swagger.annotations.*;
import java.util.Objects;


@ApiModel(description = "Conditions used for Throttling")
public class ThrottleConditionDTO   {
  

public enum TypeEnum {

    HEADERCONDITION(String.valueOf("HEADERCONDITION")), IPCONDITION(String.valueOf("IPCONDITION")), JWTCLAIMSCONDITION(String.valueOf("JWTCLAIMSCONDITION")), QUERYPARAMETERCONDITION(String.valueOf("QUERYPARAMETERCONDITION"));


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

  private Boolean invertCondition = false;

  private HeaderConditionDTO headerCondition;

  private IPConditionDTO ipCondition;

  private JWTClaimsConditionDTO jwtClaimsCondition;

  private QueryParameterConditionDTO queryParameterCondition;


  /**
   * Type of the throttling condition. Allowed values are \&quot;HEADERCONDITION\&quot;, \&quot;IPCONDITION\&quot;, \&quot;JWTCLAIMSCONDITION\&quot; and \&quot;QUERYPARAMETERCONDITION\&quot;. 
   **/
  public ThrottleConditionDTO type(TypeEnum type) {
    this.type = type;
    return this;
  }

  
  @ApiModelProperty(required = true, value = "Type of the throttling condition. Allowed values are \"HEADERCONDITION\", \"IPCONDITION\", \"JWTCLAIMSCONDITION\" and \"QUERYPARAMETERCONDITION\". ")
  @JsonProperty("type")
  @NotNull
  public TypeEnum getType() {
    return type;
  }
  public void setType(TypeEnum type) {
    this.type = type;
  }


  /**
   * Specifies whether inversion of the condition to be matched against the request.  **Note:** When you add conditional groups for advanced throttling policies, this paramater should have the same value (&#39;true&#39; or &#39;false&#39;) for the same type of conditional group. 
   **/
  public ThrottleConditionDTO invertCondition(Boolean invertCondition) {
    this.invertCondition = invertCondition;
    return this;
  }

  
  @ApiModelProperty(value = "Specifies whether inversion of the condition to be matched against the request.  **Note:** When you add conditional groups for advanced throttling policies, this paramater should have the same value ('true' or 'false') for the same type of conditional group. ")
  @JsonProperty("invertCondition")
  public Boolean getInvertCondition() {
    return invertCondition;
  }
  public void setInvertCondition(Boolean invertCondition) {
    this.invertCondition = invertCondition;
  }


  /**
   **/
  public ThrottleConditionDTO headerCondition(HeaderConditionDTO headerCondition) {
    this.headerCondition = headerCondition;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("headerCondition")
  public HeaderConditionDTO getHeaderCondition() {
    return headerCondition;
  }
  public void setHeaderCondition(HeaderConditionDTO headerCondition) {
    this.headerCondition = headerCondition;
  }


  /**
   **/
  public ThrottleConditionDTO ipCondition(IPConditionDTO ipCondition) {
    this.ipCondition = ipCondition;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("ipCondition")
  public IPConditionDTO getIpCondition() {
    return ipCondition;
  }
  public void setIpCondition(IPConditionDTO ipCondition) {
    this.ipCondition = ipCondition;
  }


  /**
   **/
  public ThrottleConditionDTO jwtClaimsCondition(JWTClaimsConditionDTO jwtClaimsCondition) {
    this.jwtClaimsCondition = jwtClaimsCondition;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("jwtClaimsCondition")
  public JWTClaimsConditionDTO getJwtClaimsCondition() {
    return jwtClaimsCondition;
  }
  public void setJwtClaimsCondition(JWTClaimsConditionDTO jwtClaimsCondition) {
    this.jwtClaimsCondition = jwtClaimsCondition;
  }


  /**
   **/
  public ThrottleConditionDTO queryParameterCondition(QueryParameterConditionDTO queryParameterCondition) {
    this.queryParameterCondition = queryParameterCondition;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("queryParameterCondition")
  public QueryParameterConditionDTO getQueryParameterCondition() {
    return queryParameterCondition;
  }
  public void setQueryParameterCondition(QueryParameterConditionDTO queryParameterCondition) {
    this.queryParameterCondition = queryParameterCondition;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    ThrottleConditionDTO throttleCondition = (ThrottleConditionDTO) o;
    return Objects.equals(type, throttleCondition.type) &&
        Objects.equals(invertCondition, throttleCondition.invertCondition) &&
        Objects.equals(headerCondition, throttleCondition.headerCondition) &&
        Objects.equals(ipCondition, throttleCondition.ipCondition) &&
        Objects.equals(jwtClaimsCondition, throttleCondition.jwtClaimsCondition) &&
        Objects.equals(queryParameterCondition, throttleCondition.queryParameterCondition);
  }

  @Override
  public int hashCode() {
    return Objects.hash(type, invertCondition, headerCondition, ipCondition, jwtClaimsCondition, queryParameterCondition);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class ThrottleConditionDTO {\n");
    
    sb.append("    type: ").append(toIndentedString(type)).append("\n");
    sb.append("    invertCondition: ").append(toIndentedString(invertCondition)).append("\n");
    sb.append("    headerCondition: ").append(toIndentedString(headerCondition)).append("\n");
    sb.append("    ipCondition: ").append(toIndentedString(ipCondition)).append("\n");
    sb.append("    jwtClaimsCondition: ").append(toIndentedString(jwtClaimsCondition)).append("\n");
    sb.append("    queryParameterCondition: ").append(toIndentedString(queryParameterCondition)).append("\n");
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

