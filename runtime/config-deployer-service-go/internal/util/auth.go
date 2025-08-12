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

package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
)

// UserContext represents authenticated user information
type UserContext struct {
	Username     string                 `json:"username"`
	UserID       *string                `json:"userId,omitempty"`
	Organization *dto.Organization      `json:"organization"`
	Claims       map[string]interface{} `json:"claims"`
}

// GetAuthenticatedUserContext extracts user context from request context
func GetAuthenticatedUserContext(cxt *gin.Context) (*UserContext, error) {
	userContextAttribute := cxt.Value(constants.ValidatedUserContext)
	if userContextAttribute == nil {
		return nil, fmt.Errorf("unauthorized Request")
	}
	userContext, ok := userContextAttribute.(*UserContext)
	if !ok {
		return nil, fmt.Errorf("unauthorized Request")
	}
	return userContext, nil
}
