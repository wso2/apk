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

package org.wso2.apk.apimgt.impl.dto;

import com.google.gson.annotations.SerializedName;

import java.io.Serializable;

public class TokenHandlingDto implements Serializable {

    @SerializedName("enable")
    private Boolean enable = true;

    public enum TypeEnum {

        @SerializedName("REFERENCE") REFERENCE(String.valueOf("REFERENCE")), @SerializedName("JWT") JWT(
                String.valueOf("JWT")), @SerializedName("CUSTOM") CUSTOM(String.valueOf("CUSTOM"));

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

        public static TypeEnum fromValue(String v) {

            for (TypeEnum b : TypeEnum.values()) {
                if (String.valueOf(b.value).equals(v)) {
                    return b;
                }
            }
            return null;
        }
    }

    @SerializedName("type")
    private TypeEnum type = null;
    @SerializedName("value")
    private Object value = null;

    public Boolean getEnable() {

        return enable;
    }

    public void setEnable(Boolean enable) {

        this.enable = enable;
    }

    public TypeEnum getType() {

        return type;
    }

    public void setType(TypeEnum type) {

        this.type = type;
    }

    public Object getValue() {

        return value;
    }

    public void setValue(Object value) {

        this.value = value;
    }
}
