-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Promos  (
    Id bigint PRIMARY KEY UNIQUE,
    Name varchar(255) NOT NULL UNIQUE,
    Currency smallint,
    Amount int NOT NULL,
    Uses int DEFAULT 1,
    Creator bigint REFERENCES Users(Id),
    ExpAt timestamp NOT NULL,
    CreatedAt timestamp DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Promos;
-- +goose StatementEnd

