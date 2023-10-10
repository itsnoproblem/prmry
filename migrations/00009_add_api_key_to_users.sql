-- +goose Up
-- +goose StatementBegin
ALTER TABLE `users`
    ADD COLUMN `api_key` VARCHAR(32) NOT NULL DEFAULT '',
    ADD INDEX  `idx_api_key` (`api_key`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `users`
    DROP INDEX `idx_api_key`,
    DROP COLUMN `api_key`;
-- +goose StatementEnd
