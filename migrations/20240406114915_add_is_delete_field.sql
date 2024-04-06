-- +goose Up
-- +goose StatementBegin
ALTER TABLE tests
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

ALTER TABLE groups
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tests
    DROP COLUMN is_deleted;

ALTER TABLE groups
    DROP COLUMN is_deleted;
-- +goose StatementEnd
