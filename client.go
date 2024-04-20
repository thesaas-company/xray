package library

import (
	"github.com/adarsh-jaiss/library/databases/mysql"
	"github.com/adarsh-jaiss/library/databases/postgres"
	"fmt"
)


//TODO: Add description
// Please fix the interface changes
type ISQL interface {
	Schema(string) (Table, error)
	Execute(string) ([]byte, error)
	Tables(string) ([]string, error)
}

//TODO: Add description
func NewClientWithConfig(dbConfig *config.Config, dbType DbType) (ISQL, error)  {
		switch dbType {
		case MySQL:
			return mysql.NewMySQLWithConfig(dbConfig)
		case Postgres:
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

func NewClient(dbClient *sql.DB, dbType DbType) (ISQL, error)  {
	switch dbType {
	case MySQL:
		return mysql.NewMySQL(dbClient)
	case Postgres:
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