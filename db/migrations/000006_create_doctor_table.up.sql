create table if not exists Doctor
(
    id                int auto_increment
    primary key,
    name              varchar(255) null,
    phone_number      varchar(512) not null,
    telegram_id       int          null,
    telegram_username varchar(255) not null,
    whats_app_url     varchar(512) null,
    constraint Doctor_id_uindex
    unique (id)
);

