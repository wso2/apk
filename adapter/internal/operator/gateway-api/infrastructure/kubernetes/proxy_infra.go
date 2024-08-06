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
	"errors"

	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes/proxy"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
)

// CreateOrUpdateProxyInfra creates the managed kube infra, if it doesn't exist.
func (i *Infra) CreateOrUpdateProxyInfra(ctx context.Context, infra *ir.Infra) error {
	if infra == nil {
		return errors.New("infra ir is nil")
	}

	if infra.Proxy == nil {
		return errors.New("infra proxy ir is nil")
	}

	r := proxy.NewResourceRender(i.Namespace, infra.GetProxyInfra())
	return i.createOrUpdate(ctx, r)
}

// DeleteProxyInfra removes the managed kube infra, if it doesn't exist.
func (i *Infra) DeleteProxyInfra(ctx context.Context, infra *ir.Infra) error {
	if infra == nil {
		return errors.New("infra ir is nil")
	}

	r := proxy.NewResourceRender(i.Namespace, infra.GetProxyInfra())
	return i.delete(ctx, r)
}
