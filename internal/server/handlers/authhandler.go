package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	AuthHandlerPath = "/login/"
)

func AuthHandler(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in auth handler: "

	}
}
