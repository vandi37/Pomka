-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS UserToPromo (
    UserId bigint REFERENCES Users(Id),
    PromoId bigint REFERENCES Promos(Id),
    activatedAt timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS UserToPromo;
-- +goose StatementEnd
