-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    datetime_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    datetime_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    criteria JSON,
    user_id INT references users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_datetime_update()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.datetime_update = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_datetime_update
    BEFORE UPDATE ON tests
    FOR EACH ROW
EXECUTE FUNCTION update_datetime_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_datetime_update ON tests;
DROP FUNCTION IF EXISTS update_datetime_update();
DROP TABLE IF EXISTS tests;
-- +goose StatementEnd
