-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_details (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    worker_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    razorpay_payment_id VARCHAR(32),
    razorpay_order_id VARCHAR(32) NOT NULL, 
    razorpay_signature VARCHAR(128),
    status VARCHAR(32) NOT NULL DEFAULT 'pending', -- pending, funded, released, refunded
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,

    CONSTRAINT payment_task_fk FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT payment_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT payment_worker_fk FOREIGN KEY (worker_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payment_task_id ON payment_details(task_id);
CREATE INDEX IF NOT EXISTS idx_payment_user_id ON payment_details(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_worker_id ON payment_details(worker_id);

-- creating a foreign key constraint to tasks table
ALTER TABLE tasks ADD CONSTRAINT task_payment_fk FOREIGN KEY (payment_id) REFERENCES payment_details(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_details CASCADE;
-- +goose StatementEnd
