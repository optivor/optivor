package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/cli"
	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/server"
	"github.com/optivor/optivor/internal/storage/s3"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "init" || cmd == "deploy" || cmd == "doctor" || cmd == "logs" || cmd == "metrics" || cmd == "driver" || cmd == "help" || cmd == "-v" || cmd == "--version" || cmd == "-h" || cmd == "--help" {
			if err := cli.Execute(); err != nil {
				os.Exit(1)
			}
			return
		}
	}

	configPath := flag.String("config", "", "path to config file (default: ./optivor.yaml)")
	provider := flag.String("provider", "", "override storage provider (e.g. s3, minio, r2)")
	flag.Parse()

	runServer(*configPath, *provider)
}

func runServer(configPath string, providerFlag string) {
	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if providerFlag != "" {
		cfg.Storage.Driver = providerFlag
	}

	logger := server.NewLogger(cfg, os.Stdout)

	drv := cfg.Storage.Driver
	if drv != "" && drv != "s3" && drv != "minio" && drv != "r2" {
		logger.Error("Unknown storage provider", "provider", drv)
		os.Exit(1)
	}

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
