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

package suite

import (
	"testing"
  "time"

	"github.com/wso2/apk/test/integration/integration/utils/kubernetes"
	"github.com/wso2/apk/test/integration/integration/utils/roundtripper"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
	gwapisuite "sigs.k8s.io/gateway-api/conformance/utils/suite"
)

// IntegrationTestSuite defines the test suite used to run Gateway API
// conformance tests.
type IntegrationTestSuite struct {
	Client            client.Client
	RoundTripper      roundtripper.RoundTripper
	GatewayClassName  string
	ControllerName    string
	Debug             bool
	Cleanup           bool
	BaseManifests     string
	Applier           kubernetes.Applier
	SupportedFeatures map[gwapisuite.SupportedFeature]bool
	TimeoutConfig     config.TimeoutConfig
	SkipTests         sets.Set[string]
}

// Options can be used to initialize a IntegrationTestSuite.
type Options struct {
	Client           client.Client
	GatewayClassName string
	Debug            bool
	RoundTripper     roundtripper.RoundTripper
	BaseManifests    string
	NamespaceLabels  map[string]string
	// ValidUniqueListenerPorts maps each listener port of each Gateway in the
	// manifests to a valid, unique port. There must be as many
	// ValidUniqueListenerPorts as there are listeners in the set of manifests.
	// For example, given two Gateways, each with 2 listeners, there should be
	// four ValidUniqueListenerPorts.
	// If empty or nil, ports are not modified.
	ValidUniqueListenerPorts []v1beta1.PortNumber

	// CleanupBaseResources indicates whether or not the base test
	// resources such as Gateways should be cleaned up after the run.
	CleanupBaseResources bool
	SupportedFeatures    map[gwapisuite.SupportedFeature]bool
	TimeoutConfig        config.TimeoutConfig
	// SkipTests contains all the tests not to be run and can be used to opt out
	// of specific tests
	SkipTests []string
}

// StandardCoreFeatures are the features that are required to be conformant with
// the Core API features that are part of the Standard release channel.
var StandardCoreFeatures = map[gwapisuite.SupportedFeature]bool{
	gwapisuite.SupportReferenceGrant: true,
}

// New returns a new IntegrationTestSuite.
func New(s Options) *IntegrationTestSuite {
	config.SetupTimeoutConfig(&s.TimeoutConfig)

	roundTripper := s.RoundTripper
	if roundTripper == nil {
		roundTripper = &roundtripper.DefaultRoundTripper{Debug: s.Debug, TimeoutConfig: s.TimeoutConfig}
	}

	if s.SupportedFeatures == nil {
		s.SupportedFeatures = StandardCoreFeatures
	} else {
		for feature, val := range StandardCoreFeatures {
			if _, ok := s.SupportedFeatures[feature]; !ok {
				s.SupportedFeatures[feature] = val
			}
		}
	}

	suite := &IntegrationTestSuite{
		Client:           s.Client,
		RoundTripper:     roundTripper,
		GatewayClassName: s.GatewayClassName,
		Debug:            s.Debug,
		Cleanup:          s.CleanupBaseResources,
		BaseManifests:    s.BaseManifests,
		Applier: kubernetes.Applier{
			NamespaceLabels:          s.NamespaceLabels,
			ValidUniqueListenerPorts: s.ValidUniqueListenerPorts,
		},
		SupportedFeatures: s.SupportedFeatures,
		TimeoutConfig:     s.TimeoutConfig,
		SkipTests:         sets.New(s.SkipTests...),
	}

	// apply defaults
	if suite.BaseManifests == "" {
		suite.BaseManifests = "base/manifests.yaml"
	}

	return suite
}

// Setup ensures the base resources required for conformance tests are installed
// in the cluster. It also ensures that all relevant resources are ready.
func (suite *IntegrationTestSuite) Setup(t *testing.T) {
	// TODO (Amila): Revisit when gateway resource support is added
	suite.ControllerName = "wso2.com/apk-gateway-default"

	suite.Applier.GatewayClass = suite.GatewayClassName
	suite.Applier.ControllerName = suite.ControllerName

	t.Logf("Test Setup: Applying base manifests")
	suite.Applier.MustApplyWithCleanup(t, suite.Client, suite.TimeoutConfig, suite.BaseManifests, suite.Cleanup)

	namespaces := []string{
		"gateway-integration-test-infra",
	}
	kubernetes.NamespacesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, namespaces)
}

// Run runs the provided set of conformance tests.
func (suite *IntegrationTestSuite) Run(t *testing.T, tests []IntegrationTest) {
	for _, test := range tests {
		t.Run(test.ShortName, func(t *testing.T) {
			test.Run(t, suite)
		})
	}
}

// IntegrationTest is used to define each individual conformance test.
type IntegrationTest struct {
	ShortName   string
	Description string
	Features    []gwapisuite.SupportedFeature
	Manifests   []string
	Slow        bool
	Parallel    bool
	Test        func(*testing.T, *IntegrationTestSuite)
}

// Run runs an individual tests, applying and cleaning up the required manifests
// before calling the Test function.
func (test *IntegrationTest) Run(t *testing.T, suite *IntegrationTestSuite) {
	if test.Parallel {
		t.Parallel()
	}

	for _, manifestLocation := range test.Manifests {
		t.Logf("Applying %s", manifestLocation)
		suite.Applier.MustApplyWithCleanup(t, suite.Client, suite.TimeoutConfig, manifestLocation, true)
	}

	test.Test(t, suite)
}

// WaitForNextMinute wait until next clock minute starts
func WaitForNextMinute(t *testing.T) {
	additionalSeconds := 5
	now := time.Now()
	nextMinute := now.Add(time.Minute).Truncate(time.Minute)
	nextTime := nextMinute.Add(time.Duration(additionalSeconds) * time.Second)
	durationToWait := nextTime.Sub(now)
	
	if int(durationToWait.Seconds()) > 15 {
		t.Logf("Not waiting for next minute as we have enough time in this minute. Current time (%v)\n", now)
		return
	}
	t.Logf("Waiting for the next minute and %d seconds to start (%v)... current time (%v)\n", additionalSeconds, nextTime.Format("15:04:05"), now)
	time.Sleep(durationToWait)
	t.Logf("Next minute and %d seconds have started! Current time: (%v)", additionalSeconds, time.Now())
}
