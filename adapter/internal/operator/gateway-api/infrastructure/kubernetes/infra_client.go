/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package kubernetes

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/loggers"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InfraClient struct {
	client.Client
}

func New(cli client.Client) *InfraClient {
	return &InfraClient{
		Client: cli,
	}
}

func (cli *InfraClient) CreateOrUpdate(ctx context.Context, key client.ObjectKey, current client.Object, specific client.Object, updateChecker func() bool) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := cli.Client.Get(ctx, key, current); err != nil {
			loggers.LoggerAPI.Debugf("Error while getting resource %+v : %v", key, err)
			if kerrors.IsNotFound(err) {
				loggers.LoggerAPI.Infof("Creating a new resource %+v", key)
				// Create if it does not exist.
				if err := cli.Client.Create(ctx, specific); err != nil {
					return fmt.Errorf("for Create: %w", err)
				}
			}
		} else {
			// Since the client.Object does not have a specific Spec field to compare
			// just perform an update for now.
			if updateChecker() {
				specific.SetUID(current.GetUID())
				if err := cli.Client.Update(ctx, specific); err != nil {
					return fmt.Errorf("for Update: %w", err)
				}
			}
		}

		return nil
	})
}

func (cli *InfraClient) Delete(ctx context.Context, object client.Object) error {
	if err := cli.Client.Delete(ctx, object); err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

// GetUID retrieves the uid of one resource.
func (cli *InfraClient) GetUID(ctx context.Context, key client.ObjectKey, current client.Object) (types.UID, error) {
	if err := cli.Client.Get(ctx, key, current); err != nil {
		return "", err
	}
	return current.GetUID(), nil
}
