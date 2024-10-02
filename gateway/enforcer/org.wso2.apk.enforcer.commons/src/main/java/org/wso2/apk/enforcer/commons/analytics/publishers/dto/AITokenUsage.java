/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package org.wso2.apk.enforcer.commons.analytics.publishers.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * AI token usage in analytics event.
 */
public class AITokenUsage {
    @JsonProperty("totalTokens")
    private Double totalTokens;

    @JsonProperty("promptTokens")
    private Double promptTokens;

    @JsonProperty("completionTokens")
    private Double completionTokens;

    @JsonProperty("hour")
    private Integer hour;

    public Double getTotalTokens() {

        return totalTokens;
    }

    public Integer getHour() {

        return hour;
    }

    public void setHour(Integer hour) {

        this.hour = hour;
    }

    public void setTotalTokens(Double totalTokens) {

        this.totalTokens = totalTokens;
    }

    public Double getPromptTokens() {

        return promptTokens;
    }

    public void setPromptTokens(Double promptTokens) {

        this.promptTokens = promptTokens;
    }

    public Double getCompletionTokens() {

        return completionTokens;
    }

    public void setCompletionTokens(Double completionTokens) {

        this.completionTokens = completionTokens;
    }
}