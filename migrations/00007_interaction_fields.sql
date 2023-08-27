-- +goose Up
-- +goose StatementBegin
ALTER TABLE `interactions`
    DROP COLUMN `err`,
    ADD COLUMN `type` VARCHAR(256) NOT NULL DEFAULT '',
    ADD COLUMN `model` VARCHAR(256) NOT NULL DEFAULT '',
    ADD COLUMN `prompt` BLOB NOT NULL,
    ADD COLUMN `completion` BLOB NOT NULL,
    ADD COLUMN `prompt_tokens` INT(11) UNSIGNED NOT NULL DEFAULT 0,
    ADD COLUMN `completion_tokens` INT(11) UNSIGNED NOT NULL DEFAULT 0;

UPDATE `interactions`
SET
    `type` = IFNULL(`response`->>'$.object', ''),
    `model` = IFNULL(`request`->>'$.model', ''),
    `prompt` = IFNULL(`request`->>'$.messages[0].content', IFNULL(`request`->>'$.prompt', '')),
    `completion` = IFNULL(`response`->>'$.choices[0].message.content', IFNULL(`response`->>'$.choices[0].text', '')),
    `prompt_tokens` = IFNULL(`response`->>'$.usage.prompt_tokens', 0),
    `completion_tokens` = IFNULL(`response`->>'$.usage.completion_tokens', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `interactions`
    ADD COLUMN `err` VARCHAR(2048) NOT NULL,
    DROP COLUMN `type`,
    DROP COLUMN `model`,
    DROP COLUMN `prompt`,
    DROP COLUMN `completion`,
    DROP COLUMN `prompt_tokens`,
    DROP COLUMN `completion_tokens`;
-- +goose StatementEnd
