package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/hash"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
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
	storage RegisterUserProvider,
	insp *authinspector.AuthInspector,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in register handler: "
		newUser := &obj.User{}
		userInbytes := c.Param("Authorization")
		err := json.Unmarshal([]byte(userInbytes), newUser)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusBadRequest)
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
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		parseT, err := time.Parse(time.RFC3339, time.Now().String())
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		err = insp.Auth(ctx, newUser, parseT)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(newUser)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("Authorization", string(b))
		c.Status(http.StatusOK)
	}
}
