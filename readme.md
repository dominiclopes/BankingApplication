
Curl Commands: 

curl localhost:8080/ping
curl localhost:8080/getCustomerList -v
curl localhost:8080/createCustomer -X POST -d '{"email":"user_email", "phoneNumber": "user_contact"}' -v
curl localhost:8080/getAccountDetails -X POST -d '{"id":"user_id", "password": "user_password"}' -v
curl localhost:8080/depositAccount -X POST -d '{"id":"user_id", "password": "user_password", "amount": "100"}' -v
curl localhost:8080/withdrawAccount -X POST -d '{"id":"user_id", "password": "user_password", "amount": "50"}' -v
curl localhost:8080/getTransactions -X POST -d '{"id":"user_id", "password": "user_password", "startDate": "10-09-2022", "endDate": "12-09-2022"}' -v
                                    


sudo -i -u postgres
access the PostgreSQL prompt immediately by typing: psql
exit out of the PostgreSQL prompt by typing: \q

Accessing a Postgres Prompt Without Switching Accounts: sudo -u postgres psql


CREATE TABLE users(
    id          VARCHAR(50) PRIMARY KEY ,
    email       VARCHAR(50) NOT NULL,
    phonenumber VARCHAR(50) NOT NULL,
    password    VARCHAR(50) NOT NULL,
    balance     DECIMAL NOT NULL
);


CREATE TABLE transactions(
    id          VARCHAR(50) PRIMARY KEY ,
    type        VARCHAR(50) NOT NULL,
    amount      DECIMAL NOT NULL,
    balance     DECIMAL NOT NULL,
    datetime    TIMESTAMP NOT NULL,
    userid      VARCHAR(50) NOT NULL
);