create table if not exists User
(
    id           int auto_increment
    primary key,
    username     varchar(256)         null,
    care         tinyint(1) default 1 null,
    phone_number varchar(255)         null,
    telegram_id  int                  null,
    chat_id      int                  null,
    first_name   varchar(255)         null
);