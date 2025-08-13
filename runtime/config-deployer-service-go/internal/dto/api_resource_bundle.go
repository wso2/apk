package dto

import "github.com/wso2/apk/config-deployer-service-go/internal/model"

// APIResourceBundle is a DTO that represents a bundle of API resources.
type APIResourceBundle struct {
	Organization      string             `json:"organization,omitempty" yaml:"organization,omitempty"`
	Namespace         string             `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	CPInitiated       bool               `json:"cpInitiated,omitempty" yaml:"cpInitiated,omitempty"`
	APKConf           *model.APKConf     `json:"apkConf,omitempty" yaml:"apkConf,omitempty"`
	CombinedResources []CombinedResource `json:"combinedResources,omitempty" yaml:"combinedResources,omitempty"`
	Definition        string             `json:"definition,omitempty" yaml:"definition,omitempty"`
}

// CombinedResource is a DTO that represents a combination of resources based on APK operations.
type CombinedResource struct {
	APKOperations []model.APKOperations `json:"apkOperations,omitempty" yaml:"apkOperations,omitempty"`
}
