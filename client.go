package library

import (
	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	"github.com/adarsh-jaiss/library/databases/mysql"
	"github.com/adarsh-jaiss/library/databases/postgres"
	"fmt"
)

type DbType int

const (
	MySQL    DbType = iota + 1 
	Postgres  
)

func (w DbType) String() string {
	return [...]string{"mysql", "postgres"}[w-1]
}

func (w DbType) EnumIndex() int {
	return int(w)
}


//TODO: Add code docs
type Config struct {
	Host         string `yaml:"host" pflag:",Database host url"`
	Username     string `yaml:"username" pflag:",Database username"`
	Password     string `yaml:"password" pflag:",Database password"`
	DatabaseName string `yaml:"database" pflag:",Database name"`
	Port         string `yaml:"port" pflag:",Database Port"`
	SSL          string `yaml:"ssl" pflag:",Database ssl enable/disable"`
	ProjectID    string `yaml:"project_id" pflag:",BigQuery project ID"`
	JSONKeyPath  string `yaml:"json_key_path" pflag:",Path to BigQuery JSON key file"`
	DBType       string `yaml:"type" pflag:",Database type"`
	Warehouse string         `yaml:"warehouse" pflag:",Snowflake warehouse"`
	Schema    string         `yaml:"schema" pflag:",Snowflake database schema"`
	Account   string         `yaml:"account" pflag:",Snowflake account ID"`
	Debug     bool           `yaml:"debug" pflag:",Debug mode"`
}

//TODO: Add description
type ISQL interface {
	Schema(string) ([]byte, error)
	Execute(string) ([]byte, error)
	Tables(string) ([]byte, error)
}

//TODO: Add description
func NewClient(dbConfig *Config, dbType DbType) (ISQL, error)  {
		switch dbType {
		case MySQL:
			return mysql.NewMySQL(dbConfig)
		case Postgres:
			return postgres.NewPostgres(dbConfig)
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