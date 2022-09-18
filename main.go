package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"example.com/banking/repositories"
	"example.com/banking/services"
)

func main() {
	fmt.Println("Starting with the banking application")

	dbConn, err := repositories.CreateDBConnection()
	if err != nil {
		panic(err)
	}
	defer repositories.CloseDBConnection(dbConn)

	bankStore := repositories.NewBankStore(dbConn)
	newBank := services.NewBank(bankStore)

	router := mux.NewRouter()

	router.HandleFunc("/ping", services.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/account", services.CreateAccountHandler(newBank)).Methods(http.MethodPost)
	router.HandleFunc("/accounts", services.GetAccountsHandler(newBank)).Methods(http.MethodGet)
	router.HandleFunc("/account/{account_id}", services.GetAccountDetailsHandler(newBank)).Methods(http.MethodGet)
	router.HandleFunc("/account/{account_id}/deposit", services.DepositAmountHandler(newBank)).Methods(http.MethodPost)
	router.HandleFunc("/account/{account_id}/withdraw", services.WithdrawAmountHandler(newBank)).Methods(http.MethodPost)
	router.HandleFunc("/account/{account_id}/transactions", services.GetTransactionDetailsHandler(newBank)).Methods(http.MethodPost)

	http.ListenAndServe("127.0.0.1:8080", router)
}
