package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eqkez0r/gophermart/internal/storage"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/luhn"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

const (
	WithdrawHandlerPath = "/withdraw"
)

type WithdrawHandlerProvider interface {
	NewWithdraw(context.Context, uint64, uint64, float64) error
}

func WithdrawHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store WithdrawHandlerProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in withdraw handler: "

		token := c.Request.Header.Get("Authorization")
		withdraw := &obj.Withdraw{}
		_, userId, _, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, withdraw)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		number, err := strconv.Atoi(withdraw.Order)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		if !luhn.Valid(uint64(number)) {
			logger.Error(e.Wrap(op, errors.New("invalid order number")))
			c.Status(http.StatusUnprocessableEntity)
			return
		}

		err = store.NewWithdraw(ctx, userId, uint64(number), withdraw.Sum)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			switch {
			case errors.Is(err, storage.ErrIsOrderIsNotExist):
				{
					c.Status(http.StatusUnprocessableEntity)
				}
			case errors.Is(err, storage.ErrBalanceIsNotEnough):
				{
					c.Status(http.StatusPaymentRequired)
				}
			default:
				{
					c.Status(http.StatusInternalServerError)
				}
			}
			return
		}

		c.Status(http.StatusOK)
	}
}
