-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS userToPromo (
    userId bigint NOT NULL FOREIGN KEY REFERENCES users(userId),
    promoName varchar(255) NOT NULL FOREIGN KEY REFERENCES promos(Name),
    actAt timestamp NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE userToPromo;
-- +goose StatementEnd
