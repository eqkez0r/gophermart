package postgres

import (
	"context"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/retry"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

const (
	queryCreateUserTable = `CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL ,
    password VARCHAR(128) NOT NULL     
)`
	queryCreateOrdersTable = `
	CREATE TABLE orders(
		order_number BIGINT UNIQUE NOT NULL, 
		order_customer INTEGER REFERENCES users(user_id) ON DELETE CASCADE NOT NULL, 
		order_accrual DOUBLE PRECISION,
		order_time TIMESTAMP WITH TIME ZONE NOT NULL,
		order_status VARCHAR(10) NOT NULL 
	)`

	//queryLastIndex = `SELECT id FROM Users ORDER BY id DESC LIMIT 1`
	queryNewUser       = `INSERT INTO users(login, password) VALUES ($1, $2)`
	queryGetUser       = `SELECT * FROM users WHERE login = $1`
	queryGetOnlyLogin  = `SELECT login FROM users WHERE login = $1`
	queryGetLastUserID = "SELECT user_id FROM users ORDER BY user_id DESC LIMIT 1"

	queryNewOrder          = `INSERT INTO orders(order_number, order_customer, order_time, order_status) VALUES ($1,$2,$3,$4)`
	queryGetOrder          = `SELECT * FROM orders WHERE order_number = $1`
	queryChangeOrderStatus = `UPDATE orders SET order_status = $1 WHERE order_number = $2`
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
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *PostgreSQLStorage) NewUser(ctx context.Context, user *obj.User) error {
	p.logger.Infof("user data %v", user)
	_, err := p.pool.Exec(ctx, queryNewUser, user.Login, user.Password)
	if err != nil {
		p.logger.Errorf("Database exec user: %s. %v", user.Login, err)
		return err
	}
	return nil
}

func (p *PostgreSQLStorage) GetUser(ctx context.Context, login string) (*obj.User, error) {
	row := p.pool.QueryRow(ctx, queryGetUser, login)
	usr := &obj.User{}
	if err := row.Scan(&usr.UserID, &usr.Login, &usr.Password); err != nil {
		p.logger.Errorf("Database scan user: %s. %v", login, err)
		return nil, err
	}
	return usr, nil
}

func (p *PostgreSQLStorage) GetLastUserID(ctx context.Context) (uint64, error) {
	row := p.pool.QueryRow(ctx, queryGetLastUserID)
	var userID uint64
	if err := row.Scan(&userID); err != nil {
		p.logger.Errorf("Database scan last user ID: %s. %v", userID, err)
		return 0, err
	}
	return userID, nil
}

func (p *PostgreSQLStorage) NewOrder(ctx context.Context, number uint64, userid uint64) error {
	t := time.Now().Format(time.RFC3339)
	_, err := p.pool.Exec(ctx, queryNewOrder, number, userid, t, obj.OrderStatusNew)
	if err != nil {
		p.logger.Errorf("Database exec order: %s. %v", number, err)
		return err
	}
	return nil
}

func (p *PostgreSQLStorage) GetOrder(ctx context.Context, number uint64) (*obj.Order, error) {
	row := p.pool.QueryRow(ctx, queryGetOrder, number)
	order := &obj.Order{}
	if err := row.Scan(order); err != nil {
		p.logger.Errorf("Database scan order: %s. %v", number, err)
		return nil, err
	}
	return nil, nil
}

func (p *PostgreSQLStorage) GracefulShutdown() {
	p.pool.Close()
}

func (p *PostgreSQLStorage) IsUserExist(ctx context.Context, login string) (bool, error) {
	row := p.pool.QueryRow(ctx, queryGetOnlyLogin, login)
	var dblogin string
	if err := row.Scan(&dblogin); err != nil {
		return false, err
	}
	return true, nil
}

func (p *PostgreSQLStorage) GetUnfinishedOrders(ctx context.Context) ([]*obj.Order, error) {
	orders := make([]*obj.Order, 0)
	rows, err := p.pool.Query(ctx, queryGetNotFinished)
	if err != nil {
		p.logger.Errorf("Database query orders: %s.", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := &obj.Order{}
		if err := rows.Scan(&order.Number, &order.UserID, &order.Accrual, &order.UploadAt, &order.Status); err != nil {
			p.logger.Errorf("Database scan order: %s.", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
