package handlers

import (
	"context"
	"errors"
	"fmt"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/eqkez0r/gophermart/pkg/jwt"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/hash"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"net/http"
)

const (
	RegisterHandlerPath = "/register"
)

var (
	errInvalidFormat = errors.New("invalid format request")
)

type NewUserProvider interface {
	NewUser(context.Context, *obj.User) error
	GetLastUserID(context.Context) (uint64, error)
}

func RegisterHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	storage NewUserProvider,
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
					logger.Error(e.Wrap(op, pgErr))
					c.Status(http.StatusConflict)
					return
				}
			}
			c.Status(http.StatusInternalServerError)
			return
		}

		userid, err := storage.GetLastUserID(ctx)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		token, err := jwt.CreateJWT(newUser.Login, userid)
		if err != nil {
			logger.Error(e.Wrap(op, err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Header("Authorization", token)
		c.Status(http.StatusOK)
	}
}
