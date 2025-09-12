package crbuilder

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	eg "github.com/envoyproxy/gateway/api/v1alpha1"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	constantscommon "github.com/wso2/apk/common-go-libs/constants"
	gqlCommon "github.com/wso2/apk/common-go-libs/graphql"
	"github.com/wso2/apk/common-go-libs/pkg/logging"
	utilscommon "github.com/wso2/apk/common-go-libs/utils"
	"github.com/wso2/apk/config-deployer-service-go/internal/config"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

var allowedMethods = map[string]gatewayv1.HTTPMethod{
	"GET":     gatewayv1.HTTPMethodGet,
	"POST":    gatewayv1.HTTPMethodPost,
	"PUT":     gatewayv1.HTTPMethodPut,
	"DELETE":  gatewayv1.HTTPMethodDelete,
	"PATCH":   gatewayv1.HTTPMethodPatch,
	"HEAD":    gatewayv1.HTTPMethodHead,
	"OPTIONS": gatewayv1.HTTPMethodOptions,
	"TRACE":   gatewayv1.HTTPMethodTrace,
}

var logger logging.Logger

func init() {
	logger = config.GetConfig().Logger
}

// extractBackendBasePath extracts the base path from the backend endpoint of an operation based on the environment.
func extractBackendBasePath(operation model.APKOperations, environment string) (string, error) {
	endpoint, err := getEndpointForEnvironment(operation, environment)
	if err != nil {
		return "", err
	}
	return extractPathFromEndpoint(endpoint)
}

// extractBackendHostname extracts the hostname from the backend endpoint of an operation based on the environment.
func extractBackendHostname(operation model.APKOperations, environment string) (string, error) {
	endpoint, err := getEndpointForEnvironment(operation, environment)
	if err != nil {
		return "", err
	}
	return extractHostFromEndpoint(endpoint)
}

// getEndpointForEnvironment retrieves the endpoint for the specified environment.
func getEndpointForEnvironment(operation model.APKOperations, environment string) (interface{}, error) {
	if environment == constants.SANDBOX_TYPE && operation.EndpointConfigurations.Sandbox != nil &&
		len(operation.EndpointConfigurations.Sandbox) > 0 {
		return operation.EndpointConfigurations.Sandbox[0].Endpoint, nil
	} else if environment == constants.PRODUCTION_TYPE && operation.EndpointConfigurations.Production != nil &&
		len(operation.EndpointConfigurations.Production) > 0 {
		return operation.EndpointConfigurations.Production[0].Endpoint, nil
	}
	return nil, fmt.Errorf("no valid endpoint configurations found for environment: %s", environment)
}

// extractPathFromEndpoint extracts the path from an endpoint interface.
func extractPathFromEndpoint(endpoint interface{}) (string, error) {
	typeofEndpoint, _ := endpointType(endpoint)
	if typeofEndpoint == "string" {
		if parsed, err := url.Parse(endpoint.(string)); err == nil {
			return parsed.Path, nil
		} else {
			return "", err
		}
	}
	return "", nil
}

// extractHostFromEndpoint extracts the hostname from an endpoint interface.
func extractHostFromEndpoint(endpoint interface{}) (string, error) {
	typeofEndpoint, _ := endpointType(endpoint)
	if typeofEndpoint == "string" {
		if parsed, err := url.Parse(endpoint.(string)); err == nil {
			// Extract hostname without port
			if strings.Contains(parsed.Host, ":") {
				hostname, _, err := net.SplitHostPort(parsed.Host)
				if err != nil {
					return "", err
				}
				return hostname, nil
			}
			return parsed.Host, nil
		} else {
			return "", err
		}
	}
	return "", nil
}

func safeHTTPMethod(verb *string) *gatewayv1.HTTPMethod {
	if verb == nil {
		return nil
	}
	if m, ok := allowedMethods[strings.ToUpper(*verb)]; ok {
		return &m
	}
	// fallback: skip invalid methods
	fallback := gatewayv1.HTTPMethodGet
	return &fallback
}

func safeGRPCMethod(verb *string, service *string) *gatewayv1.GRPCMethodMatch {
	if verb == nil {
		return nil
	}
	return &gatewayv1.GRPCMethodMatch{
		Method:  verb,
		Service: service,
	}
}

// CreateResources creates Kubernetes resources based on the provided APIResourceBundle.
func CreateResources(apiResourceBundle *dto.APIResourceBundle) ([]client.Object, error) {
	if apiResourceBundle == nil || apiResourceBundle.APKConf == nil || apiResourceBundle.APKConf.EndpointConfigurations == nil {
		return nil, fmt.Errorf("invalid APIResourceBundle")
	}
	var err error
	objects := make([]client.Object, 0)
	routeMetadataList := createRouteMetadataList(apiResourceBundle, constants.PRODUCTION_TYPE)
	// Production
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Production) > 0 {
		objectsP, err := createResourcesForEnvironment(apiResourceBundle, constants.PRODUCTION_TYPE, routeMetadataList)
		if err != nil {
			return nil, err
		}
		var routeMetadata *dpv2alpha1.RouteMetadata
		httpRouteNames := make([]string, 0)
		for _, object := range objectsP {
			if routeMeta, ok := object.(*dpv2alpha1.RouteMetadata); ok {
				routeMetadata = routeMeta
			} else {
				if httpRoute, ok := object.(*gatewayv1.HTTPRoute); ok {
					httpRouteNames = append(httpRouteNames, httpRoute.Name)
				}
				objects = append(objects, object)
			}
		}
		httpRouteAnnotations := generateHTTPRouteAnnotations(httpRouteNames)
		if routeMetadata != nil {
			routeMetadata.SetAnnotations(httpRouteAnnotations)
			objects = append(objects, routeMetadata)
		}
	}
	// Sandbox
	if len(apiResourceBundle.APKConf.EndpointConfigurations.Sandbox) > 0 {
		objectsS, err := createResourcesForEnvironment(apiResourceBundle, constants.SANDBOX_TYPE, routeMetadataList)
		if err != nil {
			return nil, err
		}
		var routeMetadata *dpv2alpha1.RouteMetadata
		httpRouteNames := make([]string, 0)
		for _, object := range objectsS {
			if routeMeta, ok := object.(*dpv2alpha1.RouteMetadata); ok {
				routeMetadata = routeMeta
			} else {
				if httpRoute, ok := object.(*gatewayv1.HTTPRoute); ok {
					httpRouteNames = append(httpRouteNames, httpRoute.Name)
				}
				objects = append(objects, object)
			}
		}
		httpRouteAnnotations := generateHTTPRouteAnnotations(httpRouteNames)
		if routeMetadata != nil {
			routeMetadata.SetAnnotations(httpRouteAnnotations)
			objects = append(objects, routeMetadata)
		}
	}

	// Collect HTTPRoute names for the APIResourceBundle
	httpRouteNames := extractHTTPRouteNames(objects)
	httpRouteAnnotations := generateHTTPRouteAnnotations(httpRouteNames)

	// Convert and append routeMetadataList
	for _, rm := range routeMetadataList {
		rm.SetAnnotations(httpRouteAnnotations)
		objects = append(objects, client.Object(rm))
	}

	if apiResourceBundle.Namespace != "" {
		for _, object := range objects {
			object.SetNamespace(apiResourceBundle.Namespace)
		}
	}
	labels := make(map[string]string)
	labels[constantscommon.LabelKGWName] = util.SanitizeOrHashLabel(apiResourceBundle.APKConf.Name)
	labels[constantscommon.LabelKGWVersion] = util.SanitizeOrHashLabel(apiResourceBundle.APKConf.Version)
	labels[constantscommon.LabelKGWOrganization] = util.SanitizeOrHashLabel(apiResourceBundle.Organization)
	labels[constantscommon.LabelKGWUUID] = util.SanitizeOrHashLabel(apiResourceBundle.APKConf.ID)
	labels[constantscommon.LabelKGWCPInitiated] = util.SanitizeOrHashLabel(strconv.FormatBool(apiResourceBundle.CPInitiated))

	for _, object := range objects {
		object.SetLabels(labels)
	}
	return objects, err

}

// createRouteMetadataList creates RouteMetadata objects for the given APIResourceBundle and environment.
func createRouteMetadataList(apiResourceBundle *dto.APIResourceBundle, environment string) []*dpv2alpha1.RouteMetadata {
	definitionCMName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
		apiResourceBundle.Organization)
	routeMetadataList := make([]*dpv2alpha1.RouteMetadata, 0)
	routeMetadata := &dpv2alpha1.RouteMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Name: util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
				apiResourceBundle.Organization),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRouteMetadata,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RouteMetadataSpec{
			API: dpv2alpha1.API{
				Name:         apiResourceBundle.APKConf.Name,
				Version:      apiResourceBundle.APKConf.Version,
				Organization: apiResourceBundle.Organization,
				Environment: func() string {
					if apiResourceBundle.APKConf.Environment != nil {
						return *apiResourceBundle.APKConf.Environment
					}
					return ""
				}(),
				Context: apiResourceBundle.APKConf.BasePath,
				DefinitionPath: func() string {
					if apiResourceBundle.APKConf.DefinitionPath != nil {
						return *apiResourceBundle.APKConf.DefinitionPath
					}
					return "definition"
				}(),
				UUID: apiResourceBundle.APKConf.ID,
				DefinitionFileRef: &gwapiv1a2.LocalObjectReference{
					Name: gatewayv1.ObjectName(definitionCMName),
					Kind: gatewayv1.Kind(constantscommon.KindConfigMap),
				},
			},
		},
	}
	routeMetadataList = append(routeMetadataList, routeMetadata)
	return routeMetadataList
}

// createResourcesForEnvironment creates the necessary Kubernetes resources for a given environment
func createResourcesForEnvironment(apiResourceBundle *dto.APIResourceBundle, environment string,
	routeMetadataList []*dpv2alpha1.RouteMetadata) ([]client.Object, error) {
	objects := make([]client.Object, 0)
	definitionCMName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
		apiResourceBundle.Organization)
	routePolicies := make(map[string]*dpv2alpha1.RoutePolicy)

	if apiResourceBundle.Definition != "" {
		cm := createConfigMapForDefinition(
			definitionCMName,
			apiResourceBundle.Definition,
		)
		objects = append(objects, cm)
	}

	// RoutePolicy
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
				apiResourceBundle.Organization),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  make([]*dpv2alpha1.Mediation, 0),
			ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
		},
	}

	if apiResourceBundle.APKConf.SubscriptionValidation {
		routePolicy.Spec.RequestMediation = append(routePolicy.Spec.RequestMediation, &dpv2alpha1.Mediation{
			PolicyName:    constantscommon.MediationSubscriptionValidation,
			PolicyID:      "",
			PolicyVersion: "",
			Parameters: []*dpv2alpha1.Parameter{
				{
					Key:   "Enabled",
					Value: "true",
				},
			},
		})
	}

	// AIProvider
	if apiResourceBundle.APKConf.AIProvider != nil {
		aiProviderRoutePolicy := &dpv2alpha1.RoutePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name: utilscommon.CreateAIProviderName(apiResourceBundle.APKConf.AIProvider.Name, apiResourceBundle.APKConf.AIProvider.APIVersion),
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       constantscommon.KindRoutePolicy,
				APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
			},
		}
		routePolicies[constants.AIProviderRoutePolicy] = aiProviderRoutePolicy
		// Do not add to objects as it will be added globally by either agent or manually
	}

	// GraphQL
	if apiResourceBundle.APKConf.Type == constants.API_TYPE_GRAPHQL {
		gqlSchemaConfigMapName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment,
			apiResourceBundle.APKConf.Version, apiResourceBundle.Organization) + "-graphql-schema"
		routePolicy.Spec.RequestMediation = append(routePolicy.Spec.RequestMediation, &dpv2alpha1.Mediation{
			PolicyName:    constantscommon.MediationGraphQL,
			PolicyID:      "",
			PolicyVersion: "",
			Parameters: []*dpv2alpha1.Parameter{
				{
					Key: constantscommon.GraphQLPolicyKeySchema,
					ValueRef: &gwapiv1a2.LocalObjectReference{
						Name: gwapiv1a2.ObjectName(gqlSchemaConfigMapName),
						Kind: constantscommon.KindConfigMap,
					},
				},
			},
		})
		gqlOperations := make([]*gqlCommon.Operation, 0)
		for _, operation := range apiResourceBundle.APKConf.Operations {
			gqlOperations = append(gqlOperations, &gqlCommon.Operation{
				Target: *operation.Target,
				Verb:   *operation.Verb,
				Scopes: operation.Scopes,
			})
		}
		yamlBytes, err := yaml.Marshal(gqlOperations)
		if err != nil {
			logger.Sugar().Errorf("Error occurred while marshalling GraphQL operations to YAML: %v", err)
		}
		yamlString := string(yamlBytes)
		cm := createConfigMapForGQlSchema(gqlSchemaConfigMapName, yamlString)
		objects = append(objects, cm)
	}
	routePolicies[constants.BaseRoutePolicy] = routePolicy
	objects = append(objects, routePolicy)
	routes := make(map[int][]client.Object)
	btpByName := make(map[string]*eg.BackendTrafficPolicy)

	if apiResourceBundle.APKConf.Type == constants.API_TYPE_GRPC {
		// Create the GRPCRoute objects
		routesGRPC, objectsForWithVersion := GenerateGRPCRoutes(apiResourceBundle, true, environment, routePolicies, routeMetadataList)
		for _, obj := range objectsForWithVersion {
			objects = append(objects, obj)
		}
		if apiResourceBundle.APKConf.DefaultVersion {
			routesL, objectsForWithoutVersion := GenerateGRPCRoutes(apiResourceBundle, false, environment, routePolicies, routeMetadataList)
			for _, obj := range objectsForWithoutVersion {
				objects = append(objects, obj)
			}
			for key, value := range routesL {
				routesGRPC[key] = append(value, routesGRPC[key]...)
			}
		}
		// Convert map[int][]gatewayv1.GRPCRoute to map[int][]client.Object
		for key, grpcRoutes := range routesGRPC {
			for _, grpcRoute := range grpcRoutes {
				routes[key] = append(routes[key], client.Object(&grpcRoute))
			}
		}
	} else {
		// Create the HTTPRoute objects
		routesHTTP, objectsForWithVersion := GenerateHTTPRoutes(apiResourceBundle, true, environment,
			routePolicies, routeMetadataList, apiResourceBundle.APKConf.KeyManagers)
		for _, obj := range objectsForWithVersion {
			objects = append(objects, obj)
		}
		if apiResourceBundle.APKConf.DefaultVersion {
			routesL, objectsForWithoutVersion := GenerateHTTPRoutes(apiResourceBundle, false, environment,
				routePolicies, routeMetadataList, apiResourceBundle.APKConf.KeyManagers)
			for _, obj := range objectsForWithoutVersion {
				objects = append(objects, obj)
			}
			for key, value := range routesL {
				routesHTTP[key] = append(value, routesHTTP[key]...)
			}
		}
		// Convert map[int][]gatewayv1.HTTPRoute to map[int][]client.Object
		for key, httpRoutes := range routesHTTP {
			for _, httpRoute := range httpRoutes {
				routes[key] = append(routes[key], client.Object(&httpRoute))
			}
		}
	}

	if apiResourceBundle.APKConf.Type == constants.API_TYPE_GRPC {
		var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
		for _, grpcRoutes := range routes {
			for _, grpcRoute := range grpcRoutes {
				targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
						Name:  gwapiv1a2.ObjectName(grpcRoute.GetName()),
						Kind:  constantscommon.KindGRPCRoute,
						Group: constantscommon.K8sGroupNetworking,
					},
				})
			}
		}
		// Generate BackendTrafficPolicy for GRPC
		btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		backendTrafficPolicy := generateBackendTrafficPolicyForGRPC(btpName, targetRefs)
		if existing, ok := btpByName[btpName]; ok {
			mergeBTP(existing, backendTrafficPolicy)
		} else {
			btpByName[btpName] = backendTrafficPolicy.DeepCopy()
		}
	}

	kindType := constantscommon.KindHTTPRoute
	if apiResourceBundle.APKConf.Type == constants.API_TYPE_GRPC {
		kindType = constantscommon.KindGRPCRoute
	}

	// AI ratelimit
	endpoints := apiResourceBundle.APKConf.EndpointConfigurations.Production
	if environment == constants.SANDBOX_TYPE {
		endpoints = apiResourceBundle.APKConf.EndpointConfigurations.Sandbox
	}
	aiRatelimit := pickFirstAIRatelimit(endpoints)
	if aiRatelimit != nil {
		var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
		for _, httpRoutes := range routes {
			for _, httpRoute := range httpRoutes {
				targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
						Name:  gwapiv1a2.ObjectName(httpRoute.GetName()),
						Kind:  constantscommon.KindHTTPRoute,
						Group: constantscommon.K8sGroupNetworking,
					},
				})
			}
		}
		// Generate BackendTrafficPolicy for AI Ratelimit
		btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		backendTrafficPolicy := generateBackendTrafficPolicyForAIRatelimit(btpName, targetRefs, aiRatelimit)
		if existing, ok := btpByName[btpName]; ok {
			mergeBTP(existing, backendTrafficPolicy) // see helpers below
		} else {
			// store a deep-copy if your generator returns a pointer that might be reused
			btpByName[btpName] = backendTrafficPolicy.DeepCopy()
		}
	}

	// Ratelimit
	if apiResourceBundle.APKConf.RateLimit != nil {
		var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
		kindType := constantscommon.KindHTTPRoute
		if apiResourceBundle.APKConf.Type == constants.API_TYPE_GRPC {
			kindType = constantscommon.KindGRPCRoute
		}
		for _, httpRoutes := range routes {
			for _, httpRoute := range httpRoutes {
				targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
						Name:  gwapiv1a2.ObjectName(httpRoute.GetName()),
						Kind:  gwapiv1a2.Kind(kindType),
						Group: constantscommon.K8sGroupNetworking,
					},
				})
			}
		}
		// Generate BackendTrafficPolicy for Ratelimit
		btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		backendTrafficPolicy := generateBackendTrafficPolicyForRatelimit(btpName, targetRefs, apiResourceBundle.APKConf.RateLimit)
		if existing, ok := btpByName[btpName]; ok {
			mergeBTP(existing, backendTrafficPolicy)
		} else {
			btpByName[btpName] = backendTrafficPolicy.DeepCopy()
		}
	}

	// Operation level ratelimit
	ratelimitToCombinedResourceIndexMap := make(map[string]map[int][]int)
	for i, combinedResource := range apiResourceBundle.CombinedResources {
		ratelimit := combinedResource.APKOperations[0].RateLimit
		if ratelimit == nil {
			continue
		}
		if _, exists := ratelimitToCombinedResourceIndexMap[ratelimit.Unit]; !exists {
			ratelimitToCombinedResourceIndexMap[ratelimit.Unit] = make(map[int][]int)
		}
		ratelimitToCombinedResourceIndexMap[ratelimit.Unit][ratelimit.RequestsPerUnit] =
			append(ratelimitToCombinedResourceIndexMap[ratelimit.Unit][ratelimit.RequestsPerUnit], i)
	}
	counter := 0
	for unit, combinedResourceIndices := range ratelimitToCombinedResourceIndexMap {
		for requestsPerUnit, indices := range combinedResourceIndices {
			if len(indices) == 0 {
				continue
			}
			var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
			for _, index := range indices {
				for _, httpRoute := range routes[index] {
					targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
						LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
							Name:  gwapiv1a2.ObjectName(httpRoute.GetName()),
							Kind:  gwapiv1a2.Kind(kindType),
							Group: constantscommon.K8sGroupNetworking,
						},
					})
				}
			}
			counter++
			btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
				apiResourceBundle.Organization)
			btpName = fmt.Sprintf("%s-%d", btpName, counter)
			backendTrafficPolicy := generateBackendTrafficPolicyForRatelimit(btpName, targetRefs, &model.RateLimit{
				Unit:            unit,
				RequestsPerUnit: requestsPerUnit,
			})
			if existing, ok := btpByName[btpName]; ok {
				mergeBTP(existing, backendTrafficPolicy)
			} else {
				btpByName[btpName] = backendTrafficPolicy.DeepCopy()
			}
		}
	}

	// Append all BackendTrafficPolicy objects
	for _, b := range btpByName {
		objects = append(objects, b)
	}

	// CORS
	cors := apiResourceBundle.APKConf.CorsConfiguration
	if cors != nil && !cors.CorsConfigurationEnabled {
		cors = nil
	}

	// Operation level security
	for i, combinedResource := range apiResourceBundle.CombinedResources {
		isSecured, scopes := pickIsSecuredAndScopes(combinedResource.APKOperations)
		if !isSecured && cors == nil {
			continue
		}
		var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
		for _, httpRoute := range routes[i] {
			targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
				LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
					Name:  gwapiv1a2.ObjectName(httpRoute.GetName()),
					Kind:  gwapiv1a2.Kind(kindType),
					Group: constantscommon.K8sGroupNetworking,
				},
			})
		}
		spName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		spName = fmt.Sprintf("%s-%d", spName, i+1)
		sp := generateSecurityPolicy(spName, isSecured, scopes, targetRefs, cors, apiResourceBundle.APKConf.KeyManagers,
			apiResourceBundle.APKConf.Authentication, apiResourceBundle.Namespace)
		objects = append(objects, sp)
	}

	return objects, nil
}

func chunkOperations(ops []model.APKOperations, size int) [][]model.APKOperations {
	var chunks [][]model.APKOperations
	for size < len(ops) {
		ops, chunks = ops[size:], append(chunks, ops[0:size:size])
	}
	chunks = append(chunks, ops)
	return chunks
}

// GenerateGRPCRoutes generates GRPCRoute objects for the give APIResourceBundle.
func GenerateGRPCRoutes(bundle *dto.APIResourceBundle, withVersion bool, environment string,
	routePolicies map[string]*dpv2alpha1.RoutePolicy, routeMetadataList []*dpv2alpha1.RouteMetadata) (
	map[int][]gatewayv1.GRPCRoute, []client.Object) {
	objects := make([]client.Object, 0)
	routesMap := make(map[int][]gatewayv1.GRPCRoute)
	backendMap := make(map[string]map[string]*eg.Backend)
	crName := util.GenerateCRName(bundle.APKConf.Name, environment, bundle.APKConf.Version, bundle.Organization)
	// Generate the GRPCRoute objects
	for i, combined := range bundle.CombinedResources {
		batches := chunkOperations(combined.APKOperations, 16)
		for j, batch := range batches {
			parentName := config.GetConfig().ParentGatewayName
			parentNamespace := bundle.Namespace
			parentSectionName := config.GetConfig().ParentGatewaySectionName
			routeName := getHTTPRouteCRName(crName, i, j, withVersion)
			route := gatewayv1.GRPCRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: constantscommon.K8sGatewayAPIV1,
					Kind:       constantscommon.KindGRPCRoute,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: routeName,
					Annotations: map[string]string{
						constants.K8sHTTPRouteEnvTypeAnnotation: strings.ToUpper(environment),
					},
				},
				Spec: gatewayv1.GRPCRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						gatewayv1.Hostname(func() string {
							gatewayHostName := config.GetConfig().GatewayHostName
							if environment == constants.SANDBOX_TYPE {
								return fmt.Sprintf("%s.%s.%s", bundle.Organization, constants.SANDBOX_TYPE, gatewayHostName)
							}
							return fmt.Sprintf("%s.%s", bundle.Organization, gatewayHostName)
						}()),
					},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        gatewayv1.ObjectName(parentName),
								Group:       ptrTo(gatewayv1.Group(constantscommon.K8sGroupNetworking)),
								Kind:        ptrTo(gatewayv1.Kind(constantscommon.KindGateway)),
								Namespace:   ptrTo(gatewayv1.Namespace(parentNamespace)),
								SectionName: ptrTo(gatewayv1.SectionName(parentSectionName)),
							},
						},
					},
					Rules: []gatewayv1.GRPCRouteRule{},
				},
			}

			for _, op := range batch {
				if op.Verb == nil {
					getMethod := string(gatewayv1.HTTPMethodGet)
					op.Verb = &getMethod
				}
				if op.Target == nil {
					continue
				}

				apiBasePath := strings.TrimPrefix(bundle.APKConf.BasePath, "/")
				if withVersion {
					version := bundle.APKConf.Version
					apiBasePath = fmt.Sprintf("%s.%s",
						strings.TrimSuffix(apiBasePath, "/"),
						strings.TrimPrefix(version, "/"),
					)
				}
				path := fmt.Sprintf("%s.%s",
					strings.TrimSuffix(apiBasePath, "/"),
					strings.TrimPrefix(*op.Target, "/"),
				)
				method := safeGRPCMethod(op.Verb, &path)

				ecs := op.EndpointConfigurations.Production
				if environment == constants.SANDBOX_TYPE {
					ecs = op.EndpointConfigurations.Sandbox
				}
				if len(ecs) == 0 {
					continue
				}
				// Create backend reference
				grpcBackendRefs := createGRPCBackendRefs(ecs, backendMap, routeName)
				rule := gatewayv1.GRPCRouteRule{
					Matches: []gatewayv1.GRPCRouteMatch{
						{
							Method: method,
						},
					},
				}
				if len(grpcBackendRefs) > 0 {
					rule.BackendRefs = grpcBackendRefs
				} else {
					logger.Sugar().Warnf("No backend references found for operation %s in API %s", *op.Target, bundle.APKConf.Name)
				}
				route.Spec.Rules = append(route.Spec.Rules, rule)
			}
			routesMap[i] = append(routesMap[i], route)
			objects = append(objects, &route)
		}
	}
	for scheme, backendMapL := range backendMap {
		for _, backend := range backendMapL {
			objects = append(objects, backend)
			// Create BackendTLSPolicy if TLS is enabled
			if scheme == "https" {
				backendTLSPolicy := generateBackendTLSPolicyWithWellKnownCerts(backend.Name, backend.Spec.Endpoints[0].FQDN.Hostname,
					[]gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
								Name:  gwapiv1a2.ObjectName(backend.Name),
								Kind:  constantscommon.KindBackend,
								Group: constantscommon.EnvoyGateway,
							},
						},
					})
				objects = append(objects, backendTLSPolicy)
			}
		}
	}
	return routesMap, objects
}

// mergeBTP merges two BackendTrafficPolicy objects.
func mergeBTP(dst, src *eg.BackendTrafficPolicy) {
	// Labels/annotations (dst wins unless src adds new)
	if dst.Labels == nil {
		dst.Labels = map[string]string{}
	}
	for k, v := range src.Labels {
		dst.Labels[k] = v
	}
	if dst.Annotations == nil {
		dst.Annotations = map[string]string{}
	}
	for k, v := range src.Annotations {
		dst.Annotations[k] = v
	}

	// Merge target refs uniquely
	dst.Spec.TargetRefs = mergeUniqueTargetRefs(dst.Spec.TargetRefs, src.Spec.TargetRefs)

	// Merge RateLimits
	if src.Spec.RateLimit != nil {
		if dst.Spec.RateLimit == nil {
			dst.Spec.RateLimit = src.Spec.RateLimit.DeepCopy()
		} else {
			mergeRateLimit(dst.Spec.RateLimit, src.Spec.RateLimit)
		}
	}

	// Merge Cluster Settings
	if src.Spec.ClusterSettings != (eg.ClusterSettings{}) {
		if dst.Spec.ClusterSettings == (eg.ClusterSettings{}) {
			dst.Spec.ClusterSettings = *src.Spec.ClusterSettings.DeepCopy()
		} else {
			mergeClusterSettings(dst.Spec.ClusterSettings, src.Spec.ClusterSettings)
		}
	}

	// merge any other spec fields here
}

// mergeUniqueTargetRefs merges two slices of LocalPolicyTargetReferenceWithSectionName, ensuring that the result contains unique entries.
func mergeUniqueTargetRefs(a, b []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName) []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName {
	type key struct{ group, kind, ns, name string }
	seen := make(map[key]struct{}, len(a))
	out := make([]gwapiv1a2.LocalPolicyTargetReferenceWithSectionName, 0, len(a)+len(b))

	add := func(tr gwapiv1a2.LocalPolicyTargetReferenceWithSectionName) {
		k := key{group: string(tr.Group), kind: string(tr.Kind), name: string(tr.Name)}
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, tr)
		}
	}
	for _, tr := range a {
		add(tr)
	}
	for _, tr := range b {
		add(tr)
	}
	return out
}

// mergeRateLimit merges two RateLimitSpec objects.
func mergeRateLimit(dst, src *eg.RateLimitSpec) {
	if src.Type != "" {
		dst.Type = src.Type
	}
	if src.Global != nil {
		dst.Global = src.Global
	}
	if src.Local != nil {
		dst.Local = src.Local
	}
}

// mergeClusterSettings merges two ClusterSettings objects.
func mergeClusterSettings(dst, src eg.ClusterSettings) {
	if src.DNS != nil {
		dst.DNS = src.DNS.DeepCopy()
	}
	if src.LoadBalancer != nil {
		dst.LoadBalancer = src.LoadBalancer.DeepCopy()
	}
	if src.Timeout != nil {
		dst.Timeout = src.Timeout.DeepCopy()
	}
}

// GenerateHTTPRoutes generates HTTPRoute objects for the given APIResourceBundle.
func GenerateHTTPRoutes(bundle *dto.APIResourceBundle, withVersion bool, environment string, routePolicies map[string]*dpv2alpha1.RoutePolicy,
	routeMetadataList []*dpv2alpha1.RouteMetadata, kms []model.KeyManager) (map[int][]gatewayv1.HTTPRoute, []client.Object) {
	objects := make([]client.Object, 0)
	routesMap := make(map[int][]gatewayv1.HTTPRoute)
	backendMap := make(map[string]map[string]*eg.Backend)
	envoyExtensionPolicyMap := make(map[string]*eg.EnvoyExtensionPolicy)
	mapOfLuaSourceCodeConfigMap := make(map[string]*corev1.ConfigMap)
	crName := util.GenerateCRName(bundle.APKConf.Name, environment, bundle.APKConf.Version, bundle.Organization)
	for i, combined := range bundle.CombinedResources {
		batches := chunkOperations(combined.APKOperations, 16)

		for j, batch := range batches {
			parentName := config.GetConfig().ParentGatewayName
			parentNamespace := bundle.Namespace
			parentSectionName := config.GetConfig().ParentGatewaySectionName
			routeName := getHTTPRouteCRName(crName, i, j, withVersion)
			route := gatewayv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: constantscommon.K8sGatewayAPIV1,
					Kind:       constantscommon.KindHTTPRoute,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: routeName,
					Annotations: map[string]string{
						constants.K8sHTTPRouteEnvTypeAnnotation: strings.ToUpper(environment),
					},
				},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						gatewayv1.Hostname(func() string {
							gatewayHostName := config.GetConfig().GatewayHostName
							if environment == constants.SANDBOX_TYPE {
								return fmt.Sprintf("%s.%s.%s", bundle.Organization, constants.SANDBOX_TYPE, gatewayHostName)
							}
							return fmt.Sprintf("%s.%s", bundle.Organization, gatewayHostName)
						}()),
					},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        gatewayv1.ObjectName(parentName),
								Group:       ptrTo(gatewayv1.Group(constantscommon.K8sGroupNetworking)),
								Kind:        ptrTo(gatewayv1.Kind(constantscommon.KindGateway)),
								Namespace:   ptrTo(gatewayv1.Namespace(parentNamespace)),
								SectionName: ptrTo(gatewayv1.SectionName(parentSectionName)),
							},
						},
					},
					Rules: []gatewayv1.HTTPRouteRule{},
				},
			}
			interceptorPolicyList := make([]*model.APKOperationPolicy, 0)

			for _, op := range batch {
				backendBasePath, err := extractBackendBasePath(op, environment)
				if err != nil {
					logger.Sugar().Errorf("Error extracting backend base path for operation %s: %v", *op.Target, err)
					continue
				}
				backendHostname, err := extractBackendHostname(op, environment)
				if err != nil {
					logger.Sugar().Errorf("Error extracting backend hostname for operation %s: %v", *op.Target, err)
					continue
				}
				if op.Verb == nil {
					getMethod := string(gatewayv1.HTTPMethodGet)
					op.Verb = &getMethod
				}
				if op.Target == nil {
					continue
				}
				method := safeHTTPMethod(op.Verb)
				apiBasePath := bundle.APKConf.BasePath
				if withVersion {
					version := bundle.APKConf.Version
					apiBasePath = fmt.Sprintf("%s/%s",
						strings.TrimSuffix(apiBasePath, "/"),
						strings.TrimPrefix(version, "/"),
					)
				}

				path := fmt.Sprintf("%s/%s",
					strings.TrimSuffix(apiBasePath, "/"),
					strings.TrimSuffix(strings.TrimPrefix(*op.Target, "/"), "*"),
				)
				serviceContractPath := fmt.Sprintf("%s/%s",
					strings.TrimSuffix(backendBasePath, "/"),
					strings.TrimSuffix(strings.TrimPrefix(*op.Target, "/"), "*"),
				)

				rule := gatewayv1.HTTPRouteRule{}

				routePoliciesL := make(map[string]*dpv2alpha1.RoutePolicy)
				for key, policy := range routePolicies {
					routePoliciesL[key] = policy.DeepCopy()
				}

				ecs := op.EndpointConfigurations.Production
				if environment == constants.SANDBOX_TYPE {
					ecs = op.EndpointConfigurations.Sandbox
				}
				if len(ecs) == 0 {
					continue
				}
				ec := ecs[0] // Pick the first endpoint configuration for the operation
				if ec.EndpointSecurity != nil && *ec.EndpointSecurity.Enabled {
					endpointSecurityType, _ := securityType(ec.EndpointSecurity.SecurityType)
					switch endpointSecurityType {
					case securityTypeBasic:
						basicSecurity, ok := ec.EndpointSecurity.SecurityType.(*model.BasicEndpointSecurity)
						if !ok {
							if securityMap, mapOk := ec.EndpointSecurity.SecurityType.(map[string]interface{}); mapOk {
								basicSecurity = convertMapToBasicEndpointSecurity(securityMap)
							} else {
								logger.Sugar().Errorf("Failed to convert endpoint security to BasicEndpointSecurity for endpoint %v", ec.Endpoint)
								continue
							}
						}
						uniqueHash := generateBasicSecurityHash(basicSecurity)
						if routePoliciesL[uniqueHash] == nil {
							basicSecurityPolicy := createBasicEndpointSecurityMediationPolicy(uniqueHash, ec.EndpointSecurity, basicSecurity)
							routePoliciesL[uniqueHash] = basicSecurityPolicy
							objects = append(objects, basicSecurityPolicy)
						}
						delete(routePoliciesL, constants.BaseRoutePolicy)
					case securityTypeAPIKey:
						apiKeySecurity, ok := ec.EndpointSecurity.SecurityType.(*model.APIKeyEndpointSecurity)
						if !ok {
							if securityMap, mapOk := ec.EndpointSecurity.SecurityType.(map[string]interface{}); mapOk {
								apiKeySecurity = convertMapToAPIKeyEndpointSecurity(securityMap)
							} else {
								logger.Sugar().Errorf("Failed to convert endpoint security to APIKeyEndpointSecurity for endpoint %v", ec.Endpoint)
								continue
							}
						}
						uniqueHash := generateAPIKeySecurityHash(apiKeySecurity)
						if routePoliciesL[uniqueHash] == nil {
							apiKeyPolicy := createAPIKeyMediationPolicy(uniqueHash, ec.EndpointSecurity, apiKeySecurity)
							routePoliciesL[uniqueHash] = apiKeyPolicy
							objects = append(objects, apiKeyPolicy)
						}
						delete(routePoliciesL, constants.BaseRoutePolicy)
					}
				}

				var requestRedirectFilter *gatewayv1.HTTPRouteFilter

				// Add operation-level policy filters
				if op.OperationPolicies != nil {
					// Aggregate filters by type to avoid duplicates
					var requestHeaderModifier *gatewayv1.HTTPHeaderFilter
					var requestMirrorFilter *gatewayv1.HTTPRouteFilter
					var responseHeaderModifier *gatewayv1.HTTPHeaderFilter

					// Process request policies
					for _, requestPolicy := range op.OperationPolicies.Request {
						activePolicy := requestPolicy.GetActivePolicy()
						if activePolicy != nil {
							switch policy := activePolicy.(type) {
							case *model.HeaderModifierPolicy:
								if requestHeaderModifier == nil {
									requestHeaderModifier = &gatewayv1.HTTPHeaderFilter{}
								}
								switch policy.PolicyName {
								case model.PolicyNameAddHeader:
									requestHeaderModifier.Add = append(requestHeaderModifier.Add, gatewayv1.HTTPHeader{
										Name:  gatewayv1.HTTPHeaderName(policy.Parameters.HeaderName),
										Value: *policy.Parameters.HeaderValue,
									})
								case model.PolicyNameSetHeader:
									requestHeaderModifier.Set = append(requestHeaderModifier.Set, gatewayv1.HTTPHeader{
										Name:  gatewayv1.HTTPHeaderName(policy.Parameters.HeaderName),
										Value: *policy.Parameters.HeaderValue,
									})
								case model.PolicyNameRemoveHeader:
									requestHeaderModifier.Remove = append(requestHeaderModifier.Remove, policy.Parameters.HeaderName)
								}
							case *model.LuaInterceptorPolicy, *model.WASMInterceptorPolicy:
								interceptorPolicyList = append(interceptorPolicyList, &policy)
							case *model.BackendJWTPolicy:
								backendJWTPolicy := createBackendJWTMediationPolicy(policy, kms)
								routePoliciesL[backendJWTPolicy.Name] = backendJWTPolicy
								objects = append(objects, backendJWTPolicy)
							case *model.RequestMirrorPolicy:
								mirrorEndpoints := []model.EndpointConfiguration{
									{
										Endpoint: policy.Parameters.URLs[0],
									},
								}
								requestMirrorBackendRefs := createBackendRefs(mirrorEndpoints, backendMap, routeName)
								requestMirrorBackendRef := requestMirrorBackendRefs[0].BackendObjectReference
								requestMirrorFilter = &gatewayv1.HTTPRouteFilter{
									Type: gatewayv1.HTTPRouteFilterRequestMirror,
									RequestMirror: &gatewayv1.HTTPRequestMirrorFilter{
										BackendRef: requestMirrorBackendRef,
										Percent:    ptrTo(int32(100)),
									},
								}
							case *model.RequestRedirectPolicy:
								scheme, host, port, err := extractSchemeHostPort(policy.Parameters.URL)
								endpointPath, _ := extractPathFromEndpoint(policy.Parameters.URL)
								redirectPath := fmt.Sprintf("%s/%s",
									strings.TrimSuffix(endpointPath, "/"),
									strings.TrimSuffix(strings.TrimPrefix(*op.Target, "/"), "*"),
								)
								isRegexPath, _, substitution := GenerateRegexPath(path, endpointPath, apiBasePath)
								if isRegexPath {
									redirectPath = substitution
								}
								if err != nil {
									logger.Sugar().Errorf("Error extracting scheme, host, and port from URL %s: %v", policy.Parameters.URL, err)
									continue
								}
								requestRedirectFilter = &gatewayv1.HTTPRouteFilter{
									Type: gatewayv1.HTTPRouteFilterRequestRedirect,
									RequestRedirect: &gatewayv1.HTTPRequestRedirectFilter{
										Scheme:   ptrTo(scheme),
										Hostname: ptrTo(gatewayv1.PreciseHostname(host)),
										Path: &gatewayv1.HTTPPathModifier{
											Type:            gatewayv1.FullPathHTTPPathModifier,
											ReplaceFullPath: &redirectPath,
										},
										Port:       ptrTo(gatewayv1.PortNumber(port)),
										StatusCode: policy.Parameters.StatusCode,
									},
								}
							case *model.ModelBasedRoundRobinPolicy:
								modelBasedRoundRobinPolicy, err := generateModelBasedRoundRobinPolicy(policy, environment)
								if err != nil {
									logger.Sugar().Errorf("Error generating ModelBasedRoundRobinPolicy: %v", err)
									continue
								}
								routePoliciesL[modelBasedRoundRobinPolicy.Name] = modelBasedRoundRobinPolicy
								objects = append(objects, modelBasedRoundRobinPolicy)
							case *model.CommonPolicy:
								aiGuardrailPolicy := generateAIGuardrailPolicy(policy, constantscommon.REQUEST_FLOW)
								routePoliciesL[aiGuardrailPolicy.Name] = aiGuardrailPolicy
								objects = append(objects, aiGuardrailPolicy)
							}
						}
					}

					// Process response policies
					for _, responsePolicy := range op.OperationPolicies.Response {
						activePolicy := responsePolicy.GetActivePolicy()
						if activePolicy != nil {
							switch policy := activePolicy.(type) {
							case *model.HeaderModifierPolicy:
								if responseHeaderModifier == nil {
									responseHeaderModifier = &gatewayv1.HTTPHeaderFilter{}
								}
								switch policy.PolicyName {
								case model.PolicyNameAddHeader:
									responseHeaderModifier.Add = append(responseHeaderModifier.Add, gatewayv1.HTTPHeader{
										Name:  gatewayv1.HTTPHeaderName(policy.Parameters.HeaderName),
										Value: *policy.Parameters.HeaderValue,
									})
								case model.PolicyNameSetHeader:
									responseHeaderModifier.Set = append(responseHeaderModifier.Set, gatewayv1.HTTPHeader{
										Name:  gatewayv1.HTTPHeaderName(policy.Parameters.HeaderName),
										Value: *policy.Parameters.HeaderValue,
									})
								case model.PolicyNameRemoveHeader:
									responseHeaderModifier.Remove = append(responseHeaderModifier.Remove, policy.Parameters.HeaderName)
								}
							case *model.LuaInterceptorPolicy, *model.WASMInterceptorPolicy:
								interceptorPolicyList = append(interceptorPolicyList, &policy)
							case *model.CommonPolicy:
								aiGuardrailPolicy := generateAIGuardrailPolicy(policy, constantscommon.RESPONSE_FLOW)
								routePoliciesL[aiGuardrailPolicy.Name] = aiGuardrailPolicy
								objects = append(objects, aiGuardrailPolicy)
							}
						}
					}

					// Add aggregated filters to the rule
					if requestHeaderModifier != nil {
						rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
							Type:                  gatewayv1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: requestHeaderModifier,
						})
					}
					if requestMirrorFilter != nil {
						rule.Filters = append(rule.Filters, *requestMirrorFilter)
					}
					if requestRedirectFilter != nil {
						rule.Filters = append(rule.Filters, *requestRedirectFilter)
					}
					if responseHeaderModifier != nil {
						rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
							Type:                   gatewayv1.HTTPRouteFilterResponseHeaderModifier,
							ResponseHeaderModifier: responseHeaderModifier,
						})
					}
				}

				for _, policy := range routePoliciesL {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constantscommon.WSO2KubernetesGateway,
							Kind:  constantscommon.KindRoutePolicy,
							Name:  gatewayv1.ObjectName(policy.Name),
						},
					})
				}
				for _, metadata := range routeMetadataList {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constantscommon.WSO2KubernetesGateway,
							Kind:  constantscommon.KindRouteMetadata,
							Name:  gatewayv1.ObjectName(metadata.Name),
						},
					})
				}

				isRegexPath, pattern, substitution := GenerateRegexPath(path, backendBasePath, apiBasePath)
				hrfName := ""
				pathMatchType := gatewayv1.PathMatchExact
				if strings.HasSuffix(*op.Target, "*") {
					pathMatchType = gatewayv1.PathMatchPathPrefix
				}

				if isRegexPath {
					pathMatchType = gatewayv1.PathMatchRegularExpression
					path = pattern
				}
				rule.Matches = []gatewayv1.HTTPRouteMatch{
					{
						Path: &gatewayv1.HTTPPathMatch{
							Type:  ptrTo(pathMatchType),
							Value: ptrTo(path),
						},
						Method: method,
					},
				}
				if bundle.APKConf.CorsConfiguration != nil && bundle.APKConf.CorsConfiguration.CorsConfigurationEnabled {
					rule.Matches = append(rule.Matches,
						gatewayv1.HTTPRouteMatch{
							Path: &gatewayv1.HTTPPathMatch{
								Type:  ptrTo(pathMatchType),
								Value: ptrTo(path),
							},
							Method: ptrTo(gatewayv1.HTTPMethodOptions),
						},
					)
				}

				if requestRedirectFilter == nil {
					// Create backend reference
					httpBackendRefs := createBackendRefs(ecs, backendMap, routeName)
					if len(httpBackendRefs) > 0 {
						rule.BackendRefs = httpBackendRefs
					} else {
						logger.Sugar().Warnf("No backend references found for operation %s in API %s", *op.Target, bundle.APKConf.Name)
					}

					if isRegexPath {
						pathMatchType = gatewayv1.PathMatchRegularExpression
						path = pattern
						// Create HTTPRouteFilter
						sum := sha256.Sum256([]byte(fmt.Sprintf("%s-%s", path, string(*method))))
						pathIdentifier := fmt.Sprintf("%x", sum[:8])
						hrfName = fmt.Sprintf("%s-%s", routeName, pathIdentifier)
						hrf := eg.HTTPRouteFilter{
							TypeMeta: metav1.TypeMeta{
								Kind:       constantscommon.KindHTTPRouteFilter,
								APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name: hrfName,
							},
							Spec: eg.HTTPRouteFilterSpec{
								URLRewrite: &eg.HTTPURLRewriteFilter{
									Path: &eg.HTTPPathModifier{
										Type: eg.RegexHTTPPathModifier,
										ReplaceRegexMatch: &eg.ReplaceRegexMatch{
											Pattern:      pattern,
											Substitution: substitution,
										},
									},
								},
							},
						}
						objects = append(objects, &hrf)
					}
				}

				if hrfName != "" {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constantscommon.EnvoyGateway,
							Kind:  constantscommon.KindHTTPRouteFilter,
							Name:  gatewayv1.ObjectName(hrfName),
						},
					})
				} else if requestRedirectFilter == nil {
					urlRewrite := &gatewayv1.HTTPURLRewriteFilter{
						Path: &gatewayv1.HTTPPathModifier{
							Type: gatewayv1.FullPathHTTPPathModifier,
						},
					}
					if strings.HasSuffix(*op.Target, "*") {
						urlRewrite.Path.ReplacePrefixMatch = &serviceContractPath
						urlRewrite.Path.Type = gatewayv1.PrefixMatchHTTPPathModifier
					} else {
						urlRewrite.Path.ReplaceFullPath = &serviceContractPath
					}
					if backendHostname != "" {
						urlRewrite.Hostname = ptrTo(gatewayv1.PreciseHostname(backendHostname))
					}
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type:       gatewayv1.HTTPRouteFilterURLRewrite,
						URLRewrite: urlRewrite,
					})
				}
				route.Spec.Rules = append(route.Spec.Rules, rule)
			}

			for _, interceptorPolicy := range interceptorPolicyList {
				envoyExtensionPolicy, err := generateEnvoyExtensionPolicy(interceptorPolicyList, routeName)
				if err != nil {
					logger.Sugar().Errorf("Error generating Lua EnvoyExtensionPolicy: %v", err)
					continue
				}
				envoyExtensionPolicyMap[envoyExtensionPolicy.Name] = envoyExtensionPolicy
				if luaPolicy, ok := (*interceptorPolicy).(*model.LuaInterceptorPolicy); ok {
					if luaPolicy.Parameters.MountInConfigMap != nil && *luaPolicy.Parameters.MountInConfigMap {
						luaSourceCodeConfigMap := createConfigMapForLuaSourceCode(luaPolicy.Parameters)
						mapOfLuaSourceCodeConfigMap[luaSourceCodeConfigMap.Name] = luaSourceCodeConfigMap
					}
				}
			}

			routesMap[i] = append(routesMap[i], route)
			objects = append(objects, &route)
		}
	}
	for _, envoyExtensionPolicy := range envoyExtensionPolicyMap {
		objects = append(objects, envoyExtensionPolicy)
	}
	for _, luaSourceCodeConfigMap := range mapOfLuaSourceCodeConfigMap {
		objects = append(objects, luaSourceCodeConfigMap)
	}
	for scheme, backendMapL := range backendMap {
		for _, backend := range backendMapL {
			objects = append(objects, backend)
			// Create BackendTLSPolicy if TLS is enabled
			if scheme == "https" {
				backendTLSPolicy := generateBackendTLSPolicyWithWellKnownCerts(backend.Name, backend.Spec.Endpoints[0].FQDN.Hostname,
					[]gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
								Name:  gwapiv1a2.ObjectName(backend.Name),
								Kind:  constantscommon.KindBackend,
								Group: constantscommon.EnvoyGateway,
							},
						},
					})
				objects = append(objects, backendTLSPolicy)
			}
		}
	}
	return routesMap, objects
}

// generateEnvoyExtensionPolicy generates an EnvoyExtensionPolicy for a given LuaInterceptorPolicy or WASMInterceptorPolicy
func generateEnvoyExtensionPolicy(interceptorPolicyList []*model.APKOperationPolicy, routeName string) (*eg.EnvoyExtensionPolicy, error) {
	luaFilterList := make([]eg.Lua, 0)
	wasmFilterList := make([]eg.Wasm, 0)
	filterNameList := make([]string, 0)
	for _, interceptorPolicy := range interceptorPolicyList {
		if luaPolicy, ok := (*interceptorPolicy).(*model.LuaInterceptorPolicy); ok {
			if !slices.Contains(filterNameList, luaPolicy.Parameters.Name) {
				luaFilter, err := createLuaFilter(luaPolicy.Parameters)
				if err != nil {
					return nil, err
				}
				luaFilterList = append(luaFilterList, *luaFilter)
				filterNameList = append(filterNameList, luaPolicy.Parameters.Name)
			}
		} else if wasmPolicy, ok := (*interceptorPolicy).(*model.WASMInterceptorPolicy); ok {
			if !slices.Contains(filterNameList, wasmPolicy.Parameters.Name) {
				wasmFilter, err := createWASMFilter(wasmPolicy.Parameters, routeName)
				if err != nil {
					return nil, err
				}
				wasmFilterList = append(wasmFilterList, *wasmFilter)
				filterNameList = append(filterNameList, wasmPolicy.Parameters.Name)
			}
		}
	}
	filterNameList = append(filterNameList, routeName)
	policyName := util.SanitizeOrHashName(strings.Join(filterNameList, "-"))
	luaEnvoyExtensionPolicy := &eg.EnvoyExtensionPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindEnvoyExtensionPolicy,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: policyName,
		},
		Spec: eg.EnvoyExtensionPolicySpec{
			Lua:  luaFilterList,
			Wasm: wasmFilterList,
			PolicyTargetReferences: eg.PolicyTargetReferences{
				TargetRefs: []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
					{
						LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
							Name:  gwapiv1a2.ObjectName(routeName),
							Kind:  constantscommon.KindHTTPRoute,
							Group: constantscommon.K8sGroupNetworking,
						},
					},
				},
			},
		},
	}
	return luaEnvoyExtensionPolicy, nil
}

// createLuaFilter creates a Lua filter configuration based on the parameters
func createLuaFilter(parameters *model.LuaInterceptorPolicyParameters) (*eg.Lua, error) {
	luaValueType, err := getLuaValueType(parameters)
	if err != nil {
		return nil, err
	}
	if luaValueType == eg.LuaValueTypeInline {
		return &eg.Lua{
			Type:   luaValueType,
			Inline: parameters.SourceCode,
		}, nil
	}

	return &eg.Lua{
		Type: luaValueType,
		ValueRef: &gatewayv1.LocalObjectReference{
			Kind: constantscommon.KindConfigMap,
			Name: gatewayv1.ObjectName(*parameters.SourceCodeRef),
		},
	}, nil
}

// getLuaValueType determines the Lua value type based on the provided parameters
func getLuaValueType(parameters *model.LuaInterceptorPolicyParameters) (eg.LuaValueType, error) {
	if parameters.MountInConfigMap != nil && *parameters.MountInConfigMap {
		if parameters.SourceCodeRef == nil || *parameters.SourceCodeRef == "" ||
			parameters.SourceCode == nil || *parameters.SourceCode == "" {
			return "", fmt.Errorf("both SourceCodeRef and SourceCode must be set when MountInConfigMap is true")
		}
		return eg.LuaValueTypeValueRef, nil
	}
	if parameters.SourceCodeRef != nil && *parameters.SourceCodeRef != "" {
		return eg.LuaValueTypeValueRef, nil
	}
	if parameters.SourceCode != nil && *parameters.SourceCode != "" {
		return eg.LuaValueTypeInline, nil
	}
	return "", fmt.Errorf("either SourceCode or SourceCodeRef must be set")
}

// createWASMFilter creates a WASM filter configuration based on the parameters
func createWASMFilter(parameters *model.WASMInterceptorPolicyParameters, routeName string) (*eg.Wasm, error) {
	// make the name unique by appending a hash of the route name
	routeNameBytes := []byte(routeName)
	hash := sha256.Sum256(routeNameBytes)
	wasmFilterName := fmt.Sprintf("%s-%x", parameters.Name, hash[:16])

	wasmFilter := &eg.Wasm{
		Name:     &wasmFilterName,
		RootID:   &parameters.RootID,
		FailOpen: parameters.FailOpen,
		Env: &eg.WasmEnv{
			HostKeys: parameters.HostKeys,
		},
	}
	if parameters.Config != nil && *parameters.Config != "" {
		configJSON := &apiextensionsv1.JSON{
			Raw: []byte(*parameters.Config),
		}
		wasmFilter.Config = configJSON
	}
	var imagePullPolicy eg.ImagePullPolicy
	if parameters.ImagePullPolicy != nil && *parameters.ImagePullPolicy != "" {
		imagePullPolicy = eg.ImagePullPolicy(*parameters.ImagePullPolicy)
	}
	wasmSourceCodeType, err := getWASMSourceCodeType(parameters)
	if err != nil {
		return nil, err
	}
	if wasmSourceCodeType == eg.HTTPWasmCodeSourceType {
		wasmFilter.Code = eg.WasmCodeSource{
			Type: wasmSourceCodeType,
			HTTP: &eg.HTTPWasmCodeSource{
				URL: *parameters.URL,
			},
			PullPolicy: &imagePullPolicy,
		}
	} else if wasmSourceCodeType == eg.ImageWasmCodeSourceType {
		wasmFilter.Code = eg.WasmCodeSource{
			Type: wasmSourceCodeType,
			Image: &eg.ImageWasmCodeSource{
				URL: *parameters.Image,
			},
			PullPolicy: &imagePullPolicy,
		}
	}
	return wasmFilter, nil
}

// getLuaValueType determines the Lua value type based on the provided parameters
func getWASMSourceCodeType(parameters *model.WASMInterceptorPolicyParameters) (eg.WasmCodeSourceType, error) {
	if parameters.URL != nil && *parameters.URL != "" {
		return eg.HTTPWasmCodeSourceType, nil
	}
	if parameters.Image != nil && *parameters.Image != "" {
		return eg.ImageWasmCodeSourceType, nil
	}
	return "", fmt.Errorf("either URL or Image must be set")
}

// createConfigMapForLuaSourceCode creates a ConfigMap to hold the Lua source code for a given LuaInterceptorPolicy
func createConfigMapForLuaSourceCode(parameters *model.LuaInterceptorPolicyParameters) *corev1.ConfigMap {
	cmName := util.SanitizeOrHashName(*parameters.SourceCodeRef)
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       constantscommon.KindConfigMap,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: cmName,
		},
		Data: map[string]string{
			"lua": *parameters.SourceCode,
		},
	}
}

func getHTTPRouteCRName(crName string, i int, j int, includesVersion bool) string {
	var withVersion string
	if includesVersion {
		withVersion = "1"
	} else {
		withVersion = "0"
	}
	return fmt.Sprintf("%s-%d-%d-%s", crName, i+1, j+1, withVersion)
}

func ptrTo[T any](v T) *T {
	return &v
}

// GenerateRegexPath analyzes the URL path and constructs a regex pattern and substitution if it contains path variables.
func GenerateRegexPath(path, basePath, apiBasePath string) (isRegexPath bool, pattern string, substitution string) {
	// Regular expression to find {var} patterns
	re := regexp.MustCompile(`\{[^{}]+\}`)

	// Check if path contains any {var}
	matches := re.FindAllStringIndex(path, -1)
	if len(matches) == 0 {
		substitution = strings.Replace(basePath+path, apiBasePath, "", 1)
		substitution = strings.ReplaceAll(substitution, "//", "/")
		return false, path, substitution
	}
	// If we have matches, we will treat this as a regex path
	isRegexPath = true

	// Replace {var} with (.*) in pattern
	pattern = re.ReplaceAllString(path, `([^/]+)`)

	// Build substitution string
	substitutionBuilder := strings.Builder{}
	substitutionBuilder.WriteString(basePath)

	start := 0
	matchIndex := 1
	for _, match := range matches {
		substitutionBuilder.WriteString(path[start:match[0]])
		substitutionBuilder.WriteString(fmt.Sprintf(`\%d`, matchIndex))
		matchIndex++
		start = match[1]
	}
	substitutionBuilder.WriteString(path[start:])

	substitution = substitutionBuilder.String()
	substitution = strings.Replace(substitution, apiBasePath, "", 1)
	substitution = strings.ReplaceAll(substitution, "//", "/")
	pattern = strings.ReplaceAll(pattern, "//", "/")
	return
}

func pickFirstAIRatelimit(ecs []model.EndpointConfiguration) *model.AIRatelimit {
	for _, ec := range ecs {
		if ec.AIRatelimit != nil {
			return ec.AIRatelimit
		}
	}
	return nil
}

// Generate BackendTrafficPolicy for AI Ratelimit
func generateBackendTrafficPolicyForAIRatelimit(name string, targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName,
	rlConf *model.AIRatelimit) *eg.BackendTrafficPolicy {
	return &eg.BackendTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindBackendTrafficPolicy,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: eg.BackendTrafficPolicySpec{
			PolicyTargetReferences: eg.PolicyTargetReferences{
				TargetRefs: targetRefs,
			},
			RateLimit: &eg.RateLimitSpec{
				Type: eg.GlobalRateLimitType,
				Global: &eg.GlobalRateLimit{
					Rules: generateAIRatelimitRules(rlConf),
				},
			},
		},
	}
}

// Generate BackendTrafficPolicy for Ratelimit
func generateBackendTrafficPolicyForRatelimit(name string, targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName,
	rlConf *model.RateLimit) *eg.BackendTrafficPolicy {
	return &eg.BackendTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindBackendTrafficPolicy,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: eg.BackendTrafficPolicySpec{
			PolicyTargetReferences: eg.PolicyTargetReferences{
				TargetRefs: targetRefs,
			},
			MergeType: ptrTo(eg.StrategicMerge),
			RateLimit: &eg.RateLimitSpec{
				Type: eg.GlobalRateLimitType,

				Global: &eg.GlobalRateLimit{
					Rules: generateRatelimitRules(rlConf),
				},
			},
		},
	}
}

// Generate BackendTrafficPolicy for GRPC
func generateBackendTrafficPolicyForGRPC(name string, targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName) *eg.BackendTrafficPolicy {
	return &eg.BackendTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindBackendTrafficPolicy,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: eg.BackendTrafficPolicySpec{
			PolicyTargetReferences: eg.PolicyTargetReferences{
				TargetRefs: targetRefs,
			},
			MergeType: ptrTo(eg.StrategicMerge),
			ClusterSettings: eg.ClusterSettings{
				Timeout: &eg.Timeout{
					HTTP: &eg.HTTPTimeout{
						RequestTimeout: ptrTo(gatewayv1.Duration("0s")),
					},
				},
			},
		},
	}
}

// generateRatelimitRules generates the rate limit rules based on the RatelimitConfiguration.
func generateRatelimitRules(rlConf *model.RateLimit) []eg.RateLimitRule {
	var ratelimitRules []eg.RateLimitRule
	ratelimitRules = append(ratelimitRules, eg.RateLimitRule{
		Limit: eg.RateLimitValue{
			Unit:     eg.RateLimitUnit(rlConf.Unit),
			Requests: uint(rlConf.RequestsPerUnit),
		},
	})
	return ratelimitRules
}

// generateRatelimitRules generates the rate limit rules based on the RatelimitConfiguration.
func generateAIRatelimitRules(rlConf *model.AIRatelimit) []eg.RateLimitRule {
	var ratelimitRules []eg.RateLimitRule

	ratelimitRules = append(ratelimitRules, eg.RateLimitRule{
		Limit: eg.RateLimitValue{
			Unit:     eg.RateLimitUnit(rlConf.Token.Unit),
			Requests: uint(rlConf.Token.PromptLimit),
		},
		// Shared: func(b bool) *bool { return &b }(true),
		Cost: &eg.RateLimitCost{
			Request: &eg.RateLimitCostSpecifier{
				From:   eg.RateLimitCostFromNumber,
				Number: func(v uint64) *uint64 { return &v }(0),
			},
			Response: &eg.RateLimitCostSpecifier{
				From: eg.RateLimitCostFromMetadata,
				Metadata: &eg.RateLimitCostMetadata{
					Namespace: constantscommon.MetadataNamespace,
					Key:       constantscommon.PromptTokenCountIDMetadataKey,
				},
			},
		},
	})

	ratelimitRules = append(ratelimitRules, eg.RateLimitRule{
		Limit: eg.RateLimitValue{
			Unit:     eg.RateLimitUnit(rlConf.Token.Unit),
			Requests: uint(rlConf.Token.CompletionLimit),
		},
		// Shared: func(b bool) *bool { return &b }(true),
		Cost: &eg.RateLimitCost{
			Request: &eg.RateLimitCostSpecifier{
				From:   eg.RateLimitCostFromNumber,
				Number: func(v uint64) *uint64 { return &v }(0),
			},
			Response: &eg.RateLimitCostSpecifier{
				From: eg.RateLimitCostFromMetadata,
				Metadata: &eg.RateLimitCostMetadata{
					Namespace: constantscommon.MetadataNamespace,
					Key:       constantscommon.CompletionTokenCountIDMetadataKey,
				},
			},
		},
	})

	ratelimitRules = append(ratelimitRules, eg.RateLimitRule{
		Limit: eg.RateLimitValue{
			Unit:     eg.RateLimitUnit(rlConf.Token.Unit),
			Requests: uint(rlConf.Token.TotalLimit),
		},
		// Shared: func(b bool) *bool { return &b }(true),
		Cost: &eg.RateLimitCost{
			Request: &eg.RateLimitCostSpecifier{
				From:   eg.RateLimitCostFromNumber,
				Number: func(v uint64) *uint64 { return &v }(0),
			},
			Response: &eg.RateLimitCostSpecifier{
				From: eg.RateLimitCostFromMetadata,
				Metadata: &eg.RateLimitCostMetadata{
					Namespace: constantscommon.MetadataNamespace,
					Key:       constantscommon.TotalTokenCountIDMetadataKey,
				},
			},
		},
	})

	ratelimitRules = append(ratelimitRules, eg.RateLimitRule{
		Limit: eg.RateLimitValue{
			Unit:     eg.RateLimitUnit(rlConf.Request.Unit),
			Requests: uint(rlConf.Request.RequestLimit),
		},
	})

	return ratelimitRules
}

func pickIsSecuredAndScopes(apkOperations []model.APKOperations) (bool, []string) {
	isSecured := false
	scopes := make([]string, 0)
	if len(apkOperations) > 0 && apkOperations[0].Secured != nil {
		isSecured = *apkOperations[0].Secured
		scopes = apkOperations[0].Scopes
	}
	return isSecured, scopes
}

// generateSecurityPolicy generates a SecurityPolicy object based on the provided parameters.
func generateSecurityPolicy(name string, isSecured bool, scopes []string, targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName,
	cors *model.CORSConfiguration, kms []model.KeyManager, auths []model.AuthenticationRequest, namespace string) *eg.SecurityPolicy {
	sp := &eg.SecurityPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindSecurityPolicy,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: eg.SecurityPolicySpec{
			PolicyTargetReferences: eg.PolicyTargetReferences{
				TargetRefs: targetRefs,
			},
		},
	}
	if cors != nil {
		corsOrigins := ConvertStringsToOrigins(cors.AccessControlAllowOrigins)
		if len(corsOrigins) == 0 {
			logger.Sugar().Warnf("No valid CORS origins found in configuration for API %s", name)
		}
		sp.Spec.CORS = &eg.CORS{
			AllowOrigins:     corsOrigins,
			AllowMethods:     cors.AccessControlAllowMethods,
			AllowHeaders:     cors.AccessControlAllowHeaders,
			ExposeHeaders:    cors.AccessControlExposeHeaders,
			MaxAge:           convertIntPtrToDurationSecondsForMaxAge(cors.AccessControlAllowMaxAge, 3600),
			AllowCredentials: cors.AccessControlAllowCredentials,
		}
	}
	headerNames := extractHeaderNames(auths)
	jwtHeaderExtractors := make([]eg.JWTHeaderExtractor, 0, len(headerNames))
	for _, h := range headerNames {
		jwtHeaderExtractors = append(jwtHeaderExtractors, eg.JWTHeaderExtractor{
			Name:        h,
			ValuePrefix: ptrTo("Bearer "),
		})
	}
	isAuthEnabled := len(auths) == 0
	for _, auth := range auths {
		authType := auth.GetAuthType()
		switch a := authType.(type) {
		case *model.OAuth2Authentication:
			if a.Enabled {
				isAuthEnabled = true
				break
			}
		case *model.JWTAuthentication:
			if a.Enabled {
				isAuthEnabled = true
				break
			}
		case *model.APIKeyAuthentication:
			if a.Enabled {
				isAuthEnabled = true
				break
			}
		case *model.MTLSAuthentication:
			if a.Enabled {
				isAuthEnabled = true
				break
			}
		}
	}

	if isSecured && isAuthEnabled {
		if len(kms) == 0 {
			defaultIDPKM := generateDefaultIDPKeyManager(namespace)
			kms = append(kms, defaultIDPKM)
		}
		for _, km := range kms {
			provider := eg.JWTProvider{
				Name: km.Name,
				RemoteJWKS: &eg.RemoteJWKS{
					URI: km.JWKSEndpoint,
				},
				ExtractFrom: &eg.JWTExtractor{
					Headers: jwtHeaderExtractors,
				},
				Issuer: km.Issuer,
			}
			if km.K8sBackend != nil && km.K8sBackend.Name != nil && km.K8sBackend.Port != nil {
				provider.RemoteJWKS.BackendRefs = []eg.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Group:     ptrTo(gatewayv1.Group(constantscommon.EnvoyGateway)),
							Kind:      ptrTo(gatewayv1.Kind(constantscommon.KindBackend)),
							Name:      gatewayv1.ObjectName(*km.K8sBackend.Name),
							Namespace: ptrTo(gatewayv1.Namespace(*km.K8sBackend.Namespace)),
							Port:      ptrTo(gatewayv1.PortNumber(*km.K8sBackend.Port)),
						},
					},
				}
			}

			if sp.Spec.JWT == nil {
				sp.Spec.JWT = &eg.JWT{
					Providers: []eg.JWTProvider{provider},
				}
			} else {
				sp.Spec.JWT.Providers = append(sp.Spec.JWT.Providers, provider)
			}
		}
		if len(scopes) > 0 {
			rules := make([]eg.AuthorizationRule, 0, len(kms))
			jwtScopes := ConvertStringsToJWTScope(scopes)
			for _, km := range kms {
				rules = append(rules, eg.AuthorizationRule{
					Name:   ptrTo(km.Name),
					Action: eg.AuthorizationActionAllow,
					Principal: eg.Principal{
						JWT: &eg.JWTPrincipal{
							Provider: km.Name,
							Scopes:   jwtScopes,
						},
					},
				})
			}
			sp.Spec.Authorization = &eg.Authorization{
				Rules: rules,
			}
		}
	}
	return sp
}

// generateDefaultIDPKeyManager generates a default KeyManager configuration for the inbuilt IDP.
func generateDefaultIDPKeyManager(namespace string) model.KeyManager {
	k8sRelease := config.GetConfig().K8sReleaseName
	k8sResourcePrefix := fmt.Sprintf("%s-%s", k8sRelease, config.GetConfig().K8sResourcePrefix)
	jwtProviderName := fmt.Sprintf("%s-idp-jwt-issuer", k8sResourcePrefix)
	jwksURI := fmt.Sprintf("https://%s-idp-ds-service.%s.svc:%s/oauth2/jwks", k8sResourcePrefix, namespace,
		config.GetConfig().ConfigDSServerPort)
	defaultDSBackendName := fmt.Sprintf("%s-oauth-ds-backend", k8sResourcePrefix)

	km := model.KeyManager{
		Name:         jwtProviderName,
		JWKSEndpoint: jwksURI,
		K8sBackend: &model.K8sBackend{
			Name: ptrTo(defaultDSBackendName),
			Port: ptrTo(func() int {
				port, err := strconv.Atoi(config.GetConfig().ConfigDSServerPort)
				if err != nil {
					return 9443
				}
				return port
			}()),
			Namespace: ptrTo(namespace),
		},
	}
	return km
}

var originPattern = regexp.MustCompile(`^(\*|https?:\/\/(\*|(\*\.)?(([\w-]+\.?)+)?[\w-]+)(:\d{1,5})?)$`)

// ConvertStringsToOrigins validates and converts a list of strings into []Origin
func ConvertStringsToOrigins(inputs []string) []eg.Origin {
	var origins []eg.Origin
	for _, s := range inputs {
		if !originPattern.MatchString(s) {
			continue
		}
		origins = append(origins, eg.Origin(s))
	}
	return origins
}

// ConvertIntPtrToDuration converts a pointer to int into a Duration string.
// If ptr is nil, it uses the default value. Both are suffixed with "s" by default
// unless you want another unit.
func convertIntPtrToDurationSecondsForMaxAge(ptr *int, defaultValue int) *gatewayv1.Duration {
	// Decide which value to use
	val := defaultValue
	if ptr != nil {
		val = *ptr
	}
	// Example: treat int as seconds and append "s"
	durStr := fmt.Sprintf("%ds", val)
	dur := gatewayv1.Duration(durStr)
	return &dur
}

func extractHeaderNames(auths []model.AuthenticationRequest) []string {
	headerNames := make([]string, 0)
	for _, auth := range auths {
		authType := auth.GetAuthType()
		switch v := authType.(type) {
		case *model.OAuth2Authentication:
			headerNames = append(headerNames, v.HeaderName)
		case *model.JWTAuthentication:
			headerNames = append(headerNames, v.HeaderName)
		case *model.APIKeyAuthentication:
			if v.HeaderEnable {
				headerNames = append(headerNames, v.HeaderName)
			}
		}
	}
	return headerNames
}

// ConvertStringsToJWTScope converts a []string into a []v1alpha1.JWTScope.
func ConvertStringsToJWTScope(scopes []string) []eg.JWTScope {
	result := make([]eg.JWTScope, len(scopes))
	for i, s := range scopes {
		result[i] = eg.JWTScope(s)
	}
	return result
}

// Generate Backend resource for a given backend endpoint of an API
func generateBackend(name string, backendEndpoint model.EndpointConfiguration) (*eg.Backend, error) {
	_, host, port, err := extractSchemeHostPort(backendEndpoint.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to extract scheme, host and port from endpoint %v: %w", backendEndpoint.Endpoint, err)
	}
	backend := &eg.Backend{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindBackend,
			APIVersion: constantscommon.EnvoyGatewayV1Alpha1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: eg.BackendSpec{
			Endpoints: []eg.BackendEndpoint{
				{
					FQDN: &eg.FQDNEndpoint{
						Hostname: host,
						Port:     int32(port),
					},
				},
			},
		},
	}
	return backend, nil

}

func extractSchemeHostPort(endpoint interface{}) (scheme, host string, port int, err error) {
	endpointType, err := endpointType(endpoint)
	switch endpointType {
	case endpointTypeString:
		endpointStr, _ := endpoint.(string)
		parsed, parseErr := url.Parse(endpointStr)
		if parseErr != nil {
			return "", "", 0, parseErr
		}
		scheme = parsed.Scheme
		host = parsed.Hostname()

		// Determine port
		if parsed.Port() != "" {
			fmt.Sscanf(parsed.Port(), "%d", &port)
		} else {
			if scheme == "http" {
				port = 80
			} else if scheme == "https" {
				port = 443
			}
		}
	case endpointTypeK8sService:
		endpointK8s, _ := endpoint.(model.K8sService)
		// If any field is nil, fill defaults
		if endpointK8s.Protocol != nil {
			scheme = *endpointK8s.Protocol
		} else {
			scheme = "http"
		}
		if endpointK8s.Name != nil && endpointK8s.Namespace != nil {
			host = fmt.Sprintf("%s.%s.svc.cluster.local", *endpointK8s.Name, *endpointK8s.Namespace)
		}
		if endpointK8s.Port != nil {
			port = *endpointK8s.Port
		} else {
			if scheme == "http" {
				port = 80
			} else if scheme == "https" {
				port = 443
			}
		}
	default:
		return "", "", 0, fmt.Errorf("unsupported endpoint type: %T", endpoint)
	}

	return scheme, host, port, nil
}

func generateBackendTLSPolicyWithWellKnownCerts(name, host string,
	targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName) *gwapiv1a3.BackendTLSPolicy {
	wellKnownCerts := gwapiv1a3.WellKnownCACertificatesSystem
	return &gwapiv1a3.BackendTLSPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindBackendTLSPolicy,
			APIVersion: constantscommon.K8sGroupNetworkingv1Alpha3,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: gwapiv1a3.BackendTLSPolicySpec{
			TargetRefs: targetRefs,
			Validation: gwapiv1a3.BackendTLSPolicyValidation{
				WellKnownCACertificates: &wellKnownCerts,
				Hostname:                gatewayv1.PreciseHostname(host),
			},
		},
	}
}

var dnsSafeRegex = regexp.MustCompile(`[^a-z0-9-]`)

func endpointCRName(endpoint interface{}) (string, error) {
	var base string

	endpointType, _ := endpointType(endpoint)
	switch endpointType {
	case endpointTypeString:
		endpointStr, _ := endpoint.(string)
		// Lowercase and replace invalid chars
		base = "backend-" + strings.ToLower(endpointStr)
	case endpointTypeK8sService:
		endpointK8s, _ := endpoint.(model.K8sService)
		name := ""
		if endpointK8s.Name != nil {
			name = *endpointK8s.Name
		}
		namespace := ""
		if endpointK8s.Namespace != nil {
			namespace = *endpointK8s.Namespace
		}
		port := ""
		if endpointK8s.Port != nil {
			port = fmt.Sprintf("%d", *endpointK8s.Port)
		}
		protocol := ""
		if endpointK8s.Protocol != nil {
			protocol = *endpointK8s.Protocol
		}
		base = fmt.Sprintf("k8s-%s-%s-%s-%s", name, namespace, port, protocol)
	default:
		return "", fmt.Errorf("unsupported endpoint type: %T", endpoint)
	}

	// Sanitize: only lowercase letters, numbers, dash
	base = strings.ToLower(base)
	base = dnsSafeRegex.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")

	// Ensure length <= 63
	if len(base) > 63 {
		// Shorten with hash suffix
		h := sha1.Sum([]byte(base))
		hash := hex.EncodeToString(h[:])[:8]
		base = base[:63-len(hash)-1] + "-" + hash
	}

	return base, nil
}

func createBackendRefs(ecs []model.EndpointConfiguration, backendMap map[string]map[string]*eg.Backend,
	routeName string) []gatewayv1.HTTPBackendRef {
	backendRefs := make([]gatewayv1.HTTPBackendRef, 0, len(ecs))
	for _, ec := range ecs {
		if ec.Endpoint == nil {
			continue
		}
		endpointType, err := endpointType(ec.Endpoint)
		if err != nil {
			logger.Sugar().Errorf("Failed to determine endpoint type for %v: %v", ec.Endpoint, err)
			continue
		}

		if endpointType == endpointTypeK8sService {
			endpoint, err := convertMapToK8sService(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to convert endpoint %v to K8sService: %v", ec.Endpoint, err)
				continue
			}
			if endpoint.Name == nil {
				continue
			}
			serviceObjectName := gatewayv1.ObjectName(*endpoint.Name)
			httpBackendRef := gatewayv1.HTTPBackendRef{
				BackendRef: gatewayv1.BackendRef{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Name: serviceObjectName,
						Port: ptrTo(gatewayv1.PortNumber(*endpoint.Port)),
					},
				},
			}
			if ec.Weight != nil {
				weight32 := int32(*ec.Weight)
				httpBackendRef.Weight = &weight32
			}
			backendRefs = append(backendRefs, httpBackendRef)
		}
		if endpointType == endpointTypeString {
			scheme, _, _, err := extractSchemeHostPort(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to extract scheme, host and port from endpoint %v: %v", ec.Endpoint, err)
				continue
			}

			if backendMap[scheme] == nil {
				backendMap[scheme] = make(map[string]*eg.Backend)
			}

			// Generate a unique ID for the backend
			backendID, err := endpointCRName(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to generate CR name for endpoint %v: %v", ec.Endpoint, err)
				continue
			}
			if backendMap[scheme][backendID] == nil {
				// Create a new backend object
				encodedBackendID := hex.EncodeToString(sha1.New().Sum([]byte(backendID)))
				backendName := fmt.Sprintf("%s-%s", routeName, encodedBackendID)
				if len(backendName) > 253 {
					backendName = backendName[:253]
				}
				backend, err := generateBackend(backendName, ec)
				if err != nil {
					logger.Sugar().Errorf("Failed to generate backend for endpoint %v: %v", ec.Endpoint, err)
					continue
				}
				backendMap[scheme][backendID] = backend
			}
			httpBackendRef := gatewayv1.HTTPBackendRef{
				BackendRef: gatewayv1.BackendRef{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Group: ptrTo(gatewayv1.Group(constantscommon.EnvoyGateway)),
						Kind:  ptrTo(gatewayv1.Kind(constantscommon.KindBackend)),
						Name:  gatewayv1.ObjectName(backendMap[scheme][backendID].Name),
					},
				},
			}
			backendRefs = append(backendRefs, httpBackendRef)

		}
	}
	return backendRefs
}

func createGRPCBackendRefs(ecs []model.EndpointConfiguration, backendMap map[string]map[string]*eg.Backend,
	routeName string) []gatewayv1.GRPCBackendRef {
	backendRefs := make([]gatewayv1.GRPCBackendRef, 0, len(ecs))
	for _, ec := range ecs {
		if ec.Endpoint == nil {
			continue
		}
		endpointType, err := endpointType(ec.Endpoint)
		if err != nil {
			logger.Sugar().Errorf("Failed to determine endpoint type for %v: %v", ec.Endpoint, err)
			continue
		}

		if endpointType == endpointTypeK8sService {
			endpoint, err := convertMapToK8sService(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to convert endpoint %v to K8sService: %v", ec.Endpoint, err)
				continue
			}
			if endpoint.Name == nil {
				continue
			}
			serviceObjectName := gatewayv1.ObjectName(*endpoint.Name)
			grpcBackendRef := gatewayv1.GRPCBackendRef{
				BackendRef: gatewayv1.BackendRef{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Name: serviceObjectName,
						Port: ptrTo(gatewayv1.PortNumber(*endpoint.Port)),
					},
				},
			}
			if ec.Weight != nil {
				weight32 := int32(*ec.Weight)
				grpcBackendRef.Weight = &weight32
			}
			backendRefs = append(backendRefs, grpcBackendRef)
		}
		if endpointType == endpointTypeString {
			scheme, _, _, err := extractSchemeHostPort(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to extract scheme, host and port from endpoint %v: %v", ec.Endpoint, err)
				continue
			}

			if backendMap[scheme] == nil {
				backendMap[scheme] = make(map[string]*eg.Backend)
			}

			// Generate a unique ID for the backend
			backendID, err := endpointCRName(ec.Endpoint)
			if err != nil {
				logger.Sugar().Errorf("Failed to generate CR name for endpoint %v: %v", ec.Endpoint, err)
				continue
			}
			if backendMap[scheme][backendID] == nil {
				// Create a new backend object
				backendName := fmt.Sprintf("%s-%s", routeName, hex.EncodeToString(sha1.New().Sum([]byte(backendID)))[:8])
				backend, err := generateBackend(backendName, ec)
				if err != nil {
					logger.Sugar().Errorf("Failed to generate backend for endpoint %v: %v", ec.Endpoint, err)
					continue
				}
				backendMap[scheme][backendID] = backend
			}
			grpcBackendRef := gatewayv1.GRPCBackendRef{
				BackendRef: gatewayv1.BackendRef{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Group: ptrTo(gatewayv1.Group(constantscommon.EnvoyGateway)),
						Kind:  ptrTo(gatewayv1.Kind(constantscommon.KindBackend)),
						Name:  gatewayv1.ObjectName(backendMap[scheme][backendID].Name),
					},
				},
			}
			backendRefs = append(backendRefs, grpcBackendRef)

		}
	}
	return backendRefs
}

const (
	endpointTypeString     = "string"
	endpointTypeK8sService = "k8sservice"
)

func endpointType(endpoint interface{}) (string, error) {
	switch endpoint.(type) {
	case string:
		return endpointTypeString, nil
	case model.K8sService:
		return endpointTypeK8sService, nil
	default:
		// Try to cast map[string]interface{} to K8sService
		if endpointMap, ok := endpoint.(map[string]interface{}); ok {
			// Check if it has the typical K8s service fields
			if _, hasName := endpointMap["name"]; hasName {
				return endpointTypeK8sService, nil
			}
			if _, hasNamespace := endpointMap["namespace"]; hasNamespace {
				return endpointTypeK8sService, nil
			}
			if _, hasPort := endpointMap["port"]; hasPort {
				return endpointTypeK8sService, nil
			}
			if _, hasProtocol := endpointMap["protocol"]; hasProtocol {
				return endpointTypeK8sService, nil
			}
		}
		return "", fmt.Errorf("unsupported endpoint type: %T", endpoint)
	}
}

const (
	securityTypeBasic  = "basic"
	securityTypeAPIKey = "apikey"
)

func securityType(securityType interface{}) (string, error) {
	switch securityType.(type) {
	case *model.BasicEndpointSecurity:
		return securityTypeBasic, nil
	case *model.APIKeyEndpointSecurity:
		return securityTypeAPIKey, nil
	default:
		// Try to cast map[string]interface{} to determine security type
		if securityMap, ok := securityType.(map[string]interface{}); ok {
			// Check for Basic security fields
			if _, hasUserNameKey := securityMap["userNameKey"]; hasUserNameKey {
				return securityTypeBasic, nil
			}
			if _, hasPasswordKey := securityMap["passwordKey"]; hasPasswordKey {
				return securityTypeBasic, nil
			}
			// Check for API Key security fields
			if _, hasSecretName := securityMap["secretName"]; hasSecretName {
				return securityTypeAPIKey, nil
			}
			if _, hasAPIKeyNameKey := securityMap["apiKeyNameKey"]; hasAPIKeyNameKey {
				return securityTypeAPIKey, nil
			}
			if _, hasAPIKeyValueKey := securityMap["apiKeyValueKey"]; hasAPIKeyValueKey {
				return securityTypeAPIKey, nil
			}
			if _, hasIn := securityMap["in"]; hasIn {
				return securityTypeAPIKey, nil
			}
		}
		return "", fmt.Errorf("unsupported security type: %T", securityType)
	}
}

// convertMapToBasicEndpointSecurity converts a map representation of BasicEndpointSecurity to the struct.
func convertMapToBasicEndpointSecurity(securityMap map[string]interface{}) *model.BasicEndpointSecurity {
	basicSecurity := &model.BasicEndpointSecurity{}

	if secretName, exists := securityMap["secretName"]; exists {
		if secretNameStr, ok := secretName.(string); ok {
			basicSecurity.SecretName = secretNameStr
		}
	}

	if userNameKey, exists := securityMap["userNameKey"]; exists {
		if apiKeyNameKeyStr, ok := userNameKey.(string); ok {
			basicSecurity.UserNameKey = apiKeyNameKeyStr
		}
	}

	if passwordKey, exists := securityMap["passwordKey"]; exists {
		if apiKeyValueKeyStr, ok := passwordKey.(string); ok {
			basicSecurity.PasswordKey = apiKeyValueKeyStr
		}
	}

	return basicSecurity
}

// generateBasicSecurityHash generates a SHA256 hash for the given APIKeyEndpointSecurity configuration.
func generateBasicSecurityHash(apiKeySecurity *model.BasicEndpointSecurity) string {
	if apiKeySecurity == nil {
		return ""
	}

	// Concatenate all fields in a consistent order
	hashInput := fmt.Sprintf("%s|%s|%s",
		apiKeySecurity.SecretName,
		apiKeySecurity.UserNameKey,
		apiKeySecurity.PasswordKey,
	)

	// Generate SHA256 hash
	hasher := sha256.New()
	hasher.Write([]byte(hashInput))
	hashBytes := hasher.Sum(nil)

	// Return hex encoded hash
	return hex.EncodeToString(hashBytes)
}

// createBasicEndpointSecurityMediationPolicy creates a Mediation policy for API Key security.
func createBasicEndpointSecurityMediationPolicy(name string, endpointSecurity *model.EndpointSecurity,
	basicSecurity *model.BasicEndpointSecurity) *dpv2alpha1.RoutePolicy {
	basicAuthMediationPolicy := &dpv2alpha1.Mediation{
		PolicyName:    constantscommon.MediationBackendBasicSecurity,
		PolicyID:      "",
		PolicyVersion: "",
		Parameters: []*dpv2alpha1.Parameter{
			{
				Key:   "Enabled",
				Value: strconv.FormatBool(*endpointSecurity.Enabled),
			},
			{
				Key:   "UserNameKey",
				Value: basicSecurity.UserNameKey,
			},
			{
				Key:   "PasswordKey",
				Value: basicSecurity.PasswordKey,
			},
			{
				Key: "BasicAuth",
				ValueRef: &gwapiv1a2.LocalObjectReference{
					Name: gwapiv1a2.ObjectName(basicSecurity.SecretName),
					Kind: constantscommon.KindSecret,
				},
			},
		},
	}
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  []*dpv2alpha1.Mediation{basicAuthMediationPolicy},
			ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
		},
	}
	return routePolicy
}

// convertMapToAPIKeyEndpointSecurity converts a map representation of APIKeyEndpointSecurity to the struct.
func convertMapToAPIKeyEndpointSecurity(securityMap map[string]interface{}) *model.APIKeyEndpointSecurity {
	apiKeySecurity := &model.APIKeyEndpointSecurity{}

	if secretName, exists := securityMap["secretName"]; exists {
		if secretNameStr, ok := secretName.(string); ok {
			apiKeySecurity.SecretName = secretNameStr
		}
	}

	if apiKeyNameKey, exists := securityMap["apiKeyNameKey"]; exists {
		if apiKeyNameKeyStr, ok := apiKeyNameKey.(string); ok {
			apiKeySecurity.APIKeyNameKey = apiKeyNameKeyStr
		}
	}

	if apiKeyValueKey, exists := securityMap["apiKeyValueKey"]; exists {
		if apiKeyValueKeyStr, ok := apiKeyValueKey.(string); ok {
			apiKeySecurity.APIKeyValueKey = apiKeyValueKeyStr
		}
	}

	if in, exists := securityMap["in"]; exists {
		if inStr, ok := in.(string); ok {
			apiKeySecurity.In = inStr
		}
	}

	return apiKeySecurity
}

// generateAPIKeySecurityHash generates a SHA256 hash for the given APIKeyEndpointSecurity configuration.
func generateAPIKeySecurityHash(apiKeySecurity *model.APIKeyEndpointSecurity) string {
	if apiKeySecurity == nil {
		return ""
	}

	// Concatenate all fields in a consistent order
	hashInput := fmt.Sprintf("%s|%s|%s|%s",
		apiKeySecurity.SecretName,
		apiKeySecurity.In,
		apiKeySecurity.APIKeyNameKey,
		apiKeySecurity.APIKeyValueKey,
	)

	// Generate SHA256 hash
	hasher := sha256.New()
	hasher.Write([]byte(hashInput))
	hashBytes := hasher.Sum(nil)

	// Return hex encoded hash
	return hex.EncodeToString(hashBytes)
}

// createAPIKeyMediationPolicy creates a Mediation policy for API Key security.
func createAPIKeyMediationPolicy(name string, endpointSecurity *model.EndpointSecurity,
	apiKeySecurity *model.APIKeyEndpointSecurity) *dpv2alpha1.RoutePolicy {
	apiKeyMediationPolicy := &dpv2alpha1.Mediation{
		PolicyName:    constantscommon.MediationBackendAPIKey,
		PolicyID:      "",
		PolicyVersion: "",
		Parameters: []*dpv2alpha1.Parameter{
			{
				Key:   "Enabled",
				Value: strconv.FormatBool(*endpointSecurity.Enabled),
			},
			{
				Key:   "In",
				Value: apiKeySecurity.In,
			},
			{
				Key:   "InValue",
				Value: apiKeySecurity.APIKeyNameKey,
			},
			{
				Key: "APIKey",
				ValueRef: &gwapiv1a2.LocalObjectReference{
					Name: gwapiv1a2.ObjectName(apiKeySecurity.SecretName),
					Kind: constantscommon.KindSecret,
				},
			},
		},
	}
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  []*dpv2alpha1.Mediation{apiKeyMediationPolicy},
			ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
		},
	}
	return routePolicy
}

// createBackendJWTMediationPolicy creates a Mediation policy for API Key security.
func createBackendJWTMediationPolicy(jwtSecurityPolicy *model.BackendJWTPolicy, kms []model.KeyManager) *dpv2alpha1.RoutePolicy {
	customClaimsJSON := convertCustomClaimsToJSON(jwtSecurityPolicy.Parameters.CustomClaims)
	claimMappingsJSON := convertClaimMappingsToJSON(kms)
	apiKeyMediationPolicy := &dpv2alpha1.Mediation{
		PolicyName:    constantscommon.MediationBackendJWT,
		PolicyID:      "",
		PolicyVersion: "",
		Parameters: []*dpv2alpha1.Parameter{
			{
				Key:   "Enabled",
				Value: "true",
			},
			{
				Key:   "Encoding",
				Value: *jwtSecurityPolicy.Parameters.Encoding,
			},
			{
				Key:   "Header",
				Value: *jwtSecurityPolicy.Parameters.Header,
			},
			{
				Key:   "SigningAlgorithm",
				Value: *jwtSecurityPolicy.Parameters.SigningAlgorithm,
			},
			{
				Key:   "TokenTTL",
				Value: strconv.Itoa(*jwtSecurityPolicy.Parameters.TokenTTL),
			},
			{
				Key:   "UseKid",
				Value: "false",
			},
			{
				Key:   "CustomClaims",
				Value: customClaimsJSON,
			},
			{
				Key:   "ClaimMapping",
				Value: claimMappingsJSON,
			},
		},
	}
	name := util.GeneratePolicyHash(jwtSecurityPolicy)
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  []*dpv2alpha1.Mediation{apiKeyMediationPolicy},
			ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
		},
	}
	return routePolicy
}

// convertCustomClaimsToJSON converts a slice of CustomClaims to a JSON string representation
func convertCustomClaimsToJSON(customClaims []model.CustomClaims) string {
	claimsMap := make(map[string]interface{})
	if customClaims != nil && len(customClaims) > 0 {
		for _, customClaim := range customClaims {
			claimValueType := strings.ToLower(customClaim.Type)
			claimValue := customClaim.Value
			claimsMap[customClaim.Claim] = customClaim.Value
			switch claimValueType {
			case "string":
				claimsMap[customClaim.Claim] = claimValue
			case "int":
				intValue, err := strconv.Atoi(claimValue)
				if err == nil {
					claimsMap[customClaim.Claim] = intValue
				}
			case "boolean":
				boolValue, err := strconv.ParseBool(claimValue)
				if err == nil {
					claimsMap[customClaim.Claim] = boolValue
				}
			case "float":
				floatValue, err := strconv.ParseFloat(claimValue, 64)
				if err == nil {
					claimsMap[customClaim.Claim] = floatValue
				}
			}
		}
	}
	customClaimsJSON, err := json.Marshal(claimsMap)
	if err != nil {
		return "{}"
	}
	return string(customClaimsJSON)
}

// convertClaimMappingsToJSON converts a slice of CustomClaims to a map for claim mappings
func convertClaimMappingsToJSON(kms []model.KeyManager) string {
	claimMappings := make(map[string]string)
	if kms != nil && len(kms) > 0 {
		for _, km := range kms {
			if km.ClaimMapping != nil && len(km.ClaimMapping) > 0 {
				for _, claimMapping := range km.ClaimMapping {
					claimMappings[claimMapping.LocalClaim] = claimMapping.RemoteClaim
				}
			}
		}
	}
	claimMappingsJSON, err := json.Marshal(claimMappings)
	if err != nil {
		return "{}"
	}
	return string(claimMappingsJSON)
}

// generateModelBasedRoundRobinPolicy generates a RoutePolicy for model-based round-robin load balancing.
func generateModelBasedRoundRobinPolicy(policy *model.ModelBasedRoundRobinPolicy, environment string) (*dpv2alpha1.RoutePolicy, error) {
	var modelRouting []model.ModelRouting
	if environment == constants.PRODUCTION_TYPE {
		modelRouting = policy.Parameters.ProductionModels
	} else {
		modelRouting = policy.Parameters.SandboxModels
	}
	modelClusterPairs, err := generateModelClusterPairs(modelRouting)
	if err != nil {
		logger.Sugar().Errorf("Failed to generate model-cluster pairs: %v", err)
		return nil, err
	}
	modelBasedRoundRobinPolicy := &dpv2alpha1.Mediation{
		PolicyName:    constantscommon.MediationAIModelBasedRoundRobin,
		PolicyID:      "",
		PolicyVersion: "",
		Parameters: []*dpv2alpha1.Parameter{
			{
				Key:   "Enabled",
				Value: "true",
			},
			{
				Key:   "OnQuotaExceedSuspendDuration",
				Value: strconv.Itoa(policy.Parameters.OnQuotaExceedSuspendDuration),
			},
			{
				Key:   "ModelsClusterPair",
				Value: modelClusterPairs,
			},
		},
	}
	name := util.GeneratePolicyHash(policy)
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  []*dpv2alpha1.Mediation{modelBasedRoundRobinPolicy},
			ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
		},
	}
	return routePolicy, nil
}

// generateModelClusterPairs generates an array of model-cluster pairs.
func generateModelClusterPairs(routing []model.ModelRouting) (string, error) {
	type ModelClusterPair struct {
		ModelName   string `json:"modelName"`
		ClusterName string `json:"clusterName"`
		Weight      int    `json:"weight"`
	}

	pairs := make([]ModelClusterPair, 0, len(routing))
	for _, route := range routing {
		pair := ModelClusterPair{
			ModelName:   route.Model,
			ClusterName: route.Endpoint,
			Weight:      route.Weight,
		}
		pairs = append(pairs, pair)
	}

	jsonBytes, err := json.Marshal(pairs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal model-cluster pairs to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// generateAIGuardrailPolicy generates a RoutePolicy for AI Guardrail.
func generateAIGuardrailPolicy(policy *model.CommonPolicy, direction string) *dpv2alpha1.RoutePolicy {
	mediationPolicyName := getAIGuardrailPolicyName(policy.PolicyName)
	parameterList := make([]*dpv2alpha1.Parameter, 0, len(policy.Parameters))
	for _, param := range policy.Parameters {
		parameterList = append(parameterList, &dpv2alpha1.Parameter{
			Key:   param.Key,
			Value: param.Value,
		})
	}
	aiGuardrailPolicy := &dpv2alpha1.Mediation{
		PolicyName:    mediationPolicyName,
		PolicyID:      "",
		PolicyVersion: "",
		Parameters:    parameterList,
	}
	name := util.GeneratePolicyHash(policy)
	var requestMediation, responseMediation []*dpv2alpha1.Mediation
	if direction == constantscommon.REQUEST_FLOW {
		requestMediation = []*dpv2alpha1.Mediation{aiGuardrailPolicy}
		responseMediation = make([]*dpv2alpha1.Mediation, 0)
	} else {
		requestMediation = make([]*dpv2alpha1.Mediation, 0)
		responseMediation = []*dpv2alpha1.Mediation{aiGuardrailPolicy}
	}
	routePolicy := &dpv2alpha1.RoutePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       constantscommon.KindRoutePolicy,
			APIVersion: constantscommon.WSO2KubernetesGatewayV2Alpha1,
		},
		Spec: dpv2alpha1.RoutePolicySpec{
			RequestMediation:  requestMediation,
			ResponseMediation: responseMediation,
		},
	}
	return routePolicy
}

// getAIGuardrailPolicyName maps the policy name to the corresponding mediation constant
func getAIGuardrailPolicyName(policyName model.PolicyName) string {
	switch policyName {
	case model.PolicyNameWordCountGuardrail:
		return constantscommon.MediationWordCountGuardrail
	case model.PolicyNameSentenceCountGuardrail:
		return constantscommon.MediationSentenceCountGuardrail
	case model.PolicyNameContentLengthGuardrail:
		return constantscommon.MediationContentLengthGuardrail
	case model.PolicyNamePIIMaskingGuardrail:
		return constantscommon.MediationPIIMaskingGuardrail
	case model.PolicyNameURLGuardrail:
		return constantscommon.MediationURLGuardrail
	case model.PolicyNameRegexGuardrail:
		return constantscommon.MediationRegexGuardrail
	default:
		return string(policyName)
	}
}

func convertMapToK8sService(endpoint interface{}) (*model.K8sService, error) {
	endpointMap, ok := endpoint.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("endpoint is not a map[string]interface{}")
	}

	service := &model.K8sService{}

	if name, exists := endpointMap["name"]; exists {
		if nameStr, ok := name.(string); ok {
			service.Name = &nameStr
		}
	}

	if namespace, exists := endpointMap["namespace"]; exists {
		if namespaceStr, ok := namespace.(string); ok {
			service.Namespace = &namespaceStr
		}
	}

	if port, exists := endpointMap["port"]; exists {
		if portInt, ok := port.(int); ok {
			service.Port = &portInt
		} else if portFloat, ok := port.(float64); ok {
			portInt := int(portFloat)
			service.Port = &portInt
		}
	}

	if protocol, exists := endpointMap["protocol"]; exists {
		if protocolStr, ok := protocol.(string); ok {
			service.Protocol = &protocolStr
		}
	}

	return service, nil
}

func createConfigMapForDefinition(name, definition string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       constantscommon.KindConfigMap,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"Definition": definition,
		},
	}
}

func createConfigMapForGQlSchema(name, schema string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       constantscommon.KindConfigMap,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"Schema": schema,
		},
	}
}

// extractHTTPRouteNames extracts HTTP route names from the provided objects.
func extractHTTPRouteNames(objectList []client.Object) []string {
	httpRouteNames := make([]string, 0)
	for _, object := range objectList {
		if httpRoute, ok := object.(*gatewayv1.HTTPRoute); ok {
			httpRouteNames = append(httpRouteNames, httpRoute.Name)
		}
	}
	return httpRouteNames
}

// generateHTTPRouteAnnotations generates annotations for HTTP routes based on the provided route names.
func generateHTTPRouteAnnotations(httpRouteNames []string) map[string]string {
	annotations := make(map[string]string)
	if len(httpRouteNames) > 0 {
		currentAnnotationValue := ""
		annotationIndex := 1
		for _, routeName := range httpRouteNames {
			separator := ""
			if currentAnnotationValue != "" {
				separator = ","
			}
			potentialValue := currentAnnotationValue + separator + routeName
			// Check if adding this route name would exceed the limit
			if len(potentialValue) > constants.K8sMaxAnnotationLength {
				annotations[fmt.Sprintf("dp.wso2.com/httproutes_%d", annotationIndex)] = currentAnnotationValue
				annotationIndex++
				currentAnnotationValue = routeName
			} else {
				currentAnnotationValue = potentialValue
			}
		}
		// Add the final annotation
		if currentAnnotationValue != "" {
			annotations[fmt.Sprintf("dp.wso2.com/httproutes_%d", annotationIndex)] = currentAnnotationValue
		}
	}
	return annotations
}
