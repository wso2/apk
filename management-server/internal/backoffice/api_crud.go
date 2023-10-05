package backoffice

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
)

// Back Office client connetion
var backOfficeClient *http.Client

type api struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Context         string `json:"context"`
	Version         string `json:"version"`
	Provider        string `json:"provider"`
	OrganizationID  string `json:"organization"`
	LifeCycleStatus string `json:"LifeCycleStatus"`
}

type definition interface{}

type requestData struct {
	APIProperties api        `json:"apiProperties"`
	Definition    definition `json:"Definition"`
}

func init() {
	_, _, truststoreLocation := tlsutils.GetKeyLocations()
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}
	backOfficeClient = &http.Client{Transport: transport}
}

func getBackOfficeURL() string {
	conf := config.ReadConfigs()
	logger.LoggerMGTServer.Infof("backoffice service: https://%s:%d%s", conf.BackOffice.Host, conf.BackOffice.Port, conf.BackOffice.ServiceBasePath)
	return fmt.Sprintf("https://%s:%d%s", conf.BackOffice.Host, conf.BackOffice.Port, conf.BackOffice.ServiceBasePath)
}

func composeRequestBody(api *apiProtos.API) requestData {
	request := new(requestData)
	request.APIProperties.ID = api.Uuid
	request.APIProperties.Name = api.Name
	request.APIProperties.Context = api.BasePath
	request.APIProperties.Version = api.Version
	request.APIProperties.Provider = api.Provider
	request.APIProperties.OrganizationID = api.OrganizationId
	json.Unmarshal([]byte(api.Definition), &request.Definition)
	return *request
}

// CreateAPI creates an API by invoking backoffice service
func CreateAPI(api *apiProtos.API) error {
	postBody, _ := json.Marshal(composeRequestBody(api))
	requestBody := bytes.NewBuffer(postBody)
	_, err := backOfficeClient.Post(getBackOfficeURL(), "application/json", requestBody)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAPI updates an API by invoking backoffice service
func UpdateAPI(api *apiProtos.API) error {
	putBody, _ := json.Marshal(composeRequestBody(api))
	requestBody := bytes.NewBuffer(putBody)
	putRequest, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", getBackOfficeURL(), api.Uuid), requestBody)
	if err != nil {
		return err
	}

	// Perform the HTTP request and check the response status code
	response, err := backOfficeClient.Do(putRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		// If the status code indicates an 404, call the create API to create the API in database.
		// This is done to handle the case where the API is not in the database due to managemnt server failure.
		CreateAPI(api)
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
