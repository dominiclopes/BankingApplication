package db

import "errors"

var (
	ErrAccountNotExist     = errors.New("no account exist in db")
	ErrUserNotExist        = errors.New("user does not exist in db")
	ErrTransactionNotExist = errors.New("transactions for the user for not exist in db")
	ErrInsufficientFunds   = errors.New("insufficient funds")
)
