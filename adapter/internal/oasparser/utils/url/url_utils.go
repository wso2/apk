package urlutils

import "github.com/wso2/apk/adapter/internal/oasparser/constants"

// GetURLType return URLType as https or http
// depending on whether tls is enabled or not.
func GetURLType(tlsEnabled bool) string {
	if tlsEnabled {
		return constants.HTTPS
	}
	return constants.HTTP
}
