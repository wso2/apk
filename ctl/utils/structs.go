/*
 * Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

type HTTPRouteConfig struct {
	ApiVersion    string        `yaml:"apiVersion"`
	Kind          string        `yaml:"kind"`
	MetaData      MetaData      `yaml:"metadata"`
	HttpRouteSpec HttpRouteSpec `yaml:"spec"`
}

type MetaData struct {
	Name      string            `yaml:"name,omitempty"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

type HttpRouteSpec struct {
	ParentRefs []ParentRef `yaml:"parentRefs"`
	HostNames  []string    `yaml:"hostnames"`
	Rules      []Rule      `yaml:"rules"`
}

type ParentRef struct {
	Name string `yaml:"name"`
}

type Rule struct {
	Matches     []Match      `yaml:"matches"`
	BackendRefs []BackendRef `yaml:"backendRefs,omitempty"`
}

type Match struct {
	Path        Path         `yaml:"path,omitempty"`
	Headers     []Header     `yaml:"headers,omitempty"`
	QueryParams []QueryParam `yaml:"queryParams,omitempty"`
	// BackendRefs BackendRef `yaml:"backendRefs,omitempty"`
}

type QueryParam struct {
	Type  string `yaml:"type"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Path struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type Header struct {
	Type  string `yaml:"type"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type BackendRef struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace,omitempty"`
	Port      int    `yaml:"port"`
	Group     string `yaml:"group"`
	Kind      string `yaml:"kind,omitempty"`
	Weight    string `yaml:"weight,omitempty"`
}

type ConfigMap struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	MetaData   MetaData          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

type SwaggerInfo struct {
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
}

type NetworkPolicy struct {
	ApiVersion        string            `yaml:"apiVersion"`
	Kind              string            `yaml:"kind"`
	MetaData          MetaData          `yaml:"metadata"`
	NetworkPolicySpec NetworkPolicySpec `yaml:"spec"`
}

type NetworkPolicySpec struct {
	PodSelector PodSelector `yaml:"podSelector"`
	Ingress     []Ingress   `yaml:"ingress"`
}

type PodSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type MatchLabels struct {
	Data map[string]string `yaml:"data"`
}

type Ingress struct {
	From []From `yaml:"from"`
}

type From struct {
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector"`
}

type NamespaceSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

// type ConfigMap struct {
// 	Name           string
// 	Namespace      string
// 	File           string
// 	SwaggerContent string
// }
