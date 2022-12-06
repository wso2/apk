package backoffice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/management-server/internal/config"
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
	return fmt.Sprintf("http://%s:%d/api/am/backoffice/internal/apis", conf.BackOffice.Host, conf.BackOffice.Port)
}

func CreateAPI(api *apiProtos.API) error {
	postBody, _ := json.Marshal(map[string]string{
		"name":            api.Context,
		"context":         api.Context,
		"version":         api.Version,
		"provider":        "apkuser",
		"lifeCycleStatus": "PUBLISHED",
	})
	responseBody := bytes.NewBuffer(postBody)
	_, err := backOfficeClient.Post(getBackOfficeURL(), "application/json", responseBody)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAPI(api *apiProtos.API) error {
	return nil
}

func DeleteAPI(api *apiProtos.API) error {
	return nil
}
