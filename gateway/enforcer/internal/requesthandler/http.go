/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package requesthandler

import (

	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// HTTP is a struct that represents an HTTP request handler
type HTTP struct {
}

// GetMatchedResource returns the matched resource
func (h *HTTP) GetMatchedResource(api *requestconfig.API, epa dto.ExternalProcessingEnvoyAttributes) *requestconfig.Resource {
	method := epa.RequestMethod
	pathTemplate := util.NormalizePath(epa.Path)
	for _, resource := range api.Resources {
		if string(resource.Method) == method && resource.Path == pathTemplate {
			return resource
		}
	}
	return nil
}
