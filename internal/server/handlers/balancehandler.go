package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	BalanceHandlerPath = "/balance"
)

func BalanceHandler(
	ctx context.Context,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
