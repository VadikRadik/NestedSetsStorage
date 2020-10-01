package configs

// Config is the configuration parameters for service
type Config struct {
	DbConnectionSting string `toml:"db_connection_string"`
	DbDriver          string `toml:"db_driver"`
	APIPort           string `toml:"api_port"`
	APIKey            string `toml:"api_key"`
}
