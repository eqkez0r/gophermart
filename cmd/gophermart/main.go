package main

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	"github.com/eqkez0r/gophermart/internal/orderfetcher"
	httpserver "github.com/eqkez0r/gophermart/internal/server"
	"github.com/eqkez0r/gophermart/internal/storage"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	suggaredLogger := logger.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		suggaredLogger.Fatal(err)
	}

	s, err := storage.NewStorage(ctx, suggaredLogger, "postgresql", cfg.DatabaseURI)
	if err != nil {
		suggaredLogger.Fatal(err)
	}

	var wg sync.WaitGroup
	of := orderfetcher.New(suggaredLogger, cfg.AccrualSystemAddress, s)

	wg.Add(1)
	go of.Run(ctx, &wg)

	server, err := httpserver.New(ctx, cfg, suggaredLogger, s, of)
	if err != nil {
		suggaredLogger.Fatal(err)
	}
	server.Run(ctx)
	wg.Wait()
	s.GracefulShutdown()
}
