package library

import (
	"database/sql"
	"fmt"
	// "os"
	"github.com/adarsh-jaiss/library/sample/sample"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Client *sql.DB
}

func NewMySQL(db *sql.DB) ISQL{
	return &MySQL{
		Client: db,
	}
}

func dbURL(dbConfig *sample.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.DatabaseName,
		dbConfig.SSL,
	)
}

func NewMysqlConnection(config *sample.Config) (ISQL, error) {
	// if os.Getenv("DB_PASSWORD") == "" || len(os.Getenv("DB_PASSWORD")) == 0 {
	// 	return nil, fmt.Errorf("please set DB_PASSWORD env variable for the database")
	// }

	db,err := sql.Open(config.DBType, dbURL(&config.Database))
	if err!= nil{
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}
	return &MySQL{
		Client: db,
	},nil
}


func (m *MySQL) Schema(table string) ([]byte, error) {
	return nil, nil
}

func (m *MySQL) Execute(query string) ([]byte, error) {
	return nil, nil
}

func (m *MySQL) Tables(database string) ([]byte, error) {
	return nil, nil
}

func (m *MySQL) NewClient(dbType string) (ISQL, error) {
	return nil, nil
}
