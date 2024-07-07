package main

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	httpserver "github.com/eqkez0r/gophermart/internal/server"
	"github.com/eqkez0r/gophermart/internal/storage"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	"go.uber.org/zap"
	"log"
	"os/signal"
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

	ins := authinspector.New(suggaredLogger, s)
	go ins.Observe(ctx)

	server, err := httpserver.New(ctx, cfg, suggaredLogger, s, ins)
	if err != nil {
		suggaredLogger.Fatal(err)
	}
	server.Run(ctx)
}
