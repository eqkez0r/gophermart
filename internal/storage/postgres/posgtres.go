package postgres

import (
	"context"
	"github.com/eqkez0r/gophermart/internal/storage"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/retry"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	queryCreateUserTable = `CREATE TABLE Users(
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    password VARCHAR(128)          
)`
	queryCreateOrdersTable  = ``
	queryCreateAccrualTable = ``

	querySelectLastUserIndex = `SELECT id FROM Users ORDER BY id DESC LIMIT 1`
	queryPutNewUser          = `INSERT INTO Users(id, login, password) VALUES $1, $2, $3`
	queryGetUser             = `SELECT * FROM Users WHERE login = $1`
)

type PostgreSQLStorage struct {
	logger *zap.SugaredLogger
	pool   *pgxpool.Pool
}

func New(
	ctx context.Context,
	logger *zap.SugaredLogger,
	uri string,
) (storage.Storage, error) {
	const op = "Initial PostreSQL user storage error: "
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = retry.Retry(logger, 3, func() error {
		if err = pool.Ping(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = retry.Retry(logger, 3, func() error {
		_, err = pool.Exec(ctx, queryCreateUserTable)
		return nil
	})

	return &PostgreSQLStorage{
		logger: logger,
	}, nil
}

func (p PostgreSQLStorage) NewUser(user *obj.User) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLStorage) GetUser(login string) (*obj.User, error) {
	//TODO implement me
	panic("implement me")
}
