package library

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type MySQL struct {
	Client *sql.DB
}

func NewMySQL(dbConfig *sample.DatabaseConfig) (types.ISQL, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("MYSQL_DB_PASSWORD") == "" || len(os.Getenv("MYSQL_DB_PASSWORD")) == 0 { // added mysql to be more verbose about the db type
		return nil, fmt.Errorf("please set MYSQL_DB_PASSWORD env variable for the database")
	}
	dsn := dbURLMySQL(dbConfig)

	db, err := sql.Open(dbConfig.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}

	return &MySQL{
		Client: db,
	}, nil

}

func dbURLMySQL(dbConfig *sample.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.DatabaseName,
		dbConfig.SSL,
	)
}
