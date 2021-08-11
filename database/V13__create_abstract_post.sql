create table abstract_post
(
    id uuid not null,
    title varchar(300) not null,
    tagline varchar(100) not null,
    preview_image text,
    view_time bigint not null,
    post_id uuid not null
        constraint abstract_post_posts_id_fk
        references posts,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index abstract_post_id_uindex
    on abstract_post (id);

create unique index abstract_post_post_id_uindex
    on abstract_post (post_id);

alter table abstract_post
    add constraint abstract_post_pk
        primary key (id);

