-- +goose Up
-- +goose StatementBegin
CREATE TABLE links
(
    id        BIGSERIAL
        CONSTRAINT id_pk PRIMARY KEY,
    original  varchar NOT NULL,
    shortened varchar NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
