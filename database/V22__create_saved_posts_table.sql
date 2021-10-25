create table saved_posts
(
    post_id uuid
        constraint saved_posts_posts_id_fk
        references posts,
    user_id uuid
        constraint saved_posts_users_id_fk
        references users,
    constraint saved_posts_pk
        primary key (user_id, post_id)
);

create unique index saved_posts_post_id_user_id_uindex
    on saved_posts (post_id, user_id);

