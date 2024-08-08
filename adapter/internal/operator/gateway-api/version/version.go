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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package version

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/wso2/apk/adapter/internal/operator/gateway-api/v1alpha1"
	"sigs.k8s.io/yaml"
)

type Info struct {
	EnvoyGatewayVersion    string `json:"envoyGatewayVersion"`
	GatewayAPIVersion      string `json:"gatewayAPIVersion"`
	EnvoyProxyVersion      string `json:"envoyProxyVersion"`
	EnforcerVersion        string `json:"enforcerVersion"`
	ShutdownManagerVersion string `json:"shutdownManagerVersion"`
	GitCommitID            string `json:"gitCommitID"`
}

func Get() Info {
	return Info{
		EnvoyGatewayVersion:    envoyGatewayVersion,
		GatewayAPIVersion:      gatewayAPIVersion,
		EnvoyProxyVersion:      envoyProxyVersion,
		EnforcerVersion:        enforcerVersion,
		ShutdownManagerVersion: shutdownManagerVersion,
		GitCommitID:            gitCommitID,
	}
}

var (
	envoyGatewayVersion    string
	gatewayAPIVersion      string
	envoyProxyVersion      = strings.Split(v1alpha1.DefaultEnvoyProxyImage, ":")[1]
	enforcerVersion        = strings.Split(v1alpha1.DefaultEnforcerImage, ":")[1]
	shutdownManagerVersion string
	gitCommitID            string
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == "sigs.k8s.io/gateway-api" {
				gatewayAPIVersion = dep.Version
			}
		}
	}
}

// Print shows the versions of the APK Gateway.
func Print(w io.Writer, format string) error {
	v := Get()
	switch format {
	case "json":
		if marshalled, err := json.MarshalIndent(v, "", "  "); err == nil {
			_, _ = fmt.Fprintln(w, string(marshalled))
		}
	case "yaml":
		if marshalled, err := yaml.Marshal(v); err == nil {
			_, _ = fmt.Fprintln(w, string(marshalled))
		}
	default:
		_, _ = fmt.Fprintf(w, "ENVOY_GATEWAY_VERSION: %s\n", v.EnvoyGatewayVersion)
		_, _ = fmt.Fprintf(w, "ENVOY_PROXY_VERSION: %s\n", v.EnvoyProxyVersion)
		_, _ = fmt.Fprintf(w, "GATEWAYAPI_VERSION: %s\n", v.GatewayAPIVersion)
		_, _ = fmt.Fprintf(w, "SHUTDOWN_MANAGER_VERSION: %s\n", v.ShutdownManagerVersion)
		_, _ = fmt.Fprintf(w, "GIT_COMMIT_ID: %s\n", v.GitCommitID)
	}

	return nil
}
