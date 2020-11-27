create table comments
(
    id bigserial not null
        constraint comments_pk
        primary key,
    post_id bigint not null
        constraint comments_posts_id_fk
        references posts,
    comments post_comment[],
    created_at timestamp default current_timestamp not null,
    updated_at timestamp,
    deleted_at timestamp
);

create unique index comments_post_id_uindex
    on comments (post_id);

