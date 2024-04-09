package sql

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// BigQueryRepo is a repository for interacting with BigQuery.
type BigQueryRepo struct {
	Client *bigquery.Client
	Config *config.Config
}

// Schema retrieves the schema information for the specified tables in BigQuery.
func (r BigQueryRepo) Schema(tables []string) ([]TableContext, int64, error) {
	var schema []TableContext

	var colCount int64 = 0

	ctx := context.Background()
	for _, table := range tables {
		tableRef := r.Client.Dataset(r.Config.Database.Database).Table(table)
		schemaInfo, err := tableRef.Metadata(ctx)
		if err != nil {
			return nil, 0, err
		}

		var rows []map[string]interface{}

		for _, field := range schemaInfo.Schema {

			colInfo := map[string]interface{}{
				"column_name": field.Name,
				"data_type":   field.Type,
			}
			if field.Type == bigquery.RecordFieldType {
				b, _ := field.Schema.ToJSONFields()
				colInfo["column_key"] = string(b)
			}
			rows = append(rows, colInfo)
		}

		b, _ := json.Marshal(rows)
		var tableStructure []ColContext
		json.Unmarshal(b, &tableStructure)
		schema = append(schema, TableContext{
			Name: table,
			Data: tableStructure,
		})
		colCount = 0
	}
	return schema, colCount, nil
}

// Execute executes a SQL-like query in BigQuery and returns the result as a list of maps.
func (r BigQueryRepo) Execute(query string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	q := r.Client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}

	for {
		var values map[string]bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		row := map[string]interface{}{}
		for key, value := range values {
			row[key] = value
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func (r BigQueryRepo) Tables(dataset string) ([]string, error) {
	result := r.Client.Dataset(dataset).Tables(context.Background())
	var tables []string
	for {
		table, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		tables = append(tables, table.TableID)
	}
	return tables, nil
}

// NewBigQueryRepo creates a new BigQuery repository based on the provided configuration.
func NewBigQueryRepo(cfg *config.Config) (ISQL, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, cfg.Database.ProjectID, option.WithCredentialsFile(cfg.Database.JSONKeyPath))
	if err != nil {
		return nil, err
	}

	return BigQueryRepo{
		Client: client,
		Config: cfg,
	}, nil
}