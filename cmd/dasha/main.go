package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saintbyte/dasha-standalone-complex/internal/config"
	"github.com/saintbyte/dasha-standalone-complex/internal/httpapi"
)

func main() {
	cfg, path, err := config.Load()
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: levelFromString(cfg.LogLevel),
	}))
	slog.SetDefault(logger)
	slog.Info("запуск dasha-standalone-complex", "config", path)

	app := httpapi.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", cfg.Server.Addr)
		if err := app.Listen(cfg.Server.Addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(shutdownCtx)
}

func levelFromString(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
