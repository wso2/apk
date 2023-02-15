package integration_test

import (
	"strings"
	"testing"

	"github.com/wso2/apk/test/integration/integration/tests"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/gateway-api/conformance/utils/flags"
	gwapisuite "sigs.k8s.io/gateway-api/conformance/utils/suite"
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

	// TODO(Amila): uncomment after operator package is moved from internal to pkg directory
	// dpv1alpha1.Install(client.Scheme())
	supportedFeatures := parseSupportedFeatures(*flags.SupportedFeatures)
	exemptFeatures := parseSupportedFeatures(*flags.ExemptFeatures)
	for feature := range exemptFeatures {
		supportedFeatures[feature] = false
	}

	t.Logf("Running conformance tests with %s GatewayClass\n cleanup: %t\n debug: %t\n supported features: [%v]\n exempt features: [%v]",
		*flags.GatewayClassName, *flags.CleanupBaseResources, *flags.ShowDebug, *flags.SupportedFeatures, *flags.ExemptFeatures)

	cSuite := suite.New(suite.Options{
		Client:               client,
		GatewayClassName:     *flags.GatewayClassName,
		Debug:                *flags.ShowDebug,
		CleanupBaseResources: *flags.CleanupBaseResources,
		SupportedFeatures:    supportedFeatures,
	})
	cSuite.Setup(t)
	cSuite.Run(t, tests.IntegrationTests)
}

// parseSupportedFeatures parses flag arguments and converts the string to
// map[suite.SupportedFeature]bool
func parseSupportedFeatures(f string) map[gwapisuite.SupportedFeature]bool {
	res := map[gwapisuite.SupportedFeature]bool{}
	for _, value := range strings.Split(f, ",") {
		res[gwapisuite.SupportedFeature(value)] = true
	}
	return res
}
