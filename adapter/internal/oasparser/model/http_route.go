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

	"github.com/google/uuid"
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
	var endPoints []Endpoint
	var policies = OperationPolicies{}
	hasPolicies := false
	for _, rule := range httpRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			hasPolicies = true
			switch filter.Type {
			case gwapiv1b1.HTTPRouteFilterURLRewrite:
				policyParameters := make(map[string]interface{})
				policyParameters[constants.RewritePathType] = filter.URLRewrite.Path.Type
				policyParameters[constants.IncludeQueryParams] = true

				switch filter.URLRewrite.Path.Type {
				case gwapiv1b1.FullPathHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = *filter.URLRewrite.Path.ReplaceFullPath
					break
				case gwapiv1b1.PrefixMatchHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = *filter.URLRewrite.Path.ReplacePrefixMatch
					break
				}

				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1b1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
			}
		}

		for _, match := range rule.Matches {
			resourcePath := *match.Path.Value
			resources = append(resources, &Resource{path: resourcePath,
				methods:       getAllowedOperations(match.Method, policies),
				pathMatchType: *match.Path.Type,
				hasPolicies:   hasPolicies})
		}
		for _, backend := range rule.BackendRefs {
			endPoints = append(endPoints,
				Endpoint{Host: fmt.Sprintf("%s.%s", backend.Name,
					operatorutils.GetNamespace(backend.Namespace, httpRoute.Namespace)),
					Port: uint32(*backend.Port)})
		}
	}
	if isProd {
		swagger.productionEndpoints = &EndpointCluster{
			EndpointPrefix: constants.ProdClustersConfigNamePrefix,
			Endpoints:      endPoints,
		}
	} else {
		swagger.sandboxEndpoints = &EndpointCluster{
			EndpointPrefix: constants.SandClustersConfigNamePrefix,
			Endpoints:      endPoints,
		}
	}
	swagger.resources = resources
	return nil
}

// getAllowedOperations retuns a list of allowed operatons, if httpMethod is not specified then all methods are allowed.
func getAllowedOperations(httpMethod *gwapiv1b1.HTTPMethod, policies OperationPolicies) []*Operation {
	if httpMethod != nil {
		return []*Operation{{iD: uuid.New().String(), method: string(*httpMethod), policies: policies}}
	}
	return []*Operation{{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodGet), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPost), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodDelete), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPatch), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPut), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodHead), policies: policies},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodOptions), policies: policies}}
}

// SetInfoAPICR populates ID, ApiType, Version and XWso2BasePath of mgwSwagger.
func (swagger *MgwSwagger) SetInfoAPICR(api dpv1alpha1.API) error {
	swagger.UUID = string(api.ObjectMeta.UID)
	//TODO (amali) why id = APIDisplayName?
	swagger.id = api.Spec.APIDisplayName
	swagger.title = api.Spec.APIDisplayName
	swagger.apiType = api.Spec.APIType
	swagger.version = api.Spec.APIVersion
	swagger.xWso2Basepath = api.Spec.Context
	swagger.OrganizationID = api.Spec.Organization
	return nil
}

// ValidateIR validates the mgwSwagger based on the data required for xDS update.
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
