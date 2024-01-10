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

	// Convert Oauth2Auth to v1alpha2.Oauth2Auth
	dst.Spec.Default.AuthTypes.Oauth2 = v1alpha2.Oauth2Auth(src.Spec.Default.AuthTypes.Oauth2)
	dst.Spec.Override.AuthTypes.Oauth2 = v1alpha2.Oauth2Auth(src.Spec.Override.AuthTypes.Oauth2)

	for _, apiKeyAuth := range src.Spec.Default.AuthTypes.APIKey {
		convertedAPIKeyAuth := v1alpha2.APIKeyAuth(apiKeyAuth)
		dst.Spec.Default.AuthTypes.APIKey = append(dst.Spec.Default.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	for _, apiKeyAuth := range src.Spec.Override.AuthTypes.APIKey {
		convertedAPIKeyAuth := v1alpha2.APIKeyAuth(apiKeyAuth)
		dst.Spec.Override.AuthTypes.APIKey = append(dst.Spec.Override.AuthTypes.APIKey, convertedAPIKeyAuth)
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
	src.Spec.Default.AuthTypes.Oauth2 = Oauth2Auth(dst.Spec.Default.AuthTypes.Oauth2)
	src.Spec.Override.AuthTypes.Oauth2 = Oauth2Auth(dst.Spec.Override.AuthTypes.Oauth2)

	for _, apiKeyAuth := range dst.Spec.Default.AuthTypes.APIKey {
		convertedAPIKeyAuth := APIKeyAuth(apiKeyAuth)
		src.Spec.Default.AuthTypes.APIKey = append(src.Spec.Default.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	for _, apiKeyAuth := range dst.Spec.Override.AuthTypes.APIKey {
		convertedAPIKeyAuth := APIKeyAuth(apiKeyAuth)
		src.Spec.Override.AuthTypes.APIKey = append(src.Spec.Override.AuthTypes.APIKey, convertedAPIKeyAuth)
	}

	// Status
	src.Status = AuthenticationStatus(dst.Status)
	return nil
}
