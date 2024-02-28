-- +goose Up
-- +goose StatementBegin
create table funnel_flows
(
    funnel_id   varchar(36)  not null default '',
    flow_id    varchar(36)  not null default '',
    created_at datetime     not null default CURRENT_TIMESTAMP,
    constraint funnel_flows_pk
        primary key (funnel_id, flow_id),
    key `idx_funnel_id` (funnel_id)
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE funnel_flows;
-- +goose StatementEnd
