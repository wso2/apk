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
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/config/enforcer"
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
	mandateSubscriptionValidation := config.Enforcer.MandateSubscriptionValidation
	mandateInternalKeyValidation := config.Enforcer.MandateInternalKeyValidation

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
	httpClient := &enforcer.HttpClient{
		SkipSSl:                config.Enforcer.Client.SkipSSL,
		HostnameVerifier:       config.Enforcer.Client.HostnameVerifier,
		MaxTotalConnections:    int32(config.Enforcer.Client.MaxTotalConnectins),
		MaxConnectionsPerRoute: int32(config.Enforcer.Client.MaxPerHostConnectins),
		ConnectTimeout:         int32(config.Enforcer.Client.ConnectionTimeout),
		SocketTimeout:          int32(config.Enforcer.Client.SocketTimeout),
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
		Cache:                         cache,
		Tracing:                       tracing,
		Metrics:                       metrics,
		Analytics:                     analytics,
		Management:                    management,
		Filters:                       filters,
		Soap:                          soap,
		MandateSubscriptionValidation: mandateSubscriptionValidation,
		MandateInternalKeyValidation:  mandateInternalKeyValidation,
		HttpClient:                    httpClient,
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
