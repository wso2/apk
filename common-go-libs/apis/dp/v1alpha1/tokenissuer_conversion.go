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
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this TokenIssuer CR to the Hub version (v1alpha2).
// src is v1alpha1.TokenIssuer and dst is v1alpha2.TokenIssuer.
func (src *TokenIssuer) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha2.TokenIssuer)
	dst.ObjectMeta = src.ObjectMeta

	// Spec
	dst.Spec.Name = src.Spec.Name
	dst.Spec.Organization = src.Spec.Organization
	dst.Spec.Issuer = src.Spec.Issuer
	dst.Spec.ConsumerKeyClaim = src.Spec.ConsumerKeyClaim
	dst.Spec.ScopesClaim = src.Spec.ScopesClaim

	sig := *src.Spec.SignatureValidation
	jwks := *sig.JWKS
	certificate := *sig.Certificate

	jwksv2 := v1alpha2.JWKS{
		URL: jwks.URL,
		TLS: &v1alpha2.CERTConfig{
			CertificateInline: jwks.TLS.CertificateInline,
			SecretRef: &v1alpha2.RefConfig{
				Name: jwks.TLS.SecretRef.Name,
				Key:  jwks.TLS.SecretRef.Key,
			},
			ConfigMapRef: &v1alpha2.RefConfig{
				Name: jwks.TLS.ConfigMapRef.Name,
				Key:  jwks.TLS.ConfigMapRef.Key,
			},
		},
	}

	certv2 := v1alpha2.CERTConfig{
		CertificateInline: certificate.CertificateInline,
		SecretRef: &v1alpha2.RefConfig{
			Name: certificate.SecretRef.Name,
			Key:  certificate.SecretRef.Key,
		},
		ConfigMapRef: &v1alpha2.RefConfig{
			Name: certificate.ConfigMapRef.Name,
			Key:  certificate.ConfigMapRef.Key,
		},
	}

	dst.Spec.SignatureValidation = &v1alpha2.SignatureValidation{
		JWKS:        &jwksv2,
		Certificate: &certv2,
	}

	var claimMappings []v1alpha2.ClaimMapping
	for _, p := range *src.Spec.ClaimMappings {
		claimMappings = append(claimMappings, v1alpha2.ClaimMapping(p))
	}
	dst.Spec.ClaimMappings = &claimMappings

	dst.Spec.TargetRef = src.Spec.TargetRef
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
// src is v1alpha1.TokenIssuer and dst is v1alpha2.TokenIssuer.
func (src *TokenIssuer) ConvertFrom(srcRaw conversion.Hub) error {

	dst := srcRaw.(*v1alpha2.TokenIssuer)
	src.ObjectMeta = dst.ObjectMeta

	// Spec
	src.Spec.Name = dst.Spec.Name
	src.Spec.Organization = dst.Spec.Organization
	src.Spec.Issuer = dst.Spec.Issuer
	src.Spec.ConsumerKeyClaim = dst.Spec.ConsumerKeyClaim
	src.Spec.ScopesClaim = dst.Spec.ScopesClaim

	sig := *dst.Spec.SignatureValidation
	jwks := *sig.JWKS
	certificate := *sig.Certificate

	jwksv1 := JWKS{
		URL: jwks.URL,
		TLS: &CERTConfig{
			CertificateInline: jwks.TLS.CertificateInline,
			SecretRef: &RefConfig{
				Name: jwks.TLS.SecretRef.Name,
				Key:  jwks.TLS.SecretRef.Key,
			},
			ConfigMapRef: &RefConfig{
				Name: jwks.TLS.ConfigMapRef.Name,
				Key:  jwks.TLS.ConfigMapRef.Key,
			},
		},
	}

	certv1 := CERTConfig{
		CertificateInline: certificate.CertificateInline,
		SecretRef: &RefConfig{
			Name: certificate.SecretRef.Name,
			Key:  certificate.SecretRef.Key,
		},
		ConfigMapRef: &RefConfig{
			Name: certificate.ConfigMapRef.Name,
			Key:  certificate.ConfigMapRef.Key,
		},
	}

	src.Spec.SignatureValidation = &SignatureValidation{
		JWKS:        &jwksv1,
		Certificate: &certv1,
	}

	var claimMappings []ClaimMapping
	for _, p := range *dst.Spec.ClaimMappings {
		claimMappings = append(claimMappings, ClaimMapping(p))
	}
	src.Spec.ClaimMappings = &claimMappings

	src.Spec.TargetRef = dst.Spec.TargetRef

	return nil
}
