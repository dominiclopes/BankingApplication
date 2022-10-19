package bank

import (
	"github.com/dgrijalva/jwt-go"
)

type PingResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID string
	Role   string
	jwt.StandardClaims
}

type CreateAccountRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type CreateAccountResponse struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID string `json:"account_id"`
}

type DepositWithdrawAmountRequest struct {
	Amount float32 `json:"amount"`
}

type GetTransactionDetailsRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
