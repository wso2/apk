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
	"context"
	"strconv"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/common-go-libs/utils"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	// CreationTimeStamp constant for annotation creationTimeStamp
	CreationTimeStamp = "creationTimeStamp"
)

// K8sArtifactDeployer is a struct that implements ArtifactDeployer interface
type K8sArtifactDeployer struct {
	client client.Client
}

// NewK8sArtifactDeployer creates a new K8sArtifactDeployer
func NewK8sArtifactDeployer(mgr manager.Manager) K8sArtifactDeployer {
	return K8sArtifactDeployer{client: mgr.GetClient()}
}

// DeployApplication deploys an application
func (k8sArtifactDeployer K8sArtifactDeployer) DeployApplication(application server.Application) error {
	crApplication := cpv1alpha2.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      application.UUID,
			Namespace: utils.GetOperatorPodNamespace(),
			Labels:    map[string]string{CreationTimeStamp: strconv.FormatInt(application.TimeStamp, 10)},
		},
		Spec: cpv1alpha2.ApplicationSpec{
			Name:         application.Name,
			Owner:        application.Owner,
			Organization: application.OrganizationID,
			Attributes:   application.Attributes,
		},
	}
	loggers.LoggerAPKOperator.Debugf("Creating Application %s", application.UUID)
	loggers.LoggerAPKOperator.Debugf("Application CR %v ", crApplication)
	err := k8sArtifactDeployer.client.Create(context.Background(), &crApplication)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to create application in k8s %v", err.Error()))
		return err
	}
	return nil
}

// UpdateApplication updates an application
func (k8sArtifactDeployer K8sArtifactDeployer) UpdateApplication(application server.Application) error {
	crApplication := cpv1alpha2.Application{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: application.UUID, Namespace: utils.GetOperatorPodNamespace()}, &crApplication)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
			return err
		}
		k8sArtifactDeployer.DeployApplication(application)
	} else {
		crApplication.Spec.Name = application.Name
		crApplication.Spec.Owner = application.Owner
		crApplication.Spec.Organization = application.OrganizationID
		crApplication.Spec.Attributes = application.Attributes
		err := k8sArtifactDeployer.client.Update(context.Background(), &crApplication)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to update application in k8s %v", err.Error()))
			return err
		}
	}

	return nil
}

// UpdateKeyMappings updates a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) UpdateKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	return nil
}

// DeploySubscription deploys a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) DeploySubscription(subscription server.Subscription) error {
	crSubscription := cpv1alpha2.Subscription{ObjectMeta: v1.ObjectMeta{Name: subscription.UUID, Namespace: utils.GetOperatorPodNamespace()},
		Spec: cpv1alpha2.SubscriptionSpec{Organization: subscription.Organization, API: cpv1alpha2.API{Name: subscription.SubscribedAPI.Name, Version: subscription.SubscribedAPI.Version}, SubscriptionStatus: subscription.SubStatus}}
	err := k8sArtifactDeployer.client.Create(context.Background(), &crSubscription)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1101, logging.CRITICAL, "Failed to create subscription in k8s %v", err.Error()))
		return err
	}
	return nil
}

// UpdateSubscription updates a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) UpdateSubscription(subscription server.Subscription) error {
	crSubscription := cpv1alpha2.Subscription{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: subscription.UUID, Namespace: utils.GetOperatorPodNamespace()}, &crSubscription)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get subscription from k8s %v", err.Error()))
			return err
		}
		k8sArtifactDeployer.DeploySubscription(subscription)
	} else {
		crSubscription.Spec.Organization = subscription.Organization
		crSubscription.Spec.API.Name = subscription.SubscribedAPI.Name
		crSubscription.Spec.API.Version = subscription.SubscribedAPI.Version
		crSubscription.Spec.SubscriptionStatus = subscription.SubStatus
		err := k8sArtifactDeployer.client.Update(context.Background(), &crSubscription)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to update subscription in k8s %v", err.Error()))
			return err
		}
	}
	return nil
}

// DeployApplicationMappings deploys an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeployApplicationMappings(applicationMapping server.ApplicationMapping) error {
	crApplicationMapping := cpv1alpha2.ApplicationMapping{ObjectMeta: v1.ObjectMeta{Name: applicationMapping.UUID, Namespace: utils.GetOperatorPodNamespace()},
		Spec: cpv1alpha2.ApplicationMappingSpec{ApplicationRef: applicationMapping.ApplicationRef, SubscriptionRef: applicationMapping.SubscriptionRef}}
	return k8sArtifactDeployer.client.Create(context.Background(), &crApplicationMapping)
}

// DeployKeyMappings deploys a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeployKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	var crApplication cpv1alpha2.Application
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: keyMapping.ApplicationUUID, Namespace: utils.GetOperatorPodNamespace()}, &crApplication)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return err
	}
	securitySchemes := cpv1alpha2.SecuritySchemes{}
	if crApplication.Spec.SecuritySchemes != nil {
		securitySchemes = *crApplication.Spec.SecuritySchemes
	}
	if keyMapping.SecurityScheme == constants.OAuth2 {
		if securitySchemes.OAuth2 == nil {
			securitySchemes.OAuth2 = &cpv1alpha2.SecurityScheme{Environments: []cpv1alpha2.Environment{generateSecurityScheme(keyMapping)}}
		} else {
			environments := make([]cpv1alpha2.Environment, 0)
			for _, environment := range securitySchemes.OAuth2.Environments {
				if environment.EnvID != keyMapping.EnvID || environment.AppID != keyMapping.ApplicationIdentifier || environment.KeyType != keyMapping.KeyType {
					environments = append(environments, environment)
				}
			}
			securitySchemes.OAuth2.Environments = append(environments, generateSecurityScheme(keyMapping))
		}
	}
	crApplication.Spec.SecuritySchemes = &securitySchemes
	loggers.LoggerAPKOperator.Infof("Updating Application %v", crApplication)
	return k8sArtifactDeployer.client.Update(context.Background(), &crApplication)
}

// DeleteAllApplicationMappings deletes all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllApplicationMappings() error {
	return nil
}

// DeleteAllKeyMappings deletes all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllKeyMappings() error {
	return nil
}

// DeleteAllSubscriptions deletes all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteAllSubscriptions() error {
	return nil
}

// DeleteApplication deletes an application
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteApplication(applicationID string) error {
	crApplication := cpv1alpha2.Application{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: applicationID, Namespace: utils.GetOperatorPodNamespace()}, &crApplication)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
			return err
		}
	} else {
		err := k8sArtifactDeployer.client.Delete(context.Background(), &crApplication)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to delete application in k8s %v", err.Error()))
			return err
		}
	}

	return nil
}

// DeleteApplicationMappings deletes an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteApplicationMappings(applicationMapping string) error {
	crApplicationMapping := cpv1alpha2.ApplicationMapping{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: applicationMapping, Namespace: utils.GetOperatorPodNamespace()}, &crApplicationMapping)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application mapping from k8s %v", err.Error()))
			return err
		}
	} else {
		err := k8sArtifactDeployer.client.Delete(context.Background(), &crApplicationMapping)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to delete application mapping in k8s %v", err.Error()))
			return err
		}
	}

	return nil
}

// UpdateApplicationMappings updates an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) UpdateApplicationMappings(applicationMapping server.ApplicationMapping) error {
	crApplicationMapping := cpv1alpha2.ApplicationMapping{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: applicationMapping.UUID, Namespace: utils.GetOperatorPodNamespace()}, &crApplicationMapping)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application mapping from k8s %v", err.Error()))
			return err
		}
		k8sArtifactDeployer.DeployApplicationMappings(applicationMapping)
	} else {
		crApplicationMapping.Spec.ApplicationRef = applicationMapping.ApplicationRef
		crApplicationMapping.Spec.SubscriptionRef = applicationMapping.SubscriptionRef
		err := k8sArtifactDeployer.client.Update(context.Background(), &crApplicationMapping)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to update application mapping in k8s %v", err.Error()))
			return err
		}
	}

	return nil
}

// DeleteKeyMappings deletes a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	var crApplication cpv1alpha2.Application
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: keyMapping.ApplicationUUID, Namespace: utils.GetOperatorPodNamespace()}, &crApplication)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return err
	}
	if crApplication.Spec.SecuritySchemes != nil {
		securitySchemes := *crApplication.Spec.SecuritySchemes
		if keyMapping.SecurityScheme == constants.OAuth2 && securitySchemes.OAuth2 != nil {
			if securitySchemes.OAuth2.Environments != nil && len(securitySchemes.OAuth2.Environments) > 0 {
				environments := make([]cpv1alpha2.Environment, 0)
				for _, environment := range securitySchemes.OAuth2.Environments {
					if environment.EnvID != keyMapping.EnvID || environment.AppID != keyMapping.ApplicationIdentifier {
						environments = append(environments, environment)
					}
				}
				securitySchemes.OAuth2.Environments = environments
			}
		}
		crApplication.Spec.SecuritySchemes = &securitySchemes
		loggers.LoggerAPKOperator.Infof("Updating Application %v", crApplication)
		return k8sArtifactDeployer.client.Update(context.Background(), &crApplication)
	}
	return nil
}

// DeleteSubscription deletes a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) DeleteSubscription(subscriptionID string) error {
	crSubscription := cpv1alpha2.Subscription{}
	err := k8sArtifactDeployer.client.Get(context.Background(), client.ObjectKey{Name: subscriptionID, Namespace: utils.GetOperatorPodNamespace()}, &crSubscription)
	if err != nil {
		if !k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get subscription from k8s %v", err.Error()))
			return err
		}
	} else {
		err := k8sArtifactDeployer.client.Delete(context.Background(), &crSubscription)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.CRITICAL, "Failed to delete subscription in k8s %v", err.Error()))
			return err
		}
	}

	return nil
}

// DeployAllApplicationMappings deploys all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllApplicationMappings(applicationMappings server.ApplicationMappingList) error {
	applicationMappingsFromK8s, _, err := k8sArtifactDeployer.retrieveAllApplicationMappings("")
	if err != nil {
		return err
	}
	clonedApplicationMappingsFromK8s := make([]cpv1alpha2.ApplicationMapping, len(applicationMappingsFromK8s))
	copy(clonedApplicationMappingsFromK8s, applicationMappingsFromK8s)
	clonedApplicationMappings := make([]server.ApplicationMapping, len(applicationMappings.List))
	copy(clonedApplicationMappings, applicationMappings.List)
	newApplicationMappings := make([]server.ApplicationMapping, 0)
	sameApplicationMappings := make([]server.ApplicationMapping, 0)
	for _, applicationMapping := range clonedApplicationMappings {
		found := false
		unFilteredApplicationMappingsInK8s := make([]cpv1alpha2.ApplicationMapping, 0)
		for _, applicationMappingFromK8s := range clonedApplicationMappingsFromK8s {
			if applicationMapping.ApplicationRef == applicationMappingFromK8s.Spec.ApplicationRef && applicationMapping.SubscriptionRef == applicationMappingFromK8s.Spec.SubscriptionRef {
				sameApplicationMappings = append(sameApplicationMappings, applicationMapping)
				found = true
				break
			}
			unFilteredApplicationMappingsInK8s = append(unFilteredApplicationMappingsInK8s, applicationMappingFromK8s)
		}
		clonedApplicationMappingsFromK8s = unFilteredApplicationMappingsInK8s
		if !found {
			newApplicationMappings = append(newApplicationMappings, applicationMapping)
		}
	}
	for _, applicationMapping := range newApplicationMappings {
		err := k8sArtifactDeployer.DeployApplicationMappings(applicationMapping)
		if err != nil {
			return err
		}
	}
	for _, applicationMapping := range sameApplicationMappings {
		err := k8sArtifactDeployer.UpdateApplicationMappings(applicationMapping)
		if err != nil {
			return err
		}
	}
	for _, applicationMappingFromK8s := range clonedApplicationMappingsFromK8s {
		err := k8sArtifactDeployer.DeleteApplicationMappings(applicationMappingFromK8s.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeployAllApplications deploys all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllApplications(applications server.ApplicationList) error {
	applicationsFromK8s, _, err := k8sArtifactDeployer.retrieveAllApplicationsFromK8s("")
	if err != nil {
		return err
	}
	clonedApplicationsFromK8s := make([]cpv1alpha2.Application, len(applicationsFromK8s))
	copy(clonedApplicationsFromK8s, applicationsFromK8s)
	clonedApplications := make([]server.Application, len(applications.List))
	copy(clonedApplications, applications.List)
	newApplications := make([]server.Application, 0)
	sameApplications := make([]server.Application, 0)
	for _, application := range clonedApplications {
		found := false
		unFilteredApplicationsInK8s := make([]cpv1alpha2.Application, 0)
		for _, applicationFromK8s := range clonedApplicationsFromK8s {
			if application.UUID == applicationFromK8s.Name {
				sameApplications = append(sameApplications, application)
				found = true
				break
			}
			unFilteredApplicationsInK8s = append(unFilteredApplicationsInK8s, applicationFromK8s)
		}
		clonedApplicationsFromK8s = unFilteredApplicationsInK8s
		if !found {
			newApplications = append(newApplications, application)
		}
	}
	for _, application := range newApplications {
		err := k8sArtifactDeployer.DeployApplication(application)
		if err != nil {
			return err
		}
	}
	for _, application := range sameApplications {
		err := k8sArtifactDeployer.UpdateApplication(application)
		if err != nil {
			return err
		}
	}
	for _, applicationFromK8s := range clonedApplicationsFromK8s {
		err := k8sArtifactDeployer.DeleteApplication(applicationFromK8s.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeployAllKeyMappings deploys all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllKeyMappings(keyMappings server.ApplicationKeyMappingList) error {
	for _, keyMapping := range keyMappings.List {
		err := k8sArtifactDeployer.DeployKeyMappings(keyMapping)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeployAllSubscriptions deploys all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) DeployAllSubscriptions(subscriptions server.SubscriptionList) error {
	subscriptionsFromK8s, _, err := k8sArtifactDeployer.retrieveAllSubscriptionsFromK8s("")
	if err != nil {
		return err
	}
	clonedSubscriptionsFromK8s := make([]cpv1alpha2.Subscription, len(subscriptionsFromK8s))
	copy(clonedSubscriptionsFromK8s, subscriptionsFromK8s)
	clonedSubscriptions := make([]server.Subscription, len(subscriptions.List))
	copy(clonedSubscriptions, subscriptions.List)
	newSubscriptions := make([]server.Subscription, 0)
	sameSubscriptions := make([]server.Subscription, 0)
	for _, subscription := range clonedSubscriptions {
		found := false
		unFilteredSubscriptionsInK8s := make([]cpv1alpha2.Subscription, 0)
		for _, subscriptionFromK8s := range clonedSubscriptionsFromK8s {
			if subscription.UUID == subscriptionFromK8s.Name {
				sameSubscriptions = append(sameSubscriptions, subscription)
				found = true
				break
			}
			unFilteredSubscriptionsInK8s = append(unFilteredSubscriptionsInK8s, subscriptionFromK8s)
		}
		clonedSubscriptionsFromK8s = unFilteredSubscriptionsInK8s
		if !found {
			newSubscriptions = append(newSubscriptions, subscription)
		}
	}
	for _, subscription := range newSubscriptions {
		err := k8sArtifactDeployer.DeploySubscription(subscription)
		if err != nil {
			return err
		}
	}
	for _, subscription := range sameSubscriptions {
		err := k8sArtifactDeployer.UpdateSubscription(subscription)
		if err != nil {
			return err
		}
	}
	for _, subscriptionFromK8s := range clonedSubscriptionsFromK8s {
		err := k8sArtifactDeployer.DeleteSubscription(subscriptionFromK8s.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllApplicationMappings returns all application mappings
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllApplicationMappings() (server.ApplicationMappingList, error) {
	return server.ApplicationMappingList{}, nil
}

// GetAllApplications returns all applications
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllApplications() (server.ApplicationList, error) {
	return server.ApplicationList{}, nil
}

// GetAllKeyMappings returns all key mappings
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllKeyMappings() (server.ApplicationKeyMappingList, error) {
	return server.ApplicationKeyMappingList{}, nil
}

// GetAllSubscriptions returns all subscriptions
func (k8sArtifactDeployer K8sArtifactDeployer) GetAllSubscriptions() (server.SubscriptionList, error) {
	return server.SubscriptionList{}, nil
}

// GetApplication returns an application
func (k8sArtifactDeployer K8sArtifactDeployer) GetApplication(applicationID string) (server.Application, error) {
	return server.Application{}, nil
}

// GetApplicationMappings returns an application mapping
func (k8sArtifactDeployer K8sArtifactDeployer) GetApplicationMappings(applicationID string) (server.ApplicationMapping, error) {
	return server.ApplicationMapping{}, nil
}

// GetKeyMappings returns a key mapping
func (k8sArtifactDeployer K8sArtifactDeployer) GetKeyMappings(applicationID string) (server.ApplicationKeyMapping, error) {
	return server.ApplicationKeyMapping{}, nil
}

// GetSubscription returns a subscription
func (k8sArtifactDeployer K8sArtifactDeployer) GetSubscription(subscriptionID string) (server.Subscription, error) {
	return server.Subscription{}, nil
}

// GenerateSecurityScheme generates a security scheme
func generateSecurityScheme(keyMapping server.ApplicationKeyMapping) cpv1alpha2.Environment {
	return cpv1alpha2.Environment{EnvID: keyMapping.EnvID, AppID: keyMapping.ApplicationIdentifier, KeyType: keyMapping.KeyType}
}

func (k8sArtifactDeployer K8sArtifactDeployer) retrieveAllApplicationsFromK8s(nextToken string) ([]cpv1alpha2.Application, string, error) {
	applicationList := cpv1alpha2.ApplicationList{}
	resolvedApplicationList := make([]cpv1alpha2.Application, 0)
	var err error
	if nextToken == "" {
		err = k8sArtifactDeployer.client.List(context.Background(), &applicationList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace()})
	} else {
		err = k8sArtifactDeployer.client.List(context.Background(), &applicationList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace(), Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedApplicationList = append(resolvedApplicationList, applicationList.Items...)
	if applicationList.Continue != "" {
		tempApplicationList, _, err := k8sArtifactDeployer.retrieveAllApplicationsFromK8s(applicationList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedApplicationList = append(resolvedApplicationList, tempApplicationList...)
	}
	return resolvedApplicationList, applicationList.Continue, nil
}

func (k8sArtifactDeployer K8sArtifactDeployer) retrieveAllSubscriptionsFromK8s(nextToken string) ([]cpv1alpha2.Subscription, string, error) {
	subscriptionList := cpv1alpha2.SubscriptionList{}
	resolvedSubscripitonList := make([]cpv1alpha2.Subscription, 0)
	var err error
	if nextToken == "" {
		err = k8sArtifactDeployer.client.List(context.Background(), &subscriptionList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace()})
	} else {
		err = k8sArtifactDeployer.client.List(context.Background(), &subscriptionList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace(), Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedSubscripitonList = append(resolvedSubscripitonList, subscriptionList.Items...)
	if subscriptionList.Continue != "" {
		tempSubscriptipnList, _, err := k8sArtifactDeployer.retrieveAllSubscriptionsFromK8s(subscriptionList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedSubscripitonList = append(resolvedSubscripitonList, tempSubscriptipnList...)
	}
	return resolvedSubscripitonList, subscriptionList.Continue, nil
}
func (k8sArtifactDeployer K8sArtifactDeployer) retrieveAllApplicationMappings(nextToken string) ([]cpv1alpha2.ApplicationMapping, string, error) {
	applicationMappingList := cpv1alpha2.ApplicationMappingList{}
	resolvedApplicationMappingList := make([]cpv1alpha2.ApplicationMapping, 0)
	var err error
	if nextToken == "" {
		err = k8sArtifactDeployer.client.List(context.Background(), &applicationMappingList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace()})
	} else {
		err = k8sArtifactDeployer.client.List(context.Background(), &applicationMappingList, &client.ListOptions{Namespace: utils.GetOperatorPodNamespace(), Continue: nextToken})
	}
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedApplicationMappingList = append(resolvedApplicationMappingList, applicationMappingList.Items...)
	if applicationMappingList.Continue != "" {
		tempApplicationMappingList, _, err := k8sArtifactDeployer.retrieveAllApplicationMappings(applicationMappingList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedApplicationMappingList = append(resolvedApplicationMappingList, tempApplicationMappingList...)
	}
	return resolvedApplicationMappingList, applicationMappingList.Continue, nil
}
