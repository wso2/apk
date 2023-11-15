/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package xds

import (
	"strconv"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/config/enforcer"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/adapter/pkg/eventhub/types"
)

var (
	// APIListMap has the following mapping label -> apiUUID -> API (Metadata)
	APIListMap map[string]map[string]*subscription.APIs
)

// EventType is a enum to distinguish Create, Update and Delete Events
type EventType int

const (
	// CreateEvent : enum
	CreateEvent EventType = iota
	// UpdateEvent : enum
	UpdateEvent
	// DeleteEvent : enum
	DeleteEvent
)

const blockedStatus string = "BLOCKED"

// MarshalConfig will marshal a Config struct - read from the config toml - to
// enfocer's CDS resource representation.
func MarshalConfig(config *config.Config) *enforcer.Config {

	keyPairs := []*enforcer.Keypair{}

	// New configuration
	for _, kp := range config.Enforcer.JwtGenerator.Keypair {
		keypair := &enforcer.Keypair{
			PublicCertificatePath: kp.PublicCertificatePath,
			PrivateKeyPath:        kp.PrivateKeyPath,
			UseForSigning:         kp.UseForSigning,
		}

		keyPairs = append(keyPairs, keypair)
	}

	authService := &enforcer.Service{
		KeepAliveTime:  config.Enforcer.AuthService.KeepAliveTime,
		MaxHeaderLimit: config.Enforcer.AuthService.MaxHeaderLimit,
		MaxMessageSize: config.Enforcer.AuthService.MaxMessageSize,
		Port:           config.Enforcer.AuthService.Port,
		ThreadPool: &enforcer.ThreadPool{
			CoreSize:      config.Enforcer.AuthService.ThreadPool.CoreSize,
			KeepAliveTime: config.Enforcer.AuthService.ThreadPool.KeepAliveTime,
			MaxSize:       config.Enforcer.AuthService.ThreadPool.MaxSize,
			QueueSize:     config.Enforcer.AuthService.ThreadPool.QueueSize,
		},
	}

	cache := &enforcer.Cache{
		Enable:      config.Enforcer.Cache.Enabled,
		MaximumSize: config.Enforcer.Cache.MaximumSize,
		ExpiryTime:  config.Enforcer.Cache.ExpiryTime,
	}

	tracing := &enforcer.Tracing{
		Enabled:          config.Tracing.Enabled,
		Type:             config.Tracing.Type,
		ConfigProperties: config.Tracing.ConfigProperties,
	}
	metrics := &enforcer.Metrics{
		Enabled: config.Enforcer.Metrics.Enabled,
		Type:    config.Enforcer.Metrics.Type,
	}
	analytics := &enforcer.Analytics{
		Enabled:            config.Analytics.Enabled,
		Properties:         config.Analytics.Properties,
		AnalyticsPublisher: marshalAnalyticsPublishers(*config),
		Service: &enforcer.Service{
			Port:           config.Analytics.Enforcer.LogReceiver.Port,
			MaxHeaderLimit: config.Analytics.Enforcer.LogReceiver.MaxHeaderLimit,
			KeepAliveTime:  config.Analytics.Enforcer.LogReceiver.KeepAliveTime,
			MaxMessageSize: config.Analytics.Enforcer.LogReceiver.MaxMessageSize,
			ThreadPool: &enforcer.ThreadPool{
				CoreSize:      config.Analytics.Enforcer.LogReceiver.ThreadPool.CoreSize,
				MaxSize:       config.Analytics.Enforcer.LogReceiver.ThreadPool.MaxSize,
				QueueSize:     config.Analytics.Enforcer.LogReceiver.ThreadPool.QueueSize,
				KeepAliveTime: config.Analytics.Enforcer.LogReceiver.ThreadPool.KeepAliveTime,
			},
		},
	}

	management := &enforcer.Management{
		Username: config.Enforcer.Management.Username,
		Password: config.Enforcer.Management.Password,
	}

	soap := &enforcer.Soap{
		SoapErrorInXMLEnabled: config.Adapter.SoapErrorInXMLEnabled,
	}

	filters := []*enforcer.Filter{}

	for _, filterConfig := range config.Enforcer.Filters {
		filter := &enforcer.Filter{
			ClassName:        filterConfig.ClassName,
			Position:         filterConfig.Position,
			ConfigProperties: filterConfig.ConfigProperties,
		}
		filters = append(filters, filter)
	}

	return &enforcer.Config{
		JwtGenerator: &enforcer.JWTGenerator{
			Keypairs: keyPairs,
		},
		AuthService: authService,
		Security: &enforcer.Security{
			ApiKey: &enforcer.APIKeyEnforcer{
				Enabled:             config.Enforcer.Security.APIkey.Enabled,
				Issuer:              config.Enforcer.Security.APIkey.Issuer,
				CertificateFilePath: config.Enforcer.Security.APIkey.CertificateFilePath,
			},
			RuntimeToken: &enforcer.APIKeyEnforcer{
				Enabled:             config.Enforcer.Security.InternalKey.Enabled,
				Issuer:              config.Enforcer.Security.InternalKey.Issuer,
				CertificateFilePath: config.Enforcer.Security.InternalKey.CertificateFilePath,
			},
			MutualSSL: &enforcer.MutualSSL{
				CertificateHeader:               config.Enforcer.Security.MutualSSL.CertificateHeader,
				EnableClientValidation:          config.Enforcer.Security.MutualSSL.EnableClientValidation,
				ClientCertificateEncode:         config.Enforcer.Security.MutualSSL.ClientCertificateEncode,
				EnableOutboundCertificateHeader: config.Enforcer.Security.MutualSSL.EnableOutboundCertificateHeader,
			},
		},
		Cache:      cache,
		Tracing:    tracing,
		Metrics:    metrics,
		Analytics:  analytics,
		Management: management,
		Filters:    filters,
		Soap:       soap,
	}
}

func marshalAnalyticsPublishers(config config.Config) []*enforcer.AnalyticsPublisher {
	analyticsPublishers := config.Analytics.Enforcer.Publisher
	resolvedAnalyticsPublishers := make([]*enforcer.AnalyticsPublisher, len(analyticsPublishers))
	for i, publisher := range analyticsPublishers {
		resolvedAnalyticsPublishers[i] = &enforcer.AnalyticsPublisher{Enabled: publisher.Enabled,
			Type:             publisher.Type,
			ConfigProperties: publisher.ConfigProperties}
	}
	return resolvedAnalyticsPublishers
}

// marshalAPIListMapToList converts the data into APIList proto type
func marshalAPIListMapToList(apiMap map[string]*subscription.APIs) *subscription.APIList {
	apis := []*subscription.APIs{}
	for _, api := range apiMap {
		apis = append(apis, api)
	}

	return &subscription.APIList{
		List: apis,
	}
}

// MarshalAPIMetataAndReturnList updates the internal APIListMap and returns the XDS compatible APIList.
// apiList is the internal APIList object (For single API, this would contain a List with just one API)
// initialAPIUUIDListMap is assigned during startup when global adapter is associated. This would be empty otherwise.
// gatewayLabel is the environment.
func MarshalAPIMetataAndReturnList(apiList *types.APIList, initialAPIUUIDListMap map[string]int, gatewayLabel string) *subscription.APIList {

	if APIListMap == nil {
		APIListMap = make(map[string]map[string]*subscription.APIs)
	}
	// var resourceMapForLabel map[string]*subscription.APIs
	if _, ok := APIListMap[gatewayLabel]; !ok {
		APIListMap[gatewayLabel] = make(map[string]*subscription.APIs)
	}
	resourceMapForLabel := APIListMap[gatewayLabel]
	for item := range apiList.List {
		api := apiList.List[item]
		// initialAPIUUIDListMap is not null if the adapter is running with global adapter enabled, and it is
		// the first method invocation.
		if initialAPIUUIDListMap != nil {
			if _, ok := initialAPIUUIDListMap[api.UUID]; !ok {
				continue
			}
		}
		newAPI := marshalAPIMetadata(&api)
		resourceMapForLabel[api.UUID] = newAPI
	}
	return marshalAPIListMapToList(resourceMapForLabel)
}

// DeleteAPIAndReturnList removes the API from internal maps and returns the marshalled API List.
// If the apiUUID is not found in the internal map under the provided environment, then it would return a
// nil value. Hence it is required to check if the return value is nil, prior to updating the XDS cache.
func DeleteAPIAndReturnList(apiUUID, organizationUUID string, gatewayLabel string) *subscription.APIList {
	if _, ok := APIListMap[gatewayLabel]; !ok {
		logger.LoggerXds.Debugf("No API Metadata is available under gateway Environment : %s", gatewayLabel)
		return nil
	}
	delete(APIListMap[gatewayLabel], apiUUID)
	return marshalAPIListMapToList(APIListMap[gatewayLabel])
}

// MarshalAPIForLifeCycleChangeEventAndReturnList updates the internal map's API instances lifecycle state only if
// stored API Instance's or input status event is a blocked event.
// If no change is applied, it would return nil. Hence the XDS cache should not be updated.
func MarshalAPIForLifeCycleChangeEventAndReturnList(apiUUID, status, gatewayLabel string) *subscription.APIList {
	if _, ok := APIListMap[gatewayLabel]; !ok {
		logger.LoggerXds.Debugf("No API Metadata is available under gateway Environment : %s", gatewayLabel)
		return nil
	}
	if _, ok := APIListMap[gatewayLabel][apiUUID]; !ok {
		logger.LoggerXds.Debugf("No API Metadata for API ID: %s is available under gateway Environment : %s",
			apiUUID, gatewayLabel)
		return nil
	}
	storedAPILCState := APIListMap[gatewayLabel][apiUUID].LcState

	// Because the adapter only required to update the XDS if it is related to blocked state.
	if !(storedAPILCState == blockedStatus || status == blockedStatus) {
		return nil
	}
	APIListMap[gatewayLabel][apiUUID].LcState = status
	return marshalAPIListMapToList(APIListMap[gatewayLabel])
}

func marshalAPIMetadata(api *types.API) *subscription.APIs {
	return &subscription.APIs{
		ApiId:            strconv.Itoa(api.APIID),
		Name:             api.Name,
		Provider:         api.Provider,
		Version:          api.Version,
		BasePath:         api.BasePath,
		Policy:           api.Policy,
		ApiType:          api.APIType,
		Uuid:             api.UUID,
		IsDefaultVersion: api.IsDefaultVersion,
		LcState:          api.APIStatus,
	}
}

// CheckIfAPIMetadataIsAlreadyAvailable returns true only if the API Metadata for the given API UUID
// is already available
func CheckIfAPIMetadataIsAlreadyAvailable(apiUUID, label string) bool {
	if _, labelAvailable := APIListMap[label]; labelAvailable {
		if _, apiAvailale := APIListMap[label][apiUUID]; apiAvailale {
			return true
		}
	}
	return false
}
