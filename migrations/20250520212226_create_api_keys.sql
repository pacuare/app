-- +goose Up
-- +goose StatementBegin
create table APIKeys (
    id serial primary key,
    email text not null,
    key text not null default md5(random()::text),
    description text,
    createdOn date default 'now'::date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table APIKeys;
-- +goose StatementEnd
