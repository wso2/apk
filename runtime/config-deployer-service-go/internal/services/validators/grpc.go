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

package validators

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
)

type GRPCAPIValidator struct{}

// ValidateGRPCAPIDefinition validates the gRPC API definition content.
func (grpcAPIValidator *GRPCAPIValidator) ValidateGRPCAPIDefinition(inputByteArray []byte) (*dto.APIDefinitionValidationResponse, error) {
	validationResponse := &dto.APIDefinitionValidationResponse{}
	if inputByteArray == nil && len(inputByteArray) == 0 {
		validationResponse.IsValid = false
		return validationResponse, fmt.Errorf("gRPC Proto Definition cannot be empty or null")
	} else {
		validationResponse.IsValid = true
		validationResponse.ProtoContent = inputByteArray
	}
	return validationResponse, nil
}
