-- +goose Up
CREATE TABLE IF NOT EXISTS releases
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    date       DATETIME     NOT NULL,
    status     VARCHAR(20)  NOT NULL,
    notes      TEXT,
    duty_users JSON         NOT NULL DEFAULT (JSON_ARRAY()),
    created_at TIMESTAMP    NULL     DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_status (status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd