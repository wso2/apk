/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.config.api;

import java.util.ArrayList;

/**
 * The model containing API Definition (OpenAPI/Swagger) Validation Information
 */
public class APIDefinitionValidationResponse {
    private boolean isValid = false;
    private String content;
    private String jsonContent;
    private byte[] protoContent;
    private String protocol;
    private Info info;
    private APIDefinition parser;
    private ArrayList<ErrorHandler> errorItems = new ArrayList<>();
    private boolean isInit = false;

    public boolean isValid() {
        return isValid;
    }

    public void setValid(boolean valid) {
        isValid = valid;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getContent() {
        return content;
    }

    public String getJsonContent() {
        return jsonContent;
    }

    public void setJsonContent(String jsonContent) {
        this.jsonContent = jsonContent;
    }

    public void setErrorItems(ArrayList<ErrorHandler> errorItems) {
        this.errorItems = errorItems;
    }

    public ArrayList<ErrorHandler> getErrorItems() {
        return errorItems;
    }

    public void setInfo(Info info) {
        this.info = info;
    }

    public Info getInfo() {
        return info;
    }

    public boolean isInit() {
        return isInit;
    }

    public void setInit(boolean init) {
        isInit = init;
    }

    public APIDefinition getParser() {
        return parser;
    }

    public void setParser(APIDefinition parser) {
        this.parser = parser;
    }

    public String getProtocol() {

        return protocol;
    }

    public void setProtocol(String protocol) {

        this.protocol = protocol;
    }

    public byte[] getProtoContent() {
        return protoContent;
    }

    public void setProtoContent(byte[] protoContent) {
        this.protoContent = protoContent;
    }
}
