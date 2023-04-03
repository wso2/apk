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
	"fmt"
	"reflect"

	constants "github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/utils/envutils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

// FilterByNamespaces takes a list of namespaces and returns a filter function
// which return true if the input object is in the given namespaces list,
// and returns false otherwise
func FilterByNamespaces(namespaces []string) func(object client.Object) bool {
	return func(object client.Object) bool {
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

// GetConfigMapValue call kubernetes client and get the configmap and key
func GetConfigMapValue(ctx context.Context, client client.Client,
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

// GetSecretValue call kubernetes client and get the secret and key
func GetSecretValue(ctx context.Context, client client.Client,
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
