package configs

// Config is the configuration parameters for service
type Config struct {
	DbConnectionSting string `toml:"db_connection_string"`
	DbDriver          string `toml:"db_driver"`
}
