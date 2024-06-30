package main

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	httpserver "github.com/eqkez0r/gophermart/internal/server"
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

	settings, err := config.NewConfig()
	if err != nil {
		suggaredLogger.Fatal(err)
	}

	server, err := httpserver.New()
	if err != nil {
		suggaredLogger.Fatal(err)
	}
	server.Run(ctx)
}
