package db

import (
	"context"
	"database/sql"

	uuidgen "github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const (
	createUserQuery                = `INSERT INTO users(email, phone_number, password, type) VALUES ($1, $2, crypt($3, gen_salt('bf')), $4) returning id`
	getUserByEmailAndPasswordQuery = `SELECT * FROM users WHERE email=$1 and password=crypt($2, password)`
	deleteUserByIDQuery            = `DELETE FROM users WHERE id=$1`

	createAccountQuery               = `INSERT INTO accounts(id, balance, user_id) VALUES ($1, $2, $3)`
	listAccountsQuery                = `SELECT accounts.id, accounts.balance, users.email, users.phone_number from accounts inner join users on accounts.user_id=users.id`
	getAccountByAccIDQuery           = `SELECT accounts.id, accounts.balance, users.email, users.phone_number from accounts inner join users on accounts.user_id=users.id where accounts.id=$1 and accounts.user_id=$2`
	updateAccountBalanceByAccIDQuery = `UPDATE accounts SET balance=$1 WHERE id=$2`
	deleteAccountByIDQuery           = `DELETE FROM accounts WHERE id=$1`

	createTransactionQuery      = `INSERT INTO transactions(id, type, amount, balance, account_id) VALUES ($1, $2, $3, $4, $5)`
	getTransactionsByAccIDQuery = `SELECT * FROM transactions WHERE account_id=$1 and created_at BETWEEN $2 and $3`
)

type User struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    string `json:"password" db:"password"`
	Type        string `json:"-" db:"type"`
}

type Account struct {
	ID      string  `json:"account_id" db:"id"`
	Balance float32 `json:"balance" db:"balance"`
	UserID  string  `json:"-" db:"user_id"`
}

type UserAccountDetails struct {
	Account
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

type Transaction struct {
	ID        string  `json:"-" db:"id"`
	Type      string  `json:"type" db:"type"`
	Amount    float32 `json:"amount" db:"amount"`
	Balance   float32 `json:"balance" db:"balance"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	AccountID string  `json:"-" db:"account_id"`
}

func (s *store) GetUserByEmailAndPassword(ctx context.Context, email string, password string) (u User, err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		err = s.db.GetContext(ctx, &u, getUserByEmailAndPasswordQuery, email, password)
		return err
	})

	if err == sql.ErrNoRows {
		return u, ErrUserNotExist
	}

	return
}

func (s *store) CreateAccount(ctx context.Context, u User, acc Account) (err error) {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errors.WithStack(e)
				return
			}
		}
		tx.Commit()
	}()

	ctxWithTx := newContext(ctx, tx)
	err = WithDefaultTimeout(ctxWithTx, func(ctx context.Context) error {
		var user_id int64

		// Create user
		if err := s.db.GetContext(ctx, &user_id, createUserQuery, u.Email, u.PhoneNumber, u.Password, u.Type); err != nil {
			return err
		}

		// Create user account
		if _, err := s.db.Exec(createAccountQuery, acc.ID, acc.Balance, user_id); err != nil {
			return err
		}

		return nil
	})
	return
}

func (s *store) GetAccountList(ctx context.Context) (accounts []UserAccountDetails, err error) {

	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		return s.db.SelectContext(ctx, &accounts, listAccountsQuery)
	})

	if err == sql.ErrNoRows {
		return accounts, ErrAccountNotExist
	}

	return
}

func (s *store) GetAccountDetails(ctx context.Context, accID, userID string) (acc UserAccountDetails, err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		return s.db.GetContext(ctx, &acc, getAccountByAccIDQuery, accID, userID)
	})

	if err == sql.ErrNoRows {
		return acc, ErrAccountNotExist
	}

	return
}

func (s *store) AddTransaction(ctx context.Context, t Transaction) (err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		_, err = s.db.Exec(createTransactionQuery, t.ID, t.Type, t.Amount, t.Balance, t.AccountID)
		return err
	})

	return
}

func (s *store) DepositAmount(ctx context.Context, accID, userID string, amount float32) (err error) {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errors.WithStack(e)
				return
			}
		}
		tx.Commit()
	}()

	ctxWithTx := newContext(ctx, tx)
	err = WithDefaultTimeout(ctxWithTx, func(ctx context.Context) error {
		// get the account details
		acc, err := s.GetAccountDetails(ctx, accID, userID)
		if err != nil {
			return err
		}

		// update the user balance
		balance := acc.Balance + amount

		// update user details
		if _, err = s.db.Exec(updateAccountBalanceByAccIDQuery, balance, accID); err != nil {
			return err
		}

		// Add a transaction
		t := Transaction{
			ID:        uuidgen.New(),
			Type:      "Credit",
			Amount:    amount,
			Balance:   balance,
			AccountID: accID,
		}
		if err = s.AddTransaction(ctx, t); err != nil {
			return err
		}

		s.logger.Infof("Credited amount: %v, in account: %v. Balance: %v\n", amount, accID, balance)
		return err
	})
	return
}

func (s *store) WithdrawAmount(ctx context.Context, accID, userID string, amount float32) (err error) {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errors.WithStack(e)
				return
			}

			tx.Commit()
		}
	}()

	ctxWithTx := newContext(ctx, tx)
	err = WithDefaultTimeout(ctxWithTx, func(ctx context.Context) error {
		acc, err := s.GetAccountDetails(ctx, accID, userID)
		if err != nil {
			return err
		}

		// verify if amount can be debited
		if acc.Balance < amount {
			s.logger.Errorf("amount %v cannot be debited from account %v. insufficient funds: %v\n",
				amount, accID, acc.Balance)
			return ErrInsufficientFunds
		}

		// update the user balance
		balance := acc.Balance - amount

		// update user details
		if _, err = s.db.Exec(updateAccountBalanceByAccIDQuery, balance, accID); err != nil {
			return err
		}

		// sdd a transaction
		t := Transaction{
			ID:        uuidgen.New(),
			Type:      "Debit",
			Amount:    amount,
			Balance:   balance,
			AccountID: accID,
		}
		if err = s.AddTransaction(ctx, t); err != nil {
			return err
		}

		s.logger.Infof("Debited amount: %v, from account: %v. Balance: %v\n", amount, accID, balance)
		return err
	})
	return
}

func (s *store) GetTransactions(ctx context.Context, accID, userID, startDate, endDate string) (transactions []Transaction, err error) {

	if _, err = s.GetAccountDetails(ctx, accID, userID); err != nil {
		return
	}

	err = s.db.SelectContext(ctx, &transactions, getTransactionsByAccIDQuery, accID, startDate, endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return transactions, ErrTransactionNotExist
		}
		return
	}

	s.logger.Info("Transactions details:", transactions)
	return
}
