package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"example.com/banking/api"
)

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	api.Success(rw, http.StatusOK, api.Response{Message: "pong"})
}

func LoginHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var uAuth LoginRequest

		err := json.NewDecoder(req.Body).Decode(&uAuth)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		if uAuth.Email == "" || uAuth.Password == "" {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: "Bad Request"})
			return
		}

		tokenString, tokenExpirationTime, err := s.Login(req.Context(), uAuth)
		if err != nil {
			if err == ErrUnauthorized {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: tokenExpirationTime,
		})

		api.Success(rw, http.StatusOK, api.Response{Message: "Successfully logged in"})
	})
}

func CreateAccountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		if claims.Role != "accountant" {
			api.Error(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		var accReq CreateAccountRequest

		err = json.NewDecoder(req.Body).Decode(&accReq)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		// Validate if the user email and phonenumber is correct
		if accReq.Email == "" || accReq.PhoneNumber == "" {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: "BadRequest"})
			return
		}

		accRes, err := s.CreateAccount(req.Context(), accReq)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, accRes)
	})
}

func GetAccountsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		if claims.Role != "accountant" {
			api.Error(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		accounts, err := s.GetAccountList(req.Context())
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, accounts)
	})
}

func GetAccountDetailsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		params := mux.Vars(req)
		accID := params["account_id"]

		acc, err := s.GetAccountDetails(req.Context(), accID)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, acc)
	})
}

func DepositAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var depositAmountRequest DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&depositAmountRequest)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		err = s.DepositAmount(req.Context(), accId, depositAmountRequest.Amount)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, api.Response{Message: fmt.Sprintf("Successfully credited account with amount %v", depositAmountRequest.Amount)})

	})
}

func WithdrawAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var withdrawAmountRequest DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&withdrawAmountRequest)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		err = s.WithdrawAmount(req.Context(), accId, withdrawAmountRequest.Amount)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, api.Response{Message: fmt.Sprintf("Successfully debited account with amount %v", withdrawAmountRequest.Amount)})
	})
}

func GetTransactionDetailsHandler(b Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
				return
			}
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var transactionDetailsRequest GetTransactionDetailsRequest
		err = json.NewDecoder(req.Body).Decode(&transactionDetailsRequest)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		// Validate the formats of the start and end date
		startDateTime, err := time.Parse("2006-01-02", transactionDetailsRequest.StartDate)
		if err != nil {
			err = fmt.Errorf("error parsing startdate: %v", transactionDetailsRequest.StartDate)
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		endDateTime, err := time.Parse("2006-01-02", transactionDetailsRequest.EndDate)
		if err != nil {
			err = fmt.Errorf("error parsing enddate %v", transactionDetailsRequest.EndDate)
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		// Validate the difference between the days should be between 1-30 days
		if startDateTime == endDateTime || startDateTime.After(endDateTime) {
			err = fmt.Errorf("start date: %v must be less than end date: %v",
				transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}
		diffStartAndEndDate := endDateTime.Sub(startDateTime)
		if diffStartAndEndDate.Hours()/24 > 30 {
			err = fmt.Errorf("difference between the start date and end date must be less than or equal to 30 days")
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		transactions, err := b.GetTransactionDetails(req.Context(), accId, transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, transactions)
	})
}

func DeleteAccountHandler(b Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		tokenString := cookie.Value

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			api.Error(rw, http.StatusUnauthorized, api.Response{Message: err.Error()})
			return
		}

		if claims.Role != "accountant" {
			api.Error(rw, http.StatusUnauthorized, api.Response{Message: "Unauthorized"})
			return
		}

		params := mux.Vars(req)
		accID := params["account_id"]

		if err = b.DeleteAccount(req.Context(), accID); err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, api.Response{Message: "Successfully deleted account"})
	})
}
