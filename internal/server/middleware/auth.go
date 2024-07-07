package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Auth(
	logger *zap.SugaredLogger,
	insp *authinspector.AuthInspector,

) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth middleware error: "
		authData := c.Param("Authorization")
		if authData == "" {
			logger.Error(e.Wrap(op, fmt.Errorf("empty field")))
			c.Status(http.StatusBadRequest)
			return
		}
		user := &obj.User{}
		err := json.Unmarshal([]byte(authData), user)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusBadRequest)
			return
		}

		ok := insp.CheckAuth(user)
		if !ok {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
