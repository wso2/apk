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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
)

func TestConcatRateLimitPolicies(t *testing.T) {
	type testItem struct {
		schemeUpSpec   dpv1alpha1.RateLimitPolicySpec
		schemeDownSpec dpv1alpha1.RateLimitPolicySpec
		result         dpv1alpha1.RateLimitPolicySpec
		message        string
	}

	schemeUp := &dpv1alpha1.RateLimitPolicy{}
	schemeDown := &dpv1alpha1.RateLimitPolicy{}
	resultScheme := &dpv1alpha1.RateLimitPolicy{}

	dataItems := []testItem{
		{
			schemeUpSpec: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			result: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			message: "When API level override and Resource level override policies both provided",
		},
		{
			schemeUpSpec: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			result: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			message: "When API level override and Resource level default policies both provided",
		},
		{
			schemeUpSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			result: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			message: "When API level default and Resource level override policies both provided",
		},
		{
			schemeUpSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Day",
					},
				},
			},
			result: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
			},
			message: "When API level default and Resource level default policies both provided",
		},
		{
			schemeUpSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    10,
						SpanUnit: "Minute",
					},
				},
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Application",
					API: dpv1alpha1.APIRateLimit{
						Count:    20,
						SpanUnit: "Second",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.RateLimitPolicySpec{
				Default: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Api",
					API: dpv1alpha1.APIRateLimit{
						Count:    30,
						SpanUnit: "Day",
					},
				},
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Subscription",
					API: dpv1alpha1.APIRateLimit{
						Count:    40,
						SpanUnit: "Hour",
					},
				},
			},
			result: dpv1alpha1.RateLimitPolicySpec{
				Override: dpv1alpha1.RateLimitAPIPolicy{
					Type: "Subscription",
					API: dpv1alpha1.APIRateLimit{
						Count:    40,
						SpanUnit: "Hour",
					},
				},
			},
			message: "When both API level and Resource level both override and default policies provided",
		},
	}

	for _, item := range dataItems {
		schemeUp.Spec = item.schemeUpSpec
		schemeDown.Spec = item.schemeDownSpec
		resultScheme.Spec = item.result
		actualResult := concatRateLimitPolicies(schemeUp, schemeDown)
		assert.Equal(t, resultScheme, actualResult, item.message)
	}
}

func TestConcatAPIPolicies(t *testing.T) {

	type testItem struct {
		schemeUpSpec   dpv1alpha1.APIPolicySpec
		schemeDownSpec dpv1alpha1.APIPolicySpec
		result         dpv1alpha1.APIPolicySpec
		message        string
	}

	schemeUp := &dpv1alpha1.APIPolicy{}
	schemeDown := &dpv1alpha1.APIPolicy{}
	resultScheme := &dpv1alpha1.APIPolicy{}

	dataItems := []testItem{
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			message: "only schemeUp override policies should be provided",
		},
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.APIPolicySpec{
				Default: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			message: "only schemeUp override policies should be provided",
		},
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Default: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			message: "only schemeDown override policies should be provided",
		},
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Default: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			schemeDownSpec: dpv1alpha1.APIPolicySpec{
				Default: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"c", "d"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q2", Value: "val2"},
						},
						RemoveAll: "false",
					},
				},
			},
			message: "only schemeDown default policies should be provided",
		},
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			message: "only schemeUp override policies is provided",
		},
		{
			schemeUpSpec: dpv1alpha1.APIPolicySpec{
				Default: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			result: dpv1alpha1.APIPolicySpec{
				Override: dpv1alpha1.PolicySpec{
					RequestQueryModifier: dpv1alpha1.RequestQueryModifier{
						Remove: []string{"a", "b"},
						Add: []dpv1alpha1.HTTPQuery{
							{Name: "q1", Value: "val1"},
						},
						RemoveAll: "true",
					},
				},
			},
			message: "only schemeUp default policies is provided",
		},
	}
	for _, item := range dataItems {
		schemeUp.Spec = item.schemeUpSpec
		schemeDown.Spec = item.schemeDownSpec
		resultScheme.Spec = item.result
		actualResult := concatAPIPolicies(schemeUp, schemeDown)
		assert.Equal(t, resultScheme, actualResult, item.message)
	}
}
