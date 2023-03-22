// /*
//  *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
//  *
//  *  Licensed under the Apache License, Version 2.0 (the "License");
//  *  you may not use this file except in compliance with the License.
//  *  You may obtain a copy of the License at
//  *
//  *  http://www.apache.org/licenses/LICENSE-2.0
//  *
//  *  Unless required by applicable law or agreed to in writing, software
//  *  distributed under the License is distributed on an "AS IS" BASIS,
//  *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  *  See the License for the specific language governing permissions and
//  *  limitations under the License.
//  *
//  */

package model

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/wso2/apk/adapter/config"
// 	"github.com/wso2/apk/adapter/internal/oasparser/utils"
// )

// // TestSetInfoAsyncAPI for adapterInternalAPI.SetInfoAsyncAPI(asyncapi)
// func TestSetInfoAsyncAPI(t *testing.T) {

// 	type setInfoAsyncAPITestItem struct {
// 		actual   AdapterInternalAPI
// 		expected AdapterInternalAPI
// 	}

// 	asyncapiFilePath := config.GetApkHome() + "/test-resources/envoycodegen/asyncapi_websocket.yaml"
// 	asyncapiByteArr, err := ioutil.ReadFile(asyncapiFilePath)
// 	assert.Nil(t, err, "Error while reading file : %v"+asyncapiFilePath)
// 	apiJsn, conversionErr := utils.ToJSON(asyncapiByteArr)
// 	assert.Nil(t, conversionErr, "YAML to JSON conversion error : %v"+asyncapiFilePath)

// 	var asyncapi AsyncAPI
// 	err = json.Unmarshal(apiJsn, &asyncapi)
// 	assert.Nil(t, err, "Error occurred while parsing api.yaml")
// 	var adapterInternalAPI AdapterInternalAPI
// 	err = adapterInternalAPI.SetInfoAsyncAPI(asyncapi)
// 	assert.Nil(t, err, "Error while populating the AdapterInternalAPI object for websocket APIs")

// 	dataItem := setInfoAsyncAPITestItem{
// 		actual: adapterInternalAPI,
// 		expected: AdapterInternalAPI{
// 			title:   "WebSocket",
// 			version: "1",
// 			productionEndpoints: &EndpointCluster{
// 				EndpointPrefix: "clusterProd",
// 				Endpoints: []Endpoint{
// 					{
// 						Host:     "ws.ifelse.io",
// 						Port:     443,
// 						URLType:  "wss",
// 						Basepath: "",
// 						RawURL:   "wss://ws.ifelse.io:443",
// 					},
// 				},
// 				EndpointType: "load_balance",
// 			},
// 			resources: []*Resource{
// 				{
// 					path: "/notifications",
// 					methods: []*Operation{
// 						{
// 							method: "GET",
// 							security: []map[string][]string{
// 								{
// 									"oauth2": {"abc"},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				{
// 					path: "/rooms/{roomID}",
// 					methods: []*Operation{
// 						{
// 							method: "GET",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	assert.Nil(t, err, "Error while populating the adapterInternalAPI object for asyncAPIs")

// 	assert.Equal(t, dataItem.expected.productionEndpoints.Endpoints[0],
// 		dataItem.actual.productionEndpoints.Endpoints[0], "AsyncAPI AdapterInternalAPI productionEndpoints mismatch")
// 	assert.Nil(t, dataItem.actual.sandboxEndpoints, "AsyncAPI AdapterInternalAPI sandboxEndpoints not nil")

// 	assert.Equal(t, dataItem.expected.resources[0].path,
// 		dataItem.actual.resources[0].path, "AsyncAPI AdapterInternalAPI /notifications path mismatch")
// 	assert.Equal(t, dataItem.expected.resources[1].path,
// 		dataItem.actual.resources[1].path, "AsyncAPI AdapterInternalAPI /rooms/{roomID} path mismatch")

// 	assert.Equal(t, dataItem.expected.resources[0].methods[0].method,
// 		dataItem.actual.resources[0].methods[0].method,
// 		"AsyncAPI AdapterInternalAPI resource method mismatch")

// 	assert.Equal(t, dataItem.expected.resources[0].methods[0].security[0]["oauth2"],
// 		dataItem.actual.resources[0].methods[0].security[0]["oauth2"],
// 		"AsyncAPI AdapterInternalAPI publish security scope mismatch")

// 	assert.Equal(t, len(dataItem.expected.resources[0].methods), 1,
// 		"AsyncAPI AdapterInternalAPI resource has more that one method")
// }
