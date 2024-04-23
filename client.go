package library

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/library/sample/config"
	"github.com/adarsh-jaiss/library/sample/databases/mysql"
	"github.com/adarsh-jaiss/library/sample/databases/postgres"
	"github.com/adarsh-jaiss/library/sample/logger"
	"github.com/adarsh-jaiss/library/sample/types"
)

//TODO: Add description

// TODO: Add description
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
