-- +goose Up
-- +goose StatementBegin
ALTER TABLE questions
    ADD COLUMN data JSON,
    ADD COLUMN type VARCHAR(255) DEFAULT 'choose';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE questions
    DROP COLUMN data,
    DROP COLUMN type;
-- +goose StatementEnd
