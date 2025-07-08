package mediation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	datastore "github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// AIModelBasedRoundRobin represents the configuration for AI Model Based Round Robin policy in the API Gateway.
type AIModelBasedRoundRobin struct {
	PolicyName                   string
	PolicyVersion                string
	PolicyID                     string
	Enabled                      bool
	OnQuotaExceedSuspendDuration int
	ModelsClusterPair            []ModelClusterPair
	logger                      *logging.Logger
}

const (
	// MediationAIModelBasedRoundRobinKeyEnabled is the key for enabling/disabling the AI Model Based Round Robin policy.
	MediationAIModelBasedRoundRobinKeyEnabled = "Enabled"
	// MediationAIModelBasedRoundRobinKeyOnQuotaExceedSuspendDuration is the key for specifying the duration to suspend on quota exceed.
	MediationAIModelBasedRoundRobinKeyOnQuotaExceedSuspendDuration = "OnQuotaExceedSuspendDuration"
	// MediationAIModelBasedRoundRobinKeyModelsClusterPair is the key for specifying the AI model and its associated cluster pairs.
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
		logger:                      &config.GetConfig().Logger,
	}
}

// Process processes the request configuration for AI Model Based Round Robin.
func (a *AIModelBasedRoundRobin) Process(requestConfigHolder *requestconfig.Holder) *Result {
	result := &Result{}
	if requestConfigHolder.ProcessingPhase == requestconfig.ProcessingPhaseRequestHeaders {
		var modelWeights []datastore.ModelWeight
		for _, model := range a.ModelsClusterPair {
			modelWeights = append(modelWeights, datastore.ModelWeight{
				Name:     model.ModelName,
				Endpoint: model.ClusterName,
				Weight:   model.Weight,
			})
		}
		a.logger.Sugar().Debugf(fmt.Sprintf("Supported Models: %+v", a.ModelsClusterPair))
		a.logger.Sugar().Debugf(fmt.Sprintf("Model Weights: %+v", modelWeights))
		a.logger.Sugar().Debugf(fmt.Sprintf("On Quota Exceed Suspend Duration: %v", a.OnQuotaExceedSuspendDuration))
		selectedModel, selectedEndpoint := datastore.GetModelBasedRoundRobinTracker().GetNextModel(requestConfigHolder.RouteMetadata.Spec.API.Name, requestConfigHolder.RouteName, modelWeights)
		a.logger.Sugar().Debug(fmt.Sprintf("Selected Model: %v", selectedModel))
		a.logger.Sugar().Debug(fmt.Sprintf("Selected Endpoint: %v", selectedEndpoint))
		if selectedModel == "" || selectedEndpoint == "" {
			a.logger.Sugar().Debug("Unable to select a model since all models are suspended. Continue with the user provided model")
		} else {
			// change request body to model to selected model
			httpBody := requestConfigHolder.RequestBody.Body
			a.logger.Sugar().Debug(fmt.Sprintf("request body before %+v\n", httpBody))
			// Define a map to hold the JSON data
			var jsonData map[string]interface{}
			// Unmarshal the JSON data into the map
			err := json.Unmarshal(httpBody, &jsonData)
			if err != nil {
				a.logger.Error(err, "Error unmarshaling JSON Reuqest Body")
			}
			a.logger.Sugar().Debug(fmt.Sprintf("jsonData %+v\n", jsonData))
			// Change the model to the selected model
			jsonData["model"] = selectedModel
			// Convert the JSON object to a []byte
			newHTTPBody, err := json.Marshal(jsonData)
			if err != nil {
				a.logger.Error(err, "Error marshaling JSON")
			}

			// Calculate the new body length
			newBodyLength := len(newHTTPBody)
			a.logger.Sugar().Debug(fmt.Sprintf("new body length: %d\n", newBodyLength))

			result.AddHeaders = map[string]string{
				"Content-Length": fmt.Sprintf("%d", newBodyLength), // Set the new Content-Length
				"x-wso2-cluster-header": selectedEndpoint, // Set the cluster header
			}
			result.ModifyBody = true
			result.Body = string(newHTTPBody)
			a.logger.Sugar().Debug(fmt.Sprintf("Modified request body by round robin logic: %+v\n", newHTTPBody))
		}
	} else if requestConfigHolder.ProcessingPhase == requestconfig.ProcessingPhaseResponseHeaders {
		headerValues := requestConfigHolder.ResponseHeaders.GetHeaders().GetHeaders()
		a.logger.Sugar().Debug(fmt.Sprintf("Header Values: %v", headerValues))
		remainingTokenCount := 100
		remainingRequestCount := 100
		remainingCount := 100
		status := 200
		err := error(nil)
		for _, headerValue := range headerValues {
			if headerValue.Key == "x-ratelimit-remaining-tokens" {
				value, err := util.ConvertStringToInt(string(headerValue.RawValue))
				if err != nil {
					a.logger.Error(err, "Unable to retrieve remaining token count by header")
				}
				remainingTokenCount = value
			}
			if headerValue.Key == "x-ratelimit-remaining-requests" {
				value, err := util.ConvertStringToInt(string(headerValue.RawValue))
				if err != nil {
					a.logger.Error(err, "Unable to retrieve remaining request count by header")
				}
				remainingRequestCount = value
			}
			if headerValue.Key == "status" {
				status, err = util.ConvertStringToInt(string(headerValue.RawValue))
				if err != nil {
					a.logger.Error(err, "Unable to retrieve status code by header")
				}
			}
			if headerValue.Key == "x-ratelimit-remaining" {
				value, err := util.ConvertStringToInt(string(headerValue.RawValue))
				if err != nil {
					a.logger.Error(err, "Unable to retrieve remaining count by header")
				}
				remainingCount = value
			}
		}
		if remainingCount <= 0 || remainingTokenCount <= 0 || remainingRequestCount <= 0 || status == 429 { // Suspend model if token/request count reaches 0 or status code is 429
			a.logger.Sugar().Debug("Token/request are exhausted. Suspending the model")
			requestConfigHolder.AI.SuspendModel = true
		}
	} else if requestConfigHolder.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		a.logger.Sugar().Debug("API Level Model Based Round Robin enabled")
		httpBody := requestConfigHolder.ResponseBody.Body
		// Define a map to hold the JSON data
		var jsonData map[string]interface{}
		// Unmarshal the JSON data into the map
		err := json.Unmarshal(httpBody, &jsonData)
		if err != nil {
			a.logger.Error(err, "Error unmarshaling JSON Response Body")
		}
		a.logger.Sugar().Debug(fmt.Sprintf("jsonData %+v\n", jsonData))
		// Retrieve Model from the JSON data
		model := ""
		if modelValue, ok := jsonData["model"].(string); ok {
			model = modelValue
		} else {
			a.logger.Error(fmt.Errorf("model is not a string"), "failed to extract model from JSON data")
		}
		a.logger.Sugar().Debug("Suspending model: " + model)
		duration := a.OnQuotaExceedSuspendDuration
		datastore.GetModelBasedRoundRobinTracker().SuspendModel(requestConfigHolder.RouteMetadata.Spec.API.Name, requestConfigHolder.RouteName, model, time.Duration(time.Duration(duration*1000*1000*1000)))
	}
	return result
}
