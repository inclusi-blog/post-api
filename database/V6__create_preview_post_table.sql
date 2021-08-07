create table PREVIEW_POSTS
(
    id bigserial not null,
    post_id BIGINT not null
        constraint preview_posts_posts_id_fk
        references posts
        on delete cascade,
    title varchar(100) not null,
    tagline varchar(100) not null,
    preview_image text,
    like_count bigint not null,
    comment_count bigint not null,
    view_time bigint not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp,
    deleted_at timestamp
);

create unique index preview_posts_id_uindex
    on PREVIEW_POSTS (id);

create unique index preview_posts_post_id_uindex
    on PREVIEW_POSTS (post_id);

alter table PREVIEW_POSTS
    add constraint preview_post_pk
        primary key (id);

