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

package gatewayapi

import (
	"sort"

	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
)

type XdsIRRoutes []*ir.HTTPRoute

func (x XdsIRRoutes) Len() int      { return len(x) }
func (x XdsIRRoutes) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x XdsIRRoutes) Less(i, j int) bool {

	// 1. Sort based on path match type
	// Exact > RegularExpression > PathPrefix
	if x[i].PathMatch != nil && x[i].PathMatch.Exact != nil {
		if x[j].PathMatch != nil {
			if x[j].PathMatch.SafeRegex != nil {
				return false
			}
			if x[j].PathMatch.Prefix != nil {
				return false
			}
		}
	}
	if x[i].PathMatch != nil && x[i].PathMatch.SafeRegex != nil {
		if x[j].PathMatch != nil {
			if x[j].PathMatch.Exact != nil {
				return true
			}
			if x[j].PathMatch.Prefix != nil {
				return false
			}
		}
	}
	if x[i].PathMatch != nil && x[i].PathMatch.Prefix != nil {
		if x[j].PathMatch != nil {
			if x[j].PathMatch.Exact != nil {
				return true
			}
			if x[j].PathMatch.SafeRegex != nil {
				return true
			}
		}
	}
	// Equal case

	// 2. Sort based on characters in a matching path.
	pCountI := pathMatchCount(x[i].PathMatch)
	pCountJ := pathMatchCount(x[j].PathMatch)
	if pCountI < pCountJ {
		return true
	}
	if pCountI > pCountJ {
		return false
	}
	// Equal case

	// 3. Sort based on the number of Header matches.
	hCountI := len(x[i].HeaderMatches)
	hCountJ := len(x[j].HeaderMatches)
	if hCountI < hCountJ {
		return true
	}
	if hCountI > hCountJ {
		return false
	}
	// Equal case

	// 4. Sort based on the number of Query param matches.
	qCountI := len(x[i].QueryParamMatches)
	qCountJ := len(x[j].QueryParamMatches)
	return qCountI < qCountJ
}

// sortXdsIR sorts the xdsIR based on the match precedence
// defined in the Gateway API spec.
// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.HTTPRouteRule
func sortXdsIRMap(xdsIR XdsIRMap) {
	for _, irItem := range xdsIR {
		for _, http := range irItem.HTTP {
			// descending order
			sort.Sort(sort.Reverse(XdsIRRoutes(http.Routes)))
		}
	}
}

func pathMatchCount(pathMatch *ir.StringMatch) int {
	if pathMatch != nil {
		if pathMatch.Exact != nil {
			return len(*pathMatch.Exact)
		}
		if pathMatch.SafeRegex != nil {
			return len(*pathMatch.SafeRegex)
		}
		if pathMatch.Prefix != nil {
			return len(*pathMatch.Prefix)
		}
	}
	return 0
}
