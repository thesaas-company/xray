// This is End to End Testing using real database connection.

package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/adarsh-jaiss/library/sample/config"
	"github.com/adarsh-jaiss/library/sample/types"
	"github.com/joho/godotenv"
)

// TODO: currently you are doing ent2end testing with real database, It is possible to mock the  database client, Not a P0 but unit test should use a mock database client
type TestMySql struct {
	Client *sql.DB
}

const (
	DBType = "mysql"
)

func NewTestMySQL() (types.ISQL, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	DatabaseConfig := &config.Config{
		Username:     os.Getenv("MYSQL_DB_USERNAME"),
		Host:         os.Getenv("MYSQL_DB_HOST"),
		DatabaseName: os.Getenv("MYSQL_DB_NAME"),
		SSL:          os.Getenv("MYSQL_DB_SSL"),
		Port:         os.Getenv("MYSQL_DB_PORT"),
	}

	dsn := dbURLMySQL(DatabaseConfig)
	db, err := sql.Open(DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}

	return &MySQL{
		Client: db,
	}, nil
}

func TestMySQLConnection(t *testing.T) {
	tClient, err := NewTestMySQL()
	if err != nil {
		t.Errorf("Error connecting client, Expected No error, got: %v", err)
	}
	if tClient == nil {
		t.Errorf("Expected a client, got nil")
	}

}

func TestMysqlSchema(t *testing.T) {
	tClient, err := NewTestMySQL()
	if err != nil {
		t.Errorf("Error creating client, Expected No error, got: %v", err)
	}

	tableName := os.Getenv("POSTGRES_TABLE_NAME")
	res, err := tClient.Schema(tableName)
	if err != nil {
		t.Errorf("Error getting schema, Expected No error, got: %v", err)
	}

	var jsonRes types.Table
	err = json.Unmarshal(res, &jsonRes)
	if err != nil {
		t.Errorf("Error unmarshalling schema, Expected No error, got: %v", err)
	}

	// fmt.Println(res)
	t.Logf("schema: %v", jsonRes)
}

func TestMysqlExecute(t *testing.T) {
	tClient, err := NewTestMySQL()
	if err != nil {
		t.Errorf("Error creating client, Expected No error, got: %v", err)
	}

	query := "SELECT * FROM $1"
	res, err := tClient.Execute(query)

	if err != nil {
		t.Errorf("Error executing query, Expected No error, got: %v", err)
	}

	t.Logf("Execute result: %v", res)
}

func TestGetTables(t *testing.T) {
	tClient, err := NewTestMySQL()
	if err != nil {
		t.Errorf("Error creating client, Expected No error, got: %v", err)

	}
	DBName := os.Getenv("MYSQL_DB_NAME")
	tables, err := tClient.Tables(DBName)
	if err != nil {
		t.Errorf("Error getting tables, Expected No error, got: %v", err)
	}

	t.Logf("Tables List: %v", tables)
}

func TestNewClient(t *testing.T) {
	os.Setenv("MYSQL_DB_TYPE", "mysql")

	DBConfig := &config.Config{
		Username:     os.Getenv("MYSQL_DB_USERNAME"),
		Host:         os.Getenv("MYSQL_DB_HOST"),
		DatabaseName: os.Getenv("MYSQL_DB_NAME"),
		SSL:          os.Getenv("MYSQL_DB_SSL"),
	}

	testCases := []struct {
		name        string
		dbType      string
		expectError bool
	}{
		{
			name:        "Valid DB Type",
			dbType:      DBType,
			expectError: false,
		},
		{
			name:        "Invalid DB Type",
			dbType:      "unsupported",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewMySQLWithConfig(DBConfig)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, got nil")
				}
				if client != nil {
					t.Errorf("Expected nil, got %v", client)
				}
			} else {
				if err != nil {
					t.Log(tc.dbType)
					t.Errorf("Error creating new client, Expected No error, got: %v", err)
				}
				_, ok := client.(types.ISQL)
				if !ok {
					t.Errorf("Expected a client implementing ISQL, got %T", client)
				}
			}
		})
	}
}
