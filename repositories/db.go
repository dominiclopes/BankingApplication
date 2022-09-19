package repositories

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	uuidgen "github.com/pborman/uuid"
)

const (
	DBHOST = "localhost"
	DBPORT = 5432
	DBUSER = "postgres"
	DBPASS = "12345"
	DBNAME = "bank"
)

func CreateDBConnection() (db *sql.DB, err error) {
	psqlConnectionString := fmt.Sprintf("host=%v password=%v port=%v user=%v dbname=%v sslmode=disable",
		DBHOST, DBPASS, DBPORT, DBUSER, DBNAME)
	fmt.Println("Using connection string:", psqlConnectionString)

	db, err = sql.Open("postgres", psqlConnectionString)
	if err != nil {
		err = fmt.Errorf("error opening database connection, err: %v", err)
		return
	}

	if err = db.Ping(); err != nil {
		err = fmt.Errorf("error creating database connection, err: %v", err)
		return
	}

	return db, nil
}

func CloseDBConnection(db *sql.DB) (err error) {
	if err = db.Close(); err != nil {
		err = fmt.Errorf("error closing database connection, err: %v", err)
		return
	}
	return
}

type User struct {
	ID          string  `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	PhoneNumber string  `json:"phone_number" db:"phonenumber"`
	Password    string  `json:"password" db:"password"`
	Balance     float32 `json:"balance" db:"balance"`
	Type        string  `json:"-" db:"type"`
}

type Transaction struct {
	ID       string  `json:"-" db:"id"`
	Type     string  `json:"type" db:"type"`
	Amount   float32 `json:"amount" db:"amount"`
	Balance  float32 `json:"balance" db:"balance"`
	DateTime string  `json:"time" db:"datetime"`
	UserID   string  `json:"-" db:"userid"`
}

type BankStorer interface {
	CreateAccount(u User) (err error)
	GetAccountList() (users []User, err error)
	GetAccountDetails(accID string) (u User, err error)
	AddTransaction(t Transaction) (err error)
	DepositAmount(accID string, amount float32) (balance float32, err error)
	WithdrawAmount(accID string, amount float32) (balance float32, err error)
	GetTransactions(accID string) (transactions []Transaction, err error)
}

const (
	createAccountQuery = `insert into users(id, email, phonenumber, password, balance, type) 
	values($1, $2, $3, $4, $5, $6)`
	listAccountsQuery             = `select * from users`
	getAccountByIDQuery           = `select * from users where id=$1`
	updateAccountBalanceByIDQuery = `update users set balance=$1 where id=$2`
	createTransactionQuery        = `insert into transactions(id, type, amount, balance, datetime, userid) 
	values ($1, $2, $3, $4, $5, $6)`
	transactAmountQuery  = `update users set balance=$1 where id=$2`
	getTransactionsQuery = `select * from transactions where userid=$1`
)

type bankStore struct {
	db *sql.DB
}

func NewBankStore(db *sql.DB) BankStorer {
	return &bankStore{db: db}
}

func (m *bankStore) CreateAccount(u User) (err error) {
	if _, err = m.db.Exec(createAccountQuery, u.ID, u.Email, u.PhoneNumber, u.Password, u.Balance, u.Type); err != nil {
		err = fmt.Errorf("error inserting user data, err: %v", err)
		return
	}
	return
}

func (m *bankStore) GetAccountList() (users []User, err error) {
	rows, err := m.db.Query(listAccountsQuery)
	if err != nil {
		err = fmt.Errorf("error getting accounts details, err: %v", err)
		return
	}
	defer rows.Close()

	users = make([]User, 0)
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Email, &u.PhoneNumber, &u.Password, &u.Balance, &u.Type); err != nil {
			err = fmt.Errorf("error scanning fetched rows: %v", err)
			return
		}
		users = append(users, u)
	}
	return
}

func (m *bankStore) GetAccountDetails(accID string) (u User, err error) {

	row := m.db.QueryRow(getAccountByIDQuery, accID)

	if err = row.Scan(&u.ID, &u.Email, &u.PhoneNumber, &u.Password, &u.Balance, &u.Type); err != nil {
		err = fmt.Errorf("error scanning the fetched row, err: %v", err)
		return
	}
	return
}

func (m *bankStore) AddTransaction(t Transaction) (err error) {

	if _, err = m.db.Exec(createTransactionQuery, t.ID, t.Type, t.Amount, t.Balance, t.DateTime, t.UserID); err != nil {
		err = fmt.Errorf("error inserting tranasaction details, err: %v", err)
		return
	}
	return
}

func (m *bankStore) DepositAmount(accID string, amount float32) (balance float32, err error) {
	u, err := m.GetAccountDetails(accID)
	if err != nil {
		err = fmt.Errorf("account %v not present", accID)
		return
	}

	// Update the user balance
	balance = u.Balance + amount

	// update user details
	if _, err = m.db.Exec(transactAmountQuery, balance, accID); err != nil {
		err = fmt.Errorf("error updating the user balance, err: %v", err)
		return 0.0, err
	}

	// Add a transaction
	t := Transaction{
		ID:       uuidgen.New(),
		UserID:   accID,
		Amount:   amount,
		Balance:  balance,
		Type:     "Credit",
		DateTime: time.Now().Format("2006-01-02 15:04:05.000"),
	}
	if err = m.AddTransaction(t); err != nil {
		return 0.0, err
	}

	fmt.Printf("Credited amount: %v, in account: %v. Balance: %v\n", amount, accID, balance)
	return
}

func (m *bankStore) WithdrawAmount(accID string, amount float32) (balance float32, err error) {
	u, err := m.GetAccountDetails(accID)
	if err != nil {
		err = fmt.Errorf("account %v not present", accID)
		return
	}

	// Verify if amount can be debited
	if u.Balance < amount {
		err = fmt.Errorf("amount %v cannot be debited from account %v. insufficient funds: %v",
			amount, accID, u.Balance)
		return
	}

	// Update the user balance
	balance = u.Balance - amount

	// update user details
	if _, err = m.db.Exec(transactAmountQuery, balance, accID); err != nil {
		err = fmt.Errorf("error updating the user balance, err: %v", err)
		return 0.0, err
	}

	// Add a transaction
	t := Transaction{
		ID:       uuidgen.New(),
		UserID:   accID,
		Amount:   amount,
		Balance:  balance,
		Type:     "Debit",
		DateTime: time.Now().Format("2006-01-02 15:04:05.000"),
	}
	if err = m.AddTransaction(t); err != nil {
		return 0.0, err
	}

	fmt.Printf("Debited amount: %v, from account: %v. Balance: %v\n", amount, accID, balance)
	return
}

func (m *bankStore) GetTransactions(accID string) (transactions []Transaction, err error) {
	if _, err = m.GetAccountDetails(accID); err != nil {
		err = fmt.Errorf("account %v not present", accID)
		return
	}

	fmt.Println(getTransactionsQuery, accID)
	rows, err := m.db.Query(getTransactionsQuery, accID)
	if err != nil {
		err = fmt.Errorf("error getting account transactions, err: %v", err)
		return
	}

	transactions = make([]Transaction, 0)
	for rows.Next() {
		var t Transaction

		if err = rows.Scan(&t.ID, &t.Type, &t.Amount, &t.Balance, &t.DateTime, &t.UserID); err != nil {
			err = fmt.Errorf("error scanning rows, err: %v", err)
			return
		}

		transactions = append(transactions, t)
	}

	fmt.Println("Transactions details:", transactions)
	return
}
