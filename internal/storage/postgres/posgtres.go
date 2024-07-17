package postgres

import (
	"context"
	"errors"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/eqkez0r/gophermart/utils/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"log"
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
		order_number VARCHAR(20) UNIQUE NOT NULL, 
		order_customer INTEGER REFERENCES users(user_id) ON DELETE CASCADE NOT NULL, 
		order_accrual DOUBLE PRECISION,
		order_time TIMESTAMP WITH TIME ZONE NOT NULL,
		order_status VARCHAR(10) NOT NULL 
	)`
	queryCreateWithdrawsTable = `CREATE TABLE withdraws(
    	order_customer INTEGER REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    	order_number VARCHAR(20) REFERENCES orders(order_number) UNIQUE NOT NULL,
    	accrual DOUBLE PRECISION NOT NULL,
    	withdraw_time TIMESTAMP WITH TIME ZONE NOT NULL
)`
	queryNewUser              = `INSERT INTO users(login, password, accrual_balance, withdrawal_balance) VALUES ($1, $2, 0, 0)`
	queryGetUser              = `SELECT * FROM users WHERE login = $1`
	queryGetOnlyLogin         = `SELECT login FROM users WHERE login = $1`
	queryGetLastUserID        = "SELECT user_id FROM users ORDER BY user_id DESC LIMIT 1"
	queryGetBalance           = `SELECT (accrual_balance, withdrawal_balance) FROM users WHERE user_id = $1`
	queryUpdateBalance        = `UPDATE users SET accrual_balance = $1, withdrawal_balance = $2 WHERE user_id = $3`
	queryUpdateAccrualBalance = `UPDATE users SET accrual_balance = users.accrual_balance + $1 WHERE user_id = $2`

	queryNewOrder = `INSERT INTO orders(order_number,
                   order_customer,
                   order_time,
                   order_status) VALUES ($1,$2,$3,$4)`
	queryGetOrderList      = `SELECT * FROM orders WHERE order_customer = $1`
	queryUpdateOrderStatus = `UPDATE orders SET order_status = $1 WHERE order_number = $2`
	queryGetNotFinished    = `SELECT order_customer, order_number FROM orders WHERE order_status = 'NEW' OR order_status = 'PROCESSING'`
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
	p.logger.Infof("initial user data %v", usr)
	if err := row.Scan(&usr.UserID, &usr.Login, &usr.Password, &usr.Balance, &usr.Withdraw); err != nil {
		p.logger.Errorf("Database scan user: %s. %v", login, err)
		return nil, err
	}
	p.logger.Infof("Get user data %v", usr)
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

func (p *PostgreSQLStorage) NewOrder(ctx context.Context, number string, userID uint64) error {
	t := time.Now().Format(time.RFC3339)
	log.Println(number, userID, t)
	_, err := p.pool.Exec(ctx, queryNewOrder, number, userID, t, obj.OrderStatusNew)
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
		if err = rows.Scan(&order.UserID, &order.Number); err != nil {
			p.logger.Errorf("Database scan order: %s.", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (p *PostgreSQLStorage) NewWithdraw(ctx context.Context, userID uint64, number string, withdraw float64) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	order := &obj.Order{}
	err = p.pool.QueryRow(ctx, queryGetOrder, number).Scan(&order)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.ErrIsOrderIsNotExist
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
		return e.ErrBalanceIsNotEnough
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

	if _, err = p.pool.Exec(ctx, queryUpdateBalance, balance.Balance, balance.Withdraw); err != nil {
		p.logger.Errorf("Database exec change account balance: %s.", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Errorf("Database commit transaction: %s.", err)
		return err
	}
	defer func() {
		if err != nil {
			p.logger.Info("Rollback transaction")
			tx.Rollback(ctx)
		}
	}()
	return nil
}

func (p *PostgreSQLStorage) Withdrawals(ctx context.Context, userID uint64) ([]*obj.Withdraw, error) {
	withdrawals := make([]*obj.Withdraw, 0)
	rows, err := p.pool.Query(ctx, queryGetNotFinished)
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

func (p *PostgreSQLStorage) UpdateAccrual(ctx context.Context, userID uint64, accrual *obj.Accrual) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, queryUpdateOrderStatus, obj.AccrualStatusToOrderStatus[accrual.Status], accrual.Order)
	if err != nil {
		p.logger.Errorf("Database exec update order status: %s.", err)
		return err
	}

	if accrual.Status == obj.AccrualStatusProcessed && accrual.Accrual != nil {
		_, err = p.pool.Exec(ctx, queryUpdateAccrualBalance, accrual.Accrual, userID)
		if err != nil {
			p.logger.Errorf("Database exec update accrual balance: %s.", userID)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		p.logger.Errorf("Database commit transaction: %s.", err)
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				p.logger.Errorf("Rollback transaction: %s.", err)
			}
		}
	}()
	return nil
}

func (p *PostgreSQLStorage) GracefulShutdown() error {
	p.pool.Close()
	return nil
}
