create table followings
(
    follower_id uuid not null
        constraint followings_users_id_fk
        references users,
    following_id uuid not null
        constraint followings_users_id_fk_2
        references users,
    created_at timestamptz default current_timestamp not null,
    constraint followings_pk
        primary key (following_id, follower_id)
);

create unique index followings_follower_id_following_id_uindex
    on followings (follower_id, following_id);

