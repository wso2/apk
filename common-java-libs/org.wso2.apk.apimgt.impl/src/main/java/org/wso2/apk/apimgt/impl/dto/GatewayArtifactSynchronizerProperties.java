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

import org.wso2.apk.apimgt.impl.APIConstants;

import java.util.HashSet;
import java.util.Set;

public class GatewayArtifactSynchronizerProperties {

    private boolean saveArtifactsEnabled = false;
    private boolean retrieveFromStorageEnabled = false;
    private String saverName = APIConstants.GatewayArtifactSynchronizer.DB_SAVER_NAME;
    private String retrieverName = APIConstants.GatewayArtifactSynchronizer.DB_RETRIEVER_NAME;
    private Set<String> gatewayLabels = new HashSet<>();
    private String artifactSynchronizerDataSource = "jdbc/WSO2AM_DB";
    private long retryDuartion = 15000 ;
    private String gatewayStartup = "sync";
    private long eventWaitingTime = 1;


    public String getSaverName() {

        return saverName;
    }

    public long getEventWaitingTime() {

        return eventWaitingTime;
    }

    public void setEventWaitingTime(long eventWaitingTime) {

        this.eventWaitingTime = eventWaitingTime;
    }


    public void setSaverName(String saverName) {

        this.saverName = saverName;
    }

    public String getRetrieverName() {

        return retrieverName;
    }

    public String getArtifactSynchronizerDataSource() {

        return artifactSynchronizerDataSource;
    }

    public void setArtifactSynchronizerDataSource(String artifactSynchronizerDataSource) {

        this.artifactSynchronizerDataSource = artifactSynchronizerDataSource;
    }

    public void setRetrieverName(String retrieverName) {

        this.retrieverName = retrieverName;
    }

    public Set<String> getGatewayLabels() {

        return gatewayLabels;
    }

    public void setGatewayLabels(Set<String> gatewayLabels) {

        this.gatewayLabels = gatewayLabels;
    }

    public void setPublishDirectlyToGatewayEnabled(boolean publishDirectlyToGatewayEnabled) {

    }

    public boolean isRetrieveFromStorageEnabled() {

        return retrieveFromStorageEnabled;
    }

    public void setRetrieveFromStorageEnabled(boolean retrieveFromStorageEnabled) {

        this.retrieveFromStorageEnabled = retrieveFromStorageEnabled;
    }

    public boolean isSaveArtifactsEnabled() {

        return saveArtifactsEnabled;
    }

    public void setSaveArtifactsEnabled(boolean saveArtifactsEnabled) {

        this.saveArtifactsEnabled = saveArtifactsEnabled;
    }

    public long getRetryDuartion() {

        return retryDuartion;
    }

    public void  setRetryDuartion(long retryDuartion) {

        this.retryDuartion = retryDuartion;
    }

    public String getGatewayStartup() {

        return gatewayStartup;
    }

    public void  setGatewayStartup(String gatewayStartup) {

        this.gatewayStartup = gatewayStartup;
    }
}
