-- +goose Up
-- +goose StatementBegin
CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    mark INT NOT NULL,
    score INT NOT NULL,
    max_score INT NOT NULL,
    datetime_complete TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pass_id INT UNIQUE REFERENCES passes(id) ON DELETE CASCADE,
    access_id INT REFERENCES accesses(id) ON DELETE CASCADE,
    student_id INT REFERENCES students(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results;
-- +goose StatementEnd
