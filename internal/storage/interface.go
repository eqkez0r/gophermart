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
	NewOrder(context.Context, string, string) error
	GetOrdersList(context.Context, string) ([]*obj.Order, error)
	GetUnfinishedOrders(context.Context) ([]*obj.Order, error)
	GetBalance(context.Context, string) (*obj.AccrualBalance, error)
	NewWithdraw(context.Context, string, string, float32) error
	Withdrawals(context.Context, string) ([]*obj.Withdraw, error)
	UpdateAccrual(context.Context, uint64, *obj.Accrual) error
	GracefulShutdown() error
}
