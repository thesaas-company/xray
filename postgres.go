package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresRepo is a repository for interacting with PostgreSQL databases.
type PostgresRepo struct {
	Client *sql.DB
}

// Schema retrieves the schema information for the specified tables.
func (r PostgresRepo) Schema(tables []string) ([]TableContext, int64, error) {
	var schema []TableContext

	var colCount int64 = 0
	query := "SELECT CAST(column_name AS TEXT), CAST(data_type AS TEXT) FROM information_schema.columns WHERE table_name::text = $1;"

	for _, table := range tables {
		rows, err := r.Client.Query(query, table)
		if err != nil {
			return nil, 0, err
		}
		data, c := rowsToJSON(rows)
		var tableStructure []ColContext
		json.Unmarshal([]byte(data), &tableStructure)
		schema = append(schema, TableContext{
			Name: table,
			Data: tableStructure,
		})
		colCount = c
	}
	return schema, colCount, nil
}

// Execute executes a SQL query and returns the result as a list of maps.
func (r PostgresRepo) Execute(query string) ([]map[string]interface{}, error) {
	rows, err := r.Client.Query(query)
	if err != nil {
		return nil, err
	}

	str, _ := rowsToJSON(rows)

	var data []map[string]interface{}
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Tables retrieves the names of tables in the specified database.
func (r PostgresRepo) Tables(_ string) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := r.Client.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tableNames []string

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tableNames, nil
}

// NewPostgresRepo creates a new PostgreSQL repository based on the provided configuration.
func NewPostgresRepo(cfg *config.Config) (ISQL, error) {
	db, err := sql.Open(cfg.Type, fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, constant.DataSherlockPassword, cfg.Database.Database, cfg.Database.SSL))
	if err != nil {
		return nil, err
	}
	return PostgresRepo{
		Client: db,
	}, nil
}