package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"example.com/banking/bank"
	"example.com/banking/config"
)

const (
	versionHeader = "Accept"
)

func initRouter(dep dependencies) (router *mux.Router) {
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router = mux.NewRouter()
	router.HandleFunc("/ping", bank.PingHandler).Methods(http.MethodGet)

	router.HandleFunc("/login", bank.LoginHandler(dep.BankService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/account", bank.CreateAccountHandler(dep.BankService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/accounts", bank.GetAccountsHandler(dep.BankService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/account/{account_id}", bank.GetAccountDetailsHandler(dep.BankService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/account/{account_id}/deposit", bank.DepositAmountHandler(dep.BankService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/account/{account_id}/withdraw", bank.WithdrawAmountHandler(dep.BankService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/account/{account_id}/transactions", bank.GetTransactionDetailsHandler(dep.BankService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/account/{account_id}", bank.DeleteAccountHandler(dep.BankService)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	sh := http.StripPrefix("/docs/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/docs/").Handler(sh)
	return
}
