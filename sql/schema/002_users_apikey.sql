-- +goose Up
Alter table users add column api_key varchar(64) unique not null default (
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
Alter TABLE users DROP COLUMN api_key;
