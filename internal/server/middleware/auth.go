package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
