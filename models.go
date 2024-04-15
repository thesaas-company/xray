package library

type TableContext struct {
	Name        string          `json:"name"`
	Data        []ColumnContext `json:"data"`
	ColumnCount int64           `json:"column_count"`
}

type ColumnContext struct {
	ColumnName  string `json:"column_name"`
	ColumnKey   string `json:"column_key"`
	DataType    string `json:"data_type"`
	// Description string `json:"description"`
	// Metatags    string `json:"metatags"`
	// Visibility  bool   `json:"visibility"`
}

type QueryResult struct {
	Columns []string `json:"columns"`
	Rows   [][]interface{}    `json:"rows"`
}

type ISQL interface {
	Schema(table string) ([]byte, error)
	Execute(query string) ([]byte, error)
	Tables(database string) ([]byte, error)
	NewClient(dbType string) (ISQL, error)
}





