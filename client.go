package library

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/library/sample/config"
	"github.com/adarsh-jaiss/library/sample/databases/mysql"
	"github.com/adarsh-jaiss/library/sample/databases/postgres"
	"github.com/adarsh-jaiss/library/sample/types"
)

//TODO: Add description

// TODO: Add description
func NewClientWithConfig(dbConfig *config.Config, dbType types.DbType) (types.ISQL, error) {
	switch dbType {
	case types.MySQL:
		return mysql.NewMySQLWithConfig(dbConfig)
	case types.Postgres:
		return postgres.NewPostgresWithConfig(dbConfig)
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
	switch dbType {
	case types.MySQL:
		return mysql.NewMySQL(dbClient)
	case types.Postgres:
		return postgres.NewPostgres(dbClient)
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
