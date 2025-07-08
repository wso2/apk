package mediation

import (
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// MediationAITokenRatelimit holds the name of the AI Token Rate Limit mediation policy.
	MediationAITokenRatelimit = "AITokenRatelimit"
	// MediationSubscriptionRatelimit holds the name of the Subscription Rate Limit mediation policy.
	MediationSubscriptionRatelimit = "SubscriptionRatelimit"
	// MediationSubscriptionValidation holds the name of the Subscription Validation mediation policy.
	MediationSubscriptionValidation = "SubscriptionValidation"
	// MediationAIModelBasedRoundRobin holds the name of the AI Model Based Round Robin mediation policy.
	MediationAIModelBasedRoundRobin = "AIModelBasedRoundRobin"
	// MediationAnalytics holds the name of the Analytics mediation policy.
	MediationAnalytics = "Analytics"
	// MediationBackendJWT holds the name of the Backend JWT mediation policy.
	MediationBackendJWT = "BackendJWT"
	// MediationGraphQL holds the name of the GraphQL mediation policy.
	MediationGraphQL = "GraphQL"
)

// MediationAndRequestHeaderProcessing defines the mediation and request header processing
var MediationAndRequestHeaderProcessing = map[string]bool{
	MediationAITokenRatelimit:       false,
	MediationSubscriptionRatelimit:  true,
	MediationSubscriptionValidation: true,
	MediationAIModelBasedRoundRobin: false,
	MediationAnalytics:              true,
	MediationBackendJWT:             true,
	MediationGraphQL:                false,
}

// MediationAndRequestBodyProcessing defines the mediation and request header processing
var MediationAndRequestBodyProcessing = map[string]bool{
	MediationAITokenRatelimit:       false,
	MediationSubscriptionRatelimit:  false,
	MediationSubscriptionValidation: false,
	MediationAIModelBasedRoundRobin: true,
	MediationAnalytics:              false,
	MediationBackendJWT:             false,
	MediationGraphQL:                true,
}

// MediationAndResponseHeaderProcessing defines the mediation and request header processing
var MediationAndResponseHeaderProcessing = map[string]bool{
	MediationAITokenRatelimit:       false,
	MediationSubscriptionRatelimit:  false,
	MediationSubscriptionValidation: false,
	MediationAIModelBasedRoundRobin: false,
	MediationAnalytics:              false,
	MediationBackendJWT:             false,
	MediationGraphQL:                false,
}

// MediationAndResponseBodyProcessing defines the mediation and response body processing
var MediationAndResponseBodyProcessing = map[string]bool{
	MediationAITokenRatelimit:       true,
	MediationSubscriptionRatelimit:  false,
	MediationSubscriptionValidation: false,
	MediationAIModelBasedRoundRobin: false,
	MediationAnalytics:              false,
	MediationBackendJWT:             false,
	MediationGraphQL:                false,
}

// Result holds the result of mediation processing.
type Result struct {
	AddHeaders                   map[string]string
	RemoveHeaders                []string
	ModifyBody                   bool
	Body                         string
	ImmediateResponse            bool
	ImmediateResponseCode        v32.StatusCode
	ImmediateResponseBody        string
	ImmediateResponseDetail      string
	ImmediateResponseHeaders     map[string]string
	ImmediateResponseContentType string
	StopFurtherProcessing        bool
	Metadata                     map[string]*structpb.Value
}

// Mediation interface defines the methods that all mediation policies must implement.
type Mediation interface {
	Process(*requestconfig.Holder) *Result
}

// CreateMediation creates a Mediation instance based on the provided Mediation object from the cluster.
func CreateMediation(mediationFromCluster *dpv2alpha1.Mediation) Mediation {
	switch mediationFromCluster.PolicyName {
	case MediationAITokenRatelimit:
		return NewAITokenRateLimit(mediationFromCluster)
	case MediationSubscriptionRatelimit:
		return NewSubscriptionRatelimit(mediationFromCluster)
	case MediationSubscriptionValidation:
		return NewSubscriptionValidation(mediationFromCluster)
	case MediationAIModelBasedRoundRobin:
		return NewAIModelBasedRoundRobin(mediationFromCluster)
	case MediationAnalytics:
		return NewAnalytics(mediationFromCluster)
	case MediationBackendJWT:
		return NewBackendJWT(mediationFromCluster)
	case MediationGraphQL:
		return NewGraphQL(mediationFromCluster)
	default:
		return nil
	}
}

func extractPolicyValue(params []dpv2alpha1.Parameter, key string) (string, bool) {
	for _, param := range params {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}
