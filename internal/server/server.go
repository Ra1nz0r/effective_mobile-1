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

// Запускает сервер.
func Run() {
	// Загружаем переменные окружения из '.env' файла.
	cfg, errLoad := config.LoadConfig(".")
	if errLoad != nil {
		log.Fatal("cannot load config", errLoad)
	}

	if errLog := logger.Initialize(cfg.LogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	// ================================================
	connect, errConn := sql.Open(cfg.DatabaseDriver, dbURL)
	if errConn != nil {
		logger.Zap.Fatal(errConn)
	}
	queries := hd.NewHandlerQueries(connect, cfg)
	// ================================================

	// Проверяем существование table в базе данных.
	exists, errExs := srv.TableExists(connect, "library")
	if errExs != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to check if table exists: %w", errExs))
	}

	// Создаём table, если он не существует.
	if !exists {
		if errRunMigr := srv.RunMigrations(dbURL, cfg.MigrationPath); errRunMigr != nil {
			logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", errConn))
		}
	}
	// ================================================

	logger.Zap.Info("Running handlers.")

	r := chi.NewRouter()

	r.Group(func(r chi.Router) { // исправить эндпойнты на другие
		r.Delete("/api/delete", queries.DeleteSong)
		r.Get("/list", queries.ListAllSongsWithFilters)
		r.Get("/songs/verse", queries.TextSongWithPagination)
		r.Post("/api/add", queries.AddSongInLibrary)
		r.Put("/api/update", queries.UpdateSong)
	})

	r.Group(func(r chi.Router) {
		//r.Use(hs.WithResponseDetails)
		// Разбить по переопределённым методам
	})

	logger.Zap.Info(fmt.Sprintf("Starting server on: '%s'", cfg.ServerHost))

	srv := http.Server{
		Addr:         cfg.ServerHost,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

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
