-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS promos  (
    Name varchar(255) NOT NULL UNIQUE,
    Value int NOT NULL,
    Creator varchar(255),
    Currency int,
    ExpAt timestamp NOT NULL,
    CreatedAt timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE promos;
-- +goose StatementEnd

