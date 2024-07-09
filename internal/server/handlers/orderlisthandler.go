package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	OrderListHandlerPath = "/orders"
)

func OrderListHandler(
	ctx context.Context,
) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
