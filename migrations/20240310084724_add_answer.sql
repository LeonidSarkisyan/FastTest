-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    text TEXT,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    question_id INT REFERENCES questions(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE answers;
-- +goose StatementEnd
