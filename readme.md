
curl localhost:8080/ping
curl localhost:8080/getCustomerList -v
curl localhost:8080/createCustomer -X POST -d '{"email":"rahul.jadhav@joshsoftware.com", "phoneNumber": "12345676"}' -v
curl localhost:8080/getAccountDetails -X POST -d '{"id":"bd8a8449-5706-4839-8d70-198b800f5b72", "password": "d5d2255a-fc98-459a-aba8-26acbd1a2090"}' -v
curl localhost:8080/depositAccount -X POST -d '{"id":"bd8a8449-5706-4839-8d70-198b800f5b72", "password": "d5d2255a", "amount": "100"}' -v
curl localhost:8080/withdrawAccount -X POST -d '{"id":"bd8a8449-5706-4839-8d70-198b800f5b72", "password": "d5d2255a", "amount": "50"}' -v
curl localhost:8080/getTransactions -X POST -d '{"id":"bd8a8449-5706-4839-8d70-198b800f5b72", "password": "d5d2255a-fc98-459a-aba8-26acbd1a2090", "startDate": "10-09-2022", "endDate": "12-09-2022"}' -v



new_bank.GetCustomerList()

	u := User{Email: "dominiclopes@abc.com", PhoneNumber: 1234567}
	accId, _ := new_bank.CreateAccount(u)

	new_bank.GetCustomerList()

	new_bank.GetCustomerDetails(accId)

	new_bank.DepositAmount(accId, 100000)
	new_bank.GetCustomerDetails(accId)

	new_bank.WithdrawAmount(accId, 3000)
	new_bank.GetCustomerDetails(accId)

	transactions, err := new_bank.GetTransactionDetails(accId, "11-09-2022", "12-09-2022")
	if err != nil {
		fmt.Println(err)
	}




// func main() {
// 	fmt.Println("Starting with the banking application")

// 	roles := map[string]Role{
// 		"Accountant": {
// 			ID:   uuidgen.New(),
// 			Name: "Accountant",
// 		},
// 		"Customer": {
// 			ID:   uuidgen.New(),
// 			Name: "Customer",
// 		},
// 	}
// 	fmt.Println(roles)

// 	accounTypes := map[string]AccountType{
// 		"Savings": {
// 			ID:   uuidgen.New(),
// 			Type: "Savings",
// 		},
// 		"Current": {
// 			ID:   uuidgen.New(),
// 			Type: "Current",
// 		},
// 	}
// 	fmt.Println(accounTypes)

// 	transactionTypes := map[string]TransactionType{
// 		"Credit": {
// 			ID:   uuidgen.New(),
// 			Type: "Credit",
// 		},
// 		"Debit": {
// 			ID:   uuidgen.New(),
// 			Type: "Debit",
// 		},
// 	}
// 	fmt.Println(transactionTypes)

// 	var users map[string]User = make(map[string]User)
// 	var usersAccountDetails map[string]UserAccountDetails = make(map[string]UserAccountDetails)
// 	var transactions map[string]Transaction = make(map[string]Transaction)

// 	newBank := bank{
// 		roles:               roles,
// 		accountTypes:        accounTypes,
// 		transactionTypes:    transactionTypes,
// 		users:               users,
// 		usersAccountDetails: usersAccountDetails,
// 		transactions:        transactions,
// 	}
// 	fmt.Println(newBank)

// 	http.HandleFunc("/ping", PingHandler)
// 	http.HandleFunc("/createAccount", CreateAccountHandler(newBank))
// 	http.ListenAndServe("127.0.0.1:8080", nil)
// }

// type PingResponse struct {
// 	Message string `json:"message"`
// }

// func PingHandler(w http.ResponseWriter, r *http.Request) {
// 	response := PingResponse{Message: "pong"}

// 	responseBytes, err := json.Marshal(response)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Add("Content-Type", "application/json")
// 	w.Write(responseBytes)
// }

// type CreateAccountRequest struct {
// 	UserEmail       string `json:"email"`
// 	UserPhoneNumber int    `json:"phoneNumber"`
// 	AccountType     string `json:"accountType"`
// }

// type CreateAccountResponse struct {
// 	UserEmail       string `json:"email"`
// 	UserPhoneNumber int    `json:"phoneNumber"`
// 	AccountType     string `json:"accountType"`
// 	BackAccountID   string `json:"bankAccID"`
// 	Password        string `json:"password"`
// }

// func CreateAccountHandler(b bank) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// todo We need to ensure its the accountant who is creating the account
// 		// We need the user details in the form of json
// 		// we must get the user details, role details from the request
// 		var uAccReq CreateAccountRequest

// 		err := json.NewDecoder(r.Body).Decode(&uAccReq)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		fmt.Println("Request data:", uAccReq)
// 		u := User{
// 			Email:       uAccReq.UserEmail,
// 			PhoneNumber: uAccReq.UserPhoneNumber,
// 			RoleID:      b.roles["Customer"].ID,
// 		}

// 		uAccResp, err := b.CreateAccount(u, b.accountTypes[uAccReq.AccountType])
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		var uAccRespBytes []byte
// 		uAccRespBytes, err = json.Marshal(uAccResp)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Add("Content-Type", "application/json")
// 		w.Write(uAccRespBytes)
// 	})

// }

// func GetAccountListHandler() {
// }

// func GetAccountDetailsHandler() {
// }

// func GetTransactionDetailsHandler() {

// }
