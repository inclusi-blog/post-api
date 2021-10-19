create table social_links
(
    id uuid not null,
    facebook text,
    linkedin text,
    twitter text,
    user_id uuid
        constraint social_links_users_id_fk
        references users
);

create unique index social_links_user_id_uindex
    on social_links (user_id);

