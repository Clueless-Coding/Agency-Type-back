-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    game_mode VARCHAR(50) NOT NULL,
    start_time TIMESTAMP DEFAULT TO_CHAR(CURRENT_TIMESTAMP AT TIME ZONE 'MSK', 'YYYY-MM-DD HH24:MI:SS')::TIMESTAMP,
    duration time NOT NULL,
    misstakes int NOT NULL,
    accuracy FLOAT NOT NULL,
    count_words int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE results;
-- +goose StatementEnd