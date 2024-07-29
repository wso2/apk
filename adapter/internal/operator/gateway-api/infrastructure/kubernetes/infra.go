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

	"github.com/wso2/apk/adapter/config"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceRender renders Kubernetes infrastructure resources
// based on Infra IR resources.
type ResourceRender interface {
	Name() string
	ServiceAccount() (*corev1.ServiceAccount, error)
	Service() (*corev1.Service, error)
	ConfigMap() (*corev1.ConfigMap, error)
	Deployment() (*appsv1.Deployment, error)
	HorizontalPodAutoscaler() (*autoscalingv2.HorizontalPodAutoscaler, error)
}

// Infra manages the creation and deletion of Kubernetes infrastructure
// based on Infra IR resources.
type Infra struct {
	// Namespace is the Namespace used for managed infra.
	Namespace string

	// Client wrap k8s client.
	Client *InfraClient
}

// NewInfra returns a new Infra.
func NewInfra(cli client.Client) *Infra {
	conf := config.ReadConfigs()
	return &Infra{
		Namespace: conf.Envoy.Namespace,
		Client:    New(cli),
	}
}

// createOrUpdate creates a ServiceAccount/ConfigMap/Deployment/Service in the kube api server based on the
// provided ResourceRender, if it doesn't exist and updates it if it does.
func (i *Infra) createOrUpdate(ctx context.Context, r ResourceRender) error {

	// certs, err := crypto.GenerateCerts()
	// if err != nil {
	// 	return fmt.Errorf("failed to generate certificates: %w", err)
	// }
	// secrets, err := gatewayapi.CreateOrUpdateSecrets(ctx, i.Client.Client, gatewayapi.CertsToSecret(i.Namespace, certs), true)

	// if err != nil {
	// 	if errors.Is(err, gatewayapi.ErrSecretExists) {
	// 		loggers.LoggerAPKOperator.Info(err.Error())
	// 	} else {
	// 		return fmt.Errorf("failed to create or update secrets: %w", err)
	// 	}
	// } else {
	// 	for i := range secrets {
	// 		s := secrets[i]
	// 		loggers.LoggerAPKOperator.Info("created secret", "namespace", s.Namespace, "name", s.Name)
	// 	}
	// }

	if err := i.createOrUpdateServiceAccount(ctx, r); err != nil {
		return fmt.Errorf("failed to create or update serviceaccount %s/%s: %w", i.Namespace, r.Name(), err)
	}

	// if err := i.createOrUpdateConfigMap(ctx, r); err != nil {
	// 	return fmt.Errorf("failed to create or update configmap %s/%s: %w", i.Namespace, r.Name(), err)
	// }

	if err := i.createOrUpdateDeployment(ctx, r); err != nil {
		return fmt.Errorf("failed to create or update deployment %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.createOrUpdateService(ctx, r); err != nil {
		return fmt.Errorf("failed to create or update service %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.createOrUpdateHPA(ctx, r); err != nil {
		return fmt.Errorf("failed to create or update hpa %s/%s: %w", i.Namespace, r.Name(), err)
	}

	return nil
}

// delete deletes the ServiceAccount/ConfigMap/Deployment/Service in the kube api server, if it exists.
func (i *Infra) delete(ctx context.Context, r ResourceRender) error {
	if err := i.deleteServiceAccount(ctx, r); err != nil {
		return fmt.Errorf("failed to delete serviceaccount %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.deleteConfigMap(ctx, r); err != nil {
		return fmt.Errorf("failed to delete configmap %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.deleteDeployment(ctx, r); err != nil {
		return fmt.Errorf("failed to delete deployment %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.deleteService(ctx, r); err != nil {
		return fmt.Errorf("failed to delete service %s/%s: %w", i.Namespace, r.Name(), err)
	}

	if err := i.deleteHPA(ctx, r); err != nil {
		return fmt.Errorf("failed to delete hpa %s/%s: %w", i.Namespace, r.Name(), err)
	}

	return nil
}
