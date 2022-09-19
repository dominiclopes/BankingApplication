package models

import (
	"github.com/dgrijalva/jwt-go"

	"example.com/banking/repositories"
)

type Claims struct {
	ID   string
	Role string
	jwt.StandardClaims
}

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	User repositories.User `json:"user"`
}

type DepositWithdrawAmountRequest struct {
	Amount float32 `json:"amount"`
}

type GetTransactionDetailsRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PingResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	TokenString string `json:"token"`
}

type CreateAccountResponse struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type DepositWithdrawAmountResponse struct {
	Balance float32 `json:"balance"`
}