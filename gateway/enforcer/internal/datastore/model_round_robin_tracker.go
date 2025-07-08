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
	"log"
	"sync"
	"time"
)

// ModelBasedRoundRobinTracker tracks API-Resource-Model round-robin counts
type ModelBasedRoundRobinTracker struct {
	sync.Mutex
	counts    map[string]map[string]map[string]int       // Nested map: API -> Resource -> Model -> Count
	suspended map[string]map[string]map[string]time.Time // Nested map: API -> Resource -> Model -> Suspension End Time
}

var (
	trackerInstance *ModelBasedRoundRobinTracker
	once            sync.Once
)

// GetModelBasedRoundRobinTracker returns the singleton instance of ModelBasedRoundRobinTracker
func GetModelBasedRoundRobinTracker() *ModelBasedRoundRobinTracker {
	once.Do(func() {
		trackerInstance = &ModelBasedRoundRobinTracker{
			counts:    make(map[string]map[string]map[string]int),
			suspended: make(map[string]map[string]map[string]time.Time),
		}
		go trackerInstance.ReactivateSuspendedModels()
	})
	return trackerInstance
}

// ModelWeight represents a model with its associated weight
type ModelWeight struct {
	Name     string
	Endpoint string
	Weight   int
}

// GetNextModel handles weighted round-robin logic and returns the next model for the given API and resource
func (r *ModelBasedRoundRobinTracker) GetNextModel(api, resource string, models []ModelWeight) (string, string) {
	r.Lock()
	defer r.Unlock()

	// Initialize maps if not already done
	if _, exists := r.counts[api]; !exists {
		r.counts[api] = make(map[string]map[string]int)
	}
	if _, exists := r.suspended[api]; !exists {
		r.suspended[api] = make(map[string]map[string]time.Time)
	}
	if _, exists := r.counts[api][resource]; !exists {
		r.counts[api][resource] = make(map[string]int)
	}
	if _, exists := r.suspended[api][resource]; !exists {
		r.suspended[api][resource] = make(map[string]time.Time)
	}

	// Filter out suspended models
	activeModels := []ModelWeight{}
	for _, model := range models {
		if suspendEnd, suspended := r.suspended[api][resource][model.Name]; !suspended || time.Now().After(suspendEnd) {
			activeModels = append(activeModels, model)
		}
	}
	log.Println("Suspended Models: ", r.suspended)
	log.Println("Active Models: ", activeModels)

	// If no active models are available, return an empty string
	if len(activeModels) == 0 {
		return "", ""
	}

	// Perform weighted round-robin on active models
	totalWeight := 0
	var selectedModel string
	var selectedEndpoint string
	minEffectiveWeight := int(^uint(0) >> 1) // Initialize with max int value

	for _, model := range activeModels {
		totalWeight += model.Weight
		effectiveWeight := r.counts[api][resource][model.Name] / model.Weight
		if effectiveWeight < minEffectiveWeight {
			selectedModel = model.Name
			selectedEndpoint = model.Endpoint
			minEffectiveWeight = effectiveWeight
		}
	}

	// Increment the count for the selected model
	r.counts[api][resource][selectedModel]++
	return selectedModel, selectedEndpoint
}

// SuspendModel suspends a model for the given API and resource
func (r *ModelBasedRoundRobinTracker) SuspendModel(api, resource, model string, duration time.Duration) {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.suspended[api]; !exists {
		r.suspended[api] = make(map[string]map[string]time.Time)
	}
	if _, exists := r.suspended[api][resource]; !exists {
		r.suspended[api][resource] = make(map[string]time.Time)
	}

	// Set the suspension end time
	log.Println("Suspended API-Resource-Model: ", api, resource, model)
	log.Println("Suspended Duration: ", duration)
	r.suspended[api][resource][model] = time.Now().Add(duration)
}

// ReactivateSuspendedModels periodically checks and removes expired suspensions
func (r *ModelBasedRoundRobinTracker) ReactivateSuspendedModels() {
	for {
		time.Sleep(1 * time.Second) // Periodic task

		r.Lock()
		for api, resources := range r.suspended {
			for resource, models := range resources {
				for model, suspendEnd := range models {
					if time.Now().After(suspendEnd) {
						delete(r.suspended[api][resource], model)
					}
				}
			}
		}
		r.Unlock()
	}
}
