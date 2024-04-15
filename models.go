package library

type TableContext struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MetaTags    string          `json:"meta_tags"`
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

type ISQL interface {
	Schema(table string) ([]byte, error)
	Execute(query string) ([]byte, error)
	Tables(database string) ([]byte, error)
	NewClient(dbType string) (ISQL, error)
}





