/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 LLC. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package database

import (
	"errors"
	"sync"
	"time"

	apkmgt "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
)

// CachedApplication is an application with an expiry timestamp
type CachedApplication struct {
	application       *apkmgt.Application
	expireAtTimestamp int64
}

// ApplicationLocalCache is data holder for cached applications
type ApplicationLocalCache struct {
	stop chan struct{}

	wg   sync.WaitGroup
	mu   sync.RWMutex
	apps map[string]CachedApplication
}

var cleanupInterval time.Duration
var ttl time.Duration

func init() {
	conf := config.ReadConfigs()
	cleanupInterval, _ = time.ParseDuration(conf.Database.DbCache.CleanupInterval)
	ttl, _ = time.ParseDuration(conf.Database.DbCache.TTL)
	DbCache = NewApplicationLocalCache(cleanupInterval)
}

// NewApplicationLocalCache creates an application local cache
func NewApplicationLocalCache(cleanupInterval time.Duration) *ApplicationLocalCache {
	lc := &ApplicationLocalCache{
		apps: make(map[string]CachedApplication),
		stop: make(chan struct{}),
	}

	lc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)
	return lc
}

func (lc *ApplicationLocalCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-lc.stop:
			return
		case <-t.C:
			lc.mu.Lock()
			for uid, cu := range lc.apps {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.apps, uid)
				}
			}
			lc.mu.Unlock()
		}
	}
}

func (lc *ApplicationLocalCache) stopCleanup() {
	close(lc.stop)
	lc.wg.Wait()
}

// Update updates the expiry timestamp of an existing application in the cache
func (lc *ApplicationLocalCache) Update(u *apkmgt.Application, expireAtTimestamp int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.apps[u.Uuid] = CachedApplication{
		application:       u,
		expireAtTimestamp: expireAtTimestamp,
	}
	logger.LoggerDatabase.Infof("Cache updated successfully.. cache: %v", lc.apps)
}

// UpdateSubscriptionInApplication updates the subscription of an application in the cache
func (lc *ApplicationLocalCache) UpdateSubscriptionInApplication(appUUID string, s *apkmgt.Subscription) error {
	if app, ok := lc.apps[appUUID]; ok {
		for i, sub := range app.application.Subscriptions {
			if sub.Uuid == s.Uuid {
				lc.apps[appUUID].application.Subscriptions[i] = s
				return nil
			}
		}
		app.expireAtTimestamp = time.Now().Unix() + ttl.Microseconds()
	} else {
		return ErrApplicationNotInCache
	}
	return ErrSubscriptionNotInAppCache
}

// AddSubscriptionForApplication adds a subscription to an application in the cache
func (lc *ApplicationLocalCache) AddSubscriptionForApplication(appUUID string, s *apkmgt.Subscription) error {
	if app, ok := lc.apps[appUUID]; ok {
		app.application.Subscriptions = append(app.application.Subscriptions, s)
		app.expireAtTimestamp = time.Now().Unix() + ttl.Microseconds()
		return nil
	}
	return ErrApplicationNotInCache
}

// DeleteSubscriptionFromApplication deletes a subscription from an application in the cache
func (lc *ApplicationLocalCache) DeleteSubscriptionFromApplication(appUUID, subUUID string) error {
	if app, ok := lc.apps[appUUID]; ok {
		for i, sub := range app.application.Subscriptions {
			if sub.Uuid == subUUID {
				app.application.Subscriptions[i] = app.application.Subscriptions[len(app.application.Subscriptions)-1]
				app.application.Subscriptions = app.application.Subscriptions[:len(app.application.Subscriptions)-1]
				return nil
			}
		}
		app.expireAtTimestamp = time.Now().Unix() + ttl.Microseconds()
	} else {
		return ErrApplicationNotInCache
	}
	return ErrSubscriptionNotInAppCache
}

var (
	// ErrApplicationNotInCache is the error returned when application is not present in the cache
	ErrApplicationNotInCache = errors.New("unable to find application in cache")
	// ErrSubscriptionNotInAppCache is the error returned when subscription is not present in the cache
	ErrSubscriptionNotInAppCache = errors.New("unable to find subscription in application cache")
)

// Read returns an applicaiton found in the in the cache with the given application id
func (lc *ApplicationLocalCache) Read(id string) (*apkmgt.Application, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	cu, ok := lc.apps[id]
	if !ok {
		return &apkmgt.Application{}, ErrApplicationNotInCache
	}

	return cu.application, nil
}

// ReadAll returns all the applications in the cache
func (lc *ApplicationLocalCache) ReadAll() ([]*apkmgt.Application, error) {
	var apps []*apkmgt.Application

	for _, app := range lc.apps {
		apps = append(apps, app.application)
	}
	return apps, nil
}

// Delete deletes an applicaton in the cache using the given id
func (lc *ApplicationLocalCache) Delete(id string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	delete(lc.apps, id)
}
