package handlers

import (
	"fmt"
	"github.com/eqkez0r/gophermart/internal/storage"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	RegisterHandlerPath = "/register"
)

func RegisterHandler(
	logger *zap.SugaredLogger,
	storage storage.Storage) gin.HandlerFunc {
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
		//TODO:hash pass here
		err = storage.NewUser(newUser)
		if err != nil {
			//switch error
			c.Status(http.StatusInternalServerError)
			return
		}
		//TODO:Auth here
		c.Status(http.StatusOK)
	}
}
