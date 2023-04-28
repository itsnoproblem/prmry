-- +goose Up
-- +goose StatementBegin
ALTER TABLE `interactions`
    ADD COLUMN `user_id` VARCHAR(36) NOT NULL DEFAULT '',
    ADD INDEX idx_user_id (`user_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `interactions`
    DROP INDEX `idx_user_id`,
    DROP COLUMN `user_id`;
-- +goose StatementEnd
