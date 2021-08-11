create table interests
(
    id uuid not null,
    name varchar(50) not null,
    cover_pic text,
    added_by uuid not null
        constraint interests_users_id_fk
        references users,
    is_blocked boolean default false not null,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index interests_id_uindex
    on interests (id);

create unique index interests_name_uindex
    on interests (name);

alter table interests
    add constraint interests_pk
        primary key (id);

