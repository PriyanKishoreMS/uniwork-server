-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    college_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    mobile VARCHAR(15) NOT NULL UNIQUE,
    avatar VARCHAR(512) NOT NULL DEFAULT 'default',
    dept VARCHAR(255) NOT NULL,
    tasks_completed INT NOT NULL DEFAULT 0, 
    earned BIGINT NOT NULL DEFAULT 0,
    rating DECIMAL NOT NULL DEFAULT 0,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT user_clg_fk FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_users_college_id ON users(college_id);
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO users (college_id, name, email, dept, mobile)
VALUES
(1, 'Priyan Kishore', '20113022@student.hindustanuniv.ac.in', 'Computer Science', '7010376477'),
(1, 'Chandana Sathwika', '20113024@student.hindustanuniv.ac.in', 'Computer Science', '7674017177'),
(6, 'Laksitha Bharani', 'laksitha2004@gmail.com', 'MBBS', '6380886960'),
(3, 'Test User', 'testuser@college.ac.in', 'Test Dept', '1234567890');
-- +goose StatementEnd




-- +goose Down
-- -- +goose StatementBegin
-- ALTER TABLE users DROP FOREIGN KEY user_clg_fk;
-- -- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
