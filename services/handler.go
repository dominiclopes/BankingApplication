package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"example.com/banking/models"
	"example.com/banking/repositories"
)

func Response(rw http.ResponseWriter, status int, response interface{}) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("error encoding response"))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(responseBytes)
}

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	Response(rw, http.StatusOK, models.PingResponse{Message: "pong"})
}

func LoginHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var uAuth models.LoginRequest

		err := json.NewDecoder(req.Body).Decode(&uAuth)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		tokenString, tokenExpirationTime, err := s.Login(uAuth)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: tokenExpirationTime,
		})

		rw.Write([]byte("Successfully logged in"))
	})
}

func CreateAccountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		if claims.Role != "accountant" {
			Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		var u repositories.User

		err = json.NewDecoder(req.Body).Decode(&u)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Validate if the user email and phonenumber is correct
		if u.Email == "" || u.PhoneNumber == "" {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: "BadRequest"})
			return
		}

		acc, err := s.CreateAccount(u)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, acc)
	})
}

func GetAccountsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		if claims.Role != "accountant" {
			Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		accounts, err := s.GetAccountList()
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, accounts)
	})
}

func GetAccountDetailsHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		params := mux.Vars(req)
		accID := params["account_id"]

		acc, err := s.GetAccountDetails(accID)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, acc)
	})
}

func DepositAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var depositAmountRequest models.DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&depositAmountRequest)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		bal, err := s.DepositAmount(accId, depositAmountRequest.Amount)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, bal)

	})
}

func WithdrawAmountHandler(s Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var withdrawAmountRequest models.DepositWithdrawAmountRequest
		err = json.NewDecoder(req.Body).Decode(&withdrawAmountRequest)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		bal, err := s.WithdrawAmount(accId, withdrawAmountRequest.Amount)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, bal)
	})
}

func GetTransactionDetailsHandler(b Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
				return
			}
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		tokenString := cookie.Value

		// Authenticate and verify the authorization
		_, err = ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		params := mux.Vars(req)
		accId := params["account_id"]

		var transactionDetailsRequest models.GetTransactionDetailsRequest
		err = json.NewDecoder(req.Body).Decode(&transactionDetailsRequest)
		if err != nil {
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Validate the formats of the start and end date
		startDateTime, err := time.Parse("2006-01-02", transactionDetailsRequest.StartDate)
		if err != nil {
			err = fmt.Errorf("error parsing startdate: %v", transactionDetailsRequest.StartDate)
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		endDateTime, err := time.Parse("2006-01-02", transactionDetailsRequest.EndDate)
		if err != nil {
			err = fmt.Errorf("error parsing enddate %v", transactionDetailsRequest.EndDate)
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Validate the difference between the days should be between 1-30 days
		if startDateTime == endDateTime || startDateTime.After(endDateTime) {
			err = fmt.Errorf("start date: %v must be less than end date: %v",
				transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		diffStartAndEndDate := endDateTime.Sub(startDateTime)
		if diffStartAndEndDate.Hours()/24 > 30 {
			err = fmt.Errorf("difference between the start date and end date must be less than or equal to 30 days")
			Response(rw, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		transactions, err := b.GetTransactionDetails(accId, transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
		if err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		Response(rw, http.StatusOK, transactions)
	})
}

func DeleteAccountHandler(b Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("token")
		if err != nil {
			Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
			return
		}

		tokenString := cookie.Value

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
			return
		}

		if claims.Role != "accountat" {
			Response(rw, http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			return
		}

		params := mux.Vars(req)
		accID := params["account_id"]

		if err = b.DeleteAccount(accID); err != nil {
			Response(rw, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}

		rw.Write([]byte("Successfully deleted account"))

	})
}
