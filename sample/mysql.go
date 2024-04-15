package sample

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

// dbURL generates a MySQL database URL based on the configuration.
func dbURL(dbConfig *DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		dbConfig.Password,
		// constant.DataSherlockPassword,
		dbConfig.Host,
		dbConfig.DatabaseName,
		dbConfig.SSL,
	)
}

// MysqlRepo is a repository for interacting with MySQL databases.
type MysqlRepo struct {
	Client *sql.DB
}

// Schema retrieves the schema information for the specified tables.
func (r *MysqlRepo) Schema(tables []string) ([]TableContext, int64, error) {
	var schema []TableContext

	var colCount int64 = 0
	query := "SELECT column_name AS column_name, data_type AS data_type, column_key AS column_key FROM information_schema.columns WHERE table_name = ?"
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
func (r *MysqlRepo) Execute(query string) ([]map[string]interface{}, error) {
	rows, err := r.Client.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	str, _ := rowsToJSON(rows)
	var data []map[string]interface{}
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Tables retrieves the names of tables in the specified database.
func (r *MysqlRepo) Tables(database string) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ?"
	rows, err := r.Client.Query(query, database)
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

// NewMysqlRepo creates a new MySQL repository based on the provided configuration.
func NewMysqlRepo(cfg *Config) (ISQL, error) {
	if os.Getenv("DB_PASSWORD") == "" || len(os.Getenv("DB_PASSWORD")) == 0 {
		return nil, fmt.Errorf("please set DB_PASSWORD env variable for the database")
	}
	db, err := sql.Open(cfg.DBType, dbURL(&cfg.Database))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &MysqlRepo{
		Client: db,
	}, nil
}