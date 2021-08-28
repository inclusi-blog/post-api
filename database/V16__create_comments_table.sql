create table comments
(
    id uuid not null,
    data text not null,
    post_id uuid not null
        constraint comments_posts_id_fk
        references posts,
    commented_by uuid not null
        constraint comments_users_id_fk
        references users,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index comments_id_uindex
    on comments (id);

alter table comments
    add constraint comments_pk
        primary key (id);

