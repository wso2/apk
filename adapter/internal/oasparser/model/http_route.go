package model

import (
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/api/v1alpha1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// SetInfoHTTPRouteCR populates resources and endpoints of mgwSwagger. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (swagger *MgwSwagger) SetInfoHTTPRouteCR(httpRoute gwapiv1b1.HTTPRoute) error {
	var resources []*Resource
	var endpointCluster EndpointCluster
	var endPoints []Endpoint
	for _, rule := range httpRoute.Spec.Rules {
		for _, match := range rule.Matches {
			resources = append(resources, &Resource{path: *match.Path.Value, methods: []*Operation{{iD: ":method", method: "GET"}}})
		}
		for _, backend := range rule.BackendRefs {
			endPoints = append(endPoints, Endpoint{Host: string(backend.Name), Port: uint32(*backend.Port)})
		}
	}
	endpointCluster = EndpointCluster{
		EndpointPrefix: constants.ProdClustersConfigNamePrefix,
		Endpoints:      endPoints,
	}
	swagger.productionEndpoints = &endpointCluster
	swagger.resources = resources
	return nil
}

// SetInfoAPICR populates ID, ApiType, Version and XWso2BasePath of mgwSwagger.
func (swagger *MgwSwagger) SetInfoAPICR(api dpv1alpha1.API) error {
	swagger.UUID = string(api.ObjectMeta.UID)
	swagger.id = api.Spec.APIDisplayName
	swagger.apiType = api.Spec.APIType
	swagger.version = api.Spec.APIVersion
	swagger.xWso2Basepath = api.Spec.Context
	return nil
}
