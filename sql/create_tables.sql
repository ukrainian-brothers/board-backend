create table users
(
    login        varchar(36),
    id           varchar(36) not null
        constraint users_pk
            primary key,
    password     varchar(97),
    name         varchar(15),
    surname      varchar(15),
    mail         varchar(45),
    phone_number varchar(15)
);

alter table users
    owner to postgres;

create table adverts
(
    id              varchar(36),
    user_id         varchar(36)
        constraint user___fk
            references users,
    created_at      timestamp default now(),
    updated_at      timestamp,
    destroyed_at    timestamp,
    title           varchar(50),
    description     varchar(250),
    type            varchar(15),
    views           integer,
    contact_details json
);

alter table adverts
    owner to postgres;

create unique index adverts_id_uindex
    on adverts (id);

create unique index users_id_uindex
    on users (id);

create unique index users_login_uindex
    on users (login);

