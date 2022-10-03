CREATE TABLE users(
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(50) NOT NULL UNIQUE,
    phone_number VARCHAR(10) NOT NULL,
    password     VARCHAR(60) NOT NULL,
    type         VARCHAR(10) NOT NULL
);

CREATE TABLE accounts(
    id      UUID PRIMARY KEY,
    balance DECIMAL NOT NULL DEFAULT 0.0,
    user_id INTEGER REFERENCES users (id)
);

CREATE TABLE transactions(
    id          UUID PRIMARY KEY,
    type        VARCHAR(6) NOT NULL,
    amount      DECIMAL NOT NULL,
    balance     DECIMAL NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    account_id  UUID REFERENCES accounts (id)
);

/* Add the accountant details */
INSERT INTO users(email, phone_number, password, type) VALUES('account@bank.com', '8655645204', crypt('josh@123', gen_salt('bf')), 'accountant');
