create table read_later
(
    post_id uuid
        constraint read_later_posts_id_fk
        references posts,
    user_id uuid
        constraint read_later_users_id_fk
        references users,
    constraint read_later_pk
        primary key (user_id, post_id)
);

create unique index read_later_post_id_user_id_uindex
    on read_later (post_id, user_id);

