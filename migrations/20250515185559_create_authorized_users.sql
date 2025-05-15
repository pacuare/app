-- +goose Up
-- +goose StatementBegin

CREATE TABLE AuthorizedUsers (
    email TEXT NOT NULL PRIMARY KEY,
    fullAccess BOOLEAN DEFAULT FALSE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE AuthorizedUsers;

-- +goose StatementEnd
