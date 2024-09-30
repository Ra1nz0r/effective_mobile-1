package config

import "github.com/spf13/viper"

type Config struct {
	ServerHost       string `mapstructure:"ADDRESS"`          // адрес сервера
	LogLevel         string `mapstructure:"LOG_LEVEL"`        // уровень логирования
	MigrationPath    string `mapstructure:"MIGRATION_PATH"`   // путь до папки с миграциями
	ExternalAPIURL   string `mapstructure:"EXTERNAL_API_URL"` // адрес до внешнего API
	PaginationLimit  int32  `mapstructure:"PAGINATION_LIMIT"` // размер пагинации по-умолчанию
	DatabaseUser     string `mapstructure:"DB_USER"`          // имя пользователя датабазы
	DatabasePassword string `mapstructure:"DB_PASSWORD"`      // пароль пользователя датабазы
	DatabaseHost     string `mapstructure:"DB_HOST"`          // адрес для подключения к датабазе
	DatabasePort     string `mapstructure:"DB_PORT"`          // порт для подключения к датабазе
	DatabaseName     string `mapstructure:"DB_NAME"`          // имя датабазы
	DatabaseDriver   string `mapstructure:"DB_DRIVER"`        // драйвер датабазы
}

// LoadConfig загружает из файла '.env' переменные окружения.
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
