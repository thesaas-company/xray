package xray

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/bigquery"

	"github.com/thesaas-company/xray/databases/mysql"
	"github.com/thesaas-company/xray/databases/postgres"
	"github.com/thesaas-company/xray/databases/redshift"
	"github.com/thesaas-company/xray/databases/snowflake"
	"github.com/thesaas-company/xray/logger"
	"github.com/thesaas-company/xray/types"
)

type clientCreator func(dbType types.DbType) (types.ISQL, error)

func newClientGeneric(dbType types.DbType, creator clientCreator) (types.ISQL, error) {
	switch dbType {
	case types.MySQL, types.Postgres, types.Snowflake, types.BigQuery, types.Redshift:
		sqlClient, err := creator(dbType)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func NewClientWithConfig(dbConfig *config.Config, dbType types.DbType) (types.ISQL, error) {
	return newClientGeneric(dbType, func(dbType types.DbType) (types.ISQL, error) {
		switch dbType {
		case types.MySQL:
			return mysql.NewMySQLWithConfig(dbConfig)
		case types.Postgres:
			return postgres.NewPostgresWithConfig(dbConfig)
		case types.Snowflake:
			return snowflake.NewSnowflakeWithConfig(dbConfig)
		case types.BigQuery:
			return bigquery.NewBigQueryWithConfig(dbConfig)
		case types.Redshift:
			return redshift.NewRedshiftWithConfig(dbConfig)
		default:
			return nil, fmt.Errorf("unsupported database type: %s", dbType)
		}
	})
}

func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {
	return newClientGeneric(dbType, func(dbType types.DbType) (types.ISQL, error) {
		switch dbType {
		case types.MySQL:
			return mysql.NewMySQL(dbClient)
		case types.Postgres:
			return postgres.NewPostgres(dbClient)
		case types.Snowflake:
			return snowflake.NewSnowflake(dbClient)
		case types.BigQuery:
			return bigquery.NewBigQuery(dbClient)
		case types.Redshift:
			return redshift.NewRedshift(dbClient)
		default:
			return nil, fmt.Errorf("unsupported database type: %s", dbType)
		}
	})
}
