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
	"fmt"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	jwtIssuerIndex       = "jwtIssuerIndex"
	secretJWTIssuerIndex = "secretJWTIssuerIndex"
	configmapIssuerIndex = "configmapIssuerIndex"
)

// JWTIssuerReconciler reconciles a JWTIssuer object
type JWTIssuerReconciler struct {
	client client.Client
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
func (r *JWTIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error

	loggers.LoggerAPKOperator.Debugf("Reconciling jwtIssuer: %v", req.NamespacedName.String())

	jwtKey := req.NamespacedName
	var jwtIssuerList = new(dpv1alpha1.JWTIssuerList)
	if err := r.client.List(ctx, jwtIssuerList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get jwtIssuer %s/%s",
			jwtKey.Namespace, jwtKey.Name)
	}
	jwtIssuerMapping, err := getJWTIssuers(ctx, r.client, jwtKey)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3111, err.Error()))
		return ctrl.Result{}, err
	}
	UpdateEnforcerJWTIssuers(jwtIssuerMapping)
	return ctrl.Result{}, nil
}

// NewJWTIssuerReconciler creates a new Application controller instance.
func NewJWTIssuerReconciler(mgr manager.Manager) error {
	r := &JWTIssuerReconciler{
		client: mgr.GetClient(),
	}
	ctx := context.Background()

	if err := addJWTIssuerIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3112, err))
		return err
	}
	c, err := controller.New(constants.JWTIssuerReconciler, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3111, err.Error()))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.JWTIssuer{}}, &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(utils.FilterByNamespaces([]string{utils.GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3112, err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("JWTIssuer Controller successfully started. Watching JWTIssuer Objects...")
	return nil
}

// addJWTIssuerIndexes adds indexers related to Gateways
func addJWTIssuerIndexes(ctx context.Context, mgr manager.Manager) error {

	// Secret to JWTIssuer indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.JWTIssuer{}, secretJWTIssuerIndex,
		func(rawObj k8client.Object) []string {
			jwtIssuer := rawObj.(*dpv1alpha1.JWTIssuer)
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
	// Configmap to JWTIssuer indexer
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.JWTIssuer{}, configmapIssuerIndex,
		func(rawObj k8client.Object) []string {
			jwtIssuer := rawObj.(*dpv1alpha1.JWTIssuer)
			var configMaps []string
			if jwtIssuer.Spec.SignatureValidation.Certificate != nil && jwtIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef != nil && len(jwtIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			if jwtIssuer.Spec.SignatureValidation.JWKS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef != nil && len(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			return configMaps
		})
	return err
}

// UpdateEnforcerJWTIssuers updates the JWT Issuers in the Enforcer
func UpdateEnforcerJWTIssuers(jwtIssuerMapping v1alpha1.JWTIssuerMapping) {
	jwtIssuerList := marshalJWTIssuerList(jwtIssuerMapping)
	xds.UpdateEnforcerJWTIssuers(jwtIssuerList)
}
func marshalJWTIssuerList(jwtIssuerMapping v1alpha1.JWTIssuerMapping) *subscription.JWTIssuerList {
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
		jwtIssuer.Certificate = certificate
		jwtIssuers = append(jwtIssuers, jwtIssuer)

	}
	jwtIssuersJSON, _ := json.Marshal(jwtIssuers)
	loggers.LoggerAPKOperator.Debugf("JwtIssuer Data: %v", string(jwtIssuersJSON))
	return &subscription.JWTIssuerList{List: jwtIssuers}
}

// getJWTIssuers returns the JWTIssuers for the given JWTIssuerMapping
func getJWTIssuers(ctx context.Context, client client.Client, namespace types.NamespacedName) (dpv1alpha1.JWTIssuerMapping, error) {
	jwtIssuerMapping := make(dpv1alpha1.JWTIssuerMapping)
	jwtIssuerList := &dpv1alpha1.JWTIssuerList{}
	if err := client.List(ctx, jwtIssuerList); err != nil {
		return nil, err
	}
	for _, jwtIssuer := range jwtIssuerList.Items {
		resolvedJwtIssuer := dpv1alpha1.ResolvedJWTIssuer{}
		resolvedJwtIssuer.Issuer = jwtIssuer.Spec.Issuer
		resolvedJwtIssuer.ConsumerKeyClaim = jwtIssuer.Spec.ConsumerKeyClaim
		resolvedJwtIssuer.ScopesClaim = jwtIssuer.Spec.ScopesClaim
		resolvedJwtIssuer.Organization = jwtIssuer.Spec.Organization
		signatureValidation := dpv1alpha1.ResolvedSignatureValidation{}
		if jwtIssuer.Spec.SignatureValidation.JWKS != nil && len(jwtIssuer.Spec.SignatureValidation.JWKS.URL) > 0 {
			jwks := &dpv1alpha1.ResolvedJWKS{}
			jwks.URL = jwtIssuer.Spec.SignatureValidation.JWKS.URL
			if jwtIssuer.Spec.SignatureValidation.JWKS.TLS != nil {
				tlsCertificate, err := utils.ResolveCertificate(ctx, client, jwtIssuer.ObjectMeta.Namespace, *&jwtIssuer.Spec.SignatureValidation.JWKS.TLS.CertificateInline, *&jwtIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef, *&jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef)
				if err != nil || tlsCertificate == "" {
					loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3113, err.Error()))
					continue
				}
				jwks.TLS = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
			}
			signatureValidation.JWKS = jwks
		}
		if jwtIssuer.Spec.SignatureValidation.Certificate != nil {
			tlsCertificate, err := utils.ResolveCertificate(ctx, client, jwtIssuer.ObjectMeta.Namespace, jwtIssuer.Spec.SignatureValidation.Certificate.CertificateInline, *&jwtIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef, *&jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef)
			if err != nil || tlsCertificate == "" {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3113, err.Error()))
				return nil, err
			}
			signatureValidation.Certificate = &dpv1alpha1.ResolvedTLSConfig{ResolvedCertificate: tlsCertificate}
		}
		resolvedJwtIssuer.SignatureValidation = signatureValidation
		jwtIssuerMappingName := types.NamespacedName{
			Name:      jwtIssuer.Name,
			Namespace: jwtIssuer.Namespace,
		}
		jwtIssuerMapping[jwtIssuerMappingName] = &resolvedJwtIssuer
	}
	return jwtIssuerMapping, nil
}
