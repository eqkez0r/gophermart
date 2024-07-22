package e

import "errors"

var (
	ErrBalanceIsNotEnough              = errors.New("balance is not enough")
	ErrIsOrderIsNotExist               = errors.New("order is not exist")
	ErrIsOrderExistWithAnotherCustomer = errors.New("order is exist with a another customer")
)
