package types

import (
	"database/sql"

	"github.com/adarsh-jaiss/library/sample/sample"
)

type TableContext struct {
	Name        string          `json:"name"`
	Data        []ColumnContext `json:"data"`
	ColumnCount int64           `json:"column_count"`
	Description string          `json:"description"`
	Metatags    string          `json:"metatags"`
}

type ColumnContext struct {
    ColumnName   string        `json:"column_name"`
    DataType     string        `json:"data_type"`
    IsNullable   string        `json:"is_nullable"`
    ColumnKey    string        `json:"column_key"`
    DefaultValue sql.NullString `json:"default_value"`
    Extra        string        `json:"extra"`
    Description  string        `json:"description"`
    Metatags     string        `json:"metatags"`
    Visibility   bool          `json:"visibility"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

// type ISQL interface {
// 	Schema(table string) ([]byte, error)
// 	Execute(query string) ([]byte, error)
// 	Tables(databaseName string) ([]byte, error)
// 	NewClient(dbType string) (ISQL, error)
// }

type ISQL interface {
	Schema(string) ([]byte, error)
	Execute(string) ([]byte, error)
	Tables(string) ([]byte, error)
	NewClient(*sample.DatabaseConfig ,string) (ISQL, error)
}