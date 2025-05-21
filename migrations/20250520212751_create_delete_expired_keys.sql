-- +goose Up
-- +goose StatementBegin
create procedure DeleteExpiredKeys()
language sql
as $$
    delete from APIKeys where createdOn < ('now'::date - '30 days'::interval);
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop procedure DeleteExpiredKeys;
-- +goose StatementEnd
