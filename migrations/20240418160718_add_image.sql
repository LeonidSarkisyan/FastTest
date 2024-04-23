-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    url TEXT,
    question_id INT REFERENCES questions(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
