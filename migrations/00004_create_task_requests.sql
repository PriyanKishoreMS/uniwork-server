-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_requests (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    task_worker_id BIGINT NOT NULL,
    -- status being pending, accepted, rejected
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT task_req_task_fk FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT task_req_task_worker_fk FOREIGN KEY (task_worker_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE task_requests ADD CONSTRAINT task_requests_unique UNIQUE (task_id, task_worker_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_requests;
-- +goose StatementEnd
