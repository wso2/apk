package org.wso2.apk.apimgt.devportal.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.swagger.annotations.ApiModelProperty;


import java.util.Objects;



public class DocumentSearchResultDTO extends SearchResultDTO  {
  

public enum DocTypeEnum {

    HOWTO(String.valueOf("HOWTO")), SAMPLES(String.valueOf("SAMPLES")), PUBLIC_FORUM(String.valueOf("PUBLIC_FORUM")), SUPPORT_FORUM(String.valueOf("SUPPORT_FORUM")), API_MESSAGE_FORMAT(String.valueOf("API_MESSAGE_FORMAT")), SWAGGER_DOC(String.valueOf("SWAGGER_DOC")), OTHER(String.valueOf("OTHER"));


    private String value;

    DocTypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static DocTypeEnum fromValue(String value) {
        for (DocTypeEnum b : DocTypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private DocTypeEnum docType;

  private String summary;


public enum SourceTypeEnum {

    INLINE(String.valueOf("INLINE")), URL(String.valueOf("URL")), FILE(String.valueOf("FILE")), MARKDOWN(String.valueOf("MARKDOWN"));


    private String value;

    SourceTypeEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static SourceTypeEnum fromValue(String value) {
        for (SourceTypeEnum b : SourceTypeEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private SourceTypeEnum sourceType;

  private String sourceUrl;

  private String otherTypeName;


public enum VisibilityEnum {

    OWNER_ONLY(String.valueOf("OWNER_ONLY")), PRIVATE(String.valueOf("PRIVATE")), API_LEVEL(String.valueOf("API_LEVEL"));


    private String value;

    VisibilityEnum(String v) {
        value = v;
    }

    public String value() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

    public static VisibilityEnum fromValue(String value) {
        for (VisibilityEnum b : VisibilityEnum.values()) {
            if (b.value.equals(value)) {
                return b;
            }
        }
        throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }
}

  private VisibilityEnum visibility;

  private String apiName;

  private String apiVersion;

  private String apiProvider;

  private String apiUUID;


  /**
   **/
  public DocumentSearchResultDTO docType(DocTypeEnum docType) {
    this.docType = docType;
    return this;
  }

  
  @ApiModelProperty(example = "HOWTO", value = "")
  @JsonProperty("docType")
  public DocTypeEnum getDocType() {
    return docType;
  }
  public void setDocType(DocTypeEnum docType) {
    this.docType = docType;
  }


  /**
   **/
  public DocumentSearchResultDTO summary(String summary) {
    this.summary = summary;
    return this;
  }

  
  @ApiModelProperty(example = "Summary of Calculator Documentation", value = "")
  @JsonProperty("summary")
  public String getSummary() {
    return summary;
  }
  public void setSummary(String summary) {
    this.summary = summary;
  }


  /**
   **/
  public DocumentSearchResultDTO sourceType(SourceTypeEnum sourceType) {
    this.sourceType = sourceType;
    return this;
  }

  
  @ApiModelProperty(example = "INLINE", value = "")
  @JsonProperty("sourceType")
  public SourceTypeEnum getSourceType() {
    return sourceType;
  }
  public void setSourceType(SourceTypeEnum sourceType) {
    this.sourceType = sourceType;
  }


  /**
   **/
  public DocumentSearchResultDTO sourceUrl(String sourceUrl) {
    this.sourceUrl = sourceUrl;
    return this;
  }

  
  @ApiModelProperty(example = "", value = "")
  @JsonProperty("sourceUrl")
  public String getSourceUrl() {
    return sourceUrl;
  }
  public void setSourceUrl(String sourceUrl) {
    this.sourceUrl = sourceUrl;
  }


  /**
   **/
  public DocumentSearchResultDTO otherTypeName(String otherTypeName) {
    this.otherTypeName = otherTypeName;
    return this;
  }

  
  @ApiModelProperty(example = "", value = "")
  @JsonProperty("otherTypeName")
  public String getOtherTypeName() {
    return otherTypeName;
  }
  public void setOtherTypeName(String otherTypeName) {
    this.otherTypeName = otherTypeName;
  }


  /**
   **/
  public DocumentSearchResultDTO visibility(VisibilityEnum visibility) {
    this.visibility = visibility;
    return this;
  }

  
  @ApiModelProperty(example = "API_LEVEL", value = "")
  @JsonProperty("visibility")
  public VisibilityEnum getVisibility() {
    return visibility;
  }
  public void setVisibility(VisibilityEnum visibility) {
    this.visibility = visibility;
  }


  /**
   * The name of the associated API
   **/
  public DocumentSearchResultDTO apiName(String apiName) {
    this.apiName = apiName;
    return this;
  }

  
  @ApiModelProperty(example = "TestAPI", value = "The name of the associated API")
  @JsonProperty("apiName")
  public String getApiName() {
    return apiName;
  }
  public void setApiName(String apiName) {
    this.apiName = apiName;
  }


  /**
   * The version of the associated API
   **/
  public DocumentSearchResultDTO apiVersion(String apiVersion) {
    this.apiVersion = apiVersion;
    return this;
  }

  
  @ApiModelProperty(example = "1.0.0", value = "The version of the associated API")
  @JsonProperty("apiVersion")
  public String getApiVersion() {
    return apiVersion;
  }
  public void setApiVersion(String apiVersion) {
    this.apiVersion = apiVersion;
  }


  /**
   **/
  public DocumentSearchResultDTO apiProvider(String apiProvider) {
    this.apiProvider = apiProvider;
    return this;
  }

  
  @ApiModelProperty(example = "admin", value = "")
  @JsonProperty("apiProvider")
  public String getApiProvider() {
    return apiProvider;
  }
  public void setApiProvider(String apiProvider) {
    this.apiProvider = apiProvider;
  }


  /**
   **/
  public DocumentSearchResultDTO apiUUID(String apiUUID) {
    this.apiUUID = apiUUID;
    return this;
  }

  
  @ApiModelProperty(value = "")
  @JsonProperty("apiUUID")
  public String getApiUUID() {
    return apiUUID;
  }
  public void setApiUUID(String apiUUID) {
    this.apiUUID = apiUUID;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    DocumentSearchResultDTO documentSearchResult = (DocumentSearchResultDTO) o;
    return Objects.equals(docType, documentSearchResult.docType) &&
        Objects.equals(summary, documentSearchResult.summary) &&
        Objects.equals(sourceType, documentSearchResult.sourceType) &&
        Objects.equals(sourceUrl, documentSearchResult.sourceUrl) &&
        Objects.equals(otherTypeName, documentSearchResult.otherTypeName) &&
        Objects.equals(visibility, documentSearchResult.visibility) &&
        Objects.equals(apiName, documentSearchResult.apiName) &&
        Objects.equals(apiVersion, documentSearchResult.apiVersion) &&
        Objects.equals(apiProvider, documentSearchResult.apiProvider) &&
        Objects.equals(apiUUID, documentSearchResult.apiUUID) &&
        super.equals(o);
  }

  @Override
  public int hashCode() {
    return Objects.hash(docType, summary, sourceType, sourceUrl, otherTypeName, visibility, apiName, apiVersion, apiProvider, apiUUID, super.hashCode());
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class DocumentSearchResultDTO {\n");
    sb.append("    ").append(toIndentedString(super.toString())).append("\n");
    sb.append("    docType: ").append(toIndentedString(docType)).append("\n");
    sb.append("    summary: ").append(toIndentedString(summary)).append("\n");
    sb.append("    sourceType: ").append(toIndentedString(sourceType)).append("\n");
    sb.append("    sourceUrl: ").append(toIndentedString(sourceUrl)).append("\n");
    sb.append("    otherTypeName: ").append(toIndentedString(otherTypeName)).append("\n");
    sb.append("    visibility: ").append(toIndentedString(visibility)).append("\n");
    sb.append("    apiName: ").append(toIndentedString(apiName)).append("\n");
    sb.append("    apiVersion: ").append(toIndentedString(apiVersion)).append("\n");
    sb.append("    apiProvider: ").append(toIndentedString(apiProvider)).append("\n");
    sb.append("    apiUUID: ").append(toIndentedString(apiUUID)).append("\n");
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

