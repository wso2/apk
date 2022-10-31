package testutils

const SampleTestData = "testData"
const SampleTestSwaggerFile = "SampleSwagger.yaml"
const SampleCTestorruptedSwaggerFile = "SampleCorruptedSwagger.yaml"
const APIVersion = "1.0.0"
const APIName = "petstore-test"
const BackendServiceURL = "http://httpbin.default.svc.cluster.local:80/api/v3"
const CorruptedBackendServiceURL = "httpbin.default.svc.cluster.local"
const HttpRouteConfigFile = "HTTPRouteConfig.yaml"
const ConfigMapFile = "ConfigMap.yaml"

const testDir = "integration"

const SampleHttpRouteConfigWithDefaultSwagger = testDir + "/testdata/SampleHttpRouteConfigWithDefaultSwagger.yaml"
const SampleConfigMapWithDefaultSwagger = testDir + "/testdata/SampleConfigMapWithDefaultSwagger.yaml"
