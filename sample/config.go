package sample

type DatabaseConfig struct {
	Host         string `yaml:"host" pflag:",Database host url"`
	Username     string `yaml:"username" pflag:",Database username"`
	Password     string `yaml:"password" pflag:",Database password"`
	DatabaseName string `yaml:"database" pflag:",Database name"`
	Port         string `yaml:"port" pflag:",Database Port"`
	SSL          string `yaml:"ssl" pflag:",Database ssl enable/disable"`
	ProjectID    string `yaml:"project_id" pflag:",BigQuery project ID"`
	JSONKeyPath  string `yaml:"json_key_path" pflag:",Path to BigQuery JSON key file"`
	DBType       string `yaml:"type" pflag:",Database type"`
}

// type DatabaseType string

type Config struct {
	Database  DatabaseConfig `yaml:"database"`
	Debug     bool           `yaml:"debug" pflag:",Debug mode"`
	Warehouse string         `yaml:"warehouse" pflag:",Snowflake warehouse"`
	Schema    string         `yaml:"schema" pflag:",Snowflake database schema"`
	Account   string         `yaml:"account" pflag:",Snowflake account ID"`
}

func NewConfig(conig Config) *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:         conig.Database.Host,
			Username:     conig.Database.Username,
			Password:     conig.Database.Password,
			DatabaseName: conig.Database.DatabaseName,
			Port:         conig.Database.Port,
			SSL:          conig.Database.SSL,
			ProjectID:    conig.Database.ProjectID,
			JSONKeyPath:  conig.Database.JSONKeyPath,
			DBType:       conig.Database.DBType,
		},
		Debug:     conig.Debug,
		Warehouse: conig.Warehouse,
		Schema:    conig.Schema,
		Account:   conig.Account,
	}

}
