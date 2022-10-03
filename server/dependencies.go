package server

import (
	"example.com/banking/app"
	"example.com/banking/bank"
	"example.com/banking/db"
)

type dependencies struct {
	BankService bank.Service
}

func initDependencies() (dependencies, error) {
	logger := app.GetLogger()

	appDB := app.GetDB()
	dbStore := db.NewStorer(appDB)

	bankService := bank.NewBankService(dbStore, logger)

	return dependencies{
		BankService: bankService,
	}, nil
}
