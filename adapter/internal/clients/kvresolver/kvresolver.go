/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package kvresolver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/wso2/apk/adapter/internal/loggers"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/pkg/rhttp"
)

var (
	kvResolverOnce sync.Once
	basePath       = "/api/v1/secrets"
)

type KVResolverClient interface {
	GetSecrets(ctx context.Context, secrets []string) (Secrets, error)
}

type KVResolverClientImpl struct {
	KVResolverClient *rhttp.RetrievableHTTPClient
	Config           *config.Config
}

var kvResolverClient = &KVResolverClientImpl{}

func InitializeKVResolverClient() *KVResolverClientImpl {
	conf := config.ReadConfigs()

	kvResolverOnce.Do(func() {
		netTransport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        20,
		}

		client := &http.Client{
			Transport: netTransport,
		}
		kvResolverClient.KVResolverClient = &rhttp.RetrievableHTTPClient{
			Client:        client,
			RetryInterval: 2,
			RetryCount:    5,
		}
	})
	kvResolverClient.Config = conf
	return kvResolverClient
}

func (k *KVResolverClientImpl) GetSecrets(ctx context.Context, secrets []string) (Secrets, error) {
	loggers.LoggerAPK.Info(fmt.Sprintf("Getting secrets for keys: %v", secrets))
	url := fmt.Sprintf("%s%s/get", k.Config.KVResolver.ServiceURL, basePath)
	payload, err := getSecretPayload(secrets)
	if err != nil {
		loggers.LoggerAPK.Error("Error while getting secret payload", err.Error())
		return Secrets{}, err
	}
	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if reqErr != nil {
		loggers.LoggerAPK.Error("Error while creating request", reqErr.Error())
		return Secrets{}, reqErr
	}
	req.Header.Set("Content-Type", "application/json")
	res, kvResolverError := k.KVResolverClient.Do(ctx, req)

	if kvResolverError != nil {
		loggers.LoggerAPK.Error("Error while invoking rudder", kvResolverError.Error())
		return Secrets{}, kvResolverError
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		payload, readErr := io.ReadAll(res.Body)

		if readErr != nil {
			loggers.LoggerAPK.Error("Error while reading response body", readErr.Error())
			return Secrets{}, readErr
		}
		loggers.LoggerAPK.Error(
			fmt.Sprintf("Error response from rudder with status code: %d response payload %s",
				res.StatusCode, string(payload)), string(payload))
		return Secrets{}, fmt.Errorf("error from upstream with status code: %d", res.StatusCode)
	}

	var secretResp Secrets
	if err := json.NewDecoder(res.Body).Decode(&secretResp); err != nil {
		loggers.LoggerAPK.Error("Error while decoding response body", err.Error())
		return Secrets{}, err
	}
	loggers.LoggerAPK.Debugf("KV Resolver returned %d secrets", len(secretResp.Secrets))
	for i, secret := range secretResp.Secrets {
		loggers.LoggerAPK.Debugf("Secret %d: ID=%s, Key=%s, ValueType=%s", i, secret.ID, secret.Key, secret.ValueType)
	}
	return secretResp, nil
}

func getSecretPayload(secrets []string) ([]byte, error) {
	var secretPayload Secrets
	for _, secret := range secrets {
		secretPayload.Secrets = append(secretPayload.Secrets, Secret{Key: secret})
	}
	payload, err := json.Marshal(secretPayload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
