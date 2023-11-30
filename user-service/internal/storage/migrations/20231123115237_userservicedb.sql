-- +goose Up
-- +goose StatementBegin
CREATE TYPE taxiType AS ENUM ('economy','comfort','business');

CREATE TYPE transactionType AS ENUM ('refill' , 'spent');

CREATE TYPE transactionStatus AS ENUM ('create', 'blocked', 'canceled', 'success');

CREATE TABLE users
(
    id            uuid PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    surname       VARCHAR(255) NOT NULL,
    phone         VARCHAR(255) NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    rating        FLOAT        NOT NULL DEFAULT 0.0,
    date          TIMESTAMP    NOT NULL DEFAULT NOW(),
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE wallets
(
    id      uuid PRIMARY KEY,
    user_id uuid         NOT NULL,
    card    VARCHAR(255) NOT NULL UNIQUE,
    balance FLOAT        NOT NULL DEFAULT 0.0 CHECK ( balance >= 0.0 )
);

CREATE TABLE family_wallets
(
    id            uuid PRIMARY KEY,
    wallet_id     uuid REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    balance       FLOAT                                          NOT NULL DEFAULT 0.0 CHECK ( balance >= 0.0 ),
    fixed_balance FLOAT                                          NOT NULL DEFAULT 0.0 CHECK ( fixed_balance >= 0.0 )
);

CREATE TABLE users_wallets
(
    user_id   uuid REFERENCES users (id) ON DELETE CASCADE          NOT NULL,
    wallet_id uuid REFERENCES family_wallets (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE trips
(
    id        uuid PRIMARY KEY,
    user_id   uuid REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    taxi_type taxiType                                     NOT NULL,
    driver_id uuid                                         NOT NULL,
    "from"    VARCHAR(255)                                 NOT NULL,
    "to"      VARCHAR(255)                                 NOT NULL,
    rating    float                                        NOT NULL CHECK ( rating >= 0.0 AND rating <= 5.0),
    date      TIMESTAMP                                    NOT NULL
);

CREATE TABLE transactions
(
    id                  SERIAL PRIMARY KEY,
    wallet_id           uuid              NOT NULL,
    money               FLOAT             NOT NULL,
    transaction_type    transactionType   NOT NULL,
    transactions_status transactionStatus NOT NULL DEFAULT 'create',
    date                TIMESTAMP         NOT NULL DEFAULT NOW()
);

CREATE FUNCTION clean_trips() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$clean_trips$
BEGIN
DELETE
FROM trips
WHERE trips.user_id = NEW.user_id
  AND trips.date = (SELECT MIN(trips.date) FROM trips WHERE trips.user_id = NEW.user_id)
  AND 20 = (SELECT COUNT(*) FROM trips WHERE trips.user_id = NEW.user_id);

RETURN NULL;
end;
$clean_trips$;

CREATE FUNCTION update_rating() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$update_rating$
BEGIN
UPDATE users
SET users.rating = (SELECT AVG(trips.rating) FROM trips WHERE trips.user_id = NEW.user_id)
WHERE users.id = NEW.user_id;

RETURN NULL;
end;
$update_rating$;

CREATE TRIGGER clean_trips
    BEFORE INSERT
    ON trips
    FOR EACH ROW
    EXECUTE FUNCTION clean_trips();

CREATE TRIGGER update_rating
    AFTER INSERT
    ON trips
    FOR EACH ROW
    EXECUTE FUNCTION update_rating();

CREATE FUNCTION update_balance() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$update_balance$
BEGIN
UPDATE family_wallets
SET fixed_balance=0.0
WHERE fixed_balance > (SELECT w.balance
                       FROM wallets w
                                INNER JOIN family_wallets fw on w.id = fw.wallet_id);

UPDATE wallets
SET balance=balance +
            + (SELECT balance
               FROM family_wallets
               WHERE family_wallets.wallet_id = wallets.id
                 AND family_wallets.balance < family_wallets.fixed_balance) -
            - (SELECT fixed_balance
               FROM family_wallets
               WHERE family_wallets.wallet_id = wallets.id
                 AND family_wallets.balance < family_wallets.fixed_balance);

UPDATE family_wallets
SET balance=fixed_balance
WHERE balance < family_wallets.fixed_balance;
end;
$update_balance$;

CREATE TRIGGER update_balance
    AFTER UPDATE
    ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_balance();

CREATE UNIQUE INDEX trips_user_idx ON trips (user_id);

CREATE UNIQUE INDEX transaction_idx ON transactions (id);
CREATE UNIQUE INDEX transaction_wallet_idx ON transactions (wallet_id);

CREATE UNIQUE INDEX id_idx ON users (id);
CREATE UNIQUE INDEX phone_idx ON users (phone);
CREATE UNIQUE INDEX wallet_idx ON wallets (id);

CREATE UNIQUE INDEX wallet_user_idx ON wallets (user_id);
CREATE UNIQUE INDEX family_wallet_idx ON family_wallets (id);
CREATE UNIQUE INDEX user_wallet_idx ON users_wallets (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_wallets;
DROP TABLE family_wallets;
DROP TRIGGER update_balance ON wallets;
DROP TABLE wallets;

DROP TRIGGER update_rating ON trips;
DROP TRIGGER clean_trips ON trips;
DROP TYPE taxiType;
DROP TABLE trips;

DROP TYPE transactionType;
DROP TYPE transactionStatus;
DROP TABLE transactions;

DROP TABLE users;
-- +goose StatementEnd
