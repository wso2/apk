package crbuilder

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
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
	var backendBasePath string
	if environment == constants.SANDBOX_TYPE && operation.EndpointConfigurations.Sandbox != nil &&
		len(operation.EndpointConfigurations.Sandbox) > 0 {
		sandboxEndpoint := operation.EndpointConfigurations.Sandbox[0].Endpoint
		fmt.Printf("Sandbox endpoint: %v\n", sandboxEndpoint)
		typeofEndpoint, _ := endpointType(sandboxEndpoint)
		fmt.Printf("endpoint type: %v\n", typeofEndpoint)
		if typeofEndpoint == "string" {
			if parsed, err := url.Parse(sandboxEndpoint.(string)); err == nil {
				backendBasePath = parsed.Path
			}
		}

	} else if environment == constants.PRODUCTION_TYPE && operation.EndpointConfigurations.Production != nil &&
		len(operation.EndpointConfigurations.Production) > 0 {
		productionEndpoint := operation.EndpointConfigurations.Production[0].Endpoint
		fmt.Printf("Prod endpoint: %v\n", productionEndpoint)
		typeofEndpoint, _ := endpointType(productionEndpoint)
		fmt.Printf("endpoint type: %v\n", typeofEndpoint)
		if typeofEndpoint == "string" {
			if parsed, err := url.Parse(productionEndpoint.(string)); err == nil {
				backendBasePath = parsed.Path
			}
		}
	} else {
		return "", fmt.Errorf("no valid endpoint configurations found for environment: %s", environment)
	}
	return backendBasePath, nil
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
			Kind:       constants.WSO2KubernetesGatewayRouteMetadataKind,
			APIVersion: constants.WSO2KubernetesGatewayRouteMetadataAPIVersion,
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
					Kind: gatewayv1.Kind(constants.K8sKindConfigMap),
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
	routePolicies := make([]*dpv2alpha1.RoutePolicy, 0)

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
			Kind:       constants.WSO2KubernetesGatewayRoutePolicyKind,
			APIVersion: constants.WSO2KubernetesGatewayRoutePolicyAPIVersion,
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
				Kind:       constants.WSO2KubernetesGatewayRoutePolicyKind,
				APIVersion: constants.WSO2KubernetesGatewayRoutePolicyAPIVersion,
			},
		}
		routePolicies = append(routePolicies, aiProviderRoutePolicy)
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
						Kind: constants.K8sKindConfigMap,
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
	routePolicies = append(routePolicies, routePolicy)
	objects = append(objects, routePolicy)

	// Create the HTTPRoute objects
	routes, objectsForWithVersion := GenerateHTTPRoutes(apiResourceBundle, true, environment, routePolicies, routeMetadataList)
	for _, obj := range objectsForWithVersion {
		objects = append(objects, obj)
	}
	if apiResourceBundle.APKConf.DefaultVersion {
		routesL, objectsForWithoutVersion := GenerateHTTPRoutes(apiResourceBundle, false, environment, routePolicies, routeMetadataList)
		for _, obj := range objectsForWithoutVersion {
			objects = append(objects, obj)
		}
		for key, value := range routesL {
			routes[key] = append(value, routes[key]...)
		}
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
						Name:  gwapiv1a2.ObjectName(httpRoute.Name),
						Kind:  constants.K8sKindHTTPRoute,
						Group: constants.K8sGroupNetworking,
					},
				})
			}
		}
		// Generate BackendTrafficPolicy for AI Ratelimit
		btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		backendTrafficPolicy := generateBackendTrafficPolicyForAIRatelimit(btpName, targetRefs, aiRatelimit)
		objects = append(objects, backendTrafficPolicy)
	}

	// Ratelimit
	if apiResourceBundle.APKConf.RateLimit != nil {
		var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
		for _, httpRoutes := range routes {
			for _, httpRoute := range httpRoutes {
				targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
						Name:  gwapiv1a2.ObjectName(httpRoute.Name),
						Kind:  constants.K8sKindHTTPRoute,
						Group: constants.K8sGroupNetworking,
					},
				})
			}
		}
		// Generate BackendTrafficPolicy for Ratelimit
		btpName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		backendTrafficPolicy := generateBackendTrafficPolicyForRatelimit(btpName, targetRefs, apiResourceBundle.APKConf.RateLimit)
		objects = append(objects, backendTrafficPolicy)
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
							Name:  gwapiv1a2.ObjectName(httpRoute.Name),
							Kind:  constants.K8sKindHTTPRoute,
							Group: constants.K8sGroupNetworking,
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
			objects = append(objects, backendTrafficPolicy)
		}
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
					Name:  gwapiv1a2.ObjectName(httpRoute.Name),
					Kind:  constants.K8sKindHTTPRoute,
					Group: constants.K8sGroupNetworking,
				},
			})
		}
		spName := util.GenerateCRName(apiResourceBundle.APKConf.Name, environment, apiResourceBundle.APKConf.Version,
			apiResourceBundle.Organization)
		spName = fmt.Sprintf("%s-%d", spName, i+1)
		sp := generateSecurityPolicy(spName, isSecured, scopes, targetRefs, cors, apiResourceBundle.APKConf.KeyManagers,
			apiResourceBundle.APKConf.Authentication)
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

// GenerateHTTPRoutes generates HTTPRoute objects for the given APIResourceBundle.
func GenerateHTTPRoutes(bundle *dto.APIResourceBundle, withVersion bool, environment string, routePolicies []*dpv2alpha1.RoutePolicy,
	routeMetadataList []*dpv2alpha1.RouteMetadata) (map[int][]gatewayv1.HTTPRoute, []client.Object) {
	objects := make([]client.Object, 0)
	routesMap := make(map[int][]gatewayv1.HTTPRoute)
	backendMap := make(map[string]map[string]*eg.Backend)
	crName := util.GenerateCRName(bundle.APKConf.Name, environment, bundle.APKConf.Version, bundle.Organization)
	for i, combined := range bundle.CombinedResources {
		batches := chunkOperations(combined.APKOperations, 16)

		for j, batch := range batches {
			parentName := config.GetConfig().ParentGatewayName
			parentNamespace := config.GetConfig().ParentGatewayNamespace
			parentSectionName := config.GetConfig().ParentGatewaySectionName
			routeName := fmt.Sprintf("%s-%d-%d", crName, i+1, j+1)
			route := gatewayv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "gateway.networking.k8s.io/v1",
					Kind:       "HTTPRoute",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: routeName,
					Annotations: map[string]string{
						constants.K8sHTTPRouteEnvTypeAnnotation: environment,
					},
				},
				Spec: gatewayv1.HTTPRouteSpec{
					Hostnames: []gatewayv1.Hostname{
						gatewayv1.Hostname(func() string {
							gatewayHostName := config.GetConfig().GatewayHostName
							if environment == constants.SANDBOX_TYPE {
								return fmt.Sprintf("sandbox.%s.%s", bundle.Organization, gatewayHostName)
							}
							return fmt.Sprintf("%s.%s", bundle.Organization, gatewayHostName)
						}()),
					},
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name:        gatewayv1.ObjectName(parentName),
								Group:       ptrTo(gatewayv1.Group(constants.K8sGroupNetworking)),
								Kind:        ptrTo(gatewayv1.Kind(constants.K8sKindGateway)),
								Namespace:   ptrTo(gatewayv1.Namespace(parentNamespace)),
								SectionName: ptrTo(gatewayv1.SectionName(parentSectionName)),
							},
						},
					},
					Rules: []gatewayv1.HTTPRouteRule{},
				},
			}

			for _, op := range batch {
				backendBasePath, err := extractBackendBasePath(op, environment)
				if err != nil {
					logger.Sugar().Errorf("Error extracting backend base path for operation %s: %v", *op.Target, err)
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
					strings.TrimPrefix(*op.Target, "/"),
				)
				serviceContractPath := fmt.Sprintf("%s/%s",
					strings.TrimSuffix(backendBasePath, "/"),
					strings.TrimPrefix(*op.Target, "/"),
				)

				isRegexPath, pattern, substitution := GenerateRegexPath(path, backendBasePath, apiBasePath)
				hrfName := ""
				pathMatchType := gatewayv1.PathMatchPathPrefix
				if isRegexPath {
					pathMatchType = gatewayv1.PathMatchRegularExpression
					path = pattern

					sum := sha256.Sum256([]byte(fmt.Sprintf("%s-%s", path, string(*method))))
					pathIdentifier := fmt.Sprintf("%x", sum[:8])
					hrfName = fmt.Sprintf("%s-%s", routeName, pathIdentifier)
					// Create HTTPRouteFilter
					hrf := eg.HTTPRouteFilter{
						TypeMeta: metav1.TypeMeta{
							Kind:       constants.EnvoyGatewayHTTPRouteFilter,
							APIVersion: constants.EnvoyGatewayHTTPRouteFilterAPIVersion,
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
				ecs := op.EndpointConfigurations.Production
				if environment == constants.SANDBOX_TYPE {
					ecs = op.EndpointConfigurations.Sandbox
				}
				if len(ecs) == 0 {
					continue
				}
				// Create backend reference
				httpBackendRefs := createBackendRefs(ecs, backendMap, routeName)
				rule := gatewayv1.HTTPRouteRule{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Type:  ptrTo(pathMatchType),
								Value: ptrTo(path),
							},
							Method: method,
						},
					},
				}
				if len(httpBackendRefs) > 0 {
					rule.BackendRefs = httpBackendRefs
				} else {
					logger.Sugar().Warnf("No backend references found for operation %s in API %s", *op.Target, bundle.APKConf.Name)
				}
				for _, policy := range routePolicies {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constants.WSO2KubernetesGatewayRoutePolicyGroup,
							Kind:  constants.WSO2KubernetesGatewayRoutePolicyKind,
							Name:  gatewayv1.ObjectName(policy.Name),
						},
					})
				}
				for _, metadata := range routeMetadataList {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constants.WSO2KubernetesGatewayRouteMetadataGroup,
							Kind:  constants.WSO2KubernetesGatewayRouteMetadataKind,
							Name:  gatewayv1.ObjectName(metadata.Name),
						},
					})
				}
				if hrfName != "" {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gatewayv1.LocalObjectReference{
							Group: constants.K8sGroupEnvoyGateway,
							Kind:  constants.EnvoyGatewayHTTPRouteFilter,
							Name:  gatewayv1.ObjectName(hrfName),
						},
					})
				} else {
					rule.Filters = append(rule.Filters, gatewayv1.HTTPRouteFilter{
						Type: gatewayv1.HTTPRouteFilterURLRewrite,
						URLRewrite: &gatewayv1.HTTPURLRewriteFilter{
							Path: &gatewayv1.HTTPPathModifier{
								ReplaceFullPath: &serviceContractPath,
								Type:            gatewayv1.FullPathHTTPPathModifier,
							},
						},
					})
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
								Kind:  constants.K8sKindBackend,
								Group: constants.K8sGroupEnvoyGateway,
							},
						},
					})
				objects = append(objects, backendTLSPolicy)
			}
		}
	}
	return routesMap, objects
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
			Kind:       constants.EnvoyGatewayBackendTrafficPolicy,
			APIVersion: constants.EnvoyGatewayBackendTrafficPolicyAPIVersion,
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
			Kind:       constants.EnvoyGatewayBackendTrafficPolicy,
			APIVersion: constants.EnvoyGatewayBackendTrafficPolicyAPIVersion,
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

func generateSecurityPolicy(name string, isSecured bool, scopes []string, targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName,
	cors *model.CORSConfiguration, kms []model.KeyManager, auths []model.AuthenticationRequest) *eg.SecurityPolicy {
	sp := &eg.SecurityPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constants.K8sKindSecurityPolicy,
			APIVersion: constants.K8sAPIVersionEnvoyGateway,
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
			Name: h,
		})
	}
	if isSecured {
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
							Group: ptrTo(gatewayv1.Group(constants.K8sGroupEnvoyGateway)),
							Kind:  ptrTo(gatewayv1.Kind(constants.K8sKindBackend)),
							Name:  gatewayv1.ObjectName(*km.K8sBackend.Name),
							Port:  ptrTo(gatewayv1.PortNumber(*km.K8sBackend.Port)),
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
		switch v := auth.(type) {
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
			Kind:       constants.K8sKindBackend,
			APIVersion: constants.K8sAPIVersionEnvoyGateway,
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
	switch v := endpoint.(type) {
	case string:
		parsed, parseErr := url.Parse(v)
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

	case model.K8sService:
		// If any field is nil, fill defaults
		if v.Protocol != nil {
			scheme = *v.Protocol
		} else {
			scheme = "http"
		}
		if v.Name != nil && v.Namespace != nil {
			host = fmt.Sprintf("%s.%s.svc.cluster.local", *v.Name, *v.Namespace)
		}
		if v.Port != nil {
			port = *v.Port
		} else {
			if scheme == "http" {
				port = 80
			} else if scheme == "https" {
				port = 443
			}
		}

	default:
		return "", "", 0, fmt.Errorf("unsupported endpoint type: %T", v)
	}

	return scheme, host, port, nil
}

func generateBackendTLSPolicyWithWellKnownCerts(name, host string,
	targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName) *gwapiv1a3.BackendTLSPolicy {
	wellKnownCerts := gwapiv1a3.WellKnownCACertificatesSystem
	return &gwapiv1a3.BackendTLSPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       constants.K8sKindBackendTLSPolicy,
			APIVersion: constants.K8sAPIVersionBackendTLSPolicy,
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

func endpointsEqualStrict(e1, e2 interface{}) (bool, error) {
	switch v1 := e1.(type) {
	case string:
		v2, ok := e2.(string)
		if !ok {
			return false, nil // one is string, the other is not
		}
		u1, err := url.Parse(v1)
		if err != nil {
			return false, err
		}
		u2, err := url.Parse(v2)
		if err != nil {
			return false, err
		}
		return u1.Scheme == u2.Scheme && u1.Host == u2.Host && u1.Path == u2.Path, nil

	case model.K8sService:
		v2, ok := e2.(model.K8sService)
		if !ok {
			return false, nil // one is K8sService, the other is not
		}

		// Compare all fields (pointer safe)
		if (v1.Name == nil) != (v2.Name == nil) || (v1.Name != nil && *v1.Name != *v2.Name) {
			return false, nil
		}
		if (v1.Namespace == nil) != (v2.Namespace == nil) || (v1.Namespace != nil && *v1.Namespace != *v2.Namespace) {
			return false, nil
		}
		if (v1.Port == nil) != (v2.Port == nil) || (v1.Port != nil && *v1.Port != *v2.Port) {
			return false, nil
		}
		if (v1.Protocol == nil) != (v2.Protocol == nil) || (v1.Protocol != nil && *v1.Protocol != *v2.Protocol) {
			return false, nil
		}
		return true, nil

	default:
		return false, fmt.Errorf("unsupported endpoint type: %T", v1)
	}
}

var dnsSafeRegex = regexp.MustCompile(`[^a-z0-9-]`)

func endpointCRName(endpoint interface{}) (string, error) {
	var base string

	switch v := endpoint.(type) {
	case string:
		// Lowercase and replace invalid chars
		base = "backend-" + strings.ToLower(v)
	case model.K8sService:
		name := ""
		if v.Name != nil {
			name = *v.Name
		}
		namespace := ""
		if v.Namespace != nil {
			namespace = *v.Namespace
		}
		port := ""
		if v.Port != nil {
			port = fmt.Sprintf("%d", *v.Port)
		}
		protocol := ""
		if v.Protocol != nil {
			protocol = *v.Protocol
		}
		base = fmt.Sprintf("k8s-%s-%s-%s-%s", name, namespace, port, protocol)
	default:
		return "", fmt.Errorf("unsupported endpoint type: %T", v)
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
			endpoint := ec.Endpoint.(model.K8sService)
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
				backendName := fmt.Sprintf("%s-%s", routeName, hex.EncodeToString(sha1.New().Sum([]byte(backendID)))[:8])
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
						Group: ptrTo(gatewayv1.Group(constants.K8sGroupEnvoyGateway)),
						Kind:  ptrTo(gatewayv1.Kind(constants.K8sKindBackend)),
						Name:  gatewayv1.ObjectName(backendMap[scheme][backendID].Name),
					},
				},
			}
			backendRefs = append(backendRefs, httpBackendRef)

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
		return "", fmt.Errorf("unsupported endpoint type: %T", endpoint)
	}
}

func createConfigMapForDefinition(name, definition string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       constants.K8sKindConfigMap,
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
			Kind:       constants.K8sKindConfigMap,
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
