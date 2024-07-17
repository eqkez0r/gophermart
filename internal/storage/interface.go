package storage

import (
	"context"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
)

type Storage interface {
	NewUser(context.Context, *obj.User) error
	GetUser(context.Context, string) (*obj.User, error)
	GetLastUserID(context.Context) (uint64, error)
	IsUserExist(context.Context, string) (bool, error)
	NewOrder(context.Context, string, uint64) error
	GetOrdersList(context.Context, uint64) ([]*obj.Order, error)
	GetUnfinishedOrders(context.Context) ([]*obj.Order, error)
	GetBalance(context.Context, uint64) (*obj.AccrualBalance, error)
	NewWithdraw(context.Context, uint64, string, float64) error
	Withdrawals(context.Context, uint64) ([]*obj.Withdraw, error)
	UpdateAccrual(context.Context, uint64, *obj.Accrual) error
	GracefulShutdown() error
}
