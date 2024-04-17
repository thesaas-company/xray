package sample

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type TableContext struct {
	Name string       `json:"name"`
	Data []ColContext `json:"data"`
}

type ColContext struct {
	ColumnName string `json:"column_name"`
	ColumnKey  string `json:"column_key"`
	DataType   string `json:"data_type"`
}

type ISQL interface {
	Schema(tables []string) ([]TableContext, int64, error)
	Execute(query string) ([]map[string]interface{}, error)
	Tables(database string) ([]string, error)
}

func NewSQL(cfg *Config) (ISQL, error) {
	switch cfg.Database.DBType {
	case "mysql":
		return NewMysqlRepo(cfg)
	case "postgres":
		return NewPostgresRepo(cfg)
	// case "redshift":
	// 	return NewRedshiftRepo(cfg)
	case "snowflake":
		return NewSnowflakeRepo(cfg)
	case "bigquery":
		return NewBigQueryRepo(cfg)
	}
	return nil, nil
}

func rowsToJSON(rows *sql.Rows) (string, int64) {
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice to hold the results
	var results []map[string]interface{}
	colCount := 0
	// Iterate through the rows and populate the results
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			colCount++
			switch v := values[i].(type) {
			case []byte:
				// If the value is a byte slice, convert it to a string or handle it accordingly.
				rowData[colName] = string(v) // Example: Convert to string
			default:
				rowData[colName] = v
			}
		}

		results = append(results, rowData)
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData), int64(colCount)
}
