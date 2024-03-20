-- +goose Up
-- +goose StatementBegin
ALTER TABLE results
    ADD COLUMN time_pass INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE results
    DROP COLUMN time_pass;
-- +goose StatementEnd
