package library

import (
	"encoding/json"
	"fmt"

	"github.com/adarsh-jaiss/library/sample/sample"
	"github.com/adarsh-jaiss/library/sample/types"
	_ "github.com/go-sql-driver/mysql"
)

// This method will accept a table name as input and return the table schema (structure).
func (m *MySQL) Schema(table string) ([]byte, error) {
	// prepare the sql statement
	// This is important to avoid overhead of parsing and compiling the SQL command each time it's executed.
	statement, err := m.Client.Prepare("DESCRIBE " + table)
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
	var columns []types.ColumnContext
	for rows.Next() {
		var column types.ColumnContext
		if err := rows.Scan(&column.ColumnName, &column.DataType, &column.IsNullable, &column.ColumnKey, &column.DefaultValue, &column.Extra); err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""  // default description
		column.Metatags = ""     // default metatags
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	tableContext := types.TableContext{
		Name:        table,
		Data:        columns,
		ColumnCount: int64(len(columns)),
		Description: "", 
		Metatags:    "", 
		
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
	queryResult := types.QueryResult{
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
func (m *MySQL) Tables(databaseName string) ([]byte, error) {
	statememt, err := m.Client.Prepare("SELECT table_name FROM information_schema.tables WHERE table_schema = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing sql statement: %v", err)
	}
	defer statememt.Close()

	// execute the sql statement
	rows,err := statememt.Query(databaseName)
	if err!= nil{
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	//scan and append the result
	var tables []string
	for rows.Next(){
		var table string
		if err := rows.Scan(&table); err!= nil{
			return nil, fmt.Errorf("error scanning rows: %v",err)
		}
		tables = append(tables, table)
	}

	// checking for errors in iterating over rows
	if err := rows.Err(); err!= nil{
		return nil, fmt.Errorf("error iterating over rows:%v",err)
	}

	// convert the result to json
	jsonData,err := json.Marshal(tables)
	if err!= nil{
		return nil,fmt.Errorf("error marshalling json: %v",err)
	}

	return jsonData, nil
}

//  Generate an interface based on the specified database type.
func(m *MySQL) NewClient(dbConfig *sample.DatabaseConfig, dbType string) (types.ISQL, error) {
	switch dbType {
	case "mysql":
		return NewMySQL(dbConfig)
	// case "postgres":
	// 	return &Postgres{},nil
	// case "snowflake":
	// 	return &SnowFlake{},nil
	// case "bigquery":
	// 	return &Bigquery{},nil
	// case "redshift":
	// 	return &RedShift{},nil
	default:
			return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
	



