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

package org.wso2.apk.common;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.fasterxml.jackson.dataformat.yaml.YAMLMapper;

/**
 * Utility class to convert to yaml.
 */
public class YamlUtil {

    /**
     * This method used to convert jsonString to yaml string.
     *
     * @param jsonString string representation of Json.
     * @return string representation of yaml.
     * @throws Exception
     */
    public String fromJsonStringToYaml(String jsonString) throws Exception {
        try {

            JsonNode jsonNodeTree = new ObjectMapper().readTree(jsonString);
            return new YAMLMapper().writeValueAsString(jsonNodeTree);
        } catch (JsonProcessingException e) {
            throw new Exception("Error occurred while converting json to yaml" + e, e);
        }
    }

    public String fromYamlStringToJson(String yamlString) throws Exception {
        try {
            ObjectMapper yamlMapper = new ObjectMapper(new YAMLFactory());
            Object yamlObject = yamlMapper.readValue(yamlString, Object.class);
            return new ObjectMapper().writerWithDefaultPrettyPrinter().writeValueAsString(yamlObject);
        } catch (JsonProcessingException e) {
            throw new Exception("Error occured converting yaml to json", e);
        }

    }
}
