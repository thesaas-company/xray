package library

import (
	"database/sql"
	"fmt"
	"github.com/adarsh-jaiss/library/sample/sample"
	"testing"
)

type TestMySql struct {
	Client *sql.DB
}

var DatabaseConfig = &sample.DatabaseConfig {
	Username:     "root",
	Password:     "root",
	Host:         "localhost:3306",
	DatabaseName: "test",
	SSL:          "false",
	DBType:       "mysql",
}

func NewTestMySQL() (ISQL, error) {
	dsn := dbURLMySQL(DatabaseConfig)
	db, err := sql.Open(DatabaseConfig.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}

	return &MySQL{
		Client: db,
	}, nil
}

func TestSchema(t *testing.T) {
	tClient,err := NewTestMySQL()
	if err!= nil{
		t.Errorf("Error creating client, Expected No error, got: %v", err)
	}

	res, err := tClient.Schema("user")
	if err != nil {
		t.Errorf("Error getting schema, Expected No error, got: %v", err)
	}

	// fmt.Println(res)
	t.Logf("schema: %v", res)
}

func TestExecute(t *testing.T) {
	tClient,err := NewTestMySQL()
	if err!= nil{
		t.Errorf("Error creating client, Expected No error, got: %v", err)
	}

	query := "SELECT * FROM user"
	res, err := tClient.Execute(query)

	if err != nil {
		t.Errorf("Error executing query, Expected No error, got: %v", err)
	}

	t.Logf("Execute result: %v", res)
}

func TestGetTables(t *testing.T) {
	tClient,err := NewTestMySQL()
	if err!= nil{
		t.Errorf("Error creating client, Expected No error, got: %v", err)
	
	}
	DBName := "test"
	tables, err := tClient.Tables(DBName)
	if err != nil {
		t.Errorf("Error getting tables, Expected No error, got: %v", err)
	}

	t.Logf("Tables List: %v", tables)
}

func TestNewClient(t *testing.T) {
	m := &MySQL{}

	client, err := m.NewClient(DatabaseConfig, "mysql")
	if err != nil {
		t.Errorf("Error creating new client, Expected No error, got: %v", err)
	}
	// _, ok := client.(ISQL)
	// if !ok {
	// 	t.Errorf("Expected a client implementing ISQL, got %T", client)
	// }

	client, err = m.NewClient(DatabaseConfig, "unsupported")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if client != nil {
		t.Errorf("Expected nil, got %v", client)
	}
}
