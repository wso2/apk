package model

import (
	"errors"

	"github.com/tetratelabs/multierror"
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
	swagger.OrganizationID = api.Spec.Organization
	return nil
}

// Validate validates the mgwSwagger based on the data required for xDS update.
func (swagger *MgwSwagger) ValidateIR() error {
	var errs error

	if swagger.UUID == "" {
		errs = multierror.Append(errs, errors.New("api UUID not found"))
	}
	if swagger.id == "" {
		errs = multierror.Append(errs, errors.New("api ID not found"))
	}
	if swagger.apiType == "" {
		errs = multierror.Append(errs, errors.New("api type not found / invalid"))
	}
	if swagger.version == "" {
		errs = multierror.Append(errs, errors.New("api version not found"))
	}
	if swagger.xWso2Basepath == "" {
		errs = multierror.Append(errs, errors.New("api basepath not found"))
	}
	if len(swagger.productionEndpoints.Endpoints) == 0 {
		errs = multierror.Append(errs, errors.New("no production endpoints provided"))
	}
	if len(swagger.resources) == 0 {
		errs = multierror.Append(errs, errors.New("no resources found"))
	}
	return errs
}
