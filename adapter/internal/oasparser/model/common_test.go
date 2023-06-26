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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostandBasepathandPort(t *testing.T) {
	type setResourcesTestItem struct {
		input   string
		result  *Endpoint
		message string
	}
	dataItems := []setResourcesTestItem{
		{
			input: "https://petstore.io:8000/api/v2",
			result: &Endpoint{
				Host:     "petstore.io",
				Basepath: "/api/v2",
				Port:     8000,
				URLType:  "https",
				RawURL:   "https://petstore.io:8000/api/v2",
			},
			message: "all the details are provided in the endpoint",
		},
		{
			input: "https://petstore.io:8000/api/v2",
			result: &Endpoint{
				Host:     "petstore.io",
				Basepath: "/api/v2",
				Port:     8000,
				URLType:  "https",
				RawURL:   "https://petstore.io:8000/api/v2",
			},
			message: "when port is not provided", //here should find a way to readi configs in tests
		},
		{
			input: "petstore.io:8000/api/v2",
			result: &Endpoint{
				Host:     "petstore.io",
				Basepath: "/api/v2",
				Port:     8000,
				URLType:  "http",
				RawURL:   "http://petstore.io:8000/api/v2",
			},
			message: "when protocol is not provided",
		},
		{
			input:   "https://{defaultHost}",
			result:  nil,
			message: "when malformed endpoint is provided",
		},
		{
			input: "  https://petstore.io:8001/api/v2 ",
			result: &Endpoint{
				Host:     "petstore.io",
				Basepath: "/api/v2",
				Port:     8001,
				URLType:  "https",
				RawURL:   "https://petstore.io:8001/api/v2",
			},
			message: "When leading and trailing spaces present",
		},
	}
	for _, item := range dataItems {
		resultResources, err := getHTTPEndpoint(item.input)
		assert.Equal(t, item.result, resultResources, item.message)
		if resultResources != nil {
			assert.Nil(t, err, "Error encountered when processing the endpoint")
		} else {
			assert.NotNil(t, err, "Should return an error upon failing to process the endpoint")
		}
	}
}

func TestMalformedUrl(t *testing.T) {

	suspectedRawUrls := []string{
		"http://#de.abc.com:80/api",
		"http://&de.abc.com:80/api",
		"http://!de.abc.com:80/api",
		"tcp://http::8900",
		"http://::80",
		"http::80",
		"-",
		"api.worldbank.org-",
		"-api.worldbank.org",
		"",
	}

	for index := range suspectedRawUrls {
		response, _ := getHTTPEndpoint(suspectedRawUrls[index])
		assert.Nil(t, response)
	}

}
