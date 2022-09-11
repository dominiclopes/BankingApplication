
Curl Commands: 

curl localhost:8080/ping
curl localhost:8080/getCustomerList -v
curl localhost:8080/createCustomer -X POST -d '{"email":"user_email", "phoneNumber": "user_contact"}' -v
curl localhost:8080/getAccountDetails -X POST -d '{"id":"user_id", "password": "user_password"}' -v
curl localhost:8080/depositAccount -X POST -d '{"id":"user_id", "password": "user_password", "amount": "100"}' -v
curl localhost:8080/withdrawAccount -X POST -d '{"id":"user_id", "password": "user_password", "amount": "50"}' -v
curl localhost:8080/getTransactions -X POST -d '{"id":"user_id", "password": "user_password", "startDate": "10-09-2022", "endDate": "12-09-2022"}' -v
                                    