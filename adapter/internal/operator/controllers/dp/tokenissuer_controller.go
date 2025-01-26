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

package dp

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// func marshalJWTIssuerList(jwtIssuerMapping dpv1alpha1.JWTIssuerMapping) *subscription.JWTIssuerList {
// 	jwtIssuers := []*subscription.JWTIssuer{}
// 	for _, internalJWTIssuer := range jwtIssuerMapping {
// 		certificate := &subscription.Certificate{}
// 		jwtIssuer := &subscription.JWTIssuer{
// 			Name:             internalJWTIssuer.Name,
// 			Organization:     internalJWTIssuer.Organization,
// 			Issuer:           internalJWTIssuer.Issuer,
// 			ConsumerKeyClaim: internalJWTIssuer.ConsumerKeyClaim,
// 			ScopesClaim:      internalJWTIssuer.ScopesClaim,
// 		}
// 		if internalJWTIssuer.SignatureValidation.Certificate != nil && internalJWTIssuer.SignatureValidation.Certificate.ResolvedCertificate != "" {
// 			certificate.Certificate = internalJWTIssuer.SignatureValidation.Certificate.ResolvedCertificate
// 		}
// 		if internalJWTIssuer.SignatureValidation.JWKS != nil {
// 			jwks := &subscription.JWKS{}
// 			jwks.Url = internalJWTIssuer.SignatureValidation.JWKS.URL
// 			if internalJWTIssuer.SignatureValidation.JWKS.TLS != nil && internalJWTIssuer.SignatureValidation.JWKS.TLS.ResolvedCertificate != "" {
// 				jwks.Tls = internalJWTIssuer.SignatureValidation.JWKS.TLS.ResolvedCertificate
// 			}
// 			certificate.Jwks = jwks
// 		}
// 		jwtIssuer.ClaimMapping = internalJWTIssuer.ClaimMappings
// 		jwtIssuer.Certificate = certificate
// 		jwtIssuer.Environments = internalJWTIssuer.Environments
// 		jwtIssuers = append(jwtIssuers, jwtIssuer)

// 	}
// 	jwtIssuersJSON, _ := json.Marshal(jwtIssuers)
// 	loggers.LoggerAPKOperator.Debugf("JwtIssuer Data: %v", string(jwtIssuersJSON))
// 	return &subscription.JWTIssuerList{List: jwtIssuers}
// }

// // GetJWTIssuers returns the JWTIssuers for the given JWTIssuerMapping
func GetJWTIssuers(ctx context.Context, client k8client.Client, gateway gwapiv1.Gateway) (map[string]*v1alpha1.ResolvedJWTIssuer, error) {
	jwtIssuerMapping := make(map[string]*v1alpha1.ResolvedJWTIssuer)
	jwtIssuerList := &dpv1alpha2.TokenIssuerList{}
	if err := client.List(ctx, jwtIssuerList); err != nil {
		return nil, err
	}
	for _, jwtIssuer := range jwtIssuerList.Items {
		if jwtIssuer.Spec.TargetRef.Kind == constants.KindGateway &&
			jwtIssuer.Spec.TargetRef.Name == v1alpha2.ObjectName(gateway.ObjectMeta.Name) &&
			jwtIssuer.Spec.TargetRef.Namespace == (*v1alpha2.Namespace)(&gateway.Namespace) {
			resolvedJwtIssuer := dpv1alpha1.ResolvedJWTIssuer{}
			resolvedJwtIssuer.Issuer = jwtIssuer.Spec.Issuer
			resolvedJwtIssuer.ConsumerKeyClaim = jwtIssuer.Spec.ConsumerKeyClaim
			resolvedJwtIssuer.ScopesClaim = jwtIssuer.Spec.ScopesClaim
			resolvedJwtIssuer.Organization = jwtIssuer.Spec.Organization
			resolvedJwtIssuer.Environments = getTokenIssuerEnvironments(jwtIssuer.Spec.Environments)

			signatureValidation := dpv1alpha1.ResolvedSignatureValidation{}
			if jwtIssuer.Spec.SignatureValidation.JWKS != nil && len(jwtIssuer.Spec.SignatureValidation.JWKS.URL) > 0 {
				jwks := &dpv1alpha1.ResolvedJWKS{}
				jwks.URL = jwtIssuer.Spec.SignatureValidation.JWKS.URL
				if jwtIssuer.Spec.SignatureValidation.JWKS.TLS != nil {
					tlsCertificate, err := utils.ResolveCertificate(ctx, client, jwtIssuer.ObjectMeta.Namespace,
						jwtIssuer.Spec.SignatureValidation.JWKS.TLS.CertificateInline,
						jwtIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef, jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef)
					if err != nil {
						loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2659, logging.MAJOR,
							"Error resolving certificate for JWKS for issuer %s in CR %s, %v", resolvedJwtIssuer.Issuer, utils.NamespacedName(&jwtIssuer).String(), err.Error()))
						continue
					}
					jwks.TLS = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
				}
				signatureValidation.JWKS = jwks
			}
			if jwtIssuer.Spec.SignatureValidation.Certificate != nil {
				tlsCertificate, err := utils.ResolveCertificate(ctx, client, jwtIssuer.ObjectMeta.Namespace,
					jwtIssuer.Spec.SignatureValidation.Certificate.CertificateInline,
					jwtIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef, jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2659, logging.MAJOR,
						"Error resolving certificate for JWKS for issuer %s in CR %s, %v", resolvedJwtIssuer.Issuer, utils.NamespacedName(&jwtIssuer).String(), err.Error()))
					continue
				}

				signatureValidation.Certificate = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: convertPemCertificatetoJWK(tlsCertificate)}
			}
			resolvedJwtIssuer.SignatureValidation = signatureValidation
			if jwtIssuer.Spec.ClaimMappings != nil {
				resolvedJwtIssuer.ClaimMappings = getResolvedClaimMapping(*jwtIssuer.Spec.ClaimMappings)
			} else {
				resolvedJwtIssuer.ClaimMappings = make(map[string]string)
			}
			jwtIssuerMappingName := strings.Join([]string{gateway.Namespace, gateway.Name}, "-")
			jwtIssuerMapping[jwtIssuerMappingName] = &resolvedJwtIssuer
		}
	}
	return jwtIssuerMapping, nil
}

func getResolvedClaimMapping(claimMappings []dpv1alpha2.ClaimMapping) map[string]string {
	resolvedClaimMappings := make(map[string]string)
	for _, claimMapping := range claimMappings {
		resolvedClaimMappings[claimMapping.RemoteClaim] = claimMapping.LocalClaim
	}
	return resolvedClaimMappings
}

func getTokenIssuerEnvironments(environments []string) []string {

	resolvedEnvironments := []string{}
	if len(environments) == 0 {
		resolvedEnvironments = append(resolvedEnvironments, defaultAllEnvironments)
	} else {
		resolvedEnvironments = environments
	}

	return resolvedEnvironments
}
func convertPemCertificatetoJWK(cert string) string {
	// Decode the PEM data
	block, _ := pem.Decode([]byte(cert))
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatalf("failed to decode PEM block containing certificate")
	}

	// Parse the certificate
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse certificate: %s", err)
	}

	// Extract the public key from the certificate
	pubKey := parsedCert.PublicKey

	// Convert the public key to a JWK
	jwkKey, err := jwk.FromRaw(pubKey)
	if err != nil {
		log.Fatalf("failed to create JWK: %s", err)
	}

	// Marshal the JWK to JSON
	jwkJSON, err := json.Marshal(jwkKey)
	if err != nil {
		log.Fatalf("failed to marshal JWK to JSON: %s", err)
	}
	return string(jwkJSON)
}
