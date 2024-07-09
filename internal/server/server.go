package httpserver

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	"github.com/eqkez0r/gophermart/internal/server/handlers"
	"github.com/eqkez0r/gophermart/internal/server/middleware"
	"github.com/eqkez0r/gophermart/internal/storage"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	engine *gin.Engine
	cfg    *config.Config
	logger *zap.SugaredLogger
}

const (
	APIUserRoute = "/api/user"
)

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	s storage.Storage,
	insp *authinspector.AuthInspector,
) (*HTTPServer, error) {
	const op = "Initial server error"

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	//middleware
	engine.Use(middleware.Logger(logger))
	//handlers
	authApi := engine.Group(APIUserRoute)
	authApi.POST(handlers.RegisterHandlerPath, handlers.RegisterHandler(ctx, logger, s, insp))
	authApi.POST(handlers.AuthHandlerPath, handlers.AuthHandler(ctx, logger, insp))

	interactApi := engine.Group(APIUserRoute)
	interactApi.Use(middleware.Logger(logger), middleware.Auth(logger, insp))
	interactApi.POST(handlers.NewOrderHandlerPath, handlers.NewOrderHandler(ctx, logger, s))

	server := &HTTPServer{
		server: &http.Server{
			Addr:    cfg.RunAddress,
			Handler: engine,
		},
		engine: engine,
		cfg:    cfg,
		logger: logger,
	}

	return server, nil
}

func (s *HTTPServer) Run(ctx context.Context) {
	const op = "Server run error: "

	go func() {
		s.logger.Infof("Server was started on %s", s.cfg.RunAddress)
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Error(e.Wrap(op, err))
		}
	}()

	<-ctx.Done()
	s.GracefulShutdown(ctx)
}

func (s *HTTPServer) GracefulShutdown(ctx context.Context) {
	const op = "Graceful shutdown error: "
	s.logger.Info("Server was stopped.")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(e.Wrap(op, err))
	}
}
