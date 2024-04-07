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

	if src.Spec.SignatureValidation != nil {
		dstSignatureValidation := v1alpha2.SignatureValidation{}
		sig := *src.Spec.SignatureValidation
		if sig.JWKS != nil {
			jwks := *sig.JWKS
			jwksv2 := v1alpha2.JWKS{
				URL: jwks.URL,
			}
			if jwks.TLS != nil {
				tlsConfig := v1alpha2.CERTConfig{}
				if jwks.TLS.CertificateInline != nil {
					tlsConfig.CertificateInline = jwks.TLS.CertificateInline
				}
				if jwks.TLS.SecretRef != nil {
					tlsConfig.SecretRef = &v1alpha2.RefConfig{
						Name: jwks.TLS.SecretRef.Name,
						Key:  jwks.TLS.SecretRef.Key,
					}
				}
				if jwks.TLS.ConfigMapRef != nil {
					tlsConfig.ConfigMapRef = &v1alpha2.RefConfig{
						Name: jwks.TLS.ConfigMapRef.Name,
						Key:  jwks.TLS.ConfigMapRef.Key,
					}
				}
			}
			dstSignatureValidation.JWKS = &jwksv2
		}
		if sig.Certificate != nil {
			certificate := *sig.Certificate
			certv2 := v1alpha2.CERTConfig{
				CertificateInline: certificate.CertificateInline,
			}
			if certificate.SecretRef != nil {
				certv2.SecretRef = &v1alpha2.RefConfig{
					Name: certificate.SecretRef.Name,
					Key:  certificate.SecretRef.Key,
				}
			}
			if certificate.ConfigMapRef != nil {
				certv2.ConfigMapRef = &v1alpha2.RefConfig{
					Name: certificate.ConfigMapRef.Name,
					Key:  certificate.ConfigMapRef.Key,
				}
			}
			dstSignatureValidation.Certificate = &certv2
		}
		dst.Spec.SignatureValidation = &dstSignatureValidation
	}
	if src.Spec.ClaimMappings != nil {

		var claimMappings []v1alpha2.ClaimMapping
		for _, p := range *src.Spec.ClaimMappings {
			claimMappings = append(claimMappings, v1alpha2.ClaimMapping(p))
		}
		dst.Spec.ClaimMappings = &claimMappings
	}
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

	if dst.Spec.SignatureValidation != nil {
		dstSignatureValidation := SignatureValidation{}
		sig := *dst.Spec.SignatureValidation
		if sig.JWKS != nil {
			jwks := *sig.JWKS
			jwksv1 := JWKS{
				URL: jwks.URL,
			}
			if jwks.TLS != nil {
				tlsConfig := CERTConfig{}
				if jwks.TLS.CertificateInline != nil {
					tlsConfig.CertificateInline = jwks.TLS.CertificateInline
				}
				if jwks.TLS.SecretRef != nil {
					tlsConfig.SecretRef = &RefConfig{
						Name: jwks.TLS.SecretRef.Name,
						Key:  jwks.TLS.SecretRef.Key,
					}
				}
				if jwks.TLS.ConfigMapRef != nil {
					tlsConfig.ConfigMapRef = &RefConfig{
						Name: jwks.TLS.ConfigMapRef.Name,
						Key:  jwks.TLS.ConfigMapRef.Key,
					}
				}
			}
			dstSignatureValidation.JWKS = &jwksv1
		}
		if sig.Certificate != nil {
			certificate := *sig.Certificate
			certv1 := CERTConfig{
				CertificateInline: certificate.CertificateInline,
			}
			if certificate.SecretRef != nil {
				certv1.SecretRef = &RefConfig{
					Name: certificate.SecretRef.Name,
					Key:  certificate.SecretRef.Key,
				}
			}
			if certificate.ConfigMapRef != nil {
				certv1.ConfigMapRef = &RefConfig{
					Name: certificate.ConfigMapRef.Name,
					Key:  certificate.ConfigMapRef.Key,
				}
			}
			dstSignatureValidation.Certificate = &certv1
		}
		src.Spec.SignatureValidation = &dstSignatureValidation
	}
	if dst.Spec.ClaimMappings != nil {

		var claimMappings []ClaimMapping
		for _, p := range *dst.Spec.ClaimMappings {
			claimMappings = append(claimMappings, ClaimMapping(p))
		}
		src.Spec.ClaimMappings = &claimMappings
	}
	src.Spec.TargetRef = dst.Spec.TargetRef
	return nil
}
