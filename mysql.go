package library

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "os"
	"github.com/adarsh-jaiss/library/sample/sample"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Client *sql.DB
}

func NewMySQL(db *sql.DB) ISQL {
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

	db, err := sql.Open(config.DBType, dbURL(&config.Database))
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}
	return &MySQL{
		Client: db,
	}, nil
}

// This method will accept a table name as input and return the table schema (structure).
func (m *MySQL) Schema(table string) ([]byte, error) {
	// prepare the sql statement
	statement, err := m.Client.Prepare("DESCRIBE" + table)
	if err != nil {
		return nil, fmt.Errorf("error preparing sql statement: %v", err)
	}

	defer statement.Close()

	// execute the sql statement
	rows, err := statement.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer rows.Close()

	// scanning the result into and append it into a varibale
	var columns []ColumnContext
	for rows.Next() {
		var column ColumnContext
		if err := rows.Scan(&column.ColumnName, &column.ColumnKey, &column.DataType); err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	tableContext := TableContext{
		Name:        table,
		Data:        columns,
		ColumnCount: int64(len(columns)),
	}

	// convert the table context to json
	jsonData, err := json.Marshal(tableContext)
	if err != nil {
		return nil, fmt.Errorf("error marshalling table schema: %v", err)
	}

	return jsonData, nil
}

// Execute a database query and return the result in JSON format
func (m *MySQL) Execute(query string) ([]byte, error) {
	// prepare the sql statement
	statement, err := m.Client.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing sql statement: %v", err)
	}
	defer statement.Close()

	// execute the sql statement
	rows, err := statement.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	// getting the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		results = append(results, values)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Convert the result to JSON
	queryResult := QueryResult{
		Columns: columns,
		Rows:    results,
	}
	jsonData, err := json.Marshal(queryResult)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %v", err)
	}

	return jsonData, nil
}

// Retrieve the names of tables in the specified database.
func (m *MySQL) Tables(database string) ([]byte, error) {
	return nil, nil
}

//  Generate an interface based on the specified database type.
func (m *MySQL) NewClient(dbType string) (ISQL, error) {
	return nil, nil
}