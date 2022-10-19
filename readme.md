# Problem Statement: 
Design a bank management system with the following features
- There are two roles Accountant and Customer
- Accountants can create a bank account for the customer by taking their email and phone number and returning them bank account details(Including bank account ID and randomly generated password)
- Customers can deposit and withdraw money to their respective bank accounts by providing their bank account number and generated password
- Customers can also see their current bank account balance by providing their bank account number and generated password 
- Customers can also see their account history for the given date range

# Assumptions
- Customers can open only a single account with one email address
- There is only one Accountant and the accountant's email address is account@bank.com. The password is “josh@123” and should be inserted into the system when we start the program execution
- Account history should include: Date, Transaction Type(Debit OR Credit), Transaction amount, The total balance remaining in the account after the transaction


# Solution: 
The Banking Application provided the following API's for the account
- create user account
- list all user accounts

The Banking Application provided the following API's for the customer
- get user account details 
- deposit amount to an account
- withdraw amount from an account
- list transactions within a given date range for an account

Jwt token is used for authentication and the token is been accessed by API's using http cookie


# Steps
- To start the application, execute: go run main.go start

- To run migrations, execute: go run main.go create_migration

- For writing unit testcases, used mockery\
    - Install mockery using commad -> docker pull vektra/mockery
    - Create mocks using commad: docker run -v "$PWD":/src -w /src vektra/mockery --dir=bank --name=Service --output=bank/mocks --with-expecter

- For test coverage, execute following commads
    - go test -coverpkg=./... -coverprofile cover.out ./...
    - go tool cover -html=cover.out -o cover.html

# Code Coverage
- Unit tests are created only in example.com/banking/bank package. Test coverage: 6.2% of statements