package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/server"
	"github.com/optivor/optivor/internal/storage/s3"
)

func main() {
	configPath := flag.String("config", "", "path to config file (default: ./optivor.yaml)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := server.NewLogger(cfg, os.Stdout)

	storageDriver, err := s3.New(cfg.Storage.S3)
	if err != nil {
		logger.Error("Failed to initialize S3 storage driver", "error", err)
		os.Exit(1)
	}

	cacheStore, err := fs.New(cfg.Cache.FS.Dir, cfg.Cache.FS.MaxSizeMB*1024*1024)
	if err != nil {
		logger.Error("Failed to initialize filesystem cache", "error", err)
		os.Exit(1)
	}

	pipe := pipeline.NewPipeline()
	defer pipeline.ShutdownVips()

	srv := server.New(cfg, storageDriver, cacheStore, pipe, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Optivor binary started successfully", "port", cfg.Server.Port)

	<-stop
	logger.Info("Received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.WriteTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Optivor stopped cleanly")
}
