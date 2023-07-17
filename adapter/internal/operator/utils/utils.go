/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package utils

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/wso2/apk/adapter/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	constants "github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/envutils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj k8client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

// FilterByNamespaces takes a list of namespaces and returns a filter function
// which return true if the input object is in the given namespaces list,
// and returns false otherwise
func FilterByNamespaces(namespaces []string) func(object k8client.Object) bool {
	return func(object k8client.Object) bool {
		if namespaces == nil {
			return true
		}
		return stringutils.StringInSlice(object.GetNamespace(), namespaces)
	}
}

// GetNamespace reads namespace with a default value
func GetNamespace(namespace *gwapiv1b1.Namespace, defaultNamespace string) string {
	if namespace != nil && *namespace != "" {
		return string(*namespace)
	}
	return defaultNamespace
}

// GetOperatorPodNamespace returns the namesapce of the operator pod
func GetOperatorPodNamespace() string {
	return envutils.GetEnv(constants.OperatorPodNamespace,
		constants.OperatorPodNamespaceDefaultValue)
}

// GroupPtr returns pointer to created v1beta1.Group object
func GroupPtr(name string) *gwapiv1b1.Group {
	group := gwapiv1b1.Group(name)
	return &group
}

// KindPtr returns a pointer to created v1beta1.Kind object
func KindPtr(name string) *gwapiv1b1.Kind {
	kind := gwapiv1b1.Kind(name)
	return &kind
}

// PathMatchTypePtr returns a pointer to created v1beta1.PathMatchType object
func PathMatchTypePtr(pType gwapiv1b1.PathMatchType) *gwapiv1b1.PathMatchType {
	return &pType
}

// StringPtr returns a pointer to created string
func StringPtr(val string) *string {
	return &val
}

// GetDefaultHostNameForBackend returns the host in <backend.name>.<namespace> format
func GetDefaultHostNameForBackend(backend gwapiv1b1.HTTPBackendRef,
	defaultNamespace string) string {
	return fmt.Sprintf("%s.%s", backend.Name,
		GetNamespace(backend.Namespace, defaultNamespace))
}

// TieBreaker breaks ties when multiple objects are present in a case only single object is expected.
// tie breaking logic is explained in https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#conflict-resolution
func TieBreaker[T metav1.Object](k8sObjects []T) *T {
	if len(k8sObjects) < 1 {
		return nil
	}
	selectedk8sObject := k8sObjects[0]
	for _, k8sObject := range k8sObjects[1:] {
		if selectedk8sObject.GetCreationTimestamp().After(k8sObject.GetCreationTimestamp().Time) {
			selectedk8sObject = k8sObject
		} else if selectedk8sObject.GetCreationTimestamp().String() == k8sObject.GetCreationTimestamp().String() &&
			(types.NamespacedName{
				Name:      selectedk8sObject.GetName(),
				Namespace: selectedk8sObject.GetNamespace(),
			}).String() > (types.NamespacedName{
				Name:      k8sObject.GetName(),
				Namespace: k8sObject.GetNamespace(),
			}).String() {
			selectedk8sObject = k8sObject
		}
	}
	return &selectedk8sObject
}

// SelectPolicy selects the policy based on the policy override and default values
func SelectPolicy[T any](policyUpOverride, policyUpDefault, policyDownOverride, policyDownDefault **T) *T {
	if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride,
			combineUpAndDownValues(**policyDownOverride,
				combineUpAndDownValues(**policyDownDefault, **policyUpDefault)))
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride,
			combineUpAndDownValues(**policyDownOverride, **policyDownDefault))
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride,
			combineUpAndDownValues(**policyDownOverride, **policyUpDefault))
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride,
			combineUpAndDownValues(**policyDownDefault, **policyUpDefault))
		return &output
	} else if policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyDownOverride,
			combineUpAndDownValues(**policyDownDefault, **policyUpDefault))
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride, **policyDownOverride)
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride, **policyDownDefault)
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyUpOverride, **policyUpDefault)
		return &output
	} else if policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() {
		output := combineUpAndDownValues(**policyDownOverride, **policyDownDefault)
		return &output
	} else if policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyDownOverride, **policyUpDefault)
		return &output
	} else if policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() &&
		policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		output := combineUpAndDownValues(**policyDownDefault, **policyUpDefault)
		return &output
	} else if policyUpOverride != nil && !reflect.ValueOf(*policyUpOverride).IsZero() {
		return *policyUpOverride
	} else if policyDownOverride != nil && !reflect.ValueOf(*policyDownOverride).IsZero() {
		return *policyDownOverride
	} else if policyDownDefault != nil && !reflect.ValueOf(*policyDownDefault).IsZero() {
		return *policyDownDefault
	} else if policyUpDefault != nil && !reflect.ValueOf(*policyUpDefault).IsZero() {
		return *policyUpDefault
	}
	return nil
}

// combineUpAndDownValues combines the up and down values recursively if the value is a struct
func combineUpAndDownValues[T any](up, down T) T {
	upValue := reflect.ValueOf(up)
	downValue := reflect.ValueOf(down)
	if upValue.Type() != downValue.Type() {
		panic("Inputs must be of the same type")
	}
	if upValue.Kind() != reflect.Struct {
		return up
	}
	combinedStructValue := reflect.New(upValue.Type()).Elem()
	for i := 0; i < upValue.NumField(); i++ {
		field := upValue.Type().Field(i)
		fieldName := field.Name
		upFieldValue := upValue.FieldByName(fieldName)
		downFieldValue := downValue.FieldByName(fieldName)
		var combinedFieldValue reflect.Value
		if !upFieldValue.IsZero() {
			combinedFieldValue = upFieldValue
		} else {
			combinedFieldValue = downFieldValue
		}
		if field.Type.Kind() == reflect.Struct {
			nestedCombinedFieldValue := combineUpAndDownValues(upFieldValue.Interface(), downFieldValue.Interface())
			combinedFieldValue = reflect.ValueOf(nestedCombinedFieldValue)
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if upFieldValue.IsNil() && !downFieldValue.IsZero() {
				nestedCombinedFieldValue := combineUpAndDownValues(reflect.New(field.Type.Elem()).Elem().Interface(),
					downFieldValue.Elem().Interface())
				combinedFieldValue = reflect.New(field.Type.Elem())
				combinedFieldValue.Elem().Set(reflect.ValueOf(nestedCombinedFieldValue))
			} else if downFieldValue.IsNil() {
				combinedFieldValue = upFieldValue
			} else {
				nestedCombinedFieldValue := combineUpAndDownValues(upFieldValue.Elem().Interface(),
					downFieldValue.Elem().Interface())
				combinedFieldValue = reflect.New(field.Type.Elem())
				combinedFieldValue.Elem().Set(reflect.ValueOf(nestedCombinedFieldValue))
			}
		}
		combinedStructValue.FieldByName(fieldName).Set(combinedFieldValue)
	}
	return combinedStructValue.Interface().(T)
}

// GetPtrSlice returns a slice which is also a slice containing pointers to the elements
// in the input slice.
func GetPtrSlice[T any](inputSlice []T) []*T {
	var outputSlice []*T
	for i := range inputSlice {
		outputSlice = append(outputSlice, &inputSlice[i])
	}
	return outputSlice
}

// getConfigMapValue call kubernetes client and get the configmap and key
func getConfigMapValue(ctx context.Context, client k8client.Client,
	namespace, configMapName, key string) (string, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(ctx, types.NamespacedName{
		Name:      configMapName,
		Namespace: namespace}, configMap)
	if err != nil {
		return "", err
	}
	return configMap.Data[key], nil
}

// getSecretValue call kubernetes client and get the secret and key
func getSecretValue(ctx context.Context, client k8client.Client,
	namespace, secretName, key string) (string, error) {
	secret := &corev1.Secret{}
	err := client.Get(ctx, types.NamespacedName{
		Name:      secretName,
		Namespace: namespace}, secret)
	if err != nil {
		return "", err
	}
	return string(secret.Data[key]), nil
}

// ResolveAndAddBackendToMapping resolves backend from reference and adds it to the backendMapping.
func ResolveAndAddBackendToMapping(ctx context.Context, client k8client.Client,
	backendMapping dpv1alpha1.BackendMapping,
	backendRef dpv1alpha1.BackendReference, interceptorServiceNamespace string, api *dpv1alpha1.API) {
	namespace := gwapiv1b1.Namespace(backendRef.Namespace)
	backendName := types.NamespacedName{
		Name:      backendRef.Name,
		Namespace: GetNamespace(&namespace, interceptorServiceNamespace),
	}
	backend := GetResolvedBackend(ctx, client, backendName, api)
	if backend != nil {
		backendMapping[backendName] = backend
	}
}

// ResolveRef this function will return k8client object and update owner
func ResolveRef(ctx context.Context, client k8client.Client, api *dpv1alpha1.API,
	namespacedName types.NamespacedName, isReplace bool, obj k8client.Object, opts ...k8client.GetOption) error {
	if err := client.Get(ctx, namespacedName, obj, opts...); err != nil {
		return err
	}
	if api != nil {
		err := UpdateOwnerReference(ctx, client, obj, *api, isReplace)
		return err
	}
	return nil
}

// GetResolvedBackend resolves backend TLS configurations.
func GetResolvedBackend(ctx context.Context, client k8client.Client,
	backendNamespacedName types.NamespacedName, api *dpv1alpha1.API) *dpv1alpha1.ResolvedBackend {
	resolvedBackend := dpv1alpha1.ResolvedBackend{}
	resolvedTLSConfig := dpv1alpha1.ResolvedTLSConfig{}
	var backend dpv1alpha1.Backend
	if err := ResolveRef(ctx, client, api, backendNamespacedName, false, &backend); err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2646, backendNamespacedName, err.Error()))
		}
		return nil
	}
	resolvedBackend.Services = backend.Spec.Services
	resolvedBackend.Protocol = backend.Spec.Protocol
	resolvedBackend.BasePath = backend.Spec.BasePath
	if backend.Spec.CircuitBreaker != nil {
		resolvedBackend.CircuitBreaker = &dpv1alpha1.CircuitBreaker{
			MaxConnections:     backend.Spec.CircuitBreaker.MaxConnections,
			MaxRequests:        backend.Spec.CircuitBreaker.MaxRequests,
			MaxRetries:         backend.Spec.CircuitBreaker.MaxRetries,
			MaxConnectionPools: backend.Spec.CircuitBreaker.MaxConnectionPools,
			MaxPendingRequests: backend.Spec.CircuitBreaker.MaxPendingRequests,
		}
	}
	if backend.Spec.Timeout != nil {
		resolvedBackend.Timeout = &dpv1alpha1.Timeout{
			MaxRouteTimeoutSeconds:  backend.Spec.Timeout.MaxRouteTimeoutSeconds,
			RouteTimeoutSeconds:     backend.Spec.Timeout.RouteTimeoutSeconds,
			RouteIdleTimeoutSeconds: backend.Spec.Timeout.RouteIdleTimeoutSeconds,
		}
	}
	if backend.Spec.Retry != nil {
		resolvedBackend.Retry = &dpv1alpha1.RetryConfig{
			Count:              backend.Spec.Retry.Count,
			BaseIntervalMillis: backend.Spec.Retry.BaseIntervalMillis,
			StatusCodes:        backend.Spec.Retry.StatusCodes,
		}
	}
	if backend.Spec.HealthCheck != nil {
		resolvedBackend.HealthCheck = &dpv1alpha1.HealthCheck{
			Timeout:            backend.Spec.HealthCheck.Timeout,
			Interval:           backend.Spec.HealthCheck.Interval,
			UnhealthyThreshold: backend.Spec.HealthCheck.UnhealthyThreshold,
			HealthyThreshold:   backend.Spec.HealthCheck.HealthyThreshold,
		}
	}
	var err error
	if backend.Spec.TLS != nil {
		resolvedTLSConfig.ResolvedCertificate, err = ResolveCertificate(ctx, client,
			backend.Namespace, backend.Spec.TLS.CertificateInline, backend.Spec.TLS.ConfigMapRef, backend.Spec.TLS.SecretRef)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2654, err.Error()))
			return nil
		}
		if resolvedTLSConfig.ResolvedCertificate == "" {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2654, "resolved certificate is empty"))
			return nil
		}
		resolvedTLSConfig.AllowedSANs = backend.Spec.TLS.AllowedSANs
		resolvedBackend.TLS = resolvedTLSConfig
	}
	if backend.Spec.Security != nil {
		resolvedBackend.Security = getResolvedBackendSecurity(ctx, client,
			backend.Namespace, *backend.Spec.Security)
	}
	return &resolvedBackend
}

// UpdateOwnerReference update the child with owner reference of the given parent.
func UpdateOwnerReference(ctx context.Context, client k8client.Client, child metav1.Object, api dpv1alpha1.API,
	isReplace bool) error {
	if isReplace {
		child.SetOwnerReferences([]metav1.OwnerReference{
			{
				APIVersion: api.APIVersion,
				Kind:       api.Kind,
				Name:       api.Name,
				UID:        api.UID,
			},
		})
	} else {
		child.SetOwnerReferences(append(child.GetOwnerReferences(), metav1.OwnerReference{
			APIVersion: api.APIVersion,
			Kind:       api.Kind,
			Name:       api.Name,
			UID:        api.UID,
		}))
	}
	for {
		if err := client.Update(ctx, child.(k8client.Object)); err != nil {
			if apierrors.IsInternalError(err) {
				loggers.LoggerAPKOperator.Warnf("Error while updating OwnerReferences of k8 object : %s in %s, %v",
					child.GetName(), child.GetNamespace(), err)
				time.Sleep(5 * time.Second)
			} else {
				return err
			}
		} else {
			return nil
		}
	}
}

// getResolvedBackendSecurity resolves backend security configurations.
func getResolvedBackendSecurity(ctx context.Context, client k8client.Client,
	namespace string, security dpv1alpha1.SecurityConfig) dpv1alpha1.ResolvedSecurityConfig {
	resolvedSecurity := dpv1alpha1.ResolvedSecurityConfig{}
	switch security.Type {
	case "Basic":
		var err error
		var username string
		var password string
		username, err = getSecretValue(ctx, client,
			namespace, security.Basic.SecretRef.Name, security.Basic.SecretRef.UsernameKey)
		password, err = getSecretValue(ctx, client,
			namespace, security.Basic.SecretRef.Name, security.Basic.SecretRef.PasswordKey)
		if err != nil || username == "" || password == "" {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2648, security.Basic.SecretRef))
		}
		resolvedSecurity = dpv1alpha1.ResolvedSecurityConfig{
			Type: "Basic",
			Basic: dpv1alpha1.ResolvedBasicSecurityConfig{
				Username: username,
				Password: password,
			},
		}
	}
	return resolvedSecurity
}

// ResolveCertificate reads the certificate from TLSConfig, first checks the certificateInline field,
// if no value then load the certificate from secretRef using util function called getSecretValue
func ResolveCertificate(ctx context.Context, client k8client.Client, namespace string, certificateInline *string, configMapRef *dpv1alpha1.RefConfig, secretRef *dpv1alpha1.RefConfig) (string, error) {
	var certificate string
	var err error
	if certificateInline != nil && len(*certificateInline) > 0 {
		certificate = *certificateInline
	} else if secretRef != nil {
		if certificate, err = getSecretValue(ctx, client,
			namespace, secretRef.Name, secretRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2642, secretRef))
		}
	} else if configMapRef != nil {
		if certificate, err = getConfigMapValue(ctx, client,
			namespace, configMapRef.Name, configMapRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2643, configMapRef))
		}
	}
	if err != nil {
		return "", err
	}
	if len(certificate) > 0 {
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2627))
			return "", nil
		}
		_, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2641, err.Error()))
			return "", err
		}
	}
	return certificate, nil
}

// RetrieveNamespaceListOptions retrieve namespace list options for the given namespaces
func RetrieveNamespaceListOptions(namespaces []string) k8client.ListOptions {
	var listOptions k8client.ListOptions
	if namespaces == nil {
		listOptions = k8client.ListOptions{}
	} else {
		listOptions = k8client.ListOptions{FieldSelector: fields.SelectorFromSet(fields.Set{"metadata.namespace": strings.Join(namespaces, ",")})}
	}
	return listOptions
}

// GetInterceptorService reads InterceptorService when interceptorReference is given
func GetInterceptorService(ctx context.Context, client k8client.Client,
	interceptorReference *dpv1alpha1.InterceptorReference, api *dpv1alpha1.API) *dpv1alpha1.InterceptorService {
	interceptorService := &dpv1alpha1.InterceptorService{}
	interceptorRef := types.NamespacedName{
		Namespace: interceptorReference.Namespace,
		Name:      interceptorReference.Name,
	}
	if err := ResolveRef(ctx, client, api, interceptorRef, false, interceptorService); err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2651, interceptorRef, err.Error()))
		}
	}
	return interceptorService
}
