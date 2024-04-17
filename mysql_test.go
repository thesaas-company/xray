package library


import (
	"database/sql"
	"testing"

	"github.com/adarsh-jaiss/library/sample/sample"
)

type TestMySql struct {
	ISQL
}

func NewTestMySql(db *sql.DB) ISQL {
	return &TestMySql{
		ISQL: &MySQL{
			Client: db,
		},
	}
}

var DatabaseConfig = &sample.DatabaseConfig{
	Username:     "root",
	Password:     "root",
	Host:         "localhost:3306",
	DatabaseName: "test",
	SSL:          "false",
	DBType: "mysql",
}

func Connect(t *testing.T) *TestMySql {
	

	client, err := NewMySQL(DatabaseConfig)
	if err!= nil {
		t.Errorf("Error connecting to MySQL, Expected No error, got: %v", err)
	}
	
	if client == nil {
		t.Errorf("Error connecting to MySQL. Expected a client, got nil.")
	}
	return client.(*TestMySql)
}

func TestSchema(t *testing.T) {
	tClient := Connect(t)
	
	res, err := tClient.Schema("test")
	if err != nil {
		t.Errorf("Error getting schema, Expected No error, got: %v", err)
	}

	// fmt.Println(res)
	t.Logf("schema: %v", res)
}

func TestExecute(t *testing.T) {
	tClient := Connect(t)
	query := "SELECT * FROM test WHERE id = 1"
	res, err := tClient.Execute(query)

	if err != nil {
		t.Errorf("Error executing query, Expected No error, got: %v", err)
	}

	t.Logf("Execute result: %v", res)
}

func TestGetTables(t *testing.T)  {
	tClient := Connect(t)
	DBName := "test"
	tables,err := tClient.Tables(DBName)
	if err != nil {
		t.Errorf("Error getting tables, Expected No error, got: %v", err)
	}

	t.Logf("Tables List: %v", tables)
}

func TestNewClient(t *testing.T){
	m := &MySQL{}

	client, err := m.NewClient(DatabaseConfig,"mysql")
	if err != nil {
		t.Errorf("Error creating new client, Expected No error, got: %v", err)
	}
	_, ok := client.(ISQL)
	if !ok {
		t.Errorf("Expected a client implementing ISQL, got %T", client)
	}

	client, err = m.NewClient(DatabaseConfig,"unsupported")
    if err == nil {
        t.Errorf("Expected an error, got nil")
    }
    if client != nil {
        t.Errorf("Expected nil, got %v", client)
    }

}