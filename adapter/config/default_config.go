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
		Keystore: keystore{
			KeyPath:  "/home/wso2/security/keystore/adapter.key",
			CertPath: "/home/wso2/security/keystore/adapter.crt",
		},
		Truststore: truststore{
			Location: "/home/wso2/security/truststore",
		},
		SoapErrorInXMLEnabled: false,
		Operator: operator{
			Namespaces: nil,
		},
		Environment: "Default",
		Metrics: Metrics{
			Enabled: false,
			Type:    "prometheus",
			Port:    18006,
		},
	},
	Envoy: envoy{
		ListenerCodecType: "AUTO",
		// todo(amali) move connect timeout to crd
		ClusterTimeoutInSeconds:          20,
		EnforcerResponseTimeoutInSeconds: 20,
		UseRemoteAddress:                 false,
		SystemHost:                       "localhost",
		KeyStore: keystore{
			KeyPath:  "/home/wso2/security/keystore/router.key",
			CertPath: "/home/wso2/security/keystore/router.crt",
		},
		// todo(amali) move to crd
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
			Retry: upstreamRetry{
				StatusCodes: []uint32{504},
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
		// todo(amali) move timeout to crd
		Connection: connection{
			Timeouts: connectionTimeouts{
				RequestTimeoutInSeconds:        0,
				RequestHeadersTimeoutInSeconds: 0,
				StreamIdleTimeoutInSeconds:     300,
				IdleTimeoutInSeconds:           3600,
			},
		},
		PayloadPassingToEnforcer: payloadPassingToEnforcer{
			MaxRequestBytes:     102400,
			AllowPartialMessage: false,
			PackAsBytes:         false,
		},
		//todo(amali) test
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
			KeyFilePath:            "/home/wso2/security/keystore/router.key",
			CertFilePath:           "/home/wso2/security/keystore/router.crt",
			CaCertFilePath:         "/home/wso2/security/truststore/ratelimiter.crt",
			SSLCertSANHostname:     "",
		},
		EnableIntelligentRouting: false,
	},
	Enforcer: enforcer{
		Management: management{
			Username: "admin",
			Password: "admin",
		},
		Client: httpClient{
			SkipSSL:              false,
			HostnameVerifier:     "BROWSER_COMPATIBLE",
			MaxTotalConnectins:   100,
			MaxPerHostConnectins: 10,
			ConnectionTimeout:    10000,
			SocketTimeout:        10000,
		},
		Security: security{
			APIkey: apiKey{
				Enabled:             true,
				Issuer:              "https://apim.wso2.com/publisher",
				CertificateFilePath: "/home/wso2/security/truststore/wso2carbon.pem",
			},
			InternalKey: internalKey{
				Enabled:             true,
				Issuer:              "https://apim.wso2.com/publisher",
				CertificateFilePath: "/home/wso2/security/truststore/wso2carbon.pem",
			},
			MutualSSL: mutualSSL{
				CertificateHeader:               "X-WSO2-CLIENT-CERTIFICATE",
				EnableClientValidation:          false,
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
			Keypair: []KeyPair{
				{
					PublicCertificatePath: "/home/wso2/security/truststore/mg.pem",
					PrivateKeyPath:        "/home/wso2/security/keystore/mg.key",
					UseForSigning:         true,
				},
			},
		},
		Cache: cache{
			Enabled:     true,
			MaximumSize: 10000,
			ExpiryTime:  15,
		},
		Metrics: Metrics{
			Enabled: false,
			Type:    "azure",
		},
		MandateSubscriptionValidation: false,
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
		Enabled:    false,
		Properties: map[string]string{},
		Adapter: analyticsAdapter{
			Enabled:             false,
			BufferFlushInterval: 1000000000,
			BufferSizeBytes:     16384,
			GRPCRequestTimeout:  20000000000,
		},
		Enforcer: analyticsEnforcer{
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
			Publisher: []analyticsPublisher{},
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
