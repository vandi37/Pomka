-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS UserToPromo (
    UserId BIGINT REFERENCES Users(Id),
    PromoId BIGINT REFERENCES Promos(Id),
    activatedAt TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS UserToPromo;
-- +goose StatementEnd
