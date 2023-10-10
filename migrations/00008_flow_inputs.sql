-- +goose Up
-- +goose StatementBegin
ALTER TABLE `flows`
ADD COLUMN `inputs` JSON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `flows`
DROP COLUMN `inputs`;
-- +goose StatementEnd
