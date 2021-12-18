create table user_blocks
(
    blocked_id uuid                                  not null
        constraint user_blocks_users_id_fk
            references users,
    blocked_by uuid                                  not null
        constraint user_blocks_users_id_fk_2
            references users,
    created_at timestamptz default current_timestamp not null,
    constraint user_blocks_pk
        primary key (blocked_id, blocked_by)
);

create unique index user_blocks_blocked_id_blocked_by_uindex
    on user_blocks (blocked_id, blocked_by);

