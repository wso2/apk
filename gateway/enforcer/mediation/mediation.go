package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

const (
	MediationAITokenRatelimit       = "AITokenRatelimit"
	MediationSubscriptionRatelimit  = "SubscriptionRatelimit"
	MediationSubscriptionValidation = "SubscriptionValidation"
	MediationAIModelBasedRoundRobin = "AIModelBasedRoundRobin"
	MediationAnalytics              = "Analytics"
	MediationBackendJWT             = "BackendJWT"
	MediationGraphQL                = "GraphQL"
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

type MediationResult struct {
	AddHeaders                   map[string]string
	RemoveHeaders                []string
	ReplaceHeaders               map[string]string
	ModifyBody                   bool
	Body                         string
	ImmediateResponse            bool
	ImmediateResponseCode        int
	ImmediateResponseBody        string
	ImmediateResponseHeaders     map[string]string
	ImmediateResponseContentType string
	StopFurtherProcessing        bool
}

type Mediation interface {
	Process(*requestconfig.Holder) *MediationResult
}

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