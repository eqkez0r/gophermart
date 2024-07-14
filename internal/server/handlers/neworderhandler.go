package handlers

import (
	"context"
	"errors"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	"github.com/eqkez0r/gophermart/utils/luhn"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

const (
	NewOrderHandlerPath = "/orders"
)

type NewOrderProvider interface {
	NewOrder(context.Context, uint64, uint64) error
}

type OrderFetcherToAccrualService interface {
	RegisterNewOrder(context.Context, uint64, uint64) error
}

func NewOrderHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store NewOrderProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in new order handler: "

		ct := c.ContentType()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if ct != "text/plain" && len(body) == 0 {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusBadRequest)
			return
		}

		number, err := strconv.ParseUint(string(body), 10, 64)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		if !luhn.Valid(number) {
			logger.Error(e.Wrap(op, errors.New("invalid order number")))
			c.Status(http.StatusUnprocessableEntity)
			return
		}
		token := c.Request.Header.Get("Authorization")
		_, userid, _, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusUnauthorized)
			return
		}
		if err = store.NewOrder(ctx, number, userid); err != nil {
			logger.Error(e.Wrap(op, err))
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				logger.Info(err, pgErr)
				if pgErr.Code == "23505" {
					logger.Info("Is order was accepted")
					c.Status(http.StatusOK)
					return
				}
			}
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusAccepted)

	}
}
