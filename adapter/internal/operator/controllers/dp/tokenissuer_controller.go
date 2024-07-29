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
	"encoding/json"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/adapter/pkg/logging"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	var err error
	loggers.LoggerAPKOperator.Debugf("Reconciling jwtIssuer: %v", req.NamespacedName.String())
	jwtKey := req.NamespacedName
	jwtIssuerMapping, err := getJWTIssuers(ctx, r.client, jwtKey)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2660, logging.CRITICAL,
			"Unable to resolve JWTIssuers after updating %s : %s", req.NamespacedName.String(), err.Error()))
		return ctrl.Result{}, nil
	}
	UpdateEnforcerJWTIssuers(jwtIssuerMapping)
	return ctrl.Result{}, nil
}

// NewTokenIssuerReconciler creates a new Application controller instance.
func NewTokenIssuerReconciler(mgr manager.Manager) error {
	r := &TokenssuerReconciler{
		client: mgr.GetClient(),
	}
	ctx := context.Background()

	if err := addTokenIssuerIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2658, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	c, err := controller.New(constants.TokenIssuerController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2657, logging.BLOCKER, "Error creating TokenIssuer controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.TokenIssuer{}), &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2656, logging.BLOCKER, "Error watching TokenIssuer resources: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{}), handler.EnqueueRequestsFromMapFunc(r.populateTokenReconcileRequestsForConfigMap),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2644, logging.BLOCKER, "Error watching ConfigMap resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{}), handler.EnqueueRequestsFromMapFunc(r.populateTokenReconcileRequestsForSecret),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2645, logging.BLOCKER, "Error watching Secret resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Debug("TokenIssuer Controller successfully started. Watching TokenIssuer Objects...")
	return nil
}

func (r *TokenssuerReconciler) populateTokenReconcileRequestsForConfigMap(ctx context.Context, obj k8client.Object) []reconcile.Request {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", configMap))
		return []reconcile.Request{}
	}
	tokenIssuerList := &dpv1alpha1.TokenIssuerList{}
	err := r.client.List(ctx, tokenIssuerList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configmapIssuerIndex, utils.NamespacedName(configMap).String()),
	})
	requests := []reconcile.Request{}
	if err == nil && len(tokenIssuerList.Items) > 0 {

		for _, tokenIssuer := range tokenIssuerList.Items {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tokenIssuer.Name,
					Namespace: tokenIssuer.Namespace},
			}
			requests = append(requests, req)
			loggers.LoggerAPKOperator.Infof("Adding reconcile request for TokenIssuer: %s/%s due to configmap change: %v",
				tokenIssuer.Namespace, tokenIssuer.Name, utils.NamespacedName(configMap).String())
		}
		return requests
	}
	return requests
}

func (r *TokenssuerReconciler) populateTokenReconcileRequestsForSecret(ctx context.Context, obj k8client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", secret))
		return []reconcile.Request{}
	}
	tokenIssuerList := &dpv1alpha1.TokenIssuerList{}
	err := r.client.List(ctx, tokenIssuerList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretTokenIssuerIndex, utils.NamespacedName(secret).String()),
	})
	requests := []reconcile.Request{}
	if err == nil && len(tokenIssuerList.Items) > 0 {

		for _, tokenIssuer := range tokenIssuerList.Items {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tokenIssuer.Name,
					Namespace: tokenIssuer.Namespace},
			}
			requests = append(requests, req)
			loggers.LoggerAPKOperator.Infof("Adding reconcile request for TokenIssuer: %s/%s due to secret change: %v",
				tokenIssuer.Namespace, tokenIssuer.Name, utils.NamespacedName(secret).String())
		}
		return requests
	}
	return requests
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

// UpdateEnforcerJWTIssuers updates the JWT Issuers in the Enforcer
func UpdateEnforcerJWTIssuers(jwtIssuerMapping dpv1alpha1.JWTIssuerMapping) {
	jwtIssuerList := marshalJWTIssuerList(jwtIssuerMapping)
	xds.UpdateEnforcerJWTIssuers(jwtIssuerList)
}
func marshalJWTIssuerList(jwtIssuerMapping dpv1alpha1.JWTIssuerMapping) *subscription.JWTIssuerList {
	jwtIssuers := []*subscription.JWTIssuer{}
	for _, internalJWTIssuer := range jwtIssuerMapping {
		certificate := &subscription.Certificate{}
		jwtIssuer := &subscription.JWTIssuer{
			Name:             internalJWTIssuer.Name,
			Organization:     internalJWTIssuer.Organization,
			Issuer:           internalJWTIssuer.Issuer,
			ConsumerKeyClaim: internalJWTIssuer.ConsumerKeyClaim,
			ScopesClaim:      internalJWTIssuer.ScopesClaim,
		}
		if internalJWTIssuer.SignatureValidation.Certificate != nil && internalJWTIssuer.SignatureValidation.Certificate.ResolvedCertificate != "" {
			certificate.Certificate = internalJWTIssuer.SignatureValidation.Certificate.ResolvedCertificate
		}
		if internalJWTIssuer.SignatureValidation.JWKS != nil {
			jwks := &subscription.JWKS{}
			jwks.Url = internalJWTIssuer.SignatureValidation.JWKS.URL
			if internalJWTIssuer.SignatureValidation.JWKS.TLS != nil && internalJWTIssuer.SignatureValidation.JWKS.TLS.ResolvedCertificate != "" {
				jwks.Tls = internalJWTIssuer.SignatureValidation.JWKS.TLS.ResolvedCertificate
			}
			certificate.Jwks = jwks
		}
		jwtIssuer.ClaimMapping = internalJWTIssuer.ClaimMappings
		jwtIssuer.Certificate = certificate
		jwtIssuer.Environments = internalJWTIssuer.Environments
		jwtIssuers = append(jwtIssuers, jwtIssuer)

	}
	jwtIssuersJSON, _ := json.Marshal(jwtIssuers)
	loggers.LoggerAPKOperator.Debugf("JwtIssuer Data: %v", string(jwtIssuersJSON))
	return &subscription.JWTIssuerList{List: jwtIssuers}
}

// getJWTIssuers returns the JWTIssuers for the given JWTIssuerMapping
func getJWTIssuers(ctx context.Context, client k8client.Client, namespace types.NamespacedName) (dpv1alpha1.JWTIssuerMapping, error) {
	jwtIssuerMapping := make(dpv1alpha1.JWTIssuerMapping)
	jwtIssuerList := &dpv1alpha2.TokenIssuerList{}
	if err := client.List(ctx, jwtIssuerList); err != nil {
		return nil, err
	}
	for _, jwtIssuer := range jwtIssuerList.Items {
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
			signatureValidation.Certificate = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
		}
		resolvedJwtIssuer.SignatureValidation = signatureValidation
		jwtIssuerMappingName := types.NamespacedName{
			Name:      jwtIssuer.Name,
			Namespace: jwtIssuer.Namespace,
		}
		if jwtIssuer.Spec.ClaimMappings != nil {
			resolvedJwtIssuer.ClaimMappings = getResolvedClaimMapping(*jwtIssuer.Spec.ClaimMappings)
		} else {
			resolvedJwtIssuer.ClaimMappings = make(map[string]string)
		}
		jwtIssuerMapping[jwtIssuerMappingName] = &resolvedJwtIssuer
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
