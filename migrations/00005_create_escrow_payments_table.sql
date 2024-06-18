-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS escrow_payments (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    worker_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending', -- pending, funded, released, refunded
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,

    CONSTRAINT escrow_task_fk FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT escrow_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT escrow_worker_fk FOREIGN KEY (worker_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_escrow_task_id ON escrow_payments(task_id);
CREATE INDEX IF NOT EXISTS idx_escrow_user_id ON escrow_payments(user_id);
CREATE INDEX IF NOT EXISTS idx_escrow_worker_id ON escrow_payments(worker_id);

-- creating a foreign key constraint to tasks table
ALTER TABLE tasks ADD CONSTRAINT task_escrow_fk FOREIGN KEY (escrow_payment_id) REFERENCES escrow_payments(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS escrow_payments CASCADE;
-- +goose StatementEnd
