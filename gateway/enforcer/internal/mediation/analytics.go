package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

// Analytics represents the configuration for Analytics policy in the API Gateway.
type Analytics struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Enabled       bool   `json:"enabled"`
	logger        *logging.Logger
	cfg           *config.Server
}

const (
	// MediationAnalyticsPolicyKeyEnabled is the key for enabling/disabling the Analytics policy.
	MediationAnalyticsPolicyKeyEnabled = "Enabled"
)

// NewAnalytics creates a new Analytics instance with default values.
func NewAnalytics(mediation *dpv2alpha1.Mediation) *Analytics {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, MediationAnalyticsPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	cfg := config.GetConfig()
	logger := cfg.Logger
	return &Analytics{
		PolicyName:    "Analytics",
		PolicyVersion: mediation.PolicyVersion,
		PolicyID:      mediation.PolicyID,
		Enabled:       enabled,
		logger:        &logger,
		cfg:           cfg,
	}
}

// Process processes the request configuration for analytics.
func (a *Analytics) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Early return if requestConfig itself is nil
	if requestConfig == nil {
		return result
	}

	// Safe helpers
	addMetadata := func(key string, val string) {
		if val == "" {
			return
		}
		v, err := structpb.NewValue(val)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for %s: %v", key, err)
			return
		}
		result.Metadata[key] = v
	}

	// Safely extract API metadata
	if requestConfig.RouteMetadata != nil {
		api := requestConfig.RouteMetadata.Spec.API
		addMetadata(analytics.APIIDKey, api.Name)
		addMetadata(analytics.APIContextKey, api.Context)
		addMetadata(analytics.APINameKey, api.Name)
		addMetadata(analytics.APIVersionKey, api.Version)
		addMetadata(analytics.APITypeKey, api.Type)
		addMetadata(analytics.APICreatorKey, api.APICreator)
		addMetadata(analytics.APICreatorTenantDomainKey, api.APICreatorTenantDomain)
		addMetadata(analytics.APIOrganizationIDKey, api.Organization)
		addMetadata(analytics.APIEnvironmentKey, api.Environment)
	}

	// Request attributes
	if requestConfig.RequestAttributes != nil {
		addMetadata(analytics.CorrelationIDKey, requestConfig.RequestAttributes.RequestID)
	}

	// Region
	addMetadata(analytics.RegionKey, a.cfg.EnforcerRegionID)

	// Application info
	if requestConfig.MatchedApplication != nil {
		app := requestConfig.MatchedApplication
		addMetadata(analytics.AppIDKey, app.UUID)
		addMetadata(analytics.AppUUIDKey, app.UUID)
		// AppKeyType currently reused from API environment
		if requestConfig.RouteMetadata != nil {
			addMetadata(analytics.AppKeyTypeKey, requestConfig.RouteMetadata.Spec.API.Environment)
		}
		addMetadata(analytics.AppNameKey, app.Name)
		addMetadata(analytics.AppOwnerKey, app.Owner)
	}

	return result
}
