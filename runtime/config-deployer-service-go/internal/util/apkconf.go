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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/google/uuid"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"net/url"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1alpha3"
	"sort"
	"strconv"
	"strings"
)

type APKConfUtil struct{}

// GetAPIName generates a unique name for the API based on its type.
func (apkConfUtil *APKConfUtil) GetAPIName(apiName string, apiType string) string {
	if strings.ToUpper(apiType) == constants.API_TYPE_GRPC {
		return apkConfUtil.GetUniqueNameForGrpcApi(apiName)
	}
	return apiName
}

// GetUniqueNameForGrpcApi generates a unique name for gRPC APIs by hashing the concatenated service names.
func (apkConfUtil *APKConfUtil) GetUniqueNameForGrpcApi(concatanatedServices string) string {
	hasher := sha1.New()
	hasher.Write([]byte(concatanatedServices))
	hashedValue := hasher.Sum(nil)
	return hex.EncodeToString(hashedValue)
}

// GetUniqueIdForAPI generates a unique identifier for an API based on its name, version, and organization.
func (apkConfUtil *APKConfUtil) GetUniqueIdForAPI(name string, version string, organization *dto.Organization) string {
	uniqueString := fmt.Sprintf("%s-%s-%s", organization.Name, name, version)
	hasher := sha1.New()
	hasher.Write([]byte(uniqueString))
	hashedValue := hasher.Sum(nil)
	return hex.EncodeToString(hashedValue)
}

// GetResourceLevelEndpointConfig extracts endpoint configurations from a list of APK operations.
func (apkConfUtil *APKConfUtil) GetResourceLevelEndpointConfig(operations []model.APKOperations) []model.EndpointConfigurations {
	endpointConfigurationsList := make([]model.EndpointConfigurations, 0)
	for _, operation := range operations {
		endpointConfigurations := operation.EndpointConfigurations
		if endpointConfigurations != nil {
			endpointConfigurationsList = append(endpointConfigurationsList, *endpointConfigurations)
		}
	}
	return endpointConfigurationsList
}

// CreateAndAddBackendServices creates backend services for the API artifact based on the provided configurations
func (apkConfUtil *APKConfUtil) CreateAndAddBackendServices(apiArtifact *model.APIArtifact, apkConf *model.APKConf,
	endpointConfigurations *model.EndpointConfigurations, apiOperation *model.APKOperations, endpointType *string,
	organization *dto.Organization) (map[string][]*dto.Endpoint, error) {
	endpointIdMap := make(map[string][]*dto.Endpoint)
	var productionEndpoints []*dto.Endpoint
	var sandboxEndpoints []*dto.Endpoint

	productionEndpointConfigs := endpointConfigurations.Production
	sandboxEndpointConfigs := endpointConfigurations.Sandbox

	// Process sandbox endpoints
	if endpointType == nil || *endpointType == constants.SANDBOX_TYPE {
		if sandboxEndpointConfigs != nil {
			for _, sandboxEndpointConfig := range sandboxEndpointConfigs {
				backendService, backendTLSPolicy, err := createBackendService(apkConf, apiOperation, constants.SANDBOX_TYPE,
					organization, &sandboxEndpointConfig)
				if err != nil {
					return nil, fmt.Errorf("failed to create sandbox backend service: %w", err)
				}
				if apiOperation == nil {
					apiArtifact.SandboxEndpointAvailable = true
				}
				apiArtifact.BackendServices[backendService.Name] = backendService
				apiArtifact.BackendTLSPolicies[backendService.Name] = backendTLSPolicy
				endpointURL := constructURLFromService(sandboxEndpointConfig.Endpoint)
				endpoint := &dto.Endpoint{
					Name:         &backendService.Name,
					ServiceEntry: false,
					URL:          &endpointURL,
					Weight:       sandboxEndpointConfig.Weight,
				}
				sandboxEndpoints = append(sandboxEndpoints, endpoint)
			}
			endpointIdMap[constants.SANDBOX_TYPE] = sandboxEndpoints
		}
	}

	// Process production endpoints
	if endpointType == nil || *endpointType == constants.PRODUCTION_TYPE {
		if productionEndpointConfigs != nil {
			for _, productionEndpointConfig := range productionEndpointConfigs {
				backendService, backendTLSPolicy, err := createBackendService(apkConf, apiOperation, constants.PRODUCTION_TYPE,
					organization, &productionEndpointConfig)
				if err != nil {
					return nil, fmt.Errorf("failed to create production backend service: %w", err)
				}
				if apiOperation == nil {
					apiArtifact.ProductionEndpointAvailable = true
				}
				apiArtifact.BackendServices[backendService.Name] = backendService
				apiArtifact.BackendTLSPolicies[backendService.Name] = backendTLSPolicy
				endpointURL := constructURLFromService(productionEndpointConfig.Endpoint)
				endpoint := &dto.Endpoint{
					Name:         &backendService.Name,
					ServiceEntry: false,
					URL:          &endpointURL,
					Weight:       productionEndpointConfig.Weight,
				}
				productionEndpoints = append(productionEndpoints, endpoint)
			}
			endpointIdMap[constants.PRODUCTION_TYPE] = productionEndpoints
		}
	}

	return endpointIdMap, nil
}

// GenerateRouteMetadata generates route metadata for the API artifact based on the APK configuration and organization.
func (apkConfUtil *APKConfUtil) GenerateRouteMetadata(apiArtifact *model.APIArtifact, apkConf *model.APKConf,
	organization *dto.Organization, cpInitiated bool) *dpv2alpha1.RouteMetadata {
	routeMetadata := &dpv2alpha1.RouteMetadata{}
	routeMetadata.Name = apiArtifact.Name + "-route-metadata"
	routeMetadata.Labels = getLabels(apkConf, organization)
	routeMetadata.Labels[constants.CP_INITIATED_HASH_LABEL] = strconv.FormatBool(cpInitiated)
	// TODO - add apiproperties and definitionfileref
	routeMetadata.Spec = dpv2alpha1.RouteMetadataSpec{
		API: dpv2alpha1.API{
			Name:           apiArtifact.Name,
			Version:        apiArtifact.Version,
			Organization:   organization.Name,
			Type:           apkConf.Type,
			Environment:    GetStringValue(apkConf.Environment),
			Context:        apkConf.BasePath,
			DefinitionPath: GetStringValue(apkConf.DefinitionPath),
		},
	}
	return routeMetadata
}

// CreateAPIResourceBundle creates a resource bundle for the API artifact based on the APK configuration and definition.
func (apkConfUtil *APKConfUtil) CreateAPIResourceBundle(apkConf *model.APKConf, organization *dto.Organization,
	cpInitiated bool, namespace string) dto.APIResourceBundle {
	apiResourceBundle := dto.APIResourceBundle{
		Organization: organization.Name,
		Namespace:    namespace,
		CPInitiated:  cpInitiated,
		APKConf:      apkConf,
	}
	// bundle apk operations into bins based on secured, rate limit and scopes and create combined resources
	combinedResources := generateCombinedResources(apkConf)
	apiResourceBundle.CombinedResources = combinedResources
	return apiResourceBundle
}

func (apkConfUtil *APKConfUtil) GenerateCRsForAPIResourceBundle(apiResourceBundle dto.APIResourceBundle) {
	// Iterate through each combined resource and generate the necessary CRs
	println("GenerateCRsForAPIResourceBundle")
}

// createBackendService creates a backend service for the API artifact based on the provided configurations.
func createBackendService(apkConf *model.APKConf, apiOperation *model.APKOperations, endpointType string,
	organization *dto.Organization, endpointConfig *model.EndpointConfiguration) (*egv1a1.Backend, *v1alpha3.BackendTLSPolicy, error) {
	backendService := &egv1a1.Backend{}
	backendService.Name = getBackendServiceUid(apkConf, apiOperation, endpointType, getHost(endpointConfig.Endpoint),
		getPath(endpointConfig.Endpoint), organization)
	backendService.Labels = getLabels(apkConf, organization)
	backendService.Spec = egv1a1.BackendSpec{
		Endpoints: []egv1a1.BackendEndpoint{
			{
				FQDN: &egv1a1.FQDNEndpoint{
					Hostname: getHost(endpointConfig.Endpoint) + "/" + getPath(endpointConfig.Endpoint),
					Port:     int32(getPort(endpointConfig.Endpoint)),
				},
			},
		},
	}
	if endpointType == constants.INTERCEPTOR_TYPE {
		backendService.Name = getInterceptorBackendUid(apkConf, endpointType, organization, endpointConfig.Endpoint)
	}

	backendTLSPolicy := &v1alpha3.BackendTLSPolicy{}
	endpointCertificate := endpointConfig.Certificate
	if endpointCertificate != nil && getProtocol(endpointConfig.Endpoint) == "https" {
		backendTLSPolicy.Name = backendService.Name + "-tls-policy"
		backendTLSPolicy.Labels = getLabels(apkConf, organization)
		backendTLSPolicy.Spec = v1alpha3.BackendTLSPolicySpec{
			TargetRefs: []v1alpha2.LocalPolicyTargetReferenceWithSectionName{
				{
					LocalPolicyTargetReference: v1alpha2.LocalPolicyTargetReference{
						Group: "gateway.envoyproxy.io",
						Kind:  "Backend",
						Name:  v1alpha2.ObjectName(backendService.Name),
					},
				},
			},
			Validation: v1alpha3.BackendTLSPolicyValidation{
				CACertificateRefs: []v1.LocalObjectReference{
					{
						Name:  v1.ObjectName(*endpointCertificate.SecretName),
						Group: "",
						Kind:  "Secret",
					},
				},
				Hostname: v1.PreciseHostname(getHost(endpointConfig.Endpoint)),
			},
		}
	}

	return backendService, backendTLSPolicy, nil
}

// constructURLFromService constructs a URL from the given endpoint, which can be either a string or a Kubernetes service model.
func constructURLFromService(endpoint interface{}) string {
	switch e := endpoint.(type) {
	case string:
		return e
	case *model.K8sService:
		return constructURLFromK8sService(e)
	default:
		return ""
	}
}

// constructURLFromK8sService constructs a URL from a Kubernetes service model.
func constructURLFromK8sService(k8sService *model.K8sService) string {
	protocol := "http"
	if k8sService.Protocol != nil {
		protocol = *k8sService.Protocol
	}

	name := ""
	if k8sService.Name != nil {
		name = *k8sService.Name
	}

	namespace := ""
	if k8sService.Namespace != nil {
		namespace = *k8sService.Namespace
	}

	port := 80
	if k8sService.Port != nil {
		port = *k8sService.Port
	}

	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d", protocol, name, namespace, port)
}

// getBackendServiceUid generates a unique identifier for backend services
func getBackendServiceUid(apkConf *model.APKConf, apiOperation *model.APKOperations, endpointType, endpointHost,
	endpointPath string, organization *dto.Organization) string {
	concatenatedString := uuid.New().String()

	if apiOperation != nil && apiOperation.EndpointConfigurations != nil {
		return "backend-" + concatenatedString + "-resource"
	} else {
		parts := []string{
			organization.Name,
			apkConf.Name,
			apkConf.Version,
			endpointType,
			endpointHost,
			endpointPath,
		}
		concatenatedString = strings.Join(parts, "-")

		// Calculate SHA1 hash
		hasher := sha1.New()
		hasher.Write([]byte(concatenatedString))
		hashedValue := hasher.Sum(nil)
		concatenatedString = hex.EncodeToString(hashedValue)

		return "backend-" + concatenatedString + "-api"
	}
}

// getHost extracts the host from a URL string or K8sService
func getHost(endpoint interface{}) string {
	var endpointURL string

	switch e := endpoint.(type) {
	case string:
		endpointURL = e
	case *model.K8sService:
		endpointURL = constructURLFromK8sService(e)
	default:
		return ""
	}

	var host string
	if strings.HasPrefix(endpointURL, "https://") {
		host = endpointURL[8:] // Remove "https://"
	} else if strings.HasPrefix(endpointURL, "http://") {
		host = endpointURL[7:] // Remove "http://"
	} else {
		return ""
	}

	// Look for port separator ":"
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		return host[:colonIndex]
	}

	// Look for path separator "/"
	if slashIndex := strings.Index(host, "/"); slashIndex != -1 {
		return host[:slashIndex]
	}

	return host
}

// getPath extracts the path from a URL string
func getPath(endpoint interface{}) string {
	var endpointURL string

	switch e := endpoint.(type) {
	case string:
		endpointURL = e
	case *model.K8sService:
		endpointURL = constructURLFromK8sService(e)
	default:
		return ""
	}

	var hostPort string
	if strings.HasPrefix(endpointURL, "https://") {
		hostPort = endpointURL[8:] // Remove "https://"
	} else if strings.HasPrefix(endpointURL, "http://") {
		hostPort = endpointURL[7:] // Remove "http://"
	} else {
		return ""
	}

	// Find the first slash which indicates the start of the path
	if slashIndex := strings.Index(hostPort, "/"); slashIndex != -1 {
		return hostPort[slashIndex:] // Return from slash to end
	}

	return ""
}

// getLabels generates labels for Kubernetes resources
func getLabels(api *model.APKConf, organization *dto.Organization) map[string]string {
	// Calculate SHA1 hash for API name
	apiNameHasher := sha1.New()
	apiNameHasher.Write([]byte(api.Name))
	apiNameHash := hex.EncodeToString(apiNameHasher.Sum(nil))

	// Calculate SHA1 hash for API version
	apiVersionHasher := sha1.New()
	apiVersionHasher.Write([]byte(api.Version))
	apiVersionHash := hex.EncodeToString(apiVersionHasher.Sum(nil))

	// Calculate SHA1 hash for organization
	organizationHasher := sha1.New()
	organizationHasher.Write([]byte(organization.Name))
	organizationHash := hex.EncodeToString(organizationHasher.Sum(nil))

	labels := map[string]string{
		constants.API_NAME_HASH_LABEL:     apiNameHash,
		constants.API_VERSION_HASH_LABEL:  apiVersionHash,
		constants.ORGANIZATION_HASH_LABEL: organizationHash,
		constants.MANAGED_BY_HASH_LABEL:   constants.MANAGED_BY_HASH_LABEL_VALUE,
	}

	return labels
}

// getPort extracts the port from a URL string or K8sService
func getPort(endpoint interface{}) int {
	var urlStr string

	switch e := endpoint.(type) {
	case string:
		urlStr = e
	case *model.K8sService:
		urlStr = constructURLFromK8sService(e)
	default:
		return -1
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return -1
	}

	port := parsedURL.Port()
	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return -1
		}
		return portNum
	}

	// Return default ports based on scheme
	switch parsedURL.Scheme {
	case "https":
		return 443
	case "http":
		return 80
	default:
		return -1
	}
}

// getProtocol determines the protocol (http or https) from the endpoint
func getProtocol(endpoint interface{}) string {
	switch e := endpoint.(type) {
	case string:
		if strings.HasPrefix(e, "https://") {
			return "https"
		} else if strings.HasPrefix(e, "http://") {
			return "http"
		}
	case *model.K8sService:
		if e.Protocol != nil {
			return *e.Protocol
		}
	}
	return "http"
}

// getInterceptorBackendUid generates a unique identifier for interceptor backend services
func getInterceptorBackendUid(apkConf *model.APKConf, endpointType string, organization *dto.Organization,
	backend interface{}) string {
	parts := []string{
		organization.Name,
		apkConf.Name,
		apkConf.Version,
		endpointType,
		constructURLFromService(backend),
	}
	concatenatedString := strings.Join(parts, "-")
	hasher := sha1.New()
	hasher.Write([]byte(concatenatedString))
	hashedValue := hasher.Sum(nil)
	concatenatedString = hex.EncodeToString(hashedValue)

	return "backend-" + concatenatedString + "-interceptor"
}

// generateCombinedResources groups APKOperations into CombinedResource buckets based on
// unique combinations of Secured, RateLimit, and Scopes attributes
func generateCombinedResources(apkConf *model.APKConf) []dto.CombinedResource {
	groupMap := make(map[string][]model.APKOperations)
	operations := apkConf.Operations
	for _, operation := range operations {
		populatedOperation := populateEndpointConfigurations(operation, apkConf)
		key := generateGroupingKey(populatedOperation)
		groupMap[key] = append(groupMap[key], populatedOperation)
	}

	var combinedResources []dto.CombinedResource
	for _, groupedOperations := range groupMap {
		combinedResource := dto.CombinedResource{
			APKOperations: groupedOperations,
		}
		combinedResources = append(combinedResources, combinedResource)
	}
	return combinedResources
}

// populateEndpointConfigurations populates missing endpoint configurations from APKConf level
func populateEndpointConfigurations(operation model.APKOperations, apkConf *model.APKConf) model.APKOperations {
	populatedOperation := operation
	// If operation doesn't have endpoint configurations but APKConf does
	if populatedOperation.EndpointConfigurations == nil {
		populatedOperation.EndpointConfigurations = apkConf.EndpointConfigurations
	} else if apkConf.EndpointConfigurations != nil {
		// If operation has endpoint configurations but some are missing, populate from APKConf
		if len(populatedOperation.EndpointConfigurations.Production) == 0 && len(apkConf.EndpointConfigurations.Production) > 0 {
			populatedOperation.EndpointConfigurations.Production = apkConf.EndpointConfigurations.Production
		}
		if len(populatedOperation.EndpointConfigurations.Sandbox) == 0 && len(apkConf.EndpointConfigurations.Sandbox) > 0 {
			populatedOperation.EndpointConfigurations.Sandbox = apkConf.EndpointConfigurations.Sandbox
		}
	}

	return populatedOperation
}

// generateGroupingKey creates a unique key for grouping operations based on
// Secured, RateLimit, and Scopes attributes
func generateGroupingKey(operation model.APKOperations) string {
	var keyParts []string

	// Handle Secured field
	if operation.Secured != nil {
		keyParts = append(keyParts, fmt.Sprintf("secured:%t", *operation.Secured))
	} else {
		keyParts = append(keyParts, "secured:false")
	}

	// Handle RateLimit field
	if operation.RateLimit != nil {
		keyParts = append(keyParts, fmt.Sprintf("rateLimit:%d-%s", operation.RateLimit.RequestsPerUnit, operation.RateLimit.Unit))
	} else {
		keyParts = append(keyParts, "rateLimit:nil")
	}

	// Handle Scopes field
	if len(operation.Scopes) > 0 {
		sortedScopes := make([]string, len(operation.Scopes))
		copy(sortedScopes, operation.Scopes)
		sort.Strings(sortedScopes)
		keyParts = append(keyParts, fmt.Sprintf("scopes:%s", strings.Join(sortedScopes, ",")))
	} else {
		keyParts = append(keyParts, "scopes:empty")
	}

	return strings.Join(keyParts, "|")
}
