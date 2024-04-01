/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package controlplane

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"io/ioutil"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
)

var (
	configOnce   sync.Once
	host         string
	port         uint16
	eventQueue   chan APICPEvent
	labelsQueue  chan APICRLabelsUpdate
	wg           sync.WaitGroup
	apisRestPath string
	skipSSL      bool
)

// EventType represents the type of event.
type EventType string

const (
	// EventTypeCreate signifies a create event.
	EventTypeCreate EventType = "CREATE"
	// EventTypeUpdate signifies an update event.
	EventTypeUpdate EventType = "UPDATE"
	// EventTypeDelete signifies a delete event.
	EventTypeDelete EventType = "DELETE"
	applicationJSON           = "application/json"
	retryInterval             = 5
)

// APICRLabelsUpdate hold the label update required for a specific API CR
type APICRLabelsUpdate struct {
	Namespace string
	Name      string
	Labels    map[string]string
}

// APICPEvent represents data for the control plane API.
type APICPEvent struct {
	Event       EventType `json:"event"`
	API         API       `json:"payload"`
	CRName      string    `json:"-"`
	CRNamespace string    `json:"-"`
}

// API holds the data that needs to be sent to agent
type API struct {
	APIUUID          string            `json:"apiUUID"`
	APIName          string            `json:"apiName"`
	APIVersion       string            `json:"apiVersion"`
	IsDefaultVersion bool              `json:"isDefaultVersion"`
	Definition       string            `json:"definition"`
	APIType          string            `json:"apiType"`
	BasePath         string            `json:"basePath"`
	Organization     string            `json:"organization"`
	SystemAPI        bool              `json:"systemAPI"`
	APIProperties    map[string]string `json:"apiProperties,omitempty"`
	Environment      string            `json:"environment,omitempty"`
	RevisionID       string            `json:"revisionID"`
	SandEndpoint     string            `json:"sandEndpoint"`
	ProdEndpoint     string            `json:"prodEndpoint"`
	EndpointProtocol string            `json:"endpointProtocol"`
	CORSPolicy       *CORSPolicy       `json:"cORSPolicy,omitempty"`
	Vhost            string            `json:"vhost"`
	SecurityScheme   []string          `json:"securityScheme"`
	AuthHeader       string            `json:"authHeader"`
}

// CORSPolicy hold cors configs
type CORSPolicy struct {
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty"`
	AccessControlExposeHeaders    []string `json:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge           *int     `json:"accessControlMaxAge,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty"`
}

// init reads the configuration and starts the worker to send data.
func init() {
	configOnce.Do(func() {
		conf := config.ReadConfigs()
		if !conf.Adapter.ControlPlane.EnableAPIPropagation {
			loggers.LoggerAPK.Info("Adapter control plane is not enabled. Not starting agent worker.")
			return
		}
		host = conf.Adapter.ControlPlane.Host
		port = conf.Adapter.ControlPlane.RestPort
		apisRestPath = fmt.Sprintf("https://%s:%d%s", host, port, conf.Adapter.ControlPlane.APIsRestPath)
		skipSSL = conf.Adapter.ControlPlane.SkipSSLVerification
		eventQueue = make(chan APICPEvent, 1000)
		labelsQueue = make(chan APICRLabelsUpdate, 1000)
		wg.Add(1)
		go sendData()
	})
}

// SendData sends data as a POST request to the control plane host.
func sendData() {
	loggers.LoggerAPK.Infof("A thread assigned to send API events to agent")
	tr := &http.Transport{}
	if !skipSSL {
		_, _, truststoreLocation := tlsutils.GetKeyLocations()
		caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Configuring the http client
	client := &http.Client{
		Transport: tr,
	}
	defer wg.Done()
	for event := range eventQueue {
		loggers.LoggerAPK.Infof("Sending api event to agent. Event: %+v", event)
		jsonData, err := json.Marshal(event)
		if err != nil {
			loggers.LoggerAPK.Errorf("Error marshalling data. Error %+v", err)
			continue
		}
		for {
			resp, err := client.Post(
				apisRestPath,
				applicationJSON,
				bytes.NewBuffer(jsonData),
			)
			if err != nil {
				loggers.LoggerAPK.Errorf("Error sending data. Error: %+v, Retrying after %d seconds", err, retryInterval)
				// Sleep for some time before retrying
				time.Sleep(time.Second * retryInterval)
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusServiceUnavailable {
				body, _ := ioutil.ReadAll(resp.Body)
				loggers.LoggerAPK.Errorf("Error: Unexpected status code: %d, received message: %s, retrying after %d seconds", resp.StatusCode, string(body), retryInterval)
				// Sleep for some time before retrying
				time.Sleep(time.Second * retryInterval)
				continue
			}
			if event.Event == EventTypeDelete {
				// If its a delete event that got propagated to CP then we do not need to update CR.
				break
			}
			var responseMap map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseMap)
			if err != nil {
				loggers.LoggerAPK.Errorf("Could not decode response body as json. body: %+v", resp.Body)
				break
			}
			// Assuming the response contains an ID field, you can extract it like this:
			id, ok := responseMap["id"].(string)
			revisionID, revisionOk := responseMap["revisionID"].(string)
			if !ok || !revisionOk {
				loggers.LoggerAPK.Errorf("Id or/both revision Id field not present in response body. encoded body: %+v", responseMap)
				id = ""
				revisionID = ""
				// break
			}
			loggers.LoggerAPK.Infof("Adding label update to API %s/%s, Lebels: apiUUID: %s", event.CRNamespace, event.CRName, id)
			labelsQueue <- APICRLabelsUpdate{
				Namespace: event.CRNamespace,
				Name:      event.CRName,
				Labels:    map[string]string{"apiUUID": id, "revisionID": revisionID},
			}
			break
		}
	}
}

// AddToEventQueue adds the api event to queue
func AddToEventQueue(data APICPEvent) {
	loggers.LoggerAPK.Debugf("Event added to CP Event queue : %+v", data)
	eventQueue <- data
}

// GetLabelQueue adds the label change to queue
func GetLabelQueue() *chan APICRLabelsUpdate {
	return &labelsQueue
}
