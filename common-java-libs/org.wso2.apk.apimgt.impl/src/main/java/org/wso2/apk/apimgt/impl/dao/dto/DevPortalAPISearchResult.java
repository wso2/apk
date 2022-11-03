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

package org.wso2.apk.apimgt.impl.dao.dto;

import java.util.ArrayList;
import java.util.List;

/**
 * Search result returned when searching APIs from the persistence layer, to be displayed in the Dev Portal.
 */
public class DevPortalAPISearchResult {
    private int returnedAPIsCount;
    private int totalAPIsCount;
    private List<DevPortalAPIInfo> devPortalAPIInfoList= new ArrayList<>();

    public int getTotalAPIsCount() {
        return totalAPIsCount;
    }

    public void setTotalAPIsCount(int totalAPIsCount) {
        this.totalAPIsCount = totalAPIsCount;
    }

    public List<DevPortalAPIInfo> getDevPortalAPIInfoList() {
        return devPortalAPIInfoList;
    }

    public void setDevPortalAPIInfoList(List<DevPortalAPIInfo> devPortalAPIInfoList) {
        this.devPortalAPIInfoList = devPortalAPIInfoList;
    }

    public int getReturnedAPIsCount() {
        return returnedAPIsCount;
    }

    public void setReturnedAPIsCount(int returnedAPIsCount) {
        this.returnedAPIsCount = returnedAPIsCount;
    }

    @Override
    public String toString() {
        return "DevPortalAPISearchResult [returnedAPIsCount=" + returnedAPIsCount + ", totalAPIsCount=" + totalAPIsCount
                + ", devPortalAPIInfoList=" + devPortalAPIInfoList + "]";
    }
}
