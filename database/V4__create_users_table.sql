create table users
(
    id uuid not null,
    email varchar(50),
    role_id uuid not null
);

create unique index users_email_uindex
    on users (email);

create unique index users_id_uindex
    on users (id);

alter table users
    add constraint users_pk
        primary key (id);

alter table users
    add constraint users_roles_id_fk
        foreign key (role_id) references roles;


