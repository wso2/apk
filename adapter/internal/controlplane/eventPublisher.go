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
	apiHashMap   map[string]interface{}
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
	APISubType       string            `json:"apiSubType"`
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
	SandVhost        string            `json:"sandVhost"`
	SecurityScheme   []string          `json:"securityScheme"`
	AuthHeader       string            `json:"authHeader"`
	APIKeyHeader     string            `json:"apiKeyHeader"`
	Operations       []Operation       `json:"operations"`
	AIConfiguration  AIConfiguration   `json:"aiConfiguration"`
	APIHash          string            `json:"-"`
	SandAIRL         *AIRL             `json:"sandAIRL"`
	ProdAIRL         *AIRL             `json:"prodAIRL"`
}

// AIRL holds AI ratelimit related data
type AIRL struct {
	PromptTokenCount     *uint32 `json:"promptTokenCount"`
	CompletionTokenCount *uint32 `json:"CompletionTokenCount"`
	TotalTokenCount      *uint32 `json:"totalTokenCount"`
	TimeUnit             string  `json:"timeUnit"`
	RequestCount         *uint32 `json:"requestCount"`
}

// AIConfiguration holds the AI configuration
type AIConfiguration struct {
	LLMProviderName       string `json:"llmProviderName"`
	LLMProviderAPIVersion string `json:"llmProviderAPIVersion"`
}

// Headers contains the request and response header modifier information
type Headers struct {
	RequestHeaders  HeaderModifier `json:"requestHeaders"`
	ResponseHeaders HeaderModifier `json:"responseHeaders"`
}

// HeaderModifier contains header modifier values
type HeaderModifier struct {
	AddHeaders    []Header `json:"addHeaders"`
	RemoveHeaders []string `json:"removeHeaders"`
}

// Header contains the header information
type Header struct {
	Name  string `json:"headerName"`
	Value string `json:"headerValue,omitempty"`
}

// Operation holds the path, verb, throttling and interceptor policy
type Operation struct {
	Path    string   `json:"path"`
	Verb    string   `json:"verb"`
	Scopes  []string `json:"scopes"`
	Headers Headers  `json:"headers"`
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
		apiHashMap = make(map[string]interface{})
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
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusServiceUnavailable {
				loggers.LoggerAPK.Errorf("Error: Unexpected status code: %d, received message: %s, retrying after %d seconds", resp.StatusCode, string(body), retryInterval)
				// Sleep for some time before retrying
				time.Sleep(time.Second * retryInterval)
				continue
			}
			if event.Event == EventTypeDelete {
				// If its a delete event that got propagated to CP then we do not need to update CR.
				break
			}
			delete(apiHashMap, event.API.APIHash)
			var responseMap map[string]interface{}
			if err := json.Unmarshal([]byte(body), &responseMap); err != nil {
				loggers.LoggerAPK.Errorf("Could not decode response body as json. body: %+v", resp.Body)
				break
			}
			// Assuming the response contains an ID field, you can extract it like this:
			id, ok := responseMap["id"].(string)
			revisionID, revisionOk := responseMap["revisionID"].(string)
			if !ok {
				loggers.LoggerAPK.Errorf("Id field not present in response body. encoded body: %+v", responseMap)
				id = ""
				// break
			}
			if !revisionOk {
				loggers.LoggerAPK.Errorf("Revision field not present in response body. encoded body: %+v", responseMap)
				revisionID = ""
				// break
			}
			loggers.LoggerAPK.Infof("Adding label update to API %s/%s, Lebels: apiUUID: %s", event.CRNamespace, event.CRName, id)
			labelsQueue <- APICRLabelsUpdate{
				Namespace: event.CRNamespace,
				Name:      event.CRName,
				Labels:    map[string]string{"apiUUID": id, "revisionID": revisionID, "apiHash": event.API.APIHash},
			}
			break
		}
	}
}

// AddToEventQueue adds the api event to queue
func AddToEventQueue(data APICPEvent) {
	loggers.LoggerAPK.Debugf("Event added to CP Event queue : %+v", data)
	eventQueue <- data
	apiHashMap[data.API.APIHash] = nil
}

// GetLabelQueue adds the label change to queue
func GetLabelQueue() *chan APICRLabelsUpdate {
	return &labelsQueue
}

// IsAPIHashQueued check whether the API related to the dpHash already in the queue for the update
func IsAPIHashQueued(dpHash string) bool {
	_, found := apiHashMap[dpHash]
	return found
}
