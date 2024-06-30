package httpserver

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	server   *http.Server
	engine   *gin.Engine
	settings *config.Config
	logger   *zap.SugaredLogger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) (*HTTPServer, error) {
	const op = "Initial server error"

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	//middleware
	engine.Use()
	//handlers

	//var err error
	//var conn *pgxpool.Pool = nil

	server := &HTTPServer{
		server: &http.Server{
			Addr:    cfg.RunAddress,
			Handler: engine,
		},
		engine:   engine,
		settings: cfg,
		logger:   logger,
	}

	return server, nil
}

func (s *HTTPServer) Run(ctx context.Context) {
	const op = "Server run error: "

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Error(e.Wrap(op, err))
		}
	}()
}

func (s *HTTPServer) GracefulShutdown(ctx context.Context) {
	s.logger.Info("Server was stopped.")

}
