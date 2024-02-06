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
 */

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
)

var dbPool *pgxpool.Pool

func ConnectToDB() {
	conf := config.ReadConfigs()
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&"+
		"pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s&pool_health_check_period=%s&pool_max_conn_lifetime_jitter=%s",
		conf.CommonController.Database.Username,
		conf.CommonController.Database.Password,
		conf.CommonController.Database.Host,
		conf.CommonController.Database.Port,
		conf.CommonController.Database.Name,
		conf.CommonController.Database.PoolOptions.PoolMaxConns,
		conf.CommonController.Database.PoolOptions.PoolMinConns,
		conf.CommonController.Database.PoolOptions.PoolMaxConnLifetime,
		conf.CommonController.Database.PoolOptions.PoolMaxConnIdleTime,
		conf.CommonController.Database.PoolOptions.PoolHealthCheckPeriod,
		conf.CommonController.Database.PoolOptions.PoolMaxConnLifetimeJitter)
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		loggers.LoggerDatabase.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Unable to connect to database: %v", err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1100,
		})
	}
}

func ExecDBQuery(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := dbPool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func IsAliveConn(ctx context.Context) (isAlive bool) {
	if err := dbPool.Ping(ctx); err != nil {
		return true
	}
	return isAlive
}

func CloseDBConn() {
	dbPool.Close()
}
