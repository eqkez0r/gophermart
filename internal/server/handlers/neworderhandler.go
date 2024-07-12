package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

const (
	NewOrderHandlerPath = "/orders"
)

type NewOrderProvider interface {
	NewOrder(ctx context.Context, number uint64) error
}

func NewOrderHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	store NewOrderProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in new order handler: "

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		log.Println(string(body))

	}
}
