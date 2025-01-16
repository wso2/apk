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
 */

package v1alpha1

import (
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// ConvertTo converts this API CR to the Hub version (v1alpha4).
// src is v1alpha1.API and dst is v1alpha2.API.
func (src *APIPolicy) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha4.APIPolicy)
	dst.ObjectMeta = src.ObjectMeta
	if src.Spec.Default != nil {
		var convertedSpec = v1alpha4.PolicySpec{}
		if src.Spec.Default.BackendJWTPolicy != nil {
			convertedSpec.BackendJWTPolicy = &v1alpha4.BackendJWTToken{
				Name: src.Spec.Default.BackendJWTPolicy.Name}
		}
		if src.Spec.Default.CORSPolicy != nil {
			convertedSpec.CORSPolicy = &v1alpha4.CORSPolicy{
				Enabled:                       true,
				AccessControlAllowCredentials: src.Spec.Default.CORSPolicy.AccessControlAllowCredentials,
				AccessControlAllowHeaders:     src.Spec.Default.CORSPolicy.AccessControlAllowHeaders,
				AccessControlAllowMethods:     src.Spec.Default.CORSPolicy.AccessControlAllowMethods,
				AccessControlAllowOrigins:     src.Spec.Default.CORSPolicy.AccessControlAllowOrigins,
				AccessControlExposeHeaders:    src.Spec.Default.CORSPolicy.AccessControlExposeHeaders,
				AccessControlMaxAge:           src.Spec.Default.CORSPolicy.AccessControlMaxAge}
		}
		if src.Spec.Default.RequestInterceptors != nil {
			convertedSpec.RequestInterceptors = []v1alpha4.InterceptorReference{}
			for _, interceptor := range src.Spec.Default.RequestInterceptors {
				convertedSpec.RequestInterceptors = append(convertedSpec.RequestInterceptors, v1alpha4.InterceptorReference{
					Name: interceptor.Name})
			}
		}
		if src.Spec.Default.ResponseInterceptors != nil {
			convertedSpec.ResponseInterceptors = []v1alpha4.InterceptorReference{}
			for _, interceptor := range src.Spec.Default.ResponseInterceptors {
				convertedSpec.ResponseInterceptors = append(convertedSpec.ResponseInterceptors, v1alpha4.InterceptorReference{
					Name: interceptor.Name})
			}
		}
		convertedSpec.SubscriptionValidation = false
		dst.Spec.Default = &convertedSpec
	}

	if src.Spec.Override != nil {
		var convertedSpec = v1alpha4.PolicySpec{}
		if src.Spec.Override.BackendJWTPolicy != nil {
			convertedSpec.BackendJWTPolicy = &v1alpha4.BackendJWTToken{
				Name: src.Spec.Override.BackendJWTPolicy.Name}
		}
		if src.Spec.Override.CORSPolicy != nil {
			convertedSpec.CORSPolicy = &v1alpha4.CORSPolicy{
				Enabled:                       true,
				AccessControlAllowCredentials: src.Spec.Override.CORSPolicy.AccessControlAllowCredentials,
				AccessControlAllowHeaders:     src.Spec.Override.CORSPolicy.AccessControlAllowHeaders,
				AccessControlAllowMethods:     src.Spec.Override.CORSPolicy.AccessControlAllowMethods,
				AccessControlAllowOrigins:     src.Spec.Override.CORSPolicy.AccessControlAllowOrigins,
				AccessControlExposeHeaders:    src.Spec.Override.CORSPolicy.AccessControlExposeHeaders,
				AccessControlMaxAge:           src.Spec.Override.CORSPolicy.AccessControlMaxAge}
		}
		if src.Spec.Override.RequestInterceptors != nil {
			convertedSpec.RequestInterceptors = []v1alpha4.InterceptorReference{}
			for _, interceptor := range src.Spec.Override.RequestInterceptors {
				convertedSpec.RequestInterceptors = append(convertedSpec.RequestInterceptors, v1alpha4.InterceptorReference{
					Name: interceptor.Name})
			}
		}
		if src.Spec.Override.ResponseInterceptors != nil {
			convertedSpec.ResponseInterceptors = []v1alpha4.InterceptorReference{}
			for _, interceptor := range src.Spec.Override.ResponseInterceptors {
				convertedSpec.ResponseInterceptors = append(convertedSpec.ResponseInterceptors, v1alpha4.InterceptorReference{
					Name: interceptor.Name})
			}
		}
		convertedSpec.SubscriptionValidation = false
		dst.Spec.Override = &convertedSpec
	}
	if src.Spec.TargetRef.Name != "" {
		dst.Spec.TargetRef = gwapiv1b1.NamespacedPolicyTargetReference{
			Name:  src.Spec.TargetRef.Name,
			Group: src.Spec.TargetRef.Group,
			Kind:  src.Spec.TargetRef.Kind}
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha4) to this version.
// src is v1alpha1.API and dst is v1alpha4.API.
func (src *APIPolicy) ConvertFrom(srcRaw conversion.Hub) error {

	dst := srcRaw.(*v1alpha4.APIPolicy)
	src.ObjectMeta = dst.ObjectMeta
	// Spec
	if dst.Spec.Default != nil {
		var convertedSpec = PolicySpec{}
		if dst.Spec.Default.BackendJWTPolicy != nil {
			convertedSpec.BackendJWTPolicy = &BackendJWTToken{
				Name: dst.Spec.Default.BackendJWTPolicy.Name}
		}
		if dst.Spec.Default.CORSPolicy != nil {
			convertedSpec.CORSPolicy = &CORSPolicy{
				AccessControlAllowCredentials: dst.Spec.Default.CORSPolicy.AccessControlAllowCredentials,
				AccessControlAllowHeaders:     dst.Spec.Default.CORSPolicy.AccessControlAllowHeaders,
				AccessControlAllowMethods:     dst.Spec.Default.CORSPolicy.AccessControlAllowMethods,
				AccessControlAllowOrigins:     dst.Spec.Default.CORSPolicy.AccessControlAllowOrigins,
				AccessControlExposeHeaders:    dst.Spec.Default.CORSPolicy.AccessControlExposeHeaders,
				AccessControlMaxAge:           dst.Spec.Default.CORSPolicy.AccessControlMaxAge}
		}
		if dst.Spec.Default.RequestInterceptors != nil {
			convertedSpec.RequestInterceptors = []InterceptorReference{}
			for _, interceptor := range dst.Spec.Default.RequestInterceptors {
				convertedSpec.RequestInterceptors = append(convertedSpec.RequestInterceptors, InterceptorReference{
					Name: interceptor.Name})
			}
		}
		if dst.Spec.Default.ResponseInterceptors != nil {
			convertedSpec.ResponseInterceptors = []InterceptorReference{}
			for _, interceptor := range dst.Spec.Default.ResponseInterceptors {
				convertedSpec.ResponseInterceptors = append(convertedSpec.ResponseInterceptors, InterceptorReference{
					Name: interceptor.Name})
			}
		}
		src.Spec.Default = &convertedSpec
	}
	if dst.Spec.Override != nil {
		var convertedSpec = PolicySpec{}
		if dst.Spec.Override.BackendJWTPolicy != nil {
			convertedSpec.BackendJWTPolicy = &BackendJWTToken{
				Name: dst.Spec.Override.BackendJWTPolicy.Name}
		}
		if dst.Spec.Override.CORSPolicy != nil {
			convertedSpec.CORSPolicy = &CORSPolicy{
				AccessControlAllowCredentials: dst.Spec.Override.CORSPolicy.AccessControlAllowCredentials,
				AccessControlAllowHeaders:     dst.Spec.Override.CORSPolicy.AccessControlAllowHeaders,
				AccessControlAllowMethods:     dst.Spec.Override.CORSPolicy.AccessControlAllowMethods,
				AccessControlAllowOrigins:     dst.Spec.Override.CORSPolicy.AccessControlAllowOrigins,
				AccessControlExposeHeaders:    dst.Spec.Override.CORSPolicy.AccessControlExposeHeaders,
				AccessControlMaxAge:           dst.Spec.Override.CORSPolicy.AccessControlMaxAge}
		}
		if dst.Spec.Override.RequestInterceptors != nil {
			convertedSpec.RequestInterceptors = []InterceptorReference{}
			for _, interceptor := range dst.Spec.Override.RequestInterceptors {
				convertedSpec.RequestInterceptors = append(convertedSpec.RequestInterceptors, InterceptorReference{
					Name: interceptor.Name})
			}
		}
		if dst.Spec.Override.ResponseInterceptors != nil {
			convertedSpec.ResponseInterceptors = []InterceptorReference{}
			for _, interceptor := range dst.Spec.Override.ResponseInterceptors {
				convertedSpec.ResponseInterceptors = append(convertedSpec.ResponseInterceptors, InterceptorReference{
					Name: interceptor.Name})
			}
		}
		src.Spec.Override = &convertedSpec
	}
	if dst.Spec.TargetRef.Name != "" {
		src.Spec.TargetRef = gwapiv1b1.NamespacedPolicyTargetReference{
			Name:  dst.Spec.TargetRef.Name,
			Group: dst.Spec.TargetRef.Group,
			Kind:  dst.Spec.TargetRef.Kind}
	}
	return nil
}
