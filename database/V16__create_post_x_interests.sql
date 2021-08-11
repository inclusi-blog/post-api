create table post_x_interests
(
    post_id uuid not null
        constraint post_x_interests_posts_id_fk
        references posts,
    interest_id uuid not null
        constraint post_x_interests_interests_id_fk
        references interests
);

create unique index post_x_interests_post_id_interest_id_uindex
    on post_x_interests (post_id, interest_id);

