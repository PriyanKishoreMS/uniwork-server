-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY UNIQUE,
    user_id BIGINT NOT NULL,
    worker_id BIGINT,
    college_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255) NOT NULL,
    price BIGINT,
    status VARCHAR(255) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expiry TIMESTAMP,
    images VARCHAR(512)[], 
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT task_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT task_clg_fk FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE
);

ALTER TABLE tasks ALTER COLUMN images SET DEFAULT array[]::varchar[];
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_college_id ON tasks(college_id);
-- +goose StatementEnd



-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
