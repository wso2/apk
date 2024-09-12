/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	constants "github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/envutils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
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
func GetNamespace(namespace *gwapiv1.Namespace, defaultNamespace string) string {
	if namespace != nil && *namespace != "" {
		return string(*namespace)
	}
	return defaultNamespace
}

// ValidateAndRetrieveNamespace checks if the child resource's namespace is the same as the parent resource's namespace
func ValidateAndRetrieveNamespace(namespace *gwapiv1.Namespace, defaultNamespace string) (string, error) {
	if namespace != nil && *namespace != "" {
		if string(*namespace) == defaultNamespace {
			return string(*namespace), nil
		}
		return "", errors.New("Namespace mismatch")
	}
	return defaultNamespace, nil
}

// GetOperatorPodNamespace returns the namesapce of the operator pod
func GetOperatorPodNamespace() string {
	return envutils.GetEnv(constants.OperatorPodNamespace,
		constants.OperatorPodNamespaceDefaultValue)
}

// GroupPtr returns pointer to created v1beta1.Group object
func GroupPtr(name string) *gwapiv1.Group {
	group := gwapiv1.Group(name)
	return &group
}

// KindPtr returns a pointer to created v1beta1.Kind object
func KindPtr(name string) *gwapiv1.Kind {
	kind := gwapiv1.Kind(name)
	return &kind
}

// PathMatchTypePtr returns a pointer to created v1beta1.PathMatchType object
func PathMatchTypePtr(pType gwapiv1.PathMatchType) *gwapiv1.PathMatchType {
	return &pType
}

// StringPtr returns a pointer to created string
func StringPtr(val string) *string {
	return &val
}

// GetDefaultHostNameForBackend returns the host in <backend.name>.<namespace> format
func GetDefaultHostNameForBackend(backend gwapiv1.HTTPBackendRef,
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

// GetService retrieves the Service object and returns its details.
func GetService(ctx context.Context, client k8client.Client, namespace, serviceName string) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := client.Get(ctx, types.NamespacedName{
		Name:      serviceName,
		Namespace: namespace,
	}, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// GetResolvedBackendFromService converts a Kubernetes Service to a Resolved Backend.
func GetResolvedBackendFromService(k8sService *corev1.Service, svcPort int) (*dpv1alpha2.ResolvedBackend, error) {

	var host string
	var port uint32

	if len(k8sService.Spec.Ports) == 0 {
		port = uint32(svcPort)
	} else {
		servicePort := k8sService.Spec.Ports[0]
		port = uint32(servicePort.Port)
	}

	switch k8sService.Spec.Type {
	case corev1.ServiceTypeClusterIP, corev1.ServiceTypeNodePort:
		// Use the internal DNS name for clusterip and nodeport
		host = fmt.Sprintf("%s.%s.svc.cluster.local", k8sService.Name, k8sService.Namespace)
	case corev1.ServiceTypeLoadBalancer:
		// Use the external IP or hostname for LB services
		if len(k8sService.Status.LoadBalancer.Ingress) > 0 {
			ingress := k8sService.Status.LoadBalancer.Ingress[0]
			if ingress.IP != "" {
				host = ingress.IP
			} else if ingress.Hostname != "" {
				host = ingress.Hostname
			} else {
				return nil, fmt.Errorf("no valid ingress found for LoadBalancer service %s", k8sService.Name)
			}
		} else {
			return nil, fmt.Errorf("no load balancer ingress found for service %s", k8sService.Name)
		}
	default:
		return nil, fmt.Errorf("unsupported service type %s", k8sService.Spec.Type)
	}

	backend := &dpv1alpha2.ResolvedBackend{Services: []dpv1alpha2.Service{{Host: host, Port: port}}, Protocol: dpv1alpha2.HTTPProtocol}
	return backend, nil
}

// ResolveAndAddBackendToMapping resolves backend from reference and adds it to the backendMapping.
func ResolveAndAddBackendToMapping(ctx context.Context, client k8client.Client,
	backendMapping map[string]*dpv1alpha2.ResolvedBackend,
	backendRef dpv1alpha1.BackendReference, interceptorServiceNamespace string, api *dpv1alpha3.API) {
	backendName := types.NamespacedName{
		Name:      backendRef.Name,
		Namespace: interceptorServiceNamespace,
	}
	backend := GetResolvedBackend(ctx, client, backendName, api)
	if backend != nil {
		backendMapping[backendName.String()] = backend
	}
}

// ResolveRef this function will return k8client object and update owner
func ResolveRef(ctx context.Context, client k8client.Client, api *dpv1alpha3.API,
	namespacedName types.NamespacedName, isReplace bool, obj k8client.Object, opts ...k8client.GetOption) error {
	err := client.Get(ctx, namespacedName, obj, opts...)
	return err
}

// GetResolvedBackend resolves backend TLS configurations.
func GetResolvedBackend(ctx context.Context, client k8client.Client,
	backendNamespacedName types.NamespacedName, api *dpv1alpha3.API) *dpv1alpha2.ResolvedBackend {
	resolvedBackend := dpv1alpha2.ResolvedBackend{}
	resolvedTLSConfig := dpv1alpha2.ResolvedTLSConfig{}
	var backend dpv1alpha2.Backend
	if err := ResolveRef(ctx, client, api, backendNamespacedName, false, &backend); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2646, logging.CRITICAL, "Error while getting backend: %v, error: %v", backendNamespacedName, err.Error()))
		return nil
	}
	resolvedBackend.Backend = backend
	resolvedBackend.Services = backend.Spec.Services
	resolvedBackend.Protocol = backend.Spec.Protocol
	resolvedBackend.BasePath = backend.Spec.BasePath
	if backend.Spec.CircuitBreaker != nil {
		resolvedBackend.CircuitBreaker = &dpv1alpha2.CircuitBreaker{
			MaxConnections:     backend.Spec.CircuitBreaker.MaxConnections,
			MaxRequests:        backend.Spec.CircuitBreaker.MaxRequests,
			MaxRetries:         backend.Spec.CircuitBreaker.MaxRetries,
			MaxConnectionPools: backend.Spec.CircuitBreaker.MaxConnectionPools,
			MaxPendingRequests: backend.Spec.CircuitBreaker.MaxPendingRequests,
		}
	}
	if backend.Spec.Timeout != nil {
		resolvedBackend.Timeout = &dpv1alpha2.Timeout{
			UpstreamResponseTimeout:      backend.Spec.Timeout.UpstreamResponseTimeout,
			DownstreamRequestIdleTimeout: backend.Spec.Timeout.DownstreamRequestIdleTimeout,
		}
	}
	if backend.Spec.Retry != nil {
		resolvedBackend.Retry = &dpv1alpha2.RetryConfig{
			Count:              backend.Spec.Retry.Count,
			BaseIntervalMillis: backend.Spec.Retry.BaseIntervalMillis,
			StatusCodes:        backend.Spec.Retry.StatusCodes,
		}
	}
	if backend.Spec.HealthCheck != nil {
		resolvedBackend.HealthCheck = &dpv1alpha2.HealthCheck{
			Timeout:            backend.Spec.HealthCheck.Timeout,
			Interval:           backend.Spec.HealthCheck.Interval,
			UnhealthyThreshold: backend.Spec.HealthCheck.UnhealthyThreshold,
			HealthyThreshold:   backend.Spec.HealthCheck.HealthyThreshold,
		}
	}
	var err error
	if backend.Spec.TLS != nil {
		resolvedTLSConfig.ResolvedCertificate, err = ResolveCertificate(ctx, client,
			backend.Namespace, backend.Spec.TLS.CertificateInline, ConvertRefConfigsV1ToV2(backend.Spec.TLS.ConfigMapRef), ConvertRefConfigsV1ToV2(backend.Spec.TLS.SecretRef))
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2654, logging.CRITICAL, "Error resolving certificate for Backend %v", err.Error()))
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

// UpdateCR updates the given CR.
// use to update owner reference of the given CR.
func UpdateCR(ctx context.Context, client k8client.Client, child metav1.Object) error {
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
	namespace string, security dpv1alpha2.SecurityConfig) dpv1alpha2.ResolvedSecurityConfig {
	resolvedSecurity := dpv1alpha2.ResolvedSecurityConfig{}
	if security.Basic != nil {
		var err error
		var username string
		var password string
		username, err = getSecretValue(ctx, client,
			namespace, security.Basic.SecretRef.Name, security.Basic.SecretRef.UsernameKey)
		password, err = getSecretValue(ctx, client,
			namespace, security.Basic.SecretRef.Name, security.Basic.SecretRef.PasswordKey)
		if err != nil || username == "" || password == "" {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2648, logging.CRITICAL, "Error while reading key from secretRef: %s", security.Basic.SecretRef))
		}
		resolvedSecurity = dpv1alpha2.ResolvedSecurityConfig{
			Type: "Basic",
			Basic: dpv1alpha2.ResolvedBasicSecurityConfig{
				Username: username,
				Password: password,
			},
		}
	} else if security.APIKey != nil {
		var err error
		var in string
		var keyName string
		var keyValue string
		in = security.APIKey.In
		keyName = security.APIKey.Name
		if security.APIKey.ValueFrom.Name != "" {
			keyValue, err = getSecretValue(ctx, client,
				namespace, security.APIKey.ValueFrom.Name, security.APIKey.ValueFrom.ValueKey)
			if err != nil || keyValue == "" {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2649, logging.CRITICAL, "Error while reading key from secretRef: %s", security.APIKey.ValueFrom))
			}
		} else {
			keyValue = security.APIKey.ValueFrom.ValueKey
		}
		resolvedSecurity = dpv1alpha2.ResolvedSecurityConfig{
			Type: "APIKey",
			APIKey: dpv1alpha2.ResolvedAPIKeySecurityConfig{
				In:    in,
				Name:  keyName,
				Value: keyValue,
			},
		}
	}
	loggers.LoggerAPKOperator.Debugf("Resolved Security %v", resolvedSecurity)
	return resolvedSecurity
}

// GetResolvedMutualSSL resolves mTLS related security configurations.
func GetResolvedMutualSSL(ctx context.Context, client k8client.Client, authentication dpv1alpha2.Authentication, resolvedMutualSSL *dpv1alpha2.MutualSSL) error {
	var mutualSSL *dpv1alpha2.MutualSSLConfig
	authSpec := SelectPolicy(&authentication.Spec.Override, &authentication.Spec.Default, nil, nil)
	if authSpec.AuthTypes != nil {
		mutualSSL = authSpec.AuthTypes.MutualSSL
	}

	if mutualSSL != nil {
		resolvedCertificates, err := ResolveAllmTLSCertificates(ctx, mutualSSL, client, authentication.Namespace)
		resolvedMutualSSL.Disabled = mutualSSL.Disabled
		resolvedMutualSSL.Required = mutualSSL.Required
		resolvedMutualSSL.ClientCertificates = append(resolvedMutualSSL.ClientCertificates, resolvedCertificates...)

		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Error in resolving mutual SSL %v in authentication", mutualSSL))
			return err
		}
	}
	return nil
}

// ResolveAllmTLSCertificates resolves all mTLS certificates
func ResolveAllmTLSCertificates(ctx context.Context, mutualSSL *dpv1alpha2.MutualSSLConfig, client k8client.Client, namespace string) ([]string, error) {
	var resolvedCertificates []string
	var err error
	var certificate string
	if mutualSSL.CertificatesInline != nil {
		for _, cert := range mutualSSL.CertificatesInline {
			certificate, err = ResolveCertificate(ctx, client, namespace, cert, nil, nil)
			resolvedCertificates = append(resolvedCertificates, certificate)
		}
	} else if mutualSSL.ConfigMapRefs != nil {
		for _, cert := range mutualSSL.ConfigMapRefs {
			certificate, err = ResolveCertificate(ctx, client, namespace, nil, cert, nil)
			resolvedCertificates = append(resolvedCertificates, certificate)
		}
	} else if mutualSSL.SecretRefs != nil {
		for _, cert := range mutualSSL.SecretRefs {
			certificate, err = ResolveCertificate(ctx, client, namespace, nil, nil, cert)
			resolvedCertificates = append(resolvedCertificates, certificate)
		}
	}
	return resolvedCertificates, err
}

// ResolveCertificate reads the certificate from TLSConfig, first checks the certificateInline field,
// if no value then load the certificate from secretRef using util function called getSecretValue
func ResolveCertificate(ctx context.Context, client k8client.Client, namespace string, certificateInline *string,
	configMapRef *dpv1alpha2.RefConfig, secretRef *dpv1alpha2.RefConfig) (string, error) {
	var certificate string
	var err error
	if certificateInline != nil && len(*certificateInline) > 0 {
		certificate = *certificateInline
	} else if secretRef != nil {
		if certificate, err = getSecretValue(ctx, client,
			namespace, secretRef.Name, secretRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2642, logging.CRITICAL,
				"Error while reading certificate from secretRef %s: %s", secretRef, err.Error()))
			return "", err
		}
	} else if configMapRef != nil {
		if certificate, err = getConfigMapValue(ctx, client,
			namespace, configMapRef.Name, configMapRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2643, logging.CRITICAL,
				"Error while reading certificate from configMapRef %s : %s", configMapRef, err.Error()))
			return "", err
		}
	}
	if len(certificate) > 0 {
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2627, logging.CRITICAL, "Failed to decode certificate PEM."))
			return "", fmt.Errorf("failed to decode certificate PEM")
		}
		_, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2641, logging.CRITICAL, "Error while parsing certificate: %s", err.Error()))
			return "", fmt.Errorf("error while parsing certificate: %s", err.Error())
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
func GetInterceptorService(ctx context.Context, client k8client.Client, namespace string,
	interceptorReference *dpv1alpha3.InterceptorReference, api *dpv1alpha3.API) *dpv1alpha1.InterceptorService {
	interceptorService := &dpv1alpha1.InterceptorService{}
	interceptorRef := types.NamespacedName{
		Namespace: namespace,
		Name:      interceptorReference.Name,
	}
	if err := ResolveRef(ctx, client, api, interceptorRef, false, interceptorService); err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2651, logging.CRITICAL, "Error while getting interceptor service %s, error: %v", interceptorRef, err.Error()))
		}
	}
	return interceptorService
}

// GetBackendJWT reads BackendJWT when backendJWTReference is given
func GetBackendJWT(ctx context.Context, client k8client.Client, namespace,
	backendJWTReference string, api *dpv1alpha3.API) *dpv1alpha1.BackendJWT {
	backendJWT := &dpv1alpha1.BackendJWT{}
	backendJWTRef := types.NamespacedName{
		Namespace: namespace,
		Name:      backendJWTReference,
	}
	if err := ResolveRef(ctx, client, api, backendJWTRef, false, backendJWT); err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2662, logging.CRITICAL, "Error while getting backendjwt %s, error: %v", backendJWTRef, err.Error()))
		}
	}
	return backendJWT
}

// GetAIProvider reads AIProvider when aiProviderReference is given
func GetAIProvider(ctx context.Context, client k8client.Client, namespace string,
	aiProviderReference string, api *dpv1alpha3.API) *dpv1alpha3.AIProvider {
	aiProvider := &dpv1alpha3.AIProvider{}
	aiProviderRef := types.NamespacedName{
		Namespace: namespace,
		Name:      aiProviderReference,
	}
	loggers.LoggerAPKOperator.Debugf("AIProviderRef: %v", aiProviderRef)
	if err := ResolveRef(ctx, client, api, aiProviderRef, false, aiProvider); err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2663, logging.CRITICAL, "Error while getting aiProvider %s, error: %v", aiProviderRef, err.Error()))
		}
	}
	loggers.LoggerAPKOperator.Debugf("AIProvider: %v", aiProvider)
	return aiProvider
}

// RetrieveAPIList retrieves API list from the given kubernetes client
func RetrieveAPIList(k8sclient k8client.Client) ([]dpv1alpha3.API, error) {
	ctx := context.Background()
	conf := config.ReadConfigs()
	namespaces := conf.Adapter.Operator.Namespaces
	var apis []dpv1alpha3.API
	if namespaces == nil {
		apiList := &dpv1alpha3.APIList{}
		if err := k8sclient.List(ctx, apiList, &k8client.ListOptions{}); err != nil {
			return nil, err
		}
		apis = make([]dpv1alpha3.API, len(apiList.Items))
		copy(apis[:], apiList.Items[:])
	} else {
		for _, namespace := range namespaces {
			apiList := &dpv1alpha3.APIList{}
			if err := k8sclient.List(ctx, apiList, &k8client.ListOptions{Namespace: namespace}); err != nil {
				return nil, err
			}
			apis = append(apis, apiList.Items...)
		}
	}
	return apis, nil
}

// ConvertRefConfigsV1ToV2 converts RefConfig v2 to v1
func ConvertRefConfigsV1ToV2(refConfig *dpv1alpha2.RefConfig) *dpv1alpha2.RefConfig {
	if refConfig != nil {
		return &dpv1alpha2.RefConfig{
			Name: refConfig.Name,
			Key:  refConfig.Key,
		}
	}
	return nil
}

// ContainsString checks whether a list contains a specific string.
// It returns true if the string is found in the list, otherwise false.
func ContainsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// GetSubscriptionToAPIIndexID returns the id which can be used to list subscriptions related to a api.
func GetSubscriptionToAPIIndexID(name string, version string) string {
	return fmt.Sprintf("%s_%s", name, version)
}
