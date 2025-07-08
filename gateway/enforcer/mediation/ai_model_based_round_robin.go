package mediation

import (
	"encoding/json"
	"strconv"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
)

type AIModelBasedRoundRobin struct {
	PolicyName                   string
	PolicyVersion                string
	PolicyID                     string
	Enabled                      bool
	OnQuotaExceedSuspendDuration int
	ModelsClusterPair            []ModelClusterPair
}

const (
	MediationAIModelBasedRoundRobinKeyEnabled = "Enabled"
	MediationAIModelBasedRoundRobinKeyOnQuotaExceedSuspendDuration = "OnQuotaExceedSuspendDuration"
	MediationAIModelBasedRoundRobinKeyModelsClusterPair = "ModelsClusterPair"
)

// ModelClusterPair represents a pair of AI model and its associated cluster.
type ModelClusterPair struct {
	ModelName   string `json:"modelName"`
	ClusterName string `json:"clusterName"`
	Weight      int    `json:"weight"`
}

// NewAIModelBasedRoundRobin creates a new AIModelBasedRoundRobin instance with default values.
func NewAIModelBasedRoundRobin(mediation *dpv2alpha1.Mediation) *AIModelBasedRoundRobin {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, MediationAIModelBasedRoundRobinKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	onQuotaExceedSuspendDuration := int(0)
	if val, ok := extractPolicyValue(mediation.Parameters, MediationAIModelBasedRoundRobinKeyOnQuotaExceedSuspendDuration); ok {
		if val != "" {
			duration, err := strconv.Atoi(val)
			if err == nil {
				onQuotaExceedSuspendDuration = duration
			}
		}
	}
	modelsClusterPair := []ModelClusterPair{}
	if val, ok := extractPolicyValue(mediation.Parameters, MediationAIModelBasedRoundRobinKeyModelsClusterPair); ok {
		if val != "" {
			// Assuming val is a JSON string representing an array of ModelClusterPair
			err := json.Unmarshal([]byte(val), &modelsClusterPair)
			if err != nil {
				config.GetConfig().Logger.Sugar().Errorf("Failed to unmarshal ModelsClusterPair: %v, error: %v", val, err)
			}
		}
	}
	return &AIModelBasedRoundRobin{
		PolicyName:                   "AIModelBasedRoundRobin",
		PolicyVersion:                mediation.PolicyVersion,
		PolicyID:                     mediation.PolicyID,
		Enabled:                      enabled,
		OnQuotaExceedSuspendDuration: onQuotaExceedSuspendDuration, // Default to no suspension
		ModelsClusterPair:            modelsClusterPair,
	}
}

// Process
func (a *AIModelBasedRoundRobin) Process(requestconfig *requestconfig.Holder) *MediationResult{
	return &MediationResult{}
}