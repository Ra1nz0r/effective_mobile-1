package config

import "github.com/spf13/viper"

type Config struct {
	ServerHost       string `mapstructure:"ADDRESS"`
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	MigrationPath    string `mapstructure:"MIGRATION_PATH"`
	ExternalApiURL   string `mapstructure:"EXTERNAL_API_URL"`
	PaginationLimit  int    `mapstructure:"PAGINATION_LIMIT"`
	DatabaseUser     string `mapstructure:"DB_USER"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseName     string `mapstructure:"DB_NAME"`
	DatabaseDriver   string `mapstructure:"DB_DRIVER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
