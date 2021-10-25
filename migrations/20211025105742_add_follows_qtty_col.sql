-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD COLUMN follow_qtty INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    DROP COLUMN follow_qtty;
-- +goose StatementEnd
