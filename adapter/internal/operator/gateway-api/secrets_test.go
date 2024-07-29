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

package gatewayapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	envoyGatewaySecret = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "envoy-gateway",
			Namespace: "apk",
		},
	}

	envoySecret = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "envoy",
			Namespace: "apk",
		},
	}

	envoyRateLimitSecret = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "envoy-rate-limit",
			Namespace: "apk",
		},
	}

	oidcHMACSecret = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "envoy-oidc-hmac",
			Namespace: "apk",
		},
	}

	existingSecretsWithoutHMAC = []client.Object{
		&envoyGatewaySecret,
		&envoySecret,
		&envoyRateLimitSecret,
	}

	existingSecretsWithHMAC = []client.Object{
		&envoyGatewaySecret,
		&envoySecret,
		&envoyRateLimitSecret,
		&oidcHMACSecret,
	}

	SecretsToCreate = []corev1.Secret{
		envoyGatewaySecret,
		envoySecret,
		envoyRateLimitSecret,
		oidcHMACSecret,
	}
)

func TestCreateSecretsWhenUpgrade(t *testing.T) {
	t.Run("create HMAC secret when it does not exist", func(t *testing.T) {
		cli := fakeclient.NewClientBuilder().WithObjects(existingSecretsWithoutHMAC...).Build()

		created, err := CreateOrUpdateSecrets(context.Background(), cli, SecretsToCreate, false)
		require.ErrorIs(t, err, ErrSecretExists)
		require.Len(t, created, 1)
		require.Equal(t, "envoy-oidc-hmac", created[0].Name)

		err = cli.Get(context.Background(), client.ObjectKeyFromObject(&oidcHMACSecret), &corev1.Secret{})
		require.NoError(t, err)
	})

	t.Run("skip HMAC secret when it exist", func(t *testing.T) {
		cli := fakeclient.NewClientBuilder().WithObjects(existingSecretsWithHMAC...).Build()

		created, err := CreateOrUpdateSecrets(context.Background(), cli, SecretsToCreate, false)
		require.ErrorIs(t, err, ErrSecretExists)
		require.Emptyf(t, created, "expected no secrets to be created, got %v", created)
	})

	t.Run("update secrets when they exist", func(t *testing.T) {
		cli := fakeclient.NewClientBuilder().WithObjects(existingSecretsWithHMAC...).Build()

		created, err := CreateOrUpdateSecrets(context.Background(), cli, SecretsToCreate, true)
		require.NoError(t, err)
		require.Len(t, created, 4)
	})
}
