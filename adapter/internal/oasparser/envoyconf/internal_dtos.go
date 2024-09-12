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

package envoyconf

import (
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
)

// routeCreateParams is the DTO used to provide information to the envoy route create function
type routeCreateParams struct {
	organizationID               string
	title                        string
	version                      string
	apiType                      string
	xWSO2BasePath                string
	vHost                        string
	endpointBasePath             string
	resource                     *model.Resource
	clusterName                  string
	authHeader                   string
	requestInterceptor           map[string]model.InterceptEndpoint
	responseInterceptor          map[string]model.InterceptEndpoint
	corsPolicy                   *model.CorsConfig
	passRequestPayloadToEnforcer bool
	isDefaultVersion             bool
	createDefaultPath            bool
	apiLevelRateLimitPolicy      *model.RateLimitPolicy
	apiProperties                []v1alpha3.Property
	environment                  string
	envType                      string
	mirrorClusterNames           map[string][]string
	isAiAPI                      bool
}

// RatelimitCriteria criterias of rate limiting
type ratelimitCriteria struct {
	level                string
	organizationID       string
	basePathForRLService string
	environment          string
	envType              string
}
