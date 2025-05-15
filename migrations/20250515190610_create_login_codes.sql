-- +goose Up
-- +goose StatementBegin
CREATE TABLE LoginCodes (
    email TEXT NOT NULL PRIMARY KEY,
    code TEXT DEFAULT UPPER(SUBSTR(MD5(RANDOM()::TEXT), 0, 7))
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE LoginCodes;
-- +goose StatementEnd
