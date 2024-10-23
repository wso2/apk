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
package database

import (
	"github.com/jackc/pgx/v5"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// DBDeployer is a struct that implements ArtifactDeployer interface
type DBDeployer struct {
	client client.Client
}

// NewDBArtifactDeployer creates a new NewDBDeployer
func NewDBArtifactDeployer(mgr manager.Manager) DBDeployer {
	populateMapFromDB()
	return DBDeployer{client: nil}
}

// populateMapFromDB populates the map from the database
func populateMapFromDB() error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, getAllApplications, getAllApplicationAttributes, getAllSubscriptions, getAllApplicationKeyMappings,
			getAllAppSubs)

		applications, err := GetAllApplications(tx)
		if err != nil {
			loggers.LoggerAPI.Error("Error while getting all applications ", err)
			return err
		}
		for _, app := range applications {
			server.AddApplication(app)
		}

		subscriptions, err := GetAllSubscription(tx)
		if err != nil {
			loggers.LoggerAPI.Error("Error while getting all subscriptions ", err)
			return err
		}
		for _, subscription := range subscriptions {
			server.AddSubscription(subscription)
		}

		appSubs, err := GetAllAppSubs(tx)
		if err != nil {
			loggers.LoggerAPI.Error("Error while getting all app subs ", err)
			return err
		}
		for _, appSub := range appSubs {
			server.AddApplicationMapping(appSub)
		}
		appKeyMappings, err := GetAllApplicationKeyMappings(tx)
		if err != nil {
			loggers.LoggerAPI.Error("Error while getting all app key mappings ", err)
			return err
		}
		for _, appKeyMapping := range appKeyMappings {
			server.AddApplicationKeyMapping(appKeyMapping)
		}
		return nil
	})
	return nil
}

// DeployApplication deploys an application
func (dbDeployer DBDeployer) DeployApplication(application server.Application) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		return deployApplicationwithAttributes(tx, application)
	})
	server.AddApplication(application)
	utils.SendApplicationEvent(constants.ApplicationCreated, application.UUID, application.Name, application.Owner,
		application.OrganizationID, application.Attributes)
	return nil
}

// UpdateApplication updates an application
func (dbDeployer DBDeployer) UpdateApplication(application server.Application) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, updateApplication, insertApplicationAttributes, deleteAllAppAttributes)
		if err := UpdateApplication(tx, application.UUID, application.Name, application.Owner, application.OrganizationID); err != nil {
			loggers.LoggerAPI.Error("Error while updating application ", err)
			return err
		}
		return updateApplicationAttributes(tx, application)
	})
	server.DeleteApplication(application.UUID)
	server.AddApplication(application)
	utils.SendApplicationEvent(constants.ApplicationUpdated, application.UUID, application.Name, application.Owner,
		application.OrganizationID, application.Attributes)
	return nil
}

// DeploySubscription deploys a subscription
func (dbDeployer DBDeployer) DeploySubscription(subscription server.Subscription) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, insertSubscription)
		return AddSubscription(tx, subscription.UUID, subscription.SubscribedAPI.Name, subscription.SubscribedAPI.Version,
			subscription.SubStatus, subscription.Organization, subscription.RatelimitTier)
	})
	server.AddSubscription(subscription)
	utils.SendSubscriptionEvent(constants.SubscriptionCreated, subscription.UUID, subscription.SubStatus, subscription.Organization,
		subscription.SubscribedAPI.Name, subscription.SubscribedAPI.Version, subscription.RatelimitTier)
	return nil
}

// UpdateSubscription updates a subscription
func (dbDeployer DBDeployer) UpdateSubscription(subscription server.Subscription) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, updateSubscription)
		return UpdateSubscription(tx, subscription.UUID, subscription.SubscribedAPI.Name, subscription.SubscribedAPI.Version,
			subscription.SubStatus, subscription.Organization, subscription.RatelimitTier)
	})
	server.DeleteSubscription(subscription.UUID)
	server.AddSubscription(subscription)
	utils.SendSubscriptionEvent(constants.SubscriptionUpdated, subscription.UUID, subscription.SubStatus, subscription.Organization,
		subscription.SubscribedAPI.Name, subscription.SubscribedAPI.Version, subscription.RatelimitTier)
	return nil
}

// DeployApplicationMappings deploys an application mapping
func (dbDeployer DBDeployer) DeployApplicationMappings(applicationMapping server.ApplicationMapping) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, insertAppSub)
		return AddAppSub(tx, applicationMapping.UUID, applicationMapping.ApplicationRef, applicationMapping.SubscriptionRef,
			applicationMapping.OrganizationID)
	})
	server.AddApplicationMapping(applicationMapping)
	utils.SendApplicationMappingEvent(constants.ApplicationMappingCreated, applicationMapping.UUID, applicationMapping.ApplicationRef,
		applicationMapping.SubscriptionRef, applicationMapping.OrganizationID)
	return nil
}

// DeployKeyMappings deploys a key mapping
func (dbDeployer DBDeployer) DeployKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, insertApplicationKeyMapping)
		if keyMapping.SecurityScheme == constants.OAuth2 {
			if err := AddApplicationKeyMapping(tx, keyMapping.ApplicationUUID, constants.OAuth2, keyMapping.ApplicationIdentifier,
				keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID); err != nil {
				loggers.LoggerAPI.Error("Error while adding application key mapping ", err)
				return err
			}
		}
		return nil
	})

	server.AddApplicationKeyMapping(keyMapping)
	utils.SendApplicationKeyMappingEvent(constants.ApplicationKeyMappingCreated, keyMapping.ApplicationUUID, keyMapping.SecurityScheme,
		keyMapping.ApplicationIdentifier, keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID)
	return nil
}

// DeleteAllApplicationMappings deletes all application mappings
func (dbDeployer DBDeployer) DeleteAllApplicationMappings() error {
	return nil
}

// DeleteAllKeyMappings deletes all key mappings
func (dbDeployer DBDeployer) DeleteAllKeyMappings() error {
	return nil
}

// DeleteAllSubscriptions deletes all subscriptions
func (dbDeployer DBDeployer) DeleteAllSubscriptions() error {
	return nil
}

// DeleteApplication deletes an application
func (dbDeployer DBDeployer) DeleteApplication(applicationID string) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteAllApplications, deleteAllAppAttributes)
		if err := DeleteApplication(tx, applicationID); err != nil {
			loggers.LoggerAPI.Error("Error while deleting application ", err)
			return err
		}
		return DeleteApplicationAttributes(tx, applicationID)
	})
	server.DeleteApplication(applicationID)
	utils.SendApplicationEvent(constants.ApplicationDeleted, applicationID, "", "", "", nil)
	return nil
}

// DeleteApplicationMappings deletes an application mapping
func (dbDeployer DBDeployer) DeleteApplicationMappings(applicationMapping string) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteAppSub)
		return DeleteAppSub(tx, applicationMapping)
	})
	server.DeleteApplicationMapping(applicationMapping)
	utils.SendApplicationMappingEvent(constants.ApplicationMappingDeleted, applicationMapping, "", "", "")
	return nil
}

// UpdateApplicationMappings updates an application mapping
func (dbDeployer DBDeployer) UpdateApplicationMappings(applicationMapping server.ApplicationMapping) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, updateAppSub)
		return UpdateAppSub(tx, applicationMapping.UUID, applicationMapping.ApplicationRef, applicationMapping.SubscriptionRef,
			applicationMapping.OrganizationID)
	})
	server.DeleteApplicationMapping(applicationMapping.UUID)
	server.AddApplicationMapping(applicationMapping)
	utils.SendApplicationMappingEvent(constants.ApplicationMappingUpdated, applicationMapping.UUID, applicationMapping.ApplicationRef,
		applicationMapping.SubscriptionRef, applicationMapping.OrganizationID)
	return nil
}

// DeleteKeyMappings deletes a key mapping
func (dbDeployer DBDeployer) DeleteKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteApplicationKeyMapping)
		return DeleteApplicationKeyMapping(tx, keyMapping.ApplicationUUID, keyMapping.SecurityScheme, keyMapping.EnvID)
	})
	server.DeleteApplicationKeyMapping(keyMapping)
	utils.SendApplicationKeyMappingEvent(constants.ApplicationKeyMappingDeleted, keyMapping.ApplicationUUID, keyMapping.SecurityScheme,
		keyMapping.ApplicationIdentifier, keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID)
	return nil
}

// DeleteSubscription deletes a subscription
func (dbDeployer DBDeployer) DeleteSubscription(subscriptionID string) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteSubscription)
		return DeleteSubscription(tx, subscriptionID)
	})
	server.DeleteSubscription(subscriptionID)
	utils.SendSubscriptionEvent(constants.SubscriptionDeleted, subscriptionID, "", "", "", "", "")
	return nil
}

// DeployAllApplicationMappings deploys all application mappings
func (dbDeployer DBDeployer) DeployAllApplicationMappings(applicationMappings server.ApplicationMappingList) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, insertAppSub, deleteAllAppSub)
		if err := DeleteAllAppSub(tx); err != nil {
			loggers.LoggerAPI.Error("Error while deleting all app sub ", err)
			return err
		}
		server.DeleteAllApplicationMappings()
		for _, applicationMapping := range applicationMappings.List {
			if err := AddAppSub(tx, applicationMapping.UUID, applicationMapping.ApplicationRef, applicationMapping.SubscriptionRef,
				applicationMapping.OrganizationID); err != nil {
				loggers.LoggerAPI.Error("Error while adding app sub ", err)
				return err
			}
		}
		return nil
	})
	for _, applicationMapping := range applicationMappings.List {
		server.AddApplicationMapping(applicationMapping)
	}
	utils.SendResetEvent()
	return nil
}

// DeployAllApplications deploys all key mappings
func (dbDeployer DBDeployer) DeployAllApplications(applications server.ApplicationList) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteAllApplications, deleteAllAppAttributes, insertApplication, insertApplicationAttributes)
		if err := DeleteAllApplications(tx); err != nil {
			loggers.LoggerAPI.Error("Error while deleting all applications ", err)
			return err
		}
		if err := DeleteAllAppAttributes(tx); err != nil {
			loggers.LoggerAPI.Error("Error while deleting all app attributes ", err)
			return err
		}

		for _, application := range applications.List {
			if err := deployApplicationwithAttributes(tx, application); err != nil {
				loggers.LoggerAPI.Error("Error while deploying application with attributes ", err)
				return err
			}
		}
		return nil
	})
	server.DeleteAllApplications()
	for _, application := range applications.List {
		server.AddApplication(application)
	}
	utils.SendResetEvent()
	return nil
}

// UpdateKeyMappings updates a key mapping
func (dbDeployer DBDeployer) UpdateKeyMappings(keyMapping server.ApplicationKeyMapping) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, updateApplicationKeyMapping)
		if keyMapping.SecurityScheme == constants.OAuth2 {
			if err := UpdateApplicationKeyMapping(tx, keyMapping.ApplicationUUID, keyMapping.SecurityScheme, keyMapping.ApplicationIdentifier,
				keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID); err != nil {
				loggers.LoggerAPI.Error("Error while updating application key mapping ", err)
				return err
			}
		}
		return nil
	})

	server.DeleteApplicationKeyMapping(keyMapping)
	server.AddApplicationKeyMapping(keyMapping)
	utils.SendApplicationKeyMappingEvent(constants.ApplicationKeyMappingUpdated, keyMapping.ApplicationUUID, keyMapping.SecurityScheme,
		keyMapping.ApplicationIdentifier, keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID)
	return nil
}

// DeployAllKeyMappings deploys all key mappings
func (dbDeployer DBDeployer) DeployAllKeyMappings(keyMappings server.ApplicationKeyMappingList) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteAllApplicationKeyMappings, insertApplicationKeyMapping)
		if err := DeleteAllApplicationKeyMappings(tx); err != nil {
			loggers.LoggerAPI.Error("Error while deleting all application key mappings ", err)
			return err
		}
		server.DeleteAllApplicationKeyMappings()
		for _, keyMapping := range keyMappings.List {
			if keyMapping.SecurityScheme == constants.OAuth2 {
				if err := AddApplicationKeyMapping(tx, keyMapping.ApplicationUUID, constants.OAuth2, keyMapping.ApplicationIdentifier,
					keyMapping.KeyType, keyMapping.EnvID, keyMapping.OrganizationID); err != nil {
					loggers.LoggerAPI.Error("Error while adding application key mapping ", err)
					return err
				}

			}
		}
		return nil
	})

	server.DeleteAllApplicationKeyMappings()
	for _, keyMapping := range keyMappings.List {
		server.AddApplicationKeyMapping(keyMapping)
	}

	utils.SendResetEvent()
	return nil
}

// DeployAllSubscriptions deploys all subscriptions
func (dbDeployer DBDeployer) DeployAllSubscriptions(subscriptions server.SubscriptionList) error {
	retryUntilTransaction(func(tx pgx.Tx) error {
		PrepareQueries(tx, deleteAllSubscriptions, insertSubscription)
		if err := DeleteAllSubscriptions(tx); err != nil {
			loggers.LoggerAPI.Error("Error while deleting all subscriptions ", err)
			return err
		}
		server.DeleteAllSubscriptions()
		for _, subscription := range subscriptions.List {
			if err := AddSubscription(tx, subscription.UUID, subscription.SubscribedAPI.Name, subscription.SubscribedAPI.Version,
				subscription.SubStatus, subscription.Organization, subscription.RatelimitTier); err != nil {
				loggers.LoggerAPI.Error("Error while adding subscription ", err)
				return err
			}

		}
		return nil
	})

	for _, subscription := range subscriptions.List {
		server.AddSubscription(subscription)
	}
	utils.SendResetEvent()
	return nil
}

// GetAllApplicationMappings returns all application mappings
func (dbDeployer DBDeployer) GetAllApplicationMappings() (server.ApplicationMappingList, error) {
	return server.ApplicationMappingList{}, nil
}

// GetAllApplications returns all applications
func (dbDeployer DBDeployer) GetAllApplications() (server.ApplicationList, error) {
	return server.ApplicationList{}, nil
}

// GetAllKeyMappings returns all key mappings
func (dbDeployer DBDeployer) GetAllKeyMappings() (server.ApplicationKeyMappingList, error) {
	return server.ApplicationKeyMappingList{}, nil
}

// GetAllSubscriptions returns all subscriptions
func (dbDeployer DBDeployer) GetAllSubscriptions() (server.SubscriptionList, error) {
	return server.SubscriptionList{}, nil
}

// GetApplication returns an application
func (dbDeployer DBDeployer) GetApplication(applicationID string) (server.Application, error) {
	return server.Application{}, nil
}

// GetApplicationMappings returns an application mapping
func (dbDeployer DBDeployer) GetApplicationMappings(applicationID string) (server.ApplicationMapping, error) {
	return server.ApplicationMapping{}, nil
}

// GetKeyMappings returns a key mapping
func (dbDeployer DBDeployer) GetKeyMappings(applicationID string) (server.ApplicationKeyMapping, error) {
	return server.ApplicationKeyMapping{}, nil
}

// GetSubscription returns a subscription
func (dbDeployer DBDeployer) GetSubscription(subscriptionID string) (server.Subscription, error) {
	return server.Subscription{}, nil
}
