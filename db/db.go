package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ctxKey int

const (
	dbKey          ctxKey = 0
	defaultTimeout        = 1 * time.Second
)

type Storer interface {
	GetUserByEmailAndPassword(ctx context.Context, email string, password string) (u User, err error)
	CreateAccount(ctx context.Context, u User, acc Account) (err error)
	GetAccountList(ctx context.Context) (accounts []UserAccountDetails, err error)
	GetAccountDetails(ctx context.Context, accID, userID string) (acc UserAccountDetails, err error)
	AddTransaction(ctx context.Context, t Transaction) (err error)
	DepositAmount(ctx context.Context, accID, userID string, amount float32) (err error)
	WithdrawAmount(ctx context.Context, accID, userID string, amount float32) (err error)
	GetTransactions(ctx context.Context, accID, userID, startDate, endDate string) (transactions []Transaction, err error)
}

type store struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewStorer(d *sqlx.DB, l *zap.SugaredLogger) Storer {
	return &store{
		db:     d,
		logger: l,
	}
}

func newContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, dbKey, tx)
}

func WithTimeout(ctx context.Context, timeout time.Duration, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return op(ctxWithTimeout)
}

func WithDefaultTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	return WithTimeout(ctx, defaultTimeout, op)
}
