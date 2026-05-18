-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feedbacks (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255)   NOT NULL,
    email            VARCHAR(255)   NOT NULL,
    subject            VARCHAR(255)   NOT NULL,
    message     TEXT           NOT NULL DEFAULT '', 
    created_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedbacks;
-- +goose StatementEnd