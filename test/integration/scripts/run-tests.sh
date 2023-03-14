#!/usr/bin/env bash

# Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.

# Create new namespace to install chart
kubectl create ns apk-integration-test

# Install wso2 apk chart with cp disabled
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add jetstack https://charts.jetstack.io
helm dependency build ../../helm-charts
helm install apk-test-setup ../../helm-charts -n apk-integration-test --set wso2.apk.cp.enabled=false \
--set dp.adapter.deployment.image=adapter:test --set dp.adapter.deployment.imagePullPolicy=IfNotPresent \
--set dp.gatewayRuntime.deployment.router.image=router:test --set dp.gatewayRuntime.deployment.router.imagePullPolicy=IfNotPresent \
--set dp.runtime.deployment.enforcer.image=enforcer:test --set dp.runtime.deployment.enforcer.imagePullPolicy=IfNotPresent \
--set wso2.apk.dp.ratelimiter.enabled=false --set wso2.apk.dp.redis.enabled=false


# Wait gateway resources to be available.
kubectl wait --timeout=5m -n gateway-system deployment/gateway-api-admission-server --for=condition=Available
kubectl wait --timeout=5m -n gateway-system job/gateway-api-admission --for=condition=Complete
kubectl wait --timeout=5m -n apk-integration-test deployment/apk-test-setup-wso2-apk-adapter-deployment --for=condition=Available

# Run tests
go test -v integration_test.go
