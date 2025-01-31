/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

// Subscription defines the desired state of Subscription
type Subscription struct {
	SubStatus     string         `json:"subStatus,omitempty"`
	UUID          string         `json:"uuid,omitempty"`
	Organization  string         `json:"organization,omitempty"`
	RatelimitTier string         `json:"ratelimitTier,omitempty"`
	SubscribedAPI *SubscribedAPI `json:"subscribedApi,omitempty"`
}

// API defines the API associated with the subscription
type API struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// SubscriptionList contains a list of Subscription
type SubscriptionList struct {
	List []Subscription `json:"list"`
}

// SubscribedAPI defines the API associated with the subscription
type SubscribedAPI struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}
