/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package controlplane

import "github.com/wso2/apk/common-controller/internal/server"

// ArtifactDeployer is an interface that defines the methods that should be implemented by an artifact deployer
type ArtifactDeployer interface {
	DeployApplication(application server.Application) error
	UpdateApplication(application server.Application) error
	DeploySubscription(subscription server.Subscription) error
	UpdateSubscription(subscription server.Subscription) error
	DeployApplicationMappings(applicationMapping server.ApplicationMapping) error
	UpdateApplicationMappings(applicationMapping server.ApplicationMapping) error
	DeployKeyMappings(keyMapping server.ApplicationKeyMapping) error
	UpdateKeyMappings(keyMapping server.ApplicationKeyMapping) error
	GetApplication(applicationID string) (server.Application, error)
	GetSubscription(subscriptionID string) (server.Subscription, error)
	GetApplicationMappings(applicationID string) (server.ApplicationMapping, error)
	GetKeyMappings(applicationID string) (server.ApplicationKeyMapping, error)
	GetAllApplications() (server.ApplicationList, error)
	GetAllSubscriptions() (server.SubscriptionList, error)
	GetAllApplicationMappings() (server.ApplicationMappingList, error)
	GetAllKeyMappings() (server.ApplicationKeyMappingList, error)
	DeleteApplication(applicationID string) error
	DeleteSubscription(subscriptionID string) error
	DeleteApplicationMappings(applicationID string) error
	DeleteKeyMappings(keyMapping server.ApplicationKeyMapping) error
	DeployAllApplications(applications server.ApplicationList) error
	DeployAllSubscriptions(subscriptions server.SubscriptionList) error
	DeployAllApplicationMappings(applicationMappings server.ApplicationMappingList) error
	DeployAllKeyMappings(keyMappings server.ApplicationKeyMappingList) error
}
