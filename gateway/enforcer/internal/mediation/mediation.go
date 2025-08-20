package mediation

import (
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	constantscommon "github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// MediationAITokenRatelimit holds the name of the AI Token Rate Limit mediation policy.
	MediationAITokenRatelimit = constantscommon.MediationAITokenRatelimit
	// MediationAIProvider holds the name of the AI Provider mediation policy.
	MediationAIProvider = constantscommon.MediationAIProvider
	// MediationSubscriptionRatelimit holds the name of the Subscription Rate Limit mediation policy.
	MediationSubscriptionRatelimit = constantscommon.MediationSubscriptionRatelimit
	// MediationSubscriptionValidation holds the name of the Subscription Validation mediation policy.
	MediationSubscriptionValidation = constantscommon.MediationSubscriptionValidation
	// MediationAIModelBasedRoundRobin holds the name of the AI Model Based Round Robin mediation policy.
	MediationAIModelBasedRoundRobin = constantscommon.MediationAIModelBasedRoundRobin
	// MediationAnalytics holds the name of the Analytics mediation policy.
	MediationAnalytics = constantscommon.MediationAnalytics
	// MediationBackendJWT holds the name of the Backend JWT mediation policy.
	MediationBackendJWT = constantscommon.MediationBackendJWT
	// MediationGraphQL holds the name of the GraphQL mediation policy.
	MediationGraphQL = constantscommon.MediationGraphQL
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
	MediationAITokenRatelimit:       true,
	MediationAIProvider:             true,
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
	MediationAIProvider:             true,
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

// NewResult creates a new Result instance with default values.
func NewResult() *Result {
	return &Result{
		AddHeaders:               make(map[string]string),
		ImmediateResponseHeaders: make(map[string]string),
		Metadata:                 make(map[string]*structpb.Value),
	}
}

// Mediation interface defines the methods that all mediation policies must implement.
type Mediation interface {
	Process(*requestconfig.Holder) *Result
}

// CreateMediation creates a Mediation instance based on the provided Mediation object from the cluster.
func CreateMediation(mediationFromCluster *dpv2alpha1.Mediation) Mediation {
	if MediationMap[mediationFromCluster] != nil {
		return MediationMap[mediationFromCluster]
	}
	switch mediationFromCluster.PolicyName {
	case MediationAIProvider:
		// Check if the mediation already exists
		mediation := NewAIProvider(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationSubscriptionRatelimit:
		mediation := NewSubscriptionRatelimit(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationSubscriptionValidation:
		mediation := NewSubscriptionValidation(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationAIModelBasedRoundRobin:
		mediation := NewAIModelBasedRoundRobin(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationAnalytics:
		mediation := NewAnalytics(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationBackendJWT:
		mediation := NewBackendJWT(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	case MediationGraphQL:
		mediation := NewGraphQL(mediationFromCluster)
		MediationMap[mediationFromCluster] = mediation
		return mediation
	default:
		return nil
	}
}

func extractPolicyValue(params []*dpv2alpha1.Parameter, key string) (string, bool) {
	for _, param := range params {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

// MediationMap holds the mapping of Mediation objects to their corresponding Mediation instances.
var MediationMap = make(map[*dpv2alpha1.Mediation]Mediation)

// DeleteMediation removes a Mediation from the MediationMap.
func DeleteMediation(mediation *dpv2alpha1.Mediation) {
	if _, exists := MediationMap[mediation]; exists {
		delete(MediationMap, mediation)
	}
}
