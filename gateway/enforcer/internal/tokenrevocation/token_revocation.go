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

package tokenrevocation

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

const (
	jtiTimeStampSeperator = "_##_"
)

// RevokedTokenFetcher fetches revoked tokens from the Redis datastore.
type RevokedTokenFetcher struct {
	user          string
	password      string
	address       string
	tlsConfig     *tls.Config
	jtiDatastore  *datastore.RevokedJTIStore
	cfg           *config.Server
	channelName   string
	retryInterval time.Duration
}

// NewRevokedTokenFetcher creates a new instance of RevokedTokenFetcher.
func NewRevokedTokenFetcher(cfg *config.Server, jtiDatastore *datastore.RevokedJTIStore, tlsConfig *tls.Config) *RevokedTokenFetcher {
	return &RevokedTokenFetcher{
		user:         cfg.RedisUsername,
		password:     cfg.RedisPassword,
		address:      cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		tlsConfig:    tlsConfig,
		jtiDatastore: jtiDatastore,
		cfg:          cfg,
		channelName:  cfg.RevokedTokensRedisChannel,
	}
}

// Start starts the revoked token fetcher.
func (r *RevokedTokenFetcher) Start() {
	go r.fetchRevokedTokens()
	go r.subscribe()
}

// fetchRevokedTokens fetches revoked tokens from the Redis datastore.
func (r *RevokedTokenFetcher) fetchRevokedTokens() {
	client := util.CreateRedisClient(r.address, r.user, r.password, r.tlsConfig)
	var cursor uint64
	for {
		scanResult, cursor, err := client.Scan(context.Background(), cursor, "", 0).Result()
		if err != nil {
			r.cfg.Logger.Error(err, "Error fetching revoked tokens")
			time.Sleep(r.retryInterval)
			r.cfg.Logger.Sugar().Debug("Retrying to fetch revoked tokens")
			r.fetchRevokedTokens()
			return
		}
		for _, key := range scanResult {
			value, err := client.Get(context.Background(), key).Result()
			if err != nil {
				r.cfg.Logger.Error(err, fmt.Sprintf("Error fetching token value for key %s", key))
				continue
			}
			// Process the key and value as needed
			expirationTime, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				r.cfg.Logger.Error(err, fmt.Sprintf("Error parsing expiration time for key %s", key))
				continue
			}
			r.cfg.Logger.Sugar().Debug(fmt.Sprintf("Fetched revoked token: key=%s, expirationTime=%d", key, expirationTime))
			r.jtiDatastore.AddJTI(key, time.Unix(expirationTime, 0))
		}
		if cursor == 0 {
			break
		}
	}
}

// subscribe subscribes to the Redis channel for revoked tokens.
func (r *RevokedTokenFetcher) subscribe() {
	ctx := context.Background()
	client := util.CreateRedisClient(r.address, r.user, r.password, r.tlsConfig)
	subscriber := client.Subscribe(ctx, r.channelName)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			r.cfg.Logger.Error(err, "Error receiving message from the Redis channel")
			client.Close()
			time.Sleep(r.retryInterval)
			r.cfg.Logger.Sugar().Debug("Retrying to subscribe to the Redis channel")
			r.subscribe()
		}
		r.cfg.Logger.Sugar().Debug(fmt.Sprintf("Received message: %s", msg.Payload))
		parts := strings.Split(msg.Payload, jtiTimeStampSeperator)
		if len(parts) != 2 {
			r.cfg.Logger.Error(fmt.Errorf("invalid message format"), "Error splitting message payload")
			continue
		}
		jti := parts[0]
		expirationTime, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			r.cfg.Logger.Error(err, fmt.Sprintf("Error parsing expiration time for JTI %s", jti))
			continue
		}
		r.cfg.Logger.Sugar().Debug(fmt.Sprintf("Received revoked token: jti=%s, expirationTime=%d", jti, expirationTime))
		r.jtiDatastore.AddJTI(jti, time.Unix(expirationTime, 0))
	}

}
