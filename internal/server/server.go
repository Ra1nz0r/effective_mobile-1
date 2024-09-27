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
	"github.com/Ra1nz0r/effective_mobile-1/internal/services"
	"github.com/go-chi/chi/v5"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Запускает сервер.
func Run() {
	config.ServerFlags()

	if errLog := logger.Initialize(config.DefLogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	dbURL := "postgres://postgres:admin@localhost:5432/library?sslmode=disable" // вынести в конфиг или енв

	// ================================================
	conn, errConn := sql.Open("pgx", dbURL) // вынести в конфиг или енв, драйвер
	if errConn != nil {
		log.Fatal(errConn)
	}
	queries := hd.NewHandleQueries(conn)
	// ================================================

	// Проверяем существование table в базе данных.
	exists, errExs := services.TableExists(conn, "library") // вынести в конфиг или енв, лайбрари
	if errExs != nil {
		log.Fatalf("Failed to check if table exists: %v", errExs)
	}

	// Создаём table, если он не существует.
	if !exists {
		if errRunMigr := services.RunMigrations(dbURL); errRunMigr != nil {
			log.Fatalf("Failed to run migrations: %v", errConn)
		}
	}
	// ================================================

	logger.Zap.Info("Running handlers.")

	r := chi.NewRouter()

	r.Group(func(r chi.Router) { // исправить эндпойнты на другие
		r.Put("/api/update", queries.UpdateSong)
		r.Delete("/api/delete", queries.DeleteSong)
		r.Post("/api/add", queries.AddSongInLibrary)
	})

	r.Group(func(r chi.Router) {
		//r.Use(hs.WithResponseDetails)
		r.Get("/list", queries.ListAllSongsWithFilters)
		r.Get("/", queries.GetAll)
		r.Get("/songs/verse", queries.TextSongWithPagination)
	})

	logger.Zap.Info(fmt.Sprintf("Starting server on: '%s'", config.DefServerHost))

	srv := http.Server{
		Addr:         config.DefServerHost,
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
