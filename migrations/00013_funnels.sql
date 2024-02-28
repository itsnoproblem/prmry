-- +goose Up
-- +goose StatementBegin
create table funnels
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table funnel_flows
(
    funnel_id   varchar(36)  not null default '',
    flow_id    varchar(36)  not null default '',
    created_at datetime     not null default CURRENT_TIMESTAMP,
    constraint funnel_flows_pk
        primary key (funnel_id, flow_id),
    key `idx_funnel_id` (funnel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table `funnels`;
drop table `funnel_flows`;
-- +goose StatementEnd
