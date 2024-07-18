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
	WithdrawalsHandlerPath = "/withdrawals"
)

type WithdrawalsProvider interface {
	Withdrawals(context.Context, string) ([]*obj.Withdraw, error)
}

func WithdrawalsHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store WithdrawalsProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in withdrawals handler: "

		token := c.Request.Header.Get("Authorization")

		login, _, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		withdrawals, err := store.Withdrawals(ctx, login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			logger.Infof("No Withdrawals for user %s", login)
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusOK, withdrawals)
	}
}
