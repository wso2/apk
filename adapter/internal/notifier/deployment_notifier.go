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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/auth"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/tlsutils"
)

const (
	deployedRevisionEP   string = "internal/data/v1/apis/deployed-revisions"
	unDeployedRevisionEP string = "internal/data/v1/apis/undeployed-revision"
	authBasic            string = "Basic "
	authHeader           string = "Authorization"
	contentTypeHeader    string = "Content-Type"
)

//UpdateDeployedRevisions create the DeployedAPIRevision object
func UpdateDeployedRevisions(apiID string, revisionID int, envs []string, vhost string) *DeployedAPIRevision {
	revisions := &DeployedAPIRevision{
		APIID:      apiID,
		RevisionID: revisionID,
		EnvInfo:    []DeployedEnvInfo{},
	}
	for _, env := range envs {
		info := DeployedEnvInfo{
			Name:  env,
			VHost: vhost,
		}
		revisions.EnvInfo = append(revisions.EnvInfo, info)
	}
	return revisions
}

//SendRevisionUpdateAck sends succeeded revision deployment acknowledgement to the control plane
func SendRevisionUpdateAck(deployedRevisionList []*DeployedAPIRevision) {
	conf, _ := config.ReadConfigs()
	cpConfigs := conf.ControlPlane

	if len(deployedRevisionList) < 1 || !cpConfigs.Enabled || !cpConfigs.SendRevisionUpdate {
		return
	}

	logger.LoggerNotifier.Debugf("Revision deployed message is sending to Control plane")

	revisionEP := cpConfigs.ServiceURL
	if strings.HasSuffix(revisionEP, "/") {
		revisionEP += deployedRevisionEP
	} else {
		revisionEP += "/" + deployedRevisionEP
	}

	jsonValue, _ := json.Marshal(deployedRevisionList)

	// Setting authorization header
	basicAuth := authBasic + auth.GetBasicAuth(cpConfigs.Username, cpConfigs.Password)

	logger.LoggerNotifier.Debugf("Revision deployed message sending to Control plane: %v", string(jsonValue))

	// Adding 3 retries for revision update sending
	retries := 0
	for retries < 3 {
		retries++

		req, _ := http.NewRequest("PATCH", revisionEP, bytes.NewBuffer(jsonValue))
		req.Header.Set(authHeader, basicAuth)
		req.Header.Set(contentTypeHeader, "application/json")
		resp, err := tlsutils.InvokeControlPlane(req, cpConfigs.SkipSSLVerification)

		success := true
		if err != nil {
			logger.LoggerNotifier.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error response from %v for attempt %v : %v", revisionEP, retries, err.Error()),
				Severity:  logging.MAJOR,
				ErrorCode: 2100,
			})
			success = false
		}
		if resp != nil && resp.StatusCode != http.StatusOK {
			logger.LoggerNotifier.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error response status code %v from %v for attempt %v", resp.StatusCode, revisionEP, retries),
				Severity:  logging.MINOR,
				ErrorCode: 2101,
			})
			success = false
		}
		if success {
			logger.LoggerNotifier.Infof("Revision deployed message sent to Control plane for attempt %v", retries)
			break
		}
	}
}

// SendRevisionUndeployAck - send the undeployed revision acknowledgement to control plane
func SendRevisionUndeployAck(apiUUID string, revisionUUID string, environment string) {
	conf, _ := config.ReadConfigs()
	cpConfigs := conf.ControlPlane
	if apiUUID == "" || revisionUUID == "" || environment == "" || !cpConfigs.Enabled || !cpConfigs.SendRevisionUpdate {
		return
	}
	revisionEP := cpConfigs.ServiceURL
	if strings.HasSuffix(revisionEP, "/") {
		revisionEP += unDeployedRevisionEP
	} else {
		revisionEP += "/" + unDeployedRevisionEP
	}

	removedRevision := UnDeployedAPIRevision{
		APIUUID:      apiUUID,
		RevisionUUID: revisionUUID,
		Environment:  environment,
	}

	jsonValue, _ := json.Marshal(removedRevision)
	basicAuth := authBasic + auth.GetBasicAuth(cpConfigs.Username, cpConfigs.Password)
	retries := 0
	for retries < 3 {
		retries++
		req, _ := http.NewRequest("POST", revisionEP, bytes.NewBuffer(jsonValue))
		req.Header.Set(authHeader, basicAuth)
		req.Header.Set(contentTypeHeader, "application/json")
		resp, err := tlsutils.InvokeControlPlane(req, cpConfigs.SkipSSLVerification)

		success := true
		if err != nil {
			logger.LoggerNotifier.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error response from %s for attempt %d : %v", revisionEP, retries, err.Error()),
				Severity:  logging.MAJOR,
				ErrorCode: 2100,
			})
			success = false
		}
		if resp != nil && resp.StatusCode != http.StatusOK {
			logger.LoggerNotifier.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error response status code %v from %s for attempt %d", resp.StatusCode, revisionEP, retries),
				Severity:  logging.MINOR,
				ErrorCode: 2101,
			})
			success = false
		}
		if success {
			logger.LoggerNotifier.Infof("Revision un-deployed message sent to Control plane for attempt %d", retries)
			break
		}
		time.Sleep(2 * time.Second)
	}
}
