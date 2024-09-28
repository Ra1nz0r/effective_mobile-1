package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ra1nz0r/effective_mobile-1/internal/config"
	hd "github.com/Ra1nz0r/effective_mobile-1/internal/handlers"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
	srv "github.com/Ra1nz0r/effective_mobile-1/internal/services"
	"github.com/go-chi/chi/v5"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Run запускает сервер.
func Run() {
	// Загружаем переменные окружения из '.env' файла.
	cfg, errLoad := config.LoadConfig(".")
	if errLoad != nil {
		log.Fatal("cannot load config", errLoad)
	}

	// Инициализируем логгер, с указанием уровня логирования.
	if errLog := logger.Initialize(cfg.LogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	// Конфигурируем путь для подключения к PostgreSQL.
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	logger.Zap.Debug("Connecting to the database.")
	// Открываем подключение к базе данных.
	connect, errConn := sql.Open(cfg.DatabaseDriver, dbURL)
	if errConn != nil {
		logger.Zap.Fatal(errConn)
	}

	// Передаём подключение и настройки приложения нашим обработчикам.
	queries := hd.NewHandlerQueries(connect, cfg)

	logger.Zap.Debug("Checking the existence of a table in the database.")
	// Проверяем существование table в базе данных.
	exists, errExs := srv.TableExists(connect, "library") //sdfsdf
	if errExs != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to check if table exists: %w", errExs))
	}

	// Создаём table, если он не существует.
	if !exists {
		logger.Zap.Info("Table not exists, creating.")
		if errRunMigr := srv.RunMigrations(dbURL, cfg.MigrationPath); errRunMigr != nil {
			logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", errConn))
		}
	}

	logger.Zap.Debug("Running handlers.")

	// Создаём router и endpoints.
	r := chi.NewRouter()

	r.Group(func(r chi.Router) { // исправить эндпойнты на другие
		r.Use(queries.WithRequestDetails)

		r.Delete("/api/delete", queries.DeleteSong)
		r.Post("/api/add", queries.AddSongInLibrary)
		r.Put("/api/update", queries.UpdateSong)
	})

	r.Group(func(r chi.Router) {
		r.Use(queries.WithResponseDetails)

		r.Get("/list", queries.ListAllSongsWithFilters)
		r.Get("/songs/verse", queries.TextSongWithPagination)
	})

	logger.Zap.Debug("Configuring and starting the server.")

	// Конфигурируем и запускаем сервер.
	srv := http.Server{
		Addr:         cfg.ServerHost,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	logger.Zap.Info(fmt.Sprintf("The server is running on: '%s'", cfg.ServerHost))

	go func() {
		if errListn := srv.ListenAndServe(); !errors.Is(errListn, http.ErrServerClosed) {
			logger.Zap.Fatal("HTTP server error:", errListn)
		}
		logger.Zap.Info("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if errShut := srv.Shutdown(shutdownCtx); errShut != nil {
		logger.Zap.Fatal("HTTP shutdown error", errShut)
	}
	logger.Zap.Info("Graceful shutdown complete.")
}
