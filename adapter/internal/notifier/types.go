/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package notifier

// DeployedAPIRevision represents Information of deployed API revision data
type DeployedAPIRevision struct {
	APIID      string            `json:"apiId"`
	RevisionID int               `json:"revisionId"`
	EnvInfo    []DeployedEnvInfo `json:"envInfo"`
}

// DeployedEnvInfo represents env Information of deployed API revision
type DeployedEnvInfo struct {
	Name  string `json:"name"`
	VHost string `json:"vhost"`
}

// UnDeployedAPIRevision info
type UnDeployedAPIRevision struct {
	APIUUID      string `json:"apiUUID"`
	RevisionUUID string `json:"revisionUUID"`
	Environment  string `json:"environment"`
}
