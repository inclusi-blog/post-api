create table drafts
(
    id uuid not null,
    data jsonb,
    preview_image text,
    tagline varchar(100),
    user_id uuid
        constraint drafts_users_id_fk
            references users,
    interests varchar(50) ARRAY[5],
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index drafts_id_uindex
    on drafts (id);

alter table drafts
    add constraint drafts_pk
        primary key (id);

