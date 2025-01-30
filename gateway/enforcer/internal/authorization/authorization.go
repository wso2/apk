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

package authorization

import (
	"fmt"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Validate performs the authorization.
func Validate(rch *requestconfig.Holder, subAppDataStore *datastore.SubscriptionApplicationDataStore, cfg *config.Server) *dto.ImmediateResponse {
	if immediateResponse := ValidateScopes(rch, subAppDataStore, cfg); immediateResponse != nil {
		return immediateResponse
	}
	cfg.Logger.Info(fmt.Sprintf("Scope validation successful for the request: %s", rch.MatchedResource.Path))
	if immediateResponse := ValidateSubscription(rch, subAppDataStore, cfg); immediateResponse != nil {
		return immediateResponse
	}
	cfg.Logger.Info(fmt.Sprintf("Subscription validation successful for the request: %s", rch.MatchedResource.Path))

	return nil
}
