package util

import (
	"fmt"
)

// PrepareAPIKey prepares the API key using the given vhost, basePath, and version.
func PrepareAPIKey(vhost, basePath, version string) string {
    return fmt.Sprintf("%s:%s:%s", vhost, basePath, version)
}