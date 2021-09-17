create table user_interests
(
    user_id uuid not null
        constraint user_interests_users_id_fk
        references users,
    interest_id uuid not null
        constraint user_interests_interests_id_fk
        references interests
);

create unique index user_interests_user_id_interest_id_uindex
    on user_interests (user_id, interest_id);

alter table user_interests
    add constraint user_interests_pk
        unique (user_id, interest_id);

