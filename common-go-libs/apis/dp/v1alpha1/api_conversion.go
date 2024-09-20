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

package v1alpha1

import (
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this API CR to the Hub version (v1alpha3).
// src is v1alpha1.API and dst is v1alpha3.API.
func (src *API) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha3.API)
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.APIName = src.Spec.APIName
	dst.Spec.APIVersion = src.Spec.APIVersion
	dst.Spec.IsDefaultVersion = src.Spec.IsDefaultVersion
	dst.Spec.DefinitionFileRef = src.Spec.DefinitionFileRef
	dst.Spec.DefinitionPath = src.Spec.DefinitionPath
	dst.Spec.APIType = src.Spec.APIType
	dst.Spec.BasePath = src.Spec.BasePath
	dst.Spec.Organization = src.Spec.Organization
	dst.Spec.SystemAPI = src.Spec.SystemAPI

	if src.Spec.Production != nil {
		dst.Spec.Production = []v1alpha3.EnvConfig{}
		for _, productionRef := range src.Spec.Production {
			dst.Spec.Production = append(dst.Spec.Production, v1alpha3.EnvConfig{
				RouteRefs: productionRef.HTTPRouteRefs,
			})
		}
	}
	if src.Spec.Sandbox != nil {
		dst.Spec.Sandbox = []v1alpha3.EnvConfig{}
		for _, sandboxRef := range src.Spec.Sandbox {
			dst.Spec.Sandbox = append(dst.Spec.Sandbox, v1alpha3.EnvConfig{
				RouteRefs: sandboxRef.HTTPRouteRefs,
			})
		}
	}

	// Convert []Property to []v1alpha2.Property
	var properties []v1alpha3.Property
	for _, p := range src.Spec.APIProperties {
		properties = append(properties, v1alpha3.Property(p))
	}
	dst.Spec.APIProperties = properties

	// Status
	dst.Status.DeploymentStatus = v1alpha3.DeploymentStatus(src.Status.DeploymentStatus)

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
// src is v1alpha1.API and dst is v1alpha3.API.
func (src *API) ConvertFrom(srcRaw conversion.Hub) error {

	dst := srcRaw.(*v1alpha3.API)
	src.ObjectMeta = dst.ObjectMeta

	// Spec
	src.Spec.APIName = dst.Spec.APIName
	src.Spec.APIVersion = dst.Spec.APIVersion
	src.Spec.IsDefaultVersion = dst.Spec.IsDefaultVersion
	src.Spec.DefinitionFileRef = dst.Spec.DefinitionFileRef
	src.Spec.DefinitionPath = dst.Spec.DefinitionPath
	src.Spec.APIType = dst.Spec.APIType
	src.Spec.BasePath = dst.Spec.BasePath
	src.Spec.Organization = dst.Spec.Organization
	src.Spec.SystemAPI = dst.Spec.SystemAPI

	if dst.Spec.Production != nil {
		src.Spec.Production = []EnvConfig{}
		for _, productionRef := range dst.Spec.Production {
			src.Spec.Production = append(src.Spec.Production, EnvConfig{
				HTTPRouteRefs: productionRef.RouteRefs,
			})
		}
	}

	if dst.Spec.Sandbox != nil {
		src.Spec.Sandbox = []EnvConfig{}
		for _, sandboxRef := range dst.Spec.Sandbox {
			src.Spec.Sandbox = append(src.Spec.Sandbox, EnvConfig{
				HTTPRouteRefs: sandboxRef.RouteRefs,
			})
		}
	}

	// Convert []Property to []v1alpha1.Property
	var properties []Property
	for _, p := range dst.Spec.APIProperties {
		properties = append(properties, Property(p))
	}
	src.Spec.APIProperties = properties

	// Status
	src.Status.DeploymentStatus = DeploymentStatus(dst.Status.DeploymentStatus)

	return nil
}
