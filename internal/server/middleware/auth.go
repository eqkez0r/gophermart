package middleware

import (
	"context"
	"fmt"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type GetUserProvider interface {
	IsUserExist(context.Context, string) (bool, error)
}

func Auth(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storage GetUserProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth middleware error: "
		token := c.Param("Authorization")
		if token == "" {
			logger.Error(e.Wrap(op, fmt.Errorf("empty field")))
			c.Status(http.StatusBadRequest)
			return
		}

		login, ttl, err := jwt.JWTPayload(token)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		ok, err := storage.IsUserExist(ctx, login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusUnauthorized)
			return
		}

		if !ok {
			logger.Error(e.Wrap(op, fmt.Errorf("user not found")))
			c.Status(http.StatusUnauthorized)
			return
		}

		if time.Now().After(ttl) {
			logger.Error(e.Wrap(op, fmt.Errorf("token expired")))
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
