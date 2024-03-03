/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	CommonController: commoncontroller{
		Server: server{
			Label: "ratelimiter",
		},
		Operator: operator{
			Namespaces: nil,
		},
		Keystore: keystore{
			KeyPath:  "/home/wso2/security/keystore/mg.key",
			CertPath: "/home/wso2/security/keystore/mg.pem",
		},
		Truststore: truststore{
			Location: "/home/wso2/security/truststore",
		},
		Environment:       "Default",
		InternalAPIServer: internalAPIServer{Port: 18003},
		ControlPlane: controlplane{
			Enabled:       false,
			Host:          "localhost",
			EventPort:     18000,
			RestPort:      18001,
			RetryInterval: 5,
			Persistence:   persistence{Type: "K8s"},
		},
		Metrics: metrics{
			Enabled: false,
			Type:    "prometheus",
			Port:    18010,
		},
	},
}
