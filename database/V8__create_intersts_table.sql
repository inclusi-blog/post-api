create table interests
(
    id uuid not null
        constraint interests_pk
            primary key,
    name varchar(50) not null,
    cover_pic text,
    added_by uuid,
    is_blocked boolean default false not null,
    approved_by uuid not null
        constraint interests_admin_id_fk
            references admin,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index interests_id_uindex
    on interests (id);

create unique index interests_name_uindex
    on interests (name);


