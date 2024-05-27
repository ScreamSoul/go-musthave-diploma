-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    registration_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE order_loading_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    number SERIAL PRIMARY KEY,
    user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
    status order_loading_status NOT NULL DEFAULT 'NEW',
    accrual INTEGER CHECK (accrual >= 0),
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE loyalty_wallets (
    user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
    balance INTEGER CHECK (balance >= 0) NOT NULL DEFAULT 0,
    spent INTEGER CHECK (spent >= 0) NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id)
);

CREATE TABLE loyalty_wallet_operations (
    operation_id UUID PRIMARY KEY,
    user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
    order_id INTEGER REFERENCES "orders"(number) ON DELETE SET NULL,
    amount_deducted INTEGER CHECK (amount_deducted >= 0) NOT NULL DEFAULT 0,
    operation_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Составной индекс для оптимизации запросов по пользователю
CREATE INDEX idx_user_login_password ON "users"(login, password_hash);

-- Индекс для оптимизации запросов по операциям с кошельком лояльности
CREATE INDEX idx_loyaltywalletoperation_user ON "loyalty_wallet_operations"(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_login_password;
DROP INDEX IF EXISTS idx_loyaltywalletoperation_user;
DROP TYPE IF EXISTS order_loading_statuы;

DROP TABLE IF EXISTS loyalty_wallets;
DROP TABLE IF EXISTS loyalty_wallet_operations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd
