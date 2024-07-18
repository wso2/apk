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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package translator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	corsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	hcmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/types"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	registerHTTPFilter(&cors{})
}

type cors struct {
}

var _ httpFilter = &cors{}

// patchHCM builds and appends the CORS Filter to the HTTP Connection Manager if
// applicable.
func (*cors) patchHCM(
	mgr *hcmv3.HttpConnectionManager,
	irListener *ir.HTTPListener) error {
	if mgr == nil {
		return errors.New("hcm is nil")
	}

	if irListener == nil {
		return errors.New("ir listener is nil")
	}

	if !listenerContainsCORS(irListener) {
		return nil
	}

	// Return early if filter already exists.
	for _, httpFilter := range mgr.HttpFilters {
		if httpFilter.Name == wellknown.CORS {
			return nil
		}
	}

	corsFilter, err := buildHCMCORSFilter()
	if err != nil {
		return err
	}

	// Ensure the CORS filter is the first one in the filter chain.
	mgr.HttpFilters = append([]*hcmv3.HttpFilter{corsFilter}, mgr.HttpFilters...)

	return nil
}

// buildHCMCORSFilter returns a CORS filter from the provided IR listener.
func buildHCMCORSFilter() (*hcmv3.HttpFilter, error) {
	corsProto := &corsv3.Cors{}

	corsAny, err := anypb.New(corsProto)
	if err != nil {
		return nil, err
	}

	return &hcmv3.HttpFilter{
		Name: wellknown.CORS,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: corsAny,
		},
	}, nil
}

// listenerContainsCORS returns true if the provided listener has CORS
// policies attached to its routes.
func listenerContainsCORS(irListener *ir.HTTPListener) bool {
	if irListener == nil {
		return false
	}

	for _, route := range irListener.Routes {
		if route.CORS != nil {
			return true
		}
	}

	return false
}

// patchRoute patches the provided route with the CORS config if applicable.
func (*cors) patchRoute(route *routev3.Route, irRoute *ir.HTTPRoute) error {
	if route == nil {
		return errors.New("xds route is nil")
	}
	if irRoute == nil {
		return errors.New("ir route is nil")
	}
	if irRoute.CORS == nil {
		return nil
	}

	filterCfg := route.GetTypedPerFilterConfig()
	if _, ok := filterCfg[wellknown.CORS]; ok {
		// This should not happen since this is the only place where the CORS
		// filter is added in a route.
		return fmt.Errorf("route already contains cors config: %+v", route)
	}

	var (
		allowOrigins     []*matcherv3.StringMatcher
		allowMethods     string
		allowHeaders     string
		exposeHeaders    string
		maxAge           string
		allowCredentials *wrappers.BoolValue
	)

	//nolint:gocritic

	for _, origin := range irRoute.CORS.AllowOrigins {
		allowOrigins = append(allowOrigins, buildXdsStringMatcher(origin))
	}

	allowMethods = strings.Join(irRoute.CORS.AllowMethods, ", ")
	allowHeaders = strings.Join(irRoute.CORS.AllowHeaders, ", ")
	exposeHeaders = strings.Join(irRoute.CORS.ExposeHeaders, ", ")
	if irRoute.CORS.MaxAge != nil {
		maxAge = strconv.Itoa(int(irRoute.CORS.MaxAge.Seconds()))
	}
	allowCredentials = &wrappers.BoolValue{Value: irRoute.CORS.AllowCredentials}

	routeCfgProto := &corsv3.CorsPolicy{
		AllowOriginStringMatch: allowOrigins,
		AllowMethods:           allowMethods,
		AllowHeaders:           allowHeaders,
		ExposeHeaders:          exposeHeaders,
		MaxAge:                 maxAge,
		AllowCredentials:       allowCredentials,
	}

	routeCfgAny, err := anypb.New(routeCfgProto)
	if err != nil {
		return err
	}

	if filterCfg == nil {
		route.TypedPerFilterConfig = make(map[string]*anypb.Any)
	}

	route.TypedPerFilterConfig[wellknown.CORS] = routeCfgAny

	return nil
}

func (c *cors) patchResources(*types.ResourceVersionTable, []*ir.HTTPRoute) error {
	return nil
}
