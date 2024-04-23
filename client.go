package library

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/xray/config"
	"github.com/adarsh-jaiss/xray/databases/mysql"
	"github.com/adarsh-jaiss/xray/databases/postgres"
	"github.com/adarsh-jaiss/xray/logger"
	"github.com/adarsh-jaiss/xray/types"
)

// NewClientWithConfig creates a new SQL client with the given configuration and database type.
// It returns an error if the database type is not supported or if there is a problem creating the client.
func NewClientWithConfig(dbConfig *config.Config, dbType types.DbType) (types.ISQL, error) {
	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQLWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgresWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// NewClient creates a new SQL client with the given database client and database type.
// It returns an error if the database type is not supported or if there is a problem creating the client.
func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {

	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQL(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgres(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
