-- +goose Up
-- +goose StatementBegin
create table permissions
(
    client_id uuid not null,
    user_id   uuid not null,
    can_use   boolean default true,
    constraint permissions_client_user_unique unique (client_id, user_id),
    constraint permissions_client_fk foreign key (client_id) references clients (id) on delete cascade,
    constraint permissions_user_fk foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table permissions
-- +goose StatementEnd
