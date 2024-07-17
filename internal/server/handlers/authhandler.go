package handlers

import (
	"context"
	"fmt"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
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

		u := &obj.User{}
		if c.ContentType() != "application/json" {
			logger.Error(e.Wrap(op, errInvalidFormat))
			c.Status(http.StatusBadRequest)
			return
		}
		err := c.BindJSON(u)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		if u.Login == "" || u.Password == "" {
			logger.Error(e.Wrap(op, fmt.Errorf("empty field")))
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := storage.GetUser(ctx, u.Login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)) != nil {
			logger.Error(e.Wrap(op, fmt.Errorf("invalid password")))
			c.Status(http.StatusUnauthorized)
			return
		}
		logger.Infof("Authorized user: %v", user)
		token, err := jwt.CreateJWT(u.Login)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Authorization", token)

		c.Status(http.StatusOK)
	}
}
