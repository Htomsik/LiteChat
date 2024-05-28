package apiServer

// Config ...
type Config struct {
	Port                  string `toml:"port"`
	LogLevel              string `toml:"logLevel"`
	DatabaseType          string `toml:"databaseType"`
	DatabaseURL           string `toml:"databaseURL"`
	DatabaseAutoMigration bool   `toml:"databaseAutoMigration"`
	SessionsKey           string `toml:"sessionKey"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
