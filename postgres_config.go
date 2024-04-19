package library

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	"github.com/joho/godotenv"
)

type Postgres struct {
	Client *sql.DB
}



func PostgresDBURL(dbConfig *sample.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.DatabaseName,
		dbConfig.SSL,
	)
}

func NewPostgres(dbConfig *sample.DatabaseConfig) (types.ISQL, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("POSTGRES_DB_PASSWORD") == "" || len(os.Getenv("POSTGRES_DB_PASSWORD")) == 0 {
		return nil, fmt.Errorf("please set POSTGRES_DB_PASSWORD env varibale for the database")
	}

	dsn := PostgresDBURL(dbConfig)

	db, err := sql.Open(dbConfig.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}

	return &Postgres{
		Client: db,
	}, nil

}

