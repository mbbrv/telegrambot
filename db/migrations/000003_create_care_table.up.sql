create table if not exists Care
(
    id               int auto_increment
    primary key,
    user_id          int          null,
    description      text         null,
    url              varchar(255) null,
    photo_dictionary varchar(255) null,
    time             time         null,
    day_time         varchar(255) null,
    constraint Care_id_uindex
    unique (id),
    constraint Care_User_id_fk
    foreign key (user_id) references User (id)
);

