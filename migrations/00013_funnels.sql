-- +goose Up
-- +goose StatementBegin
create table `funnels`
(
    id         varchar(36) not null default '',
    user_id    varchar(36)  not null default '',
    name       varchar(255) not null default '',
    path       varchar(255) not null default '',
    created_at datetime     not null default CURRENT_TIMESTAMP,
    updated_at datetime     not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint api_keys_pk
        primary key (id),
    constraint idx_path
        unique key (path),
    key idx_user_path (user_id, path)
) ENGINE=InnoDB;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table `funnels`;
-- +goose StatementEnd
