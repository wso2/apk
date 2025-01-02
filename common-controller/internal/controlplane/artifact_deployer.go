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

import (
	"github.com/wso2/apk/common-go-libs/pkg/server/model"
)


// ArtifactDeployer is an interface that defines the methods that should be implemented by an artifact deployer
type ArtifactDeployer interface {
	DeployApplication(application model.Application) error
	UpdateApplication(application model.Application) error
	DeploySubscription(subscription model.Subscription) error
	UpdateSubscription(subscription model.Subscription) error
	DeployApplicationMappings(applicationMapping model.ApplicationMapping) error
	UpdateApplicationMappings(applicationMapping model.ApplicationMapping) error
	DeployKeyMappings(keyMapping model.ApplicationKeyMapping) error
	UpdateKeyMappings(keyMapping model.ApplicationKeyMapping) error
	GetApplication(applicationID string) (model.Application, error)
	GetSubscription(subscriptionID string) (model.Subscription, error)
	GetApplicationMappings(applicationID string) (model.ApplicationMapping, error)
	GetKeyMappings(applicationID string) (model.ApplicationKeyMapping, error)
	GetAllApplications() (model.ApplicationList, error)
	GetAllSubscriptions() (model.SubscriptionList, error)
	GetAllApplicationMappings() (model.ApplicationMappingList, error)
	GetAllKeyMappings() (model.ApplicationKeyMappingList, error)
	DeleteApplication(applicationID string) error
	DeleteSubscription(subscriptionID string) error
	DeleteApplicationMappings(applicationID string) error
	DeleteKeyMappings(keyMapping model.ApplicationKeyMapping) error
	DeployAllApplications(applications model.ApplicationList) error
	DeployAllSubscriptions(subscriptions model.SubscriptionList) error
	DeployAllApplicationMappings(applicationMappings model.ApplicationMappingList) error
	DeployAllKeyMappings(keyMappings model.ApplicationKeyMappingList) error
}
