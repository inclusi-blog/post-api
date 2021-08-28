create table likes
(
    post_id uuid not null
        constraint likes_posts_id_fk
        references posts,
    liked_by uuid not null
        constraint likes_users_id_fk
        references users
);

create unique index likes_post_id_liked_by_uindex
    on likes (post_id, liked_by);

alter table likes
    add constraint likes_pk
        primary key (post_id, liked_by);

