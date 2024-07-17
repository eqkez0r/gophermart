package handlers

import (
	"context"
	"errors"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	OrderListHandlerPath = "/orders"
)

type OrderListProvider interface {
	GetOrdersList(ctx context.Context, login string) ([]*obj.Order, error)
}

func OrderListHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store OrderListProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in new order list handler: "

		token := c.Request.Header.Get("Authorization")

		login, _, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		orders, err := store.GetOrdersList(ctx, login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			logger.Error(e.Wrap(op, errors.New("no orders found")))
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}
