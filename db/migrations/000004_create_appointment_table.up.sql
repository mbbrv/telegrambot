create table if not exists Appointment
(
    id          int auto_increment
    primary key,
    user_id     int          not null,
    doctor_id   int          not null,
    datetime    datetime     null,
    place       varchar(512) null,
    description text         null,
    cost        int          null,
    constraint Appointment_Doctor_id_fk
    foreign key (doctor_id) references Doctor (id),
    constraint FK_UserId
    foreign key (user_id) references User (id)
);

create index Appointment_doctor_id_index
    on Appointment (doctor_id);

create index user_id
    on Appointment (user_id);

