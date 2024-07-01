package storage

import obj "github.com/eqkez0r/gophermart/pkg/objects"

type Storage interface {
	NewUser(*obj.User) error
	GetUser(string) (*obj.User, error)
}
