package types

import (
	"database/sql"
)

type ISQL interface {
	Schema(string) ([]byte, error)
	Execute(string) ([]byte, error)
	Tables(string) ([]byte, error)
}

type Table struct {
	Name        string   `json:"name"`
	Columns     []Column `json:"columns"`
	ColumnCount int64    `json:"column_count"`
	Description string   `json:"description"`
	Metatags    []string `json:"metatags"` // Add all Column name in metatags
}

type Column struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	IsNullable   string         `json:"is_nullable"`
	Key          string         `json:"key"`
	DefaultValue sql.NullString `json:"default_value"`
	Extra        string         `json:"extra"` // Add more description about the field extra
	Description  string         `json:"description"`
	Metatags     []string       `json:"metatags"` // Add Column name in metatags, ["name"]
	Visibility   bool           `json:"visibility"`
	// TODO: Add more datapoints like (Not P0)
	// isIndex, IsPrimary, Foreign Key,
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Time    int64           `json:"time"`
	Error   string          `json:"error"`
}

type DbType int

const (
	MySQL DbType = iota + 1
	Postgres
)

func (w DbType) String() string {
	return [...]string{"mysql", "postgres"}[w-1]
}

func (w DbType) Index() int {
	return int(w)
}
