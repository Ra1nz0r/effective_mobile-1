package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ra1nz0r/effective_mobile-1/internal/config"
	hs "github.com/Ra1nz0r/effective_mobile-1/internal/handlers"
	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
	"github.com/go-chi/chi/v5"
)

// Запускает агент, который будет принимать метрики от агента.
func Run() {
	config.ServerFlags()

	r := chi.NewRouter()

	if errLog := logger.Initialize(config.DefLogLevel); errLog != nil {
		log.Fatal(errLog)
	}

	logger.Zap.Info("Running handlers.")

	r.Group(func(r chi.Router) {
		//r.Use(hs.WithRequestDetails)
		//r.Post("/update/{type}/{name}/{value}", hs.UpdateMetrics)
	})

	r.Group(func(r chi.Router) {
		//r.Use(hs.WithResponseDetails)
		r.Get("/", hs.GetAllMetrics)
		//r.Get("/value/{type}/{name}", hs.GetMetricByName)
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
