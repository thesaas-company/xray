package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/adarsh-jaiss/library/sample/config"
	"github.com/adarsh-jaiss/library/sample/types"
)

var DB_PASSWORD = "DB_PASSWORD"

const (
	POSTGRES_SCHEMA_QUERY = "SELECT column_name, data_type, character_maximum_length FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1;"
	POSTGRES_TABLE_LIST_QUERY= "SELECT table_name FROM information_schema.tables WHERE table_schema= $1 AND table_type='BASE TABLE';"
)
type Postgres struct {
	Client *sql.DB
}

func NewPostgres(dbClient *sql.DB) (types.ISQL, error) {
	return &Postgres{
		Client: dbClient,
	}, nil

}

func NewPostgresWithConfig(dbConfig *config.Config) (types.ISQL, error) {
	// TODO: Add check for env variable DB_PASSWORD, same as mysql  --> done!
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 { // added mysql to be more verbose about the db type
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}
	if os.Getenv(DB_PASSWORD) != "" || len(os.Getenv(DB_PASSWORD)) != 0 {
		DB_PASSWORD = os.Getenv(DB_PASSWORD)
	}
	
	
	dbtype := types.Postgres
	db, err := sql.Open(dbtype.String(), fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, DB_PASSWORD, dbConfig.DatabaseName, dbConfig.SSL))
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}
	return &Postgres{
		Client: db,
	}, nil
}


func (p *Postgres) Schema(table string) ([]byte, error) {
	// TODO: Extract More datapoint if possible
	statement, err := p.Client.Prepare(POSTGRES_SCHEMA_QUERY)
	if err != nil {
		return nil, fmt.Errorf("error preparing sql statement: %v", err)
	}

	defer statement.Close()

	// execute the sql statement
	rows, err := statement.Query(table)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer rows.Close()

	// scanning the result into and append it into a varibale
	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(&column.Name, &column.Type, &column.IsNullable, &column.Key, &column.DefaultValue, &column.Extra); err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""  // default description 
		column.Metatags = []string{}     // default metatags as an empty slice
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	tableContext := types.Table{
		Name:        table,
		Columns:        columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}

	// convert the table context to json
	jsonData, err := json.Marshal(tableContext)
	if err != nil {
		return nil, fmt.Errorf("error marshalling table schema: %v", err)
	}

	return jsonData, nil
}

func (p *Postgres) Execute(query string) ([]byte, error) {
	// prepare the sql statement
	statement, err := p.Client.Prepare(query)
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

func (p *Postgres) Tables(databaseName string) ([]byte, error) {
	statememt, err := p.Client.Prepare(POSTGRES_TABLE_LIST_QUERY)
	if err!= nil{
		return nil,fmt.Errorf("error preparing sql statement: %v",err)
	}
	defer statememt.Close()

	rows,err := statememt.Query(databaseName)
	if err!=nil{
		return nil, fmt.Errorf("error executing sql statement: %v",err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next(){
		var table string
		if err := rows.Scan(&table); err!= nil{
			return nil, fmt.Errorf("error scanning database")
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err!=nil{
		return nil, fmt.Errorf("error interating over rows: %v",err)
	}

	jsonData,err := json.Marshal(tables)
	if err != nil{
		return nil, fmt.Errorf("error marshalling json: %v",err)
	}

	return jsonData,nil

}

