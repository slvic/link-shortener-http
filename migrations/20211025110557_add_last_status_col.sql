-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD COLUMN last_status INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    DROP COLUMN last_status
-- +goose StatementEnd
