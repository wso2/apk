package util

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/wso2/apk/config-deployer-service-go/internal/model"
)

// SanitizeOrHashName ensures the name is Kubernetes valid.
// If invalid, returns a hashed version (sha1, hex-encoded, first 10 chars for brevity).
func SanitizeOrHashName(name string) string {
	name = strings.ToLower(name) // ensure lowercase
	if k8sNameRegex.MatchString(name) && len(name) <= 253 {
		return name
	}
	h := sha1.Sum([]byte(name))
	return hex.EncodeToString(h[:])[:10] // short hash
}

// Kubernetes CR name regex: ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
var k8sNameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// GenerateRouteMetadataName generates a Kubernetes RouteMetadata name based on the API name, environment, version, and organization.
// It sanitizes the name to ensure it is Kubernetes valid.
// The name is constructed as "apiName-env-version-organization" and sanitized.
func GenerateRouteMetadataName(apiName, env, version, organization string) string {
	name := apiName + "-" + env + "-" + version + "-" + organization
	name = strings.ReplaceAll(name, ".", "-")
	return SanitizeOrHashName(name)
}

// GenerateCRName generates a Kubernetes CR name based on the API name, environment, version, and organization.
// It sanitizes the name to ensure it is Kubernetes valid.
// The name is constructed as "apiName-env-version-organization" and sanitized.
// This is used for generating names for various Kubernetes resources.
func GenerateCRName(apiName, env, version, organization string) string {
	name := apiName + "-" + env + "-" + version + "-" + organization
	name = strings.ReplaceAll(name, ".", "-")
	return SanitizeOrHashName(name)
}

func IsSameRatelimit(r1 model.RateLimit, r2 model.RateLimit) bool {
	return r1.RequestsPerUnit == r2.RequestsPerUnit && r1.Unit == r2.Unit 
}
