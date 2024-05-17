package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

var DB_PASSWORD = "DB_PASSWORD"

const (
	Redshift_Schema_query = `SELECT "column", type, encoding, distkey, sortkey, "notnull"  FROM pg_table_def WHERE schemaname = '%s' AND tablename = '%s';`
	Redshift_Tables_query = "SHOW TABLES FROM SCHEMA %s.public;"
)

type Redshift struct {
	Client *sql.DB
	Config config.Config
}

func NewRedshift(client *sql.DB) (types.ISQL, error) {
	return &Redshift{
		Client: client,
		Config: config.Config{},
	}, nil
}

func NewRedshiftWithConfig(cfg *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 {
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}
	DB_PASSWORD = os.Getenv(DB_PASSWORD)

	dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, DB_PASSWORD, cfg.DatabaseName, cfg.SSL)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error creating a new session : %v", err)
	}

	return &Redshift{
		Client: db,
		Config: *cfg,
	}, nil
}

func (r *Redshift) Schema(table string) (types.Table, error) {
	if len(r.Config.Schema) == 0 {
		r.Config.Schema = "public"
	}

	query := fmt.Sprintf(Redshift_Schema_query, r.Config.Schema, table)
	ctx := context.Background()
	rows, err := r.Client.QueryContext(ctx, query)
	if err != nil {
		return types.Table{}, fmt.Errorf("error executing query: %v", err)
	}

	var columns []types.Column
	for rows.Next() {
		var column types.Column
		var encoding string
		var distkey bool
		var sortkey int
		var notnull bool
		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&encoding,
			&distkey,
			&sortkey,
			&notnull,
		); err != nil {
			return types.Table{}, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Metatags = []string{encoding, fmt.Sprintf("distkey:%v", distkey), fmt.Sprintf("sortkey:%d", sortkey), fmt.Sprintf("notnull:%v", notnull)}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return types.Table{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil
}

type TableResponse struct {
	TableName    string
	DatabaseName string
	SchemaName   string
	TableType    string
	TableAcl     sql.NullString
	Remarks      sql.NullString
}

func (r *Redshift) Tables(databaseName string) ([]string, error) {
	// ctx := context.Background()
	query := fmt.Sprintf(Redshift_Tables_query, databaseName)

	res, err := r.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	var tables []string

	for res.Next() {
		var table TableResponse
		if err := res.Scan(&table.DatabaseName, &table.SchemaName, &table.TableName, &table.TableType, &table.TableAcl, &table.Remarks); err != nil {
			return nil, fmt.Errorf("error scanning result: %v", err)
		}
		fmt.Println(table)
		tables = append(tables, table.TableName)
	}
	fmt.Println(tables)

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over result: %v", err)
	}

	return tables, nil

}

func (r *Redshift) Execute(query string) ([]byte, error) {
	ctx := context.Background()
	rows, err := r.Client.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// getting the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		// create a slice of values and pointers
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			//  create a slice of pointers to the values
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

func (r *Redshift) GenerateCreateTableQuery(table types.Table) string {
	query := fmt.Sprintf("CREATE TABLE %s.%s.%s (", r.Config.DatabaseName, r.Config.Schema, table.Name)
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + colType

		if column.IsPrimary {
			query += " PRIMARY KEY"
			if column.AutoIncrement {
				query += fmt.Sprintf(" IDENTITY(%v, %v)", column.IdentitySeed, column.IdentityStep)
			}
		}

		if column.IsNullable == "NO" {
			query += " NOT NULL"
		}

		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ");"
	return query
}
