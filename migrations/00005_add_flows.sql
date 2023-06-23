-- +goose Up
-- +goose StatementBegin
CREATE TABLE `flows` (
    `id` varchar(36) NOT NULL DEFAULT '',
    `user_id` varchar(36) NOT NULL DEFAULT '',
    `name` varchar(255) NOT NULL DEFAULT '',
    `rules` json NOT NULL,
    `require_all` tinyint(1) UNSIGNED NOT NULL DEFAULT 0,
    `prompt` varchar(8192) NOT NULL DEFAULT '',
    `prompt_args` json NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `flows`;
-- +goose StatementEnd
