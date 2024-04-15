package library

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sf "github.com/snowflakedb/gosnowflake" // Import Snowflake driver
)

// SnowflakeRepo is a repository for interacting with Snowflake databases.
type SnowflakeRepo struct {
	Client *sql.DB
	Config *Config
}

// Schema retrieves the schema information for the specified tables in Snowflake.
func (r SnowflakeRepo) Schema(tables []string) ([]TableContext, int64, error) {
	var schema []TableContext

	var colCount int64 = 0
	query := "SELECT column_name::TEXT, data_type::TEXT FROM information_schema.columns WHERE table_name::TEXT = ?;"

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

// Execute executes a SQL query in Snowflake and returns the result as a list of maps.
func (r SnowflakeRepo) Execute(query string) ([]map[string]interface{}, error) {
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

// Tables retrieves the names of tables in the specified Snowflake database.
func (r SnowflakeRepo) Tables(_ string) ([]string, error) {
	query := fmt.Sprintf("USE WAREHOUSE %s", r.Config.Warehouse)

	_, err := r.Client.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	query = "SHOW TERSE TABLES"

	rows, err := r.Client.Query(query)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return nil, err
	}
	return tableNames, nil
}

// NewSnowflakeRepo creates a new Snowflake repository based on the provided configuration.
func NewSnowflakeRepo(cfg *Config) (ISQL, error) {
	dns, err := sf.DSN(&sf.Config{
		Account:   cfg.Account,
		User:      cfg.Database.Username,
		// Password:  constant.DataSherlockPassword,
		Password: cfg.Database.Password,
		Database:  cfg.Database.DatabaseName,
		Warehouse: cfg.Warehouse,
	})
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(string(cfg.DBType), dns)
	if err != nil {
		return nil, err
	}
	return SnowflakeRepo{
		Client: db,
		Config: cfg,
	}, nil
}