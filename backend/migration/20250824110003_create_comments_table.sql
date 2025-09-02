-- +goose Up
CREATE TABLE IF NOT EXISTS comments
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    release_id INT UNSIGNED NOT NULL,
    author     VARCHAR(255) NOT NULL,
    message    TEXT         NOT NULL,
    created_at TIMESTAMP    NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_comments_release FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd