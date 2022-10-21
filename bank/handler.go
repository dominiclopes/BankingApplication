package bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dominiclopes/BankingApplication/api"
	"github.com/dominiclopes/BankingApplication/db"
	"github.com/dominiclopes/BankingApplication/utils"
)

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	api.APIResponse(rw, http.StatusOK, api.Response{Message: "pong"})
}

func LoginHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var uAuth LoginRequest

		err := json.NewDecoder(req.Body).Decode(&uAuth)
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		// Validate if the login request
		err = uAuth.Validate()
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		logRes, err := s.Login(req.Context(), uAuth)
		if err != nil {
			if err == ErrUnauthorized {
				api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
				return
			}
			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: "Err - Internal Server Error"})
			return
		}

		api.APIResponse(rw, http.StatusOK, logRes)
	})
}

func ReadToken(req *http.Request) (token string, err error) {
	authHeader := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		err = errors.New("malformed header")
		return
	}
	token = authHeader[1]
	return
}

func CreateAccountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil || claims.Role != "accountant" {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		var accReq CreateAccountRequest

		err = json.NewDecoder(req.Body).Decode(&accReq)
		if err != nil {
			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: "Internal Server Error"})
			return
		}

		// Validate if the login request
		err = accReq.Validate()
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		accRes, err := s.CreateAccount(req.Context(), accReq)
		if err != nil {
			if err.Error() == "account exists for the given email" {
				api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: "Err - Account exists for the given email"})
				return
			}
			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: "Err - Internal Server Error - Failure creating user account"})
			return
		}

		api.APIResponse(rw, http.StatusOK, accRes)
	})
}

func GetAccountsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil || claims.Role != "accountant" {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		accounts, err := s.GetAccountList(req.Context())
		if err != nil {
			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.APIResponse(rw, http.StatusOK, accounts)
	})
}

func GetAccountDetailsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		params := mux.Vars(req)
		accID := params["account_id"]

		acc, err := s.GetAccountDetails(req.Context(), accID, claims.UserID)
		if err != nil {
			if err == db.ErrAccountNotExist {
				api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
				return
			}

			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.APIResponse(rw, http.StatusOK, acc)
	})
}

func DepositAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var depositAmountRequest DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&depositAmountRequest)
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		err = s.DepositAmount(req.Context(), accId, claims.UserID, depositAmountRequest.Amount)
		if err != nil {
			if err == db.ErrAccountNotExist {
				api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
				return
			}

			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.APIResponse(rw, http.StatusOK, api.Response{Message: fmt.Sprintf("Successfully credited account with amount %v", depositAmountRequest.Amount)})

	})
}

func WithdrawAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}
		params := mux.Vars(req)
		accId := params["account_id"]

		var withdrawAmountRequest DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&withdrawAmountRequest)
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		err = s.WithdrawAmount(req.Context(), accId, claims.UserID, withdrawAmountRequest.Amount)
		if err != nil {
			if err == db.ErrInsufficientFunds {
				api.APIResponse(rw, http.StatusOK, api.Response{Message: err.Error()})
				return
			}

			if err == db.ErrAccountNotExist {
				api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
				return
			}

			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.APIResponse(rw, http.StatusOK, api.Response{Message: fmt.Sprintf("Successfully debited account with amount %v", withdrawAmountRequest.Amount)})
	})
}

func GetTransactionDetailsHandler(b Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Authenticate the user
		token, err := ReadToken(req)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		// Authorize the user
		claims, err := utils.Decode(token)
		if err != nil {
			api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var transactionDetailsRequest GetTransactionDetailsRequest
		err = json.NewDecoder(req.Body).Decode(&transactionDetailsRequest)
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		err = transactionDetailsRequest.Validate()
		if err != nil {
			api.APIResponse(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		transactions, err := b.GetTransactionDetails(req.Context(), accId, claims.UserID, transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
		if err != nil {
			if err == db.ErrAccountNotExist {
				api.APIResponse(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
				return
			}

			api.APIResponse(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		if len(transactions) == 0 {
			api.APIResponse(rw, http.StatusOK, api.Response{Message: "No transactions found within the date range"})
			return
		}
		api.APIResponse(rw, http.StatusOK, transactions)
	})
}
