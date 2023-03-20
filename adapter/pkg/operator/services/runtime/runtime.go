package runtime

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
)

// Runtime client connetion
var runtimeClient *http.Client

func init() {
	_, _, truststoreLocation := tlsutils.GetKeyLocations()
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}
	runtimeClient = &http.Client{Transport: transport}
}

func getInternalRuntimeServiceURL() string {
	conf := config.ReadConfigs()
	serviceURL := fmt.Sprintf("https://%s:%d%s",
		conf.Runtime.Host,
		conf.Runtime.Port,
		conf.Runtime.ServiceBasePath)
	loggers.LoggerAPKOperator.Debugf("runtime service: %s", serviceURL)
	return serviceURL
}

// GetAPIDefinition gets the API defintion using API UUID and return it as a string
func GetAPIDefinition(apiUUID string, organizationUUID string) string {
	definitionString := "{}"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/apis/%s/definition",
		getInternalRuntimeServiceURL(), apiUUID), nil)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating api definition request: %v", err)
	}
	req.Header.Set(constants.OrganizationHeader, organizationUUID)

	response, err := runtimeClient.Do(req)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error retrieving api definition: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error reading api definition: %v", err)
		}
		definitionString = string(bodyBytes)
	}
	return definitionString
}
