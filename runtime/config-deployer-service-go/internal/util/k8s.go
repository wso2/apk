/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	constantscommon "github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/config"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// GetRouteMetadataList retrieves all RouteMetadata Custom Resources from the Kubernetes cluster based on API ID and namespace.
func GetRouteMetadataList(apiID string, namespace string, k8sClient client.Client) (*v2alpha1.RouteMetadataList, error) {
	routeMetadataList := &v2alpha1.RouteMetadataList{}
	err := k8sClient.List(context.Background(), routeMetadataList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.SelectorFromSet(map[string]string{constantscommon.LabelKGWUUID: apiID}),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list RouteMetadata CRs: %w", err)
	}
	return routeMetadataList, nil
}

// GetCRsUsedByRouteMetadataNotInAPIArtifact retrieves CRs that are used by a given RouteMetadata but not present in the provided APIArtifact.
func GetCRsUsedByRouteMetadataNotInAPIArtifact(routeMetadata v2alpha1.RouteMetadata,
	apiArtifact *dto.APIArtifact, namespace string, k8sClient client.Client) (*unstructured.UnstructuredList, error) {
	routeLabels := GetFilteredLabels(routeMetadata.GetLabels())
	objectList, err := GetCRsFromLabels(routeLabels, namespace, k8sClient)
	if err != nil {
		return nil, fmt.Errorf("unable to get custom resources with labels %v: %w", routeLabels, err)
	}
	currentCRs := make(map[string]unstructured.Unstructured)
	for _, object := range objectList.Items {
		currentCRs[object.GetName()] = object
	}
	artifactCRNames := make(map[string]bool)
	for _, k8sArtifact := range apiArtifact.K8sArtifacts {
		artifactCRNames[k8sArtifact.GetName()] = true
	}
	var orphanedCRs unstructured.UnstructuredList
	for name, cr := range currentCRs {
		if !artifactCRNames[name] {
			orphanedCRs.Items = append(orphanedCRs.Items, cr)
		}
	}

	return &orphanedCRs, nil
}

// UndeployK8sRouteMetadataCR removes specific RouteMetadata CR from the Kubernetes cluster based on RouteMetadata name.
func UndeployK8sRouteMetadataCR(k8sClient client.Client, k8sRouteMetadata v2alpha1.RouteMetadata) error {
	err := k8sClient.Delete(context.Background(), &k8sRouteMetadata, &client.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete RouteMetadata CR: %w", err)
	}
	return nil
}

// GetFilteredLabels filters the RouteMetadata labels to only include API name, version, and organization.
func GetFilteredLabels(routeMetadataLabels map[string]string) map[string]string {
	filteredLabels := make(map[string]string)

	if apiName, exists := routeMetadataLabels[constantscommon.LabelKGWName]; exists {
		filteredLabels[constantscommon.LabelKGWName] = apiName
	}
	if apiVersion, exists := routeMetadataLabels[constantscommon.LabelKGWVersion]; exists {
		filteredLabels[constantscommon.LabelKGWVersion] = apiVersion
	}
	if organization, exists := routeMetadataLabels[constantscommon.LabelKGWOrganization]; exists {
		filteredLabels[constantscommon.LabelKGWOrganization] = organization
	}
	return filteredLabels
}

// GetCRsFromLabels retrieves specific custom resources that match the provided labels.
func GetCRsFromLabels(filteredLabels map[string]string, namespace string,
	k8sClient client.Client) (*unstructured.UnstructuredList, error) {

	var allObjects unstructured.UnstructuredList

	// Define the custom resource types you want to search
	resourceTypes := []schema.GroupVersionKind{
		{Group: constantscommon.WSO2KubernetesGateway, Version: "v2alpha1", Kind: constantscommon.KindRouteMetadata},
		{Group: constantscommon.WSO2KubernetesGateway, Version: "v2alpha1", Kind: constantscommon.KindRoutePolicy},
		{Group: constantscommon.K8sGroupNetworking, Version: "v1", Kind: constantscommon.KindHTTPRoute},
		{Group: constantscommon.K8sGroupNetworking, Version: "v1", Kind: constantscommon.KindGRPCRoute},
		{Group: constantscommon.K8sGroupNetworking, Version: "v1alpha3", Kind: constantscommon.KindBackendTLSPolicy},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindBackend},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindHTTPRouteFilter},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindSecurityPolicy},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindBackendTrafficPolicy},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindClientTrafficPolicy},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindEnvoyExtensionPolicy},
		{Group: constantscommon.EnvoyGateway, Version: "v1alpha1", Kind: constantscommon.KindEnvoyPatchPolicy},
		{Group: "", Version: "v1", Kind: constantscommon.KindService},
		{Group: "", Version: "v1", Kind: constantscommon.KindConfigMap},
		{Group: "", Version: "v1", Kind: constantscommon.KindSecret},
	}

	for _, gvk := range resourceTypes {
		objectList := &unstructured.UnstructuredList{}
		objectList.SetGroupVersionKind(gvk)

		err := k8sClient.List(context.Background(), objectList, &client.ListOptions{
			Namespace:     namespace,
			LabelSelector: labels.SelectorFromSet(filteredLabels),
		})
		if err != nil {
			// Continue with other resource types if this one fails
			continue
		}

		// Add matching resources to the result
		allObjects.Items = append(allObjects.Items, objectList.Items...)
	}

	return &allObjects, nil
}

// UndeployCR removes a specific custom resource from the Kubernetes cluster.
func UndeployCR(k8sClient client.Client, object unstructured.Unstructured) error {
	err := k8sClient.Delete(context.Background(), &object, &client.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete custom resource %s: %w", object.GetName(), err)
	}
	return nil
}

// ApplyK8sResource applies a Kubernetes resource to the cluster using the provided client.
func ApplyK8sResource(k8sClient client.Client, namespace string, object client.Object) error {
	// Set the namespace if it's provided and the object doesn't already have one
	if namespace != "" && object.GetNamespace() == "" {
		object.SetNamespace(namespace)
	}

	// Check if the resource already exists
	existingResource := &unstructured.Unstructured{}
	existingResource.SetGroupVersionKind(object.GetObjectKind().GroupVersionKind())
	err := k8sClient.Get(context.Background(), client.ObjectKey{
		Name:      object.GetName(),
		Namespace: object.GetNamespace(),
	}, existingResource)

	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return fmt.Errorf("failed to get existing resource: %w", err)
		}
		// Resource doesn't exist, create it
		return k8sClient.Create(context.Background(), object)
	}

	// Resource exists, update it
	// Copy resource version for proper updates
	object.SetResourceVersion(existingResource.GetResourceVersion())
	return k8sClient.Update(context.Background(), object)
}

// GetNamespace retrieves the namespace to be used for Kubernetes operations.
func GetNamespace(c *gin.Context) string {
	namespace := config.GetConfig().DefaultNamespace
	currentNamespace, err := getCurrentNamespace()
	if err == nil && currentNamespace != "" {
		namespace = currentNamespace
	}
	queryNamespace := c.Query("namespace")
	if queryNamespace != "" {
		namespace = queryNamespace
	}
	return namespace
}

// getCurrentNamespace retrieves the current namespace of the pod from the service account token file.
func getCurrentNamespace() (string, error) {
	// Try to read the namespace from the service account token file
	// This file is mounted by Kubernetes in every pod
	namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", fmt.Errorf("failed to read namespace from service account: %w", err)
	}

	namespace := strings.TrimSpace(string(namespaceBytes))
	if namespace == "" {
		return "", fmt.Errorf("namespace is empty")
	}

	return namespace, nil
}

// GeneratePolicyHash generates a SHA256 hash for any policy that implements APKOperationPolicy interface.
func GeneratePolicyHash[T model.APKOperationPolicy](policy T) string {
	// Use reflection to get the actual policy struct for consistent serialization
	policyValue := reflect.ValueOf(policy)
	if policyValue.Kind() == reflect.Ptr && policyValue.IsNil() {
		return ""
	}

	// Serialize the entire policy to JSON for consistent hashing
	jsonBytes, err := json.Marshal(policy)
	if err != nil {
		// Fallback to string representation if JSON marshaling fails
		hashInput := fmt.Sprintf("%s|%v", policy.GetPolicyName(), policy)
		hasher := sha256.New()
		hasher.Write([]byte(hashInput))
		hashBytes := hasher.Sum(nil)
		return hex.EncodeToString(hashBytes)[:16]
	}

	hasher := sha256.New()
	hasher.Write(jsonBytes)
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)[:16]
}
