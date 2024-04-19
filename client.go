package library

import (
	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	"fmt"
)

func NewClient(client types.ISQL) func(*sample.DatabaseConfig, string) (types.ISQL, error) {
	return func(dbConfig *sample.DatabaseConfig, dbType string) (types.ISQL, error) {
		switch dbType {
		case "mysql":
			return NewMySQL(dbConfig)
		case "postgres":
			return NewPostgres(dbConfig)
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
}