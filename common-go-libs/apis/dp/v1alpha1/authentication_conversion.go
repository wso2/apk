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

	dst.Spec.Default.Disabled = src.Spec.Default.Disabled
	dst.Spec.Override.Disabled = src.Spec.Override.Disabled

	// Convert Oauth2Auth default to v1alpha2.Oauth2Auth : Required field added as mandatory for OAuth2
	dst.Spec.Default.AuthTypes.Oauth2 = v1alpha2.Oauth2Auth{
		Required:            "mandatory",
		Disabled:            src.Spec.Default.AuthTypes.Oauth2.Disabled,
		Header:              src.Spec.Default.AuthTypes.Oauth2.Header,
		SendTokenToUpstream: src.Spec.Default.AuthTypes.Oauth2.SendTokenToUpstream,
	}

	// Convert Oauth2Auth override to v1alpha2.Oauth2Auth : Required field added as mandatory for OAuth2
	dst.Spec.Override.AuthTypes.Oauth2 = v1alpha2.Oauth2Auth{
		Required:            "mandatory",
		Disabled:            src.Spec.Default.AuthTypes.Oauth2.Disabled,
		Header:              src.Spec.Default.AuthTypes.Oauth2.Header,
		SendTokenToUpstream: src.Spec.Default.AuthTypes.Oauth2.SendTokenToUpstream,
	}

	// Convert Oauth2Auth Default to v1alpha2.APIKey : Required field added as optional for APIKey
	for _, apiKeyAuth := range src.Spec.Default.AuthTypes.APIKey {
		convertedAPIKeyAuth := v1alpha2.APIKeyAuth{
			In:                  apiKeyAuth.In,
			Name:                apiKeyAuth.Name,
			SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
		}
		dst.Spec.Default.AuthTypes.APIKey = append(dst.Spec.Default.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	// Convert Oauth2Auth Override to v1alpha2.APIKey : Required field added as optional for APIKey
	for _, apiKeyAuth := range src.Spec.Override.AuthTypes.APIKey {
		convertedAPIKeyAuth := v1alpha2.APIKeyAuth{
			In:                  apiKeyAuth.In,
			Name:                apiKeyAuth.Name,
			SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
		}
		dst.Spec.Override.AuthTypes.APIKey = append(dst.Spec.Override.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	// Convert testConsoleKey Override to v1alpha2.JWT
	if src.Spec.Override.AuthTypes.TestConsoleKey != (TestConsoleKeyAuth{}) {
		dst.Spec.Override.AuthTypes.JWT = v1alpha2.JWT{
			Header:              src.Spec.Override.AuthTypes.TestConsoleKey.Header,
			SendTokenToUpstream: src.Spec.Override.AuthTypes.TestConsoleKey.SendTokenToUpstream,
		}
	}

	// Convert testConsoleKey Default to v1alpha2.JWT
	if src.Spec.Default.AuthTypes.TestConsoleKey != (TestConsoleKeyAuth{}) {
		dst.Spec.Default.AuthTypes.JWT = v1alpha2.JWT{
			Header:              src.Spec.Default.AuthTypes.TestConsoleKey.Header,
			SendTokenToUpstream: src.Spec.Default.AuthTypes.TestConsoleKey.SendTokenToUpstream,
		}
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

	src.Spec.Default.Disabled = dst.Spec.Default.Disabled
	src.Spec.Override.Disabled = dst.Spec.Override.Disabled
	src.Spec.Default.AuthTypes.Oauth2 = Oauth2Auth{
		Disabled:            src.Spec.Default.AuthTypes.Oauth2.Disabled,
		Header:              src.Spec.Default.AuthTypes.Oauth2.Header,
		SendTokenToUpstream: src.Spec.Default.AuthTypes.Oauth2.SendTokenToUpstream,
	}
	src.Spec.Override.AuthTypes.Oauth2 = Oauth2Auth{
		Disabled:            src.Spec.Override.AuthTypes.Oauth2.Disabled,
		Header:              src.Spec.Override.AuthTypes.Oauth2.Header,
		SendTokenToUpstream: src.Spec.Override.AuthTypes.Oauth2.SendTokenToUpstream,
	}

	for _, apiKeyAuth := range dst.Spec.Default.AuthTypes.APIKey {
		convertedAPIKeyAuth := APIKeyAuth{
			In:                  apiKeyAuth.In,
			Name:                apiKeyAuth.Name,
			SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
		}
		src.Spec.Default.AuthTypes.APIKey = append(src.Spec.Default.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	for _, apiKeyAuth := range dst.Spec.Override.AuthTypes.APIKey {
		convertedAPIKeyAuth := APIKeyAuth{
			In:                  apiKeyAuth.In,
			Name:                apiKeyAuth.Name,
			SendTokenToUpstream: apiKeyAuth.SendTokenToUpstream,
		}
		src.Spec.Override.AuthTypes.APIKey = append(src.Spec.Override.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	// Convert testConsoleKey Override to v1alpha1.TestConsoleKey
	src.Spec.Override.AuthTypes.TestConsoleKey = TestConsoleKeyAuth{
		Header:              dst.Spec.Override.AuthTypes.JWT.Header,
		SendTokenToUpstream: dst.Spec.Override.AuthTypes.JWT.SendTokenToUpstream,
	}

	// Convert testConsoleKey Default to v1alpha1.TestConsoleKey
	src.Spec.Default.AuthTypes.TestConsoleKey = TestConsoleKeyAuth{
		Header:              dst.Spec.Default.AuthTypes.JWT.Header,
		SendTokenToUpstream: dst.Spec.Default.AuthTypes.JWT.SendTokenToUpstream,
	}

	// Status
	src.Status = AuthenticationStatus(dst.Status)
	return nil
}
