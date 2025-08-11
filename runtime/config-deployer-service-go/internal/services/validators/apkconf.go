package validators

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"github.com/xeipuuv/gojsonschema"
)

var GlobalAPKConfValidator *APKConfValidator

type APKConfValidator struct {
	SchemaContent string
}

func NewAPKConfValidator(schemaContent string) *APKConfValidator {
	return &APKConfValidator{
		SchemaContent: schemaContent,
	}
}

// ValidateAPKConf validates the APK configuration JSON string and returns a validation response.
func (apkConfValidator *APKConfValidator) ValidateAPKConf(apkConfJson string) (*dto.APKConfValidationResponse, error) {
	schemaLoader := gojsonschema.NewStringLoader(GlobalAPKConfValidator.SchemaContent)
	documentLoader := gojsonschema.NewStringLoader(apkConfJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}
	response := &dto.APKConfValidationResponse{
		Validated: result.Valid(),
	}
	if !result.Valid() {
		var errors []dto.ErrorHandler
		for _, desc := range result.Errors() {
			errorHandler := dto.ErrorHandler{
				ErrorMessage:     desc.Type(),
				ErrorDescription: desc.Description(),
			}
			errors = append(errors, errorHandler)
		}
		response.ErrorItems = errors
	}
	return response, nil
}

// ValidateEndpointConfigurations validates endpoint configurations for APK operations
func (apkConfValidator *APKConfValidator) ValidateEndpointConfigurations(apkConf *model.APKConf) map[string]string {
	var errors map[string]string
	productionEndpointAvailable := false
	sandboxEndpointAvailable := false

	if apkConf.EndpointConfigurations != nil {
		sandboxEndpointAvailable = apkConf.EndpointConfigurations.Sandbox != nil &&
			len(apkConf.EndpointConfigurations.Sandbox) > 0
		productionEndpointAvailable = apkConf.EndpointConfigurations.Production != nil &&
			len(apkConf.EndpointConfigurations.Production) > 0
	}

	if apkConf.Operations != nil {
		for _, operation := range apkConf.Operations {
			operationLevelProductionEndpointAvailable := false
			operationLevelSandboxEndpointAvailable := false

			if operation.EndpointConfigurations != nil {
				operationLevelProductionEndpointAvailable = operation.EndpointConfigurations.Production != nil &&
					len(operation.EndpointConfigurations.Production) > 0
				operationLevelSandboxEndpointAvailable = operation.EndpointConfigurations.Sandbox != nil &&
					len(operation.EndpointConfigurations.Sandbox) > 0
			}

			if (!operationLevelProductionEndpointAvailable && !productionEndpointAvailable) &&
				(!operationLevelSandboxEndpointAvailable && !sandboxEndpointAvailable) {
				target := "unknown"
				if operation.Target != nil {
					target = *operation.Target
				}
				errors["endpoint"] = fmt.Sprintf("production/sandbox endpoint not available for %s", target)
			}
		}
	}
	return errors
}

// ValidateRateLimit validates the rate limit configuration for APK operations
func (apkConfValidator *APKConfValidator) ValidateRateLimit(apiRateLimit *model.RateLimit, operations []model.APKOperations) error {
	if apiRateLimit == nil {
		return nil
	} else {
		for _, operation := range operations {
			operationRateLimit := operation.RateLimit
			if operationRateLimit != nil {
				return fmt.Errorf("presence of both resource level and API level rate limits is not allowed")
			}
		}
	}
	return nil
}
