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

// GenerateCRName generates a Kubernetes CR name based on the API name, environment, version, and organization.
// It sanitizes the name to ensure it is Kubernetes valid.
// The name is constructed as "apiName-env-version-organization" and sanitized.
// This is used for generating names for various Kubernetes resources.
func GenerateCRName(apiName, env, version, organization string) string {
	// concatenate the parts to be hashed
	toHash := "-" + env + "-" + version + "-" + organization

	// compute SHA-1 hash
	h := sha1.New()
	h.Write([]byte(toHash))
	hashBytes := h.Sum(nil)

	// convert to hex string
	hashHex := hex.EncodeToString(hashBytes)

	// take last 10 characters of hash
	last10 := hashHex[len(hashHex)-10:]

	// replace any spaces or dots in apiName
	cleanAPIName := strings.ReplaceAll(apiName, ".", "-")
	cleanAPIName = strings.ReplaceAll(cleanAPIName, " ", "-")
	cleanAPIName = strings.ToLower(cleanAPIName)

	return cleanAPIName + last10
}

func IsSameRatelimit(r1 model.RateLimit, r2 model.RateLimit) bool {
	return r1.RequestsPerUnit == r2.RequestsPerUnit && r1.Unit == r2.Unit 
}

// HashLast50SHA1 returns the last 50 characters of the SHA-1 hash of the input string.
// Since SHA-1 is only 40 hex chars, this will just return the whole hash.
func HashLast50SHA1(input string) string {
	hash := sha1.Sum([]byte(input))
	hexStr := hex.EncodeToString(hash[:]) // SHA-1 produces 40 hex chars
	if len(hexStr) <= 50 {
		return hexStr
	}
	return hexStr[len(hexStr)-50:]
}