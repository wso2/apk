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
 */

package tests

import (
	"context"
	"testing"
	"time"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func init() {
	// //IntegrationTests = append(IntegrationTests, UpdateOwnerReferences)
}

// UpdateOwnerReferences test
var UpdateOwnerReferences = suite.IntegrationTest{
	ShortName:   "UpdateOwnerReference",
	Description: "Test owner reference functionality",
	Manifests: []string{
		"tests/owner-reference.yaml",
	},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		time.Sleep(5 * time.Second)
		namespace := "apk-integration-test"
		// Get hr1 route check the owner reference
		hr1Key := types.NamespacedName{
			Name:      "hr1",
			Namespace: namespace,
		}
		var hr1 gwapiv1b1.HTTPRoute
		if err := suite.Client.Get(context.TODO(), hr1Key, &hr1); err != nil {
			t.Fatalf("Unable to load http route with key %+v error: %+v", hr1Key, err)
		}
		// check the owner reference
		if len(hr1.ObjectMeta.OwnerReferences) != 1 || hr1.ObjectMeta.OwnerReferences[0].Name != "api1" {
			t.Fatalf("Unexpected owner reference found in http route hr1:  %+v", hr1)
		}

		// Get hr2 route check the owner reference
		hr2Key := types.NamespacedName{
			Name:      "hr2",
			Namespace: namespace,
		}
		var hr2 gwapiv1b1.HTTPRoute
		if err := suite.Client.Get(context.TODO(), hr2Key, &hr2); err != nil {
			t.Fatalf("Unable to load http route with key %+v error: %+v", hr2Key, err)
		}
		// check the owner reference
		if len(hr2.ObjectMeta.OwnerReferences) != 2 {
			t.Fatalf("Unexpected owner reference found in http route hr2:  %+v", hr2)
		}
		namesToCheck := []string{"api1", "api2"}
		found := false
		for _, ownerRef := range hr2.ObjectMeta.OwnerReferences {
			foundLocal := false
			for _, name := range namesToCheck {
				if ownerRef.Name == name {
					foundLocal = true
				}
			}
			found = foundLocal
		}
		if !found {
			t.Fatalf("Unexpected owner reference found in http route hr2:  %+v", hr2)
		}
		// Load backend1
		backend1Key := types.NamespacedName{
			Name:      "backend1",
			Namespace: namespace,
		}
		var backend1 dpv1alpha1.Backend
		if err := suite.Client.Get(context.TODO(), backend1Key, &backend1); err != nil {
			t.Fatalf("Unable to load backend with key %+v error: %+v", backend1Key, err)
		}
		// check the owner reference
		if len(backend1.ObjectMeta.OwnerReferences) != 1 || backend1.ObjectMeta.OwnerReferences[0].Name != "api1" {
			t.Fatalf("Unexpected owner reference found in backend1:  %+v", backend1)
		}

		// Ge backend2 check the owner reference
		backend2Key := types.NamespacedName{
			Name:      "backend2",
			Namespace: namespace,
		}
		var backend2 dpv1alpha1.Backend
		if err := suite.Client.Get(context.TODO(), backend2Key, &backend2); err != nil {
			t.Fatalf("Unable to load backend with key %+v error: %+v", backend2Key, err)
		}
		// check the owner reference
		if len(hr2.ObjectMeta.OwnerReferences) != 2 {
			t.Fatalf("Unexpected owner reference found in backend2:  %+v", backend2)
		}
		found = false
		for _, ownerRef := range hr2.ObjectMeta.OwnerReferences {
			foundLocal := false
			for _, name := range namesToCheck {
				if ownerRef.Name == name {
					foundLocal = true
				}
			}
			found = foundLocal
		}
		if !found {
			t.Fatalf("Unexpected owner reference found in backend2:  %+v", backend2)
		}

		// Get rl1
		rl1Key := types.NamespacedName{
			Name:      "rl1",
			Namespace: namespace,
		}
		var rl1 dpv1alpha1.RateLimitPolicy
		if err := suite.Client.Get(context.TODO(), rl1Key, &rl1); err != nil {
			t.Fatalf("Unable to load ratelimit with key %+v error: %+v", rl1Key, err)
		}
		// check the owner reference
		if len(rl1.ObjectMeta.OwnerReferences) != 1 || rl1.ObjectMeta.OwnerReferences[0].Name != "api1" {
			t.Fatalf("Unexpected owner reference found in ratelimit rl1:  %+v", rl1)
		}

		// Get Both apis
		api1Key := types.NamespacedName{
			Name:      "api1",
			Namespace: namespace,
		}
		var api1 dpv1alpha1.API
		if err := suite.Client.Get(context.TODO(), api1Key, &api1); err != nil {
			t.Fatalf("Unable to load api with key %+v error: %+v", api1Key, err)
		}

		api2Key := types.NamespacedName{
			Name:      "api2",
			Namespace: namespace,
		}
		var api2 dpv1alpha1.API
		if err := suite.Client.Get(context.TODO(), api2Key, &api2); err != nil {
			t.Fatalf("Unable to load api with key %+v error: %+v", api2Key, err)
		}

		// Delete api1
		if err := suite.Client.Delete(context.TODO(), &api1); err != nil {
			t.Fatalf("Unable to delete api with key %+v error: %+v", api1Key, err)
		}
		// Wait 5 seconds
		time.Sleep(5 * time.Second)

		// Verify hr2 has only one parent in the ownerReferences
		if err := suite.Client.Get(context.TODO(), hr2Key, &hr2); err != nil {
			t.Fatalf("Unable to load http route with key %+v error: %+v", hr2Key, err)
		}
		// check the owner reference
		if len(hr2.ObjectMeta.OwnerReferences) != 1 {
			t.Fatalf("Unexpected owner reference found in http route hr2:  %+v", hr2)
		}
		namesToCheck = []string{"api2"}
		found = false
		for _, ownerRef := range hr2.ObjectMeta.OwnerReferences {
			foundLocal := false
			for _, name := range namesToCheck {
				if ownerRef.Name == name {
					foundLocal = true
				}
			}
			found = foundLocal
		}
		if !found {
			t.Fatalf("Unexpected owner reference found in http route hr2:  %+v", hr2)
		}

	},
}
