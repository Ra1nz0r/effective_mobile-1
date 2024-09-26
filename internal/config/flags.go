package config

import (
	"flag"
	"os"
)

var (
	DefServerHost = "0.0.0.0:7654" // стандартный адрес для сервера
	DefLogLevel   = "info"         // стандартный уровень логирования
)

// Создаёт флаги для запуска сервера, если в терминале переданы переменные окружения,
// то приоритет будет отдаваться им.
func ServerFlags() {
	flag.StringVar(&DefServerHost, "a", DefServerHost, "address and port to run server")
	flag.StringVar(&DefLogLevel, "l", DefLogLevel, "set log level")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		DefServerHost = envServerAddress
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		DefLogLevel = envLogLevel
	}
}
