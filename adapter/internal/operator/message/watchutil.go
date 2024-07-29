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

package message

import (
	"github.com/telepresenceio/watchable"
	"github.com/wso2/apk/adapter/internal/loggers"
)

type Update[K comparable, V any] watchable.Update[K, V]

type Metadata struct {
	Runner  string
	Message string
}

// HandleSubscription takes a channel returned by
// watchable.Map.Subscribe() (or .SubscribeSubset()), and calls the
// given function for each initial value in the map, and for any
// updates.
//
// This is better than simply iterating over snapshot.Updates because
// it handles the case where the watchable.Map already contains
// entries before .Subscribe is called.
func HandleSubscription[K comparable, V any](
	meta Metadata,
	subscription <-chan watchable.Snapshot[K, V],
	handle func(updateFunc Update[K, V], errChans chan error),
) {
	//TODO: find a suitable value
	errChans := make(chan error, 10)
	go func() {
		for err := range errChans {
			loggers.LoggerAPKOperator.Error(err, "observed an error")
		}
	}()

	if snapshot, ok := <-subscription; ok {
		for k, v := range snapshot.State {
			handle(Update[K, V]{
				Key:   k,
				Value: v,
			}, errChans)
		}
	}
	for snapshot := range subscription {
		for _, update := range snapshot.Updates {
			handle(Update[K, V](update), errChans)
		}
	}
}
