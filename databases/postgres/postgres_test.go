package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	"github.com/joho/godotenv"
)

type TestPostgres struct {
	Client *sql.DB
}

func NewTestPostgres() (types.ISQL, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConfig := &sample.DatabaseConfig{
		Username:     os.Getenv("POSTGRES_DB_USERNAME"),
		Password:     os.Getenv("POSTGRES_DB_PASSWORD"),
		Host:         os.Getenv("POSTGRES_DB_HOST"),
		DatabaseName: os.Getenv("POSTGRES_DB_NAME"),
		SSL:          os.Getenv("POSTGRES_DB_SSLMODE"),
		DBType:       os.Getenv("POSTGRES_DB_TYPE"),
		Port:         os.Getenv("POSTGRES_DB_PORT"),
	}
	
	db,err := sql.Open("postgres", fmt.Sprintf( "user=%s password=%s dbname=%s sslmode=%s",dbConfig.Username,dbConfig.Password,dbConfig.DatabaseName,dbConfig.SSL))
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}

	return &Postgres{
		Client: db,
	}, nil
}

func TestPostgresConnection(t *testing.T) {
	tClient, err := NewTestPostgres()
	if err != nil {
		t.Errorf("Error connecting client, Expected No error, got: %v", err)
	}
	if tClient == nil {
		t.Errorf("Expected a client, got nil")
	}

}

func TestPostgresSchema(t *testing.T) {
	tClient, err := NewTestPostgres()
	if err != nil {
		t.Errorf("Error connecting client, Expected No error, got: %v", err)
	}

	tableName := os.Getenv("POSTGRES_TABLE_NAME")
	res, err := tClient.Schema(tableName)
	if err != nil {
		t.Errorf("Error retrieving schema, Expected No error, got: %v", err)
	}

	if res == nil {
		t.Errorf("Expected schema, got nil")
	}

	t.Logf("schema: %v", res)
}

func TestPostgresExecute(t *testing.T) {
	tClient, err := NewTestPostgres()
	if err != nil {
		t.Errorf("Error connecting client, Expected No error, got: %v", err)
	}
	table := "user"
	query := "SELECT * FROM " + table
	res, err := tClient.Execute(query)

	if err != nil {
		t.Errorf("Error executing query, Expected No error, got: %v", err)
	}

	t.Logf("Execute result: %v", res)
}

func TestPostgresGetTables(t *testing.T) {

	tClient, err := NewTestPostgres()
	if err != nil {
		t.Errorf("Error creating client, Expected No error, got: %v", err)

	}
	DBName := os.Getenv("POSTGRES_DB_NAME")

	
	tables, err := tClient.Tables(DBName)
	if err != nil {
		t.Errorf("Error getting tables, Expected No error, got: %v", err)
	}

	t.Logf("Tables List: %v", tables)
}

func TestPostgresNewclient(t *testing.T) {
	os.Setenv("POSTGRES_DB_TYPE", "postgres")

	DBConfig := &sample.DatabaseConfig{
		Username:     os.Getenv("POSTGRES_DB_USERNAME"),
		Password:     os.Getenv("POSTGRES_DB_PASSWORD"),
		Host:         os.Getenv("POSTGRES_DB_HOST"),
		DatabaseName: os.Getenv("POSTGRES_DB_NAME"),
		SSL:          os.Getenv("POSTGRES_DB_SSL"),
		DBType:       os.Getenv("POSTGRES_DB_TYPE"),
	}

	p := &Postgres{}

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
			client, err := p.NewClient(DBConfig, tc.dbType)

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
