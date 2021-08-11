create table roles
(
    id uuid not null,
    name varchar(10) not null,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index roles_id_uindex
    on roles (id);

create unique index roles_name_uindex
    on roles (name);

alter table roles
    add constraint roles_pk
        primary key (id);

