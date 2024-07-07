package postgres

import (
	"context"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/retry"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	queryCreateUserTable = `CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    password VARCHAR(128)          
)`
	queryCreateOrdersTable = `
	CREATE TABLE orders(
		order_id SERIAL PRIMARY KEY,
		order_customer INTEGER REFERENCES users(user_id),
		order_sum DOUBLE PRECISION,
		order_accrual DOUBLE PRECISION,
		order_time timestamp,
		order_status INTEGER
	)`

	//queryLastIndex = `SELECT id FROM Users ORDER BY id DESC LIMIT 1`
	queryNewUser = `INSERT INTO users(login, password) VALUES ($1, $2)`
	queryGetUser = `SELECT * FROM users WHERE login = $1`

	queryNewOrder = `INSERT INTO orders(order_customer,
                   order_sum,
                   order_time,
                   order_status) VALUES ($1,$2,$3,
                                     (
                                     	CASE $4
                                     		WHEN 'NEW' THEN 0
                                     		WHEN 'PROCESSING' THEN 1
                                     		WHEN 'INVALID' THEN 2
                                     		WHEN 'PROCESSED' THEN 3
                                         END
                                     ))`
	queryGetOrder          = `SELECT * FROM orders WHERE order_id = $1`
	queryChangeOrderStatus = `UPDATE orders SET order_status = $1 WHERE order_id = $2`
	queryGetNotFinished    = `SELECT * FROM orders WHERE order_status = 'NEW' OR order_status = 'PROCESSING'`
)

type PostgreSQLStorage struct {
	logger *zap.SugaredLogger
	pool   *pgxpool.Pool
}

func New(
	ctx context.Context,
	logger *zap.SugaredLogger,
	uri string,
) (*PostgreSQLStorage, error) {
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

	err = retry.Retry(logger, 3, func() error {
		_, err = pool.Exec(ctx, queryCreateOrdersTable)
		return nil
	})

	return &PostgreSQLStorage{
		logger: logger,
	}, nil
}

func (p *PostgreSQLStorage) NewUser(ctx context.Context, user *obj.User) error {
	_, err := p.pool.Exec(ctx, queryNewUser, user.Login, user.Password)
	if err != nil {
		p.logger.Errorf("Database scan user: %s. %v", user.Login, err)
		return err
	}
	return nil
}

func (p *PostgreSQLStorage) GetUser(ctx context.Context, login string) (*obj.User, error) {
	row := p.pool.QueryRow(ctx, queryGetUser, login)
	usr := &obj.User{}
	if err := row.Scan(usr); err != nil {
		p.logger.Errorf("Database scan user: %s. %v", login, err)
		return nil, err
	}
	return usr, nil
}

func (p *PostgreSQLStorage) NewOrder(ctx context.Context, order *obj.Order) error {
	return nil
}

func (p *PostgreSQLStorage) GetOrder(ctx context.Context, orderID string) (*obj.Order, error) {
	return nil, nil
}

func (p *PostgreSQLStorage) Close() {
	p.pool.Close()
}
