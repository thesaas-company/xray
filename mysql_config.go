package library

import (
	"database/sql"
	"os"
	"github.com/adarsh-jaiss/library/sample/sample"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
)

type MySQL struct {
	Client *sql.DB
}

func NewMySQL(dbConfig *sample.DatabaseConfig) (ISQL, error) {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }


	if os.Getenv("DB_PASSWORD") == "" || len(os.Getenv("DB_PASSWORD")) == 0 {
		return nil, fmt.Errorf("please set DB_PASSWORD env variable for the database")
	}
	dsn := dbURLMySQL(dbConfig)
	
	db, err := sql.Open(dbConfig.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}
	return &MySQL{
		Client: db,
	},nil
	
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