package bank

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuidgen "github.com/pborman/uuid"
	"go.uber.org/zap"

	"example.com/banking/db"
)

var secretKey = []byte("I'mGoingToBeAGolangDeveloper")

type Service interface {
	Login(ctx context.Context, lReq LoginRequest) (tokenString string, tokenExpirationTime time.Time, err error)
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

func generateJWT(userID string, role string) (tokenString string, tokenExpirationTime time.Time, err error) {
	tokenExpirationTime = time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		err = fmt.Errorf("error generating token, err: %v", err)
		return
	}
	return
}

func ValidateJWT(tokenString string) (claims *Claims, err error) {
	claims = &Claims{}

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		err = fmt.Errorf("unauthorized, err: %v", err)
		return
	}
	return
}

func (b *bankService) Login(ctx context.Context, u LoginRequest) (tokenString string, tokenExpirationTime time.Time, err error) {
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
	tokenString, tokenExpirationTime, err = generateJWT(user.ID, user.Type)
	if err != nil {
		return
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
	fmt.Printf("Depositing amount: %v in account: %v\n", amount, accId)

	err = b.store.DepositAmount(ctx, accId, userID, amount)
	if err != nil {
		return
	}

	return
}

func (b *bankService) WithdrawAmount(ctx context.Context, accId, userID string, amount float32) (err error) {
	fmt.Printf("Withdrawing amount: %v from account: %v\n", amount, accId)

	err = b.store.WithdrawAmount(ctx, accId, userID, amount)
	if err != nil {
		return
	}
	return
}

func (b *bankService) GetTransactionDetails(ctx context.Context, accId, userID, startDate, endDate string) (transactions []db.Transaction, err error) {
	fmt.Printf("Getting transactions details for account: %v, from %v to %v\n", accId, startDate, endDate)
	allTransactions, err := b.store.GetTransactions(ctx, accId, userID)
	if err != nil {
		return
	}

	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		err = fmt.Errorf("error parsing startdate: %v", startDate)
		return
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		err = fmt.Errorf("error parsing endate %v", endDate)
		return
	}

	transactions = make([]db.Transaction, 0)
	var tDateTimeAsNeeded time.Time
	var errFormat1, errFormat2, errFormat3 error
	for _, t := range allTransactions {
		tDateTimeAsNeeded, errFormat1 = time.Parse("2006-01-02T15:04:05.000Z", t.CreatedAt)
		if errFormat1 != nil {
			tDateTimeAsNeeded, errFormat2 = time.Parse("2006-01-02T15:04:05.00Z", t.CreatedAt)

			if errFormat2 != nil {
				tDateTimeAsNeeded, errFormat3 = time.Parse("2006-01-02T15:04:05.0Z", t.CreatedAt)
			}

			if errFormat3 != nil {
				err = fmt.Errorf("error parsing transaction date: %v", t.CreatedAt)
				return nil, err
			}
		}

		// Transaction time must be greater than or equal to startdate and
		// less than or equal to end date
		if (tDateTimeAsNeeded == startDateTime || tDateTimeAsNeeded.After(startDateTime)) &&
			(tDateTimeAsNeeded == endDateTime || tDateTimeAsNeeded.Before(endDateTime)) {
			transactions = append(transactions, t)
		}
	}
	return
}
