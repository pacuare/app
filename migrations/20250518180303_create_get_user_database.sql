-- +goose Up
-- +goose StatementBegin
create function GetUserDatabase(email text) returns text language plpgsql as
$$ begin
    return ('user_' || replace(replace(email, '@', '__'), '.', '_'));
end $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function GetUserDatabase;
-- +goose StatementEnd
