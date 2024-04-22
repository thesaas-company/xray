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
	MySQL_Client,_ := mysql.NewMySQLWithConfig(dbConfig)
	MySQL_Client = middleware.NewLogMiddleware(MySQL_Client)

	Postgres_Client,_ := postgres.NewPostgresWithConfig(dbConfig)
	Postgres_Client = middleware.NewLogMiddleware(Postgres_Client)
	
	
	switch dbType {
	case types.MySQL:
		return MySQL_Client,nil
	case types.Postgres:
		return Postgres_Client,nil
	// case "snowflake":
	// 	return &SnowFlake{},nil
	// case "bigquery":
	// 	return &Bigquery{},nil
	// case "redshift":
	// 	return &RedShift{},nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {
	MySQL_Client,_ := mysql.NewMySQL(dbClient)
	MySQL_Client = middleware.NewLogMiddleware(MySQL_Client)

	Postgres_Client,_ := postgres.NewPostgres(dbClient)
	Postgres_Client = middleware.NewLogMiddleware(Postgres_Client)
	
	
	switch dbType {
	case types.MySQL:
		return MySQL_Client,nil
	case types.Postgres:
		return Postgres_Client,nil
	// case "snowflake":
	// 	return &SnowFlake{},nil
	// case "bigquery":
	// 	return &Bigquery{},nil
	// case "redshift":
	// 	return &RedShift{},nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
