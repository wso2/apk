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

package config

import (
	"testing"
)

func TestValidateConfigValuesForValidValues(t *testing.T) {

	validValues := []string{gatewayTypeDefault, gatewayTypeChoreo, gatewayTypeChoreoPDP}
	conf := ReadConfigs()

	for _, config := range validValues {
		propertyMap := make(map[string]string)
		propertyMap[gatewayTypeValue] = config
		conf.Analytics.Properties = propertyMap

		err := validateAnalyticsConfigs(conf)
		if err != nil {
			t.Errorf("Expected validation of '%s' to be successful, but got an error: %v", config, err)
		}
	}

	propertyMap := make(map[string]string)
	conf.Analytics.Properties = propertyMap

	err := validateAnalyticsConfigs(conf)
	if err != nil {
		t.Errorf("Expected validation of '%s' to be successful, but got an error: %v", "empty property map", err)
	}

	propertyMap["notDefinedGatewayTypeProperty"] = "dummy_value"
	conf.Analytics.Properties = propertyMap

	err = validateAnalyticsConfigs(conf)
	if err != nil {
		t.Errorf("Expected validation of '%s' to be successful, but got an error: %v", "not defined gateway type property", err)
	}
}

func TestValidateConfigValuesForInvalidValues(t *testing.T) {

	validValues := []string{gatewayTypeDefault + "-test", "choreo", "choreo-pdp"}
	conf := ReadConfigs()

	for _, config := range validValues {
		propertyMap := make(map[string]string)
		propertyMap[gatewayTypeValue] = config
		conf.Analytics.Properties = propertyMap

		err := validateAnalyticsConfigs(conf)
		if err == nil {
			t.Errorf("Expected validation of '%s' to result in an error, but got nil", config)
		}
	}

}
