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

package rhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type RetrievableHTTPClient struct {
	Client        *http.Client
	RetryInterval int
	RetryCount    int
}

func (r *RetrievableHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	respChannel := make(chan *http.Response, 1)
	errChannel := make(chan error, 1)

	go func() {
		var resp *http.Response
		var err error
		for i := 1; i <= r.RetryCount; i++ {
			if ctx.Err() != nil {
				errChannel <- ctx.Err()
				return
			}

			logrus.Debug("Sending request, attempt: ", i)
			resp, err = r.Client.Do(req.WithContext(ctx))
			if err == nil {
				logrus.Debugf("Request succeeded in attempt %d", i)
				respChannel <- resp
				return
			}

			logrus.Errorf(
				"Request attempt %d failed due to %s, retrying in %ds. Remaining retries: %d",
				i, err.Error(), r.RetryInterval, r.RetryCount-i)

			select {
			case <-ctx.Done():
				errChannel <- ctx.Err()
				return
			case <-time.After(time.Duration(r.RetryInterval) * time.Second):
			}
		}

		if err != nil {
			logrus.Errorf("Request failed after %d attempts. Backedoff", r.RetryCount)
			errChannel <- fmt.Errorf("request failed after %d attempts: %w", r.RetryCount, err)
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respChannel:
		return resp, nil
	case err := <-errChannel:
		return nil, err
	}
}
