create table clients
(
    id          uuid primary key      default gen_random_uuid(),
    class       varchar(50)  not null default 'client',
    name        varchar(50)  not null,
    description varchar,
    logo        varchar,
    image       varchar,
    secret      varchar      not null,
    callback    varchar      not null,
    can_use     boolean      not null default true,
    created_at  timestamp(6) not null default now(),
    updated_at  timestamp(6) not null default now()
);

insert into clients(class, name, logo, image, secret, callback, can_use)
values ('manager', 'SSO Manager', '', 'https://alnovi.ru/drive/webman/3rdparty/SynologyDrive-Drive/images/_Asset/2x/_Drive/wallpaper_drive.jpg', substr(md5(random()::text), 0, 25), 'https://ya.ru', false),
       ('profile', 'SSO Profile', '', 'https://alnovi.ru/drive/webman/3rdparty/SynologyDrive-Drive/images/_Asset/2x/_Drive/wallpaper_drive.jpg', substr(md5(random()::text), 0, 25), 'https://ya.ru', true);