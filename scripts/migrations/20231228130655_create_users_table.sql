-- +goose Up
-- +goose StatementBegin
create table users
(
    id         uuid primary key      default gen_random_uuid(),
    image      varchar,
    name       varchar      not null,
    login      varchar(100) not null,
    email      varchar(100) not null,
    password   varchar      not null,
    created_at timestamp(6) not null default now(),
    updated_at timestamp(6) not null default now()
);

insert into users (name, login, email, password)
values ('Admin', 'admin', 'admin@example.com', 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users
-- +goose StatementEnd
