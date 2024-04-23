package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adarsh-jaiss/library/sample/types"
)

// setting up a mock db connection
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred while creating a new mock database connection")
	}

	// Add dummy data for testing
	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John").
		AddRow(2, "Alice").
		AddRow(3, "Bob")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM user`)).WillReturnRows(mockRows)

	return db, mock
}

func TestSchema(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableName := "user"
	query := "DESCRIBE " + tableName

	columns := []string{"Field", "Type", "Null", "Key", "Default", "Extra"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", "NO", "PRI", nil, "auto_increment")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	// we then create a new instance of our MySQL object and test the function
	m := MySQL{Client: db}
	res, err := m.Schema(tableName)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}

	var response types.Table
	if err := json.Unmarshal(res, &response); err != nil {
		t.Errorf("error was not expected while recording stats: %s", err)
	}

	fmt.Printf("Table schema : %+v\n", response)

	// we make sure that all expectations were met, otherwise an error will be reported
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT * FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	m := MySQL{Client: db}
	res, err := m.Execute(query)
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

func TestGetTableName(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableList := []string{"user", "product", "order"}

	// Retrieve the list of tables
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(MYSQL_TABLES_LIST_QUERY).WillReturnRows(rows)

	m := MySQL{Client: db}
	tables, err := m.Tables("test")
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	var res []string
	err = json.Unmarshal(tables, &res)
	if err != nil {
		t.Errorf("error unmarshalling the result : %s", err)
	}

	fmt.Printf("Table names: %+v\n", res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
