package e

import "errors"

var (
	ErrBalanceIsNotEnough = errors.New("balance is not enough")
	ErrIsOrderIsNotExist  = errors.New("order is not exist")
)
