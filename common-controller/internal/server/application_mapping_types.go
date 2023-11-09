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

package server

// ApplicationMapping defines the desired state of ApplicationMapping
type ApplicationMapping struct {
	UUID            string `json:"uuid"`
	ApplicationRef  string `json:"applicationRef"`
	SubscriptionRef string `json:"subscriptionRef"`
}

// ApplicationMappingList contains a list of ApplicationMapping
type ApplicationMappingList struct {
	List []ApplicationMapping `json:"list"`
}
