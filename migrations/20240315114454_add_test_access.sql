-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accesses (
    id SERIAL PRIMARY KEY,
    date_start DATE NOT NULL,
    date_end DATE NOT NULL,
    passage_time INTEGER NOT NULL,
    criteria JSON NOT NULL,
    group_id INT REFERENCES groups(id) ON DELETE SET NULL,
    test_id INT REFERENCES tests(id) ON DELETE SET NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS passes (
    id SERIAL PRIMARY KEY,
    code BIGINT NOT NULL,
    is_activated BOOLEAN DEFAULT FALSE,
    datetime_activate TIMESTAMP,
    access_id INT REFERENCES accesses(id) ON DELETE CASCADE,
    student_id INT REFERENCES students(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passes;
DROP TABLE accesses;
-- +goose StatementEnd
