create table tokens
(
    id         uuid primary key      default gen_random_uuid(),
    class      varchar(30)  not null,
    hash       varchar(255) not null,
    user_id    uuid,
    client_id  uuid,
    meta       jsonb,
    not_before timestamp(6) not null default now(),
    expiration timestamp(6) not null default now(),
    created_at timestamp(6) not null default now(),
    updated_at timestamp(6) not null default now(),
    constraint tokens_class_hash_unique unique (class, hash),
    constraint tokens_user_fk foreign key (user_id) references users (id) on delete cascade,
    constraint tokens_client_fk foreign key (client_id) references clients (id) on delete cascade
);