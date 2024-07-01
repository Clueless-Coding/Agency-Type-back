-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(250) NOT NULL,
    is_admin boolean DEFAULT false,
    registred_date TIMESTAMP DEFAULT TO_CHAR(CURRENT_TIMESTAMP AT TIME ZONE 'MSK', 'YYYY-MM-DD HH24:MI:SS')::TIMESTAMP,
    token VARCHAR(250)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd