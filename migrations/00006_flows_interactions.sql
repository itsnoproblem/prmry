-- +goose Up
-- +goose StatementBegin
ALTER TABLE `interactions`
ADD COLUMN `flow_id` VARCHAR(36) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `interactions`
DROP COLUMN `flow_id`;
-- +goose StatementEnd
