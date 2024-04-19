package library

import (
	"database/sql"
	"fmt"
	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
)

type Postgres struct {
	Client *sql.DB
}

func NewPostgres(dbConfig *sample.DatabaseConfig) (types.ISQL, error) {
	db, err := sql.Open(dbConfig.DBType, fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DatabaseName, dbConfig.SSL))
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}

	return &Postgres{
		Client: db,
	}, nil

}
