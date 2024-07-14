package postgres

import (
	"context"
	"errors"
	"github.com/eqkez0r/gophermart/internal/storage"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

const (
	queryCreateUserTable = `CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL ,
    password VARCHAR(128) NOT NULL,
    accrual_balance DOUBLE PRECISION NOT NULL,
    withdrawal_balance DOUBLE PRECISION NOT NULL    
)`
	queryCreateOrdersTable = `
	CREATE TABLE orders(
		order_number BIGINT UNIQUE NOT NULL, 
		order_customer INTEGER REFERENCES users(user_id) ON DELETE CASCADE NOT NULL, 
		order_accrual DOUBLE PRECISION,
		order_time TIMESTAMP WITH TIME ZONE NOT NULL,
		order_status VARCHAR(10) NOT NULL 
	)`
	queryCreateWithdrawsTable = `
	CREATE TABLE withdrawals(
	    order_customer INTEGER REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
	    order_number BIGINT REFERENCES orders(order_number) ON DELETE CASCADE NOT NULL,
	    sum DOUBLE PRECISION REFERENCES orders(order_accrual) ON DELETE CASCADE NOT NULL,
	    withdraw_time TIMESTAMP WITH TIME ZONE NOT NULL
	)
`

	queryNewUser       = `INSERT INTO users(login, password, accrual_balance, withdraw) VALUES ($1, $2, 0, 0)`
	queryGetUser       = `SELECT * FROM users WHERE login = $1`
	queryGetOnlyLogin  = `SELECT login FROM users WHERE login = $1`
	queryGetLastUserID = "SELECT user_id FROM users ORDER BY user_id DESC LIMIT 1"

	queryGetBalance    = `SELECT (accrual_balance, withdrawal_balance) FROM users WHERE user_id = $1`
	queryChangeBalance = `UPDATE users SET accrual_balance = $1, withdrawal_balance = $2 WHERE user_id = $3`

	queryNewOrder = `INSERT INTO orders(order_number,
                   order_customer,
                   order_time,
                   order_status) VALUES ($1,$2,$3,$4)`
	queryGetOrderList      = `SELECT * FROM orders WHERE order_customer = $1`
	queryChangeOrderStatus = `UPDATE orders SET order_status = $1 WHERE order_number = $2`
	queryGetNotFinished    = `SELECT * FROM orders WHERE order_status = 'NEW' OR order_status = 'PROCESSING'`
	queryGetOrder          = `SELECT (order_number) FROM orders WHERE order_number = $1`

	queryNewWithdraw     = `INSERT INTO withdrawals(order_customer, order_number, sum, withdraw_time) VALUES ($1, $2, $3, $4)`
	queryGetWithdrawList = `SELECT * FROM withdrawals WHERE order_customer = $1`
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

	err = retry.Retry(logger, 3, func() error {
		_, err = pool.Exec(ctx, queryCreateWithdrawsTable)
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

func (p *PostgreSQLStorage) GetOrdersList(ctx context.Context, userID uint64) ([]*obj.Order, error) {
	orders := make([]*obj.Order, 0)
	rows, err := p.pool.Query(ctx, queryGetOrderList, int64(userID))
	if err != nil {
		p.logger.Errorf("Database query orders list: %s. %v", userID, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		order := &obj.Order{}
		if err = rows.Scan(&order.Number, &order.UserID, &order.Accrual, &order.UploadAt, &order.Status); err != nil {
			p.logger.Errorf("Database query orders list: %s. %v", userID, err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (p *PostgreSQLStorage) GetBalance(ctx context.Context, userID uint64) (*obj.AccrualBalance, error) {
	balance := &obj.AccrualBalance{}
	row := p.pool.QueryRow(ctx, queryGetBalance, int64(userID))
	if err := row.Scan(&balance.Balance, &balance.Withdraw); err != nil {
		p.logger.Errorf("Database query account balance: %s. %v", userID, err)
		return nil, err
	}
	return balance, nil
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

func (p *PostgreSQLStorage) NewWithdraw(ctx context.Context, userID, orderID uint64, withdraw float64) error {
	order := &obj.Order{}
	err := p.pool.QueryRow(ctx, queryGetOrder, orderID).Scan(&order)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrIsOrderIsNotExist
		}
		return err
	}

	balance, err := p.GetBalance(ctx, userID)
	if err != nil {
		p.logger.Errorf("Database query account balance: %s.", err)
		return err
	}

	if balance.Balance < withdraw {
		p.logger.Errorf("Not enough balance for user: %s.", userID)
		return storage.ErrBalanceIsNotEnough
	}

	p.logger.Infof("Withdrawing user: %s. Before balance: %d, NewWithdraw: %d.",
		userID, balance.Balance, balance.Withdraw)
	balance.Balance -= withdraw
	balance.Withdraw += withdraw
	p.logger.Infof("Withdrawing user: %s. After balance: %d, NewWithdraw: %d.",
		userID, balance.Balance, balance.Withdraw)

	if _, err = p.pool.Exec(ctx, queryNewWithdraw,
		userID, balance.Balance, balance.Withdraw, time.Now().Format(time.RFC3339)); err != nil {
		p.logger.Errorf("Database exec new withdraw: %s.", userID)
		return err
	}

	if _, err = p.pool.Exec(ctx, queryChangeBalance, balance.Balance, balance.Withdraw); err != nil {
		p.logger.Errorf("Database exec change account balance: %s.", err)
		return err
	}

	return nil
}

func (p *PostgreSQLStorage) Withdrawals(ctx context.Context, userID uint64) ([]*obj.Withdraw, error) {
	withdrawals := make([]*obj.Withdraw, 0)
	rows, err := p.pool.Query(ctx, queryGetNotFinished, int64(userID))
	if err != nil {
		p.logger.Errorf("Database query orders: %s.", err)
		return nil, err
	}
	for rows.Next() {
		withdraw := &obj.Withdraw{}
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			p.logger.Errorf("Database scan orders: %s.", err)
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}
	return withdrawals, nil
}

func (p *PostgreSQLStorage) GracefulShutdown() error {
	p.pool.Close()
	return nil
}
