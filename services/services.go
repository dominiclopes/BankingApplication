package services

import (
	"fmt"
	"time"

	uuidgen "github.com/pborman/uuid"

	"example.com/banking/models"
	"example.com/banking/repositories"
)

type Service interface {
	CreateAccount(u repositories.User) (acc models.CreateAccountResponse, err error)
	GetAccountList() (accounts []repositories.User, err error)
	GetAccountDetails(accId string) (acc repositories.User, err error)
	DepositAmount(accId string, amount float32) (bal models.DepositWithdrawAmountResponse, err error)
	WithdrawAmount(accId string, amount float32) (bal models.DepositWithdrawAmountResponse, err error)
	GetTransactionDetails(accId string, startDate, endDate string) (transactions []repositories.Transaction, err error)
}

type bankService struct {
	bankStore repositories.BankStorer
}

func NewBank(bs repositories.BankStorer) Service {
	return &bankService{bankStore: bs}
}

func (b *bankService) CreateAccount(u repositories.User) (acc models.CreateAccountResponse, err error) {
	fmt.Printf("Creating an account for user email: %v, phone number: %v\n", u.Email, u.PhoneNumber)

	// Create the user ID, password and update the balance
	u.ID = uuidgen.New()
	u.Password = uuidgen.New()
	u.Balance = 0.0

	// Save the user in the bank
	err = b.bankStore.CreateAccount(u)
	if err != nil {
		return
	}

	// Create the response
	acc = models.CreateAccountResponse{
		ID:       u.ID,
		Password: u.Password,
	}

	fmt.Printf("Created account with details: %v, for user with email: %v, phone number: %v. Opening balance: %v\n",
		acc, u.Email, u.Email, u.Balance)

	return
}

func (b *bankService) GetAccountList() (accounts []repositories.User, err error) {
	fmt.Println("Getting the list of accounts in the bank")
	accounts, err = b.bankStore.GetAccountList()
	return
}

func (b *bankService) GetAccountDetails(accId string) (acc repositories.User, err error) {
	fmt.Printf("Getting the customer details for account: %v\n", accId)
	acc, err = b.bankStore.GetAccountDetails(accId)
	return
}

func (b *bankService) DepositAmount(accId string, amount float32) (bal models.DepositWithdrawAmountResponse, err error) {
	fmt.Printf("Depositing amount: %v in account: %v\n", amount, accId)

	balance, err := b.bankStore.DepositAmount(accId, amount)
	if err != nil {
		return
	}

	bal = models.DepositWithdrawAmountResponse{Balance: balance}
	return
}

func (b *bankService) WithdrawAmount(accId string, amount float32) (bal models.DepositWithdrawAmountResponse, err error) {
	fmt.Printf("Withdrawing amount: %v from account: %v\n", amount, accId)

	balance, err := b.bankStore.WithdrawAmount(accId, amount)
	if err != nil {
		return
	}

	bal = models.DepositWithdrawAmountResponse{Balance: balance}
	return
}

func (b *bankService) GetTransactionDetails(accId string, startDate, endDate string) (transactions []repositories.Transaction, err error) {
	fmt.Printf("Getting transactions details for account: %v, from %v to %v\n", accId, startDate, endDate)
	allTransactions, err := b.bankStore.GetTransactions(accId)
	if err != nil {
		return
	}

	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		err = fmt.Errorf("error parsing startdate: %v", startDate)
		return
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		err = fmt.Errorf("error parsing endate %v", endDate)
		return
	}

	for _, t := range allTransactions {
		if t.UserID == accId {
			tDateTimeAsNeeded, err := time.Parse("2006-01-02T15:04:05.000Z", t.DateTime)
			if err != nil {
				err = fmt.Errorf("error parsing transaction date: %v", t.DateTime)
				return nil, err
			}

			// tDateTimeAsNeeded, err := time.Parse("02-01-2006",
			// 	fmt.Sprintf("%02d-%02d-%d", tDateTime.Day(), tDateTime.Month(), tDateTime.Year()))
			// if err != nil {
			// 	err = fmt.Errorf("Error parsing transaction date: %v", tDateTime)
			// 	return nil, err
			// }
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
	return
}
