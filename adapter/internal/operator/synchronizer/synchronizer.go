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

package synchronizer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// APIEvent holds the data structure used for passing API
// lifecycle events from controller go routine to synchronizer
// go routine.
type APIEvent struct {
	EventType     string
	Event         APIState
	UpdatedEvents []string
}

// SuccessEvent holds the data structure used for aknowledgement of a successful API deployment
type SuccessEvent struct {
	APINamespacedName types.NamespacedName
	State             string
	Events            []string
}

var (
	// TODO: Decide on a buffer size and add to config.
	paritionCh chan APIEvent
)

func init() {
	paritionCh = make(chan APIEvent, 10)
}

// HandleAPILifeCycleEvents handles the API events generated from OperatorDataStore
func HandleAPILifeCycleEvents(ch *chan APIEvent, successChannel *chan SuccessEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for API lifecycle events...")
	for event := range *ch {
		if event.Event.APIDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2628, logging.CRITICAL, "API Event is nil"))
		}
		loggers.LoggerAPKOperator.Infof("%s event received for %v", event.EventType, event.Event.APIDefinition.Name)
		var err error
		switch event.EventType {
		case constants.Delete:
			err = undeployAPIInGateway(event.Event)
		case constants.Create:
			err = deployAPIInGateway(event.Event)
		case constants.Update:
			err = deployAPIInGateway(event.Event)
		}
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.MAJOR, "API deployment failed for %s event : %v", event.EventType, err))
		} else {
			if event.EventType != constants.Delete {
				*successChannel <- SuccessEvent{
					APINamespacedName: utils.NamespacedName(event.Event.APIDefinition),
					State:             event.EventType,
					Events:            event.UpdatedEvents,
				}
			}
			if config.ReadConfigs().PartitionServer.Enabled {
				paritionCh <- event
			}
		}
	}
}

// deleteAPIInGateway undeploys the related API in CREATE and UPDATE events.
func undeployAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.ProdHTTPRoute.HTTPRouteCombined, apiState)
	}
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error2630, logging.MAJOR, "Error undeploying prod httpRoute of API : %v in Organization %v from environments %v."+
			" Hence not checking on deleting the sand httpRoute of the API", string(apiState.APIDefinition.ObjectMeta.UID), apiState.APIDefinition.Spec.Organization,
			getLabelsForAPI(apiState.ProdHTTPRoute.HTTPRouteCombined)))
		return err
	}
	if apiState.SandHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.SandHTTPRoute.HTTPRouteCombined, apiState)
	}
	return err
}

func deleteAPIFromEnv(httpRoute *gwapiv1b1.HTTPRoute, apiState APIState) error {
	labels := getLabelsForAPI(httpRoute)
	org := apiState.APIDefinition.Spec.Organization
	uuid := string(apiState.APIDefinition.ObjectMeta.UID)
	return xds.DeleteAPICREvent(labels, uuid, org)
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) error {
	var err error
	if len(apiState.OldOrganizationID) != 0 {
		xds.RemoveAPIFromOrgAPIMap(string((*apiState.APIDefinition).ObjectMeta.UID), apiState.OldOrganizationID)
	}
	if apiState.ProdHTTPRoute == nil {
		var adapterInternalAPI model.AdapterInternalAPI
		adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
		xds.RemoveAPICacheForEnv(adapterInternalAPI, constants.Production)
	}
	if apiState.SandHTTPRoute == nil {
		var adapterInternalAPI model.AdapterInternalAPI
		adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
		xds.RemoveAPICacheForEnv(adapterInternalAPI, constants.Sandbox)
	}
	if apiState.ProdHTTPRoute != nil {
		_, err = GenerateAdapterInternalAPI(apiState, apiState.ProdHTTPRoute, constants.Production)
	}
	if err != nil {
		return err
	}
	if apiState.SandHTTPRoute != nil {
		_, err = GenerateAdapterInternalAPI(apiState, apiState.SandHTTPRoute, constants.Sandbox)
	}
	return err
}

// GenerateAdapterInternalAPI this will populate a AdapterInternalAPI representation for an HTTPRoute
func GenerateAdapterInternalAPI(apiState APIState, httpRoute *HTTPRouteState, envType string) (*model.AdapterInternalAPI, error) {
	var adapterInternalAPI model.AdapterInternalAPI
	adapterInternalAPI.SetIsDefaultVersion(apiState.APIDefinition.Spec.IsDefaultVersion)
	adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
	adapterInternalAPI.SetAPIDefinitionFile(apiState.APIDefinitionFile)
	adapterInternalAPI.SetAPIDefinitionEndpoint(apiState.APIDefinition.Spec.DefinitionPath)
	adapterInternalAPI.EnvType = envType

	environment := apiState.APIDefinition.Spec.Environment
	if environment == "" {
		conf := config.ReadConfigs()
		environment = conf.Adapter.Environment
	}
	adapterInternalAPI.SetEnvironment(environment)

	resourceParams := model.ResourceParams{
		AuthSchemes:               apiState.Authentications,
		ResourceAuthSchemes:       apiState.ResourceAuthentications,
		BackendMapping:            httpRoute.BackendMapping,
		APIPolicies:               apiState.APIPolicies,
		ResourceAPIPolicies:       apiState.ResourceAPIPolicies,
		ResourceScopes:            httpRoute.Scopes,
		InterceptorServiceMapping: apiState.InterceptorServiceMapping,
		BackendJWTMapping:         apiState.BackendJWTMapping,
		RateLimitPolicies:         apiState.RateLimitPolicies,
		ResourceRateLimitPolicies: apiState.ResourceRateLimitPolicies,
	}
	if err := adapterInternalAPI.SetInfoHTTPRouteCR(httpRoute.HTTPRouteCombined, resourceParams); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2631, logging.MAJOR, "Error setting HttpRoute CR info to adapterInternalAPI. %v", err))
		return nil, err
	}
	if err := adapterInternalAPI.Validate(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2632, logging.MAJOR, "Error validating adapterInternalAPI intermediate representation. %v", err))
		return nil, err
	}
	vHosts := getVhostsForAPI(httpRoute.HTTPRouteCombined)
	labels := getLabelsForAPI(httpRoute.HTTPRouteCombined)
	listeners := getListenersForAPI(httpRoute.HTTPRouteCombined, adapterInternalAPI.UUID)

	err := xds.UpdateAPICache(vHosts, labels, listeners, adapterInternalAPI)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2633, logging.MAJOR, "Error updating the API : %s:%s in vhosts: %s, API_UUID: %v. %v",
			adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), vHosts, adapterInternalAPI.UUID, err))
	}
	return &adapterInternalAPI, nil
}

// getVhostForAPI returns the vHosts related to an API.
func getVhostsForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var vHosts []string
	for _, hostName := range httpRoute.Spec.Hostnames {
		vHosts = append(vHosts, string(hostName))
	}
	fmt.Println("vhosts size: ", len(vHosts))
	return vHosts
}

// getLabelsForAPI returns the labels related to an API.
func getLabelsForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var labels []string
	var err error
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		err = xds.SanitizeGateway(string(parentRef.Name), false)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2653, logging.CRITICAL, "Gateway Label is invalid: %s", string(parentRef.Name)))
		} else {
			labels = append(labels, string(parentRef.Name))
		}
	}
	return labels
}

// getListenersForAPI returns the listeners related to an API.
func getListenersForAPI(httpRoute *gwapiv1b1.HTTPRoute, apiUUID string) []string {
	var listeners []string
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		listeners = append(listeners, string(*parentRef.SectionName))
	}
	return listeners
}

// Runtime client connetion
var partitionClient *http.Client

func init() {
	conf := config.ReadConfigs()

	_, _, truststoreLocation := tlsutils.GetKeyLocations()
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: conf.PartitionServer.DisableSslVerification},
	}
	partitionClient = &http.Client{Transport: transport}
}

// SendEventToPartitionServer sends the API create/update/delete event to the partition server.
func SendEventToPartitionServer() {
	conf := config.ReadConfigs()
	for apiEvent := range paritionCh {
		if !apiEvent.Event.APIDefinition.Spec.SystemAPI {
			apiDefinition := apiEvent.Event.APIDefinition
			loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v, API_UUID: %v", apiDefinition.Spec.APIName, string(apiDefinition.ObjectMeta.UID))
			api := apiEvent.Event
			eventType := apiEvent.EventType
			basePath := api.APIDefinition.Spec.BasePath
			organization := api.APIDefinition.Spec.Organization
			version := api.APIDefinition.Spec.APIVersion
			apiName := api.APIDefinition.Spec.APIName
			apiUUID := string(api.APIDefinition.Name)
			var hostNames []string
			httpRoute := api.ProdHTTPRoute
			if httpRoute == nil {
				httpRoute = api.SandHTTPRoute
			}
			for _, hostName := range httpRoute.HTTPRouteCombined.Spec.Hostnames {
				hostNames = append(hostNames, string(hostName))
			}
			data := PartitionEvent{
				EventType:    eventType,
				BasePath:     basePath,
				Organization: organization,
				APIVersion:   version,
				APIName:      apiName,
				APIUUID:      apiUUID,
				Vhosts:       hostNames,
				Partition:    conf.PartitionServer.PartitionName,
			}
			payload, err := json.Marshal(data)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating Event: %v, API_UUID: %v", err, apiUUID)
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", conf.PartitionServer.Host, conf.PartitionServer.ServiceBasePath, "/api-deployment"), bytes.NewBuffer(payload))
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating api definition request: %v, API_UUID: %v", err, apiUUID)
			}
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			resp, err := partitionClient.Do(req)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error sending API Event: %v, API_UUID: %v", err, apiUUID)
			}
			if resp.StatusCode == http.StatusAccepted {
				loggers.LoggerAPKOperator.Info("API Event Accepted", resp.Status)
			}
		}
	}
}

// PartitionEvent is the event sent to the partition server.
type PartitionEvent struct {
	EventType    string   `json:"eventType"`
	APIName      string   `json:"apiName"`
	APIVersion   string   `json:"apiVersion"`
	BasePath     string   `json:"basePath"`
	Organization string   `json:"organization"`
	Partition    string   `json:"partition"`
	APIUUID      string   `json:"apiId"`
	Vhosts       []string `json:"vhosts"`
}
