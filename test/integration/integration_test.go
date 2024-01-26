/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 */

package integration

import (
	"testing"

	"github.com/wso2/apk/test/integration/integration/tests"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
	
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"

)

func TestIntegration(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("Error loading Kubernetes config: %v", err)
	}
	client, err := client.New(cfg, client.Options{})
	if err != nil {
		t.Fatalf("Error initializing Kubernetes client: %v", err)
	}

	v1alpha2.Install(client.Scheme())
	v1beta1.Install(client.Scheme())
	dpv1alpha1.AddToScheme(client.Scheme())


	// TODO(Amila): Uncomment after operator package in adaptor is moved from internal to pkg directory
	// dpv1alpha1.Install(client.Scheme())

	cSuite := suite.New(suite.Options{
		Client:               client,
		Debug:                true,
		CleanupBaseResources: true,
	})
	cSuite.Setup(t)
	cSuite.Run(t, tests.IntegrationTests)
}
