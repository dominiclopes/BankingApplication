package bank

import (
	"context"
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
	GetAccountList(ctx context.Context) (accounts []db.Account, err error)
	GetAccountDetails(ctx context.Context, accId string) (acc db.Account, err error)
	DepositAmount(ctx context.Context, accId string, amount float32) (err error)
	WithdrawAmount(ctx context.Context, accId string, amount float32) (err error)
	GetTransactionDetails(ctx context.Context, accId string, startDate, endDate string) (transactions []db.Transaction, err error)
	DeleteAccount(ctx context.Context, accID string) (err error)
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

func generateJWT(email string, role string) (tokenString string, tokenExpirationTime time.Time, err error) {
	tokenExpirationTime = time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: email,
		Role:  role,
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

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			err = fmt.Errorf("unauthorized, err: %v", err)
			return
		}
		err = fmt.Errorf("bad request, err: %v", err)
		return
	}

	if !token.Valid {
		err = fmt.Errorf("token expired - unauthorized, err: %v", err)
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
	tokenString, tokenExpirationTime, err = generateJWT(user.Email, user.Type)
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

func (b *bankService) GetAccountList(ctx context.Context) (accounts []db.Account, err error) {
	b.logger.Info("Getting the list of accounts in the bank")
	accounts, err = b.store.GetAccountList(ctx)
	return
}

func (b *bankService) GetAccountDetails(ctx context.Context, accId string) (acc db.Account, err error) {
	b.logger.Infof("Getting the customer details for account: %v\n", accId)
	acc, err = b.store.GetAccountDetails(ctx, accId)
	return
}

func (b *bankService) DepositAmount(ctx context.Context, accId string, amount float32) (err error) {
	fmt.Printf("Depositing amount: %v in account: %v\n", amount, accId)

	err = b.store.DepositAmount(ctx, accId, amount)
	if err != nil {
		return
	}

	return
}

func (b *bankService) WithdrawAmount(ctx context.Context, accId string, amount float32) (err error) {
	fmt.Printf("Withdrawing amount: %v from account: %v\n", amount, accId)

	err = b.store.WithdrawAmount(ctx, accId, amount)
	if err != nil {
		return
	}
	return
}

func (b *bankService) GetTransactionDetails(ctx context.Context, accId, startDate, endDate string) (transactions []db.Transaction, err error) {
	fmt.Printf("Getting transactions details for account: %v, from %v to %v\n", accId, startDate, endDate)
	allTransactions, err := b.store.GetTransactions(ctx, accId)
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
	for _, t := range allTransactions {
		tDateTimeAsNeeded, err := time.Parse("2006-01-02T15:04:05.000Z", t.CreatedAt)
		if err != nil {
			err = fmt.Errorf("error parsing transaction date: %v", t.CreatedAt)
			return nil, err
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

func (b *bankService) DeleteAccount(ctx context.Context, accID string) (err error) {
	fmt.Printf("Deleting account: %v\n", accID)
	err = b.store.DeleteAccount(ctx, accID)
	return
}
