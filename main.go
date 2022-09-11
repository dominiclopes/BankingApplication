package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	uuidgen "github.com/pborman/uuid"
)

type User struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Password    string  `json:"password"`
	Balance     float32 `json:"balance"`
}

type Transaction struct {
	ID              string  `json:"-"`
	TransactionType string  `json:"type"`
	Amount          float32 `json:"amount"`
	Balance         float32 `json:"balance"`
	DateTime        string  `json:"time"`
	UserID          string  `json:"-"`
}

type BankService interface {
	CreateCustomer(user User) (accountId string, password string)
	GetCustomerList() (customers []User)
	GetCustomerDetails(accountId string) (err error)
	DepositAmount(accountId string, amount float32) (balance float32, err error)
	WithdrawAmount(accountId string, amount float32) (balance float32, err error)
	GetTransactionDetails(accountId string, startDate, endDate string) (transactions []Transaction, err error)
}

type bank struct {
	users        map[string]User
	transactions map[string]Transaction
}

func (b bank) CreateAccount(user User) User {
	fmt.Printf("Creating an account for user email: %v, phone number: %v\n", user.Email, user.PhoneNumber)

	// Create the user ID
	user.ID = uuidgen.New()
	// Create the user password
	user.Password = uuidgen.New()
	// Set the user account balance
	user.Balance = 0.0

	// Save the user in the bank
	b.users[user.ID] = user

	fmt.Printf("Created account %v with password %v for user email: %v, phone number: %v. Opening balance: %v\n",
		user.ID, user.Password, user.Email, user.Email, user.Balance)

	return user
}

func (b bank) GetCustomerList() (customers []User) {
	fmt.Println("Getting the customer list")

	customers = make([]User, 0)
	for _, u := range b.users {
		customers = append(customers, u)
	}

	fmt.Printf("The customer list is %v\n", customers)
	return customers
}

func (b bank) GetCustomerDetails(accountId string) (userDetails User, err error) {
	fmt.Printf("Getting the customer details for account: %v\n", accountId)
	userDetails, userPresent := b.users[accountId]
	if !userPresent {
		err = fmt.Errorf("User with account: %v not present", accountId)
		return userDetails, err
	}
	fmt.Printf("Customer details for account: %v are %v\n", accountId, userDetails)
	return userDetails, err
}

func (b bank) DepositAmount(accountId string, amount float32) (balance float32, err error) {
	fmt.Printf("Depositing amount: %v in account: %v\n", amount, accountId)
	userDetails, userPresent := b.users[accountId]
	if !userPresent {
		err = fmt.Errorf("User with account: %v not present", accountId)
		return 0.0, err
	}

	// Update the user balance
	userDetails.Balance = userDetails.Balance + amount

	// Create a transaction
	transaction := Transaction{
		ID:              uuidgen.New(),
		UserID:          userDetails.ID,
		Amount:          amount,
		Balance:         userDetails.Balance,
		TransactionType: "Credit",
		DateTime:        time.Now().Format("02-01-2006 15:04:05.000"),
	}
	fmt.Println("Reporting transaction:", transaction)

	// Update the bank
	b.users[accountId] = userDetails
	b.transactions[transaction.ID] = transaction

	fmt.Println("Bank transactions:", b.transactions)
	fmt.Printf("Credited amount: %v, in account: %v. Balance: %v\n", amount, accountId, userDetails.Balance)
	return userDetails.Balance, nil
}

func (b bank) WithdrawAmount(accountId string, amount float32) (balance float32, err error) {
	fmt.Printf("Withdrawing amount: %v from account: %v\n", amount, accountId)
	userDetails, userPresent := b.users[accountId]
	if !userPresent {
		err = fmt.Errorf("User with account: %v not present", accountId)
		return 0.0, err
	}

	// Verify if amount can be debited
	if userDetails.Balance < amount {
		err = fmt.Errorf("Amount %v cannot be debited from account %v. Balance in account: %v",
			amount, accountId, userDetails.Balance)
		return 0.0, err
	}

	// Update the user balance
	userDetails.Balance = userDetails.Balance - amount

	// Create a transaction
	transaction := Transaction{
		ID:              uuidgen.New(),
		UserID:          userDetails.ID,
		Amount:          amount,
		Balance:         userDetails.Balance,
		TransactionType: "Debit",
		DateTime:        time.Now().Format("02-01-2006 15:04:05.000"),
	}
	fmt.Println("Reporting transaction:", transaction)

	b.users[accountId] = userDetails
	b.transactions[transaction.ID] = transaction

	fmt.Println("Bank transactions:", b.transactions)
	fmt.Printf("Debited amount: %v, from account: %v. Balance: %v\n", amount, accountId, userDetails.Balance)
	return userDetails.Balance, nil
}

func (b bank) GetTransactionDetails(accountId string, startDate, endDate string) (transactions []Transaction, err error) {
	fmt.Printf("Getting transactions for account: %v, from %v to %v\n", accountId, startDate, endDate)
	_, userPresent := b.users[accountId]
	if !userPresent {
		err = fmt.Errorf("User with account: %v not present", accountId)
		return nil, err
	}

	startDateTime, err := time.Parse("02-01-2006", startDate)
	if err != nil {
		err = fmt.Errorf("Error parsing startdate: %v", startDate)
		return nil, err
	}
	endDateTime, err := time.Parse("02-01-2006", endDate)
	if err != nil {
		err = fmt.Errorf("Error parsing startdate %v", endDate)
		return nil, err
	}
	if startDateTime.After(endDateTime) || startDateTime == endDateTime {
		err = fmt.Errorf("Start Date: %v must be less than end date: %v", startDate, endDate)
		return nil, err
	}

	transactions = make([]Transaction, 0)
	for _, t := range b.transactions {
		if t.UserID == accountId {
			tDateTime, err := time.Parse("02-01-2006 15:04:05.000", t.DateTime)
			if err != nil {
				err = fmt.Errorf("Error parsing transaction date: %v", t.DateTime)
				return nil, err
			}

			tDateTimeAsNeeded, err := time.Parse("02-01-2006",
				fmt.Sprintf("%02d-%02d-%d", tDateTime.Day(), tDateTime.Month(), tDateTime.Year()))
			if err != nil {
				err = fmt.Errorf("Error parsing transaction date: %v", tDateTime)
				return nil, err
			}
			fmt.Printf("Transaction time %v, startDate: %v, endDate: %v\n", tDateTimeAsNeeded, startDateTime, endDateTime)
			fmt.Printf("Is transactiondate same as start date?: %v\nIs transactiondate greater than start date?: %v\nIs transactiondate same as end date?: %v\nIs transactiondate lesser than end date?: %v\n",
				tDateTimeAsNeeded == startDateTime, tDateTimeAsNeeded.After(startDateTime),
				tDateTimeAsNeeded == endDateTime, tDateTimeAsNeeded.Before(endDateTime))

			// Transaction time must be greater than or equal to startdate and
			// less than or equal to end date
			if (tDateTimeAsNeeded == startDateTime || tDateTimeAsNeeded.After(startDateTime)) &&
				(tDateTimeAsNeeded == endDateTime || tDateTimeAsNeeded.Before(endDateTime)) {
				transactions = append(transactions, t)
				fmt.Println("Adding transaction:", t)
			}

		}
	}

	fmt.Println("Transactions details:", transactions)
	return transactions, nil

}

func main() {
	fmt.Println("Starting with the banking application")
	users := make(map[string]User, 0)
	transactions := make(map[string]Transaction, 0)
	new_bank := bank{
		users:        users,
		transactions: transactions,
	}

	router := mux.NewRouter()

	router.HandleFunc("/ping", PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/createCustomer", CreateCustomerHandler(new_bank)).Methods(http.MethodPost)
	router.HandleFunc("/getAccountDetails", GetAccountDetailsHandler(new_bank)).Methods(http.MethodPost)
	router.HandleFunc("/getCustomerList", GetCustomerListHandler(new_bank)).Methods(http.MethodGet)
	router.HandleFunc("/depositAccount", DepositAmountHandler(new_bank)).Methods(http.MethodPost)
	router.HandleFunc("/withdrawAccount", WithdrawAmountHandler(new_bank)).Methods(http.MethodPost)
	router.HandleFunc("/getTransactions", GetTransactionDetailsHandler(new_bank)).Methods(http.MethodPost)
	http.ListenAndServe("127.0.0.1:8080", router)
}

type PingResponse struct {
	Message string `json:"message"`
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{
		Message: "pong",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBytes)
}

func CreateCustomerHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u User

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		u = b.CreateAccount(u)
		responseBytes, err := json.Marshal(u)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

type AccountDetailsRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func GetAccountDetailsHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accountDetailsRequest AccountDetailsRequest
		err := json.NewDecoder(r.Body).Decode(&accountDetailsRequest)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		u, err := b.GetCustomerDetails(accountDetailsRequest.ID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.Marshal(u)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

func GetCustomerListHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var customers = make([]User, 0)
		customers = b.GetCustomerList()
		responseBytes, err := json.Marshal(customers)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}

type DepositWithdrawAmountRequest struct {
	AccountDetailsRequest
	Amount string `json:"amount"`
}

func DepositAmountHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var depositAmountRequest DepositWithdrawAmountRequest
		err := json.NewDecoder(r.Body).Decode(&depositAmountRequest)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		amount, err := strconv.ParseFloat(depositAmountRequest.Amount, 32)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = b.DepositAmount(depositAmountRequest.ID, float32(amount))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("Deposited amount successfully"))

	})
}

func WithdrawAmountHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var withdrawAmountRequest DepositWithdrawAmountRequest
		err := json.NewDecoder(r.Body).Decode(&withdrawAmountRequest)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		amount, err := strconv.ParseFloat(withdrawAmountRequest.Amount, 32)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = b.WithdrawAmount(withdrawAmountRequest.ID, float32(amount))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("Debited amount successfully"))
	})
}

type GetTransactionDetailsRequest struct {
	AccountDetailsRequest
	StartDate string `json:"startDate"`
	EndDate   string `json:"EndDate"`
}

func GetTransactionDetailsHandler(b bank) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var transactionDetailsRequest GetTransactionDetailsRequest
		err := json.NewDecoder(r.Body).Decode(&transactionDetailsRequest)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		transactions := make([]Transaction, 0)
		transactions, err = b.GetTransactionDetails(transactionDetailsRequest.ID,
			transactionDetailsRequest.StartDate, transactionDetailsRequest.EndDate)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.Marshal(transactions)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBytes)
	})
}
