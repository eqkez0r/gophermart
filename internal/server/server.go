package httpserver

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/config"
	"github.com/eqkez0r/gophermart/internal/orderfetcher"
	"github.com/eqkez0r/gophermart/internal/server/handlers"
	"github.com/eqkez0r/gophermart/internal/server/middleware"
	"github.com/eqkez0r/gophermart/internal/storage"
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
	s      storage.Storage
}

const (
	APIUserRoute    = "/api/user"
	APIBalanceRoute = "/balance"
)

func New(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	s storage.Storage,
	of *orderfetcher.OrderFetcher,
) (*HTTPServer, error) {
	//const op = "Initial server error"

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectFixedPath = true

	//middleware
	engine.Use(middleware.Logger(logger))
	//handlers
	authAPI := engine.Group(APIUserRoute)
	authAPI.POST(handlers.RegisterHandlerPath, handlers.RegisterHandler(ctx, logger, s))
	authAPI.POST(handlers.AuthHandlerPath, handlers.AuthHandler(ctx, logger, s))

	userAPI := engine.Group(APIUserRoute)
	userAPI.Use(middleware.Logger(logger), middleware.Auth(ctx, logger, s), middleware.Gzip(logger))
	userAPI.POST(handlers.NewOrderHandlerPath, handlers.NewOrderHandler(ctx, logger, s))
	userAPI.GET(handlers.OrderListHandlerPath, handlers.OrderListHandler(ctx, logger, s))
	userAPI.GET(handlers.WithdrawalsHandlerPath, handlers.WithdrawalsHandler(ctx, logger, s))

	balanceAPI := userAPI.Group(APIBalanceRoute)
	balanceAPI.GET(handlers.BalanceHandlerPath, handlers.BalanceHandler(ctx, logger, s))
	balanceAPI.POST(handlers.WithdrawHandlerPath, handlers.WithdrawHandler(ctx, logger, s))

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

}

func (s *HTTPServer) GracefulShutdown(ctx context.Context) {
	const op = "Graceful shutdown error: "
	s.logger.Info("Server was stopped.")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(e.Wrap(op, err))
	}
	if err := s.s.GracefulShutdown(); err != nil {
		s.logger.Error(e.Wrap(op, err))
	}
}
