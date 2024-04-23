package library

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/library/sample/config"
	"github.com/adarsh-jaiss/library/sample/databases/mysql"
	"github.com/adarsh-jaiss/library/sample/databases/postgres"
	"github.com/adarsh-jaiss/library/sample/middleware"
	"github.com/adarsh-jaiss/library/sample/types"
)

//TODO: Add description

// TODO: Add description
func NewClientWithConfig(dbConfig *config.Config, dbType types.DbType) (types.ISQL, error) {	
	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQLWithConfig(dbConfig)
		// TODO: Handle error
		return  middleware.NewLogMiddleware(sqlClient),nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgresWithConfig(dbConfig)
		// TODO: Handle error
		return middleware.NewLogMiddleware(sqlClient),nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {

	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQL(dbConfig)
		// TODO: Handle error
		return  middleware.NewLogMiddleware(sqlClient),nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgres(dbConfig)
		// TODO: Handle error
		return middleware.NewLogMiddleware(sqlClient),nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
