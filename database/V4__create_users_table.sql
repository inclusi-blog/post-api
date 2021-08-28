create table users
(
    id uuid not null,
    email varchar(50) not null,
    password varchar(64) not null,
    username varchar(50) not null,
    is_active boolean not null,
    role_id uuid not null
        constraint users_roles_id_fk
            references roles,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index users_email_uindex
    on users (email);

create unique index users_id_uindex
    on users (id);

create unique index users_username_uindex
    on users (username);

alter table users
    add constraint users_pk
        primary key (id);

