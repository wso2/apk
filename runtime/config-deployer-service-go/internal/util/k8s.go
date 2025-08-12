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
	"fmt"
	"github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// GetRouteMetadataList retrieves all RouteMetadata Custom Resources from the Kubernetes cluster based on API ID and namespace.
func GetRouteMetadataList(apiID string, namespace string, k8sClient client.Client) (*v2alpha1.RouteMetadataList, error) {
	routeMetadataList := &v2alpha1.RouteMetadataList{}
	err := k8sClient.List(context.Background(), routeMetadataList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.SelectorFromSet(map[string]string{"apiUUID": apiID}),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list RouteMetadata CRs: %w", err)
	}
	return routeMetadataList, nil
}

// GetCRsUsedByRouteMetadataNotInAPIArtifact retrieves CRs that are currently used by the RouteMetadata
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

	if apiName, exists := routeMetadataLabels[constants.API_NAME_HASH_LABEL]; exists {
		filteredLabels[constants.API_NAME_HASH_LABEL] = apiName
	}
	if apiVersion, exists := routeMetadataLabels[constants.API_VERSION_HASH_LABEL]; exists {
		filteredLabels[constants.API_VERSION_HASH_LABEL] = apiVersion
	}
	if organization, exists := routeMetadataLabels[constants.ORGANIZATION_HASH_LABEL]; exists {
		filteredLabels[constants.ORGANIZATION_HASH_LABEL] = organization
	}
	return filteredLabels
}

// GetCRsFromLabels retrieves all custom resources in the specified namespace that match the provided labels.
func GetCRsFromLabels(filteredLabels map[string]string, namespace string,
	k8sClient client.Client) (*unstructured.UnstructuredList, error) {
	objectList := &unstructured.UnstructuredList{}
	err := k8sClient.List(context.Background(), objectList, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.SelectorFromSet(filteredLabels),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list objects with labels %v: %w", filteredLabels, err)
	}
	return objectList, nil
}

// UndeployCR removes a specific custom resource from the Kubernetes cluster.
func UndeployCR(k8sClient client.Client, object unstructured.Unstructured) error {
	err := k8sClient.Delete(context.Background(), &object, &client.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete custom resource %s: %w", object.GetName(), err)
	}
	return nil
}

// DeleteCRsByLabels deletes all custom resources in the specified namespace that match all provided labels
func DeleteCRsByLabels(namespace string, labels map[string]string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %w", err)
	}

	// Get all API resources
	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to discover API resources: %w", err)
	}

	// Create label selector from provided labels
	var labelSelectors []string
	for key, value := range labels {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", key, value))
	}
	labelSelector := strings.Join(labelSelectors, ",")

	var deletionErrors []error

	// Iterate through all API resources
	for _, apiResourceList := range apiResourceLists {
		if apiResourceList == nil {
			continue
		}

		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			// Skip subresources and resources that don't support list/delete operations
			if strings.Contains(apiResource.Name, "/") {
				continue
			}

			if !contains(apiResource.Verbs, "list") || !contains(apiResource.Verbs, "delete") {
				continue
			}

			// Skip built-in Kubernetes resources, focus on CRs
			if gv.Group == "" || isBuiltInResource(gv.Group) {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: apiResource.Name,
			}

			// List resources with label selector
			list, err := dynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: labelSelector,
			})
			if err != nil {
				// Log error but continue with other resources
				deletionErrors = append(deletionErrors, fmt.Errorf("failed to list %s: %w", gvr.String(), err))
				continue
			}

			// Delete each matching resource
			for _, item := range list.Items {
				err := dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), item.GetName(), metav1.DeleteOptions{})
				if err != nil {
					deletionErrors = append(deletionErrors, fmt.Errorf("failed to delete %s/%s: %w", gvr.String(), item.GetName(), err))
				}
			}
		}
	}

	if len(deletionErrors) > 0 {
		var errorMessages []string
		for _, err := range deletionErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return fmt.Errorf("deletion errors occurred: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isBuiltInResource checks if a group is a built-in Kubernetes resource group
func isBuiltInResource(group string) bool {
	builtInGroups := []string{
		"apps",
		"extensions",
		"networking.k8s.io",
		"rbac.authorization.k8s.io",
		"authorization.k8s.io",
		"autoscaling",
		"batch",
		"certificates.k8s.io",
		"coordination.k8s.io",
		"discovery.k8s.io",
		"events.k8s.io",
		"node.k8s.io",
		"policy",
		"scheduling.k8s.io",
		"storage.k8s.io",
		"metrics.k8s.io",
		"apiregistration.k8s.io",
		"admissionregistration.k8s.io",
	}

	for _, builtIn := range builtInGroups {
		if group == builtIn {
			return true
		}
	}
	return false
}

// ApplyK8sResource applies a Kubernetes resource to the cluster using the provided client.
func ApplyK8sResource(k8sClient client.Client, namespace string, object client.Object) error {
	//// Check if the resource already exists
	//existingResource := &unstructured.Unstructured{}
	//existingResource.SetGroupVersionKind(object.GetObjectKind().GroupVersionKind())
	//err := k8sClient.Get(context.Background(), client.ObjectKey{
	//	Name:      object.GetName(),
	//	Namespace: object.GetNamespace(),
	//}, existingResource)
	//
	//if err != nil {
	//	if client.IgnoreNotFound(err) != nil {
	//		return fmt.Errorf("failed to get existing resource: %w", err)
	//	}
	//	return k8sClient.Create(context.Background(), object)
	//}
	return k8sClient.Update(context.Background(), object)
}
