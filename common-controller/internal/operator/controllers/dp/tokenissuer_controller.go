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
 *  distributed under the License is distributed on an "AS IS" 	,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package dp

import (
	"context"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	tokenIssuerIndex       = "tokenIssuerIndex"
	secretTokenIssuerIndex = "secretTokenIssuerIndex"
	configmapIssuerIndex   = "configmapIssuerIndex"
	defaultAllEnvironments = "*"
)

// TokenssuerReconciler reconciles a TokenIssuer object
type TokenssuerReconciler struct {
	client k8client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=jwtissuers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=jwtissuers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=jwtissuers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JWTIssuer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TokenssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = log.FromContext(ctx)
	tokenIssuerKey := req.NamespacedName

	loggers.LoggerAPKOperator.Debugf("Reconciling tokenIssuer: %v", tokenIssuerKey.String())
	var tokenIssuer dpv1alpha2.TokenIssuer
	if err := r.client.Get(ctx, req.NamespacedName, &tokenIssuer); err != nil {
		if k8error.IsNotFound(err) {
			tokenIssuerSpec, found := r.ods.GetTokenIssuerFromStore(tokenIssuerKey)
			loggers.LoggerAPKOperator.Debugf("TokenIssuer cr not available in k8s")
			loggers.LoggerAPKOperator.Debugf("cached TokenIssuer spec: %v,%v", tokenIssuerSpec, found)
			if found {
				resolvedTokenIssuer, err := getResolvedTokenIssuer(ctx, r.client, tokenIssuerKey.Namespace, tokenIssuerSpec)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2661, logging.CRITICAL, "Error resolving tokenIssuer: %v", err.Error()))
					return ctrl.Result{}, err
				}
				utils.SendDeleteTokenIssuerEvent(*resolvedTokenIssuer)
				r.ods.DeleteTokenIssuerFromStore(tokenIssuerKey)
				server.DeleteTokenIssuer(tokenIssuerKey.String())
			} else {
				loggers.LoggerAPKOperator.Debugf("TokenIssuer %s/%s does not exist in k8s", tokenIssuerKey.Namespace, tokenIssuerKey.Name)
			}
		}
	} else {
		loggers.LoggerAPKOperator.Debugf("TokenIssuer cr available in k8s")
		oldTokenIssuerSpec, found := r.ods.GetTokenIssuerFromStore(tokenIssuerKey)
		resolvedTokenIssuer, err := getResolvedTokenIssuer(ctx, r.client, tokenIssuerKey.Namespace, tokenIssuer.Spec)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2661, logging.BLOCKER, "Error resolving tokenIssuer: %v", err.Error()))
			return ctrl.Result{}, nil
		}
		if found {
			resolvedOldTokenIssuer, err := getResolvedTokenIssuer(ctx, r.client, tokenIssuerKey.Namespace, oldTokenIssuerSpec)
			if err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2661, logging.BLOCKER, "Error resolving tokenIssuer: %v", err.Error()))
				return ctrl.Result{}, nil
			}
			// update
			loggers.LoggerAPKOperator.Debugf("TokenIssuer in ods")
			utils.SendUpdateTokenIssuerEvent(*resolvedOldTokenIssuer, *resolvedTokenIssuer)
		} else {
			loggers.LoggerAPKOperator.Debugf("TokenIssuer in ods consider as update")
			utils.SendAddTokenIssuerEvent(*resolvedTokenIssuer)
		}
		r.ods.AddorUpdateTokenIssuerToStore(tokenIssuerKey, tokenIssuer.Spec)
		r.sendTokenIssuerUpdates(tokenIssuerKey, *resolvedTokenIssuer, found)
	}
	return ctrl.Result{}, nil

}
func (r *TokenssuerReconciler) sendTokenIssuerUpdates(tokenIssuerKey types.NamespacedName, tokenIssuer dpv1alpha1.ResolvedJWTIssuer, update bool) {
	resolvedTokenIssuer := marshalTokenIssuer(tokenIssuer)
	if update {
		server.DeleteTokenIssuer(tokenIssuerKey.String())
	}
	server.AddTokenIssuer(tokenIssuerKey.String(), resolvedTokenIssuer)
}

func marshalTokenIssuer(tokenIssuer dpv1alpha1.ResolvedJWTIssuer) server.TokenIssuer {
	resolvedTokenIssuer := server.TokenIssuer{
		Name:             tokenIssuer.Name,
		Issuer:           tokenIssuer.Issuer,
		Organization:     tokenIssuer.Organization,
		ConsumerKeyClaim: tokenIssuer.ConsumerKeyClaim,
		ScopesClaim:      tokenIssuer.ScopesClaim,
		ClaimMappings:    tokenIssuer.ClaimMappings,
		Environments:     tokenIssuer.Environments,
	}
	signatureValidation := server.ResolvedSignatureValidation{}
	if tokenIssuer.SignatureValidation.JWKS != nil {
		signatureValidation.JWKS = &server.ResolvedJWKS{}
		if len(tokenIssuer.SignatureValidation.JWKS.URL) > 0 {
			signatureValidation.JWKS.URL = tokenIssuer.SignatureValidation.JWKS.URL
		}
		if tokenIssuer.SignatureValidation.JWKS.TLS != nil {
			signatureValidation.JWKS.TLS = &server.ResolvedTLSConfig{ResolvedCertificate: tokenIssuer.SignatureValidation.JWKS.TLS.ResolvedCertificate}
		}
	} else if tokenIssuer.SignatureValidation.Certificate != nil {
		signatureValidation.Certificate = &server.ResolvedTLSConfig{ResolvedCertificate: tokenIssuer.SignatureValidation.Certificate.ResolvedCertificate}
	}
	resolvedTokenIssuer.SignatureValidation = signatureValidation
	return resolvedTokenIssuer
}

// NewTokenIssuerReconciler creates a new Application controller instance.
func NewTokenIssuerReconciler(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &TokenssuerReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	ctx := context.Background()
	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.CommonController.Operator.Namespaces))}

	if err := addTokenIssuerIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2658, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	c, err := controller.New(constants.TokenIssuerReconSiler, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2657, logging.BLOCKER, "Error creating TokenIssuer controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.TokenIssuer{}), &handler.EnqueueRequestForObject{}, predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2656, logging.BLOCKER, "Error watching TokenIssuer resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("TokenIssuer Controller successfully started. Watching TokenIssuer Objects...")
	return nil
}

// addTokenIssuerIndexes adds indexers related to Gateways
func addTokenIssuerIndexes(ctx context.Context, mgr manager.Manager) error {

	// Secret to TokenIssuer indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.TokenIssuer{}, secretTokenIssuerIndex,
		func(rawObj k8client.Object) []string {
			jwtIssuer := rawObj.(*dpv1alpha1.TokenIssuer)
			var secrets []string
			if jwtIssuer.Spec.SignatureValidation.Certificate != nil && jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef != nil && len(jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			if jwtIssuer.Spec.SignatureValidation.JWKS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef != nil && len(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			return secrets
		}); err != nil {
		return err
	}
	// Configmap to TokenIssuer indexer
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.TokenIssuer{}, configmapIssuerIndex,
		func(rawObj k8client.Object) []string {
			tokenIssuer := rawObj.(*dpv1alpha1.TokenIssuer)
			var configMaps []string
			if tokenIssuer.Spec.SignatureValidation.Certificate != nil && tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef != nil && len(tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name),
						Namespace: tokenIssuer.Namespace,
					}.String())
			}
			if tokenIssuer.Spec.SignatureValidation.JWKS != nil && tokenIssuer.Spec.SignatureValidation.JWKS.TLS != nil && tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef != nil && len(tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name),
						Namespace: tokenIssuer.Namespace,
					}.String())
			}
			return configMaps
		})
	return err
}

// getResolvedTokenIssuer returns the resolved tokenIssuer
func getResolvedTokenIssuer(ctx context.Context, client k8client.Client, namespace string, jwtIssuer dpv1alpha2.TokenIssuerSpec) (*dpv1alpha1.ResolvedJWTIssuer, error) {
	resolvedJwtIssuer := dpv1alpha1.ResolvedJWTIssuer{}
	resolvedJwtIssuer.Issuer = jwtIssuer.Issuer
	resolvedJwtIssuer.ConsumerKeyClaim = jwtIssuer.ConsumerKeyClaim
	resolvedJwtIssuer.ScopesClaim = jwtIssuer.ScopesClaim
	resolvedJwtIssuer.Organization = jwtIssuer.Organization
	resolvedJwtIssuer.Environments = getTokenIssuerEnvironments(jwtIssuer.Environments)

	signatureValidation := dpv1alpha1.ResolvedSignatureValidation{}
	if jwtIssuer.SignatureValidation.JWKS != nil && len(jwtIssuer.SignatureValidation.JWKS.URL) > 0 {
		jwks := &dpv1alpha1.ResolvedJWKS{}
		jwks.URL = jwtIssuer.SignatureValidation.JWKS.URL
		if jwtIssuer.SignatureValidation.JWKS.TLS != nil {

			var tlsConfigMapRef *dpv1alpha1.RefConfig
			var tlsSecretRef *dpv1alpha1.RefConfig
			if jwtIssuer.SignatureValidation.JWKS.TLS.ConfigMapRef != nil {
				tlsConfigMapRef = utils.ConvertRefConfigsV2ToV1(jwtIssuer.SignatureValidation.JWKS.TLS.ConfigMapRef)
			}
			if jwtIssuer.SignatureValidation.JWKS.TLS.SecretRef != nil {
				tlsSecretRef = utils.ConvertRefConfigsV2ToV1(jwtIssuer.SignatureValidation.JWKS.TLS.SecretRef)
			}

			tlsCertificate, err := utils.ResolveCertificate(ctx, client, namespace, jwtIssuer.SignatureValidation.JWKS.TLS.CertificateInline, tlsConfigMapRef, tlsSecretRef)
			if err != nil || tlsCertificate == "" {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2659, logging.MAJOR, "Error resolving certificate for JWKS %v", err.Error()))
				return nil, err
			}
			jwks.TLS = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
		}
		signatureValidation.JWKS = jwks
	}
	if jwtIssuer.SignatureValidation.Certificate != nil {

		var tlsConfigMapRef *dpv1alpha1.RefConfig
		var tlsSecretRef *dpv1alpha1.RefConfig
		if jwtIssuer.SignatureValidation.Certificate.ConfigMapRef != nil {
			tlsConfigMapRef = utils.ConvertRefConfigsV2ToV1(jwtIssuer.SignatureValidation.Certificate.ConfigMapRef)
		}
		if jwtIssuer.SignatureValidation.Certificate.SecretRef != nil {
			tlsSecretRef = utils.ConvertRefConfigsV2ToV1(jwtIssuer.SignatureValidation.Certificate.SecretRef)
		}

		tlsCertificate, err := utils.ResolveCertificate(ctx, client, namespace, jwtIssuer.SignatureValidation.Certificate.CertificateInline, tlsConfigMapRef, tlsSecretRef)
		if err != nil || tlsCertificate == "" {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2659, logging.MAJOR, "Error resolving certificate for JWKS %v", err.Error()))
			return nil, err
		}
		signatureValidation.Certificate = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
	}
	resolvedJwtIssuer.SignatureValidation = signatureValidation
	if jwtIssuer.ClaimMappings != nil {
		resolvedJwtIssuer.ClaimMappings = getResolvedClaimMapping(*jwtIssuer.ClaimMappings)
	} else {
		resolvedJwtIssuer.ClaimMappings = make(map[string]string)
	}
	return &resolvedJwtIssuer, nil
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
