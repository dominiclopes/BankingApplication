package bank

import (
	"context"
	"errors"

	uuidgen "github.com/pborman/uuid"
	"go.uber.org/zap"

	"github.com/dominiclopes/BankingApplication/db"
	"github.com/dominiclopes/BankingApplication/utils"
)

type Service interface {
	Login(ctx context.Context, lReq LoginRequest) (loginRes LoginResponse, err error)
	CreateAccount(ctx context.Context, accReq CreateAccountRequest) (accRes CreateAccountResponse, err error)
	GetAccountList(ctx context.Context) (accounts []db.UserAccountDetails, err error)
	GetAccountDetails(ctx context.Context, accId, userID string) (acc db.UserAccountDetails, err error)
	DepositAmount(ctx context.Context, accId, userID string, amount float32) (err error)
	WithdrawAmount(ctx context.Context, accId, userID string, amount float32) (err error)
	GetTransactionDetails(ctx context.Context, accId, userID string, startDate, endDate string) (transactions []db.Transaction, err error)
}

type bankService struct {
	store  db.Storer
	logger *zap.SugaredLogger
}

func NewBankService(s db.Storer, l *zap.SugaredLogger) Service {
	return &bankService{
		store:  s,
		logger: l,
	}
}

func (b *bankService) Login(ctx context.Context, u LoginRequest) (loginRes LoginResponse, err error) {
	// Verify if user is present
	user, err := b.store.GetUserByEmailAndPassword(ctx, u.Email, u.Password)
	if err != nil {
		if err == db.ErrUserNotExist {
			err = ErrUnauthorized
			return
		}
		return
	}

	// Create the JWT token
	tokenString, err := utils.Encode(user.ID, user.Type)
	if err != nil {
		return
	}

	loginRes = LoginResponse{
		Token: tokenString,
	}

	return
}

func (b *bankService) CreateAccount(ctx context.Context, accReq CreateAccountRequest) (accRes CreateAccountResponse, err error) {
	b.logger.Infof("Creating an account for user email: %v, phone number: %v\n", accReq.Email, accReq.PhoneNumber)

	// Create the user ID, password and update the balance
	u := db.User{
		Email:       accReq.Email,
		PhoneNumber: accReq.PhoneNumber,
		Password:    uuidgen.New(),
		Type:        "customer",
	}

	acc := db.Account{
		ID:      uuidgen.New(),
		Balance: 0.0,
		UserID:  u.ID,
	}

	// Save the user in the bank
	err = b.store.CreateAccount(ctx, u, acc)
	if err != nil {
		b.logger.Errorf("Err creating user account: %v", err.Error())
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			err = errors.New("account exists for the given email")
			return
		}
		return
	}

	// Create the response
	accRes = CreateAccountResponse{
		Email:     u.Email,
		Password:  u.Password,
		AccountID: acc.ID,
	}

	b.logger.Infof("Created account with details: %v. Opening balance: %v\n", accRes, acc.Balance)
	return
}

func (b *bankService) GetAccountList(ctx context.Context) (accounts []db.UserAccountDetails, err error) {
	b.logger.Info("Getting the list of accounts in the bank")
	accounts, err = b.store.GetAccountList(ctx)
	return
}

func (b *bankService) GetAccountDetails(ctx context.Context, accId, userID string) (acc db.UserAccountDetails, err error) {
	b.logger.Infof("Getting the customer details for account: %v\n", accId)
	acc, err = b.store.GetAccountDetails(ctx, accId, userID)
	return
}

func (b *bankService) DepositAmount(ctx context.Context, accId, userID string, amount float32) (err error) {
	b.logger.Infof("Depositing amount: %v in account: %v\n", amount, accId)

	err = b.store.DepositAmount(ctx, accId, userID, amount)
	if err != nil {
		return
	}

	return
}

func (b *bankService) WithdrawAmount(ctx context.Context, accId, userID string, amount float32) (err error) {
	b.logger.Infof("Withdrawing amount: %v from account: %v\n", amount, accId)

	err = b.store.WithdrawAmount(ctx, accId, userID, amount)
	if err != nil {
		return
	}
	return
}

func (b *bankService) GetTransactionDetails(ctx context.Context, accId, userID, startDate, endDate string) (transactions []db.Transaction, err error) {
	b.logger.Infof("Getting transactions details for account: %v, from %v to %v\n", accId, startDate, endDate)

	transactions, err = b.store.GetTransactions(ctx, accId, userID, startDate, endDate)
	if err != nil {
		return
	}
	return
}
