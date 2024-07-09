package storage

import (
	"context"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
)

type Storage interface {
	NewUser(context.Context, *obj.User) error
	GetUser(context.Context, string) (*obj.User, error)
	IsUserExist(context.Context, string) (bool, error)
	NewOrder(context.Context, uint64) error
	GetOrder(context.Context, string) (*obj.Order, error)
	Close()
}
