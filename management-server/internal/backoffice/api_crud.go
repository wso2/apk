package backoffice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
)

// Back Office client connetion
var backOfficeClient *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: nil,
	}
	backOfficeClient = &http.Client{Transport: transport}
}

func getBackOfficeURL() string {
	conf := config.ReadConfigs()
	logger.LoggerMGTServer.Infof("backoffice service: http://%s:%d%s", conf.BackOffice.Host, conf.BackOffice.Port, conf.BackOffice.ServiceBasePath)
	return fmt.Sprintf("http://%s:%d%s", conf.BackOffice.Host, conf.BackOffice.Port, conf.BackOffice.ServiceBasePath)
}

// CreateAPI creates an API by invoking backoffice service
func CreateAPI(api *apiProtos.API) error {
	postBody, _ := json.Marshal(api)
	requestBody := bytes.NewBuffer(postBody)
	_, err := backOfficeClient.Post(getBackOfficeURL(), "application/json", requestBody)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAPI updates an API by invoking backoffice service
func UpdateAPI(api *apiProtos.API) error {
	putBody, _ := json.Marshal(api)
	requestBody := bytes.NewBuffer(putBody)
	putRequest, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", getBackOfficeURL(), api.Uuid), requestBody)
	_, err = backOfficeClient.Do(putRequest)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAPI deletes an API by invoking backoffice service
func DeleteAPI(api *apiProtos.API) error {
	deleteRequest, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", getBackOfficeURL(), api.Uuid), nil)
	_, err = backOfficeClient.Do(deleteRequest)
	if err != nil {
		return err
	}
	return nil
}
