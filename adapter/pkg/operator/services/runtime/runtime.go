package runtime

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
)

// Runtime client connetion
var runtimeClient *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: nil,
	}
	runtimeClient = &http.Client{Transport: transport}
}

func getRuntimeServiceURL() string {
	conf := config.ReadConfigs()
	serviceURL := fmt.Sprintf("http://%s:%d%s",
		conf.Runtime.Host,
		conf.Runtime.Port,
		conf.Runtime.ServiceBasePath)
	loggers.LoggerAPKOperator.Debugf("runtime service: %s", serviceURL)
	return serviceURL
}

// GetAPIDefinition gets the API defintion using API UUID and return it as a string
func GetAPIDefinition(apiUUID string) string {
	definitionString := "{}"
	response, err := runtimeClient.Get(fmt.Sprintf("%s/apis/%s/definition",
		getRuntimeServiceURL(), apiUUID))
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
