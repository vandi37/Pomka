-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Checks"  (
    "Id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "CreatorId" BIGINT REFERENCES "Users"("Id"),
    "Key" TEXT NOT NULL UNIQUE,
    "Currency" SMALLINT CHECK ("Currency" = 0 OR "Currency" = 1 OR "Currency" = 2),
    "Amount" INT NOT NULL CHECK ("Amount" > 0),
    "CreatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "Checks";
-- +goose StatementEnd
