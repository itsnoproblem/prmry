-- +goose Up
-- +goose StatementBegin
ALTER TABLE `flows`
    ADD COLUMN `model` VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN `temperature` DECIMAL(1,1) NOT NULL DEFAULT 0.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `flows`
    DROP COLUMN `model`,
    DROP COLUMN `temperature`;
-- +goose StatementEnd
