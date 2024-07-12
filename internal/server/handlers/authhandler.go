package handlers

import (
	"context"
	"fmt"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	AuthHandlerPath = "/login/"
)

type GetUserProvider interface {
	GetUser(context.Context, string) (*obj.User, error)
}

func AuthHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storage GetUserProvider,

) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in auth handler: "

		newUser := &obj.User{}
		if c.ContentType() != "application/json" {
			logger.Error(e.Wrap(op, errInvalidFormat))
			c.Status(http.StatusBadRequest)
			return
		}
		err := c.BindJSON(newUser)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		if newUser.Login == "" || newUser.Password == "" {
			logger.Error(e.Wrap(op, fmt.Errorf("empty field")))
			c.Status(http.StatusBadRequest)
			return
		}
		newUser.Password, err = hash.HashPassword(newUser.Password)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		user, err := storage.GetUser(ctx, newUser.Login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password)) != nil {
			logger.Error(e.Wrap(op, fmt.Errorf("invalid password")))
			c.Status(http.StatusUnauthorized)
			return
		}

		token, err := jwt.CreateJWT(newUser.Login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Authorization", token)

		c.Status(http.StatusOK)
	}
}
