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
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"text/template"

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"gopkg.in/yaml.v2"
)

const (
	policyCCGateway string = "ChoreoConnect"
)

var (
	// policyDefFuncMap is a map of functions used in policy definitions
	policyDefFuncMap = template.FuncMap{
		// isParamExists checks the key is exists in the params map, this will not consider the value of the param
		// if the go templated "{{ if .param }}" is used, that will consider the
		// value of the param (if value is a zero value, it consider as not exists)
		"isParamExists": func(m map[string]interface{}, key string) (ok bool) {
			_, ok = m[key]
			return
		},
	}
)

// PolicyFlow holds list of Policies in a operation (in one flow: In, Out or Fault)
type PolicyFlow string

const (
	policyInFlow    PolicyFlow = "request"
	policyOutFlow   PolicyFlow = "response"
	policyFaultFlow PolicyFlow = "fault"
)

// PolicyContainerMap maps PolicyName -> PolicyContainer
type PolicyContainerMap map[string]PolicyContainer

// PolicyContainer holds the definition and specification of policy
type PolicyContainer struct {
	Specification PolicySpecification
	Definition    PolicyDefinition
}

// PolicySpecification holds policy specification from ./Policy/<policy>.yaml files
type PolicySpecification struct {
	Type    string `yaml:"type" json:"type"`
	Version string `yaml:"version" json:"version"`
	Data    struct {
		Name              string   `yaml:"name"`
		Version           string   `yaml:"version"`
		ApplicableFlows   []string `yaml:"applicableFlows"`
		SupportedGateways []string `yaml:"supportedGateways"`
		SupportedAPITypes []string `yaml:"supportedApiTypes"`
		MultipleAllowed   bool     `yaml:"multipleAllowed"`
		PolicyAttributes  []struct {
			Name            string `yaml:"name"`
			ValidationRegex string `yaml:"validationRegex,omitempty"`
			Type            string `yaml:"type"`
			DefaultValue    string `yaml:"defaultValue"`
			Required        bool   `yaml:"required,omitempty"`
		} `yaml:"policyAttributes"`
	}
}

// PolicyDefinition holds the content of policy definition which is rendered from ./Policy/<policy>.gotmpl files
type PolicyDefinition struct {
	Definition struct {
		Action     string                 `yaml:"action"`
		Parameters map[string]interface{} `yaml:"parameters"`
	} `yaml:"definition"`
	RawData []byte `yaml:"-"`
}

// GetFormattedOperationalPolicies returns formatted, Choreo Connect policy from a user templated policy
// here, the struct swagger is only used for logging purpose, in case if we introduce logger context to get org ID, API ID, we can remove it from here
func (p PolicyContainerMap) GetFormattedOperationalPolicies(policies OperationPolicies, swagger *MgwSwagger) (OperationPolicies, error) {
	fmtPolicies := OperationPolicies{}

	for _, policy := range policies.Request {
		if fmtPolicy, err := p.getFormattedPolicyFromTemplated(policy, policyInFlow, swagger); err == nil {
			fmtPolicies.Request = append(fmtPolicies.Request, fmtPolicy)
			loggers.LoggerOasparser.Debugf("Applying operation policy %q in request flow, for API %q in org %q, formatted policy %v",
				policy.GetFullName(), swagger.GetID(), swagger.OrganizationID, fmtPolicy)
		} else {
			return fmtPolicies, err
		}
	}

	for _, policy := range policies.Response {
		if fmtPolicy, err := p.getFormattedPolicyFromTemplated(policy, policyOutFlow, swagger); err == nil {
			fmtPolicies.Response = append(fmtPolicies.Response, fmtPolicy)
			loggers.LoggerOasparser.Debugf("Applying operation policy %q in response flow, for API %q in org %q, formatted policy %v",
				policy.GetFullName(), swagger.GetID(), swagger.OrganizationID, fmtPolicy)
		} else {
			return fmtPolicies, err
		}
	}

	for _, policy := range policies.Fault {
		if fmtPolicy, err := p.getFormattedPolicyFromTemplated(policy, policyFaultFlow, swagger); err == nil {
			fmtPolicies.Fault = append(fmtPolicies.Fault, fmtPolicy)
			loggers.LoggerOasparser.Debugf("Applying operation policy %q in fault flow, for API %q in org %q, formatted policy %v",
				policy.GetFullName(), swagger.GetID(), swagger.OrganizationID, fmtPolicy)
		} else {
			return fmtPolicies, err
		}
	}

	return fmtPolicies, nil
}

// getFormattedPolicyFromTemplated returns formatted, Choreo Connect policy from a user templated policy
func (p PolicyContainerMap) getFormattedPolicyFromTemplated(policy Policy, flow PolicyFlow, swagger *MgwSwagger) (Policy, error) {
	policyFullName := policy.GetFullName()
	spec := p[policyFullName].Specification
	if err := spec.validatePolicy(policy, flow); err != nil {
		swagger.GetID()
		loggers.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Operation policy validation failed for API %q in org %q:, policy %q: %v", swagger.GetID(), swagger.OrganizationID, policyFullName, err),
			Severity:  logging.MINOR,
			ErrorCode: 2204,
		})
		return policy, err
	}

	defRaw := p[policyFullName].Definition.RawData
	t, err := template.New("policy-def").Funcs(policyDefFuncMap).Parse(string(defRaw))
	if err != nil {
		loggers.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error parsing the operation policy definition %q into go template of the API %q in org %q: %v", policyFullName, swagger.GetID(), swagger.OrganizationID, err),
			Severity:  logging.MINOR,
			ErrorCode: 2205,
		})
		return Policy{}, err
	}

	var out bytes.Buffer
	err = t.Execute(&out, policy.Parameters)
	if err != nil {
		loggers.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error parsing operation policy definition %q of the API %q in org %q: %v", policyFullName, swagger.GetID(), swagger.OrganizationID, err),
			Severity:  logging.MINOR,
			ErrorCode: 2206,
		})
		return Policy{}, err
	}

	def := PolicyDefinition{}
	if err := yaml.Unmarshal(out.Bytes(), &def); err != nil {
		loggers.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error parsing formalized operation policy definition %q into yaml of the API %q in org %q: %v", policyFullName, swagger.GetID(), swagger.OrganizationID, err),
			Severity:  logging.MINOR,
			ErrorCode: 2207,
		})
		return Policy{}, err
	}

	// Update templated policy itself and return, not updating a pointer to keep the original template values as it is.
	policy.Parameters = def.Definition.Parameters
	policy.Action = def.Definition.Action

	// Fill default values
	spec.fillDefaultsInPolicy(&policy)

	// Check the API Policy supported by Choreo Connect
	// Required params may be comming from default values as defined in the policy specification
	// Hence do the validation after filling default values
	if err := validatePolicyAction(&policy); err != nil {
		loggers.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("API policy validation failed, policy: %q of the API %q in org %q: %v", policyFullName, swagger.GetID(), swagger.OrganizationID, err),
			Severity:  logging.MINOR,
			ErrorCode: 0,
		})
		return Policy{}, err
	}
	return policy, nil
}

// validatePolicy validates the given policy against the spec
func (spec *PolicySpecification) validatePolicy(policy Policy, flow PolicyFlow) error {
	if spec.Data.Name != policy.PolicyName || spec.Data.Version != policy.PolicyVersion {
		return fmt.Errorf("invalid policy specification, spec name %q:%q and policy name %q:%q mismatch",
			spec.Data.Name, spec.Data.Version, policy.PolicyName, policy.PolicyVersion)
	}
	if !arrayContains(spec.Data.ApplicableFlows, string(flow)) {
		return fmt.Errorf("policy flow %q not supported", flow)
	}
	if !arrayContains(spec.Data.SupportedGateways, policyCCGateway) {
		return errors.New("choreo connect gateway not supported")
	}

	policyPrams, ok := policy.Parameters.(map[string]interface{})
	if ok {
		for _, attrib := range spec.Data.PolicyAttributes {
			val, found := policyPrams[attrib.Name]
			if attrib.Required && !found {
				return fmt.Errorf("required paramater %q not found", attrib.Name)
			}

			switch v := val.(type) {
			case string:
				regexStr := attrib.ValidationRegex
				if regexStr != "" {
					reg, err := regexp.Compile(regexStr)
					if err != nil {
						return fmt.Errorf("invalid regex expression in policy spec %s, regex: %q", spec.Data.Name, attrib.ValidationRegex)
					}
					if !reg.MatchString(v) {
						return fmt.Errorf("invalid parameter value of attribute %q, regex match failed", attrib.Name)
					}
				}
			}
		}
	}

	return nil
}

// fillDefaultsInPolicy updates the policy with default values defined in the spec if the key is not found in the policy
func (spec *PolicySpecification) fillDefaultsInPolicy(policy *Policy) {
	if paramMap, isMap := policy.Parameters.(map[string]interface{}); isMap {
		for _, attrib := range spec.Data.PolicyAttributes {
			if _, ok := paramMap[attrib.Name]; !ok && attrib.DefaultValue != "" {
				paramMap[attrib.Name] = attrib.DefaultValue
				loggers.LoggerOasparser.Debugf("Update with policy attribute %q of policy %q with default value from spec",
					attrib.Name, policy.PolicyName)
			}
		}
		policy.Parameters = paramMap
	}
}
