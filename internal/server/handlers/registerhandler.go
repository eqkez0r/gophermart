package handlers

import (
	"context"
	"fmt"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	RegisterHandlerPath = "/register"
)

type RegisterUserProvider interface {
	NewUser(context.Context, *obj.User) error
}

func RegisterHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storage RegisterUserProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in register handler: "
		newUser := &obj.User{}
		err := c.BindJSON(newUser)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Infof("user %v", newUser)
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
		err = storage.NewUser(ctx, newUser)
		if err != nil {
			//switch error
			c.Status(http.StatusInternalServerError)
			return
		}
		//TODO:Auth here
		c.Status(http.StatusOK)
	}
}
