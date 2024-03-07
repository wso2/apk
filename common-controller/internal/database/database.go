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
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
)

var dbPool *pgxpool.Pool

// ConnectToDB connects to the database
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

// ExecDBQuery executes a database query
func ExecDBQuery(tx pgx.Tx, query string, args ...interface{}) error {
	_, err := tx.Exec(context.Background(), query, args...)
	return err
}

// ExecDBQueryRows executes a database query and returns a row
func ExecDBQueryRows(tx pgx.Tx, query string, args ...interface{}) (pgx.Rows, error) {
	return tx.Query(context.Background(), query, args...)
}

// IsAliveConn checks if the database connection is alive
func IsAliveConn(ctx context.Context) (isAlive bool) {
	err := dbPool.Ping(ctx)
	return err == nil
}

// CloseDBConn closes the database connection
func CloseDBConn() {
	conf := config.ReadConfigs()
	if conf.CommonController.Database.Enabled {
		dbPool.Close()
	}
}

// PrepareQueries prepares the queries
func PrepareQueries(tx pgx.Tx, queries ...string) {
	for _, query := range queries {
		_, err := tx.Prepare(context.Background(), query, query)
		if err != nil {
			loggers.LoggerAPI.Errorf("Error while preparing query: %s, %s", query, err.Error())
		}
	}
}

// performTransaction performs a transaction
func performTransaction(fn func(tx pgx.Tx) error) error {
	con := context.Background()
	tx, err := dbPool.BeginTx(con, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error while begining the transaction %v", err)
	}
	defer func() {
		if err != nil {
			loggers.LoggerAPI.Error("Rollback due to error: ", err)
			err = tx.Rollback(con)
		} else {
			err = tx.Commit(con)
		}
		if err != nil {
			loggers.LoggerAPI.Error("Error while commiting the transaction ", err)
		}
	}()
	err = fn(tx)
	return err
}

// retryTransaction retries a transaction
func retryUntilTransaction(fn func(tx pgx.Tx) error) error {
	if err := performTransaction(fn); err != nil {
		loggers.LoggerAPI.Warn("Retrying because of the error: ", err)
		if strings.Contains(err.Error(), "conn closed") {
			loggers.LoggerAPI.Info("Reconnecting to DB...")
			ConnectToDB()
		}
		return performTransaction(fn)
	}
	return nil
}
