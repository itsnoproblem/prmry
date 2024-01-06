-- +goose Up
-- +goose StatementBegin
create table api_keys
(
    user_id    varchar(36)  default ''                not null,
    name       varchar(255) default ''                not null,
    value      varchar(32)  default ''                not null,
    created_at datetime     default CURRENT_TIMESTAMP not null,
    constraint api_keys_pk
        primary key (value)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create index idx_user_id
    on api_keys (user_id);

insert into api_keys (user_id, name, value)
select id, 'untitled', api_key from users;

alter table users
    drop column api_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table api_keys;
alter table users add column api_key varchar(32) default '' not null;
-- +goose StatementEnd
