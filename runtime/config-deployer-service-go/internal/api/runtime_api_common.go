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

package api

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/parsers"
)

type RuntimeAPICommonUtil struct{}

func (runtimeAPIUtil *RuntimeAPICommonUtil) GetAPIFromDefinition(definition string, apiType string) (*dto.API, error) {
	parser := parsers.GetParser(apiType)
	if parser != nil {
		api, err := parser.GetAPIFromDefinition(definition)
		if err != nil {
			return nil, err
		}
		return api, nil
	}
	return nil, fmt.Errorf("definition parser not found: %s", apiType)
}
