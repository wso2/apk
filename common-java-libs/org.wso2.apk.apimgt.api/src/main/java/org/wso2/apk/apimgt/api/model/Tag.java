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
package org.wso2.apk.apimgt.api.model;

/**
 * Represents a tag(s) used to categorize an API.
 *
 * @see API
 */
@SuppressWarnings("unused")
public class Tag {

    private String name;
    private String description;
    private boolean thumbnailExists;
    private int noOfOccurrences;
    private String thumbnailUrl;

    public Tag(String name) {
        this.name = name;
    }

    public Tag(String name, int noOfOccurrences) {
        this.name = name;
        this.noOfOccurrences = noOfOccurrences;
    }

    public String getName() {
        return name;
    }

    public int getNoOfOccurrences() {
        return noOfOccurrences;
    }

    public void setNoOfOccurrences(int noOfOccurrences) {
        this.noOfOccurrences = noOfOccurrences;
    }

    public String getDescription() {

        if(description != null){
            return description;
        }
        return ""; // Return empty string if there is no description.
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public boolean isThumbnailExists() {
        return thumbnailExists;
    }

    public void setThumbnailExists(boolean thumbnailExists) {
        this.thumbnailExists = thumbnailExists;
    }

    public String getThumbnailUrl() {
        return thumbnailUrl;
    }

    public void setThumbnailUrl(String thumbnailUrl) {
        this.thumbnailUrl = thumbnailUrl;
    }

    @Override public boolean equals(Object o) {
        if (this == o)
            return true;
        if (o == null || getClass() != o.getClass())
            return false;

        Tag tag = (Tag) o;

        return name.equals(tag.name);

    }

    @Override public int hashCode() {
        return name.hashCode();
    }
}
