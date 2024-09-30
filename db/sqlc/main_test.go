package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"fmt"

	"github.com/Ra1nz0r/effective_mobile-1/internal/config"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	cfg, errLoad := config.LoadConfig("../..")
	if errLoad != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", errLoad))
	}

	// Конфигурируем путь для подключения к PostgreSQL.
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	connect, errConn := sql.Open(cfg.DatabaseDriver, dbURL)
	if errConn != nil {
		logger.Zap.Fatal(fmt.Errorf("unable to create connection to database: %w", errConn))
	}

	testQueries = New(connect)

	os.Exit(m.Run())

}
