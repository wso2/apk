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

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
)

var dbPool *pgxpool.Pool

// ConnectToDB creates the DB connection
func ConnectToDB() {
	conf := config.ReadConfigs()
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&"+
		"pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s&pool_health_check_period=%s&pool_max_conn_lifetime_jitter=%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
		conf.Database.PoolOptions.PoolMaxConns,
		conf.Database.PoolOptions.PoolMinConns,
		conf.Database.PoolOptions.PoolMaxConnLifetime,
		conf.Database.PoolOptions.PoolMaxConnIdleTime,
		conf.Database.PoolOptions.PoolHealthCheckPeriod,
		conf.Database.PoolOptions.PoolMaxConnLifetimeJitter)
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Unable to connect to database: %v", err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1100,
		})
	}
}

// ExecDBQuery executes a given database query with the arguments provided
func ExecDBQuery(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := dbPool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// IsAliveConn checks whether the DB connections pool is alive
func IsAliveConn(ctx context.Context) (isAlive bool) {
	if err := dbPool.Ping(ctx); err != nil {
		return true
	}
	return isAlive
}

// CloseDBConn closes the DB connections pool
func CloseDBConn() {
	dbPool.Close()
}
