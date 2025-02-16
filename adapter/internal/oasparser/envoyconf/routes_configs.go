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

package envoyconf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	extAuthService "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	extProcessorv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/any"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	opConstants "github.com/wso2/apk/adapter/internal/operator/constants"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Constants for Rate Limiting
const (
	DescriptorKeyForSubscription = "subscription"
	DescriptorKeyForPolicy       = "policy"
	DescriptorKeyForOrganization = "organization"
	extAuthzFilterName           = "envoy.filters.http.ext_authz"
	extProcFilterName            = "envoy.filters.http.ext_proc"

	descriptorMetadataKeyForSubscription          = "ratelimit:subscription"
	descriptorMetadataKeyForUsagePolicy           = "ratelimit:usage-policy"
	descriptorMetadataKeyForOrganization          = "ratelimit:organization"
	descriptorMetadataKeyForBurstCtrlSubscription = "burstCtrl:subscription"
	descriptorMetadataKeyForBurstCtrlUsagePolicy  = "burstCtrl:usage-policy"
	descriptorMetadataKeyForBurstCtrlOrganization = "burstCtrl:organization"
	// DescriptorKeyForAIRequestTokenCount is the descriptor key for AI request token count ratelimit
	DescriptorKeyForAIRequestTokenCount = "airequesttokencount"
	// DescriptorKeyForAIResponseTokenCount is the descriptor key for AI response token count ratelimit
	DescriptorKeyForAIResponseTokenCount = "airesponsetokencount"
	// DescriptorKeyForAITotalTokenCount is the descriptor key for AI total token count ratelimit
	DescriptorKeyForAITotalTokenCount = "aitotaltokencount"
	// DescriptorKeyForAIRequestCount is the descriptor key for AI request count ratelimit
	DescriptorKeyForAIRequestCount = "airequestcount"
	// DescriptorKeyForAIRequestTokenCountForSubscriptionBasedAIRL is the descriptor key for AI request token count ratelimit
	DescriptorKeyForAIRequestTokenCountForSubscriptionBasedAIRL = "airequesttokencountsubs"
	// DescriptorKeyForAIResponseTokenCountForSubscriptionBasedAIRL is the descriptor key for AI response token count ratelimit
	DescriptorKeyForAIResponseTokenCountForSubscriptionBasedAIRL = "airesponsetokencountsubs"
	// DescriptorKeyForAITotalTokenCountForSubscriptionBasedAIRL is the descriptor key for AI total token count ratelimit
	DescriptorKeyForAITotalTokenCountForSubscriptionBasedAIRL = "aitotaltokencountsubs"
	// DescriptorKeyForAIRequestCountForSubscriptionBasedAIRL is the descriptor key for AI request count ratelimit
	DescriptorKeyForAIRequestCountForSubscriptionBasedAIRL = "airequestcountsubs"
	DynamicMetadataKeyForOrganizationAndAIRLPolicy         = "ratelimit:organization-and-rlpolicy"
	DynamicMetadataKeyForSubscription                      = "ratelimit:subscription"
	DescriptorKeyForAISubscription                         = "subscription"
)

func generateRouteConfig(apiType string, routeName string, method *string, match *routev3.RouteMatch, action *routev3.Route_Route, redirectAction *routev3.Route_Redirect,
	metadata *corev3.Metadata, decorator *routev3.Decorator, typedPerFilterConfig map[string]*anypb.Any,
	requestHeadersToAdd []*corev3.HeaderValueOption, requestHeadersToRemove []string,
	responseHeadersToAdd []*corev3.HeaderValueOption, responseHeadersToRemove []string, authentication *model.Authentication) *routev3.Route {
	cloneTypedPerFilterConfig := cloneTypedPerFilterConfig(typedPerFilterConfig)
	//todo: need to fix it in proper way
	if apiType == constants.REST && (authentication == nil || (authentication != nil && (authentication.Disabled || authentication.Oauth2 == nil)) || (method != nil && strings.ToUpper(*method) == "OPTIONS")) {
		logger.LoggerOasparser.Infof("routename%v", routeName)
		logger.LoggerOasparser.Infof("authentication is nill %v", authentication == nil)
		if authentication != nil {
			logger.LoggerOasparser.Infof("authentication.JWT is nill%v", authentication.JWT == nil)
			logger.LoggerOasparser.Infof("authentication.Oauth2 is nill%v", authentication.Oauth2 == nil)
		}
		delete(cloneTypedPerFilterConfig, EnvoyJWT)
	}
	route := &routev3.Route{
		Name:                 routeName,
		Match:                match,
		Metadata:             metadata,
		Decorator:            decorator,
		TypedPerFilterConfig: cloneTypedPerFilterConfig,

		// headers common to all routes are removed at the Route Configuration level in listener.go
		// x-envoy headers are removed using the SuppressEnvoyHeaders param in http_filters.go
		RequestHeadersToAdd:     requestHeadersToAdd,
		RequestHeadersToRemove:  requestHeadersToRemove,
		ResponseHeadersToAdd:    responseHeadersToAdd,
		ResponseHeadersToRemove: responseHeadersToRemove,
	}

	if redirectAction != nil {
		route.Action = redirectAction
	} else if action != nil {
		route.Action = action
	}

	return route
}

func generateRouteMatch(routeRegex string) *routev3.RouteMatch {
	match := &routev3.RouteMatch{
		PathSpecifier: &routev3.RouteMatch_SafeRegex{
			SafeRegex: &envoy_type_matcherv3.RegexMatcher{
				Regex: routeRegex,
			},
		},
	}
	return match
}

func generateRouteAction(apiType string, routeConfig *model.EndpointConfig, ratelimitCriteria *ratelimitCriteria, mirrorClusterNames []string, isBackendBasedAIRatelimitEnabled bool, descriptorValueForBackendBasedAIRatelimit string, weightedCluster *routev3.WeightedCluster_ClusterWeight, isWeighted bool) (action *routev3.Route_Route) {

	if isWeighted {
		// check if weightedCluster is already in the list
		exists := false
		for i, existingWeightedCluster := range weightedClusters {
			if existingWeightedCluster.Name == weightedCluster.Name {
				if existingWeightedCluster.Weight.GetValue() == weightedCluster.Weight.GetValue() {
					exists = true
				} else {
					// Remove the existing entry with the same name but different weight
					weightedClusters = append(weightedClusters[:i], weightedClusters[i+1:]...)
				}
			}
		}

		// if not existing, add to the list
		if !exists {
			weightedClusters = append(weightedClusters, weightedCluster)
		}
		action = &routev3.Route_Route{
			Route: &routev3.RouteAction{
				HostRewriteSpecifier: &routev3.RouteAction_AutoHostRewrite{
					AutoHostRewrite: &wrapperspb.BoolValue{
						Value: true,
					},
				},
				UpgradeConfigs:    getUpgradeConfig(apiType),
				MaxStreamDuration: getMaxStreamDuration(apiType),
				ClusterSpecifier: &routev3.RouteAction_WeightedClusters{
					WeightedClusters: &routev3.WeightedCluster{
						Clusters: weightedClusters,
					},
				},
			},
		}
	} else {
		action = &routev3.Route_Route{
			Route: &routev3.RouteAction{
				HostRewriteSpecifier: &routev3.RouteAction_AutoHostRewrite{
					AutoHostRewrite: &wrapperspb.BoolValue{
						Value: true,
					},
				},
				UpgradeConfigs:    getUpgradeConfig(apiType),
				MaxStreamDuration: getMaxStreamDuration(apiType),
				ClusterSpecifier: &routev3.RouteAction_ClusterHeader{
					ClusterHeader: clusterHeaderName,
				},
			},
		}
	}

	if routeConfig != nil {
		action.Route.IdleTimeout = durationpb.New(time.Duration(routeConfig.IdleTimeoutInSeconds) * time.Second)
	}

	if routeConfig != nil && routeConfig.RetryConfig != nil {
		retryPolicy := &routev3.RetryPolicy{
			RetryBackOff: &routev3.RetryPolicy_RetryBackOff{
				BaseInterval: durationpb.New(time.Duration(routeConfig.RetryConfig.BaseIntervalInMillis) * time.Millisecond),
			},
			RetryOn:              "retriable-status-codes",
			RetriableStatusCodes: routeConfig.RetryConfig.StatusCodes,
			NumRetries:           &wrapperspb.UInt32Value{Value: uint32(routeConfig.RetryConfig.Count)},
		}
		action.Route.RetryPolicy = retryPolicy
	}
	if routeConfig != nil && routeConfig.TimeoutInMillis != 0 {
		action.Route.Timeout = durationpb.New(time.Duration(routeConfig.TimeoutInMillis) * time.Millisecond)
	}

	if ratelimitCriteria != nil && ratelimitCriteria.level != "" {
		action.Route.RateLimits = generateRateLimitPolicy(ratelimitCriteria)
	}
	if isBackendBasedAIRatelimitEnabled {
		action.Route.RateLimits = append(action.Route.RateLimits, generateBackendBasedAIRatelimit(descriptorValueForBackendBasedAIRatelimit)...)
	}

	// Add request mirroring configurations
	if mirrorClusterNames != nil && len(mirrorClusterNames) > 0 {
		mirrorPolicies := []*routev3.RouteAction_RequestMirrorPolicy{}
		for _, clusterName := range mirrorClusterNames {
			mirrorPolicy := &routev3.RouteAction_RequestMirrorPolicy{
				Cluster: clusterName,
			}
			mirrorPolicies = append(mirrorPolicies, mirrorPolicy)
		}
		action.Route.RequestMirrorPolicies = mirrorPolicies
	}

	return action
}

func generateRequestRedirectRoute(route string, policyParams interface{}) (*routev3.Route_Redirect, error) {
	policyParameters, _ := policyParams.(map[string]interface{})
	scheme, _ := policyParameters[constants.RedirectScheme].(string)
	hostname, _ := policyParameters[constants.RedirectHostname].(string)
	port, _ := policyParameters[constants.RedirectPort].(int)
	statusCode, _ := policyParameters[constants.RedirectStatusCode].(int)
	replaceFullPath, _ := policyParameters[constants.RedirectPath].(string)
	redirectActionStatusCode := mapStatusCodeToEnum(statusCode)
	if redirectActionStatusCode == -1 {
		return nil, fmt.Errorf("Invalid status code provided")
	}

	action := &routev3.Route_Redirect{
		Redirect: &routev3.RedirectAction{
			SchemeRewriteSpecifier: &routev3.RedirectAction_HttpsRedirect{
				HttpsRedirect: scheme == "https",
			},
			HostRedirect: hostname,
			PortRedirect: uint32(port),
			PathRewriteSpecifier: &routev3.RedirectAction_PathRedirect{
				PathRedirect: replaceFullPath,
			},
			ResponseCode: routev3.RedirectAction_RedirectResponseCode(redirectActionStatusCode),
		},
	}
	return action, nil
}

func mapStatusCodeToEnum(statusCode int) int {
	switch statusCode {
	case 301:
		return 0
	case 302:
		return 1
	default:
		return -1
	}
}

func generateBackendBasedAIRatelimit(descValue string) []*routev3.RateLimit {
	rateLimitForRequestTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForAIRequestTokenCount,
						DescriptorValue: descValue,
					},
				},
			},
		},
	}
	rateLimitForResponseTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForAIResponseTokenCount,
						DescriptorValue: descValue,
					},
				},
			},
		},
	}
	rateLimitForTotalTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForAITotalTokenCount,
						DescriptorValue: descValue,
					},
				},
			},
		},
	}
	rateLimitForRequestCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForAIRequestCount,
						DescriptorValue: descValue,
					},
				},
			},
		},
	}
	return []*routev3.RateLimit{&rateLimitForRequestTokenCount, &rateLimitForResponseTokenCount, &rateLimitForRequestCount, &rateLimitForTotalTokenCount}
}

func generateRateLimitPolicy(ratelimitCriteria *ratelimitCriteria) []*routev3.RateLimit {
	environmentValue := ratelimitCriteria.environment
	if ratelimitCriteria.level != RateLimitPolicyAPILevel && ratelimitCriteria.envType == opConstants.Sandbox {
		environmentValue += "_sandbox"
	}

	rateLimit := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForOrg,
						DescriptorValue: ratelimitCriteria.organizationID,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForEnvironment,
						DescriptorValue: environmentValue,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   DescriptorKeyForPath,
						DescriptorValue: ratelimitCriteria.basePathForRLService,
					},
				},
			},
		},
	}

	if ratelimitCriteria.level == RateLimitPolicyAPILevel {
		rateLimit.Actions = append(rateLimit.Actions, &routev3.RateLimit_Action{
			ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
				GenericKey: &routev3.RateLimit_Action_GenericKey{
					DescriptorKey:   DescriptorKeyForMethod,
					DescriptorValue: DescriptorValueForAPIMethod,
				},
			},
		})

	} else {
		rateLimit.Actions = append(rateLimit.Actions, &routev3.RateLimit_Action{
			ActionSpecifier: &routev3.RateLimit_Action_RequestHeaders_{
				RequestHeaders: &routev3.RateLimit_Action_RequestHeaders{
					DescriptorKey: DescriptorKeyForMethod,
					HeaderName:    DescriptorValueForOperationMethod,
				},
			},
		})
	}

	ratelimits := []*routev3.RateLimit{&rateLimit}
	return ratelimits
}

func generateHTTPMethodMatcher(methodRegex string, sandClusterName string) []*routev3.HeaderMatcher {
	headerMatcher := generateHeaderMatcher(httpMethodHeader, methodRegex)
	headerMatcherArray := []*routev3.HeaderMatcher{headerMatcher}
	return headerMatcherArray
}

func generateQueryParamMatcher(queryParamName, value string) []*routev3.QueryParameterMatcher {
	queryParamMatcher := &routev3.QueryParameterMatcher{
		Name: queryParamName,
		QueryParameterMatchSpecifier: &routev3.QueryParameterMatcher_StringMatch{
			StringMatch: &envoy_type_matcherv3.StringMatcher{
				MatchPattern: &envoy_type_matcherv3.StringMatcher_Exact{
					Exact: value,
				},
			},
		},
	}
	queryParamArray := []*routev3.QueryParameterMatcher{queryParamMatcher}
	return queryParamArray
}

func generateHeaderMatcher(headerName, valueRegex string) *routev3.HeaderMatcher {
	headerMatcherArray := &routev3.HeaderMatcher{
		Name: headerName,
		HeaderMatchSpecifier: &routev3.HeaderMatcher_StringMatch{
			StringMatch: &envoy_type_matcherv3.StringMatcher{
				MatchPattern: &envoy_type_matcherv3.StringMatcher_SafeRegex{
					SafeRegex: &envoy_type_matcherv3.RegexMatcher{
						Regex: "^" + valueRegex + "$",
					},
				},
			},
		},
	}
	return headerMatcherArray
}

func generateRegexMatchAndSubstitute(routePath, endpointResourcePath string,
	pathMatchType gwapiv1.PathMatchType) *envoy_type_matcherv3.RegexMatchAndSubstitute {
	substitutionString := generateSubstitutionString(endpointResourcePath, pathMatchType)
	return &envoy_type_matcherv3.RegexMatchAndSubstitute{
		Pattern: &envoy_type_matcherv3.RegexMatcher{
			Regex: routePath,
		},
		Substitution: substitutionString,
	}
}

// Router configs for Operational Policies

// generateHeaderToAddRouteConfig returns Router config for SET_HEADER
func generateHeaderToAddRouteConfig(policyParams interface{}) (*corev3.HeaderValueOption, error) {
	var paramsToSetHeader map[string]interface{}
	var ok bool
	var headerName, headerValue string
	if paramsToSetHeader, ok = policyParams.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("error while processing policy parameter map. Map: %v", policyParams)
	}
	if headerName, ok = paramsToSetHeader[constants.HeaderName].(string); !ok || strings.TrimSpace(headerName) == "" {
		return nil, errors.New("policy parameter map must include headerName")
	}
	if headerValue, ok = paramsToSetHeader[constants.HeaderValue].(string); !ok || strings.TrimSpace(headerValue) == "" {
		return nil, errors.New("policy parameter map must include headerValue")
	}
	headerToAdd := corev3.HeaderValueOption{
		Header: &corev3.HeaderValue{
			Key:   headerName,
			Value: headerValue,
		},
		AppendAction: *corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD.Enum(),
	}
	return &headerToAdd, nil
}

func generateHeaderToRemoveString(policyParams interface{}) (string, error) {
	var paramsToRemoveHeader map[string]interface{}
	var ok bool
	var requestHeaderToRemove string
	if paramsToRemoveHeader, ok = policyParams.(map[string]interface{}); !ok {
		return "", fmt.Errorf("error while processing policy parameter map. Map: %v", policyParams)
	}
	if requestHeaderToRemove, ok = paramsToRemoveHeader[constants.HeaderName].(string); !ok ||
		requestHeaderToRemove == "" {
		return "", errors.New("policy parameter map must include headerName")
	}
	return requestHeaderToRemove, nil
}

func generateRewritePathRouteConfig(routePath string, policyParams interface{}, pathMatchType gwapiv1.PathMatchType,
	isDefaultVersion bool) (*envoy_type_matcherv3.RegexMatchAndSubstitute, error) {

	var paramsToSetHeader map[string]interface{}
	var ok bool
	var rewritePath string
	var rewritePathType gwapiv1.HTTPPathModifierType
	if paramsToSetHeader, ok = policyParams.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("error while processing policy parameter map. Map: %v", policyParams)
	}
	if rewritePath, ok = paramsToSetHeader[constants.RewritePathResourcePath].(string); !ok {
		return nil, errors.New("policy parameter map must include rewritePath")
	}
	if rewritePathType, ok = paramsToSetHeader[constants.RewritePathType].(gwapiv1.HTTPPathModifierType); !ok ||
		string(rewritePathType) == "" {
		return nil, errors.New("policy parameter map must include rewritePathType")
	}

	substitutionString := generateSubstitutionStringWithRewritePathType(rewritePath,
		pathMatchType, rewritePathType, isDefaultVersion)

	return &envoy_type_matcherv3.RegexMatchAndSubstitute{
		Pattern: &envoy_type_matcherv3.RegexMatcher{
			Regex: routePath,
		},
		Substitution: substitutionString,
	}, nil
}

func generateSubstitutionStringWithRewritePathType(rewritePath string,
	pathMatchType gwapiv1.PathMatchType, rewritePathType gwapiv1.HTTPPathModifierType, isDefaultVersion bool) string {
	var resourceRegex string
	switch pathMatchType {
	case gwapiv1.PathMatchExact:
		resourceRegex = rewritePath
	case gwapiv1.PathMatchPathPrefix:
		switch rewritePathType {
		case gwapiv1.FullPathHTTPPathModifier:
			resourceRegex = strings.TrimSuffix(rewritePath, "/")
		case gwapiv1.PrefixMatchHTTPPathModifier:
			pathPrefix := "%s\\1"
			resourceRegex = fmt.Sprintf(pathPrefix, strings.TrimSuffix(rewritePath, "/"))
		}
	case gwapiv1.PathMatchRegularExpression:
		resourceRegex = rewritePath
	}
	return resourceRegex
}

func generateFilterConfigToSkipEnforcer() map[string]*anypb.Any {
	perFilterConfig := extAuthService.ExtAuthzPerRoute{
		Override: &extAuthService.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}

	data, _ := proto.Marshal(&perFilterConfig)
	filter := &any.Any{
		TypeUrl: extAuthzPerRouteName,
		Value:   data,
	}
	perFilterConfigExtProc := extProcessorv3.ExtProcPerRoute{
		Override: &extProcessorv3.ExtProcPerRoute_Disabled{
			Disabled: true,
		},
	}

	dataExtProc, _ := proto.Marshal(&perFilterConfigExtProc)
	filterExtProc := &any.Any{
		TypeUrl: extProcPerRouteName,
		Value:   dataExtProc,
	}

	return map[string]*any.Any{
		wellknown.HTTPExternalAuthorization: filter,
		HTTPExternalProcessor:               filterExtProc,
	}
}

func generateMetadataMatcherForInternalRoutes(metadataValue string) (dynamicMetadata []*envoy_type_matcherv3.MetadataMatcher) {
	path := &envoy_type_matcherv3.MetadataMatcher_PathSegment{
		Segment: &envoy_type_matcherv3.MetadataMatcher_PathSegment_Key{
			Key: methodRewrite,
		},
	}
	metadataMatcher := &envoy_type_matcherv3.MetadataMatcher{
		Filter: wellknown.HTTPExternalAuthorization,
		Path:   []*envoy_type_matcherv3.MetadataMatcher_PathSegment{path},
		Value: &envoy_type_matcherv3.ValueMatcher{
			MatchPattern: &envoy_type_matcherv3.ValueMatcher_StringMatch{
				StringMatch: &envoy_type_matcherv3.StringMatcher{
					MatchPattern: &envoy_type_matcherv3.StringMatcher_Exact{
						Exact: metadataValue,
					},
				},
			},
		},
	}
	return []*envoy_type_matcherv3.MetadataMatcher{
		metadataMatcher,
	}
}

func generateMetadataMatcherForExternalRoutes() (dynamicMetadata []*envoy_type_matcherv3.MetadataMatcher) {
	path := &envoy_type_matcherv3.MetadataMatcher_PathSegment{
		Segment: &envoy_type_matcherv3.MetadataMatcher_PathSegment_Key{
			Key: methodRewrite,
		},
	}
	metadataMatcher := &envoy_type_matcherv3.MetadataMatcher{
		Filter: wellknown.HTTPExternalAuthorization,
		Path:   []*envoy_type_matcherv3.MetadataMatcher_PathSegment{path},
		Value: &envoy_type_matcherv3.ValueMatcher{
			MatchPattern: &envoy_type_matcherv3.ValueMatcher_PresentMatch{
				PresentMatch: true,
			},
		},
		Invert: true,
	}
	return []*envoy_type_matcherv3.MetadataMatcher{
		metadataMatcher,
	}
}

// getRewriteRegexFromPathTemplate returns a regex with capture groups for given rewritePathTemplate.
// It replaces {uri.var.petId} included in rewritePath of the path rewrite policy
// with indexes such as \1 \2 that are expected in the substitution string
func getRewriteRegexFromPathTemplate(pathTemplate, rewritePathTemplate string) (string, error) {
	rewriteRegex := "/" + strings.TrimSuffix(strings.TrimPrefix(rewritePathTemplate, "/"), "/")
	pathParamToIndexMap := getPathParamToIndexMap(pathTemplate)
	r := regexp.MustCompile(`{uri.var.([^{}]+)}`) // define a capture group to catch the path param
	matches := r.FindAllStringSubmatch(rewritePathTemplate, -1)
	for _, match := range matches {
		// match is slice always with length two (since one capture group is defined in the regex)
		// hence we do not want to explicitly validate the slice length
		templatedParam := match[0]
		param := match[1]
		if index, ok := pathParamToIndexMap[param]; ok {
			rewriteRegex = strings.ReplaceAll(rewriteRegex, templatedParam, fmt.Sprintf(`\%d`, index))
		} else {
			return "", fmt.Errorf("invalid path param %q in rewrite path", param)
		}
	}

	// validate rewriteRegex
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9~/_.\-\\]*$`, rewriteRegex); !matched {
		logger.LoggerOasparser.Error("Rewrite path includes invalid characters")
		return "", fmt.Errorf("rewrite path regex includes invalid characters, regex %q", rewriteRegex)
	}

	return rewriteRegex, nil
}

// getPathParamToIndexMap returns a map of path params to its index (map of path param -> index)
func getPathParamToIndexMap(pathTemplate string) map[string]int {
	indexMap := make(map[string]int)
	r := regexp.MustCompile(`{([^{}]+)}`) // define a capture group to catch the path param
	matches := r.FindAllStringSubmatch(pathTemplate, -1)
	for i, paramMatches := range matches {
		// paramMatches is slice always with length two (since one capture group is defined in the regex)
		// hence we do not want to explicitly validate the slice length
		indexMap[paramMatches[1]] = i + 1
	}
	return indexMap
}

// cloneTypedPerFilterConfig clones a map[string]*anypb.Any
func cloneTypedPerFilterConfig(original map[string]*anypb.Any) map[string]*anypb.Any {
	clone := make(map[string]*anypb.Any)
	for key, value := range original {
		// Deep copy the value
		clone[key] = &anypb.Any{
			TypeUrl: value.TypeUrl,
			Value:   value.Value,
		}
	}
	return clone
}
