package library

import (
	"database/sql"

	"github.com/adarsh-jaiss/library/sample/sample"
)

type Table struct {
	Name        string          `json:"name"`
	Data        []Column `json:"data"`
	ColumnCount int64           `json:"column_count"`
	Description string          `json:"description"`
	Metatags    string          `json:"metatags"`
}

type Column struct {
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


type DbType int

const (
	MySQL    DbType = iota + 1 
	Postgres  
)

func (w DbType) String() string {
	return [...]string{"mysql", "postgres"}[w-1]
}

func (w DbType) EnumIndex() int {
	return int(w)
}