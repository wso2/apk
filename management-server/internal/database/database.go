package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wso2/apk/APKManagementServer/internal/config"
	"github.com/wso2/apk/APKManagementServer/internal/logger"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var dbPool *pgxpool.Pool

func ConnectToDB() {
	conf := config.ReadConfigs()
	var err error
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", conf.Database.Username, conf.Database.Password,
		conf.Database.Host, conf.Database.Port, conf.Database.Name)
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.LoggerServer.ErrorC(logging.ErrorDetails{
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
