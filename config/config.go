package config

// Add Logging, You can use any lib
// TODO: Add code docs
// Once we are done with mysql and postgres, Let's rethink about the config structure
type Config struct {
	Host         string `yaml:"host" pflag:",Database host url"`
	Username     string `yaml:"username" pflag:",Database username"`
	DatabaseName string `yaml:"database" pflag:",Database name"`
	Port         string `yaml:"port" pflag:",Database Port"`
	SSL          string `yaml:"ssl" pflag:",Database ssl enable/disable"`
	ProjectID    string `yaml:"project_id" pflag:",BigQuery project ID"`
	JSONKeyPath  string `yaml:"json_key_path" pflag:",Path to BigQuery JSON key file"`
	Warehouse    string `yaml:"warehouse" pflag:",Snowflake warehouse"`
	Schema       string `yaml:"schema" pflag:",Snowflake database schema"`
	Account      string `yaml:"account" pflag:",Snowflake account ID"`
	Debug        bool   `yaml:"debug" pflag:",Debug mode"`
}
