package library

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	_ "github.com/go-sql-driver/mysql"

	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/joho/godotenv"
)

type TestMySql struct {
	Client *sql.DB
}

func NewTestMySQL() (ISQL, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DatabaseConfig := &sample.DatabaseConfig{
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		DatabaseName: os.Getenv("DB_NAME"),
		SSL:          os.Getenv("DB_SSL"),
		DBType:       os.Getenv("DB_TYPE"),
	}

    dsn := dbURLMySQL(DatabaseConfig)
    db, err := sql.Open(DatabaseConfig.DBType, dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening connection to database: %v", err)
    }

    return &MySQL{
        Client: db,
    }, nil
}

func TestConnection(t *testing.T) {
	tClient,err := NewTestMySQL()
	if err!= nil{
		t.Errorf("Error connecting client, Expected No error, got: %v", err)
	}
	if tClient == nil {
		t.Errorf("Expected a client, got nil")
	}

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
	os.Setenv("DB_TYPE", "mysql")

	DBConfig := &sample.DatabaseConfig{
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		DatabaseName: os.Getenv("DB_NAME"),
		SSL:          os.Getenv("DB_SSL"),
		DBType:       os.Getenv("DB_TYPE"),
	}

    m := MySQL{}

    testCases := []struct {
        name        string
        dbType      string
        expectError bool
    }{
        {
            name:        "Valid DB Type",
            dbType:      DBConfig.DBType,
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
		client, err := m.NewClient(DBConfig ,tc.dbType)

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
                _, ok := client.(ISQL)
                if !ok {
                	t.Errorf("Expected a client implementing ISQL, got %T", client)
                }
            }
        })
    }
}
