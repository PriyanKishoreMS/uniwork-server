-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_details (
    id BIGSERIAL PRIMARY KEY,
    task_request_id BIGINT NOT NULL UNIQUE,
    task_owner_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    razorpay_payment_id VARCHAR(32),
    razorpay_order_id VARCHAR(32) NOT NULL, 
    razorpay_signature VARCHAR(128),
    status VARCHAR(32) NOT NULL DEFAULT 'pending', -- pending, funded, released, refunded
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,

    CONSTRAINT payment_task_request_fk FOREIGN KEY (task_request_id) REFERENCES task_requests(id) ON DELETE CASCADE,
    CONSTRAINT payment_task_owner_fk FOREIGN KEY (task_owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_payment_task_request_id ON payment_details(task_request_id);
CREATE INDEX IF NOT EXISTS idx_payment_task_owner_id ON payment_details(task_owner_id);

-- creating a foreign key constraint to tasks table
ALTER TABLE tasks ADD CONSTRAINT task_payment_fk FOREIGN KEY (payment_id) REFERENCES payment_details(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_details CASCADE;
-- +goose StatementEnd
