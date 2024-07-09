package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eqkez0r/gophermart/pkg/authinspector"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/hash"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	RegisterHandlerPath = "/register"
)

var (
	errInvalidFormat = errors.New("invalid format request")
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
		err = storage.NewUser(ctx, newUser)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				logger.Info(err, pgErr)
				if pgErr.Code == "23505" {
					c.Status(http.StatusConflict)
					return
				}
			}
			c.Status(http.StatusInternalServerError)
			return
		}
		t := time.Now()
		formattedTime := t.Format(time.RFC3339)
		err = insp.Auth(ctx, newUser, formattedTime)
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
