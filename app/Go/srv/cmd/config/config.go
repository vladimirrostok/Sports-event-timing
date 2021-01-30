package config

// config declares connection details.
type Config struct {
	DBHost     string `mapstructure:"db_host"`
	DBDriver   string `mapstructure:"db_driver"`
	DBUsername string `mapstructure:"db_username"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	DBPort     string `mapstructure:"db_port"`

	APIAddress     string `mapstructure:"api_address"`
	TestAPIAddress string `mapstructure:"test_api_address"`
}
