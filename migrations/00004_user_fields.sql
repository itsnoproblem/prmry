-- +goose Up
-- +goose StatementBegin
ALTER TABLE `users`
    ADD COLUMN `nickname` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    ADD COLUMN `avatar_url` varchar(4096) NOT NULL,
    ADD INDEX `idx_email` (`email`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `users`
    DROP COLUMN `nickname`,
    DROP COLUMN `avatar_url`,
    DROP INDEX `idx_email`;
-- +goose StatementEnd
