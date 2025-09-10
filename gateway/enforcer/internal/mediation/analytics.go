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
	// Implement the logic to process the requestConfig for analytics
	// This is a placeholder implementation
	result := NewResult()
	var err error
	result.Metadata[analytics.APIIDKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Name)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APIIDKey: %v", err)
	}
	result.Metadata[analytics.APIContextKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Context)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APIContextKey: %v", err)
	}
	// result.Metadata[organizationMetadataKey], err = structpb.NewValue(requestConfigHolder.MatchedAPI.OrganizationID)
	result.Metadata[analytics.APINameKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Name)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APINameKey: %v", err)
	}
	result.Metadata[analytics.APIVersionKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Version)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APIVersionKey: %v", err)
	}
	result.Metadata[analytics.APITypeKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Type)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APITypeKey: %v", err)
	}
	result.Metadata[analytics.APICreatorKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.APICreator)
	result.Metadata[analytics.APICreatorTenantDomainKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.APICreatorTenantDomain)
	result.Metadata[analytics.APIOrganizationIDKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Organization)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APIOrganizationIDKey: %v", err)
	}

	result.Metadata[analytics.CorrelationIDKey], err = structpb.NewValue(requestConfig.RequestAttributes.RequestID)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for CorrelationIDKey: %v", err)
	}
	result.Metadata[analytics.RegionKey], err = structpb.NewValue(a.cfg.EnforcerRegionID)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for RegionKey: %v", err)
	}
	// result.Metadata[analytics.UserAgentKey], err = structpb.NewValue(s.requestConfigHolder.Metadata.UserAgent)
	// result.Metadata[analytics.ClientIpKey], err = structpb.NewValue(s.requestConfigHolder.Metadata.ClientIP)
	// result.Metadata[analytics.ApiResourceTemplateKey], err = structpb.NewValue(s.requestConfigHolder.ApiResourceTemplate)
	// result.Metadata[analytics.Destination], err = structpb.NewValue(s.requestConfigHolder.Metadata.Destination)
	result.Metadata[analytics.APIEnvironmentKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Environment)
	if err != nil {
		a.logger.Sugar().Errorf("Error creating structpb value for APIEnvironmentKey: %v", err)
	}
	if requestConfig.MatchedApplication != nil {
		result.Metadata[analytics.AppIDKey], err = structpb.NewValue(requestConfig.MatchedApplication.UUID)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for AppIDKey: %v", err)
		}
		result.Metadata[analytics.AppUUIDKey], err = structpb.NewValue(requestConfig.MatchedApplication.UUID)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for AppUUIDKey: %v", err)
		}
		result.Metadata[analytics.AppKeyTypeKey], err = structpb.NewValue(requestConfig.RouteMetadata.Spec.API.Environment)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for AppKeyTypeKey: %v", err)
		}
		result.Metadata[analytics.AppNameKey], err = structpb.NewValue(requestConfig.MatchedApplication.Name)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for AppNameKey: %v", err)
		}
		result.Metadata[analytics.AppOwnerKey], err = structpb.NewValue(requestConfig.MatchedApplication.Owner)
		if err != nil {
			a.logger.Sugar().Errorf("Error creating structpb value for AppOwnerKey: %v", err)
		}
	}
	return result
}
