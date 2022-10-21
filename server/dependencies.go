package server

import (
	"github.com/dominiclopes/BankingApplication/app"
	"github.com/dominiclopes/BankingApplication/bank"
	"github.com/dominiclopes/BankingApplication/db"
)

type dependencies struct {
	BankService bank.Service
}

func initDependencies() (dependencies, error) {
	logger := app.GetLogger()

	appDB := app.GetDB()
	dbStore := db.NewStorer(appDB, logger)

	bankService := bank.NewBankService(dbStore, logger)

	return dependencies{
		BankService: bankService,
	}, nil
}
