/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tetratelabs/multierror"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	operatorutils "github.com/wso2/apk/adapter/internal/operator/utils"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// SetInfoHTTPRouteCR populates resources and endpoints of mgwSwagger. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (swagger *MgwSwagger) SetInfoHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, isProd bool) error {
	var resources []*Resource
	var endpointCluster EndpointCluster
	var endPoints []Endpoint
	for _, rule := range httpRoute.Spec.Rules {
		for _, match := range rule.Matches {
			resourcePath, err := swagger.trimBasePath(*match.Path.Value)
			if err != nil {
				return fmt.Errorf("error parsing resource path: %v", err)
			}
			resources = append(resources, &Resource{path: resourcePath, methods: []*Operation{{iD: ":method", method: "GET"}}})
		}
		for _, backend := range rule.BackendRefs {
			endPoints = append(endPoints,
				Endpoint{Host: fmt.Sprintf("%s.%s", backend.Name,
					operatorutils.GetNamespace(backend.Namespace, httpRoute.Namespace)),
					Port: uint32(*backend.Port)})
		}
	}
	endpointCluster = EndpointCluster{
		EndpointPrefix: constants.ProdClustersConfigNamePrefix,
		Endpoints:      endPoints,
	}
	if isProd {
		swagger.productionEndpoints = &endpointCluster
	} else {
		swagger.sandboxEndpoints = &endpointCluster
	}
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
	if (swagger.productionEndpoints != nil && len(swagger.productionEndpoints.Endpoints) == 0) ||
		(swagger.sandboxEndpoints != nil && len(swagger.sandboxEndpoints.Endpoints) == 0) {
		errs = multierror.Append(errs, errors.New("no endpoints provided"))
	}
	if len(swagger.resources) == 0 {
		errs = multierror.Append(errs, errors.New("no resources found"))
	}
	return errs
}

func (swagger *MgwSwagger) trimBasePath(path string) (string, error) {
	if strings.Compare(path, swagger.GetXWso2Basepath()) == 0 {
		return "/", nil
	}
	resourcePath := strings.TrimPrefix(path, swagger.GetXWso2Basepath())
	if strings.Compare(path, resourcePath) == 0 {
		return "", fmt.Errorf("basepath mismatch: %v:%v", swagger.GetXWso2Basepath(), path)
	}
	return resourcePath, nil
}
