/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package datastore

import (
	"sync"

	config_from_adapter "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/config/enforcer"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
)

// ConfigStore is a thread-safe store for APIs.
type ConfigStore struct {
	configs []*config.EnforcerConfig
	mu      sync.RWMutex
}

// NewConfigStore creates a new instance of ConfigStore.
func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		configs: make([]*config.EnforcerConfig, 0),
	}
}

// AddConfigs adds a list of config to the store.
// This method is thread-safe.
func (s *ConfigStore) AddConfigs(configs []*config_from_adapter.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configs = parseConfig(configs)
}

// GetConfigs retrieves the list of Config from the store.
// This method is thread-safe.
func (s *ConfigStore) GetConfigs() []*config.EnforcerConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.configs
}

func parseConfig(conf []*config_from_adapter.Config) []*config.EnforcerConfig {
	enforcerConfigs := make([]*config.EnforcerConfig, 0)
	for _, c := range conf {
		enforcerConfigs = append(enforcerConfigs, &config.EnforcerConfig{
			Analytics: &config_from_adapter.Analytics{
				Enabled: c.Analytics.Enabled,
			},
		})
	}
	return enforcerConfigs
}
