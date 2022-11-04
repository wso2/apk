package xds

import (
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/config/enforcer"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
)

var (
	// APIListMap has the following mapping label -> apiUUID -> API (Metadata)
	APIListMap map[string]map[string]*subscription.APIs
	// SubscriptionMap contains the subscriptions recieved from API Manager Control Plane
	SubscriptionMap map[int32]*subscription.Subscription
	// ApplicationMap contains the applications recieved from API Manager Control Plane
	ApplicationMap map[string]*subscription.Application
	// ApplicationKeyMappingMap contains the application key mappings recieved from API Manager Control Plane
	ApplicationKeyMappingMap map[string]*subscription.ApplicationKeyMapping
	// ApplicationPolicyMap contains the application policies recieved from API Manager Control Plane
	ApplicationPolicyMap map[int32]*subscription.ApplicationPolicy
	// SubscriptionPolicyMap contains the subscription policies recieved from API Manager Control Plane
	SubscriptionPolicyMap map[int32]*subscription.SubscriptionPolicy
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
	issuers := []*enforcer.Issuer{}
	urlGroups := []*enforcer.TMURLGroup{}

	for _, issuer := range config.Enforcer.Security.TokenService {
		claimMaps := []*enforcer.ClaimMapping{}
		for _, claimMap := range issuer.ClaimMapping {
			claim := &enforcer.ClaimMapping{
				RemoteClaim: claimMap.RemoteClaim,
				LocalClaim:  claimMap.LocalClaim,
			}
			claimMaps = append(claimMaps, claim)
		}
		jwtConfig := &enforcer.Issuer{
			CertificateAlias:     issuer.CertificateAlias,
			ConsumerKeyClaim:     issuer.ConsumerKeyClaim,
			Issuer:               issuer.Issuer,
			Name:                 issuer.Name,
			ValidateSubscription: issuer.ValidateSubscription,
			JwksURL:              issuer.JwksURL,
			CertificateFilePath:  issuer.CertificateFilePath,
			ClaimMapping:         claimMaps,
		}
		issuers = append(issuers, jwtConfig)
	}

	jwtUsers := []*enforcer.JWTUser{}
	for _, user := range config.Enforcer.JwtIssuer.JwtUser {
		jwtUser := &enforcer.JWTUser{
			Username: user.Username,
			Password: user.Password,
		}
		jwtUsers = append(jwtUsers, jwtUser)
	}

	for _, urlGroup := range config.Enforcer.Throttling.Publisher.URLGroup {
		group := &enforcer.TMURLGroup{
			AuthURLs:     urlGroup.AuthURLs,
			ReceiverURLs: urlGroup.ReceiverURLs,
			Type:         urlGroup.Type,
		}
		urlGroups = append(urlGroups, group)
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
		Enabled:          config.Analytics.Enabled,
		Type:             config.Analytics.Type,
		ConfigProperties: config.Analytics.Enforcer.ConfigProperties,
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

	restServer := &enforcer.RestServer{
		Enable: config.Enforcer.RestServer.Enabled,
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
			Enable:                config.Enforcer.JwtGenerator.Enabled,
			Encoding:              config.Enforcer.JwtGenerator.Encoding,
			ClaimDialect:          config.Enforcer.JwtGenerator.ClaimDialect,
			ConvertDialect:        config.Enforcer.JwtGenerator.ConvertDialect,
			Header:                config.Enforcer.JwtGenerator.Header,
			SigningAlgorithm:      config.Enforcer.JwtGenerator.SigningAlgorithm,
			EnableUserClaims:      config.Enforcer.JwtGenerator.EnableUserClaims,
			GatewayGeneratorImpl:  config.Enforcer.JwtGenerator.GatewayGeneratorImpl,
			ClaimsExtractorImpl:   config.Enforcer.JwtGenerator.ClaimsExtractorImpl,
			PublicCertificatePath: config.Enforcer.JwtGenerator.PublicCertificatePath,
			PrivateKeyPath:        config.Enforcer.JwtGenerator.PrivateKeyPath,
			TokenTtl:              config.Enforcer.JwtGenerator.TokenTTL,
		},
		JwtIssuer: &enforcer.JWTIssuer{
			Enabled:               config.Enforcer.JwtIssuer.Enabled,
			Issuer:                config.Enforcer.JwtIssuer.Issuer,
			Encoding:              config.Enforcer.JwtIssuer.Encoding,
			ClaimDialect:          config.Enforcer.JwtIssuer.ClaimDialect,
			SigningAlgorithm:      config.Enforcer.JwtIssuer.SigningAlgorithm,
			PublicCertificatePath: config.Enforcer.JwtIssuer.PublicCertificatePath,
			PrivateKeyPath:        config.Enforcer.JwtIssuer.PrivateKeyPath,
			ValidityPeriod:        config.Enforcer.JwtIssuer.ValidityPeriod,
			JwtUsers:              jwtUsers,
		},
		AuthService: authService,
		Security: &enforcer.Security{
			TokenService: issuers,
			AuthHeader: &enforcer.AuthHeader{
				EnableOutboundAuthHeader: config.Enforcer.Security.AuthHeader.EnableOutboundAuthHeader,
				AuthorizationHeader:      config.Enforcer.Security.AuthHeader.AuthorizationHeader,
				TestConsoleHeaderName:    config.Enforcer.Security.AuthHeader.TestConsoleHeaderName,
			},
			MutualSSL: &enforcer.MutualSSL{
				CertificateHeader:               config.Enforcer.Security.MutualSSL.CertificateHeader,
				EnableClientValidation:          config.Enforcer.Security.MutualSSL.EnableClientValidation,
				ClientCertificateEncode:         config.Enforcer.Security.MutualSSL.ClientCertificateEncode,
				EnableOutboundCertificateHeader: config.Enforcer.Security.MutualSSL.EnableOutboundCertificateHeader,
			},
		},
		Cache:     cache,
		Tracing:   tracing,
		Metrics:   metrics,
		Analytics: analytics,
		Throttling: &enforcer.Throttling{
			EnableGlobalEventPublishing:        config.Enforcer.Throttling.EnableGlobalEventPublishing,
			EnableHeaderConditions:             config.Enforcer.Throttling.EnableHeaderConditions,
			EnableQueryParamConditions:         config.Enforcer.Throttling.EnableQueryParamConditions,
			EnableJwtClaimConditions:           config.Enforcer.Throttling.EnableJwtClaimConditions,
			JmsConnectionInitialContextFactory: config.Enforcer.Throttling.JmsConnectionInitialContextFactory,
			JmsConnectionProviderUrl:           config.Enforcer.Throttling.JmsConnectionProviderURL,
			Publisher: &enforcer.BinaryPublisher{
				Username: config.Enforcer.Throttling.Publisher.Username,
				Password: config.Enforcer.Throttling.Publisher.Password,
				UrlGroup: urlGroups,
				Pool: &enforcer.PublisherPool{
					InitIdleObjectDataPublishingAgents: config.Enforcer.Throttling.Publisher.Pool.InitIdleObjectDataPublishingAgents,
					MaxIdleDataPublishingAgents:        config.Enforcer.Throttling.Publisher.Pool.MaxIdleDataPublishingAgents,
					PublisherThreadPoolCoreSize:        config.Enforcer.Throttling.Publisher.Pool.PublisherThreadPoolCoreSize,
					PublisherThreadPoolKeepAliveTime:   config.Enforcer.Throttling.Publisher.Pool.PublisherThreadPoolKeepAliveTime,
					PublisherThreadPoolMaximumSize:     config.Enforcer.Throttling.Publisher.Pool.PublisherThreadPoolMaximumSize,
				},
				Agent: &enforcer.ThrottleAgent{
					BatchSize:                  config.Enforcer.Throttling.Publisher.Agent.BatchSize,
					Ciphers:                    config.Enforcer.Throttling.Publisher.Agent.Ciphers,
					CorePoolSize:               config.Enforcer.Throttling.Publisher.Agent.CorePoolSize,
					EvictionTimePeriod:         config.Enforcer.Throttling.Publisher.Agent.EvictionTimePeriod,
					KeepAliveTimeInPool:        config.Enforcer.Throttling.Publisher.Agent.KeepAliveTimeInPool,
					MaxIdleConnections:         config.Enforcer.Throttling.Publisher.Agent.MaxIdleConnections,
					MaxPoolSize:                config.Enforcer.Throttling.Publisher.Agent.MaxPoolSize,
					MaxTransportPoolSize:       config.Enforcer.Throttling.Publisher.Agent.MaxTransportPoolSize,
					MinIdleTimeInPool:          config.Enforcer.Throttling.Publisher.Agent.MinIdleTimeInPool,
					QueueSize:                  config.Enforcer.Throttling.Publisher.Agent.QueueSize,
					ReconnectionInterval:       config.Enforcer.Throttling.Publisher.Agent.ReconnectionInterval,
					SecureEvictionTimePeriod:   config.Enforcer.Throttling.Publisher.Agent.SecureEvictionTimePeriod,
					SecureMaxIdleConnections:   config.Enforcer.Throttling.Publisher.Agent.SecureMaxIdleConnections,
					SecureMaxTransportPoolSize: config.Enforcer.Throttling.Publisher.Agent.SecureMaxTransportPoolSize,
					SecureMinIdleTimeInPool:    config.Enforcer.Throttling.Publisher.Agent.SecureMinIdleTimeInPool,
					SocketTimeoutMS:            config.Enforcer.Throttling.Publisher.Agent.SocketTimeoutMS,
					SslEnabledProtocols:        config.Enforcer.Throttling.Publisher.Agent.SslEnabledProtocols,
				},
			},
		},
		Management:          management,
		RestServer:          restServer,
		Filters:             filters,
		Soap:                soap,
		ControlPlaneEnabled: config.ControlPlane.Enabled,
	}
}
