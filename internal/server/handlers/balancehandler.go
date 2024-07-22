package handlers

import (
	"context"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	BalanceHandlerPath = ""
)

type BalanceProvider interface {
	GetBalance(ctx context.Context, login string) (*obj.AccrualBalance, error)
}

func BalanceHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store BalanceProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Balance handler error: "

		token := c.Request.Header.Get("Authorization")
		logger.Infof("call balance. token %s", token)
		login, _, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		balance, err := store.GetBalance(ctx, login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		logger.Infof("getting balance: %v", balance)
		c.JSON(http.StatusOK, balance)
	}
}
