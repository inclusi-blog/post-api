create table post_views
(
    post_id    uuid
        constraint post_views_posts_id_fk
            references posts,
    user_id    uuid
        constraint post_views_users_id_fk
            references users,
    created_at timestamptz default current_timestamp
);

create unique index post_views_post_id_user_id_uindex
    on post_views (post_id, user_id);

alter table post_views
    add constraint post_views_pk
        primary key (post_id, user_id);

