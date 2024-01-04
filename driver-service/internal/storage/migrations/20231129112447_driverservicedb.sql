-- +goose Up
-- +goose StatementBegin
-- CREATE TYPE TaxiType AS ENUM (0,1,2);

-- CREATE TYPE DriverStatus AS ENUM (0, 1, 2);

CREATE TABLE drivers
(
    id            uuid PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    surname       VARCHAR(255) NOT NULL,
    phone         VARCHAR(255) NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    balance       float        NOT NULL,
    taxi_type     INT     NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    rating        FLOAT        NOT NULL DEFAULT 0.0,
    status        INT NOT NULL,
    registration  TIMESTAMP    NOT NULL DEFAULT NOW(),
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE trips
(
    id        uuid PRIMARY KEY,
    driver_id uuid REFERENCES drivers (id) ON DELETE CASCADE NOT NULL,
    user_id   uuid                                           NOT NULL,
    "from"    VARCHAR(255)                                   NOT NULL,
    "to"      VARCHAR(255)                                   NOT NULL,
    rating    float                                          NOT NULL CHECK ( rating >= 0.0 AND rating <= 5.0),
    date      TIMESTAMP                                      NOT NULL
);

CREATE FUNCTION update_rating() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$update_rating$
BEGIN
    UPDATE drivers
    SET rating = (SELECT AVG(trips.rating) FROM trips WHERE trips.driver_id = NEW.driver_id)
    WHERE drivers.id = NEW.driver_id;

    RETURN NULL;
end;
$update_rating$;

CREATE TRIGGER update_rating
    AFTER INSERT
    ON trips
    FOR EACH ROW
EXECUTE FUNCTION update_rating();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_rating ON drivers;

DROP TABLE trips;
DROP TABLE drivers;

-- DROP TYPE TaxiType;
-- DROP TYPE DriverStatus;
-- +goose StatementEnd
