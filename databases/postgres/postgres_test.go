// This is unit testing of postgres using a mock DB

package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adarsh-jaiss/library/sample/types"
)

func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred while creating a new mock database connection")
	}

	return db, mock
}

func TestSchema(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()
	POSTGRES_SCHEMA_QUERY := "SELECT column_name, data_type, character_maximum_length FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1;"

	table_name := "user"

	columns := []string{"name", "type", "IsNullable", "key", "Description", "Extra"}

	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", "No", "PRIMARY", "This is the primary key of the table to identify users", "auto_increment")
	mock.ExpectPrepare(regexp.QuoteMeta(POSTGRES_SCHEMA_QUERY))
	mock.ExpectQuery(regexp.QuoteMeta(POSTGRES_SCHEMA_QUERY)).WithArgs(table_name).WillReturnRows(mockRows)

	p := Postgres{Client: db}
	res, err := p.Schema(table_name)
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	var response types.Table
	if err := json.Unmarshal(res, &response); err != nil {
		t.Errorf("error was not expected while recording stats: %s", err)
	}

	fmt.Printf("Table schema: %+v\n", response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT * FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan")

	mock.ExpectPrepare(regexp.QuoteMeta(query))
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	p := Postgres{Client: db}
	res, err := p.Execute(query)
	if err != nil {
		t.Errorf("error executing the query: %s", err)
	}

	var result types.QueryResult
	if err := json.Unmarshal(res, &result); err != nil {
		t.Errorf("error unmarshalling the result: %s", err)
	}

	fmt.Printf("Query result: %+v\n", result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


// TODO : NEED HELP!!!!!
func TestGetTableName(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableList := []string{"user", "Credit", "Debit"}
	TableName := "test"
	POSTGRES_TABLE_LIST_QUERY := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1 AND table_type = 'BASE TABLE';"
	rows := sqlmock.NewRows([]string{TableName}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectPrepare(regexp.QuoteMeta(POSTGRES_TABLE_LIST_QUERY))
	db.Prepare(POSTGRES_TABLE_LIST_QUERY)
	mock.ExpectQuery(regexp.QuoteMeta(POSTGRES_TABLE_LIST_QUERY)).WithArgs(TableName).WillReturnRows(rows)

	m := Postgres{Client: db}
	tables, err := m.Tables(TableName)
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	var res []string
	if err := json.Unmarshal(tables, &res); err != nil {
		t.Errorf("error unmarshalling the result: %s", err)
	}

	expected := []string{"user", "Credit", "Debit"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected: %v, got: %v", expected, res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}