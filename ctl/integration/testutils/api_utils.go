package testutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BLasan/APKCTL-Demo/CTL/integration/base"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func AddNewAPIWithSwagger(t *testing.T, swagerPath string) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "api"
	apiVersion := APIVersion
	out, err := deployAPIWithSwagger(t, apiName, apiVersion, swagerPath)

	assert.Nil(t, err, "Error while deploying API")
	assert.Contains(t, out, "Successfully deployed")

	args := []string{"get", "httproute"}

	apiOut, err := base.ExecuteKubernetesCommands(args...)

	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, apiOut, apiName+"-"+apiVersion)

	t.Cleanup(func() {
		removeAPI(t, apiName, apiVersion)
	})
}

func CreateNewAPIFromSwaggerWithDryRun(t *testing.T, swagerPath string) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := createAPIWithSwagger(t, apiName, apiVersion, swagerPath)

	assert.Nil(t, err, "Error while creating API from Swagger File")
	assert.Contains(t, out, "Successfully created API project with HttpRouteConfig and ConfigMap files!")

	apiProjectDir := base.GetExportedPathFromOutput(out)

	httprouteconfig := filepath.Join(base.RelativeBinaryPath, apiProjectDir, HttpRouteConfigFile)

	configmap := filepath.Join(base.RelativeBinaryPath, apiProjectDir, ConfigMapFile)

	assert.True(t, base.IsFileAvailable(t, httprouteconfig), "HttpRouteConfig is not available")
	assert.True(t, base.IsFileAvailable(t, configmap), "ConfigMap is not available")

	// removeAll(t, filepath.Join(base.RelativeBinaryPath, apiProjectDir+"../../../"))

	// removeFile(t, httprouteconfig)
	// removeFile(t, configmap)
}

func AddNewAPIWithBackendServiceURL(t *testing.T) {
	t.Helper()
	apiName := APIName
	apiVersion := APIVersion
	out, err := deployAPIWithBackendServiceURL(t, apiName, apiVersion, BackendServiceURL)

	assert.Nil(t, err, "Error while deploying API")
	assert.Contains(t, out, "Successfully deployed")
}

func CreateNewAPIFromBackendServiceURLWithDryRun(t *testing.T) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := createAPIWithBackendServiceURL(t, apiName, apiVersion, BackendServiceURL)

	assert.Nil(t, err, "Error while creating AP from Backend Service URL")
	assert.Contains(t, out, "Successfully created")

	apiProjectDir := base.GetExportedPathFromOutput(out)

	httprouteconfig := filepath.Join(base.RelativeBinaryPath, apiProjectDir, HttpRouteConfigFile)
	configmap := filepath.Join(base.RelativeBinaryPath, apiProjectDir, ConfigMapFile)

	assert.True(t, base.IsFileAvailable(t, httprouteconfig), "HttpRouteConfig is not available")
	assert.True(t, base.IsFileAvailable(t, configmap), "ConfigMap is not available")

	// absPath := filepath.Join(base.RelativeBinaryPath, apiProjectDir+"../../../")

	// t.Cleanup(func() {
	// 	removeAll(t, absPath)
	// })

	// removeFile(t, httprouteconfig)
	// removeFile(t, configmap)

}

func CreateAPIWithoutBackendService(t *testing.T) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	_, err := createAPIWithoutBackendService(t, apiName, apiVersion)

	assert.NotNil(t, err, "Either Swagger Definition or Backend Service URL should be provided")
}

func CreateAPIWithCorruptedSwaggerDefinition(t *testing.T, swaggerpath string) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	_, err := createAPIWithSwagger(t, apiName, apiVersion, swaggerpath)

	assert.NotNil(t, err, "Swagger Definition is corrupted")
}

func CreateAPIWithCorruptedBackendServiceURL(t *testing.T) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	_, err := createAPIWithBackendServiceURL(t, apiName, apiVersion, CorruptedBackendServiceURL)

	assert.NotNil(t, err, "Valid Backend Service URL should be provided")
}

func ValidateInstallAPKComponents(t *testing.T) {
	t.Helper()

	out, err := installAPK(t)

	if err != nil {
		time.Sleep(5 * time.Second)
		out, err = installAPK(t)
	}

	assert.Nil(t, err, "Error while installing APK components")
	assert.Contains(t, out, "All Done! We have configured APK to help you build and manage APIs with ease.")
}

func ValidateUninstallAPKComponents(t *testing.T) {
	t.Helper()

	out, err := uninstallAPK(t)

	assert.Nil(t, err, "Error while uninstalling APK components")
	assert.Contains(t, out, "Uninstallation completed!")
}

func ValidateAPIConfigFiles(t *testing.T) {
	t.Helper()
	apiName := APIName
	apiVersion := APIVersion
	out, err := createAPIWithBackendServiceURL(t, apiName, apiVersion, BackendServiceURL)

	assert.Nil(t, err, "Error while creating API from Backend Service URL")
	assert.Contains(t, out, "Successfully created")

	apiProjectDir := base.GetExportedPathFromOutput(out)
	httprouteconfig := filepath.Join(apiProjectDir, HttpRouteConfigFile)
	configmap := filepath.Join(apiProjectDir, ConfigMapFile)

	validateAPIRelatedFiles(t, httprouteconfig, configmap)

	// t.Cleanup(func() {
	// 	absPath := filepath.Join(base.RelativeBinaryPath, apiProjectDir+"../../../")
	// 	removeAll(t, absPath)
	// })

}

func installAPK(t *testing.T) (string, error) {
	output, err := base.Execute(t, "install", "platform", "--verbose")
	return output, err
}

func uninstallAPK(t *testing.T) (string, error) {
	output, err := base.Execute(t, "uninstall", "platform", "--verbose")
	return output, err
}

// Creates API from swagger file
func deployAPIWithSwagger(t *testing.T, apiName, apiversion, swagger string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--version", apiversion, "-f", swagger, "--verbose")
	return output, err
}

// Creates API from the Backend Service URL
func deployAPIWithBackendServiceURL(t *testing.T, apiName, apiversion, backendURL string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--service-url", backendURL, "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

func validateAPIRelatedFiles(t *testing.T, httprouteconfig, configmap string) {
	httprouteconfigContent := readAPIRelatedFiles(t, httprouteconfig)
	configmapContent := readAPIRelatedFiles(t, configmap)

	httprouteconfigContentExpected := readAPIRelatedFiles(t, SampleHttpRouteConfigWithDefaultSwagger)
	configmapContentExpected := readAPIRelatedFiles(t, SampleConfigMapWithDefaultSwagger)

	assert.Equal(t, httprouteconfigContent, httprouteconfigContentExpected)
	assert.Equal(t, configmapContent, configmapContentExpected)
}

// Creates API from swagger file
func createAPIWithSwagger(t *testing.T, apiName, apiversion, swagger string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--version", apiversion, "-f", swagger, "--dry-run", "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

// Creates API from the Backend Service URL
func createAPIWithBackendServiceURL(t *testing.T, apiName, apiversion, backendURL string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--service-url", backendURL, "--dry-run", "--verbose")
	// t.Cleanup(func() {
	// 	removeAPI(t, apiName, apiversion)
	// })
	return output, err
}

// Creates API without providing a swagger or backend service URL
func createAPIWithoutBackendService(t *testing.T, apiName, apiversion string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--dry-run", "--verbose")
	return output, err
}

func removeAPI(t *testing.T, apiname, version string) {
	base.Execute(t, "delete", "api", apiname, "--version", version)
}

// func removeAll(t *testing.T, dirname string) {
// 	t.Log("testutils.removeAll() - dir path:", dirname)
// 	if _, err := os.Stat(dirname); err == nil {
// 		err := os.RemoveAll(dirname)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }

func removeFile(t *testing.T, filename string) {
	t.Log("testutils.removeFile() - file path:", filename)
	if _, err := os.Stat(filename); err == nil {
		err := os.Remove(filename)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func readAPIRelatedFiles(t *testing.T, filename string) map[string]interface{} {

	content, err := ioutil.ReadFile(filepath.Join(base.RelativeBinaryPath, filename))

	if err != nil {
		t.Fatal(err)
	}

	yamlContent := make(map[string]interface{})
	err = yaml.Unmarshal(content, &yamlContent)

	if err != nil {
		t.Fatal(err)
	}

	return yamlContent

}
