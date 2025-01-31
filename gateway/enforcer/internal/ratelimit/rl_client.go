package ratelimit

import (
	"context"
	"fmt"
	"time"

	v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	rls_svc "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// client is a client for managing gRPC connections to the Rate Limit Service (RLS).
type client struct {
	log       logging.Logger
	cfg       *config.Server
	rlsClient rls_svc.RateLimitServiceClient
}

// keyValueHitsAddend is a struct that holds the key, value, and hits addend for a rate limit.
type keyValueHitsAddend struct {
	Key                string
	Value              string
	HitsAddend         int
	KeyValueHitsAddend *keyValueHitsAddend
}

// newClient creates a new instance of the Rate Limit Client.
func newClient(cfg *config.Server) *client {
	return &client{
		log: cfg.Logger,
		cfg: cfg,
	}
}

// start initializes the Rate Limit Client by creating a gRPC connection to the RLS.
func (r *client) start() {
	r.log.Info("Starting the rate limit client")

	clientCert, err := util.LoadCertificates(r.cfg.EnforcerPublicKeyPath, r.cfg.EnforcerPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// Load the trusted CA certificates
	certPool, err := util.LoadCACertificates(r.cfg.TrustedAdapterCertsPath)
	if err != nil {
		panic(err)
	}

	// Create the TLS configuration
	tlsConfig := util.CreateTLSConfig(clientCert, certPool)
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(context.TODO(), r.cfg.RatelimiterHost, r.cfg.RatelimiterPort, tlsConfig, r.cfg.XdsMaxRetries, time.Duration(r.cfg.XdsRetryPeriod)*time.Millisecond)
	r.rlsClient = rls_svc.NewRateLimitServiceClient(grpcConn)
	r.log.Info("Rate limit client started successfully")
}

// shouldRatelimit checks if the request should be rate limited based on the given configurations.
func (r *client) shouldRatelimit(configs []*keyValueHitsAddend) {
	for _, config := range configs {
		descriptorEntries := []*v3.RateLimitDescriptor_Entry{
			{
				Key:   config.Key,
				Value: config.Value,
			},
		}

		internalConfig := config.KeyValueHitsAddend
		hitsAddend := config.HitsAddend
		for internalConfig != nil {
			descriptorEntries = append(descriptorEntries, &v3.RateLimitDescriptor_Entry{
				Key:   internalConfig.Key,
				Value: internalConfig.Value,
			})
			hitsAddend = internalConfig.HitsAddend
			internalConfig = internalConfig.KeyValueHitsAddend
		}

		rateLimitRequest := &rls_svc.RateLimitRequest{
			Descriptors: []*v3.RateLimitDescriptor{
				{
					Entries: descriptorEntries,
				},
			},
			Domain:     "Default",
			HitsAddend: uint32(hitsAddend),
		}

		response, err := r.rlsClient.ShouldRateLimit(context.Background(), rateLimitRequest)
		if err != nil {
			r.log.Info(fmt.Sprintf("Error while calling rate limiter: %v", err))
			continue
		}

		r.log.Info(fmt.Sprintf("Rate limit response: %v", response))
	}
}
