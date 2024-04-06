-- +goose Up
-- +goose StatementBegin
ALTER TABLE accesses
    ADD COLUMN datetime_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accesses
    DROP COLUMN datetime_create;
-- +goose StatementEnd
