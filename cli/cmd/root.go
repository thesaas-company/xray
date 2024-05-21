package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	xrayTypes "github.com/thesaas-company/xray/types"
	"gopkg.in/yaml.v3"
)

// Command line flags
var (
	verbose bool
	cfgFile string
	dbType  string
)

// QueryResult represents the result of a database query.
type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Time    float64         `json:"time"`
	Error   string          `json:"error"`
}

// Table represents a table with headers and rows.
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new Table with the given headers.
func NewTable(headers []string) *Table {
	return &Table{headers: headers}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

// String returns a string representation of the table.
func (t *Table) String() string {
	// Find the maximum width of each column
	columnWidths := make([]int, len(t.headers))
	for i, header := range t.headers {
		columnWidths[i] = len(header)
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	// Create a format string based on the column widths
	var formatBuilder strings.Builder
	for _, width := range columnWidths {
		formatBuilder.WriteString(fmt.Sprintf("%%-%ds ", width))
	}
	formatString := formatBuilder.String()

	// Print the headers
	var result strings.Builder
	result.WriteString(fmt.Sprintf(formatString, toInterfaceSlice(t.headers)...))
	result.WriteRune('\n')

	// Print a separator line
	for _, width := range columnWidths {
		result.WriteString(strings.Repeat("-", width) + " ")
	}
	result.WriteRune('\n')

	// Print the rows
	for _, row := range t.rows {
		result.WriteString(fmt.Sprintf(formatString, toInterfaceSlice(row)...))
		result.WriteRune('\n')
	}

	return result.String()
}

// toInterfaceSlice converts a slice of strings to a slice of interfaces.
func toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// Command for interacting with databases
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with databases",
	Long: `
	This command provides an interactive shell to execute SQL queries on various types of databases. 
	It supports MySQL, PostgreSQL, MSSQL, Redshift, Bigquery and Snowflake. 
	To use this command, you need to provide a configuration file with the --config flag or -c flag,
	and a database type with the --type flag or -t flag. 

	The configuration file should be in YAML format and contain the necessary database connection parameters 
	such as host, username, database name, port, and SSL settings.

	You can also control the verbosity of the command's output with the --verbose or -v flag. 
	When the verbose mode is on, the command will log additional information about its operation.

	In the interactive shell, you can type SQL queries and press Enter to execute them. 
	The results will be displayed in the console. Type 'exit' to leave the shell`,

	Run: func(cmd *cobra.Command, args []string) {
		// Set up logging
		if !verbose {
			logrus.SetOutput(io.Discard)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if cfgFile == "" {
			fmt.Println("Error: Configuration file path is missing. Please use the --config flag to specify the path to your configuration file.")
			return
		}

		// Read the YAML file
		configData, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("Error: Failed to read YAML file: %v\n", err)
			return
		}
		var cfg config.Config
		err = yaml.Unmarshal(configData, &cfg)
		if err != nil {
			fmt.Printf("Error: Failed to unmarshal YAML: %v\n", err)
			return
		}

		db, err := xray.NewClientWithConfig(&cfg, parseDbType(dbType))
		if err != nil {
			fmt.Printf("Error: Failed to connect to database: %s: %v\n", dbType, err)
			return
		}

		fmt.Println("Welcome to database shell!")

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			query, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading query:", err)
				continue
			}
			query = strings.TrimSpace(query)
			if query == "exit" {
				fmt.Println("Exiting shell.")
				break
			}

			// Check if the query is a PostgreSQL meta command
			// Check if the query is a PostgreSQL meta command
			if dbType == "postgres" {
				query = PostgresMetaCommands(query)
			}

			b, err := db.Execute(query)
			if err != nil {
				fmt.Println("Error executing query:", err)
				continue
			}

			var result QueryResult
			err = json.Unmarshal(b, &result)
			if err != nil {
				fmt.Println("Error parsing query result:", err)
				continue
			}

			if len(result.Rows) == 0 {
				fmt.Println("No results found.")
				continue
			}

			table := NewTable(result.Columns)
			for _, row := range result.Rows {
				stringRow := make([]string, len(row))
				for i, v := range row {
					switch value := v.(type) {
					case string:
						// Check if the string is base64 encoded
						if isBase64(value) {
							decodedValue, err := base64.StdEncoding.DecodeString(value)
							if err != nil {
								fmt.Println("Error decoding base64 value:", err)
								stringRow[i] = value // Use original value if decoding fails
							} else {
								stringRow[i] = string(decodedValue)
							}
						} else {
							stringRow[i] = value
						}
					default:
						stringRow[i] = fmt.Sprintf("%v", value)
					}
				}
				table.AddRow(stringRow)
			}

			// Print the table
			fmt.Println(table.String())
		}
	},
}

// isBase64 checks if a string is base64 encoded
func isBase64(s string) bool {
	if len(s)%4 != 0 {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// Execute runs the command line interface.
func Execute() {
	rootCmd := &cobra.Command{Use: "xray"}

	rootCmd.AddCommand(shellCmd)
	shellCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	shellCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config.yaml")
	shellCmd.PersistentFlags().StringVarP(&dbType, "type", "t", "mysql", "Database type like mysql, postgres, bigquery")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
}

// ParseDbType parses a string and returns the corresponding DbType.
func parseDbType(s string) xrayTypes.DbType {
	switch strings.ToLower(s) {
	case "mysql":
		return xrayTypes.MySQL
	case "postgres":
		return xrayTypes.Postgres
	case "snowflake":
		return xrayTypes.Snowflake
	case "bigquery":
		return xrayTypes.BigQuery
	case "redshift":
		return xrayTypes.Redshift
	case "mssql":
		return xrayTypes.MSSQL
	default:
		return xrayTypes.MySQL
	}
}

// PostgresMetaCommands translates PostgreSQL meta commands to SQL queries
func PostgresMetaCommands(query string) string {
	switch query {
	case "\\l":
		return "SELECT datname FROM pg_database WHERE datistemplate = false;"
	case "\\dt":
		return "SELECT * FROM pg_catalog.pg_tables;"
	case "\\d":
		return "SELECT * FROM pg_catalog.pg_tables;"
	case "\\c":
		return "switch_database"
	case "\\q":
		return "exit"
	case "\\?":
		return "help"
	case "\\h":
		return "help"
	case "\\du":
		return "SELECT * FROM pg_catalog.pg_roles;"
	case "\\conninfo":
		return "SELECT * FROM pg_stat_activity WHERE pid = pg_backend_pid();"
	default:
		// Handle meta commands with parameters
		if strings.HasPrefix(query, "\\c ") {
			dbName := strings.TrimPrefix(query, "\\c ")
			return fmt.Sprintf("switch_database %s", dbName)
		} else if strings.HasPrefix(query, "\\d ") {
			tableName := strings.TrimPrefix(query, "\\d ")
			return fmt.Sprintf("SELECT * FROM %s;", tableName)
		} else if strings.HasPrefix(query, "\\dn ") {
			schemaName := strings.TrimPrefix(query, "\\dn ")
			return fmt.Sprintf("SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname = '%s';", schemaName)
		} else if strings.HasPrefix(query, "\\dp ") {
			tableName := strings.TrimPrefix(query, "\\dp ")
			return fmt.Sprintf("SELECT * FROM pg_catalog.pg_statio_all_tables WHERE relname = '%s';", tableName)
		}
	}
	// If the query doesn't match any known meta commands, return it unchanged
	return query
}
