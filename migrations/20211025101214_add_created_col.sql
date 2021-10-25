-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD COLUMN created_at TIMESTAMP;
ALTER TABLE links ALTER COLUMN created_at SET DEFAULT (now() at time zone 'utc');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    DROP COLUMN created_at
-- +goose StatementEnd
