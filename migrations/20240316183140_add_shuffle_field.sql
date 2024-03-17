-- +goose Up
-- +goose StatementBegin
ALTER TABLE accesses
    ADD COLUMN shuffle BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accesses
    DROP COLUMN shuffle;
-- +goose StatementEnd
