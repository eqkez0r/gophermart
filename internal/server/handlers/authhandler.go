package handlers

import (
	"context"
	"encoding/json"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	AuthHandlerPath = "/login/"
)

func AuthHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	insp *authinspector.AuthInspector,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Error in auth handler: "

		user := &obj.User{}
		userInbytes := c.Param("Authorization")
		err := json.Unmarshal([]byte(userInbytes), user)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusBadRequest)
			return
		}
		t := time.Now()
		formattedTime := t.Format(time.RFC3339)
		err = insp.Auth(ctx, user, formattedTime)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	}
}
