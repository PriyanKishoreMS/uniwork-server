-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    CONSTRAINT fcm_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd


-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS fcm_tokens_users_id_idx ON fcm_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fcm_tokens;
-- +goose StatementEnd
