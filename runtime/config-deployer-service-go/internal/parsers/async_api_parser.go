/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package parsers

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
)

type AsyncApiParser struct{}

// GetAPIFromDefinition is a method that should be implemented by AsyncApiParser to parse an AsyncAPI definition
func (asyncParser *AsyncApiParser) GetAPIFromDefinition(definition string) (*dto.API, error) {
	return nil, fmt.Errorf("unimplemented method 'getAPIFromDefinition' in AsyncApiParser")
}
