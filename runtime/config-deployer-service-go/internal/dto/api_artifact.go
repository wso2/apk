/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
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

package dto

import "sigs.k8s.io/controller-runtime/pkg/client"

type APIArtifact struct {
	Name         string          `json:"name" yaml:"name"`
	Version      string          `json:"version" yaml:"version"`
	K8sArtifacts []client.Object `json:"k8sArtifacts" yaml:"k8sArtifacts"`
}
