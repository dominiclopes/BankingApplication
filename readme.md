Problem Statement: 
Design a bank management system with the following features
- There are two roles Accountant and Customer
- Accountants can create a bank account for the customer by taking their email and phone number and returning them bank account details(Including bank account ID and randomly generated password)
- Customers can deposit and withdraw money to their respective bank accounts by providing their bank account number and generated password
- Customers can also see their current bank account balance by providing their bank account number and generated password 
- Customers can also see their account history for the given date range

Assumptions
- Customers can open only a single account with one email address
- There is only one Accountant and the accountant's email address is account@bank.com. The password is “josh@123” and should be inserted into the system when we start the program execution
- Account history should include: Date, Transaction Type(Debit OR Credit), Transaction amount, The total balance remaining in the account after the transaction


Solution: 
The Banking Application provided the following features
- create user account
- list all user accounts
- get user account details 
- credit amount to an account
- debit amount from an account
- list transactions for an account


To start the application, execute: go run main.go start

To run migrations, execute: go run main.go create_migration

For writing unit testcases, used mockery
docker pull vektra/mockery
