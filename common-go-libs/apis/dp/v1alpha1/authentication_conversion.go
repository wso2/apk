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
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Authentication CR to the Hub version (v1alpha2).
// src is v1alpha1.Authentication and dst is v1alpha2.Authentication.
func (src *Authentication) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha2.Authentication)
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.TargetRef = src.Spec.TargetRef
	if src.Spec.Default != nil {
		defaultAuthv1Spec := src.Spec.Default
		defaultAuthenticationSpec := v1alpha2.AuthSpec{}
		defaultAuthenticationSpec.Disabled = defaultAuthv1Spec.Disabled
		if defaultAuthv1Spec.AuthTypes != nil {
			v1alpha2authTypes := v1alpha2.APIAuth{}
			v1alpha2authTypes.OAuth2 = v1alpha2.OAuth2Auth{
				Required:            "mandatory",
				Disabled:            defaultAuthv1Spec.AuthTypes.Oauth2.Disabled,
				Header:              defaultAuthv1Spec.AuthTypes.Oauth2.Header,
				SendTokenToUpstream: defaultAuthv1Spec.AuthTypes.Oauth2.SendTokenToUpstream,
			}
			var apiKeys []v1alpha2.APIKey
			// Convert APIKeyAuth Default to v1alpha2.APIKey : Required field added as optional for APIKey
			if defaultAuthv1Spec.AuthTypes.APIKey != nil {
				for _, apiKeyAuth := range defaultAuthv1Spec.AuthTypes.APIKey {
					convertedAPIKeyAuth := v1alpha2.APIKey{
						In:                  apiKeyAuth.In,
						Name:                apiKeyAuth.Name,
						SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
					}
					apiKeys = append(apiKeys, convertedAPIKeyAuth)
				}
			}

			if len(apiKeys) > 0 {
				v1alpha2authTypes.APIKey = &v1alpha2.APIKeyAuth{
					Required: "optional",
					Keys:     apiKeys,
				}
			}
			if defaultAuthv1Spec.AuthTypes.TestConsoleKey != (TestConsoleKeyAuth{}) {
				v1alpha2authTypes.JWT = v1alpha2.JWT{
					Header:              defaultAuthv1Spec.AuthTypes.TestConsoleKey.Header,
					SendTokenToUpstream: defaultAuthv1Spec.AuthTypes.TestConsoleKey.SendTokenToUpstream,
				}
			}
			defaultAuthenticationSpec.AuthTypes = &v1alpha2authTypes
		}
		dst.Spec.Default = &defaultAuthenticationSpec
	}

	if src.Spec.Override != nil {
		overrideAuthv1Spec := src.Spec.Override
		overrideAuthenticationSpec := v1alpha2.AuthSpec{}
		overrideAuthenticationSpec.Disabled = overrideAuthv1Spec.Disabled
		if overrideAuthv1Spec.AuthTypes != nil {
			v1alpha2authTypes := v1alpha2.APIAuth{}
			v1alpha2authTypes.OAuth2 = v1alpha2.OAuth2Auth{
				Required:            "mandatory",
				Disabled:            overrideAuthv1Spec.AuthTypes.Oauth2.Disabled,
				Header:              overrideAuthv1Spec.AuthTypes.Oauth2.Header,
				SendTokenToUpstream: overrideAuthv1Spec.AuthTypes.Oauth2.SendTokenToUpstream,
			}
			var apiKeys []v1alpha2.APIKey
			// Convert Oauth2Auth Override to v1alpha2.APIKey : Required field added as optional for APIKey
			if overrideAuthv1Spec.AuthTypes.APIKey != nil {
				for _, apiKeyAuth := range overrideAuthv1Spec.AuthTypes.APIKey {
					convertedAPIKeyAuth := v1alpha2.APIKey{
						In:                  apiKeyAuth.In,
						Name:                apiKeyAuth.Name,
						SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
					}
					apiKeys = append(apiKeys, convertedAPIKeyAuth)
				}
			}

			if len(apiKeys) > 0 {
				v1alpha2authTypes.APIKey = &v1alpha2.APIKeyAuth{
					Required: "optional",
					Keys:     apiKeys,
				}
			}
			if overrideAuthv1Spec.AuthTypes.TestConsoleKey != (TestConsoleKeyAuth{}) {
				v1alpha2authTypes.JWT = v1alpha2.JWT{
					Header:              overrideAuthv1Spec.AuthTypes.TestConsoleKey.Header,
					SendTokenToUpstream: overrideAuthv1Spec.AuthTypes.TestConsoleKey.SendTokenToUpstream,
				}
			}
			overrideAuthenticationSpec.AuthTypes = &v1alpha2authTypes
		}
		dst.Spec.Override = &overrideAuthenticationSpec
	}

	// Status
	dst.Status = v1alpha2.AuthenticationStatus(src.Status)
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
// src is v1alpha1.Authentication and dst is v1alpha2.Authentication.
func (src *Authentication) ConvertFrom(srcRaw conversion.Hub) error {

	dst := srcRaw.(*v1alpha2.Authentication)
	src.ObjectMeta = dst.ObjectMeta

	// Spec
	src.Spec.TargetRef = dst.Spec.TargetRef
	if dst.Spec.Default != nil {
		defaultAuthv2Spec := dst.Spec.Default
		defaultAuthenticationSpec := AuthSpec{}
		defaultAuthenticationSpec.Disabled = defaultAuthv2Spec.Disabled
		if defaultAuthv2Spec.AuthTypes != nil {
			v1alpha1authTypes := APIAuth{}
			v1alpha1authTypes.Oauth2 = Oauth2Auth{
				Disabled:            defaultAuthv2Spec.AuthTypes.OAuth2.Disabled,
				Header:              defaultAuthv2Spec.AuthTypes.OAuth2.Header,
				SendTokenToUpstream: defaultAuthv2Spec.AuthTypes.OAuth2.SendTokenToUpstream,
			}
			// Convert APIKeyAuth Default to v1alpha2.APIKey : Required field added as optional for APIKey
			if defaultAuthv2Spec.AuthTypes.APIKey != nil && defaultAuthv2Spec.AuthTypes.APIKey.Keys != nil {
				for _, apiKeyAuth := range defaultAuthv2Spec.AuthTypes.APIKey.Keys {
					convertedAPIKeyAuth := APIKeyAuth{
						In:                  apiKeyAuth.In,
						Name:                apiKeyAuth.Name,
						SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
					}
					v1alpha1authTypes.APIKey = append(v1alpha1authTypes.APIKey, convertedAPIKeyAuth)
				}
			}
			v1alpha1authTypes.TestConsoleKey = TestConsoleKeyAuth{
				Header:              defaultAuthv2Spec.AuthTypes.JWT.Header,
				SendTokenToUpstream: defaultAuthv2Spec.AuthTypes.JWT.SendTokenToUpstream,
			}
			defaultAuthenticationSpec.AuthTypes = &v1alpha1authTypes
		}
		src.Spec.Default = &defaultAuthenticationSpec
	}

	if dst.Spec.Override != nil {
		overrideAuthv2Spec := dst.Spec.Override
		overrideAuthenticationSpec := AuthSpec{}
		overrideAuthenticationSpec.Disabled = overrideAuthv2Spec.Disabled
		if overrideAuthv2Spec.AuthTypes != nil {
			v1alpha1authTypes := APIAuth{}
			v1alpha1authTypes.Oauth2 = Oauth2Auth{
				Disabled:            overrideAuthv2Spec.AuthTypes.OAuth2.Disabled,
				Header:              overrideAuthv2Spec.AuthTypes.OAuth2.Header,
				SendTokenToUpstream: overrideAuthv2Spec.AuthTypes.OAuth2.SendTokenToUpstream,
			}
			// Convert Oauth2Auth Default to v1alpha2.APIKey : Required field added as optional for APIKey
			if overrideAuthv2Spec.AuthTypes.APIKey != nil && overrideAuthv2Spec.AuthTypes.APIKey.Keys != nil {
				for _, apiKeyAuth := range overrideAuthv2Spec.AuthTypes.APIKey.Keys {
					convertedAPIKeyAuth := APIKeyAuth{
						In:                  apiKeyAuth.In,
						Name:                apiKeyAuth.Name,
						SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
					}
					v1alpha1authTypes.APIKey = append(v1alpha1authTypes.APIKey, convertedAPIKeyAuth)
				}
			}
			v1alpha1authTypes.TestConsoleKey = TestConsoleKeyAuth{
				Header:              overrideAuthv2Spec.AuthTypes.JWT.Header,
				SendTokenToUpstream: overrideAuthv2Spec.AuthTypes.JWT.SendTokenToUpstream,
			}
			overrideAuthenticationSpec.AuthTypes = &v1alpha1authTypes
		}
		src.Spec.Override = &overrideAuthenticationSpec
	}
	// Status
	src.Status = AuthenticationStatus(dst.Status)

	return nil
}
