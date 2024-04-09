-- +goose Up
-- +goose StatementBegin
ALTER TABLE accesses
    ADD COLUMN questions JSON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accesses
    DROP COLUMN questions;
-- +goose StatementEnd
