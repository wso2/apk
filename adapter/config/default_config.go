/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package config

// Configuration object which is populated with default values.
var defaultConfig = &Config{
	Adapter: adapter{
		Consul: consul{
			Enabled:            false,
			URL:                "https://169.254.1.1:8501",
			PollInterval:       5,
			ACLToken:           "d3a2a719-4221-8c65-5212-58d4727427ac",
			ApkServiceName:     "wso2",
			ServiceMeshEnabled: false,
			CaCertFile:         "/home/wso2/security/truststore/consul/consul-agent-ca.pem",
			CertFile:           "/home/wso2/security/truststore/consul/local-dc-client-consul-0.pem",
			KeyFile:            "/home/wso2/security/truststore/consul/local-dc-client-consul-0-key.pem",
		},
		Keystore: keystore{
			KeyPath:  "/home/wso2/security/keystore/mg.key",
			CertPath: "/home/wso2/security/keystore/mg.pem",
		},
		Truststore: truststore{
			Location: "/home/wso2/security/truststore",
		},
		SoapErrorInXMLEnabled: false,
		Operator: operator{
			Namespaces: nil,
		},
	},
	Envoy: envoy{
		ListenerCodecType:                "AUTO",
		ClusterTimeoutInSeconds:          20,
		EnforcerResponseTimeoutInSeconds: 20,
		UseRemoteAddress:                 false,
		KeyStore: keystore{
			KeyPath:  "/ambassador/security/keystore/mg.key",
			CertPath: "/ambassador/security/keystore/mg.pem",
		},
		SystemHost: "localhost",
		Upstream: envoyUpstream{
			TLS: upstreamTLS{
				MinimumProtocolVersion: "TLS1_1",
				MaximumProtocolVersion: "TLS1_2",
				Ciphers: "ECDHE-ECDSA-AES128-GCM-SHA256, ECDHE-RSA-AES128-GCM-SHA256, ECDHE-ECDSA-AES128-SHA, ECDHE-RSA-AES128-SHA, " +
					"AES128-GCM-SHA256, AES128-SHA, ECDHE-ECDSA-AES256-GCM-SHA384, ECDHE-RSA-AES256-GCM-SHA384, " +
					"ECDHE-ECDSA-AES256-SHA, ECDHE-RSA-AES256-SHA, AES256-GCM-SHA384, AES256-SHA",
				TrustedCertPath:        "/etc/ssl/certs/ca-certificates.crt",
				VerifyHostName:         true,
				DisableSslVerification: false,
			},
			Timeouts: upstreamTimeout{
				MaxRouteTimeoutInSeconds:  60,
				RouteTimeoutInSeconds:     60,
				RouteIdleTimeoutInSeconds: 300,
			},
			Health: upstreamHealth{
				Timeout:            1,
				Interval:           10,
				UnhealthyThreshold: 2,
				HealthyThreshold:   2,
			},
			Retry: upstreamRetry{
				MaxRetryCount:        5,
				BaseIntervalInMillis: 25,
				StatusCodes:          []uint32{504},
			},
			DNS: upstreamDNS{
				DNSRefreshRate: 5000,
				RespectDNSTtl:  false,
			},
			HTTP2: upstreamHTTP2Options{
				HpackTableSize:       4096,
				MaxConcurrentStreams: 2147483647,
			},
		},
		Downstream: envoyDownstream{
			TLS: downstreamTLS{
				TrustedCertPath: "/etc/ssl/certs/ca-certificates.crt",
				MTLSAPIsEnabled: false,
			},
		},
		Connection: connection{
			Timeouts: connectionTimeouts{
				RequestTimeoutInSeconds:        0,
				RequestHeadersTimeoutInSeconds: 0,
				StreamIdleTimeoutInSeconds:     300,
				IdleTimeoutInSeconds:           3600,
			},
		},
		PayloadPassingToEnforcer: payloadPassingToEnforcer{
			PassRequestPayload:  false,
			MaxRequestBytes:     102400,
			AllowPartialMessage: false,
			PackAsBytes:         false,
		},
		Filters: filters{
			Compression: compression{
				Enabled: true,
				Library: "gzip",
				RequestDirection: requestDirection{
					Enabled:              false,
					MinimumContentLength: 30,
					ContentType:          []string{"application/javascript", "application/json", "application/xhtml+xml", "image/svg+xml", "text/css", "text/html", "text/plain", "text/xml"},
				},
				ResponseDirection: responseDirection{
					Enabled:              true,
					MinimumContentLength: 30,
					ContentType:          []string{"application/javascript", "application/json", "application/xhtml+xml", "image/svg+xml", "text/css", "text/html", "text/plain", "text/xml"},
					EnableForEtagHeader:  true,
				},
				LibraryProperties: map[string]interface{}{
					"memoryLevel":         uint32(3),
					"windowBits":          uint32(12),
					"compressionLevel":    uint32(9),
					"compressionStrategy": "defaultStrategy",
					"chunkSize":           uint32(4096),
				},
			},
		},
		RateLimit: rateLimit{
			Enabled: false,
			Host:    "ratelimiter",
			Port:    8091,
			XRateLimitHeaders: xRateLimitHeaders{
				Enabled:    true,
				RFCVersion: "DRAFT_VERSION_03",
			},
			FailureModeDeny:        false,
			RequestTimeoutInMillis: 80,
			KeyFilePath:            "/ambassador/security/keystore/router.key",
			CertFilePath:           "/ambassador/security/keystore/router.crt",
			CaCertFilePath:         "/ambassador/security/truststore/ratelimiter.crt",
			SSLCertSANHostname:     "",
		},
	},
	Enforcer: enforcer{
		Management: management{
			Username: "admin",
			Password: "admin",
		},
		RestServer: restServer{
			Enabled: true,
			Enable:  true,
		},
		Security: security{
			APIkey: apiKey{
				Enabled:             true,
				Issuer:              "",
				CertificateFilePath: "/home/wso2/security/truststore/wso2carbon.pem",
			},
			InternalKey: internalKey{
				Enabled:             true,
				Issuer:              "https://localhost:9443/publisher",
				CertificateFilePath: "/home/wso2/security/truststore/wso2carbon.pem",
			},
			MutualSSL: mutualSSL{
				CertificateHeader:               "X-WSO2-CLIENT-CERTIFICATE",
				EnableClientValidation:          true,
				ClientCertificateEncode:         false,
				EnableOutboundCertificateHeader: false,
			},
		},
		AuthService: authService{
			Port:           8081,
			MaxMessageSize: 1000000000,
			MaxHeaderLimit: 8192,
			KeepAliveTime:  600,
			ThreadPool: threadPool{
				CoreSize:      400,
				MaxSize:       500,
				KeepAliveTime: 600,
				QueueSize:     1000,
			},
		},
		JwtGenerator: jwtGenerator{
			PublicCertificatePath: "/home/wso2/security/truststore/mg.pem",
			PrivateKeyPath:        "/home/wso2/security/keystore/mg.key",
		},
		Cache: cache{
			Enabled:     true,
			MaximumSize: 10000,
			ExpiryTime:  15,
		},
		Metrics: metrics{
			Enabled: false,
			Type:    "azure",
		},
		JwtIssuer: jwtIssuer{
			Enabled:               true,
			Issuer:                "https://localhost:9095/testkey",
			Encoding:              "base64",
			ClaimDialect:          "",
			SigningAlgorithm:      "SHA256withRSA",
			PublicCertificatePath: "/home/wso2/security/truststore/mg.pem",
			PrivateKeyPath:        "/home/wso2/security/keystore/mg.key",
			ValidityPeriod:        3600,
			JwtUser: []JwtUser{
				{
					Username: "admin",
					Password: "$env{enforcer_admin_pwd}",
				},
			},
		},
	},
	ManagementServer: managementServer{
		Enabled:   false,
		Host:      "management-server",
		XDSPort:   18000,
		NodeLabel: "default",
		GRPCClient: gRPCClient{
			Port:                  8765,
			MaxAttempts:           5,
			BackOffInMilliSeconds: 1000,
		},
	},
	PartitionServer: partitionServer{
		Enabled:                false,
		Host:                   "partition-server",
		Port:                   9443,
		ServiceBasePath:        "/partition-service",
		PartitionName:          "default",
		DisableSslVerification: false,
	},
	Analytics: analytics{
		Enabled: false,
		Type:    "Default",
		Adapter: analyticsAdapter{
			BufferFlushInterval: 1000000000,
			BufferSizeBytes:     16384,
			GRPCRequestTimeout:  20000000000,
		},
		Enforcer: analyticsEnforcer{
			ConfigProperties: map[string]string{
				"authURL":   "$env{analytics_authURL}",
				"authToken": "$env{analytics_authToken}",
			},
			LogReceiver: authService{
				Port:           18090,
				MaxMessageSize: 1000000000,
				MaxHeaderLimit: 8192,
				KeepAliveTime:  600,
				ThreadPool: threadPool{
					CoreSize:      10,
					MaxSize:       100,
					KeepAliveTime: 600,
					QueueSize:     1000,
				},
			},
		},
	},
	Tracing: tracing{
		Enabled: false,
		Type:    "zipkin",
		ConfigProperties: map[string]string{
			"libraryName":            "APK",
			"maximumTracesPerSecond": "2",
			"maxPathLength":          "256",
			"host":                   "jaeger",
			"port":                   "9411",
			"endpoint":               "/api/v2/spans",
		},
	},
}
